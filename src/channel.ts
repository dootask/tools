import type { Channel } from "./types"

export const PROTOCOL_VERSION = 2

export type EnvelopeKind = "hello" | "welcome" | "call" | "result" | "event" | "reply" | "callback" | "callback_result"

/** 通道里传的消息 */
export interface Envelope {
  dootask: typeof PROTOCOL_VERSION
  kind: EnvelopeKind
  [key: string]: any
}

/** 一种通信方式：只负责把信封送到 DooTask，并把 DooTask 发来的信封交出来 */
export interface Transport {
  channel: Channel
  /** 等 welcome 的时间（毫秒），不填用默认 */
  helloTimeout?: number
  send(message: Envelope): void
  /** fail：通道建立不起来（如桥接页换不到凭证），带上原因 */
  listen(receive: (message: Envelope) => void, fail: (message: string) => void): void
}

export function isEnvelope(value: unknown): value is Envelope {
  return (
    !!value &&
    typeof value === "object" &&
    (value as Envelope).dootask === PROTOCOL_VERSION &&
    typeof (value as Envelope).kind === "string"
  )
}

/**
 * frame：DooTask 窗口里的 iframe
 * - 只认父窗口发来的消息；收到 welcome 后记下 DooTask 的 origin，之后只发给它、只认它
 * - hello 里没有敏感内容，发给任意 origin
 */
function frameTransport(): Transport {
  let origin = ""
  return {
    channel: "frame",
    send(message) {
      window.parent.postMessage(message, origin || "*")
    },
    listen(receive) {
      window.addEventListener("message", event => {
        if (event.source !== window.parent || !isEnvelope(event.data)) {
          return
        }
        if (!origin && event.data.kind === "welcome") {
          origin = event.origin
        }
        if (event.origin === origin) {
          receive(event.data)
        }
      })
    },
  }
}

/** Electron / App 打开的页面里注入的宿主对象 */
interface InjectedHost {
  channel: "electron" | "app"
  postMessage(message: string): void
}

declare global {
  interface Window {
    __dootaskHost?: InjectedHost
    __dootaskReceive?: (message: string) => void
  }
}

/**
 * electron / app：DooTask 客户端打开的页面里注入了 __dootaskHost
 * - 发：__dootaskHost.postMessage(JSON 字符串)；收：DooTask 调用 window.__dootaskReceive(JSON 字符串)
 * - 只有 origin 和应用登记的一致时 DooTask 才注入，这里不用再校验
 */
function injectTransport(host: InjectedHost): Transport {
  return {
    channel: host.channel,
    send(message) {
      host.postMessage(JSON.stringify(message))
    },
    listen(receive) {
      window.__dootaskReceive = raw => {
        let message: unknown
        try {
          message = JSON.parse(raw)
        } catch {
          return
        }
        if (isEnvelope(message)) {
          receive(message)
        }
      }
    },
  }
}

/** DooTask 在浏览器里打开外部应用时加在地址上的参数 */
const LAUNCH_PARAM = "dootask_launch"
const ORIGIN_PARAM = "dootask_origin"
/** 记下打开本页的 DooTask，页面内跳转、刷新后沿用 */
const ORIGIN_STORAGE_KEY = "__dootask_bridge_origin__"

/** 只有轻提示时桥接页占顶部的高度 */
const TOAST_BAND_HEIGHT = 120

function readStorage(key: string): string {
  try {
    return localStorage.getItem(key) ?? ""
  } catch {
    return ""
  }
}

function writeStorage(key: string, value: string) {
  try {
    localStorage.setItem(key, value)
  } catch {
    // 存不了就只在本次页面有效
  }
}

function normalizeOrigin(value: string): string {
  try {
    const url = new URL(value)
    return /^https?:$/.test(url.protocol) ? url.origin : ""
  } catch {
    return ""
  }
}

/**
 * 取出地址上的 launch_code 和 DooTask origin，并从地址栏去掉（不留在历史、书签和分享的链接里）；
 * 没有时沿用上次记下的 origin（code 为空，桥接页用自己存的凭证）
 */
