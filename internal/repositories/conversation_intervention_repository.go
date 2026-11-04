package repositories

import (
	"agent-desk/internal/models"
	"agent-desk/internal/pkg/enums"

	"gorm.io/gorm"
)

func (r *conversationRepository) AssignAIServing(db *gorm.DB, conversationID int64, columns map[string]any) (bool, error) {
	result := db.Model(&models.Conversation{}).
		Where("id = ? AND status = ?", conversationID, enums.IMConversationStatusAIServing).
		Updates(columns)
	return result.RowsAffected > 0, result.Error
}
