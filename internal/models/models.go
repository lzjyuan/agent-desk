package models

import (
	"agent-desk/internal/pkg/enums"
	"time"
)

// Models 注册所有需要迁移和代码生成的模型。
var Models = []any{
	&Migration{},
	&Organization{},
	&OrganizationMember{},
	&User{},
	&UserIdentity{},
	&Company{},
	&Customer{},
	&CustomerIdentity{},
	&CustomerContact{},
	&Role{},
	&Permission{},
	&UserRole{},
	&RolePermission{},
	&UserPermission{},
	&LoginSession{},
	&LoginCredentialLog{},
	&Asset{},
	&Tag{},
	&Conversation{},
	&ConversationParticipant{},
	&ConversationReadState{},
	&Message{},
	&WxWorkKFSyncState{},
	&WxWorkKFConversation{},
	&WxWorkKFMessageRef{},
	&ChannelMessageOutbox{},
	&ConversationAssignment{},
	&ConversationTag{},
	&QuickReply{},
	&ConversationEventLog{},
	&Ticket{},
	&TicketTag{},
	&TicketProgress{},
	&TicketView{},
	&TicketNoSequence{},
	&Notification{},
	&AIAgent{},
	&Channel{},
	&AgentProfile{},
	&AgentTeam{},
	&AgentTeamSchedule{},
	&AIConfig{},
	&KnowledgeBase{},
	&KnowledgeDirectory{},
	&KnowledgeDocument{},
	&KnowledgeFAQ{},
	&KnowledgeChunk{},
	&KnowledgeRetrieveLog{},
	&KnowledgeRetrieveHit{},
	&KnowledgeFeedback{},
	&DocPage{},
	&Category{},
	&Post{},
	&Comment{},
	&Reaction{},
	&CommentReport{},
	&SkillDefinition{},
	&AgentRevision{},
	&AgentRun{},
	&AgentStep{},
	&AgentToolCall{},
	&AgentToolInvocation{},
	&AgentRunQualityFeedback{},
	&AIWorkflow{},
	&AIWorkflowVersion{},
	&AIAgentWorkflowBinding{},
	&AIWorkflowRun{},
	&AIWorkflowNodeRun{},
	&ConversationInterrupt{},
	&SystemConfig{},
}

// AgentToolInvocation persists the idempotency boundary for a business tool.
// It is intentionally independent of AgentRun audit rows so a retry after a
// process interruption cannot repeat an external write.
type AgentToolInvocation struct {
	ID             int64  `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConversationID int64  `gorm:"type:bigint;not null;index;uniqueIndex:uk_agent_tool_invocation;comment:会话 ID"`
	AIAgentID      int64  `gorm:"type:bigint;not null;default:0;index;comment:AI Agent ID"`
	ToolCode       string `gorm:"type:varchar(128);not null;default:'';uniqueIndex:uk_agent_tool_invocation;comment:工具编码"`
	IdempotencyKey string `gorm:"type:varchar(160);not null;default:'';uniqueIndex:uk_agent_tool_invocation;comment:幂等键"`
	Status         string `gorm:"type:varchar(20);not null;default:'running';index;comment:状态"`
	ResultData     string `gorm:"type:text;comment:结果数据"`
	ErrorMessage   string `gorm:"type:text;comment:错误信息"`
	AuditFields
}

type Migration struct {
	ID         int64     `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	Version    int64     `gorm:"type:bigint;not null;uniqueIndex;comment:版本号"`
	Remark     string    `gorm:"type:text;comment:备注"`
	Success    bool      `gorm:"not null;default:false;comment:是否成功"`
	ErrorInfo  string    `gorm:"type:text;comment:错误详情"`
	RetryCount int       `gorm:"type:int;not null;default:0;comment:重试次数"`
	CreatedAt  time.Time `gorm:"comment:创建时间"`
	UpdatedAt  time.Time `gorm:"comment:更新时间"`
}

// SystemConfig 运营侧系统配置项；具体有哪些 config_key 由业务代码约定，表内一行一项。
type SystemConfig struct {
	ID          int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConfigKey   string       `gorm:"column:config_key;type:varchar(128);not null;uniqueIndex:uk_system_config_group_key;comment:配置键"`
	ConfigValue string       `gorm:"column:config_value;type:text;not null;comment:配置值"`
	GroupCode   string       `gorm:"column:group_code;type:varchar(64);not null;default:'';index;uniqueIndex:uk_system_config_group_key;comment:分组编码"`
	Title       string       `gorm:"type:varchar(200);not null;default:'';comment:标题"`
	Description string       `gorm:"type:text;comment:描述"`
	Status      enums.Status `gorm:"type:int;not null;default:0;index;comment:状态"`
	AuditFields
}

// TicketNoSequence 工单号日序列表。
//
// 每天一条记录，NextSeq 表示当日下一次可分配的序号。
type TicketNoSequence struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	DateKey   string    `gorm:"column:date_key;type:varchar(8);not null;uniqueIndex;comment:日期键"`
	NextSeq   int64     `gorm:"column:next_seq;type:bigint;not null;default:1;comment:下一个序号"`
	CreatedAt time.Time `gorm:"not null;index;comment:创建时间"`
	UpdatedAt time.Time `gorm:"not null;index;comment:更新时间"`
}

// TicketView 工单工作台个人保存视图。
type TicketView struct {
	ID          int64  `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	UserID      int64  `gorm:"column:user_id;type:bigint;not null;index;comment:用户 ID"`
	Name        string `gorm:"column:name;type:varchar(100);not null;default:'';index;comment:名称"`
	FiltersJSON string `gorm:"column:filters_json;type:text;not null;comment:筛选条件 JSON"`
	SortNo      int    `gorm:"column:sort_no;type:int;not null;default:0;index;comment:排序号"`
	AuditFields
}

// Notification 站内通知。
type Notification struct {
	ID               int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	RecipientUserID  int64        `gorm:"type:bigint;not null;default:0;index;comment:接收用户 ID"`
	Title            string       `gorm:"type:varchar(255);not null;default:'';comment:标题"`
	Content          string       `gorm:"type:text;comment:内容"`
	NotificationType string       `gorm:"type:varchar(50);not null;default:'';index;comment:通知类型"`
	BizType          string       `gorm:"type:varchar(50);not null;default:'';index;comment:业务类型"`
	BizID            int64        `gorm:"type:bigint;not null;default:0;index;comment:业务 ID"`
	ActionURL        string       `gorm:"type:varchar(255);not null;default:'';comment:操作链接"`
	ReadAt           *time.Time   `gorm:"index;comment:阅读时间"`
	Status           enums.Status `gorm:"type:int;not null;default:0;index;comment:状态"`
	CreatedAt        time.Time    `gorm:"not null;index;comment:创建时间"`
}

// AuditFields 定义涉及用户操作数据的统一审计字段。
// 该结构记录数据创建与更新的时间、操作者ID和操作者名称。
type AuditFields struct {
	CreatedAt      time.Time `gorm:"not null;index;comment:CreatedAt 记录数据创建时间。"`                                          // CreatedAt 记录数据创建时间。
	CreateUserID   int64     `gorm:"type:bigint;not null;default:0;index;comment:CreateUserID 记录创建人用户ID；系统任务写0。"`         // CreateUserID 记录创建人用户ID；系统任务写0。
	CreateUserName string    `gorm:"type:varchar(100);not null;default:'';comment:CreateUserName 记录创建人名称；系统任务写system。"`   // CreateUserName 记录创建人名称；系统任务写system。
	UpdatedAt      time.Time `gorm:"not null;index;comment:UpdatedAt 记录数据最近更新时间。"`                                        // UpdatedAt 记录数据最近更新时间。
	UpdateUserID   int64     `gorm:"type:bigint;not null;default:0;index;comment:UpdateUserID 记录最后更新人用户ID；系统任务写0。"`       // UpdateUserID 记录最后更新人用户ID；系统任务写0。
	UpdateUserName string    `gorm:"type:varchar(100);not null;default:'';comment:UpdateUserName 记录最后更新人名称；系统任务写system。"` // UpdateUserName 记录最后更新人名称；系统任务写system。
}

// Organization 组织 / 工作区（Tenant / Workspace）。
type Organization struct {
	ID     int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	Code   string       `gorm:"type:varchar(100);not null;uniqueIndex:uk_org_code;comment:组织唯一编码 / Slug，如 org_123456"` // 组织唯一编码 / Slug，如 org_123456
	Name   string       `gorm:"type:varchar(200);not null;default:'';comment:组织名称"`                                    // 组织名称
	Logo   string       `gorm:"type:varchar(255);not null;default:'';comment:组织 Logo URL"`                             // 组织 Logo URL
	Plan   string       `gorm:"type:varchar(50);not null;default:'free';comment:订阅计划（free, pro, enterprise）"`          // 订阅计划（free, pro, enterprise）
	Status enums.Status `gorm:"type:int;not null;default:0;index;comment:组织状态"`                                        // 组织状态
	Remark string       `gorm:"type:text;comment:备注"`                                                                  // 备注
	AuditFields
}

// OrganizationMember 组织成员关联表。
type OrganizationMember struct {
	ID             int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	OrganizationID int64        `gorm:"type:bigint;not null;index;uniqueIndex:uk_org_member;comment:组织 ID"`                 // 组织 ID
	UserID         int64        `gorm:"type:bigint;not null;index;uniqueIndex:uk_org_member;comment:用户 ID"`                 // 用户 ID
	Role           string       `gorm:"type:varchar(50);not null;default:'MEMBER';index;comment:角色：OWNER | ADMIN | MEMBER"` // 角色：OWNER | ADMIN | MEMBER
	Status         enums.Status `gorm:"type:int;not null;default:0;index;comment:成员状态"`                                     // 成员状态
	AuditFields
}

// User 后台用户账号。
type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	Username     string         `gorm:"type:varchar(100);not null;uniqueIndex;comment:用户名"`
	Nickname     string         `gorm:"type:varchar(100);not null;default:'';index;comment:昵称"`
	Avatar       string         `gorm:"type:varchar(255);not null;default:'';comment:头像地址"`
	Mobile       *string        `gorm:"type:varchar(32);uniqueIndex;comment:手机号"`
	Email        *string        `gorm:"type:varchar(100);uniqueIndex;comment:邮箱"`
	Password     string         `gorm:"type:varchar(255);not null;default:'';comment:密码哈希"`
	PasswordSalt string         `gorm:"type:varchar(64);not null;default:'';comment:密码盐"`
	UserType     enums.UserType `gorm:"column:user_type;type:varchar(20);not null;default:'user';index;comment:用户类型"`
	ActiveOrgID  int64          `gorm:"column:active_org_id;type:bigint;not null;default:0;index;comment:当前组织 ID"`
	Status       enums.Status   `gorm:"type:int;not null;default:0;index;comment:状态"`
	LastLoginAt  *time.Time     `gorm:"comment:最近登录时间"`
	LastLoginIP  string         `gorm:"type:varchar(64);not null;default:'';comment:最近登录 IP"`
	Remark       string         `gorm:"type:text;comment:备注"`
	DeletedAt    *time.Time     `gorm:"index;comment:删除时间"`
	AuditFields
}

// UserIdentity 第三方身份绑定信息。
type UserIdentity struct {
	ID              int64               `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	UserID          int64               `gorm:"type:bigint;not null;index;uniqueIndex:uk_provider_user;comment:用户 ID"`
	Provider        enums.ThirdProvider `gorm:"type:varchar(50);not null;default:'';index;uniqueIndex:uk_provider_user;uniqueIndex:uk_provider_union;comment:服务提供方"`
	ProviderUserID  string              `gorm:"type:varchar(128);not null;default:'';uniqueIndex:uk_provider_user;comment:提供方用户 ID"`
	ProviderUnionID *string             `gorm:"type:varchar(128);uniqueIndex:uk_provider_union;comment:提供方联合 ID"`
	ProviderCorpID  string              `gorm:"type:varchar(128);not null;default:'';index;comment:提供方企业 ID"`
	ProviderName    string              `gorm:"type:varchar(100);not null;default:'';comment:提供方名称"`
	RawProfile      string              `gorm:"type:text;comment:原始档案数据"`
	Status          enums.Status        `gorm:"type:int;not null;default:0;index;comment:状态"`
	LastAuthAt      *time.Time          `gorm:"comment:最近认证时间"`
	AuditFields
}

