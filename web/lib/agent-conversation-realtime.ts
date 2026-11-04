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
  isConversationVisible,
}: {
  senderType: string
  isActive: boolean
  currentAssigneeId: number
  currentUserId: number
  isConversationVisible: boolean
}) {
  return senderType === "customer" &&
    isActive &&
    currentAssigneeId > 0 &&
    currentAssigneeId === currentUserId &&
    !isConversationVisible
}
