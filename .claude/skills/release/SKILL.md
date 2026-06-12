---
name: release
description: 从 main 分支发布 dootask-tools 新版本：版本号 → CHANGELOG → commit → tag → push，Action 自动出三包（npm @dootask/tools + @dootask/cli + PyPI dootask-tools）与 doo 五平台二进制。刚性顺序、每步确认、失败即停。
---

# dootask-tools 发布流程

**刚性技能**——严格按顺序执行，每步向用户确认，任何一步失败立即停止。

## 核心原则

skill 在本地只做"准备 + 打 tag + push"；真正的构建发布全部由 `.github/workflows/release.yml` 在 GitHub Actions 上完成。**push 完不等于发布完成**——必须最后确认 Action 全绿。

三个包统一版本号：
- npm `@dootask/tools`（SDK，根 `package.json`）
- npm `@dootask/cli`（含 doo 二进制，`server/cli/package.json` + 5 个平台子包）
- PyPI `dootask-tools`（`server/python/setup.py` + `dootask/__init__.py`）
- Go `github.com/dootask/tools/server/go`（不发布，靠 git tag 被 proxy.golang.org 索引）

## 前置检查（全部通过才能继续）

执行任何发布步骤前，依次检查：

1. **分支**：必须是 `main`，否则停止，提示用户切换
2. **工作区**：`git status` 必须干净（无未提交变更、无未跟踪文件），否则**停止**并交由用户处理
3. **远程同步**：`git fetch origin && git status -sb` 显示无 ahead/behind，否则停止让用户 pull/push
4. **Node.js**：`node --version` 必须 ≥ 14（doo.js shim 用到 spawnSync）
5. **gh CLI**：`gh --version` 可用且 `gh auth status` 已登录（push 后要查 Action 状态）
6. **secrets 提示**：检查 `gh secret list -R dootask/tools` 是否含 `NPM_TOKEN` 与 `PYPI_API_TOKEN`（仅提示，没有就让用户在 GitHub 仓库设置里加，否则 Action 会失败）

检查通过后汇报结果，用户确认后再开始执行。

## 发布步骤

**每步执行前**向用户确认；**每步执行后**报告结果。

开始前先把这份清单复制到你的回复里，逐项勾选、跟踪进度：

```
发布进度：
- [ ] 前置检查（main / 干净 / 同步 / node≥14 / gh / secrets 提示）
- [ ] Step 1 决定新版本号
- [ ] Step 2 写入版本号到所有包文件
- [ ] Step 3 撰写 CHANGELOG
- [ ] 汇总变更 → 用户确认 → commit + tag + push
- [ ] 确认 GitHub Actions Release 工作流 success
- [ ] 验证三包真实可用（npm/pip/Release 页）
```

---

### Step 1: 决定新版本号

读最新 tag：

```shell
git tag -l --sort=-version:refname | head -3
```

向用户确认新版本号（默认 patch++；首次跨包统一版本时可能要直接定 `1.3.0`）。**必须符合 `vX.Y.Z` 语义化版本格式**。

> 注意 PyPI 不允许重发同号；如本号已发过（哪怕 yank），必须 bump 到下一个。

---

### Step 2: 写入版本号

要改的位置（所有用 `0.0.0-dev` 占位 → 改成正式版本号；python 是同步双写）：

```
package.json                                    "version"
server/cli/package.json                         "version" + 所有 optionalDependencies 值
server/cli/platforms/linux-x64/package.json     "version"
server/cli/platforms/linux-arm64/package.json   "version"
server/cli/platforms/darwin-x64/package.json    "version"
server/cli/platforms/darwin-arm64/package.json  "version"
server/cli/platforms/win32-x64/package.json     "version"
server/python/setup.py                          version="..."
server/python/dootask/__init__.py               __version__ = "..."
```

> 实际发布时 Action 也会用 `npm version --no-git-tag-version` 注入一遍——这一步是给"人/IDE/`go install` 等"看的，保持文件与 tag 一致。

完成后用 `git diff --stat` 给用户展示改动文件清单。

---

