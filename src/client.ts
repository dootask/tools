/**
 * 和 DooTask 的连接：握手、调用、事件、回调（协议第 2 版）
 * - 加载后自动选择通道并发 hello，收到 welcome 后就绪；HELLO_TIMEOUT 内没有 welcome 视为不在 DooTask 里
 * - 参数里的函数换成 {__fn: id}，DooTask 用 callback 调回来
 * - welcome 里的状态（主题、语言、安全距离……）存在本地，随事件更新
 */
import { DooTaskError, fromPayload, toPayload } from "./errors"
import { detectTransport, PROTOCOL_VERSION, type Envelope, type EnvelopeKind, type Transport } from "./channel"
import type { Env } from "./types"

/** 工具库版本，hello 时告诉 DooTask */
export const SDK_VERSION = "2.0.0"

const HELLO_TIMEOUT = 10000

/** 离过期不到这么久就换新令牌 */
const TOKEN_RENEW_BEFORE = 10 * 60 * 1000

export interface Session {
  appid: string
  key: string
  env: Env
  state: Record<string, any>
  capabilities: string[]
}

type Listener = (data: any) => unknown

interface Pending {
  resolve: (value: any) => void
  reject: (error: unknown) => void
}

const transport: Transport | null = detectTransport()
const pendingCalls = new Map<string, Pending>()
const functions = new Map<string, (...args: any[]) => unknown>()
const listeners = new Map<string, Set<Listener>>()
let seq = 0
let session: Session | null = null
let token = { value: "", expiredAt: 0 }
let readyPromise: Promise<Session> | null = null

function nextId(prefix: string): string {
  seq += 1
  return `${prefix}_${seq}`
}

function send(kind: EnvelopeKind, body: Record<string, any>) {
  transport?.send({ dootask: PROTOCOL_VERSION, kind, ...body })
}

/** 参数里的函数换成 {__fn: id}，记下来等 DooTask 调回 */
function serialize(value: any): any {
  if (typeof value === "function") {
    const id = nextId("fn")
    functions.set(id, value)
    return { __fn: id }
  }
  if (Array.isArray(value)) {
    return value.map(serialize)
  }
  if (value && typeof value === "object" && Object.getPrototypeOf(value) === Object.prototype) {
    const result: Record<string, any> = {}
    Object.keys(value).forEach(key => {
      if (value[key] !== undefined) {
        result[key] = serialize(value[key])
      }
    })
    return result
  }
  return value
}

async function onCallback(message: Envelope) {
  const fn = functions.get(String(message.fn))
  try {
    if (!fn) {
      throw new DooTaskError("unsupported", "callback is gone")
    }
    const value = await fn(...(Array.isArray(message.args) ? message.args : []))
    send("callback_result", { id: message.id, value })
  } catch (error) {
    send("callback_result", { id: message.id, error: toPayload(error) })
  }
}

async function onEvent(message: Envelope) {
  const name = String(message.name)
  if (session) {
    if (name === "theme" || name === "language" || name === "safe_area" || name === "presentation") {
      session.state[name] = message.data
      if (name === "presentation") {
        session.env.presentation = message.data
      }
    } else if (name === "token" && message.data) {
      token = { value: String(message.data.token), expiredAt: Number(message.data.expired_at) || 0 }
    }
  }
  const handlers = Array.from(listeners.get(name) ?? [])
  if (message.id === undefined) {
    handlers.forEach(handler => {
      try {
        handler(message.data)
      } catch (error) {
        console.error("[dootask]", error)
      }
    })
    return
  }
  // 需要回复的事件（before_close）：依次问，有一个返回 false 或提示文字就回它（DooTask 弹确认），抛错的跳过
  let value: unknown = true
  for (const handler of handlers) {
    try {
      const result = await handler(message.data)
      if (result === false || typeof result === "string") {
        value = result
        break
      }
    } catch (error) {
      console.error("[dootask]", error)
    }
  }
  send("reply", { id: message.id, value })
}

function connect(): Promise<Session> {
  if (!transport) {
    return Promise.reject(new DooTaskError("unsupported", "not running in DooTask"))
  }
  const channel = transport
  return new Promise((resolve, reject) => {
    const timer = setTimeout(() => {
      reject(new DooTaskError("unsupported", "no response from DooTask"))
    }, channel.helloTimeout ?? HELLO_TIMEOUT)
    const fail = (reason: string) => {
      clearTimeout(timer)
      reject(new DooTaskError("auth_failed", reason))
    }
    channel.listen(message => {
      switch (message.kind) {
        case "welcome":
          clearTimeout(timer)
          if (message.protocol !== PROTOCOL_VERSION) {
            reject(new DooTaskError("version_mismatch", `DooTask speaks protocol ${message.protocol}`))
            return
          }
          token = { value: String(message.token ?? ""), expiredAt: Number(message.token_expired_at) || 0 }
          session = {
            appid: String(message.appid ?? ""),
            key: String(message.key ?? ""),
            env: { ...message.env, channel: channel.channel },
            state: message.state ?? {},
            capabilities: Array.isArray(message.capabilities) ? message.capabilities : [],
          }
          resolve(session)
          break
        case "result": {
          const pending = pendingCalls.get(String(message.id))
          if (pending) {
            pendingCalls.delete(String(message.id))
            if (message.error !== undefined) {
              pending.reject(fromPayload(message.error))
            } else {
              pending.resolve(message.value)
            }
          }
          break
        }
        case "event":
          void onEvent(message)
          break
        case "callback":
          void onCallback(message)
          break
      }
    }, fail)
    send("hello", { sdk: SDK_VERSION })
  })
}

/** 等待握手完成 */
export function ready(): Promise<Session> {
  if (!readyPromise) {
    readyPromise = connect()
  }
  return readyPromise
}

/** 调用 DooTask 的接口 */
export async function call<T = any>(method: string, ...args: any[]): Promise<T> {
  await ready()
  return new Promise<T>((resolve, reject) => {
    const id = nextId("c")
    pendingCalls.set(id, { resolve, reject })
    send("call", { id, method, args: serialize(args) })
  })
}

/** 登记事件监听，返回取消函数 */
export function listen(name: string, handler: Listener): () => void {
  let set = listeners.get(name)
  if (!set) {
    set = new Set()
    listeners.set(name, set)
  }
  set.add(handler)
  return () => {
    set?.delete(handler)
  }
}

/** 某个事件当前有几个监听 */
export function listenerCount(name: string): number {
  return listeners.get(name)?.size ?? 0
}

/** 应用令牌，快过期时先换新的 */
export async function appToken(): Promise<string> {
  await ready()
  if (token.expiredAt - Date.now() < TOKEN_RENEW_BEFORE) {
    const result = await call<{ token: string; expired_at: number }>("getAppToken")
    token = { value: result.token, expiredAt: Number(result.expired_at) || 0 }
  }
  return token.value
}
