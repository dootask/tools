# DooTask Tools - Vite 示例

这是一个使用 Vite 构建的示例项目，展示如何在 Vue 3 项目中使用 `@dootask/tools` 库。

## 功能特性

- 🚀 基于 Vite 的快速开发体验
- 🎨 现代化的 UI 设计
- 📱 响应式布局
- 🔧 完整的 dootask-tools 功能演示
- 📝 详细的代码注释和使用说明

## 快速开始

### 安装依赖

```bash
cd example
npm install
```

### 开发模式

```bash
npm run dev
```

项目将在 `http://localhost:5173` 启动。

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 功能演示

这个示例展示了 `@dootask/tools` 的主要功能：

### 1. 应用状态检测

- 检测是否为微前端应用
- 获取用户信息、主题、语言等配置
- 显示窗口类型和运行环境

### 2. 窗口管理

- **打开独立窗口**: 使用 `popoutWindow()` 创建独立窗口
- **打开新窗口**: 使用 `openWindow()` 创建新窗口（仅在 Electron 环境有效）

### 3. 用户交互

- **选择用户**: 使用 `selectUsers()` 打开用户选择器
- **API 请求**: 使用 `requestAPI()` 发送服务器请求

### 4. 应用控制

- **关闭应用**: 使用 `closeApp()` 关闭当前应用
- **返回操作**: 使用 `backApp()` 执行返回操作

### 5. 提示框

- **成功提示**: 使用 `modalSuccess()` 显示成功消息
- **错误提示**: 使用 `modalError()` 显示错误消息
- **警告提示**: 使用 `modalWarning()` 显示警告消息
- **信息提示**: 使用 `modalInfo()` 显示信息消息
- **系统提示**: 使用 `modalAlert()` 显示系统提示框

## 技术栈

- **构建工具**: Vite 5.0
- **框架**: Vue 3.3
- **语言**: TypeScript
- **样式**: CSS3 + Grid/Flexbox
- **工具库**: @dootask/tools

## 注意事项

1. 这个示例需要在 DooTask 微前端环境中运行才能获得完整功能
2. 在独立环境中运行时，某些功能可能不可用或会显示默认值
3. 确保 `@dootask/tools` 库已正确安装和配置
4. 所有 API 都是异步的，返回 Promise 对象
5. 某些功能只在特定环境中有效（如 `openWindow` 只在 Electron 环境中有效）

## 开发建议

1. 使用 TypeScript 获得更好的类型支持
2. 遵循 Vue 3 Composition API 的最佳实践
3. 利用 Vite 的热重载功能提高开发效率
4. 在生产环境中启用代码分割和优化
5. 在使用任何方法之前，建议先调用 `appReady()` 确保应用已准备就绪

## 相关链接

- [主项目 README](../README.md) - 查看完整的 API 文档和使用说明
- [@dootask/tools 源码](../src/) - 查看工具库源码

## 许可证

MIT License
