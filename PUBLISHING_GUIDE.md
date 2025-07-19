# 发布指南

本文档指导如何将 DooTask Tools 发布到 npm 包管理器。

## 发布到 npm

### 前置条件

1. 确保你有 npm 账号并已登录
2. 确保你有发布权限（对于 scoped package `@dootask/tools`）

### 发布步骤

1. **构建项目**
   ```bash
   yarn build
   ```

2. **登录 npm（如果未登录）**
   ```bash
   yarn login
   ```

3. **发布包**
   ```bash
   yarn publish
   ```

### 版本管理

发布前记得更新版本号：

```bash
# 补丁版本 (1.1.2 -> 1.1.3)
yarn version --patch

# 次要版本 (1.1.2 -> 1.2.0)
yarn version --minor

# 主要版本 (1.1.2 -> 2.0.0)
yarn version --major
```

### 发布配置

在 `package.json` 中已配置：

```json
{
  "publishConfig": {
    "access": "public"
  }
}
```

这确保包以公开方式发布。

### 使用方式

发布后，用户可以通过以下方式安装：

```bash
# 使用 yarn
yarn add @dootask/tools

# 使用 npm
npm install @dootask/tools
```

### 注意事项

1. 确保 `dist` 目录已正确构建
2. 检查 `package.json` 中的 `files` 字段包含所有必要文件
3. 发布前测试构建结果
4. 遵循语义化版本控制规范

