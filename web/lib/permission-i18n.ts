const PERMISSION_ACTION_LABELS: Record<string, string> = {
  view: "View",
  create: "Create",
  update: "Update",
  delete: "Delete",
  assignRole: "Assign roles to",
  assignPermission: "Assign permissions to",
  sync: "Sync",
  revoke: "Revoke",
  assign: "Assign",
  transfer: "Transfer",
  close: "Close",
  send: "Send",
  tag: "Manage tags for",
  handover: "Handle handoffs for",
  recycle: "Recycle",
  linkCustomer: "Link customers to",
  changeStatus: "Change status for",
  progress: "Update progress for",
  resetUserTokenSecret: "Reset user token secret for",
  updateStatus: "Update status for",
  config: "Configure service rules for",
  batchGenerate: "Batch generate",
  call: "Call",
}

const PERMISSION_ACTION_LABELS_ZH: Record<string, string> = {
  view: "查看",
  create: "创建",
  update: "更新",
  delete: "删除",
  assignRole: "分配",
  assignPermission: "分配",
  sync: "同步",
  revoke: "踢除",
  assign: "分配",
  transfer: "转接",
  close: "关闭",
  send: "发送",
  tag: "管理",
  handover: "处理",
  recycle: "回收",
  linkCustomer: "关联",
  changeStatus: "变更",
  progress: "更新",
  resetUserTokenSecret: "重置",
  updateStatus: "更新",
  config: "配置",
  batchGenerate: "批量生成",
  call: "调用",
}

const PERMISSION_RESOURCE_LABELS: Record<string, { singular: string; plural: string }> = {
  user: { singular: "user", plural: "users" },
  role: { singular: "role", plural: "roles" },
  permission: { singular: "permission", plural: "permissions" },
  session: { singular: "session", plural: "sessions" },
  conversation: { singular: "conversation", plural: "conversations" },
  ticket: { singular: "ticket", plural: "tickets" },
  notification: { singular: "notification", plural: "notifications" },
  quickReply: { singular: "quick reply", plural: "quick replies" },
  tag: { singular: "tag", plural: "tags" },
  company: { singular: "company", plural: "companies" },
  channel: { singular: "channel", plural: "channels" },
  wxworkOutbox: { singular: "WeCom outbox record", plural: "WeCom outbox records" },
  customer: { singular: "customer", plural: "customers" },
  agent: { singular: "agent", plural: "agents" },
  agentTeam: { singular: "agent team", plural: "agent teams" },
  agentTeamSchedule: { singular: "agent team schedule", plural: "agent team schedules" },
  asset: { singular: "file asset", plural: "file assets" },
  aiAgent: { singular: "AI Agent", plural: "AI Agents" },
  aiConfig: { singular: "AI configuration", plural: "AI configurations" },
  knowledgeBase: { singular: "knowledge base", plural: "knowledge bases" },
  knowledgeDocument: { singular: "knowledge document", plural: "knowledge documents" },
  knowledgeFAQ: { singular: "knowledge FAQ", plural: "knowledge FAQs" },
  docPage: { singular: "documentation page", plural: "documentation pages" },
  supportCommunity: { singular: "support community", plural: "support communities" },
  supportConfig: { singular: "support configuration", plural: "support configurations" },
  skillDefinition: { singular: "Skill definition", plural: "Skill definitions" },
  mcp: { singular: "MCP tool", plural: "MCP tools" },
}

const PERMISSION_RESOURCE_LABELS_ZH: Record<string, string> = {
  user: "用户",
  role: "角色",
  permission: "权限",
  session: "会话",
  conversation: "会话",
  ticket: "工单",
  notification: "通知",
  quickReply: "快捷回复",
  tag: "标签",
  company: "公司",
  channel: "接入渠道",
  wxworkOutbox: "企微 Outbox",
  customer: "客户",
  agent: "客服",
  agentTeam: "客服组",
  agentTeamSchedule: "客服组排班",
  asset: "文件资源",
  aiAgent: "AI Agent",
  aiConfig: "AI 配置",
  knowledgeBase: "知识库",
  knowledgeDocument: "知识文档",
  knowledgeFAQ: "知识FAQ",
  docPage: "文档中心页面",
  supportCommunity: "支持中心社区",
  supportConfig: "支持中心配置",
  skillDefinition: "技能定义",
  mcp: "MCP工具",
}

