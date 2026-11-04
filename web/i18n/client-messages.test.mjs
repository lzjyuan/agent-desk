import assert from "node:assert/strict"
import test from "node:test"
import { readFile } from "node:fs/promises"

const directory = new URL("../messages/", import.meta.url)
const english = JSON.parse(await readFile(new URL("en-US.json", directory), "utf8"))
const indonesian = JSON.parse(await readFile(new URL("id-ID.json", directory), "utf8"))

test("Indonesian covers customer chat and support center strings", () => {
  function check(source, translated, path) {
    for (const [key, value] of Object.entries(source)) {
      const currentPath = `${path}.${key}`
      if (typeof value === "string") {
        assert.equal(typeof translated?.[key], "string", currentPath)
        assert.deepEqual(
          [...value.matchAll(/\{\w+\}/g)].map(([placeholder]) => placeholder).sort(),
          [...translated[key].matchAll(/\{\w+\}/g)].map(([placeholder]) => placeholder).sort(),
          `${currentPath} placeholders`
        )
      } else if (value && typeof value === "object") {
        check(value, translated?.[key], currentPath)
      }
    }
  }

  for (const namespace of ["supportChat", "supportPublic", "legal", "editor", "upload", "lightbox", "common", "api"]) {
    check(english[namespace], indonesian[namespace], namespace)
  }
})
