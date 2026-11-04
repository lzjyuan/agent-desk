let audio: HTMLAudioElement | undefined
type MessageSound = "message" | "notice" | "customer"
let playingSound: MessageSound | undefined
const playedMessages = new Set<string>()

export async function playMessageSound(messageKey: string, sound: MessageSound = "message") {
  messageKey = `${sound}:${messageKey}`
  if (typeof Audio === "undefined" || playedMessages.has(messageKey)) return

  const source = sound === "notice"
    ? "/sounds/notice.mp3"
    : sound === "customer" ? "/sounds/customer-message.mp3" : "/sounds/news.mp3"
  audio ??= new Audio(source)
  playedMessages.add(messageKey)
  if (playedMessages.size > 500) {
    playedMessages.delete(playedMessages.values().next().value!)
  }
  if (!audio.paused) {
    if (sound !== "notice" || playingSound === "notice") return
    audio.pause()
  }

  try {
    if (playingSound !== sound) audio.src = source
    playingSound = sound
    audio.currentTime = 0
    await audio.play()
  } catch (error) {
    playedMessages.delete(messageKey)
    if (error instanceof DOMException && error.name === "AbortError") return
    throw error
  }
}