const PERMISSION_NAME_OVERRIDES: Record<string, string> = {
  "user.assignRole": "Assign user roles",
  "role.assignPermission": "Assign role permissions",
  "session.revoke": "Revoke sessions",
  "conversation.linkCustomer": "Link conversation customer",
  "ticket.changeStatus": "Change ticket status",
  "ticket.progress": "Update ticket progress",
  "channel.resetUserTokenSecret": "Reset channel user token secret",
  "agent.config": "Configure agent service rules",
  "agentTeamSchedule.batchGenerate": "Batch generate agent team schedules",
  "wxworkOutbox.update": "Handle WeCom outbox records",
  "asset.create": "Upload file assets",
  "supportCommunity.update": "Manage support communities",
  "mcp.view": "View MCP debug information",
  "mcp.call": "Call MCP tools",
}

const PERMISSION_NAME_OVERRIDES_ZH: Record<string, string> = {
  "user.assignRole": "分配用户角色",
  "role.assignPermission": "分配角色权限",
  "conversation.send": "发送会话消息",
  "conversation.tag": "管理会话标签",
  "conversation.handover": "处理会话交接",
  "conversation.linkCustomer": "关联会话客户",
  "ticket.assign": "指派工单",
  "ticket.changeStatus": "变更工单状态",
  "ticket.progress": "更新工单进展",
  "agent.updateStatus": "更新客服状态",
  "agent.config": "配置客服服务规则",
  "channel.resetUserTokenSecret": "重置接入渠道用户Token Secret",
  "asset.create": "上传文件资源",
  "wxworkOutbox.update": "处置企微 Outbox",
  "supportCommunity.update": "管理支持中心社区",
  "mcp.view": "查看MCP调试信息",
  "mcp.call": "调用MCP工具",
}

export function getPermissionDisplayName(
  code: string | undefined,
  fallbackName: string,
  locale: string
) {
  const normalizedCode = code?.trim() ?? ""
  if (!normalizedCode) {
    return fallbackName
  }
  const isChinese = locale === "zh-CN"
  if (!isChinese && locale !== "en-US") {
    return fallbackName
  }
  const override = (isChinese ? PERMISSION_NAME_OVERRIDES_ZH : PERMISSION_NAME_OVERRIDES)[normalizedCode]
  if (override) {
    return override
  }
  const [resourceKey, actionKey] = normalizedCode.split(".")
  if (isChinese) {
    const resource = PERMISSION_RESOURCE_LABELS_ZH[resourceKey]
    const action = PERMISSION_ACTION_LABELS_ZH[actionKey]
    const separator = resourceKey === "aiAgent" || resourceKey === "aiConfig" ? " " : ""
    return resource && action ? `${action}${separator}${resource}` : fallbackName
  }
  const resource = PERMISSION_RESOURCE_LABELS[resourceKey]
  const action = PERMISSION_ACTION_LABELS[actionKey]
  if (!resource || !action) {
    return fallbackName
  }
  return `${action} ${resource.plural}`
}

export function getPermissionGroupName(groupName: string | undefined, locale: string) {
  const normalizedGroupName = groupName?.trim() ?? ""
  if (locale === "zh-CN") {
    return PERMISSION_RESOURCE_LABELS_ZH[normalizedGroupName] || normalizedGroupName
  }
  if (locale !== "en-US") {
    return normalizedGroupName
  }
  const resource = PERMISSION_RESOURCE_LABELS[normalizedGroupName]
  if (!resource) {
    return normalizedGroupName
  }
  return sentenceCase(resource.plural)
}

function sentenceCase(value: string) {
  if (value === "AI Agents" || value === "MCP tools") {
    return value
  }
  return value.charAt(0).toUpperCase() + value.slice(1)
}
