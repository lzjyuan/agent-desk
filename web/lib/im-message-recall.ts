export function applyRecalledMessage<T extends {
  id: number
  content: string
  payload?: string
  sendStatus: number
  recalledAt?: string
}>(messages: T[], messageId: number, recalledAt?: string): T[] {
  if (messageId <= 0 || !messages.some((message) => message.id === messageId)) {
    return messages
  }
  return messages.map((message) =>
    message.id === messageId
      ? { ...message, content: "", payload: "", sendStatus: 6, recalledAt: recalledAt ?? message.recalledAt }
      : message
  )
}

export function preserveRecalledMessage<T extends {
  content: string
  payload?: string
  sendStatus: number
  recalledAt?: string
}>(current: T | undefined, incoming: T): T {
  if (!current || (current.sendStatus !== 6 && !current.recalledAt) || incoming.sendStatus === 6 || incoming.recalledAt) {
    return incoming
  }
  return { ...incoming, content: "", payload: "", sendStatus: 6, recalledAt: current.recalledAt }
}
