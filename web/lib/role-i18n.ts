const SEEDED_ROLE_LABELS: Record<string, string> = {
  super_admin: "Super admin",
  admin: "Admin",
  cs_team_leader: "Support team lead",
  cs_user: "Support agent",
}

const SEEDED_ROLE_LABELS_ZH: Record<string, string> = {
  super_admin: "超级管理员",
  admin: "管理员",
  cs_team_leader: "客服组长",
  cs_user: "客服",
}

export function getRoleDisplayName(
  code: string | undefined,
  fallbackName: string,
  locale: string
) {
  const labels = locale === "zh-CN" ? SEEDED_ROLE_LABELS_ZH : locale === "en-US" ? SEEDED_ROLE_LABELS : null
  const label = labels?.[code?.trim() ?? ""]
  return label || fallbackName
}
