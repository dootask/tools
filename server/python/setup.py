"""
DooTask Tools 包设置文件
"""

from setuptools import setup, find_packages
import os

# 读取 README 文件
def read_readme():
    readme_path = os.path.join(os.path.dirname(__file__), "README.md")
    with open(readme_path, "r", encoding="utf-8") as f:
        return f.read()

# 读取 requirements.txt
def read_requirements():
    requirements_path = os.path.join(os.path.dirname(__file__), "requirements.txt")
    with open(requirements_path, "r", encoding="utf-8") as f:
        return [line.strip() for line in f if line.strip() and not line.startswith("#")]

setup(
    name="dootask-tools",
    version="0.0.5",
    author="DooTask Team",
    author_email="support@dootask.com",
    description="DooTask Tools 客户端库",
    long_description=read_readme(),
    long_description_content_type="text/markdown",
    url="https://github.com/dootask/dootask-tools",
    packages=["dootask"],
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Communications :: Chat",
        "Topic :: Office/Business :: Groupware",
        "Topic :: Internet :: WWW/HTTP :: Dynamic Content",
    ],
    python_requires=">=3.7",
    install_requires=read_requirements(),
    include_package_data=True,
    keywords="dootask api client chat project management collaboration",
    project_urls={
        "Bug Reports": "https://github.com/dootask/dootask-tools/issues",
        "Source": "https://github.com/dootask/dootask-tools",
        "Documentation": "https://github.com/dootask/dootask-tools#readme",
    },
) 