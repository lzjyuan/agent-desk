package bootstrap

import (
	"fmt"

	"agent-desk/internal/migration"
	"agent-desk/internal/models"

	"github.com/mlogclub/simple/sqls"
)

func InitMigrations() error {
	db := sqls.DB()
	if db.Dialector.Name() == "mysql" {
		// MySQL 建表时写入表备注，并由 AutoMigrate 同步字段备注。
		for _, model := range models.Models {
			comment := models.TableComment(model)
			quotedComment := db.Dialector.Explain("?", comment)
			if err := db.Set("gorm:table_options", " COMMENT="+quotedComment).AutoMigrate(model); err != nil {
				return err
			}
			tableType, err := db.Migrator().TableType(model)
			if err != nil {
				return err
			}
			// 已存在的表不会重新建表，备注不同时单独补齐。
			if existing, _ := tableType.Comment(); existing != comment {
				name := tableType.Name()
				if err := db.Exec(fmt.Sprintf("ALTER TABLE `%s` COMMENT=%s", name, quotedComment)).Error; err != nil {
					return err
				}
			}
		}
	} else if err := db.AutoMigrate(models.Models...); err != nil {
		// SQLite 和 PostgreSQL 的默认迁移失败时，直接返回错误。
		return err
	}
	return migration.Migrate()
}
