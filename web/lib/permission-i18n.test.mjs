import assert from "node:assert/strict"
import test from "node:test"
import ts from "typescript"
import { readFile } from "node:fs/promises"
import vm from "node:vm"

async function loadModule() {
  const source = await readFile(new URL("./permission-i18n.ts", import.meta.url), "utf8")
  const compiled = ts.transpileModule(source, {
    compilerOptions: {
      target: ts.ScriptTarget.ES2017,
      module: ts.ModuleKind.CommonJS,
    },
    fileName: "permission-i18n.ts",
  })
  const sandbox = {
    exports: {},
    module: { exports: {} },
  }
  sandbox.exports = sandbox.module.exports
  vm.runInNewContext(compiled.outputText, sandbox)
  return sandbox.module.exports
}

test("localizes seeded permission display names to English", async () => {
  const { getPermissionDisplayName, getPermissionGroupName } = await loadModule()

  assert.equal(getPermissionDisplayName("user.assignRole", "\u5206\u914d\u7528\u6237\u89d2\u8272", "en-US"), "Assign user roles")
  assert.equal(getPermissionDisplayName("agentTeamSchedule.batchGenerate", "\u6279\u91cf\u751f\u6210\u5ba2\u670d\u7ec4\u6392\u73ed", "en-US"), "Batch generate agent team schedules")
  assert.equal(getPermissionGroupName("agentTeamSchedule", "en-US"), "Agent team schedules")
})

test("localizes every built-in permission without relying on database names", async () => {
  const { getPermissionDisplayName, getPermissionGroupName } = await loadModule()
  const source = await readFile(new URL("../../internal/pkg/constants/auth.go", import.meta.url), "utf8")
  const permissions = [...source.matchAll(/Permission\{Name: "([^"]+)", Code: "([^"]+)", Type: "api", GroupName: "([^"]+)"/g)]

  assert.ok(permissions.length > 100)
  for (const [, name, code, group] of permissions) {
    assert.equal(getPermissionDisplayName(code, "stale database name", "zh-CN"), name, code)
    assert.notEqual(getPermissionDisplayName(code, name, "en-US"), name, code)
    assert.notEqual(getPermissionGroupName(group, "zh-CN"), group, group)
    assert.notEqual(getPermissionGroupName(group, "en-US"), group, group)
  }

  assert.equal(getPermissionDisplayName("custom.export", "Custom export", "zh-CN"), "Custom export")
  assert.equal(getPermissionGroupName("custom", "zh-CN"), "custom")
})
