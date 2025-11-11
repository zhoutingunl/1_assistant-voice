package mysql

import (
	"time"
)

// User 用户表结构体
type User struct {
	ID           int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:用户唯一ID"`
	Username     string    `gorm:"column:username;type:varchar(50);not null;uniqueIndex:idx_username;comment:用户名（唯一）"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null;comment:密码哈希"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;default:current_timestamp;comment:创建时间"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;default:current_timestamp on update current_timestamp;comment:更新时间"`
	IsActive     int8      `gorm:"column:is_active;type:tinyint;default:1;comment:账号状态（1=正常，0=禁用）"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// ChatHistory 历史对话表结构体（对应 chat_histories 表）
type ChatHistory struct {
	ID        int64     `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:对话组唯一ID"`
	UserID    int64     `gorm:"column:user_id;type:bigint;not null;index:idx_user_id;comment:关联用户ID"` // 关联用户表的外键
	UserName  string    `gorm:"column:username;type:varchar(50);default:null;comment:关联用户名称"`         // 允许为 NULL，用指针类型
	Filepath  string    `gorm:"column:filepath;type:varchar(100);default:null;comment:文件地址"`          // 允许为 NULL，用指针类型
	StartTime time.Time `gorm:"column:start_time;type:datetime;default:current_timestamp;comment:对话开始时间"`
	EndTime   time.Time `gorm:"column:end_time;type:datetime;default:null;comment:对话结束时间"` // 允许为 NULL，用指针类型
	IsDeleted int8      `gorm:"column:is_deleted;type:tinyint;default:0;comment:软删除标识（1=已删除，0=正常）"`
}

// TableName 指定表名（与数据库表名一致）
func (ChatHistory) TableName() string {
	return "chat_histories"
}

// Application 应用信息表结构体
type Application struct {
	ID          int64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement;comment:应用唯一ID"`
	Name        string `gorm:"column:name;type:varchar(50);not null;uniqueIndex:idx_name;comment:应用名称"`
	AppKey      string `gorm:"column:app_key;type:varchar(100);not null;uniqueIndex:idx_app_key;comment:应用地址（代码调用用）"`
	Description string `gorm:"column:description;type:text;default:null;comment:应用描述"` // 允许为 NULL
	IsEnabled   int8   `gorm:"column:is_enabled;type:tinyint;default:1;comment:启用状态（1=可调用，0=禁用）"`
}

// TableName 指定表名
func (Application) TableName() string {
	return "applications"
}
