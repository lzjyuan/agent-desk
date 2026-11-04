package repositories

import (
	"time"

	"agent-desk/internal/models"

	"gorm.io/gorm"
)

func (r *agentProfileRepository) TouchOnlineAt(db *gorm.DB, userID int64, now time.Time, interval time.Duration) error {
	return db.Model(&models.AgentProfile{}).
		Where("user_id = ?", userID).
		Where("last_online_at IS NULL OR last_online_at <= ?", now.Add(-interval)).
		UpdateColumn("last_online_at", now).Error
}
