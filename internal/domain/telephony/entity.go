package telephony

import (
	"time"
)

// RoutingRule matches outbound calls to SIP trunks based on prefix/time/priority.
type RoutingRule struct {
	ID          int64     `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	MatchType   string    `db:"match_type" json:"match_type"`     // prefix, regex, time_of_day
	MatchValue  string    `db:"match_value" json:"match_value"`   // e.g. "+86" or "09:00-18:00"
	SIPTrunkID  int64     `db:"sip_trunk_id" json:"sip_trunk_id"`
	Priority    int       `db:"priority" json:"priority"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// CLIStrategy defines how outbound caller ID is selected.
type CLIStrategy string

const (
	CLIStrategyFixed      CLIStrategy = "fixed"       // always use a specific number
	CLIStrategyRoundRobin CLIStrategy = "round_robin"  // rotate through numbers
	CLIStrategyRandom     CLIStrategy = "random"       // random selection
	CLIStrategyMatchArea  CLIStrategy = "match_area"   // match callee area code
)

// CLIPolicy defines caller ID selection rules for outbound calls.
type CLIPolicy struct {
	ID            int64       `db:"id" json:"id"`
	TenantID      int64       `db:"tenant_id" json:"tenant_id"`
	Name          string      `db:"name" json:"name"`
	Strategy      CLIStrategy `db:"strategy" json:"strategy"`
	FixedNumberID *int64      `db:"fixed_number_id" json:"fixed_number_id,omitempty"`
	NumberPoolIDs string      `db:"number_pool_ids" json:"number_pool_ids"` // comma-separated phone_number IDs
	IsDefault     bool        `db:"is_default" json:"is_default"`
	CreatedAt     time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time   `db:"updated_at" json:"updated_at"`
}

type Carrier struct {
	ID          int64     `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	Protocol    string    `db:"protocol" json:"protocol"`
	Host        string    `db:"host" json:"host"`
	Port        int       `db:"port" json:"port"`
	Status      string    `db:"status" json:"status"`
	MaxChannels int       `db:"max_channels" json:"max_channels"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type SIPTrunk struct {
	ID          int64     `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	CarrierID   int64     `db:"carrier_id" json:"carrier_id"`
	Name        string    `db:"name" json:"name"`
	Username    string    `db:"username" json:"username"`
	Password    string    `db:"password" json:"-"`
	Domain      string    `db:"domain" json:"domain"`
	Transport   string    `db:"transport" json:"transport"`
	Codecs      string    `db:"codecs" json:"codecs"`
	MaxChannels int       `db:"max_channels" json:"max_channels"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type PhoneNumber struct {
	ID              int64     `db:"id" json:"id"`
	TenantID        int64     `db:"tenant_id" json:"tenant_id"`
	Number          string    `db:"number" json:"number"`
	DisplayName     string    `db:"display_name" json:"display_name"`
	Usage           string    `db:"usage" json:"usage"`
	SIPTrunkID      *int64    `db:"sip_trunk_id" json:"sip_trunk_id,omitempty"`
	IVRFlowID       *int64    `db:"ivr_flow_id" json:"ivr_flow_id,omitempty"`
	SkillGroupID    *int64    `db:"skill_group_id" json:"skill_group_id,omitempty"`
	Status          string    `db:"status" json:"status"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type CallNumberTag struct {
	ID        int64     `db:"id" json:"id"`
	TenantID  int64     `db:"tenant_id" json:"tenant_id"`
	Number    string    `db:"number" json:"number"`
	Tag       string    `db:"tag" json:"tag"`
	Source    string    `db:"source" json:"source"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type AutoTagRule struct {
	ID          int64     `db:"id" json:"id"`
	TenantID    int64     `db:"tenant_id" json:"tenant_id"`
	Name        string    `db:"name" json:"name"`
	MatchType   string    `db:"match_type" json:"match_type"`
	MatchValue  string    `db:"match_value" json:"match_value"`
	Tag         string    `db:"tag" json:"tag"`
	Priority    int       `db:"priority" json:"priority"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
