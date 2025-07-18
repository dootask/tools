#!/bin/bash

# DooTask Tools 一键发布脚本
# 作者: DooTask Team
# 日期: $(date +%Y-%m-%d)

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 获取当前版本号
get_current_version() {
    grep -o '"[0-9]*\.[0-9]*\.[0-9]*"' setup.py | head -1 | tr -d '"'
}

# 更新版本号
update_version() {
    local new_version=$1
    log_info "更新版本号到 $new_version"
    
    # 更新 setup.py
    sed -i.bak "s/version=\"[0-9]*\.[0-9]*\.[0-9]*\"/version=\"$new_version\"/" setup.py
    
    # 更新 __init__.py
    sed -i.bak "s/__version__ = \"[0-9]*\.[0-9]*\.[0-9]*\"/__version__ = \"$new_version\"/" dootask/__init__.py
    
    # 删除备份文件
    rm -f setup.py.bak dootask/__init__.py.bak
    
    log_success "版本号已更新"
}

# 验证版本号格式
validate_version() {
    local version=$1
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        log_error "版本号格式不正确，请使用 MAJOR.MINOR.PATCH 格式 (如: 1.0.0)"
        return 1
    fi
    return 0
}

# 主函数
main() {
    echo "=============================================="
    echo "🚀 DooTask Tools 一键发布脚本"
    echo "=============================================="
    echo
    
    # 检查是否在正确的目录
    if [[ ! -f "setup.py" ]]; then
        log_error "请在项目根目录运行此脚本"
        exit 1
    fi
    
    # 获取当前版本
    current_version=$(get_current_version)
    log_info "当前版本: $current_version"
    
    # 步骤1: 创建虚拟环境
    echo
    log_info "步骤1: 创建虚拟环境..."
    if [[ ! -d "venv" ]]; then
        if ! command_exists python3; then
            log_error "Python3 未安装，请先安装 Python3"
            exit 1
        fi
        python3 -m venv venv
        log_success "虚拟环境创建成功"
    else
        log_info "虚拟环境已存在，跳过创建"
    fi
    
    # 步骤2: 激活虚拟环境
    echo
    log_info "步骤2: 激活虚拟环境..."
    source venv/bin/activate
    log_success "虚拟环境已激活"
    
    # 步骤3: 检测安装发布工具
    echo
    log_info "步骤3: 检测安装发布工具..."
    
    # 升级pip
    python -m pip install --upgrade pip
    
    # 检查并安装build工具
    if ! python -c "import build" 2>/dev/null; then
        log_info "安装 build 工具..."
        pip install build
    fi
    
    # 检查并安装twine工具
    if ! python -c "import twine" 2>/dev/null; then
        log_info "安装 twine 工具..."
        pip install twine
    fi
    
    log_success "发布工具已准备就绪"
    
    # 步骤4: 修改版本号
    echo
    log_info "步骤4: 修改版本号..."
    
    while true; do
        read -p "请输入新的版本号 (当前: $current_version): " new_version
        
        if [[ -z "$new_version" ]]; then
            log_warning "版本号不能为空"
            continue
        fi
        
        if ! validate_version "$new_version"; then
            continue
        fi
        
        if [[ "$new_version" == "$current_version" ]]; then
            log_warning "新版本号与当前版本相同"
            read -p "确定要继续吗？(y/N): " confirm
            if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
                continue
            fi
        fi
        
        break
    done
    
    update_version "$new_version"
    
    # 步骤5: 清理旧文件
    echo
    log_info "步骤5: 清理旧文件..."
    rm -rf build/ dist/ *.egg-info/ dootask_tools.egg-info/
    log_success "旧文件已清理"
    
    # 步骤6: 构建包
    echo
    log_info "步骤6: 构建包..."
    python -m build
    log_success "包构建完成"
    
    # 显示构建结果
    echo
    log_info "构建文件:"
    ls -la dist/
    
    # 步骤7: 正式发布
    echo
    log_info "步骤7: 正式发布..."
    
    # 检查是否配置了PyPI认证
    if [[ ! -f ~/.pypirc ]]; then
        log_warning "未找到 ~/.pypirc 文件"
        log_info "请确保已配置 PyPI API Token"
        echo
        echo "配置方法:"
        echo "1. 访问 https://pypi.org/manage/account/token/"
        echo "2. 创建新的 API Token"
        echo "3. 创建 ~/.pypirc 文件:"
        echo "   [distutils]"
        echo "   index-servers = pypi"
        echo "   "
        echo "   [pypi]"
        echo "   repository = https://upload.pypi.org/legacy/"
        echo "   username = __token__"
        echo "   password = pypi-your-api-token-here"
        echo
        read -p "已配置完成，按回车继续..."
    fi
    
    # 询问是否先发布到测试环境
    echo
    read -p "是否先发布到测试环境？(Y/n): " test_publish
    if [[ "$test_publish" =~ ^[Yy]$|^$ ]]; then
        log_info "发布到测试环境..."
        twine upload --repository-url https://test.pypi.org/legacy/ dist/*
        log_success "已发布到测试环境"
        log_info "测试安装命令: pip install --index-url https://test.pypi.org/simple/ dootask-tools==$new_version"
        echo
        read -p "测试通过，继续发布到正式环境？(Y/n): " prod_publish
        if [[ ! "$prod_publish" =~ ^[Yy]$|^$ ]]; then
            log_info "发布已取消"
            exit 0
        fi
    fi
    
    # 发布到正式环境
    log_info "发布到正式环境..."
    twine upload dist/*
    log_success "发布成功！"
    
    # 验证发布
    echo
    log_info "验证发布..."
    log_info "安装命令: pip install --no-cache-dir dootask-tools==$new_version"
    log_info "PyPI 链接: https://pypi.org/project/dootask-tools/$new_version/"
    
    echo
    echo "=============================================="
    log_success "🎉 发布完成！版本 $new_version 已成功发布到 PyPI"
    echo "=============================================="
}

# 捕获中断信号
trap 'log_error "发布被中断"; exit 1' INT

# 运行主函数
main "$@" 