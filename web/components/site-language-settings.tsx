"use client"

import { LanguagesIcon } from "lucide-react"
import { useRef, useState, type FormEvent } from "react"
import { toast } from "sonner"

import { useAuth } from "@/components/auth-provider"
import { useConfirm } from "@/components/confirm-provider"
import { OptionCombobox } from "@/components/option-combobox"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { DEFAULT_CLIENT_LOCALE, isSupportedLocale } from "@/i18n/config"
import { translateMessage } from "@/i18n/messages"
import { useAppLocale, useI18n } from "@/i18n/provider"
import { fetchPublicConfig, saveSiteLanguages } from "@/lib/api/config"

export function SiteLanguageSettings() {
  const t = useI18n()
  const { session } = useAuth()
  const confirm = useConfirm()
  const { locale, setLocale } = useAppLocale()
  const [open, setOpen] = useState(false)
  const [language, setLanguage] = useState(locale)
  const [currentLanguage, setCurrentLanguage] = useState(locale)
  const [clientLanguage, setClientLanguage] = useState(DEFAULT_CLIENT_LOCALE)
  const [currentClientLanguage, setCurrentClientLanguage] = useState(DEFAULT_CLIENT_LOCALE)
  const [loading, setLoading] = useState(false)
  const [ready, setReady] = useState(false)
  const loadId = useRef(0)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState("")

  if (!session?.permissions.includes("systemConfig.update")) return null

  async function close() {
    if (saving) return
    if (ready && (language !== currentLanguage || clientLanguage !== currentClientLanguage)) {
      const discard = await confirm({
        title: t("siteLanguage.unsavedTitle"),
        description: t("siteLanguage.unsavedDescription"),
        confirmText: t("siteLanguage.discardChanges"),
        cancelText: t("siteLanguage.continueEditing"),
        variant: "destructive",
      })
      if (!discard) return
    }
    loadId.current += 1
    setOpen(false)
  }

  async function save(event: FormEvent) {
    event.preventDefault()
    if (saving || !ready || (language === currentLanguage && clientLanguage === currentClientLanguage)) return
    setSaving(true)
    setError("")
    try {
      await saveSiteLanguages(language, clientLanguage)
      setCurrentLanguage(language)
      setCurrentClientLanguage(clientLanguage)
      setLocale(language)
      setOpen(false)
      toast.success(translateMessage(language, "siteLanguage.saved"))
    } catch (error) {
      setError(error instanceof Error ? error.message : t("siteLanguage.failed"))
    } finally {
      setSaving(false)
    }
  }

  return (
    <Dialog
      open={open}
      onOpenChange={(nextOpen) => {
        if (saving) return
        if (!nextOpen) {
          void close()
          return
        }
        if (nextOpen) {
          const requestId = ++loadId.current
          setLanguage(locale)
          setError("")
          setLoading(true)
          setReady(false)
          void fetchPublicConfig(true)
            .then((config) => {
              if (requestId !== loadId.current) return
              const nextLanguage = isSupportedLocale(config.language) ? config.language : locale
              setLanguage(nextLanguage)
              setCurrentLanguage(nextLanguage)
              const nextClientLanguage = isSupportedLocale(config.clientLanguage) ? config.clientLanguage : DEFAULT_CLIENT_LOCALE
              setClientLanguage(nextClientLanguage)
              setCurrentClientLanguage(nextClientLanguage)
              setReady(true)
            })
            .catch((error) => {
              if (requestId === loadId.current) setError(error instanceof Error ? error.message : t("siteLanguage.failed"))
            })
            .finally(() => {
              if (requestId === loadId.current) setLoading(false)
            })
        }
        setOpen(nextOpen)
      }}
    >
      <DialogTrigger
        render={<Button variant="outline" size="sm" className="rounded-md text-muted-foreground" />}
        aria-label={t("siteLanguage.title")}
        title={t("siteLanguage.title")}
      >
        <LanguagesIcon />
      </DialogTrigger>
      <DialogContent className="rounded-md sm:max-w-md" showCloseButton={false}>
        <form onSubmit={save} className="space-y-5">
          <DialogHeader>
            <DialogTitle>{t("siteLanguage.title")}</DialogTitle>
            <DialogDescription>{t("siteLanguage.description")}</DialogDescription>
          </DialogHeader>
          {loading ? <p role="status" className="text-sm text-muted-foreground">{t("siteLanguage.loading")}</p> : null}
          <fieldset className="space-y-2" disabled={saving || loading}>
            <legend className="mb-2 font-medium">{t("siteLanguage.label")}</legend>
            <OptionCombobox
              value={language}
              onChange={(value) => {
                if (isSupportedLocale(value)) setLanguage(value)
              }}
              options={[
                { value: "zh-CN", label: t("siteLanguage.chinese") },
                { value: "en-US", label: t("siteLanguage.english") },
              ]}
              placeholder={t("siteLanguage.label")}
              disabled={saving || loading}
            />
          </fieldset>
          <fieldset className="space-y-2" disabled={saving || loading}>
            <legend className="mb-2 font-medium">{t("siteLanguage.clientLabel")}</legend>
            <OptionCombobox
              value={clientLanguage}
              onChange={(value) => {
                if (isSupportedLocale(value)) setClientLanguage(value)
              }}
              options={[
                { value: "zh-CN", label: t("siteLanguage.chinese") },
                { value: "en-US", label: t("siteLanguage.english") },
                { value: "id-ID", label: t("siteLanguage.indonesian") },
              ]}
              placeholder={t("siteLanguage.clientLabel")}
              disabled={saving || loading}
            />
          </fieldset>
          {error ? <p role="alert" className="text-sm text-destructive">{error}</p> : null}
          <DialogFooter className="rounded-b-md">
            <Button type="button" variant="outline" disabled={saving} onClick={() => void close()}>
              {t("common.cancel")}
            </Button>
            <Button type="submit" disabled={saving || !ready || (language === currentLanguage && clientLanguage === currentClientLanguage)}>
              {t(saving ? "siteLanguage.saving" : "siteLanguage.save")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
