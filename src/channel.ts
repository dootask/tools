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
  send(message: Envelope): void
  listen(receive: (message: Envelope) => void): void
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

/** 按所在环境选择通道（注入优先于父窗口）；不在 DooTask 里返回 null */
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
  return null
}
