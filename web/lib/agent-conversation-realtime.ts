import type { RealtimeConversationPatch } from "@/lib/im-realtime-state"

const listMembershipFields = new Set<keyof RealtimeConversationPatch>([
  "status",
  "currentAssigneeId",
  "currentTeamId",
])

export function shouldReloadConversationListForRealtimePatch(
  patch: RealtimeConversationPatch | null | undefined
) {
  if (!patch) {
    return false
  }
  return Object.keys(patch).some((key) =>
    listMembershipFields.has(key as keyof RealtimeConversationPatch)
  )
}

export function shouldNotifyAssignedCustomerMessage({
  senderType,
  isActive,
  currentAssigneeId,
  currentUserId,
}: {
  senderType: string
  isActive: boolean
  currentAssigneeId: number
  currentUserId: number
}) {
  return senderType === "customer" &&
    isActive &&
    currentAssigneeId > 0 &&
    currentAssigneeId === currentUserId
}

export function shouldNotifyConversationArrival({
  eventType,
  isPending,
  isActive,
  currentAssigneeId,
  currentUserId,
  canViewConversations,
}: {
  eventType: string
  isPending: boolean
  isActive: boolean
  currentAssigneeId: number
  currentUserId: number
  canViewConversations: boolean
}) {
  if (currentUserId <= 0 || !canViewConversations) return false
  if (eventType === "conversation.pending") return isPending && currentAssigneeId === 0
  return (eventType === "conversation.assigned" || eventType === "conversation.transferred") &&
    isActive && currentAssigneeId === currentUserId
}
