# 发布指南

本文档指导如何使用 DooTask Tools (Go) 模块。

## Go 模块使用方式

Go 模块不需要像其他语言那样发布到包管理器。Go 使用 `go get` 命令直接从 Git 仓库获取代码。

### 基本使用

```bash
# 获取最新版本
go get github.com/dootask/tools/server/go

# 获取特定版本（如果打了 tag）
go get github.com/dootask/tools/server/go@v1.0.0

# 获取特定提交
go get github.com/dootask/tools/server/go@commit-hash
```

### 在项目中使用

在你的 `go.mod` 文件中添加依赖：

```go
module your-project

go 1.21

require (
    github.com/dootask/tools/server/go v0.0.0-20231201000000-abcdef123456
)
```

### 版本管理

- **最新版本**: 直接使用 `go get` 获取最新代码
- **特定提交**: 使用 `go get module@commit-hash` 获取特定提交
- **标签版本**: 如果仓库打了 Git 标签，可以使用 `go get module@v1.0.0`

### 注意事项

1. 确保你的 Go 模块有正确的 `go.mod` 文件
2. 模块路径应该与 Git 仓库路径一致
3. 如果需要私有仓库，确保有适当的访问权限
4. 建议在 `go.mod` 中使用 `replace` 指令进行本地开发测试

### 本地开发

```bash
# 使用 replace 指令指向本地路径
go mod edit -replace github.com/dootask/tools/server/go=../dootask-tools/server/go
```

