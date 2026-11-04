import assert from "node:assert/strict"
import test from "node:test"

import { applyRecalledMessage, preserveRecalledMessage } from "./im-message-recall.ts"

test("recalled websocket event hides only its message content immediately", () => {
  const messages = [
    { id: 1, content: "private text", payload: "image-url", sendStatus: 1 },
    { id: 2, content: "keep", payload: "", sendStatus: 1 },
  ]
  const updated = applyRecalledMessage(messages, 1, "2026-11-04 15:22:53")

  assert.deepEqual(updated[0], {
    id: 1, content: "", payload: "", sendStatus: 6, recalledAt: "2026-11-04 15:22:53",
  })
  assert.equal(updated[1], messages[1])
  assert.equal(messages[0].content, "private text")
  assert.equal(applyRecalledMessage(messages, 99), messages)
  assert.deepEqual(
    preserveRecalledMessage(updated[0], messages[0]),
    updated[0],
    "a stale message fetch must not restore withdrawn text"
  )
})