// Company 客户公司（组织）表。
//
//	用于存储公司主体信息；Customer（人）可通过 CompanyID 关联到所属公司。
type Company struct {
	ID     int64        `gorm:"primaryKey;autoIncrement;comment:公司主键。"`                                   // ID 为公司主键。
	Name   string       `gorm:"type:varchar(200);not null;uniqueIndex:uk_company_name;comment:公司名称（唯一）。"` // Name 为公司名称（唯一）。
	Code   string       `gorm:"type:varchar(64);not null;index;comment:公司编码/统一社会信用代码（可空语义用空串表示）。"`        // Code 为公司编码/统一社会信用代码（可空语义用空串表示）。
	Status enums.Status `gorm:"type:int;not null;default:0;comment:公司状态。"`                                // Status 为公司状态。
	Remark string       `gorm:"type:text;comment:备注。"`                                                    // Remark 为备注。
	AuditFields
}

// Customer 客户主表。
//
//	用于存储客户稳定画像信息，不包含平台身份映射和多联系方式明细。
type Customer struct {
	ID            int64        `gorm:"primaryKey;autoIncrement;comment:客户主键。"`                              // ID 为客户主键。
	UserID        int64        `gorm:"type:bigint;not null;default:0;index;comment:关联的统一登录账号ID，0表示未绑定。"`    // UserID 为关联的统一登录账号ID，0表示未绑定。
	Name          string       `gorm:"type:varchar(100);not null;default:'';index;comment:客户姓名或展示名称。"`      // Name 为客户姓名或展示名称。
	Gender        enums.Gender `gorm:"type:int;not null;default:0;;comment:性别：0未知 1男 2女。"`                  // Gender 为性别：0未知 1男 2女。
	CompanyID     int64        `gorm:"type:bigint;not null;default:0;index;comment:所属公司ID；0表示无所属公司（个人客户）。"` // CompanyID 为所属公司ID；0表示无所属公司（个人客户）。
	LastActiveAt  *time.Time   `gorm:"comment:最近活跃时间。"`                                                     // LastActiveAt 为最近活跃时间。
	PrimaryMobile string       `gorm:"type:varchar(32);not null;default:'';index;comment:主手机号（冗余展示字段）。"`    // PrimaryMobile 为主手机号（冗余展示字段）。
	PrimaryEmail  string       `gorm:"type:varchar(100);not null;default:'';index;comment:主邮箱（冗余展示字段）。"`    // PrimaryEmail 为主邮箱（冗余展示字段）。
	Status        enums.Status `gorm:"type:int;not null;default:0;;comment:客户状态。"`                          // Status 为客户状态。
	Remark        string       `gorm:"type:text;comment:备注。"`                                               // Remark 为备注。
	AuditFields
}

// CustomerIdentity 客户第三方身份映射表。
type CustomerIdentity struct {
	ID             int64                `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	CustomerID     int64                `gorm:"type:bigint;not null;uniqueIndex:uk_customer_external;comment:所属客户ID。"`                                         // 为所属客户ID。
	ExternalSource enums.ExternalSource `gorm:"type:varchar(30);uniqueIndex:uk_customer_external;comment:外部身份来源"`                                              // 为外部身份来源
	ExternalID     string               `gorm:"type:varchar(128);index:idx_external_id;uniqueIndex:uk_customer_external;comment:平台侧用户唯一ID，与访客 ExternalID 对齐。"` // 为平台侧用户唯一ID，与访客 ExternalID 对齐。
	RawProfile     string               `gorm:"type:text;comment:第三方原始资料JSON。"`                                                                                // 为第三方原始资料JSON。
	Status         enums.Status         `gorm:"type:int;not null;default:0;index;comment:映射状态。"`                                                               // 为映射状态。
	AuditFields
}

// CustomerContact 客户联系方式表。
//
//	用于维护客户的一对多联系方式，支持主联系方式、验证状态与失效标记。
type CustomerContact struct {
	ID           int64             `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	CustomerID   int64             `gorm:"type:bigint;not null;index;uniqueIndex:uk_customer_contact;comment:所属客户ID。"`                                           // CustomerID 为所属客户ID。
	ContactType  enums.ContactType `gorm:"type:varchar(30);not null;default:'';index;uniqueIndex:uk_customer_contact;comment:联系方式类型：mobile/email/wechat/other。"` // ContactType 为联系方式类型：mobile/email/wechat/other。
	ContactValue string            `gorm:"type:varchar(200);not null;default:'';index;uniqueIndex:uk_customer_contact;comment:联系方式值。"`                           // ContactValue 为联系方式值。
	IsPrimary    bool              `gorm:"not null;default:false;index;comment:是否主联系方式。"`                                                                        // IsPrimary 表示是否主联系方式。
	IsVerified   bool              `gorm:"not null;default:false;index;comment:是否已验证。"`                                                                          // IsVerified 表示是否已验证。
	VerifiedAt   *time.Time        `gorm:"comment:验证时间。"`                                                                                                        // VerifiedAt 为验证时间。
	Source       string            `gorm:"type:varchar(30);not null;default:'';index;comment:来源：manual/import/system。"`                                          // Source 为来源：manual/import/system。
	Status       enums.Status      `gorm:"type:int;not null;default:0;index;comment:联系方式状态。"`                                                                    // Status 为联系方式状态。
	Remark       string            `gorm:"type:varchar(255);not null;default:'';comment:备注。"`                                                                    // Remark 为备注。
	AuditFields
}

// Role 角色定义。
type Role struct {
	ID       int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	Name     string       `gorm:"type:varchar(100);not null;default:'';index;comment:名称"`
	Code     string       `gorm:"type:varchar(100);not null;uniqueIndex;comment:编码"`
	Status   enums.Status `gorm:"type:int;not null;default:0;index;comment:状态"`
	IsSystem bool         `gorm:"not null;default:false;index;comment:是否系统内置"`
	SortNo   int          `gorm:"type:int;not null;default:0;index;comment:排序号"`
	Remark   string       `gorm:"type:text;comment:备注"`
	AuditFields
}

// Permission 权限点定义。
type Permission struct {
	ID        int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	Name      string       `gorm:"type:varchar(100);not null;default:'';comment:名称"`
	Code      string       `gorm:"type:varchar(150);not null;uniqueIndex;comment:编码"`
	Type      string       `gorm:"type:varchar(20);not null;default:'';index;comment:类型"`
	GroupName string       `gorm:"type:varchar(100);not null;default:'';index;comment:分组名称"`
	ParentID  int64        `gorm:"type:bigint;not null;default:0;index;comment:父级 ID"`
	Path      string       `gorm:"type:varchar(255);not null;default:'';comment:路径"`
	Method    string       `gorm:"type:varchar(20);not null;default:'';comment:请求方法"`
	APIPath   string       `gorm:"type:varchar(255);not null;default:'';comment:API 路径"`
	SortNo    int          `gorm:"type:int;not null;default:0;index;comment:排序号"`
	Status    enums.Status `gorm:"type:int;not null;default:0;index;comment:状态"`
	IsBuiltin bool         `gorm:"not null;default:true;index;comment:是否内置"`
	Remark    string       `gorm:"type:text;comment:备注"`
	AuditFields
}

// UserRole 用户和角色关联。
type UserRole struct {
	ID     int64 `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	UserID int64 `gorm:"type:bigint;not null;index;uniqueIndex:uk_user_role;comment:用户 ID"`
	RoleID int64 `gorm:"type:bigint;not null;index;uniqueIndex:uk_user_role;comment:角色 ID"`
	AuditFields
}

// RolePermission 角色和权限关联。
type RolePermission struct {
	ID           int64 `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	RoleID       int64 `gorm:"type:bigint;not null;index;uniqueIndex:uk_role_permission;comment:角色 ID"`
	PermissionID int64 `gorm:"type:bigint;not null;index;uniqueIndex:uk_role_permission;comment:权限 ID"`
	AuditFields
}

// UserPermission 用户级例外权限。
//
//	用于处理少量临时授权或拒绝授权场景。
type UserPermission struct {
	ID           int64      `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	UserID       int64      `gorm:"type:bigint;not null;index;uniqueIndex:uk_user_permission;comment:用户 ID"`
	PermissionID int64      `gorm:"type:bigint;not null;index;uniqueIndex:uk_user_permission;comment:权限 ID"`
	Effect       int        `gorm:"type:int;not null;default:1;index;comment:权限生效方式：1允许 -1拒绝。"` // Effect 表示权限生效方式：1允许 -1拒绝。
	ExpiredAt    *time.Time `gorm:"comment:过期时间"`
	Remark       string     `gorm:"type:text;comment:备注"`
	AuditFields
}

// LoginSession 表示一次后台登录会话。
type LoginSession struct {
	ID         int64      `gorm:"primaryKey;autoIncrement;comment:登录会话主键。"`                                        // ID 为登录会话主键。
	UserID     int64      `gorm:"type:bigint;not null;index;comment:登录用户 ID。"`                                     // UserID 为登录用户 ID。
	Token      string     `gorm:"type:varchar(128);not null;uniqueIndex;comment:随机不透明登录凭证，使用 ak_ 前缀。"`             // Token 为随机不透明登录凭证，使用 ak_ 前缀。
	ClientType string     `gorm:"type:varchar(50);not null;default:'';index;comment:客户端类型，后台 Web 端固定为 admin_web。"` // ClientType 为客户端类型，后台 Web 端固定为 admin_web。
	ClientIP   string     `gorm:"type:varchar(64);not null;default:'';comment:登录请求来源 IP。"`                         // ClientIP 为登录请求来源 IP。
	UserAgent  string     `gorm:"type:varchar(255);not null;default:'';comment:登录请求浏览器或客户端 UA。"`                   // UserAgent 为登录请求浏览器或客户端 UA。
	ExpiredAt  time.Time  `gorm:"not null;index;comment:token 过期时间。"`                                              // ExpiredAt 为 token 过期时间。
	RevokedAt  *time.Time `gorm:"index;comment:主动注销或踢下线时间，非空表示已失效。"`                                               // RevokedAt 为主动注销或踢下线时间，非空表示已失效。
	LastSeenAt *time.Time `gorm:"comment:最近一次成功鉴权时间。"`                                                             // LastSeenAt 为最近一次成功鉴权时间。
	AuditFields
}

// LoginCredentialLog 记录一次后台登录凭证校验结果。
type LoginCredentialLog struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:登录凭证日志主键。"`                      // ID 为登录凭证日志主键。
	Principal string    `gorm:"type:varchar(100);not null;default:'';index;comment:用户输入的登录名。"`   // Principal 为用户输入的登录名。
	UserID    int64     `gorm:"type:bigint;not null;default:0;index;comment:匹配到的用户 ID，未匹配时为 0。"` // UserID 为匹配到的用户 ID，未匹配时为 0。
	Success   bool      `gorm:"not null;default:false;index;comment:本次凭证校验是否成功。"`                // Success 表示本次凭证校验是否成功。
	ClientIP  string    `gorm:"type:varchar(64);not null;default:'';comment:登录请求来源 IP。"`         // ClientIP 为登录请求来源 IP。
	UserAgent string    `gorm:"type:varchar(255);not null;default:'';comment:登录请求浏览器或客户端 UA。"`   // UserAgent 为登录请求浏览器或客户端 UA。
	Reason    string    `gorm:"type:varchar(255);not null;default:'';comment:校验结果原因。"`           // Reason 为校验结果原因。
	CreatedAt time.Time `gorm:"not null;index;comment:日志创建时间。"`                                  // CreatedAt 为日志创建时间。
}

