package call

import (
	"context"
	"sync"
)

// MockWebRTCQualityRepo is an in-memory mock for WebRTCQualityRepository.
type MockWebRTCQualityRepo struct {
	mu   sync.RWMutex
	logs map[int64]*WebRTCQualityLog
	seq  int64
}

func NewMockWebRTCQualityRepo() *MockWebRTCQualityRepo {
	return &MockWebRTCQualityRepo{logs: make(map[int64]*WebRTCQualityLog)}
}

func (m *MockWebRTCQualityRepo) Create(_ context.Context, l *WebRTCQualityLog) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	l.ID = m.seq
	clone := *l
	m.logs[l.ID] = &clone
	return nil
}

func (m *MockWebRTCQualityRepo) ListByCallID(_ context.Context, callID int64) ([]WebRTCQualityLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []WebRTCQualityLog
	for _, l := range m.logs {
		if l.CallID == callID {
			out = append(out, *l)
		}
	}
	return out, nil
}

func (m *MockWebRTCQualityRepo) ListByAgent(_ context.Context, tenantID, agentID int64, limit int) ([]WebRTCQualityLog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []WebRTCQualityLog
	for _, l := range m.logs {
		if l.TenantID == tenantID && l.AgentID == agentID {
			out = append(out, *l)
			if len(out) >= limit {
				break
			}
		}
	}
	return out, nil
}
