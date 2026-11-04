package dashboard

import (
	"testing"

	"agent-desk/internal/models"

	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
)

func TestBuildConversationResponsesIncludesChannelName(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:conversation_channel_names?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get db: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	if err := db.AutoMigrate(&models.Channel{}, &models.ConversationReadState{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	sqls.SetDB(db)
	if err := db.Create(&models.Channel{ID: 2, Name: "slotSweet3"}).Error; err != nil {
		t.Fatalf("create channel: %v", err)
	}

	results := buildConversationResponses([]models.Conversation{
		{ID: 1, ChannelID: 2},
		{ID: 2, ChannelID: 99},
	}, "zh-CN")
	if len(results) != 2 || results[0].ChannelName != "slotSweet3" || results[1].ChannelName != "" {
		t.Fatalf("unexpected channel names: %+v", results)
	}
}
