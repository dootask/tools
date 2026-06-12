#!/usr/bin/env node
/**
 * @dootask/cli 的 bin shim：转发 argv/stdio 到当前平台对应的 doo 二进制。
 *
 * 平台子包通过 optionalDependencies 引入；npm/yarn 只会装与当前 os/cpu 匹配的那一个，
 * 所以这里只需查 require.resolve 的解析结果即可拿到二进制路径。
 */
'use strict';

const { spawnSync } = require('child_process');

const PLATFORM_MAP = {
    'linux-x64': '@dootask/cli-linux-x64',
    'linux-arm64': '@dootask/cli-linux-arm64',
    'darwin-x64': '@dootask/cli-darwin-x64',
    'darwin-arm64': '@dootask/cli-darwin-arm64',
    'win32-x64': '@dootask/cli-win32-x64',
};

function main() {
    const key = `${process.platform}-${process.arch}`;
    const subPkg = PLATFORM_MAP[key];

    if (!subPkg) {
        console.error(
            `[doo] 当前平台 ${key} 暂不提供预编译二进制。\n` +
            `      支持: ${Object.keys(PLATFORM_MAP).join(', ')}\n` +
            `      可从源码构建：https://github.com/dootask/tools/tree/main/server/go/cmd/doo`
        );
        process.exit(1);
    }

    let binPath;
    try {
        binPath = require.resolve(
            `${subPkg}/bin/doo${process.platform === 'win32' ? '.exe' : ''}`
        );
    } catch (e) {
        console.error(
            `[doo] 缺少子包 ${subPkg}。请重新安装：\n` +
            `      npm i -g @dootask/cli\n` +
            `      或 npm i ${subPkg}`
        );
        process.exit(1);
    }

    const result = spawnSync(binPath, process.argv.slice(2), {
        stdio: 'inherit',
        windowsHide: false,
    });

    if (result.error) {
        console.error(`[doo] 执行失败：${result.error.message}`);
        process.exit(1);
    }

    // 透传子进程退出码（信号杀死按惯例转 128+signo）
    if (result.signal) {
        process.exit(128 + (require('os').constants.signals[result.signal] || 0));
    }
    process.exit(result.status === null ? 1 : result.status);
}

main();
