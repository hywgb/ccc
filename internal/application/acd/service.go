// Package acd implements the Automatic Call Distribution dispatcher.
//
// Background: prior to this package, calls were routed to FreeSWITCH's
// mod_callcenter (`callcenter:{skill_group_id}@default`) but no
// callcenter.conf.xml was deployed, so calls entered the queue and were never
// distributed. This package replaces that integration with a server-side
// dispatcher that:
//
//  1. Accepts Enqueue requests when a call leaves IVR (skill group, priority).
//  2. Stores queued calls in a Redis sorted set per skill group.
//  3. Polls each active skill group on a ticker; for each head-of-queue call,
//     selects an idle agent according to the configured routing policy and
//     transitions the call to ringing via lifecycle.Service.
//
// The dispatcher runs as a single goroutine. For multi-instance deployments,
// callers can scope the loop to a subset of skill groups; the Redis state
// itself is shared so any instance can drain.
package acd

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/divord97/ccc/internal/application/lifecycle"
	"github.com/divord97/ccc/internal/domain/call"
	"github.com/divord97/ccc/internal/domain/identity"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	queueKeyPrefix    = "acd:queue:"       // ZSET: score = priority-adjusted enqueue time
	activeSGKey       = "acd:active_sg"    // SET of skill_group_ids that have ever been used
	agentClaimPrefix  = "acd:agent_claim:" // SETNX agent_id during dispatch to avoid double-route
	roundRobinPrefix  = "acd:rr_cursor:"   // INCR cursor per skill group
	defaultPollPeriod = 500 * time.Millisecond
	agentClaimTTL     = 30 * time.Second
)

// LifecycleService is the subset of lifecycle.Service required by the dispatcher.
type LifecycleService interface {
	TransitionCallToRinging(ctx context.Context, callID, agentUserID int64) (*call.Call, error)
}

var _ LifecycleService = (*lifecycle.Service)(nil)

// PresenceRepo exposes the lookups the dispatcher needs.
type PresenceRepo interface {
	GetByAgentID(ctx context.Context, agentID int64) (*identity.AgentPresence, error)
}

// MembersRepo lists agents in a skill group.
type MembersRepo interface {
	GetBySkillGroup(ctx context.Context, skillGroupID int64) ([]*identity.SkillGroupMember, error)
}

// SkillGroups resolves the routing policy and tenant for a skill group.
type SkillGroups interface {
	GetByID(ctx context.Context, id int64) (*identity.SkillGroup, error)
}

// Service is the ACD dispatcher.
type Service struct {
	rdb        *redis.Client
	lifecycle  LifecycleService
	presence   PresenceRepo
	members    MembersRepo
	skillGroup SkillGroups
	logger     zerolog.Logger
	pollPeriod time.Duration
	rng        *rand.Rand
}

// Config groups the dependencies for NewService.
type Config struct {
	Redis      *redis.Client
	Lifecycle  LifecycleService
	Presence   PresenceRepo
	Members    MembersRepo
	SkillGroup SkillGroups
	Logger     zerolog.Logger
	PollPeriod time.Duration
}

// NewService wires the ACD dispatcher. The returned service is inert until Run is called.
func NewService(cfg Config) *Service {
	pp := cfg.PollPeriod
	if pp <= 0 {
		pp = defaultPollPeriod
	}
	return &Service{
		rdb:        cfg.Redis,
		lifecycle:  cfg.Lifecycle,
		presence:   cfg.Presence,
		members:    cfg.Members,
		skillGroup: cfg.SkillGroup,
		logger:     cfg.Logger,
		pollPeriod: pp,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Enqueue appends a call to the skill group queue. Higher priority values are
// served first; equal priorities are FIFO.
func (s *Service) Enqueue(ctx context.Context, callID, skillGroupID int64, priority int) error {
	if s.rdb == nil {
		return errors.New("acd: redis client not configured")
	}
	score := scoreFor(priority, time.Now())
	if err := s.rdb.ZAdd(ctx, queueKey(skillGroupID), redis.Z{Score: score, Member: strconv.FormatInt(callID, 10)}).Err(); err != nil {
		return fmt.Errorf("acd: enqueue zadd: %w", err)
	}
	if err := s.rdb.SAdd(ctx, activeSGKey, skillGroupID).Err(); err != nil {
		return fmt.Errorf("acd: register sg: %w", err)
	}
	s.logger.Debug().Int64("call_id", callID).Int64("sg", skillGroupID).Int("priority", priority).Msg("acd: enqueued")
	return nil
}

// Run drives the dispatcher loop until ctx is canceled.
func (s *Service) Run(ctx context.Context) {
	if s.rdb == nil {
		s.logger.Warn().Msg("acd: redis not configured, dispatcher disabled")
		return
	}
	t := time.NewTicker(s.pollPeriod)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.tick(ctx)
		}
	}
}

