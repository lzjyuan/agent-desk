import assert from "node:assert/strict"
import test from "node:test"
import { playMessageSound } from "./message-sound.ts"

test("message sound reuses the MP3, deduplicates messages, avoids overlap and allows retry", async () => {
  const original = globalThis.Audio
  const instances = []
  let playCount = 0
  let pauseCount = 0
  let failure
  globalThis.Audio = class {
    paused = true
    currentTime = 7
    constructor(src) {
      this.src = src
      instances.push(this)
    }
    async play() {
      playCount++
      if (failure) throw failure
      this.paused = false
    }
    pause() {
      pauseCount++
      this.paused = true
    }
  }

  try {
    await playMessageSound("101:1:1")
    const [instance] = instances
    assert.equal(instance.src, "/sounds/customer-message.mp3")
    assert.equal(instance.currentTime, 0)
    assert.equal(playCount, 1)

    await playMessageSound("101:1:2")
    assert.equal(playCount, 1)
    instance.paused = true
    await playMessageSound("101:1:1")
    await playMessageSound("101:1:2")
    assert.equal(playCount, 1)

    failure = new DOMException("User interaction required", "NotAllowedError")
    await assert.rejects(playMessageSound("101:1:3"), { name: "NotAllowedError" })
    failure = undefined
    await playMessageSound("101:1:3")
    assert.equal(playCount, 3)

    instance.paused = true
    await playMessageSound("102:1:1")
    assert.equal(playCount, 4)

    await playMessageSound("102:1:5", "notice")
    assert.equal(instance.src, "/sounds/customer-message.mp3")
    assert.equal(playCount, 5)
    assert.equal(pauseCount, 1)
    await playMessageSound("102:1:5", "notice")
    await playMessageSound("102:1:6")
    assert.equal(playCount, 5)
    assert.equal(instance.src, "/sounds/customer-message.mp3")

    instance.paused = true
    await playMessageSound("102:1:7")
    assert.equal(instance.src, "/sounds/customer-message.mp3")
    assert.equal(playCount, 6)

    instance.paused = true
    await playMessageSound("customer:1:8", "customer")
    assert.equal(instance.src, "/sounds/customer-message.mp3")
    assert.equal(playCount, 7)
    await playMessageSound("customer:1:8", "customer")
    assert.equal(playCount, 7)
  } finally {
    if (original === undefined) delete globalThis.Audio
    else globalThis.Audio = original
  }
})