// Asset 存储的文件资源，如上传的附件等。
type Asset struct {
	ID         int64               `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	AssetID    string              `gorm:"type:varchar(64);not null;uniqueIndex;comment:资源 ID"`
	Provider   enums.AssetProvider `gorm:"type:varchar(50);not null;default:'';index;comment:服务提供方"`
	StorageKey string              `gorm:"type:varchar(255);not null;default:'';uniqueIndex:uk_storage_key;comment:存储键"`
	Filename   string              `gorm:"type:varchar(255);not null;default:'';comment:文件名"`
	FileSize   int64               `gorm:"type:bigint;not null;default:0;comment:文件大小"`
	MimeType   string              `gorm:"type:varchar(100);not null;default:'';comment:媒体类型"`
	Status     enums.AssetStatus   `gorm:"type:int;not null;default:1;index;comment:状态"`
	AuditFields
}

type Tag struct {
	ID       int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ParentID int64        `gorm:"type:bigint;not null;index;comment:父级 ID"`
	Name     string       `gorm:"type:varchar(50);not null;;comment:名称"`
	Remark   string       `gorm:"type:text;;comment:备注"`
	SortNo   int          `gorm:"type:int;not null;default:0;comment:排序号"`
	Status   enums.Status `gorm:"type:int;not null;default:0;comment:状态"`
	AuditFields
}

// Conversation 客服会话。
type Conversation struct {
	ID                  int64                           `gorm:"primaryKey;autoIncrement;comment:会话主键。"`                                  // ID 为会话主键。
	AIAgentID           int64                           `gorm:"type:bigint;not null;default:0;index;comment:当前会话绑定的 AI Agent ID。"`       // AIAgentID 为当前会话绑定的 AI Agent ID。
	ChannelID           int64                           `gorm:"type:bigint;not null;default:0;index;comment:该会话来源接入渠道ID。"`               // ChannelID 为该会话来源接入渠道ID。
	CustomerID          int64                           `gorm:"type:bigint;not null;default:0;index;comment:会话所属客户 ID。"`                 // CustomerID 为会话所属客户 ID。
	CustomerName        string                          `gorm:"type:varchar(100);not null;default:'';index;comment:客户名称冗余字段，用于列表展示和搜索。"` // CustomerName 为客户名称冗余字段，用于列表展示和搜索。
	Status              enums.IMConversationStatus      `gorm:"type:int;not null;default:1;index;comment:会话状态，如待接入、处理中、已关闭。"`            // Status 为会话状态，如待接入、处理中、已关闭。
	ServiceMode         enums.IMConversationServiceMode `gorm:"type:int;not null;default:3;index;comment:服务模式，如仅AI、仅人工、AI优先人工接管。"`       // ServiceMode 为服务模式，如仅AI、仅人工、AI优先人工接管。
	Priority            int                             `gorm:"type:int;not null;default:0;index;comment:会话优先级。"`                        // Priority 为会话优先级。
	CurrentAssigneeID   int64                           `gorm:"type:bigint;not null;default:0;index;comment:当前接待客服ID。"`                  // CurrentAssigneeID 为当前接待客服ID。
	CurrentTeamID       int64                           `gorm:"type:bigint;not null;default:0;index;comment:当前处理客服组ID。"`                 // CurrentTeamID 为当前处理客服组ID。
	LastMessageID       int64                           `gorm:"type:bigint;not null;default:0;index;comment:最后一条消息ID。"`                  // LastMessageID 为最后一条消息ID。
	LastMessageAt       time.Time                       `gorm:"index;comment:最后消息时间。"`                                                   // LastMessageAt 为最后消息时间。
	LastActiveAt        time.Time                       `gorm:"index;comment:会话最近活跃时间。"`                                                 // LastActiveAt 为会话最近活跃时间。
	LastMessageSummary  string                          `gorm:"type:varchar(255);not null;default:'';comment:最后一条消息摘要。"`                 // LastMessageSummary 为最后一条消息摘要。
	CustomerUnreadCount int                             `gorm:"type:int;not null;default:0;comment:用户侧未读数。"`                             // CustomerUnreadCount 为用户侧未读数。
	AgentUnreadCount    int                             `gorm:"type:int;not null;default:0;comment:客服侧未读数。"`                             // AgentUnreadCount 为客服侧未读数。
	HandoffAt           *time.Time                      `gorm:"index;comment:最近一次转人工时间。"`                                                // HandoffAt 为最近一次转人工时间。
	HandoffReason       string                          `gorm:"type:varchar(255);not null;default:'';comment:最近一次转人工原因。"`                // HandoffReason 为最近一次转人工原因。
	AIReplyRounds       int                             `gorm:"type:int;not null;default:0;comment:当前会话内 AI 已成功回复次数。"`                   // AIReplyRounds 为当前会话内 AI 已成功回复次数。
	ClosedAt            *time.Time                      `gorm:"index;comment:会话关闭时间。"`                                                   // ClosedAt 为会话关闭时间。
	ClosedBy            int64                           `gorm:"type:bigint;not null;default:0;index;comment:关闭人用户ID，访客关闭时写0。"`           // ClosedBy 为关闭人用户ID，访客关闭时写0。
	CloseReason         string                          `gorm:"type:varchar(255);not null;default:'';comment:关闭原因。"`                     // CloseReason 为关闭原因。
	AuditFields
}

// ConversationParticipant 会话参与方。
type ConversationParticipant struct {
	ID                    int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConversationID        int64        `gorm:"type:bigint;not null;index;uniqueIndex:uk_conversation_participant;comment:会话 ID"`
	ParticipantType       string       `gorm:"type:varchar(30);not null;default:'';index;uniqueIndex:uk_conversation_participant;comment:参与方类型"`
	ParticipantID         int64        `gorm:"type:bigint;not null;default:0;uniqueIndex:uk_conversation_participant;comment:参与方 ID"`
	ExternalParticipantID string       `gorm:"type:varchar(128);not null;default:'';comment:外部参与方 ID"`
	JoinedAt              *time.Time   `gorm:"comment:加入时间"`
	LeftAt                *time.Time   `gorm:"comment:离开时间"`
	Status                enums.Status `gorm:"type:int;not null;default:0;index;comment:状态"`
	AuditFields
}

// ConversationReadState 会话读游标。
type ConversationReadState struct {
	ID                int64              `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConversationID    int64              `gorm:"type:bigint;not null;index;uniqueIndex:uk_conversation_reader;comment:会话 ID"`
	ReaderType        enums.IMSenderType `gorm:"type:varchar(30);not null;default:'';index;uniqueIndex:uk_conversation_reader;comment:阅读者类型"`
	ReaderID          int64              `gorm:"type:bigint;not null;default:0;uniqueIndex:uk_conversation_reader;comment:阅读者 ID"`
	ExternalReaderID  string             `gorm:"type:varchar(128);not null;default:'';uniqueIndex:uk_conversation_reader;comment:外部阅读者 ID"`
	LastReadMessageID int64              `gorm:"type:bigint;not null;default:0;index;comment:最后已读消息 ID"`
	LastReadAt        *time.Time         `gorm:"comment:最后已读时间"`
	AuditFields
}

// Message 会话消息。
type Message struct {
	ID              int64                 `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConversationID  int64                 `gorm:"type:bigint;not null;index;uniqueIndex:uk_conversation_client_msg;comment:会话 ID"`
	RequestID       string                `gorm:"type:varchar(128);not null;default:'';index;comment:请求 ID"`
	WorkflowRunID   int64                 `gorm:"type:bigint;not null;default:0;index;comment:工作流运行 ID"`
	ClientMsgID     string                `gorm:"type:varchar(128);not null;default:'';uniqueIndex:uk_conversation_client_msg;comment:客户端消息 ID"`
	SenderType      enums.IMSenderType    `gorm:"type:varchar(30);not null;default:'';index;comment:发送者类型"`
	SenderID        int64                 `gorm:"type:bigint;not null;default:0;index;comment:发送者 ID"`
	ReceiverType    string                `gorm:"type:varchar(30);not null;default:'';index;comment:接收者类型"`
	MessageType     enums.IMMessageType   `gorm:"type:varchar(30);not null;default:'';index;comment:消息类型"`
	Content         string                `gorm:"type:text;comment:内容"`
	Payload         string                `gorm:"type:text;comment:消息载荷"`
	SendStatus      enums.IMMessageStatus `gorm:"type:int;not null;default:2;index;comment:发送状态"`
	SentAt          *time.Time            `gorm:"index;comment:发送时间"`
	DeliveredAt     *time.Time            `gorm:"comment:送达时间"`
	ReadAt          *time.Time            `gorm:"comment:阅读时间"`
	RecalledAt      *time.Time            `gorm:"comment:撤回时间"`
	QuotedMessageID int64                 `gorm:"type:bigint;not null;default:0;index;comment:引用消息 ID"`
	AuditFields
}

// WxWorkKFSyncState 企业微信客服消息同步状态表。
//
//	按 open_kfid 记录企业微信客服消息同步游标，用于 SyncMsg 增量拉取。
type WxWorkKFSyncState struct {
	ID         int64        `gorm:"primaryKey;autoIncrement;comment:同步状态主键。"`                             // ID 为同步状态主键。
	OpenKfID   string       `gorm:"type:varchar(64);not null;default:'';uniqueIndex;comment:企业微信客服账号ID。"` // OpenKfID 为企业微信客服账号ID。
	NextCursor string       `gorm:"type:varchar(128);not null;default:'';comment:下一次增量同步使用的游标。"`          // NextCursor 为下一次增量同步使用的游标。
	LastSyncAt *time.Time   `gorm:"index;comment:最近一次成功同步时间。"`                                            // LastSyncAt 为最近一次成功同步时间。
	Status     enums.Status `gorm:"type:int;not null;default:0;index;comment:同步状态记录状态。"`                  // Status 为同步状态记录状态。
	Remark     string       `gorm:"type:text;comment:同步异常、人工备注等补充信息。"`                                    // Remark 为同步异常、人工备注等补充信息。
	AuditFields
}

// WxWorkKFConversation 企业微信客服渠道会话映射表。
//
//	维护平台会话与企业微信客服会话上下文的对应关系，供入站同步和下行发送复用。
type WxWorkKFConversation struct {
	ID             int64        `gorm:"primaryKey;autoIncrement;comment:渠道会话映射主键。"`                                    // ID 为渠道会话映射主键。
	ConversationID int64        `gorm:"type:bigint;not null;uniqueIndex;comment:平台会话ID，一条平台会话仅对应一条当前有效渠道映射。"`          // ConversationID 为平台会话ID，一条平台会话仅对应一条当前有效渠道映射。
	ChannelID      int64        `gorm:"type:bigint;not null;default:0;index;comment:所属接入渠道ID，用于标识该会话来自哪个企业微信渠道配置。"`    // ChannelID 为所属接入渠道ID，用于标识该会话来自哪个企业微信渠道配置。
	OpenKfID       string       `gorm:"type:varchar(64);not null;default:'';index:idx_openkf_ext;comment:企业微信客服账号ID。"` // OpenKfID 为企业微信客服账号ID。
	ExternalUserID string       `gorm:"type:varchar(128);not null;default:'';index:idx_openkf_ext;comment:企业微信客户ID。"`  // ExternalUserID 为企业微信客户ID。
	ServicerUserID string       `gorm:"type:varchar(128);not null;default:'';index;comment:企业微信当前接待客服成员UserID。"`       // ServicerUserID 为企业微信当前接待客服成员UserID。
	SessionStatus  string       `gorm:"type:varchar(30);not null;default:'';index;comment:微信侧会话状态快照，如接入中、转接中、已结束。"`    // SessionStatus 为微信侧会话状态快照，如接入中、转接中、已结束。
	LastWxMsgID    string       `gorm:"type:varchar(64);not null;default:'';index;comment:最近一次同步到的微信消息ID。"`            // LastWxMsgID 为最近一次同步到的微信消息ID。
	LastWxMsgTime  *time.Time   `gorm:"index;comment:最近一次微信消息时间。"`                                                     // LastWxMsgTime 为最近一次微信消息时间。
	RawProfile     string       `gorm:"type:text;comment:微信侧原始会话补充信息JSON。"`                                            // RawProfile 为微信侧原始会话补充信息JSON。
	Status         enums.Status `gorm:"type:int;not null;default:0;index;comment:渠道会话映射状态。"`                           // Status 为渠道会话映射状态。
	AuditFields
}

// WxWorkKFMessageRef 企业微信客服消息映射表。
//
//	用于实现微信消息幂等消费，并保存平台消息与微信消息的双向映射关系。
type WxWorkKFMessageRef struct {
	ID             int64        `gorm:"primaryKey;autoIncrement;comment:消息映射主键。"`                                   // ID 为消息映射主键。
	ConversationID int64        `gorm:"type:bigint;not null;default:0;index;comment:所属平台会话ID。"`                     // ConversationID 为所属平台会话ID。
	MessageID      int64        `gorm:"type:bigint;not null;default:0;index;comment:所属平台消息ID；仅渠道消息尚未生成平台消息时可暂为0。"`  // MessageID 为所属平台消息ID；仅渠道消息尚未生成平台消息时可暂为0。
	WxMsgID        string       `gorm:"type:varchar(64);not null;default:'';uniqueIndex;comment:企业微信消息ID，用于幂等去重。"`  // WxMsgID 为企业微信消息ID，用于幂等去重。
	Direction      string       `gorm:"type:varchar(20);not null;default:'';index;comment:消息方向，如 in/out。"`          // Direction 为消息方向，如 in/out。
	Origin         int          `gorm:"type:int;not null;default:0;index;comment:企业微信消息来源值，如客户发送、系统事件、企微客户端发送。"`    // Origin 为企业微信消息来源值，如客户发送、系统事件、企微客户端发送。
	OpenKfID       string       `gorm:"type:varchar(64);not null;default:'';index;comment:发送或接收该消息的客服账号ID。"`        // OpenKfID 为发送或接收该消息的客服账号ID。
	ExternalUserID string       `gorm:"type:varchar(128);not null;default:'';index;comment:消息对应的企业微信客户ID。"`         // ExternalUserID 为消息对应的企业微信客户ID。
	SendStatus     string       `gorm:"type:varchar(30);not null;default:'';index;comment:渠道发送状态快照，如 sent、failed。"` // SendStatus 为渠道发送状态快照，如 sent、failed。
	FailReason     string       `gorm:"type:text;comment:渠道发送失败原因或补偿说明。"`                                           // FailReason 为渠道发送失败原因或补偿说明。
	RawPayload     string       `gorm:"type:text;comment:企业微信原始消息JSON。"`                                            // RawPayload 为企业微信原始消息JSON。
	Status         enums.Status `gorm:"type:int;not null;default:0;index;comment:消息映射状态。"`                          // Status 为消息映射状态。
	AuditFields
}

// ChannelMessageOutbox 外部渠道消息投递任务表。
//
//	用于记录平台消息提交后的渠道发送任务，保证第三方发送动作与主事务解耦。
type ChannelMessageOutbox struct {
	ID             int64      `gorm:"primaryKey;autoIncrement;comment:投递任务主键。"`                                                              // ID 为投递任务主键。
	ChannelType    string     `gorm:"type:varchar(30);not null;default:'';index;uniqueIndex:uk_channel_message;comment:目标渠道类型，如 wxwork_kf。"` // ChannelType 为目标渠道类型，如 wxwork_kf。
	ConversationID int64      `gorm:"type:bigint;not null;default:0;index;comment:所属平台会话ID。"`                                                // ConversationID 为所属平台会话ID。
	MessageID      int64      `gorm:"type:bigint;not null;default:0;uniqueIndex:uk_channel_message;comment:待投递的平台消息ID。"`                     // MessageID 为待投递的平台消息ID。
	Payload        string     `gorm:"type:text;comment:渠道发送所需的标准化请求数据JSON。"`                                                                 // Payload 为渠道发送所需的标准化请求数据JSON。
	SendStatus     string     `gorm:"type:varchar(30);not null;default:'';index;comment:当前投递状态，如 pending、sending、sent、failed。"`              // SendStatus 为当前投递状态，如 pending、sending、sent、failed。
	RetryCount     int        `gorm:"type:int;not null;default:0;comment:已重试次数。"`                                                            // RetryCount 为已重试次数。
	NextRetryAt    *time.Time `gorm:"index;comment:下一次允许重试时间。"`                                                                              // NextRetryAt 为下一次允许重试时间。
	LastError      string     `gorm:"type:text;comment:最近一次发送失败信息。"`                                                                         // LastError 为最近一次发送失败信息。
	SentAt         *time.Time `gorm:"index;comment:最终发送成功时间。"`                                                                               // SentAt 为最终发送成功时间。
	AuditFields
}

// ConversationAssignment 会话接待关系。
type ConversationAssignment struct {
	ID             int64                    `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConversationID int64                    `gorm:"type:bigint;not null;index;comment:会话 ID"`
	FromUserID     int64                    `gorm:"type:bigint;not null;default:0;index;comment:转出用户 ID"`
	ToUserID       int64                    `gorm:"type:bigint;not null;default:0;index;comment:转入用户 ID"`
	AssignType     string                   `gorm:"type:varchar(30);not null;default:'';index;comment:分配类型"`
	Reason         string                   `gorm:"type:varchar(255);not null;default:'';comment:原因"`
	Status         enums.IMAssignmentStatus `gorm:"type:int;not null;index;comment:状态"`
	CreatedAt      time.Time                `gorm:"not null;index;comment:创建时间"`
	FinishedAt     *time.Time               `gorm:"comment:结束时间"`
	OperatorID     int64                    `gorm:"type:bigint;not null;default:0;index;comment:操作人 ID"`
}

