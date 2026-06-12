# @dootask/cli

DooTask 的命令行工具 `doo`（gh 式 CLI）的 npm 分发渠道。

## 安装

```bash
# 全局
npm i -g @dootask/cli
# 或
yarn global add @dootask/cli

# 项目内
npm i @dootask/cli
```

npm/yarn 会按 `optionalDependencies` 与 `os`/`cpu` 字段自动只装与当前平台匹配的二进制子包（约 8 MB）。

## 使用

```bash
doo --help
doo auth login --server https://your-dootask.com --email you@example.com
doo task list --project 130
doo app catalog --search 客户管理
```

详见 [doo CLI 文档](https://github.com/dootask/tools/blob/main/server/go/cmd/doo/README.md)。

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
