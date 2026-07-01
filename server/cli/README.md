# @dootask/cli

DooTask 的命令行工具 `doo`（gh 式 CLI）的 npm 分发渠道。既可供人使用，也**面向 AI Agent 设计**——结构化输出、明确退出码、环境变量注入认证、危险操作护栏。

## 安装

```bash
# 全局
npm i -g @dootask/cli
# 或
yarn global add @dootask/cli

# 项目内
npm i @dootask/cli

# 免安装一次性运行
npx @dootask/cli --help
```

npm/yarn 会按 `optionalDependencies` 与 `os`/`cpu` 字段自动只装与当前平台匹配的二进制子包（约 8 MB）。

## 使用

```bash
doo --help
doo auth login --server https://your-dootask.com --email you@example.com
doo task list --project 130
doo app catalog --search 客户管理
```

**认证优先级**：`--token` / `--server` 命令行参数 > `DOO_TOKEN` / `DOO_SERVER` 环境变量 > 配置文件（`doo auth login` 写入 `~/.config/doo/config.json`）。

更完整的命令说明见 [doo CLI 文档](https://github.com/dootask/tools/blob/main/server/go/cmd/doo/README.md)。

---

## For AI Agents

`doo` 适合作为工具交给 LLM Agent。把下面这段贴进 system prompt 即可接入：

```
你可以使用 `doo` 命令行调用 DooTask。约定：
- 认证已通过 DOO_SERVER / DOO_TOKEN 环境变量注入，无需登录。
- 命令都加 --json 取机器可读输出（需纯净再加 --quiet）。
- 退出码：0=成功；3=未认证；1=其它错误。数据在 stdout，错误在 stderr（"doo:" 前缀）。
- 删除等危险操作先征得用户同意，再带 --yes 重发（非交互不加会被拒绝）。
- 用 `doo --help` / `doo <命令> --help` 自查参数，不要臆造。
先运行 `doo user info --json` 确认认证可用。
```

封装成工具时：argv 传数组、认证仅经环境变量、不要开放 `auth` 子命令。生产参考实现见 DooTask AI 插件的 [`doo_tool.py`](https://github.com/dootask/system-plugins/blob/main/ai/src/helper/doo_tool.py)。

---

## 平台支持

| 子包 | OS | CPU |
|---|---|---|
| `@dootask/cli-linux-x64` | linux | x64 |
| `@dootask/cli-linux-arm64` | linux | arm64 |
| `@dootask/cli-darwin-x64` | darwin | x64 (Intel Mac) |
| `@dootask/cli-darwin-arm64` | darwin | arm64 (Apple Silicon) |
| `@dootask/cli-win32-x64` | win32 | x64 |

其他平台请从 [源码构建](https://github.com/dootask/tools/tree/main/server/go/cmd/doo)。

## 源码

`doo` 本身用 Go 编写，源码在 [`server/go/cmd/doo/`](../go/cmd/doo)。本 npm 包仅做分发。