// ConversationTag 会话标签关联
type ConversationTag struct {
	ID             int64 `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConversationID int64 `gorm:"type:bigint;not null;index;uniqueIndex:uk_conversation_tag;comment:会话 ID"`
	TagID          int64 `gorm:"type:bigint;not null;index;uniqueIndex:uk_conversation_tag;comment:标签 ID"`
	AuditFields
}

// QuickReply 快捷回复。
type QuickReply struct {
	ID        int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	GroupName string       `gorm:"type:varchar(50);not null;default:'';index;comment:分组名称"`
	Title     string       `gorm:"type:varchar(100);not null;default:'';index;comment:标题"`
	Content   string       `gorm:"type:text;comment:内容"`
	Status    enums.Status `gorm:"type:int;not null;index;comment:状态"`
	SortNo    int          `gorm:"type:int;not null;index;comment:排序号"`
	AuditFields
}

// AIAgent AI 接待实例。
type AIAgent struct {
	ID                     int64                           `gorm:"primaryKey;autoIncrement;comment:AI Agent 主键。"`                                       // ID 为 AI Agent 主键。
	Name                   string                          `gorm:"type:varchar(100);not null;default:'';index;comment:AI Agent 名称。"`                    // Name 为 AI Agent 名称。
	Description            string                          `gorm:"type:varchar(255);not null;default:'';comment:AI Agent 描述。"`                          // Description 为 AI Agent 描述。
	Status                 enums.Status                    `gorm:"type:int;not null;index;comment:AI Agent"`                                            // Status 为 AI Agent
	AIConfigID             int64                           `gorm:"type:bigint;not null;default:0;index;comment:关联的 AI 配置ID。"`                           // AIConfigID 为关联的 AI 配置ID。
	MaxSteps               int                             `gorm:"type:int;not null;default:6;comment:一次 Agent Loop 允许的最大推理步骤数。"`                       // MaxSteps 为一次 Agent Loop 允许的最大推理步骤数。
	ContextWindow          int                             `gorm:"type:int;not null;default:0;comment:会话上下文消息窗口，0 表示使用运行时默认值。"`                         // ContextWindow 为会话上下文消息窗口，0 表示使用运行时默认值。
	ToolPolicy             string                          `gorm:"type:text;comment:工具风险与确认策略JSON。"`                                                    // ToolPolicy 为工具风险与确认策略JSON。
	KnowledgePolicy        string                          `gorm:"type:text;comment:知识检索与无依据回答策略JSON。"`                                                 // KnowledgePolicy 为知识检索与无依据回答策略JSON。
	ServiceMode            enums.IMConversationServiceMode `gorm:"type:int;not null;default:3;index;comment:服务模式，如仅AI、仅人工、AI优先人工接管。"`                   // ServiceMode 为服务模式，如仅AI、仅人工、AI优先人工接管。
	SystemPrompt           string                          `gorm:"type:text;comment:该 Agent 的系统提示词。"`                                                   // SystemPrompt 为该 Agent 的系统提示词。
	WelcomeMessage         string                          `gorm:"type:text;comment:该 Agent 的欢迎语或首响模板。"`                                                // WelcomeMessage 为该 Agent 的欢迎语或首响模板。
	ReplyTimeoutSeconds    int                             `gorm:"type:int;not null;default:180;comment:异步自动回复超时秒数。"`                                   // ReplyTimeoutSeconds 为异步自动回复超时秒数。
	RolloutPercent         int                             `gorm:"type:int;not null;default:100;comment:该 Agent 的会话灰度百分比，100 表示全量。"`                    // RolloutPercent 为该 Agent 的会话灰度百分比，100 表示全量。
	PreviousRolloutPercent int                             `gorm:"type:int;not null;default:0;comment:PreviousRolloutPercent 保存上一次生效的灰度比例，0 表示尚无可回滚值。"` // PreviousRolloutPercent 保存上一次生效的灰度比例，0 表示尚无可回滚值。
	TeamIDs                string                          `gorm:"type:varchar(500);not null;default:'';comment:转人工时可路由的客服组ID列表，多个之间使用逗号分隔。"`           // TeamIDs 为转人工时可路由的客服组ID列表，多个之间使用逗号分隔。
	HandoffMode            enums.AIAgentHandoffMode        `gorm:"type:int;not null;default:1;comment:转人工执行方式，如进入待接入池、进入默认客服组待接入池。"`                    // HandoffMode 为转人工执行方式，如进入待接入池、进入默认客服组待接入池。
	FallbackMode           enums.AIAgentFallbackMode       `gorm:"type:int;not null;default:1;comment:知识不足时的回复策略。"`                                     // FallbackMode 为知识不足时的回复策略。
	FallbackMessage        string                          `gorm:"type:text;comment:知识不足回复文案。"`                                                         // FallbackMessage 为知识不足回复文案。
	KnowledgeIDs           string                          `gorm:"type:varchar(500);not null;default:'';comment:绑定的知识库ID列表，按顺序表示优先级。"`                  // KnowledgeIDs 为绑定的知识库ID列表，按顺序表示优先级。
	SkillIDs               string                          `gorm:"type:varchar(500);not null;default:'';comment:绑定的技能ID列表，按顺序表示允许路由的范围。"`               // SkillIDs 为绑定的技能ID列表，按顺序表示允许路由的范围。
	AllowedMCPTools        string                          `gorm:"type:text;comment:Agent 允许调用的 MCP 工具白名单配置 JSON。"`                                     // AllowedMCPTools 为 Agent 允许调用的 MCP 工具白名单配置 JSON。
	PublishedRevisionID    int64                           `gorm:"type:bigint;not null;default:0;index;comment:当前已发布 Agent 配置快照ID。"`                    // PublishedRevisionID 为当前已发布 Agent 配置快照ID。
	SortNo                 int                             `gorm:"type:int;not null;default:0;index;comment:后台展示排序号。"`                                  // SortNo 为后台展示排序号。
	AuditFields
}

// AgentRevision stores an immutable published Agent configuration snapshot.
type AgentRevision struct {
	ID              int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	AgentID         int64        `gorm:"type:bigint;not null;index;uniqueIndex:uk_agent_revision;comment:客服或 Agent ID"`
	Revision        int          `gorm:"type:int;not null;uniqueIndex:uk_agent_revision;comment:修订版本号"`
	Status          enums.Status `gorm:"type:int;not null;default:0;index;comment:状态"`
	Definition      string       `gorm:"type:text;comment:配置定义"`
	DefinitionHash  string       `gorm:"type:varchar(64);not null;default:'';index;comment:配置摘要哈希"`
	PublishedAt     *time.Time   `gorm:"index;comment:发布时间"`
	PublishedByID   int64        `gorm:"type:bigint;not null;default:0;index;comment:发布人 ID"`
	PublishedByName string       `gorm:"type:varchar(100);not null;default:'';comment:发布人名称"`
	AuditFields
}

// AgentRun is the parent audit record for one Agent Loop execution.
type AgentRun struct {
	ID               int64      `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConversationID   int64      `gorm:"type:bigint;not null;default:0;index;comment:会话 ID"`
	AIAgentID        int64      `gorm:"type:bigint;not null;default:0;index;comment:AI Agent ID"`
	AgentRevisionID  int64      `gorm:"type:bigint;not null;default:0;index;comment:Agent 发布版本 ID"`
	SourceMessageID  int64      `gorm:"type:bigint;not null;default:0;index;comment:来源消息 ID"`
	WorkflowRunID    int64      `gorm:"type:bigint;not null;default:0;index;comment:工作流运行 ID"`
	Status           string     `gorm:"type:varchar(30);not null;default:'';index;comment:状态"`
	PromptTokens     int        `gorm:"type:int;not null;default:0;comment:输入 Token 数"`
	CompletionTokens int        `gorm:"type:int;not null;default:0;comment:输出 Token 数"`
	StartedAt        time.Time  `gorm:"not null;index;comment:开始时间"`
	EndedAt          *time.Time `gorm:"index;comment:结束时间"`
	ErrorMessage     string     `gorm:"type:text;comment:错误信息"`
	TraceData        string     `gorm:"type:text;comment:运行轨迹数据"`
	CreatedAt        time.Time  `gorm:"not null;index;comment:创建时间"`
	UpdatedAt        time.Time  `gorm:"not null;index;comment:更新时间"`
}

