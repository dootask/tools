# doo — DooTask 命令行工具

`doo` 是 DooTask 的 `gh` 式命令行工具，建在官方 Go SDK（`github.com/dootask/tools/server/go`）之上，覆盖任务、项目、看板列、对话、消息、群组、用户、机器人、文件、报告、统一搜索等操作。每条命令以**当前用户身份与权限**执行（权限由 DooTask 主程序校验）。

## 构建

```bash
make build           # 产出 dist/doo
make build-all       # 交叉编译 linux/darwin/windows
make test            # 单元测试
VERSION=1.0.0 make build
```

需要 Go 1.22+。

## 配置与登录

凭证优先级：命令行 flag > 环境变量 > 配置文件 `~/.config/doo/config.json`。

```bash
# 方式一：登录并保存（写入 0600 配置文件）
doo auth login --server https://your-dootask.com --email you@example.com
doo auth status
doo auth logout

# 方式二：环境变量（适合 CI / 容器，不落盘）
export DOO_SERVER=https://your-dootask.com
export DOO_TOKEN=<token>
```

> 若实例开启了登录验证码，`auth login` 无法完成，请在浏览器登录后用 `--token` / `DOO_TOKEN` 直接传入。

## 全局参数

| flag | 说明 |
|---|---|
| `--server` | DooTask 服务器地址 |
| `--token` | API token |
| `--json` | 以 JSON 信封输出（适合脚本/程序解析，见下文「机器调用规范」） |
| `--yes, -y` | 跳过危险操作确认 |
| `--quiet, -q` | 精简输出 |
| `--dry-run` | 预演：校验参数并返回将执行的操作，不实际执行（只读命令忽略） |
| `--fields a,b` | 只输出指定字段（列表逐项裁剪；分页对象裁剪 `data` 内各项） |
| `--limit N` | 列表最多输出 N 条（`file fetch` 自身的 `--limit` 为内容长度上限） |
| `--jq 表达式` | 用内置 jq 处理输出（在 `--fields`/`--limit` 之后执行） |

默认输出为人类可读表格；列表过宽的单元格会折叠换行并截断，完整数据请用 `--json`。

## 机器调用规范（`--json`）

**输出信封**（stdout 恰好一行 JSON）：

```json
{"ok":true,"data":{...},"message":"✓ 已删除任务 #123"}
{"ok":false,"code":"not_found","message":"任务不存在"}
{"ok":false,"code":"needs_confirm","message":"…","confirm":{"summary":"删除任务 #123","command":"task delete","risk":"danger","targets":[{"type":"task","id":123}]}}
```

`message` 仅在写操作有提示文案时出现；`--dry-run` 成功时 `data` 为 `{"dry_run":true,"command","risk","args","flags"}`。

**错误码**（固定清单）：

| code | 含义 | 退出码 |
|---|---|---|
| `invalid_args` | 参数/命令不合法 | 2 |
| `auth_required` | 未登录或 token 失效 | 3 |
| `needs_confirm` | 危险操作需确认（见 `confirm` 字段） | 4 |
| `permission_denied` | 无权限 | 1 |
| `not_found` | 对象不存在或不可见 | 1 |
| `api_error` | 主程序业务错误（`message` 为原因） | 1 |
| `server_error` | 主程序 HTTP 非 200 | 1 |
| `network_error` / `timeout` | 网络不可达 / 请求超时 | 1 |
| `cancelled` | 交互确认时用户取消 | 1 |
| `error` | 其他本地错误 | 1 |

**风险级别**：每个命令为 `read` / `write` / `danger` 之一，`doo meta` 输出完整清单（命令、用法、参数、风险级别）。
`danger` 命令在 `--json` 或非交互环境下未带 `--yes` 时**不执行**，返回 `needs_confirm`；调用方向用户确认后带 `--yes` 重跑。
带 `--recover` 的恢复类操作按 `write` 处理。

## 命令一览

```
doo auth      login | status | logout
doo task      list | view | files | create | subtask | update | done | undone | dialog | notify | archive | delete
doo project   list | view | create | update | exit | delete
doo column    list | create | update | delete
doo dialog    list | search | view | users
doo message   send | send-user | list | search | view | withdraw | forward | todo | done
doo group     create | edit | add-user | remove-user | exit | transfer | disband
doo user      info | departments | basic | search
doo bot       list | view | create | update | delete
doo file      list | search | view | fetch | history | create | save | rename | restore | delete   (实验性)
doo report    received | my | view | template | submit | mark   (实验性)
doo search    <关键词> [--types ...]                 (实验性)
doo page      context | action | element             (需 --session <fd>)
doo app       list | updates | catalog [--search 词] | fields <ID> | upload <zip> [--appid] | install <ID> [...] | update <ID> [...] | reinstall <ID> [...]
              | uninstall <ID> [--delete-data] | remove <ID> | logs <ID> | containers <ID> | container-logs <ID> --service | refresh
              | call <ID> [指令] [K=V ...]
doo system    version | settings
doo meta                                             (机器可读命令清单)
```

