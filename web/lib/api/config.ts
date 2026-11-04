import { request } from "@/lib/api/client"

export type PublicConfig = {
  language: string
  clientLanguage: string
  companyName?: string
  companyLogoUrl?: string
  passwordLoginEnabled?: boolean
  wxworkEnabled: boolean
  oidcEnabled: boolean
}

let publicConfigPromise: Promise<PublicConfig> | null = null

export function fetchPublicConfig(fresh = false) {
  if (fresh) publicConfigPromise = null
  publicConfigPromise ??= request<PublicConfig>("/api/config", {
    skipAuth: true,
  }).catch((error) => {
    publicConfigPromise = null
    throw error
  })
  return publicConfigPromise
}

export async function saveSiteLanguages(language: string, clientLanguage: string) {
  await request<void>("/api/dashboard/system-config/language", {
    method: "POST",
    body: JSON.stringify({ language, clientLanguage }),
  })
  publicConfigPromise = null
}