// AgentStep records a normalized model, tool, workflow, or policy transition.
type AgentStep struct {
	ID            int64      `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	AgentRunID    int64      `gorm:"type:bigint;not null;index;comment:Agent 运行 ID"`
	WorkflowRunID int64      `gorm:"type:bigint;not null;default:0;index;comment:工作流运行 ID"`
	StepType      string     `gorm:"type:varchar(50);not null;default:'';index;comment:步骤类型"`
	StepCode      string     `gorm:"type:varchar(100);not null;default:'';index;comment:步骤编码"`
	Status        string     `gorm:"type:varchar(30);not null;default:'';index;comment:状态"`
	InputPreview  string     `gorm:"type:text;comment:输入预览"`
	OutputPreview string     `gorm:"type:text;comment:输出预览"`
	ErrorMessage  string     `gorm:"type:text;comment:错误信息"`
	StartedAt     time.Time  `gorm:"not null;index;comment:开始时间"`
	EndedAt       *time.Time `gorm:"index;comment:结束时间"`
	DurationMS    int        `gorm:"type:int;not null;default:0;comment:耗时毫秒数"`
	CreatedAt     time.Time  `gorm:"not null;index;comment:创建时间"`
}

// AgentToolCall records the safety-relevant details of a normalized tool call.
type AgentToolCall struct {
	ID               int64     `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	AgentRunID       int64     `gorm:"type:bigint;not null;index;comment:Agent 运行 ID"`
	AgentStepID      int64     `gorm:"type:bigint;not null;default:0;index;comment:Agent 步骤 ID"`
	ToolCode         string    `gorm:"type:varchar(150);not null;default:'';index;comment:工具编码"`
	RiskLevel        string    `gorm:"type:varchar(30);not null;default:'';index;comment:风险等级"`
	RequireConfirm   bool      `gorm:"not null;default:false;comment:是否需要确认"`
	Status           string    `gorm:"type:varchar(30);not null;default:'';index;comment:状态"`
	ArgumentsPreview string    `gorm:"type:text;comment:参数预览"`
	ResultPreview    string    `gorm:"type:text;comment:结果预览"`
	ErrorMessage     string    `gorm:"type:text;comment:错误信息"`
	DurationMS       int       `gorm:"type:int;not null;default:0;comment:耗时毫秒数"`
	CreatedAt        time.Time `gorm:"not null;index;comment:创建时间"`
}

// AgentRunQualityFeedback is an operator-provided quality review for one
// normalized Agent run. Runtime completion must not be treated as resolution.
type AgentRunQualityFeedback struct {
	ID               int64                          `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	AgentRunID       int64                          `gorm:"type:bigint;not null;uniqueIndex;comment:Agent 运行 ID"`
	ResolutionStatus enums.AgentRunResolutionStatus `gorm:"type:varchar(20);not null;default:'unknown';index;comment:问题解决状态"`
	EvidenceStatus   enums.AgentRunEvidenceStatus   `gorm:"type:varchar(20);not null;default:'unknown';index;comment:证据状态"`
	Comment          string                         `gorm:"type:text;comment:评论或备注"`
	AuditFields
}

// AIWorkflow 表示客服 AI Agent 可编辑会话流程主表。
type AIWorkflow struct {
	ID                 int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	Name               string       `gorm:"type:varchar(100);not null;default:'';index;comment:名称"`
	Description        string       `gorm:"type:text;comment:描述"`
	Status             enums.Status `gorm:"type:int;not null;default:0;index;comment:状态"`
	DraftDefinition    string       `gorm:"type:text;comment:草稿定义"`
	PublishedVersionID int64        `gorm:"type:bigint;not null;default:0;index;comment:已发布版本 ID"`
	SortNo             int          `gorm:"type:int;not null;default:0;index;comment:排序号"`
	AuditFields
}

// AIWorkflowVersion 表示 AI 会话流程的不可变发布版本。
type AIWorkflowVersion struct {
	ID              int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	WorkflowID      int64        `gorm:"type:bigint;not null;default:0;index;comment:工作流 ID"`
	Version         int          `gorm:"type:int;not null;default:0;index;comment:版本号"`
	Status          enums.Status `gorm:"type:int;not null;default:0;index;comment:状态"`
	Definition      string       `gorm:"type:text;comment:配置定义"`
	DefinitionHash  string       `gorm:"type:varchar(64);not null;default:'';index;comment:配置摘要哈希"`
	PublishedAt     *time.Time   `gorm:"index;comment:发布时间"`
	PublishedByID   int64        `gorm:"type:bigint;not null;default:0;index;comment:发布人 ID"`
	PublishedByName string       `gorm:"type:varchar(100);not null;default:'';comment:发布人名称"`
	AuditFields
}

// AIAgentWorkflowBinding binds an Agent to one immutable, published workflow
// version. Workflows are maintained independently and can be reused by many
// Agents; Agent publication snapshots the binding set for reproducible runs.
type AIAgentWorkflowBinding struct {
	ID                 int64  `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	AIAgentID          int64  `gorm:"type:bigint;not null;index;uniqueIndex:uk_agent_workflow_binding;comment:AI Agent ID"`
	WorkflowID         int64  `gorm:"type:bigint;not null;index;comment:工作流 ID"`
	WorkflowVersionID  int64  `gorm:"type:bigint;not null;index;uniqueIndex:uk_agent_workflow_binding;comment:工作流版本 ID"`
	ToolName           string `gorm:"type:varchar(100);not null;default:'';comment:工具名称"`
	TriggerInstruction string `gorm:"type:text;comment:触发指令"`
	Priority           int    `gorm:"type:int;not null;default:0;index;comment:优先级"`
	Enabled            bool   `gorm:"not null;default:true;index;comment:是否启用"`
	AuditFields
}

// AIWorkflowRun 表示一次会话 workflow 执行记录。
type AIWorkflowRun struct {
	ID                int64      `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	WorkflowID        int64      `gorm:"type:bigint;not null;default:0;index;comment:工作流 ID"`
	WorkflowVersionID int64      `gorm:"type:bigint;not null;default:0;index;comment:工作流版本 ID"`
	ConversationID    int64      `gorm:"type:bigint;not null;default:0;index;comment:会话 ID"`
	AIAgentID         int64      `gorm:"type:bigint;not null;default:0;index;comment:AI Agent ID"`
	MessageID         int64      `gorm:"type:bigint;not null;default:0;index;comment:消息 ID"`
	Status            int        `gorm:"type:int;not null;default:0;index;comment:状态"`
	StartedAt         time.Time  `gorm:"not null;index;comment:开始时间"`
	EndedAt           *time.Time `gorm:"index;comment:结束时间"`
	InterruptType     string     `gorm:"type:varchar(50);not null;default:'';index;comment:中断类型"`
	InterruptNodeID   string     `gorm:"type:varchar(100);not null;default:'';index;comment:中断节点 ID"`
	ErrorMessage      string     `gorm:"type:text;comment:错误信息"`
	AuditFields
}

// AIWorkflowNodeRun 表示 workflow 执行中的单节点审计记录。
type AIWorkflowNodeRun struct {
	ID            int64      `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	WorkflowRunID int64      `gorm:"type:bigint;not null;default:0;index;comment:工作流运行 ID"`
	NodeID        string     `gorm:"type:varchar(100);not null;default:'';index;comment:节点 ID"`
	NodeType      string     `gorm:"type:varchar(50);not null;default:'';index;comment:节点类型"`
	Status        int        `gorm:"type:int;not null;default:0;index;comment:状态"`
	InputPreview  string     `gorm:"type:text;comment:输入预览"`
	OutputPreview string     `gorm:"type:text;comment:输出预览"`
	ErrorMessage  string     `gorm:"type:text;comment:错误信息"`
	StartedAt     time.Time  `gorm:"not null;index;comment:开始时间"`
	EndedAt       *time.Time `gorm:"index;comment:结束时间"`
	DurationMS    int        `gorm:"type:int;not null;default:0;comment:耗时毫秒数"`
}

// Channel 接入渠道配置。
//
//	用于统一描述系统的外部接入入口。不同渠道类型共享统一的接入配置骨架，
//	例如网页客服渠道（web）和企业微信客服渠道（wxwork_kf）。
//	渠道本身负责定义“入口如何识别、默认接入哪个 AI Agent、渠道专属配置是什么”，
//	而具体消息收发、会话映射等运行时数据由各自的渠道业务表承载。
type Channel struct {
	ID                            int64  `gorm:"primaryKey;autoIncrement;comment:渠道主键。"`                                                                                                 // ID 为渠道主键。
	Name                          string `gorm:"type:varchar(100);not null;default:'';index;comment:渠道名称，用于后台展示和业务识别，例如“官网客服”“企业微信主客服”。"`                                                // Name 为渠道名称，用于后台展示和业务识别，例如“官网客服”“企业微信主客服”。
	ChannelType                   string `gorm:"type:varchar(30);not null;default:'';index;comment:渠道类型，决定该渠道的接入方式和配置解释规则。当前规划的典型取值包括：web、wxwork_kf。"`                                   // ChannelType 为渠道类型，决定该渠道的接入方式和配置解释规则。当前规划的典型取值包括：web、wxwork_kf。
	ChannelID                     string `gorm:"type:varchar(64);not null;default:'';uniqueIndex;comment:渠道入口标识，由系统自动生成。对 web 渠道，该字段用于前端通过 X-Channel-Id 标识接入来源；对其他渠道，作为统一的系统内稳定渠道标识保留。"` // ChannelID 为渠道入口标识，由系统自动生成。对 web 渠道，该字段用于前端通过 X-Channel-Id 标识接入来源；对其他渠道，作为统一的系统内稳定渠道标识保留。
	AIAgentID                     int64  `gorm:"type:bigint;not null;default:0;index;comment:该渠道默认接入的 AI Agent。 当外部客户通过该渠道首次进入系统且尚未命中现有未结束会话时，系统会使用该 AI Agent 作为会话默认接待实例。"`              // AIAgentID 为该渠道默认接入的 AI Agent。 当外部客户通过该渠道首次进入系统且尚未命中现有未结束会话时，系统会使用该 AI Agent 作为会话默认接待实例。
	AIAgentRolloutPercent         int    `gorm:"type:int;not null;default:100;comment:该渠道对 AI 自动回复施加的灰度百分比，100 表示不额外限制。"`                                                                // AIAgentRolloutPercent 为该渠道对 AI 自动回复施加的灰度百分比，100 表示不额外限制。
	PreviousAIAgentRolloutPercent int    `gorm:"type:int;not null;default:0;comment:PreviousAIAgentRolloutPercent 保存渠道上一次生效的 Agent 灰度比例，0 表示尚无可回滚值。"`                                    // PreviousAIAgentRolloutPercent 保存渠道上一次生效的 Agent 灰度比例，0 表示尚无可回滚值。
	// ConfigJSON 为渠道专属扩展配置，使用 JSON 存储。
	// 例如：
	// 1. web 渠道可记录允许域名、品牌配置等；
	// 2. wxwork_kf 渠道可记录 openKfId、欢迎语策略等。
	// 该字段只存储渠道类型私有配置，不承载通用主字段。
	ConfigJSON string       `gorm:"type:text;comment:配置 JSON"`
	Status     enums.Status `gorm:"type:int;not null;default:0;index;comment:渠道状态。禁用后，该渠道不再允许新会话接入；删除时采用软删除状态保留历史关联数据。"` // Status 为渠道状态。禁用后，该渠道不再允许新会话接入；删除时采用软删除状态保留历史关联数据。
	Remark     string       `gorm:"type:text;comment:渠道备注，用于记录接入说明、维护说明和内部运维信息。"`                                        // Remark 为渠道备注，用于记录接入说明、维护说明和内部运维信息。
	AuditFields
}

// ConversationEventLog 会话事件日志。
type ConversationEventLog struct {
	ID             int64              `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConversationID int64              `gorm:"type:bigint;not null;index;comment:会话 ID"`
	RequestID      string             `gorm:"type:varchar(128);not null;default:'';index;comment:请求 ID"`
	EventType      enums.IMEventType  `gorm:"type:varchar(50);not null;default:'';index;comment:事件类型"`
	OperatorType   enums.IMSenderType `gorm:"type:varchar(30);not null;default:'';index;comment:操作人类型"`
	OperatorID     int64              `gorm:"type:bigint;not null;default:0;index;comment:操作人 ID"`
	Content        string             `gorm:"type:text;comment:内容"`
	Payload        string             `gorm:"type:text;comment:消息载荷"`
	CreatedAt      time.Time          `gorm:"not null;index;comment:创建时间"`
}

// Ticket 客服问题记录。
type Ticket struct {
	ID                int64              `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	TicketNo          string             `gorm:"type:varchar(64);not null;default:'';uniqueIndex;comment:工单编号"`
	Title             string             `gorm:"type:varchar(255);not null;default:'';index;comment:标题"`
	Description       string             `gorm:"type:text;comment:描述"`
	Source            enums.TicketSource `gorm:"type:varchar(50);not null;default:'';index;comment:来源"`
	Channel           string             `gorm:"type:varchar(50);not null;default:'';index;comment:渠道"`
	CustomerID        int64              `gorm:"type:bigint;not null;default:0;index;comment:客户 ID"`
	ConversationID    int64              `gorm:"type:bigint;not null;default:0;index;comment:会话 ID"`
	Status            enums.TicketStatus `gorm:"type:varchar(50);not null;default:'pending';index;comment:状态"`
	CurrentAssigneeID int64              `gorm:"type:bigint;not null;default:0;index;comment:当前处理人 ID"`
	HandledAt         *time.Time         `gorm:"index;comment:处理时间"`
	AuditFields
}

