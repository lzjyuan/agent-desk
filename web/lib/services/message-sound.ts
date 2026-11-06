let audio: HTMLAudioElement | undefined
type MessageSound = "message" | "notice" | "customer"
let playingSound: MessageSound | undefined
const playedMessages = new Set<string>()

export async function playMessageSound(messageKey: string, sound: MessageSound = "message") {
  messageKey = `${sound}:${messageKey}`
  if (typeof Audio === "undefined" || playedMessages.has(messageKey)) return

  audio ??= new Audio("/sounds/customer-message.mp3")
  playedMessages.add(messageKey)
  if (playedMessages.size > 500) {
    playedMessages.delete(playedMessages.values().next().value!)
  }
  if (!audio.paused) {
    if (sound !== "notice" || playingSound === "notice") return
    audio.pause()
  }

  try {
    playingSound = sound
    audio.currentTime = 0
    await audio.play()
  } catch (error) {
    playedMessages.delete(messageKey)
    if (error instanceof DOMException && error.name === "AbortError") return
    throw error
  }
}