用 `doo <名词> --help`、`doo <名词> <动词> --help` 逐层查看参数。

## 示例

```bash
doo task list --project 130 --status uncompleted
doo task create --project 130 --name "写周报" --owner 3 --end "2026-06-20 18:00:00"
doo task done 38001
doo task update 38001 --content "进展更新"        # 仅提交改动字段，不会清空其它字段
doo project list --jq '.data[].name'
doo message send --dialog 2889 --text "下班啦" --silence
doo search 财务 --types task,project

# 应用插件（AppStore）
doo app catalog --search 客户管理                  # 中文/英文/tag 模糊匹配 id/name/description/tags
doo app fields community_kuaifan_crm               # 装前先查参数定义（必填/默认/可选项）
doo app install ai --param EMBEDDING_MODEL=bge-m3  # 已装应用未传 --param 自动 sticky，不会清空令牌
doo app install ai --cpu-limit 1.0 --memory-limit 2G --pull
doo app reinstall ai                               # 按当前已装版本与参数重部署
doo app uninstall community_kuaifan_mcp --delete-data --yes

# 调用应用自定义指令（基于应用暴露的 OpenAPI/Swagger）
doo app call approve                               # 列出该应用可用指令
doo app call approve api-list --describe           # 查看某指令的参数
doo app call approve api-list page=2               # 调用指令（也可 --arg page=2）
doo app call approve create-item title=foo --data '{"score":9}'
doo app call approve --raw GET /api/custom         # 直接调用未在规范中声明的端点
```

## 说明

- 危险/不可逆操作（删除、解散群、撤回消息等）默认需要确认；非交互环境请显式加 `--yes`（`--json` 下返回 `needs_confirm`）。
- `file` / `report` / `search` 暂走通用端点（SDK 尚无对应类型），标记为实验性，输出字段以 `--json` 为准。
- `app`（应用插件）走 AppStore 微服务（主程序反代 `/appstore/api/v1`，响应 `{code,message,data}`，与主程序 `{ret,msg,data}` 不同；请求自动带 `Version` 头供 AppStore 校验 `require_version`）：
  - `install`/`update`/`reinstall`/`uninstall`/`remove`/`refresh` 需**管理员**权限，安装/卸载会触发 docker compose、可能耗时；`list`/`catalog`/`fields`/`logs`/`containers` 普通用户即可。
  - `updates`：列出**已安装且有新版可升级**的应用（取 `/list?include=all` 中 `upgradeable=true`，与网页「可升级」徽标同源，含 community 应用）。区别于 `refresh`——`refresh` 是刷新远程源/包，`updates` 看的是已装应用能否升级。
  - `upload <zip>`：上传本地应用压缩包导入应用市场（multipart 传到 `/internal/apps/upload`，`--appid` 留空则从文件名自动提取）。**仅导入 + 合规校验，不会自动部署**；导入成功后用 `doo app install <ID>` 才真正安装（与网页「上传本地应用」一致：先导入再安装）。`--appid` 留空时后端从**文件名**提取（去结尾版本段与扩展名，如 `myapp-1.0.0.tar.gz` → `myapp`）；若文件名不等于目标应用 ID，请显式传 `--appid`。
  - `catalog --search <kw>` 在 `id`/`name`/`description`/`tags` 上做大小写不敏感的子串匹配，覆盖中英文 tag（如「客户管理」）。
  - 装/升前先用 `doo app fields <ID>` 查参数定义；`--param K=V` 可重复；fields 中不存在的 key 直接报错，必填字段缺失且无默认值时拒绝提交。
  - **sticky**：已安装应用未传 `--param` 自动沿用当前 `params`、未传 `--cpu-limit`/`--memory-limit` 自动沿用当前 `resources`，与网页表单"初值即当前值"行为一致，避免 `KB_INGEST_TOKEN` 等令牌被误清。
  - `call <ID> [指令]`：调用应用**自己的后端接口**。AppStore 读取应用暴露的 OpenAPI/Swagger 规范（`config.yml` 的 `openapi` 字段声明，或自动探测），归一化成指令目录;调用经 AppStore 内网代理到应用容器，自动透传当前用户 `Token` 与身份头（`X-Doo-User-*`），由应用后端自鉴权。无指令时列出目录,`--describe` 看参数,`--raw <METHOD> <path>` 调未声明端点,`--data` 传完整 JSON 体,`--refresh` 跳过规范缓存。详见插件开发文档「应用指令」章节。
- `doo page`（获取页面上下文 / 执行业务操作 / 操作页面元素）经主程序常驻 WebSocket（`/ws`）派发到用户浏览器执行：CLI 调 `assistant/operation/dispatch` 派发后轮询 `assistant/operation/result` 取结果，对调用者表现为同步命令。需用 `--session <fd>`（或环境变量 `DOO_SESSION`）指定目标会话；fd 为用户当前在线的 WebSocket 连接，归属与在线由主程序校验。