// TicketTag 工单标签关联。
type TicketTag struct {
	ID       int64 `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	TicketID int64 `gorm:"type:bigint;not null;index;uniqueIndex:uk_ticket_tag;comment:工单 ID"`
	TagID    int64 `gorm:"type:bigint;not null;index;uniqueIndex:uk_ticket_tag;comment:标签 ID"`
	AuditFields
}

// TicketProgress 工单处理进展。
type TicketProgress struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	TicketID  int64     `gorm:"type:bigint;not null;index;comment:工单 ID"`
	Content   string    `gorm:"type:text;comment:内容"`
	AuthorID  int64     `gorm:"type:bigint;not null;default:0;index;comment:作者 ID"`
	CreatedAt time.Time `gorm:"not null;index;comment:创建时间"`
}

// AgentProfile 客服档案。
type AgentProfile struct {
	ID                    int64               `gorm:"primaryKey;autoIncrement;comment:客服档案主键。"`                                 // ID 为客服档案主键。
	UserID                int64               `gorm:"type:bigint;not null;uniqueIndex;comment:UserID 关联后台用户，一名用户只允许一份客服档案。"`    // UserID 关联后台用户，一名用户只允许一份客服档案。
	TeamID                int64               `gorm:"type:bigint;not null;default:0;index;comment:客服所属客服组。"`                    // TeamID 为客服所属客服组。
	AgentCode             string              `gorm:"type:varchar(64);not null;default:'';uniqueIndex;comment:客服工号，用于业务侧识别客服。"` // AgentCode 为客服工号，用于业务侧识别客服。
	DisplayName           string              `gorm:"type:varchar(100);not null;default:'';index;comment:客服展示名，可区别于后台昵称。"`      // DisplayName 为客服展示名，可区别于后台昵称。
	Avatar                string              `gorm:"type:varchar(1024);not null;default:'';comment:客服头像 URL。"`                 // Avatar 为客服头像 URL。
	ServiceStatus         enums.ServiceStatus `gorm:"type:int;not null;default:0;index;comment:客服服务状态：0空闲 1忙碌。"`                // ServiceStatus 表示客服服务状态：0空闲 1忙碌。
	MaxConcurrentCount    int                 `gorm:"type:int;not null;default:0;comment:客服最大并发接待数。"`                           // MaxConcurrentCount 表示客服最大并发接待数。
	PriorityLevel         int                 `gorm:"type:int;not null;default:0;index;comment:自动分配优先级，值越大越优先。"`                // PriorityLevel 表示自动分配优先级，值越大越优先。
	AutoAssignEnabled     bool                `gorm:"not null;default:true;index;comment:是否参与自动分配。"`                            // AutoAssignEnabled 表示是否参与自动分配。
	ReceiveOfflineMessage bool                `gorm:"not null;default:false;comment:离线时是否仍接收离线消息或转接消息。"`                        // ReceiveOfflineMessage 表示离线时是否仍接收离线消息或转接消息。
	LastOnlineAt          *time.Time          `gorm:"index;comment:LastOnlineAt 记录最近一次在线时间。"`                                   // LastOnlineAt 记录最近一次在线时间。
	LastStatusAt          *time.Time          `gorm:"index;comment:LastStatusAt 记录最近一次状态变更时间。"`                                 // LastStatusAt 记录最近一次状态变更时间。
	Status                enums.Status        `gorm:"type:int;not null;default:0;index;comment:客服档案状态"`                         // Status 表示客服档案状态
	Remark                string              `gorm:"type:text;comment:Remark 记录客服备注信息。"`                                       // Remark 记录客服备注信息。
	AuditFields
}

// AgentTeam 客服组。
type AgentTeam struct {
	ID           int64        `gorm:"primaryKey;autoIncrement;comment:客服组主键。"`                       // ID 为客服组主键。
	Name         string       `gorm:"type:varchar(100);not null;default:'';index;comment:客服组名称。"`    // Name 为客服组名称。
	LeaderUserID int64        `gorm:"type:bigint;not null;default:0;index;comment:组长用户ID，0 表示暂未设置。"` // LeaderUserID 为组长用户ID，0 表示暂未设置。
	Status       enums.Status `gorm:"type:int;not null;default:0;index;comment:客服组状态"`               // Status 表示客服组状态
	Description  string       `gorm:"type:varchar(255);not null;default:'';comment:客服组简介，用于说明职责边界。"` // Description 为客服组简介，用于说明职责边界。
	Remark       string       `gorm:"type:text;comment:Remark 记录客服组内部备注。"`                           // Remark 记录客服组内部备注。
	AuditFields
}

// AgentTeamSchedule 客服组排班。
type AgentTeamSchedule struct {
	ID      int64        `gorm:"primaryKey;autoIncrement;comment:组排班主键。"`                      // ID 为组排班主键。
	TeamID  int64        `gorm:"type:bigint;not null;index;comment:被排班的客服组ID。"`                // TeamID 为被排班的客服组ID。
	StartAt time.Time    `gorm:"not null;index;comment:班次开始时间。"`                               // StartAt 为班次开始时间。
	EndAt   time.Time    `gorm:"not null;index;comment:班次结束时间。"`                               // EndAt 为班次结束时间。
	Remark  string       `gorm:"type:varchar(255);not null;default:'';comment:Remark 记录排班备注。"` // Remark 记录排班备注。
	Status  enums.Status `gorm:"type:int;not null;default:0;index;comment:组排班记录状态。"`           // Status 表示组排班记录状态。
	AuditFields
}

// AIConfig AI 统一配置。
// 一条记录表示一个可直接调用的 AI 配置实例，
// 同时包含厂商接入信息、模型信息和调用参数，不再拆分 endpoint/model 两层概念。
type AIConfig struct {
	ID               int64             `gorm:"primaryKey;autoIncrement;comment:配置主键。"`                                                     // ID 为配置主键。
	Name             string            `gorm:"type:varchar(100);not null;default:'';index;comment:配置名称，用于后台识别和展示。"`                        // Name 为配置名称，用于后台识别和展示。
	Provider         enums.AIProvider  `gorm:"type:varchar(50);not null;default:'';index;comment:供应商标识，例如 openai、azure_openai、dashscope。"` // Provider 为供应商标识，例如 openai、azure_openai、dashscope。
	BaseURL          string            `gorm:"type:varchar(255);not null;default:'';comment:模型服务基础地址，例如 https://api.openai.com/v1。"`       // BaseURL 为模型服务基础地址，例如 https://api.openai.com/v1。
	APIKey           string            `gorm:"type:varchar(255);not null;default:'';comment:服务端请求模型接口所需密钥。"`                               // APIKey 为服务端请求模型接口所需密钥。
	ModelType        enums.AIModelType `gorm:"type:varchar(30);not null;default:'';index;comment:模型类型，例如 llm、embedding、rerank。"`           // ModelType 为模型类型，例如 llm、embedding、rerank。
	ModelName        string            `gorm:"type:varchar(100);not null;default:'';index;comment:实际请求时传给上游的模型名。"`                         // ModelName 为实际请求时传给上游的模型名。
	Dimension        int               `gorm:"type:int;not null;default:0;comment:向量维度，仅 embedding 模型通常需要填写。"`                             // Dimension 为向量维度，仅 embedding 模型通常需要填写。
	MaxContextTokens int               `gorm:"type:int;not null;default:0;comment:模型支持的最大上下文 token 数。"`                                    // MaxContextTokens 为模型支持的最大上下文 token 数。
	MaxOutputTokens  int               `gorm:"type:int;not null;default:0;comment:模型建议的最大输出 token 数。"`                                     // MaxOutputTokens 为模型建议的最大输出 token 数。
	TimeoutMS        int               `gorm:"type:int;not null;default:30000;comment:调用该配置的默认超时时间，单位毫秒。"`                                 // TimeoutMS 为调用该配置的默认超时时间，单位毫秒。
	MaxRetryCount    int               `gorm:"type:int;not null;default:0;comment:默认最大重试次数。"`                                              // MaxRetryCount 为默认最大重试次数。
	RPMLimit         int               `gorm:"type:int;not null;default:0;comment:每分钟请求数限制，0 表示未显式配置。"`                                    // RPMLimit 为每分钟请求数限制，0 表示未显式配置。
	TPMLimit         int               `gorm:"type:int;not null;default:0;comment:每分钟 token 数限制，0 表示未显式配置。"`                               // TPMLimit 为每分钟 token 数限制，0 表示未显式配置。
	Status           enums.Status      `gorm:"type:int;not null;index;comment:Status 状态；同一 modelType 仅允许一条启用记录。"`                          // Status 状态；同一 modelType 仅允许一条启用记录。
	SortNo           int               `gorm:"type:int;not null;index;comment:排序号，用于后台展示和人工调整顺序。"`                                         // SortNo 为排序号，用于后台展示和人工调整顺序。
	Remark           string            `gorm:"type:text;comment:备注，用于记录用途、成本、限制和切换说明等补充信息。"`                                               // Remark 为备注，用于记录用途、成本、限制和切换说明等补充信息。
	AuditFields
}

// KnowledgeBase 知识库主表。
type KnowledgeBase struct {
	ID                    int64        `gorm:"primaryKey;autoIncrement;comment:知识库主键。"`                                        // ID 为知识库主键。
	Name                  string       `gorm:"type:varchar(100);not null;default:'';index;comment:知识库名称。"`                     // Name 为知识库名称。
	Description           string       `gorm:"type:text;comment:知识库描述。"`                                                       // Description 为知识库描述。
	KnowledgeType         string       `gorm:"type:varchar(20);not null;default:'document';index;comment:知识库类型：document/faq。"` // KnowledgeType 为知识库类型：document/faq。
	Status                enums.Status `gorm:"type:int;not null;index;comment:状态"`                                             // Status 为状态
	DefaultTopK           int          `gorm:"type:int;not null;default:10;comment:默认召回数量。"`                                   // DefaultTopK 为默认召回数量。
	DefaultScoreThreshold float64      `gorm:"type:decimal(5,4);not null;default:0.5;comment:默认相似度阈值。"`                        // DefaultScoreThreshold 为默认相似度阈值。
	DefaultRerankLimit    int          `gorm:"type:int;not null;default:5;comment:默认重排后保留数量。"`                                 // DefaultRerankLimit 为默认重排后保留数量。
	ChunkProvider         string       `gorm:"type:varchar(30);not null;default:'structured';comment:知识库分块策略 provider。"`       // ChunkProvider 为知识库分块策略 provider。
	ChunkTargetTokens     int          `gorm:"type:int;not null;default:300;comment:目标 chunk token 数。"`                        // ChunkTargetTokens 为目标 chunk token 数。
	ChunkMaxTokens        int          `gorm:"type:int;not null;default:400;comment:单 chunk 最大 token 数。"`                      // ChunkMaxTokens 为单 chunk 最大 token 数。
	ChunkOverlapTokens    int          `gorm:"type:int;not null;default:40;comment:相邻 chunk 重叠 token 数。"`                      // ChunkOverlapTokens 为相邻 chunk 重叠 token 数。
	AnswerMode            int          `gorm:"type:int;not null;default:1;comment:回答模式：1严格知识库模式 2辅助解释模式。"`                     // AnswerMode 为回答模式：1严格知识库模式 2辅助解释模式。
	SortNo                int          `gorm:"type:int;not null;default:0;index;comment:排序号，用于后台展示和知识库的人工排序管理。"`               // SortNo 为排序号，用于后台展示和知识库的人工排序管理。
	Remark                string       `gorm:"type:text;comment:备注。"`                                                          // Remark 为备注。
	AuditFields
}

// KnowledgeDirectory 知识库内部目录表。
type KnowledgeDirectory struct {
	ID              int64        `gorm:"primaryKey;autoIncrement;comment:目录主键。"`                                           // ID 为目录主键。
	KnowledgeBaseID int64        `gorm:"type:bigint;not null;index:idx_kb_parent_sort;comment:所属知识库 ID。"`                  // KnowledgeBaseID 为所属知识库 ID。
	ParentID        int64        `gorm:"type:bigint;not null;default:0;index:idx_kb_parent_sort;comment:父目录 ID，0 表示一级目录。"` // ParentID 为父目录 ID，0 表示一级目录。
	Name            string       `gorm:"type:varchar(100);not null;default:'';index;comment:目录名称。"`                        // Name 为目录名称。
	SortNo          int          `gorm:"type:int;not null;default:0;index:idx_kb_parent_sort;comment:同级排序号。"`              // SortNo 为同级排序号。
	Status          enums.Status `gorm:"type:int;not null;default:0;index;comment:状态。"`                                    // Status 为状态。
	Remark          string       `gorm:"type:text;comment:备注。"`                                                            // Remark 为备注。
	AuditFields
}

// KnowledgeDocument 知识文档主表。
type KnowledgeDocument struct {
	ID              int64                              `gorm:"primaryKey;autoIncrement;comment:文档主键。"`                                                 // ID 为文档主键。
	KnowledgeBaseID int64                              `gorm:"type:bigint;not null;index;comment:所属知识库ID。"`                                            // KnowledgeBaseID 为所属知识库ID。
	DirectoryID     int64                              `gorm:"type:bigint;not null;default:0;index;comment:所属知识库内部目录 ID，0 表示根目录。"`                     // DirectoryID 为所属知识库内部目录 ID，0 表示根目录。
	Title           string                             `gorm:"type:varchar(255);not null;default:'';index;comment:文档标题。"`                              // Title 为文档标题。
	ContentType     enums.KnowledgeDocumentContentType `gorm:"type:varchar(20);not null;default:'html';comment:内容类型：html/markdown。"`                   // ContentType 为内容类型：html/markdown。
	Content         string                             `gorm:"type:text;comment:文档内容。"`                                                                // Content 为文档内容。
	Status          enums.Status                       `gorm:"type:int;not null;default:0;index;comment:状态"`                                           // Status 为状态
	IndexStatus     enums.KnowledgeDocumentIndexStatus `gorm:"type:varchar(20);not null;default:'pending';index;comment:索引状态：pending/indexed/failed。"` // IndexStatus 为索引状态：pending/indexed/failed。
	IndexedAt       *time.Time                         `gorm:"index;comment:最近一次索引成功时间。"`                                                              // IndexedAt 为最近一次索引成功时间。
	IndexError      string                             `gorm:"type:text;comment:最近一次索引失败信息。"`                                                          // IndexError 为最近一次索引失败信息。
	ContentHash     string                             `gorm:"type:varchar(64);not null;default:'';index;comment:内容哈希，用于变更检测。"`                        // ContentHash 为内容哈希，用于变更检测。
	AuditFields
}

// KnowledgeFAQ FAQ 条目主表。
type KnowledgeFAQ struct {
	ID               int64                              `gorm:"primaryKey;autoIncrement;comment:FAQ 主键。"`                                               // ID 为 FAQ 主键。
	KnowledgeBaseID  int64                              `gorm:"type:bigint;not null;index;comment:所属 FAQ 知识库 ID。"`                                      // KnowledgeBaseID 为所属 FAQ 知识库 ID。
	DirectoryID      int64                              `gorm:"type:bigint;not null;default:0;index;comment:所属知识库内部目录 ID，0 表示根目录。"`                     // DirectoryID 为所属知识库内部目录 ID，0 表示根目录。
	Question         string                             `gorm:"type:varchar(500);not null;default:'';index;comment:标准问题。"`                              // Question 为标准问题。
	Answer           string                             `gorm:"type:text;comment:标准答案。"`                                                                // Answer 为标准答案。
	SimilarQuestions string                             `gorm:"type:text;comment:相似问 JSON 数组。"`                                                         // SimilarQuestions 为相似问 JSON 数组。
	Status           enums.Status                       `gorm:"type:int;not null;default:0;index;comment:状态。"`                                          // Status 为状态。
	IndexStatus      enums.KnowledgeDocumentIndexStatus `gorm:"type:varchar(20);not null;default:'pending';index;comment:索引状态：pending/indexed/failed。"` // IndexStatus 为索引状态：pending/indexed/failed。
	IndexedAt        *time.Time                         `gorm:"index;comment:最近一次索引成功时间。"`                                                              // IndexedAt 为最近一次索引成功时间。
	IndexError       string                             `gorm:"type:text;comment:最近一次索引失败信息。"`                                                          // IndexError 为最近一次索引失败信息。
	Remark           string                             `gorm:"type:text;comment:备注。"`                                                                  // Remark 为备注。
	AuditFields
}

// KnowledgeChunk 切片元数据表。
type KnowledgeChunk struct {
	ID              int64        `gorm:"primaryKey;autoIncrement;comment:切片主键。"`                             // ID 为切片主键。
	KnowledgeBaseID int64        `gorm:"type:bigint;not null;index;comment:知识库ID。"`                          // KnowledgeBaseID 为知识库ID。
	DocumentID      int64        `gorm:"type:bigint;not null;default:0;index;comment:文档ID。"`                 // DocumentID 为文档ID。
	FaqID           int64        `gorm:"type:bigint;not null;default:0;index;comment:FAQ ID。"`               // FaqID 为 FAQ ID。
	ChunkNo         int          `gorm:"type:int;not null;default:0;index;comment:切片序号。"`                    // ChunkNo 为切片序号。
	Title           string       `gorm:"type:varchar(255);not null;default:'';comment:切片标题。"`                // Title 为切片标题。
	Content         string       `gorm:"type:text;comment:切片内容。"`                                            // Content 为切片内容。
	ContentHash     string       `gorm:"type:varchar(64);not null;default:'';index;comment:内容哈希。"`           // ContentHash 为内容哈希。
	CharCount       int          `gorm:"type:int;not null;default:0;comment:字符数。"`                           // CharCount 为字符数。
	TokenCount      int          `gorm:"type:int;not null;default:0;comment:token数。"`                        // TokenCount 为token数。
	ChunkType       string       `gorm:"type:varchar(30);not null;default:'';comment:切片类型。"`                 // ChunkType 为切片类型。
	SectionPath     string       `gorm:"type:text;comment:章节路径。"`                                            // SectionPath 为章节路径。
	Provider        string       `gorm:"type:varchar(30);not null;default:'';comment:分块 provider。"`          // Provider 为分块 provider。
	Status          enums.Status `gorm:"type:int;not null;default:0;index;comment:状态：1有效 2已删除。"`             // Status 为状态：1有效 2已删除。
	VectorID        string       `gorm:"type:varchar(100);not null;default:'';index;comment:向量库中的point ID。"` // VectorID 为向量库中的point ID。
	CreatedAt       time.Time    `gorm:"not null;index;comment:创建时间"`
	UpdatedAt       time.Time    `gorm:"not null;index;comment:更新时间"`
}

// DocPage is both a document-center navigation node and a content page.
type DocPage struct {
	ID                        int64               `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ParentID                  int64               `gorm:"type:bigint;not null;default:0;index;uniqueIndex:idx_doc_page_parent_slug;comment:父级 ID"`
	Title                     string              `gorm:"type:varchar(255);not null;default:'';index;comment:标题"`
	Slug                      string              `gorm:"type:varchar(160);not null;default:'';uniqueIndex:idx_doc_page_parent_slug;comment:短链接标识"`
	Summary                   string              `gorm:"type:text;comment:摘要"`
	ContentType               string              `gorm:"type:varchar(20);not null;default:'markdown';comment:内容类型"`
	Content                   string              `gorm:"type:text;comment:内容"`
	CoverURL                  string              `gorm:"type:varchar(255);not null;default:'';comment:封面链接"`
	TagsJSON                  string              `gorm:"type:text;comment:标签 JSON"`
	Status                    enums.DocPageStatus `gorm:"type:varchar(20);not null;default:'draft';index;comment:状态"`
	SortNo                    int                 `gorm:"type:int;not null;default:0;index;comment:排序号"`
	ViewCount                 int64               `gorm:"type:bigint;not null;default:0;comment:浏览次数"`
	HelpfulCount              int64               `gorm:"type:bigint;not null;default:0;comment:有帮助次数"`
	UnhelpfulCount            int64               `gorm:"type:bigint;not null;default:0;comment:无帮助次数"`
	PublishedAt               *time.Time          `gorm:"index;comment:发布时间"`
	SyncedKnowledgeDocumentID int64               `gorm:"type:bigint;not null;default:0;index;comment:同步知识文档 ID"`
	Remark                    string              `gorm:"type:text;comment:备注"`
	AuditFields
}

