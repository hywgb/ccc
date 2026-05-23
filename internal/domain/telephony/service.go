package telephony

import (
	"context"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// RoutingService matches outbound calls to SIP trunks via RoutingRules.
type RoutingService struct {
	rules RoutingRuleRepository
}

func NewRoutingService(rules RoutingRuleRepository) *RoutingService {
	return &RoutingService{rules: rules}
}

// MatchRule finds the highest-priority active rule matching the callee number and current time.
func (s *RoutingService) MatchRule(ctx context.Context, tenantID int64, callee string) (*RoutingRule, error) {
	rules, err := s.rules.ListActive(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var best *RoutingRule

	for _, r := range rules {
		if !matchesRule(r, callee, now) {
			continue
		}
		if best == nil || r.Priority < best.Priority {
			best = r
		}
	}

	if best == nil {
		return nil, ErrNoMatchingRoute
	}
	return best, nil
}

func matchesRule(r *RoutingRule, callee string, now time.Time) bool {
	switch r.MatchType {
	case "prefix":
		return strings.HasPrefix(callee, r.MatchValue)
	case "time_of_day":
		return matchTimeOfDay(r.MatchValue, now)
	case "any":
		return true
	default:
		return false
	}
}

// matchTimeOfDay checks if now falls within "HH:MM-HH:MM" range.
func matchTimeOfDay(value string, now time.Time) bool {
	parts := strings.SplitN(value, "-", 2)
	if len(parts) != 2 {
		return false
	}

	start, err1 := time.Parse("15:04", strings.TrimSpace(parts[0]))
	end, err2 := time.Parse("15:04", strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return false
	}

	current := time.Date(0, 1, 1, now.Hour(), now.Minute(), 0, 0, time.UTC)
	s := time.Date(0, 1, 1, start.Hour(), start.Minute(), 0, 0, time.UTC)
	e := time.Date(0, 1, 1, end.Hour(), end.Minute(), 0, 0, time.UTC)

	if s.Before(e) {
		return !current.Before(s) && current.Before(e)
	}
	// overnight range (e.g. 22:00-06:00)
	return !current.Before(s) || current.Before(e)
}

// CLIPolicyService selects outbound caller ID based on policy.
type CLIPolicyService struct {
	policies CLIPolicyRepository
	phones   PhoneNumberRepository
	counter  uint64
}

func NewCLIPolicyService(policies CLIPolicyRepository, phones PhoneNumberRepository) *CLIPolicyService {
	return &CLIPolicyService{policies: policies, phones: phones}
}

// SelectCLI picks a caller ID number based on the given policy (or tenant default).
func (s *CLIPolicyService) SelectCLI(ctx context.Context, tenantID int64, policyID *int64, callee string) (*PhoneNumber, error) {
	var policy *CLIPolicy
	var err error

	if policyID != nil {
		policy, err = s.policies.GetByID(ctx, *policyID)
	} else {
		policy, err = s.policies.GetDefault(ctx, tenantID)
	}
	if err != nil || policy == nil {
		return nil, ErrNoCLIPolicy
	}

	switch policy.Strategy {
	case CLIStrategyFixed:
		if policy.FixedNumberID == nil {
			return nil, ErrNoCLINumber
		}
		return s.phones.GetByID(ctx, *policy.FixedNumberID)

	case CLIStrategyRoundRobin:
		pool := parseNumberPool(policy.NumberPoolIDs)
		if len(pool) == 0 {
			return nil, ErrNoCLINumber
		}
		idx := atomic.AddUint64(&s.counter, 1) % uint64(len(pool))
		return s.phones.GetByID(ctx, pool[idx])

	case CLIStrategyRandom:
		pool := parseNumberPool(policy.NumberPoolIDs)
		if len(pool) == 0 {
			return nil, ErrNoCLINumber
		}
		idx := rand.Intn(len(pool))
		return s.phones.GetByID(ctx, pool[idx])

	case CLIStrategyMatchArea:
		pool := parseNumberPool(policy.NumberPoolIDs)
		if len(pool) == 0 {
			return nil, ErrNoCLINumber
		}
		// Try to match area code prefix
		for _, pid := range pool {
			pn, err := s.phones.GetByID(ctx, pid)
			if err != nil || pn == nil {
				continue
			}
			if len(callee) >= 4 && len(pn.Number) >= 4 && pn.Number[:4] == callee[:4] {
				return pn, nil
			}
		}
		// Fallback to first number
		return s.phones.GetByID(ctx, pool[0])

	default:
		return nil, ErrNoCLIPolicy
	}
}

// TrunkHealthService manages trunk health and failover.
type TrunkHealthService struct {
	trunks SIPTrunkRepository
	groups SIPTrunkGroupRepository
	mu     sync.RWMutex
	health map[int64]*TrunkHealthStatus
}

func NewTrunkHealthService(trunks SIPTrunkRepository, groups SIPTrunkGroupRepository) *TrunkHealthService {
	return &TrunkHealthService{
		trunks: trunks,
		groups: groups,
		health: make(map[int64]*TrunkHealthStatus),
	}
}

// RecordHealthCheck updates trunk health after an OPTIONS keepalive result.
func (s *TrunkHealthService) RecordHealthCheck(trunkID int64, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	h, ok := s.health[trunkID]
	if !ok {
		h = &TrunkHealthStatus{TrunkID: trunkID, Status: TrunkStatusActive}
		s.health[trunkID] = h
	}
	h.LastCheckAt = time.Now()

	if success {
		h.Status = TrunkStatusActive
		h.FailCount = 0
	} else {
		h.FailCount++
		if h.FailCount >= 3 {
			h.Status = TrunkStatusDown
		}
	}
}

// GetHealthStatus returns a snapshot of the current health status of a trunk.
func (s *TrunkHealthService) GetHealthStatus(trunkID int64) *TrunkHealthStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.health[trunkID]
	if !ok {
		return &TrunkHealthStatus{TrunkID: trunkID, Status: TrunkStatusActive}
	}
	cp := *h
	return &cp
}

// SelectHealthyTrunk picks the best available trunk from a group, skipping down trunks.
func (s *TrunkHealthService) SelectHealthyTrunk(ctx context.Context, groupID int64) (*SIPTrunk, error) {
	members, err := s.groups.ListMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, ErrTrunkGroupNotFound
	}

	// Sort by priority (lower = higher priority) — members are already in insert order
	for _, m := range members {
		h := s.GetHealthStatus(m.SIPTrunkID)
		if h.Status == TrunkStatusDown {
			continue
		}
		trunk, err := s.trunks.GetByID(ctx, m.SIPTrunkID)
		if err != nil || trunk == nil {
			continue
		}
		return trunk, nil
	}
	return nil, ErrNoHealthyTrunk
}

func parseNumberPool(ids string) []int64 {
	if ids == "" {
		return nil
	}
	parts := strings.Split(ids, ",")
	var result []int64
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		var id int64
		valid := true
		for _, c := range p {
			if c < '0' || c > '9' {
				valid = false
				break
			}
			id = id*10 + int64(c-'0')
		}
		if valid {
			result = append(result, id)
		}
	}
	return result
}
