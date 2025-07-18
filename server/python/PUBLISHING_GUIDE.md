# 发布指南

本文档指导如何将 DooTask Tools 发布到 PyPI。

## 目录

- [快速开始](#快速开始)
- [一键发布脚本（推荐）](#一键发布脚本推荐)
- [手动发布（高级用户）](#手动发布高级用户)
- [自动化发布](#自动化发布)
- [发布脚本详细说明](#发布脚本详细说明)
- [维护建议](#维护建议)

## 快速开始

**推荐使用一键发布脚本**：

```bash
# 1. 给脚本添加执行权限
chmod +x publish.sh

# 2. 运行发布脚本
./publish.sh
```

脚本会自动引导您完成整个发布过程，包括版本更新、包构建和发布到 PyPI。

## 一键发布脚本（推荐）

项目提供了一个一键发布脚本 `publish.sh`，可以自动化完成整个发布流程。

### 使用方法

```bash
# 运行发布脚本
./publish.sh
```

### 脚本功能

发布脚本会自动执行以下步骤：

1. **创建虚拟环境** - 如果不存在，会自动创建 Python 虚拟环境
2. **激活虚拟环境** - 自动激活虚拟环境
3. **检测安装发布工具** - 自动安装 `build` 和 `twine` 工具
4. **修改版本号** - 交互式输入新版本号（遵循语义化版本控制）
5. **清理旧文件** - 清理之前的构建文件
6. **构建包** - 生成源码包和wheel包
7. **正式发布** - 发布到 PyPI（支持先发布到测试环境）

### 脚本特性

- 🎨 **彩色输出** - 使用颜色区分不同类型的信息
- 🔒 **安全验证** - 版本号格式验证，防止意外发布
- 🧪 **测试环境** - 支持先发布到测试环境验证
- 📦 **自动构建** - 自动生成源码包和wheel包
- ⚡ **错误处理** - 完善的错误处理和中断支持

### 使用示例

```bash
$ ./publish.sh
==============================================
🚀 DooTask Tools 一键发布脚本
==============================================

[INFO] 当前版本: 0.0.2

[INFO] 步骤1: 创建虚拟环境...
[SUCCESS] 虚拟环境已激活

[INFO] 步骤4: 修改版本号...
请输入新的版本号 (当前: 0.0.2): 1.0.0
[SUCCESS] 版本号已更新

[INFO] 步骤7: 正式发布...
是否先发布到测试环境？(Y/n): y
[SUCCESS] 已发布到测试环境
测试通过，继续发布到正式环境？(Y/n): y
[SUCCESS] 发布成功！

==============================================
🎉 发布完成！版本 1.0.0 已成功发布到 PyPI
==============================================
```

## 手动发布（高级用户）

如果您需要手动控制发布过程，可以按照以下步骤操作：

### 发布前准备

### 1. 安装发布工具

```bash
pip install build twine
```

### 2. 检查文件结构

确保以下文件存在：
- `setup.py` - 包配置
- `README.md` - 说明文档
- `LICENSE` - 许可证
- `requirements.txt` - 依赖列表
- `MANIFEST.in` - 打包配置
- `.gitignore` - 版本控制忽略文件

### 3. 更新版本号

在 `setup.py` 和 `dootask/__init__.py` 中更新版本号：

```python
version="1.0.0"  # 更新为新版本
```

## 构建包

### 1. 清理旧的构建文件

```bash
rm -rf build/ dist/ *.egg-info/
```

### 2. 构建包

```bash
python -m build
```

这将创建：
- `dist/dootask-tools-1.0.0.tar.gz` (源码包)
- `dist/dootask_tools-1.0.0-py3-none-any.whl` (wheel包)

### 3. 检查包内容

```bash
# 检查源码包
tar -tzf dist/dootask-tools-1.0.0.tar.gz

# 检查wheel包
unzip -l dist/dootask_tools-1.0.0-py3-none-any.whl
```

## 发布到 PyPI

### 1. 注册 PyPI 账户

- 访问 https://pypi.org/account/register/
- 注册账户并验证邮箱
- 启用双重认证（推荐）

### 2. 创建 API Token

- 访问 https://pypi.org/manage/account/token/
- 创建新的 API Token
- 复制并保存 Token

### 3. 配置认证

创建 `~/.pypirc` 文件：

```ini
[distutils]
index-servers = pypi

[pypi]
repository = https://upload.pypi.org/legacy/
username = __token__
password = pypi-your-api-token-here
```

### 4. 测试发布（可选）

首先发布到测试环境：

```bash
# 发布到测试PyPI
twine upload --repository-url https://test.pypi.org/legacy/ dist/*

# 从测试PyPI安装
pip install --index-url https://test.pypi.org/simple/ dootask-tools
```

### 5. 正式发布

```bash
twine upload dist/*
```

### 6. 验证发布

```bash
# 安装发布的包
pip install dootask-tools

# 测试导入
python -c "from dootask import DooTaskClient; print('发布成功！')"
```

## 版本管理

### 语义化版本控制

使用 `MAJOR.MINOR.PATCH` 格式：

- `MAJOR`: 不兼容的API变更
- `MINOR`: 向后兼容的功能新增
- `PATCH`: 向后兼容的问题修复

### 版本发布流程

1. 更新 `setup.py` 和 `__init__.py` 中的版本号
2. 更新 `README.md` 中的变更说明
3. 创建 Git 标签：`git tag v1.0.0`
4. 推送标签：`git push origin v1.0.0`
5. 构建并发布包

## 常见问题

### 1. 包名冲突

如果包名已存在，需要：
- 选择新的包名
- 更新 `setup.py` 中的 `name` 参数
- 更新所有相关文档

### 2. 依赖问题

确保 `requirements.txt` 中的依赖版本正确：
```
requests>=2.28.0
```

### 3. 文件缺失

如果构建时提示文件缺失，检查 `MANIFEST.in`：
```
include README.md
include LICENSE
include requirements.txt
```

### 4. 权限问题

确保 PyPI Token 有正确的权限：
- 对于新包，需要上传权限
- 对于已存在包，需要维护者权限

## 自动化发布

### 使用一键发布脚本

推荐使用项目提供的 `publish.sh` 脚本进行自动化发布：

```bash
# 简单快速的发布流程
./publish.sh
```

### 使用 GitHub Actions（可选）

如果您需要通过GitHub Actions进行发布，可以创建 `.github/workflows/publish.yml`：

```yaml
name: Publish to PyPI

on:
  release:
    types: [published]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Set up Python
      uses: actions/setup-python@v4
      with:
        python-version: '3.8'
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install build twine
    - name: Build package
      run: python -m build
    - name: Publish to PyPI
      env:
        TWINE_USERNAME: __token__
        TWINE_PASSWORD: ${{ secrets.PYPI_API_TOKEN }}
      run: twine upload dist/*
```

### 本地自动化发布

使用发布脚本的优势：
- 🚀 **简单易用** - 一条命令完成所有步骤
- 🔍 **实时反馈** - 彩色输出显示进度和状态
- 🧪 **测试验证** - 支持先发布到测试环境
- 🛡️ **安全保障** - 版本号验证和错误处理

## 发布脚本详细说明

### 环境要求

- Python 3.7+
- Bash 环境（macOS/Linux）
- Git（可选，用于创建标签）

### 配置要求

在使用发布脚本前，请确保：

1. **PyPI 账户配置**
   - 已注册 PyPI 账户
   - 创建了 API Token
   - 配置了 `~/.pypirc` 文件

2. **项目文件完整**
   - `setup.py` 配置正确
   - `dootask/__init__.py` 包含版本信息
   - `requirements.txt` 依赖列表完整
   - `README.md` 和 `LICENSE` 文件存在

### 脚本执行流程

1. **环境检查** - 检查 Python 和必要工具
2. **虚拟环境** - 创建或激活虚拟环境
3. **依赖安装** - 安装 build 和 twine 工具
4. **版本更新** - 交互式更新版本号
5. **文件清理** - 清理旧的构建文件
6. **包构建** - 生成分发包
7. **发布验证** - 检查 PyPI 配置
8. **包发布** - 上传到 PyPI

### 错误处理

脚本包含完善的错误处理：
- 版本号格式验证
- 文件存在性检查
- 网络连接验证
- 中断信号处理

### 故障排除

**常见问题**：

1. **权限错误**
   ```bash
   chmod +x publish.sh
   ```

2. **PyPI 认证失败**
   - 检查 `~/.pypirc` 文件配置
   - 确认 API Token 有效

3. **包名冲突**
   - 检查 PyPI 是否已存在同名包
   - 修改 `setup.py` 中的包名

## 维护建议

1. **定期更新依赖版本**
   - 检查 `requirements.txt` 中的版本约束
   - 测试新版本的兼容性

2. **监控 PyPI 下载统计**
   - 访问 PyPI 项目页面查看下载量
   - 分析用户使用情况

3. **及时处理用户反馈**
   - 关注 GitHub Issues
   - 响应用户问题和建议

4. **保持文档更新**
   - 更新 README.md 中的使用示例
   - 同步 API 文档变更

5. **遵循语义化版本控制**
   - MAJOR: 不兼容的 API 变更
   - MINOR: 向后兼容的功能新增
   - PATCH: 向后兼容的问题修复

6. **使用发布脚本**
   - 推荐使用 `./publish.sh` 进行发布
   - 避免手动操作可能出现的错误 