// Category groups public FAQ community posts.
type Category struct {
	ID          int64        `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	Name        string       `gorm:"type:varchar(100);not null;default:'';index;comment:名称"`
	Slug        string       `gorm:"type:varchar(120);not null;default:'';uniqueIndex;comment:短链接标识"`
	Description string       `gorm:"type:text;comment:描述"`
	SortNo      int          `gorm:"type:int;not null;default:0;index;comment:排序号"`
	Status      enums.Status `gorm:"type:int;not null;default:0;index;comment:状态"`
	Remark      string       `gorm:"type:text;comment:备注"`
	AuditFields
}

// Post is a logged-in-customer post in the public FAQ community.
type Post struct {
	ID                  int64                   `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	CategoryID          int64                   `gorm:"type:bigint;not null;default:0;index;comment:分类 ID"`
	UserID              int64                   `gorm:"type:bigint;not null;index;comment:用户 ID"`
	Title               string                  `gorm:"type:varchar(255);not null;default:'';index;comment:标题"`
	ContentType         string                  `gorm:"type:varchar(20);not null;default:'markdown';comment:内容类型"`
	Content             string                  `gorm:"type:text;comment:内容"`
	TagsJSON            string                  `gorm:"type:text;comment:标签 JSON"`
	Status              enums.PostStatus        `gorm:"type:varchar(20);not null;default:'normal';index;comment:状态"`
	AcceptedCommentID   int64                   `gorm:"type:bigint;not null;default:0;index;comment:已采纳评论 ID"`
	CommentCount        int64                   `gorm:"type:bigint;not null;default:0;comment:评论次数"`
	ReactionCount       int64                   `gorm:"type:bigint;not null;default:0;comment:互动次数"`
	ViewCount           int64                   `gorm:"type:bigint;not null;default:0;comment:浏览次数"`
	LastCommentedAt     *time.Time              `gorm:"index;comment:最近评论时间"`
	LastCommentUserType enums.CommentAuthorType `gorm:"type:varchar(20);not null;default:'';comment:最近评论用户类型"`
	LastCommentUserID   int64                   `gorm:"type:bigint;not null;default:0;comment:最近评论用户 ID"`
	AuditFields
}

// Comment is a comment from either a customer or a backend user.
type Comment struct {
	ID            int64                   `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	PostID        int64                   `gorm:"type:bigint;not null;index;comment:帖子 ID"`
	ParentID      int64                   `gorm:"type:bigint;not null;default:0;index;comment:父级 ID"`
	AuthorType    enums.CommentAuthorType `gorm:"type:varchar(20);not null;default:'customer';index;comment:作者类型"`
	AuthorID      int64                   `gorm:"type:bigint;not null;index;comment:作者 ID"`
	ContentType   string                  `gorm:"type:varchar(20);not null;default:'markdown';comment:内容类型"`
	Content       string                  `gorm:"type:text;comment:内容"`
	Status        enums.CommentStatus     `gorm:"type:varchar(20);not null;default:'normal';index;comment:状态"`
	ReactionCount int64                   `gorm:"type:bigint;not null;default:0;comment:互动次数"`
	ReplyCount    int64                   `gorm:"type:bigint;not null;default:0;comment:回复次数"`
	ReportCount   int64                   `gorm:"type:bigint;not null;default:0;comment:举报次数"`
	IsAccepted    bool                    `gorm:"not null;default:false;index;comment:是否被采纳"`
	AuditFields
}

type Reaction struct {
	ID           int64                `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	TargetType   enums.ReactionTarget `gorm:"type:varchar(20);not null;index;uniqueIndex:uk_reaction;comment:目标类型"`
	TargetID     int64                `gorm:"type:bigint;not null;index;uniqueIndex:uk_reaction;comment:目标 ID"`
	UserID       int64                `gorm:"type:bigint;not null;index;uniqueIndex:uk_reaction;comment:用户 ID"`
	ReactionType enums.ReactionType   `gorm:"type:varchar(20);not null;default:'like';uniqueIndex:uk_reaction;comment:互动类型"`
	CreatedAt    time.Time            `gorm:"not null;index;comment:创建时间"`
	UpdatedAt    time.Time            `gorm:"not null;index;comment:更新时间"`
}

type CommentReport struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	CommentID int64     `gorm:"type:bigint;not null;index;uniqueIndex:uk_comment_report;comment:评论 ID"`
	UserID    int64     `gorm:"type:bigint;not null;index;uniqueIndex:uk_comment_report;comment:用户 ID"`
	Reason    string    `gorm:"type:varchar(255);not null;default:'';comment:原因"`
	CreatedAt time.Time `gorm:"type:datetime;not null;index;comment:创建时间"`
}

