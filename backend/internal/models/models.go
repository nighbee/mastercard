package models

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	FullName     string    `gorm:"not null" json:"full_name"`
	RoleID       *uint     `gorm:"index" json:"role_id"`
	Role         *Role     `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Role model
type Role struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `gorm:"uniqueIndex;not null" json:"name"`
	Description string     `json:"description,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Permission model
type Permission struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Resource   string    `gorm:"uniqueIndex:idx_resource_action" json:"resource"`
	Action     string    `gorm:"uniqueIndex:idx_resource_action" json:"action"`
	Conditions string    `gorm:"type:jsonb" json:"conditions,omitempty"`
	Roles      []Role    `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Transaction model
type Transaction struct {
	ID                        uint      `gorm:"primaryKey" json:"id"`
	CardNo                    string    `gorm:"not null" json:"card_no"`
	Date                      time.Time `gorm:"type:date;not null;index" json:"date"`
	ProcessDate               time.Time `gorm:"type:date;not null;index" json:"process_date"`
	TrxAmountUSD              *float64  `gorm:"type:decimal(15,2)" json:"trx_amount_usd,omitempty"`
	TrxAmountEUR              *float64  `gorm:"type:decimal(15,2)" json:"trx_amount_eur,omitempty"`
	TrxAmountLocal            *float64  `gorm:"type:decimal(15,2)" json:"trx_amount_local,omitempty"`
	TrxCntUSD                 int       `gorm:"default:0" json:"trx_cnt_usd"`
	TrxCntEUR                 int       `gorm:"default:0" json:"trx_cnt_eur"`
	TrxCntLocal               int       `gorm:"default:0" json:"trx_cnt_local"`
	InterchangeFee            *float64  `gorm:"type:decimal(15,2)" json:"interchange_fee,omitempty"`
	MerchName                 *string   `gorm:"index" json:"merch_name,omitempty"`
	AggMerchName              *string   `gorm:"index" json:"agg_merch_name,omitempty"`
	IssuerCode                *string   `gorm:"index" json:"issuer_code,omitempty"`
	IssuerCountry             *string   `gorm:"index" json:"issuer_country,omitempty"`
	Bin6Code                  *string   `gorm:"type:varchar(6)" json:"bin6_code,omitempty"`
	AcquirerCode              *string   `gorm:"index" json:"acquirer_code,omitempty"`
	AcquirerCountry           *string   `json:"acquirer_country,omitempty"`
	TrxType                   *string   `gorm:"index" json:"trx_type,omitempty"`
	TrxDirection              *string   `gorm:"type:varchar(10)" json:"trx_direction,omitempty"` // plus or minus
	MCC                       *string   `gorm:"index" json:"mcc,omitempty"`
	MCCGroup                  *string   `gorm:"index" json:"mcc_group,omitempty"`
	InputMode                 *string   `json:"input_mode,omitempty"`
	WalletType                *string   `json:"wallet_type,omitempty"`
	ProductType               *string   `json:"product_type,omitempty"`
	AuthorizationStatus       *string   `gorm:"index" json:"authorization_status,omitempty"`
	AuthorizationResponseCode *string   `json:"authorization_response_code,omitempty"`
	LocationID                *string   `json:"location_id,omitempty"`
	LocationCity              *string   `gorm:"index" json:"location_city,omitempty"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// Conversation model
type Conversation struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	UserID             uint       `gorm:"not null;index" json:"user_id"`
	User               User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title              *string    `json:"title,omitempty"`
	ParentBranchID     *uint      `gorm:"index" json:"parent_branch_id,omitempty"`
	ParentBranch       *Conversation `gorm:"foreignKey:ParentBranchID" json:"parent_branch,omitempty"`
	BranchPointMessageID *uint    `json:"branch_point_message_id,omitempty"`
	Messages           []Message  `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// Message model
type Message struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ConversationID uint       `gorm:"not null;index" json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
	UserMessage    string     `gorm:"type:text;not null" json:"user_message"`
	SQLQuery       *string    `gorm:"type:text" json:"sql_query,omitempty"`
	ResultData     *string    `gorm:"type:jsonb" json:"result_data,omitempty"`
	ResultFormat   *string    `gorm:"type:varchar(20)" json:"result_format,omitempty"` // text, table, chart, error
	ErrorMessage   *string    `gorm:"type:text" json:"error_message,omitempty"`
	ExecutionTimeMs *int      `json:"execution_time_ms,omitempty"`
	CreatedAt      time.Time  `gorm:"index" json:"created_at"`
}

// AuditLog model
type AuditLog struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	UserID          *uint      `gorm:"index" json:"user_id,omitempty"`
	User            *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action          string     `gorm:"not null;index" json:"action"`
	Resource        *string    `gorm:"index" json:"resource,omitempty"`
	QueryText       *string    `gorm:"type:text" json:"query_text,omitempty"`
	SQLExecuted     *string    `gorm:"type:text" json:"sql_executed,omitempty"`
	ResultCount     *int       `json:"result_count,omitempty"`
	IPAddress       *string    `gorm:"type:inet" json:"ip_address,omitempty"`
	UserAgent       *string    `gorm:"type:text" json:"user_agent,omitempty"`
	Status          *string    `gorm:"type:varchar(20);index" json:"status,omitempty"` // success, error, denied
	ErrorMessage    *string    `gorm:"type:text" json:"error_message,omitempty"`
	ExecutionTimeMs *int       `json:"execution_time_ms,omitempty"`
	Timestamp       time.Time  `gorm:"index" json:"timestamp"`
}

// TableName overrides
func (User) TableName() string {
	return "users"
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (Transaction) TableName() string {
	return "transactions"
}

func (Conversation) TableName() string {
	return "conversations"
}

func (Message) TableName() string {
	return "messages"
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// BeforeCreate hook for User
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

