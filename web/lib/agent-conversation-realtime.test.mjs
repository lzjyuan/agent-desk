import assert from "node:assert/strict"
import test from "node:test"

import {
  shouldNotifyAssignedCustomerMessage,
  shouldReloadConversationListForRealtimePatch,
} from "./agent-conversation-realtime.ts"

test("reloads conversation list when realtime patch changes list membership fields", () => {
  assert.equal(
    shouldReloadConversationListForRealtimePatch({
      conversationId: 1,
      status: 3,
      currentAssigneeId: 101,
    }),
    true
  )
})

test("keeps local patching for message summary and unread-only changes", () => {
  assert.equal(
    shouldReloadConversationListForRealtimePatch({
      conversationId: 1,
      lastMessageId: 9,
      lastMessageSummary: "hello",
      agentUnreadCount: 1,
    }),
    false
  )
})

test("notifies the assigned agent outside the visible conversation", () => {
  const assignedMessage = {
    senderType: "customer",
    isActive: true,
    currentAssigneeId: 101,
    currentUserId: 101,
    isConversationVisible: false,
  }

  assert.equal(shouldNotifyAssignedCustomerMessage(assignedMessage), true)
  assert.equal(shouldNotifyAssignedCustomerMessage({ ...assignedMessage, isConversationVisible: true }), false)
  assert.equal(shouldNotifyAssignedCustomerMessage({ ...assignedMessage, isActive: false }), false)
  assert.equal(shouldNotifyAssignedCustomerMessage({ ...assignedMessage, currentAssigneeId: 0 }), false)
  assert.equal(shouldNotifyAssignedCustomerMessage({ ...assignedMessage, currentUserId: 102 }), false)
  assert.equal(shouldNotifyAssignedCustomerMessage({ ...assignedMessage, senderType: "agent" }), false)
})
