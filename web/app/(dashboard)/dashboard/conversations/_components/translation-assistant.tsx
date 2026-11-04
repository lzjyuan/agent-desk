"use client"

import {
  createContext,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react"
import { LanguagesIcon } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { OptionCombobox } from "@/components/option-combobox"
import { ImMessageHTML } from "@/components/im-message-html"
import { useI18n } from "@/i18n/provider"
import { readSession } from "@/lib/auth"
import {
  translateConversation,
  type AgentMessage,
  type ConversationTranslation,
  type TranslationLanguage,
} from "@/lib/api/agent"

type TranslationSettings = {
  reading: boolean
  writing: boolean
  readingLanguage: TranslationLanguage
  writingLanguage: TranslationLanguage
  display: "both" | "translation"
}
const defaults: TranslationSettings = {
  reading: false,
  writing: false,
  readingLanguage: "zh-CN",
  writingLanguage: "id",
  display: "both",
}
const languages: TranslationLanguage[] = ["zh-CN", "en", "id"]
const TranslationContext = createContext<{
  settings: TranslationSettings
  writingVersion: number
  update: (value: Partial<TranslationSettings>) => void
}>({ settings: defaults, writingVersion: 0, update: () => {} })
export const useTranslationSettings = () => useContext(TranslationContext)

export function TranslationProvider({
  accountId,
  children,
}: {
  accountId: number
  children: ReactNode
}) {
  const key = `agent-desk-translation:${accountId}`
  const [writingVersion, setWritingVersion] = useState(0)
  const [settings, setSettings] = useState<TranslationSettings>(() => {
    try {
      const saved = JSON.parse(localStorage.getItem(key) || "{}")
      return {
        reading: saved.reading === true,
        writing: saved.writing === true,
        readingLanguage: languages.includes(saved.readingLanguage)
          ? saved.readingLanguage
          : defaults.readingLanguage,
        writingLanguage: languages.includes(saved.writingLanguage)
          ? saved.writingLanguage
          : defaults.writingLanguage,
        display: saved.display === "translation" ? "translation" : "both",
      }
    } catch {
      return defaults
    }
  })
  function update(next: Partial<TranslationSettings>) {
    if ("writing" in next || "writingLanguage" in next)
      setWritingVersion((version) => version + 1)
    setSettings((previous) => {
      const value = { ...previous, ...next }
      try {
        localStorage.setItem(key, JSON.stringify(value))
      } catch {
        /* Private browsing may disable storage. */
      }
      return value
    })
  }
  return (
    <TranslationContext.Provider value={{ settings, writingVersion, update }}>
      {children}
    </TranslationContext.Provider>
  )
}

export function TranslationToolbar() {
  const t = useI18n()
  const { settings, update } = useTranslationSettings()
  const options = languages.map((value) => ({
    value,
    label: t(`translation.language.${value}`),
  }))
  return (
    <div className="flex shrink-0 justify-end border-b px-3 py-1">
      <Popover>
        <PopoverTrigger render={<Button size="sm" variant="ghost" />}>
          <LanguagesIcon className="size-4" />
          {t("translation.title")}
          {settings.reading || settings.writing ? (
            <span className="size-1.5 rounded-full bg-emerald-500" />
          ) : null}
        </PopoverTrigger>
        <PopoverContent
          align="end"
          className="w-[min(22rem,calc(100vw-2rem))] space-y-4"
        >
          <p className="font-medium">{t("translation.title")}</p>
          <label className="flex items-center gap-2 text-sm">
            <Checkbox
              checked={settings.reading}
              onCheckedChange={(checked) =>
                update({ reading: checked === true })
              }
            />
            {t("translation.autoReading")}
          </label>
          <div className="space-y-1 text-sm">
            <p>{t("translation.readingLanguage")}</p>
            <OptionCombobox
              options={options}
              value={settings.readingLanguage}
              onChange={(value) =>
                update({ readingLanguage: value as TranslationLanguage })
              }
              placeholder={t("translation.readingLanguage")}
            />
          </div>
          <div className="space-y-1 text-sm">
            <p>{t("translation.display")}</p>
            <OptionCombobox
              options={[
                { value: "both", label: t("translation.both") },
                {
                  value: "translation",
                  label: t("translation.onlyTranslation"),
                },
              ]}
              value={settings.display}
              onChange={(value) =>
                update({ display: value as TranslationSettings["display"] })
              }
              placeholder={t("translation.display")}
            />
          </div>
          <label className="flex items-center gap-2 text-sm">
            <Checkbox
              checked={settings.writing}
              onCheckedChange={(checked) =>
                update({ writing: checked === true })
              }
            />
            {t("translation.autoWriting")}
          </label>
          <div className="space-y-1 text-sm">
            <p>{t("translation.writingLanguage")}</p>
            <OptionCombobox
              options={options}
              value={settings.writingLanguage}
              onChange={(value) =>
                update({ writingLanguage: value as TranslationLanguage })
              }
              placeholder={t("translation.writingLanguage")}
            />
          </div>
          <p className="text-xs text-muted-foreground">
            {t("translation.settingsHint")}
          </p>
        </PopoverContent>
      </Popover>
    </div>
  )
}

type ReadingJob = {
  controller: AbortController
  users: number
  start: () => void
  promise: Promise<ConversationTranslation>
}
const readingCache = new Map<string, ConversationTranslation>()
const readingJobs = new Map<string, ReadingJob>()
const readingQueue: ReadingJob[] = []
let running = 0
let cacheSession = ""
function drainReadingQueue() {
  while (running < 2 && readingQueue.length) {
    const job = readingQueue.shift()!
    if (!job.controller.signal.aborted) job.start()
  }
}
function acquireReading(message: AgentMessage, language: TranslationLanguage) {
  const session = readSession()?.accessToken ?? ""
  if (cacheSession !== session) {
    readingCache.clear()
    for (const job of readingJobs.values()) job.controller.abort()
    readingJobs.clear()
    cacheSession = session
  }
  const key = JSON.stringify([
    message.conversationId,
    message.id,
    message.content,
    message.payload,
    language,
  ])
  const cached = readingCache.get(key)
  if (cached) return { promise: Promise.resolve(cached), release: () => {} }
  let job = readingJobs.get(key)
  if (!job) {
    const controller = new AbortController()
    let resolve!: (value: ConversationTranslation) => void
    let reject!: (reason: unknown) => void
    const promise = new Promise<ConversationTranslation>((yes, no) => {
      resolve = yes
      reject = no
    })
    const created: ReadingJob = {
      controller,
      users: 0,
      promise,
      start: () => {
        running++
        void translateConversation(
          {
            conversationId: message.conversationId,
            messageId: message.id,
            targetLanguage: language,
          },
          controller.signal,
        )
          .then((value) => {
            if (!controller.signal.aborted && cacheSession === session) {
              if (readingCache.size >= 500)
                readingCache.delete(readingCache.keys().next().value!)
              readingCache.set(key, value)
            }
            resolve(value)
          }, reject)
          .finally(() => {
            if (readingJobs.get(key) === created) readingJobs.delete(key)
            running--
            drainReadingQueue()
          })
      },
    }
    controller.signal.addEventListener(
      "abort",
      () => reject(new DOMException("Cancelled", "AbortError")),
      { once: true },
    )
    readingJobs.set(key, created)
    readingQueue.push(created)
    job = created
  }
  job.users++
  drainReadingQueue()
  return {
    promise: job.promise,
    release: () => {
      if (--job.users === 0 && readingJobs.get(key) === job) {
        readingJobs.delete(key)
        job.controller.abort()
      }
    },
  }
}

export function visibleTranslationText(html: string) {
  const doc = new DOMParser().parseFromString(html, "text/html")
  return doc.body.textContent?.trim() ?? ""
}

// Reuse the original image resources for display; model output never supplies HTML or URLs.
function readingHTML(translated: string, original: string) {
  const doc = new DOMParser().parseFromString(translated, "text/html")
  const source = new DOMParser().parseFromString(original, "text/html")
  const images = Array.from(source.querySelectorAll("img"))
  doc.querySelectorAll("img").forEach((image, index) => {
    if (images[index]) image.replaceWith(images[index].cloneNode(true))
  })
  return doc.body.innerHTML
}

export function ReadingTranslation({
  message,
  html,
  children,
  onImageSettled,
  onImageClick,
}: {
  message: AgentMessage
  html: string
  children: ReactNode
  onImageSettled: () => void
  onImageClick: (src: string, alt?: string) => void
}) {
  const t = useI18n()
  const { settings } = useTranslationSettings()
  const ref = useRef<HTMLDivElement>(null)
  const [visible, setVisible] = useState(false)
  const [attempt, setAttempt] = useState(0)
  const [state, setState] = useState<{
    key: string
    value?: ConversationTranslation
    error?: string
  }>()
  const eligible =
    settings.reading &&
    message.senderType === "customer" &&
    !message.recalledAt &&
    message.sendStatus !== 6 &&
    ["text", "html"].includes(message.messageType)
  const key = JSON.stringify([
    message.conversationId,
    message.id,
    message.content,
    message.payload,
    settings.readingLanguage,
  ])
  const current = state?.key === key ? state : undefined
  useEffect(() => {
    if (!ref.current) return
    const observer = new IntersectionObserver(([entry]) =>
      setVisible(entry.isIntersecting),
    )
    observer.observe(ref.current)
    return () => observer.disconnect()
  }, [])
  useEffect(() => {
    if (!eligible || !visible || !visibleTranslationText(html)) return
    let cancelled = false
    const task = acquireReading(message, settings.readingLanguage)
    void task.promise
      .then((value) => {
        if (!cancelled) setState({ key, value })
      })
      .catch((error: unknown) => {
        if (!cancelled)
          setState({
            key,
            error:
              error instanceof Error ? error.message : t("translation.failed"),
          })
      })
    return () => {
      cancelled = true
      task.release()
    }
  }, [
    eligible,
    visible,
    key,
    attempt,
    html,
    message,
    settings.readingLanguage,
    t,
  ])
  const value = eligible ? current?.value : undefined
  return (
    <div ref={ref}>
      {value && settings.display === "translation" ? (
        <>
          {value.translatedHTML ? (
            <ImMessageHTML
              html={readingHTML(value.translatedHTML, html)}
              onImageSettled={onImageSettled}
              onImageClick={onImageClick}
            />
          ) : (
            <p className="break-words whitespace-pre-wrap">{value.translatedText}</p>
          )}
          <details className="mt-2 text-xs">
            <summary className="cursor-pointer text-muted-foreground">
              {t("translation.original")}
            </summary>
            <div className="mt-1 text-sm">{children}</div>
          </details>
        </>
      ) : (
        children
      )}
      {eligible && (value || current?.error || visible) && (
        <div
          className="mt-2 border-t border-border/60 pt-2 text-xs text-muted-foreground"
          aria-live="polite"
        >
          {value ? (
            settings.display === "both" ? (
              <>
                <p>{t("translation.translated")}</p>
                  <p className="mt-1 break-words whitespace-pre-wrap text-sm text-foreground">
                  {value.translatedText}
                </p>
              </>
            ) : null
          ) : current?.error ? (
            <>
              <span>{current.error}</span>
              <Button
                size="sm"
                variant="ghost"
                onClick={() => {
                  setState(undefined)
                  setAttempt((n) => n + 1)
                }}
              >
                {t("translation.retry")}
              </Button>
            </>
          ) : visibleTranslationText(html) ? (
            t("translation.translating")
          ) : null}
        </div>
      )}
    </div>
  )
}

export function useDraftTranslation(conversationId: number) {
  const t = useI18n()
  const { settings, writingVersion } = useTranslationSettings()
  const [draft, setDraft] = useState({ html: "", composing: false })
  const [attempt, setAttempt] = useState(0)
  const [state, setState] = useState<{
    key: string
    value?: ConversationTranslation
    error?: string
  }>()
  const language = settings.writingLanguage
  const key = JSON.stringify([
    conversationId,
    draft.html,
    language,
    writingVersion,
  ])
  const enabled = settings.writing && Boolean(draft.html)
  const hasText =
    enabled &&
    typeof window !== "undefined" &&
    Boolean(visibleTranslationText(draft.html))
  const current = state?.key === key ? state : undefined
  useEffect(() => {
    if (!hasText || draft.composing) return
    const controller = new AbortController()
    const timer = setTimeout(() => {
      void translateConversation(
        {
          conversationId,
          text: draft.html,
          format: "html",
          targetLanguage: language,
        },
        controller.signal,
      )
        .then((value) => {
          if (!controller.signal.aborted) setState({ key, value })
        })
        .catch((error: unknown) => {
          if (!controller.signal.aborted)
            setState({
              key,
              error:
                error instanceof Error
                  ? error.message
                  : t("translation.failed"),
            })
        })
    }, 700)
    return () => {
      clearTimeout(timer)
      controller.abort()
    }
  }, [
    hasText,
    draft.html,
    draft.composing,
    conversationId,
    language,
    key,
    attempt,
    t,
  ])
  return {
    onDraftChange: (html: string, composing: boolean) =>
      setDraft((old) =>
        old.html === html && old.composing === composing
          ? old
          : { html, composing },
      ),
    sendOverride: settings.writing
      ? (html: string) => {
          if (!visibleTranslationText(html)) return html
          if (
            draft.composing ||
            html !== draft.html ||
            !current?.value?.translatedHTML
          )
            return null
          return current.value.translatedHTML
        }
      : undefined,
    sendDisabled:
      hasText && (draft.composing || !current?.value?.translatedHTML),
    sendLabel: settings.writing ? t("translation.sendTranslated") : undefined,
    footer: hasText ? (
      <div
        className="max-h-32 shrink-0 overflow-y-auto border-t bg-muted/30 px-3 py-2 text-xs"
        aria-live="polite"
      >
        <p className="mb-1 text-muted-foreground">
          {t("translation.preview")} · {t(`translation.language.${language}`)}
        </p>
        {current?.value ? (
          <p className="break-words whitespace-pre-wrap text-sm">
            {current.value.translatedText}
          </p>
        ) : current?.error ? (
          <>
            <span className="text-destructive">{current.error}</span>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => {
                setState(undefined)
                setAttempt((n) => n + 1)
              }}
            >
              {t("translation.retry")}
            </Button>
          </>
        ) : (
          <span>{t("translation.translating")}</span>
        )}
      </div>
    ) : null,
  }
}