### Step 3: 撰写 CHANGELOG

读取本次区间的提交：

```shell
LAST_TAG=$(git tag -l --sort=-version:refname | head -1)
git log ${LAST_TAG}..HEAD --stat
```

按 `CHANGELOG.md` 现有格式（若没有就新建），在文件顶部 `# Changelog` 说明段之后、紧挨上一个 `## [...]` 之前，插入新版本区段：

```markdown
## [<version>]

### Features

- ...

### Bug Fixes

- ...

### Documentation

- ...
```

撰写要求：
- 小节标题用**英文 Title Case**：`Features` / `Bug Fixes` / `Performance` / `Documentation` / `Security` / `Miscellaneous`，**不要译成中文**；**没有内容的小节整段省略**。
- 条目正文用**通俗友好的简体中文**，面向开发者/使用者描述更新带来的好处。
- 过滤掉对用户无意义的提交（纯构建/依赖/CI/合并提交、本技能自身的脚手架改动等）。
- 仅凭提交标题无法判断时，结合提交的完整描述正文和实际代码改动（`git show <hash>`）再决定。
- 合并相似项；每个小节内按用户价值与影响范围排序，重要的在前。

展示新版本号与你写的 changelog 区段，请用户过目。

---

## 最终：commit + tag + push

所有步骤完成后：

1. `git diff` + `git status` 汇总所有变更，向用户报告摘要
2. **询问用户是否提交并推送**
3. 用户明确确认后才执行：
   ```shell
   git add package.json server/cli server/python CHANGELOG.md
   git commit -m "release: v<新版本号>"
   git tag v<新版本号>
   git push origin main v<新版本号>
   ```
4. 未确认一律不执行

提交规范：
- commit message 用 `release: v<版本号>`（与历史一致）
- **只 add 本次发布相关改动**，按文件名/目录显式添加，不要用 `git add -A` / `git add .`
- 推完才打 tag，不要单独 push tag——一次 push 同时推 main 和 tag，保证 Action 看到一致状态

## push 之后：确认发布工作流（Action 才是真正出包）

push tag 只是触发器，**push 成功 ≠ 发布完成**：

- `.github/workflows/release.yml` 7 个 job 全绿才算发布完成（含 `build-doo / test-go / publish-npm-platforms × 5 / publish-npm-cli / publish-npm-tools / publish-pypi / github-release`）。
- push 完成后**主动确认**：
  ```shell
  gh run list --workflow=release.yml -R dootask/tools -L 1
  gh run view <run-id> -R dootask/tools --json status,conclusion,url,jobs
  ```
- 工作流仍在跑时，挂后台轮询、结束即通知用户，**不要在前台死等**。
- 全绿后向用户报告：
  - npm: `https://www.npmjs.com/package/@dootask/tools`、`https://www.npmjs.com/package/@dootask/cli`
  - PyPI: `https://pypi.org/project/dootask-tools/`
  - GitHub Release: `https://github.com/dootask/tools/releases/tag/v<版本号>`
- 任一 job 失败时报告错误信息与 run url，交用户决定后续（已发布的子包**不能撤回同号**，失败要 bump 到 `vX.Y.(Z+1)` 重新走一遍）。

## 失败处理

任何步骤失败立即停止、报告错误信息，交用户决定；不要自动重试或跳过。

### 常见失败

- **NPM_TOKEN 缺失/无权限**：去 `https://www.npmjs.com/settings/<org>/tokens` 重新生成 Automation token，scope 必须含 `@dootask`；secret 名必须是 `NPM_TOKEN`（仓库 Settings → Secrets and variables → Actions）。
- **PyPI 403 同名版本**：已发过此版本号，不能重发；必须 bump 到下一个。
- **Go vet/test 失败**：tag 已推但代码有问题——快速修代码、bump 到下一个 tag 重发；旧 tag 留着无害，Go modules 会拿能编译的最新。
- **二进制子包发了主包没发**：主包 `optionalDependencies` 引用的子包版本已在 npm 上，但 `@dootask/cli` 本身没 publish。bump 版本号重发一轮即可。
