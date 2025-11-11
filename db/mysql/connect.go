package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	return db
}

func init() {
	// 1. 配置 MySQL 连接参数（DSN）
	// 格式：user:password@tcp(ip:port)/dbname?param=value
	dsn := fmt.Sprintf("root:123456@tcp(192.168.171.129:13306)/chatserver?charset=utf8&parseTime=True&loc=Local")
	//dsn := fmt.Sprintf("root:123456@tcp(192.168.190.128:13306)/chatserver?charset=utf8&parseTime=True&loc=Local")
	// 参数说明：
	// - root: 数据库用户名
	// - edge#Admin: 数据库密码
	// - 10.40.10.64:30306: 数据库地址和端口
	// - sqlName: 要连接的数据库名
	// - charset=utf8mb4: 支持 emoji 表情
	// - parseTime=true: 允许 GORM 将 datetime 解析为 time.Time
	// - loc=Local: 时区设置为本地时区（避免时间偏移）

	// 2. 连接数据库并创建 GORM 实例
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 可选配置：开启日志（开发环境建议开启，生产环境关闭）
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Printf("连接数据库失败: %v\n", err)
		return
	}

	// 3. 测试连接（可选，验证是否能正常访问数据库）
	sqlDB, err := db.DB() // 获取底层 sql.DB 对象
	if err != nil {
		fmt.Printf("获取数据库连接失败: %v\n", err)
		return
	}
	// 尝试 ping 数据库
	if err = sqlDB.Ping(); err != nil {
		fmt.Printf("ping 数据库失败: %v\n", err)
		return
	}
}
