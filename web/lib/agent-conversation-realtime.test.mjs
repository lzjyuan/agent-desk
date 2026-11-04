import assert from "node:assert/strict"
import test from "node:test"

import {
  shouldNotifyAssignedCustomerMessage,
  shouldNotifyConversationArrival,
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

test("arrival sound is limited to pending pool entries and assignments to this agent", () => {
  const pending = { eventType: "conversation.pending", isPending: true, isActive: false, currentAssigneeId: 0, currentUserId: 101, canViewConversations: true }
  assert.equal(shouldNotifyConversationArrival(pending), true)
  assert.equal(shouldNotifyConversationArrival({ ...pending, eventType: "conversation.updated" }), false)
  assert.equal(shouldNotifyConversationArrival({ ...pending, eventType: "conversation.created" }), false)
  assert.equal(shouldNotifyConversationArrival({ ...pending, eventType: "conversation.read" }), false)
  assert.equal(shouldNotifyConversationArrival({ ...pending, isPending: false }), false)
  assert.equal(shouldNotifyConversationArrival({ ...pending, currentAssigneeId: 102 }), false)
  assert.equal(shouldNotifyConversationArrival({ ...pending, canViewConversations: false }), false)
  assert.equal(shouldNotifyConversationArrival({ ...pending, currentUserId: 0 }), false)

  const assigned = { ...pending, eventType: "conversation.assigned", isPending: false, isActive: true, currentAssigneeId: 101 }
  assert.equal(shouldNotifyConversationArrival(assigned), true)
  assert.equal(shouldNotifyConversationArrival({ ...assigned, eventType: "conversation.transferred" }), true)
  assert.equal(shouldNotifyConversationArrival({ ...assigned, currentAssigneeId: 102 }), false)
  assert.equal(shouldNotifyConversationArrival({ ...assigned, isActive: false }), false)
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

test("notifies the assigned agent for active customer messages", () => {
  const assignedMessage = {
    senderType: "customer",
    isActive: true,
    currentAssigneeId: 101,
    currentUserId: 101,
  }

  assert.equal(shouldNotifyAssignedCustomerMessage(assignedMessage), true)
  assert.equal(shouldNotifyAssignedCustomerMessage({ ...assignedMessage, isActive: false }), false)
  assert.equal(shouldNotifyAssignedCustomerMessage({ ...assignedMessage, currentAssigneeId: 0 }), false)
  assert.equal(shouldNotifyAssignedCustomerMessage({ ...assignedMessage, currentUserId: 102 }), false)
  assert.equal(shouldNotifyAssignedCustomerMessage({ ...assignedMessage, senderType: "agent" }), false)
})