// KnowledgeRetrieveLog 检索日志表。
type KnowledgeRetrieveLog struct {
	ID                 int64     `gorm:"primaryKey;autoIncrement;comment:日志主键。"`                                                          // ID 为日志主键。
	KnowledgeBaseID    int64     `gorm:"type:bigint;not null;index;comment:知识库ID。"`                                                       // KnowledgeBaseID 为知识库ID。
	Channel            string    `gorm:"type:varchar(30);not null;default:'';index;comment:渠道：im会话, agent_assist坐席辅助, api开放接口, debug调试。"` // Channel 为渠道：im会话, agent_assist坐席辅助, api开放接口, debug调试。
	Scene              string    `gorm:"type:varchar(50);not null;default:'';index;comment:场景：first_response首响, assist辅助, qa问答。"`         // Scene 为场景：first_response首响, assist辅助, qa问答。
	SessionID          string    `gorm:"type:varchar(64);not null;default:'';index;comment:会话ID。"`                                        // SessionID 为会话ID。
	ConversationID     int64     `gorm:"type:bigint;not null;default:0;index;comment:会话ID。"`                                              // ConversationID 为会话ID。
	RequestID          string    `gorm:"type:varchar(64);not null;default:'';index;comment:请求ID。"`                                        // RequestID 为请求ID。
	Question           string    `gorm:"type:text;comment:原始问题。"`                                                                         // Question 为原始问题。
	RewriteQuestion    string    `gorm:"type:text;comment:改写后问题。"`                                                                        // RewriteQuestion 为改写后问题。
	Answer             string    `gorm:"type:text;comment:生成的答案。"`                                                                        // Answer 为生成的答案。
	AnswerStatus       int       `gorm:"type:int;not null;default:1;index;comment:答案状态：1正常 2无答案 3兜底 4风控拦截。"`                              // AnswerStatus 为答案状态：1正常 2无答案 3兜底 4风控拦截。
	HitCount           int       `gorm:"type:int;not null;default:0;comment:命中数量。"`                                                       // HitCount 为命中数量。
	TopScore           float64   `gorm:"type:decimal(5,4);not null;default:0;comment:最高相似度分数。"`                                           // TopScore 为最高相似度分数。
	ChunkProvider      string    `gorm:"type:varchar(30);not null;default:'';index;comment:分块 provider。"`                                 // ChunkProvider 为分块 provider。
	ChunkTargetTokens  int       `gorm:"type:int;not null;default:0;comment:目标 token 数。"`                                                 // ChunkTargetTokens 为目标 token 数。
	ChunkMaxTokens     int       `gorm:"type:int;not null;default:0;comment:最大 token 数。"`                                                 // ChunkMaxTokens 为最大 token 数。
	ChunkOverlapTokens int       `gorm:"type:int;not null;default:0;comment:重叠 token 数。"`                                                 // ChunkOverlapTokens 为重叠 token 数。
	RerankEnabled      bool      `gorm:"not null;default:false;index;comment:否启用 rerank。"`                                                // RerankEnabled 是否启用 rerank。
	RerankLimit        int       `gorm:"type:int;not null;default:0;comment:rerank 条数。"`                                                  // RerankLimit 为 rerank 条数。
	CitationCount      int       `gorm:"type:int;not null;default:0;comment:最终引用条数。"`                                                     // CitationCount 为最终引用条数。
	UsedChunkCount     int       `gorm:"type:int;not null;default:0;comment:进入上下文的 chunk 数。"`                                             // UsedChunkCount 为进入上下文的 chunk 数。
	LatencyMs          int64     `gorm:"type:bigint;not null;default:0;comment:总耗时毫秒。"`                                                   // LatencyMs 为总耗时毫秒。
	RetrieveMs         int64     `gorm:"type:bigint;not null;default:0;comment:检索耗时毫秒。"`                                                  // RetrieveMs 为检索耗时毫秒。
	GenerateMs         int64     `gorm:"type:bigint;not null;default:0;comment:生成耗时毫秒。"`                                                  // GenerateMs 为生成耗时毫秒。
	PromptTokens       int       `gorm:"type:int;not null;default:0;comment:prompt token数。"`                                              // PromptTokens 为prompt token数。
	CompletionTokens   int       `gorm:"type:int;not null;default:0;comment:completion token数。"`                                          // CompletionTokens 为completion token数。
	ModelName          string    `gorm:"type:varchar(100);not null;default:'';comment:使用的模型名称。"`                                          // ModelName 为使用的模型名称。
	TraceData          string    `gorm:"type:text;comment:链路追踪数据JSON。"`                                                                   // TraceData 为链路追踪数据JSON。
	CreatedAt          time.Time `gorm:"not null;index;comment:创建时间"`
}

// KnowledgeRetrieveHit 检索命中详情表。
type KnowledgeRetrieveHit struct {
	ID              int64     `gorm:"primaryKey;autoIncrement;comment:命中记录主键。"`                  // ID 为命中记录主键。
	RetrieveLogID   int64     `gorm:"type:bigint;not null;index;comment:检索日志ID。"`                // RetrieveLogID 为检索日志ID。
	KnowledgeBaseID int64     `gorm:"type:bigint;not null;default:0;index;comment:命中来源知识库ID。"`   // KnowledgeBaseID 为命中来源知识库ID。
	ChunkID         int64     `gorm:"type:bigint;not null;index;comment:切片ID。"`                  // ChunkID 为切片ID。
	DocumentID      int64     `gorm:"type:bigint;not null;index;comment:文档ID。"`                  // DocumentID 为文档ID。
	DocumentTitle   string    `gorm:"type:varchar(255);not null;default:'';comment:文档标题。"`       // DocumentTitle 为文档标题。
	FaqID           int64     `gorm:"type:bigint;not null;default:0;index;comment:FAQ ID。"`      // FaqID 为 FAQ ID。
	FaqQuestion     string    `gorm:"type:varchar(500);not null;default:'';comment:FAQ 问题。"`     // FaqQuestion 为 FAQ 问题。
	ChunkNo         int       `gorm:"type:int;not null;default:0;comment:切片序号。"`                 // ChunkNo 为切片序号。
	Title           string    `gorm:"type:varchar(255);not null;default:'';comment:切片标题。"`       // Title 为切片标题。
	SectionPath     string    `gorm:"type:text;comment:章节路径。"`                                   // SectionPath 为章节路径。
	ChunkType       string    `gorm:"type:varchar(30);not null;default:'';comment:切片类型。"`        // ChunkType 为切片类型。
	Provider        string    `gorm:"type:varchar(30);not null;default:'';comment:分块 provider。"` // Provider 为分块 provider。
	RankNo          int       `gorm:"type:int;not null;default:0;comment:排名。"`                   // RankNo 为排名。
	Score           float64   `gorm:"type:decimal(5,4);not null;default:0;comment:相似度分数。"`       // Score 为相似度分数。
	RerankScore     float64   `gorm:"type:decimal(5,4);not null;default:0;comment:重排分数。"`        // RerankScore 为重排分数。
	UsedInAnswer    bool      `gorm:"not null;default:false;comment:否用于生成答案。"`                   // UsedInAnswer 是否用于生成答案。
	IsCitation      bool      `gorm:"not null;default:false;comment:否作为引用返回。"`                   // IsCitation 是否作为引用返回。
	Snippet         string    `gorm:"type:text;comment:内容片段。"`                                   // Snippet 为内容片段。
	CreatedAt       time.Time `gorm:"not null;index;comment:创建时间"`
}

// KnowledgeFeedback 问答反馈表。
type KnowledgeFeedback struct {
	ID             int64     `gorm:"primaryKey;autoIncrement;comment:反馈主键。"`                                 // ID 为反馈主键。
	RetrieveLogID  int64     `gorm:"type:bigint;not null;index;comment:检索日志ID。"`                             // RetrieveLogID 为检索日志ID。
	FeedbackType   int       `gorm:"type:int;not null;default:1;index;comment:反馈类型：1点赞 2点踩 3无帮助 4引用错误 5其他。"` // FeedbackType 为反馈类型：1点赞 2点踩 3无帮助 4引用错误 5其他。
	FeedbackReason string    `gorm:"type:varchar(500);not null;default:'';comment:反馈原因。"`                    // FeedbackReason 为反馈原因。
	UserID         int64     `gorm:"type:bigint;not null;default:0;index;comment:用户ID。"`                     // UserID 为用户ID。
	AgentID        int64     `gorm:"type:bigint;not null;default:0;index;comment:坐席ID。"`                     // AgentID 为坐席ID。
	Remark         string    `gorm:"type:text;comment:备注。"`                                                  // Remark 为备注。
	CreatedAt      time.Time `gorm:"not null;index;comment:创建时间"`
}

// SkillDefinition 表示可由后台配置并参与运行时路由的 Skill 定义。
type SkillDefinition struct {
	ID            int64        `gorm:"primaryKey;autoIncrement;comment:Skill 主键。"`                                         // ID 为 Skill 主键。
	Name          string       `gorm:"type:varchar(100);not null;default:'';index;comment:Skill 的展示名称，用于后台列表、配置页和人工选择场景。"` // Name 为 Skill 的展示名称，用于后台列表、配置页和人工选择场景。
	Description   string       `gorm:"type:varchar(255);not null;default:'';comment:Skill 的简要说明，用于描述该 Skill 的适用场景和职责边界。"`  // Description 为 Skill 的简要说明，用于描述该 Skill 的适用场景和职责边界。
	Instruction   string       `gorm:"type:text;comment:Skill 的主体说明文档存储字段，使用 Markdown 编写，供 Agent 理解任务目标、步骤和工具使用要求。"`       // Instruction 为 Skill 的主体说明文档存储字段，使用 Markdown 编写，供 Agent 理解任务目标、步骤和工具使用要求。
	Examples      string       `gorm:"type:text;comment:示例问法 JSON 数组字符串。"`                                                 // Examples 为示例问法 JSON 数组字符串。
	ToolWhitelist string       `gorm:"type:text;comment:允许使用的工具编码 JSON 数组字符串。"`                                            // ToolWhitelist 为允许使用的工具编码 JSON 数组字符串。
	Status        enums.Status `gorm:"type:int;not null;default:0;index;comment:Skill 当前状态，使用全局通用状态：0启用 1禁用 2删除。"`         // Status 为 Skill 当前状态，使用全局通用状态：0启用 1禁用 2删除。
	Remark        string       `gorm:"type:text;comment:后台备注，用于记录配置说明、维护信息或内部协作信息。"`                                       // Remark 为后台备注，用于记录配置说明、维护信息或内部协作信息。
	AuditFields
}

// ConversationInterrupt 表示会话级待恢复中断记录。
type ConversationInterrupt struct {
	ID                  int64      `gorm:"primaryKey;autoIncrement;comment:主键 ID"`
	ConversationID      int64      `gorm:"type:bigint;not null;default:0;index;comment:会话 ID"`
	AIAgentID           int64      `gorm:"type:bigint;not null;default:0;index;comment:AI Agent ID"`
	AgentRunID          int64      `gorm:"type:bigint;not null;default:0;index;comment:Agent 运行 ID"`
	AgentStepID         int64      `gorm:"type:bigint;not null;default:0;index;comment:Agent 步骤 ID"`
	SourceMessageID     int64      `gorm:"type:bigint;not null;default:0;index;comment:来源消息 ID"`
	LastResumeMessageID int64      `gorm:"type:bigint;not null;default:0;index;comment:最近恢复消息 ID"`
	WorkflowRunID       int64      `gorm:"type:bigint;not null;default:0;index;comment:工作流运行 ID"`
	WorkflowNodeID      string     `gorm:"type:varchar(100);not null;default:'';index;comment:工作流节点 ID"`
	CheckPointID        string     `gorm:"type:varchar(128);not null;default:'';uniqueIndex;comment:检查点 ID"`
	InterruptID         string     `gorm:"type:varchar(255);not null;default:'';index;comment:中断 ID"`
	InterruptType       string     `gorm:"type:varchar(50);not null;default:'';index;comment:中断类型"`
	Status              string     `gorm:"type:varchar(30);not null;default:'';index;comment:状态"`
	PromptText          string     `gorm:"type:text;comment:提示词文本"`
	RequestData         string     `gorm:"type:text;comment:请求数据"`
	CheckPointData      string     `gorm:"type:text;comment:检查点数据"`
	ResumeCount         int        `gorm:"type:int;not null;default:0;comment:恢复次数"`
	ExpiresAt           *time.Time `gorm:"index;comment:过期时间"`
	CreatedAt           time.Time  `gorm:"not null;index;comment:创建时间"`
	UpdatedAt           time.Time  `gorm:"not null;index;comment:更新时间"`
}
