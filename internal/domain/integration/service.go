package integration

import (
	"context"
	"time"

	"github.com/divord97/ccc/pkg/snowflake"
)

type DNCService struct {
	repo DNCRepository
}

func NewDNCService(repo DNCRepository) *DNCService {
	return &DNCService{repo: repo}
}

// CheckDNC returns ErrDNCBlocked if the number is on the DNC list and not expired.
func (s *DNCService) CheckDNC(ctx context.Context, tenantID int64, number string) error {
	entry, err := s.repo.GetByNumber(ctx, tenantID, number)
	if err != nil {
		return err
	}
	if entry == nil {
		return nil
	}
	if entry.ExpiresAt != nil && entry.ExpiresAt.Before(time.Now()) {
		return nil
	}
	return ErrDNCBlocked
}

// CheckBatch checks multiple numbers against DNC. Returns blocked numbers.
func (s *DNCService) CheckBatch(ctx context.Context, tenantID int64, numbers []string) ([]string, error) {
	return s.repo.CheckNumbers(ctx, tenantID, numbers)
}

// AddEntry adds a number to the DNC list.
func (s *DNCService) AddEntry(ctx context.Context, entry *DNCEntry) error {
	entry.ID = snowflake.NextID()
	entry.CreatedAt = time.Now()
	return s.repo.Create(ctx, entry)
}

// RemoveEntry removes a DNC entry.
func (s *DNCService) RemoveEntry(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

type CallTagService struct {
	repo CallTagAssignmentRepository
}

func NewCallTagService(repo CallTagAssignmentRepository) *CallTagService {
	return &CallTagService{repo: repo}
}

func (s *CallTagService) AssignTag(ctx context.Context, a *CallTagAssignment) error {
	existing, err := s.repo.ListByCallID(ctx, a.CallID)
	if err != nil {
		return err
	}
	for _, e := range existing {
		if e.TagID == a.TagID {
			return ErrTagAlreadyExists
		}
	}
	a.ID = snowflake.NextID()
	a.CreatedAt = time.Now()
	return s.repo.Create(ctx, a)
}

func (s *CallTagService) RemoveTag(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *CallTagService) GetCallTags(ctx context.Context, callID int64) ([]*CallTagAssignment, error) {
	return s.repo.ListByCallID(ctx, callID)
}
