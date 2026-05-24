package trunk

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/divord97/ccc/internal/domain/telephony"
	"github.com/divord97/ccc/internal/infrastructure/esl"
	"github.com/rs/zerolog"
)

// HealthMonitor periodically checks SIP trunk status via FreeSWITCH ESL and manages failover.
type HealthMonitor struct {
	trunks    telephony.SIPTrunkRepository
	healthSvc *telephony.TrunkHealthService
	logger    zerolog.Logger
	esl       *esl.Client
	interval  time.Duration
	stopCh    chan struct{}
}

func NewHealthMonitor(
	trunks telephony.SIPTrunkRepository,
	healthSvc *telephony.TrunkHealthService,
	logger zerolog.Logger,
	eslClient *esl.Client,
) *HealthMonitor {
	return &HealthMonitor{
		trunks:    trunks,
		healthSvc: healthSvc,
		logger:    logger,
		esl:       eslClient,
		interval:  30 * time.Second,
		stopCh:    make(chan struct{}),
	}
}

// Start begins periodic health checks for all trunks.
func (m *HealthMonitor) Start(ctx context.Context, tenantID int64) {
	go m.loop(ctx, tenantID)
	m.logger.Info().Int64("tenant_id", tenantID).Msg("trunk health monitor started")
}

// Stop terminates the health check loop.
func (m *HealthMonitor) Stop() {
	close(m.stopCh)
}

func (m *HealthMonitor) loop(ctx context.Context, tenantID int64) {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopCh:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.checkAll(ctx, tenantID)
		}
	}
}

func (m *HealthMonitor) checkAll(ctx context.Context, tenantID int64) {
	trunks, _, err := m.trunks.List(ctx, tenantID, 0, 100)
	if err != nil {
		m.logger.Error().Err(err).Msg("trunk health: failed to list trunks")
		return
	}

	for _, t := range trunks {
		success := m.sendOPTIONS(t)
		m.healthSvc.RecordHealthCheck(t.ID, success)

		h := m.healthSvc.GetHealthStatus(t.ID)
		if h.Status == telephony.TrunkStatusDown {
			m.logger.Warn().Int64("trunk_id", t.ID).Str("name", t.Name).Msg("trunk marked DOWN")
		}
	}
}

// sendOPTIONS checks SIP trunk status via FreeSWITCH ESL sofia commands.
func (m *HealthMonitor) sendOPTIONS(t *telephony.SIPTrunk) bool {
	if m.esl == nil {
		return true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	gatewayName := fmt.Sprintf("trunk_%d", t.ID)
	resp, err := m.esl.SendCommand(ctx, fmt.Sprintf("sofia status gateway %s", gatewayName))
	if err != nil {
		m.logger.Debug().Int64("trunk_id", t.ID).Err(err).Msg("trunk ESL status check failed")
		return false
	}
	return strings.Contains(resp, "REGED") || strings.Contains(resp, "NOREG")
}

// SelectTrunk picks a healthy trunk from a group with automatic failover.
func (m *HealthMonitor) SelectTrunk(ctx context.Context, groupID int64) (*telephony.SIPTrunk, error) {
	return m.healthSvc.SelectHealthyTrunk(ctx, groupID)
}