func (s *Service) tick(ctx context.Context) {
	sgIDs, err := s.rdb.SMembers(ctx, activeSGKey).Result()
	if err != nil {
		s.logger.Warn().Err(err).Msg("acd: list skill groups")
		return
	}
	for _, raw := range sgIDs {
		sgID, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			continue
		}
		s.dispatchOne(ctx, sgID)
	}
}

// dispatchOne attempts to assign the head-of-queue call for a skill group to an
// idle agent. At most one assignment per tick per skill group to keep the loop
// fair across groups.
func (s *Service) dispatchOne(ctx context.Context, sgID int64) {
	head, err := s.rdb.ZRangeWithScores(ctx, queueKey(sgID), 0, 0).Result()
	if err != nil || len(head) == 0 {
		return
	}
	callIDStr, _ := head[0].Member.(string)
	callID, err := strconv.ParseInt(callIDStr, 10, 64)
	if err != nil {
		_ = s.rdb.ZRem(ctx, queueKey(sgID), head[0].Member).Err()
		return
	}

	sg, err := s.skillGroup.GetByID(ctx, sgID)
	if err != nil || sg == nil {
		return
	}

	agentID, err := s.pickAgent(ctx, sg)
	if err != nil || agentID == 0 {
		return
	}

	if !s.tryClaim(ctx, agentID) {
		return
	}

	removed, err := s.rdb.ZRem(ctx, queueKey(sgID), callIDStr).Result()
	if err != nil || removed == 0 {
		s.releaseClaim(ctx, agentID)
		return
	}

	if _, err := s.lifecycle.TransitionCallToRinging(ctx, callID, agentID); err != nil {
		s.logger.Warn().Err(err).Int64("call_id", callID).Int64("agent_id", agentID).Msg("acd: transition to ringing failed")
		s.releaseClaim(ctx, agentID)
		// best-effort requeue with original priority so a transient failure
		// does not demote a high-priority call. We keep the original score
		// (priority window + original enqueue timestamp) which preserves both
		// ordering and FIFO-within-priority semantics.
		_ = s.rdb.ZAdd(ctx, queueKey(sgID), redis.Z{Score: head[0].Score, Member: callIDStr}).Err()
		return
	}
	s.logger.Info().Int64("call_id", callID).Int64("agent_id", agentID).Int64("sg", sgID).Msg("acd: routed call to agent")
}

func (s *Service) pickAgent(ctx context.Context, sg *identity.SkillGroup) (int64, error) {
	members, err := s.members.GetBySkillGroup(ctx, sg.ID)
	if err != nil {
		return 0, err
	}
	type idleAgent struct {
		ID       int64
		LastIdle time.Time
	}
	var candidates []idleAgent
	for _, m := range members {
		p, err := s.presence.GetByAgentID(ctx, m.AgentID)
		if err != nil || p == nil {
			continue
		}
		if p.Status != identity.PresenceIdle {
			continue
		}
		candidates = append(candidates, idleAgent{ID: m.AgentID, LastIdle: p.LastStatusAt})
	}
	if len(candidates) == 0 {
		return 0, nil
	}

	switch sg.RoutingPolicy {
	case identity.RoutingPolicyRandom:
		return candidates[s.rng.Intn(len(candidates))].ID, nil
	case identity.RoutingPolicyRoundRobin:
		idx, err := s.rdb.Incr(ctx, roundRobinPrefix+strconv.FormatInt(sg.ID, 10)).Result()
		if err != nil {
			return candidates[0].ID, nil
		}
		return candidates[int((idx-1)%int64(len(candidates)))].ID, nil
	default:
		// longest-idle (least_recent / skill_weight / familiar all fall back to longest-idle).
		best := candidates[0]
		for _, c := range candidates[1:] {
			if c.LastIdle.Before(best.LastIdle) {
				best = c
			}
		}
		return best.ID, nil
	}
}

func (s *Service) tryClaim(ctx context.Context, agentID int64) bool {
	ok, err := s.rdb.SetNX(ctx, agentClaimPrefix+strconv.FormatInt(agentID, 10), 1, agentClaimTTL).Result()
	return err == nil && ok
}

func (s *Service) releaseClaim(ctx context.Context, agentID int64) {
	_ = s.rdb.Del(ctx, agentClaimPrefix+strconv.FormatInt(agentID, 10)).Err()
}

func queueKey(sgID int64) string { return queueKeyPrefix + strconv.FormatInt(sgID, 10) }

// scoreFor encodes priority + timestamp into a single ZSet score so higher
// priority always sorts before older entries with lower priority.
func scoreFor(priority int, ts time.Time) float64 {
	// Priority window: -1e6 per priority point, then +seconds since epoch.
	return float64(-priority)*1e10 + float64(ts.UnixMilli())
}