function takeBridgeLaunch(): { origin: string; code: string } | null {
  const url = new URL(window.location.href)
  const code = url.searchParams.get(LAUNCH_PARAM) ?? ""
  const origin = normalizeOrigin(url.searchParams.get(ORIGIN_PARAM) ?? "")
  if (url.searchParams.has(LAUNCH_PARAM) || url.searchParams.has(ORIGIN_PARAM)) {
    url.searchParams.delete(LAUNCH_PARAM)
    url.searchParams.delete(ORIGIN_PARAM)
    window.history.replaceState(window.history.state, "", url.href)
  }
  if (code && origin) {
    writeStorage(ORIGIN_STORAGE_KEY, origin)
    return { origin, code }
  }
  const saved = readStorage(ORIGIN_STORAGE_KEY)
  return saved ? { origin: saved, code: "" } : null
}

function whenBody(run: (body: HTMLElement) => void) {
  if (document.body) {
    run(document.body)
  } else {
    document.addEventListener("DOMContentLoaded", () => run(document.body), { once: true })
  }
}

/**
 * bridge：浏览器里由 DooTask 打开的外部页面
 * - 插入隐藏的 iframe（DooTask 的 /bridge），加载后先发 launch（带 launch_code），协议消息之后照常收发
 * - 只和这个 iframe、这个 origin 通信；launch_code 只经 postMessage 交给它，不出现在 iframe 地址上
 * - 桥接页告诉显示方式（overlay）：有弹窗时铺满，只有轻提示时占顶部，否则隐藏且不挡点击
 * - 桥接页弹不出新窗口（会被拦截），要打开的地址（下载、主程序深度链接）由本页打开
 */
function bridgeTransport(origin: string, code: string): Transport {
  const frame = document.createElement("iframe")
  frame.src = `${origin}/bridge`
  frame.title = "DooTask"
  frame.setAttribute("aria-hidden", "true")
  Object.assign(frame.style, {
    position: "fixed",
    top: "0",
    left: "0",
    width: "100%",
    height: "100%",
    border: "0",
    margin: "0",
    padding: "0",
    zIndex: "2147483647",
    background: "transparent",
    colorScheme: "normal",
  })
  const setOverlay = (mode: unknown) => {
    const shown = mode === "full" || mode === "toast"
    frame.style.visibility = shown ? "visible" : "hidden"
    frame.style.pointerEvents = shown ? "auto" : "none"
    frame.style.height = mode === "toast" ? `${TOAST_BAND_HEIGHT}px` : "100%"
    frame.setAttribute("aria-hidden", shown ? "false" : "true")
  }
  setOverlay("hidden")

  let loaded = false
  const queue: Record<string, any>[] = []
  const post = (message: Record<string, any>) => frame.contentWindow?.postMessage(message, origin)
  frame.addEventListener("load", () => {
    if (loaded) {
      return
    }
    loaded = true
    post({ dootask: PROTOCOL_VERSION, kind: "bridge", action: "launch", code })
    queue.splice(0).forEach(post)
  })
  whenBody(body => body.appendChild(frame))

  return {
    channel: "bridge",
    helloTimeout: 30000,
    send(message) {
      if (loaded) {
        post(message)
      } else {
        queue.push(message)
      }
    },
    listen(receive, fail) {
      window.addEventListener("message", event => {
        if (event.source !== frame.contentWindow || event.origin !== origin) {
          return
        }
        const data = event.data
        if (data?.dootask !== PROTOCOL_VERSION) {
          return
        }
        if (data.kind !== "bridge") {
          if (isEnvelope(data)) {
            receive(data)
          }
          return
        }
        switch (data.action) {
          case "overlay":
            setOverlay(data.mode)
            break
          case "open":
            if (typeof data.url === "string" && /^https?:\/\//i.test(data.url)) {
              window.open(data.url, "_blank", "noopener")
            }
            break
          case "fail":
            setOverlay("hidden")
            fail(String(data.message ?? ""))
            break
        }
      })
    },
  }
}

/** 按所在环境选择通道（注入优先于父窗口，再看是不是浏览器里由 DooTask 打开的）；不在 DooTask 里返回 null */
export function detectTransport(): Transport | null {
  if (typeof window === "undefined") {
    return null
  }
  if (window.__dootaskHost) {
    return injectTransport(window.__dootaskHost)
  }
  if (window.parent && window.parent !== window) {
    return frameTransport()
  }
  const launch = takeBridgeLaunch()
  if (launch) {
    return bridgeTransport(launch.origin, launch.code)
  }
  return null
}
