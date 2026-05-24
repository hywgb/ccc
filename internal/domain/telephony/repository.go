package telephony

import "context"

type RoutingRuleRepository interface {
	Create(ctx context.Context, r *RoutingRule) error
	GetByID(ctx context.Context, id int64) (*RoutingRule, error)
	Update(ctx context.Context, r *RoutingRule) error
	ListActive(ctx context.Context, tenantID int64) ([]*RoutingRule, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*RoutingRule, int64, error)
	Delete(ctx context.Context, id int64) error
}

type CLIPolicyRepository interface {
	Create(ctx context.Context, p *CLIPolicy) error
	GetByID(ctx context.Context, id int64) (*CLIPolicy, error)
	GetDefault(ctx context.Context, tenantID int64) (*CLIPolicy, error)
	Update(ctx context.Context, p *CLIPolicy) error
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*CLIPolicy, int64, error)
}

type CarrierRepository interface {
	Create(ctx context.Context, c *Carrier) error
	GetByID(ctx context.Context, id int64) (*Carrier, error)
	Update(ctx context.Context, c *Carrier) error
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*Carrier, int64, error)
}

type SIPTrunkRepository interface {
	Create(ctx context.Context, t *SIPTrunk) error
	GetByID(ctx context.Context, id int64) (*SIPTrunk, error)
	Update(ctx context.Context, t *SIPTrunk) error
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*SIPTrunk, int64, error)
	ListAll(ctx context.Context, offset, limit int) ([]*SIPTrunk, int64, error)
}

type PhoneNumberRepository interface {
	Create(ctx context.Context, p *PhoneNumber) error
	GetByID(ctx context.Context, id int64) (*PhoneNumber, error)
	GetByNumber(ctx context.Context, tenantID int64, number string) (*PhoneNumber, error)
	Update(ctx context.Context, p *PhoneNumber) error
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*PhoneNumber, int64, error)
}

type SIPTrunkGroupRepository interface {
	Create(ctx context.Context, g *SIPTrunkGroup) error
	GetByID(ctx context.Context, id int64) (*SIPTrunkGroup, error)
	Update(ctx context.Context, g *SIPTrunkGroup) error
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*SIPTrunkGroup, int64, error)
	AddMember(ctx context.Context, m *SIPTrunkGroupMember) error
	ListMembers(ctx context.Context, groupID int64) ([]*SIPTrunkGroupMember, error)
}

type CallNumberTagRepository interface {
	Create(ctx context.Context, t *CallNumberTag) error
	ListByNumber(ctx context.Context, tenantID int64, number string) ([]*CallNumberTag, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*CallNumberTag, int64, error)
	Delete(ctx context.Context, id int64) error
}

type AutoTagRuleRepository interface {
	Create(ctx context.Context, r *AutoTagRule) error
	GetByID(ctx context.Context, id int64) (*AutoTagRule, error)
	Update(ctx context.Context, r *AutoTagRule) error
	ListActive(ctx context.Context, tenantID int64) ([]*AutoTagRule, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]*AutoTagRule, int64, error)
	Delete(ctx context.Context, id int64) error
}
