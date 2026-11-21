/** 错误码，见 README「错误」 */
export type DooTaskErrorCode =
  | "unsupported"
  | "cancelled"
  | "permission_denied"
  | "api_error"
  | "invalid_args"
  | "timeout"
  | "version_mismatch"
  | "auth_failed"

const CODES: DooTaskErrorCode[] = [
  "unsupported",
  "cancelled",
  "permission_denied",
  "api_error",
  "invalid_args",
  "timeout",
  "version_mismatch",
  "auth_failed",
]

/** 工具库抛出的所有错误，用 code 区分 */
export class DooTaskError extends Error {
  readonly code: DooTaskErrorCode
  readonly data?: any

  constructor(code: DooTaskErrorCode, message = "", data?: any) {
    super(message || code)
    this.name = "DooTaskError"
    this.code = code
    this.data = data
    Object.setPrototypeOf(this, DooTaskError.prototype)
  }
}

/** 通道里收到的错误 {code, message, data} 转成 DooTaskError */
export function fromPayload(raw: unknown): DooTaskError {
  const error = raw && typeof raw === "object" ? (raw as Record<string, unknown>) : {}
  const code = CODES.includes(error.code as DooTaskErrorCode) ? (error.code as DooTaskErrorCode) : "api_error"
  return new DooTaskError(code, typeof error.message === "string" ? error.message : "", error.data)
}

/** 应用回调里抛出的错误转成通道里的格式 */
export function toPayload(error: unknown): { code: DooTaskErrorCode; message: string; data?: any } {
  if (error instanceof DooTaskError) {
    return error.data === undefined
      ? { code: error.code, message: error.message }
      : { code: error.code, message: error.message, data: error.data }
  }
  return { code: "api_error", message: error instanceof Error ? error.message : String(error ?? "") }
}
