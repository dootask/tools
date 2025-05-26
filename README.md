# DooTask 应用工具库

这是一个为DooTask应用开发设计的工具库，提供了一系列实用的API，用于与主应用进行交互。它封装了微前端通信、用户数据获取、页面交互等多种功能，让您的DooTask应用开发变得简单高效。

## 特点

- **零配置** - 导入后自动初始化，无需手动配置（服务器渲染不支持）
- **类型支持** - 完善的TypeScript类型定义，提供智能提示
- **双重访问方式** - 支持对象式和函数式两种调用风格
- **强大功能** - 提供完整的用户、系统信息和交互方法

## 安装

```bash
npm install @dootask/tools --save
```

## 使用方法

### 对象式调用 (推荐)

直接引入并使用，无需手动初始化：

```typescript
import {props, methods, isMicroApp} from '@dootask/tools';

// 使用直接提供的属性
const theme = props.themeName;
const userId = props.userId;
const userInfo = props.userInfo;

// 使用直接提供的方法
methods.openWindow({url: 'https://example.com'});
methods.close(); // 关闭当前应用

// 检查是否在微前端环境中
if (isMicroApp()) {
    console.log('当前在微前端环境中运行');
}
```

### 函数式调用 (兼容模式)

可以继续使用函数式调用风格：

```typescript
import {
    getThemeName,
    getUserInfo,
    openWindow,
    closeApp,
    isMicroApp,
    // 更多API...
} from '@dootask/tools';

// 获取当前主题
const theme = getThemeName();

// 获取用户信息
const user = getUserInfo();

// 打开新窗口
openWindow({url: 'https://example.com'});

// 关闭当前应用
closeApp();
```

## API 文档

### props 属性

| 属性名                      | 类型                  | 说明                   |
|--------------------------|---------------------|----------------------|
| `themeName`              | `string`            | 当前主题名称               |
| `userId`                 | `number`            | 当前用户ID，0 表示未登录       |
| `userToken`              | `string`            | 当前用户Token            |
| `userInfo`               | `object`            | 当前用户信息对象             |
| `baseUrl`                | `string`            | 基础URL                |
| `systemInfo`             | `object`            | 系统信息对象               |
| `windowType`             | 'popout' \| 'embed' | 页面打开类型               |
| `isEEUIApp`              | `boolean`           | 是否为EEUI应用（App客户端）    |
| `isElectron`             | `boolean`           | 是否为Electron应用（电脑客户端） |
| `isMainElectron`         | `boolean`           | 是否为主Electron窗口       |
| `isSubElectron`          | `boolean`           | 是否为子Electron窗口       |
| `languageList`           | `array`             | 语言列表                 |
| `languageName`           | `string`            | 当前语言名称               |
| `get(key, defaultValue)` | `function`          | 获取原始属性字段             |

### methods 方法

| 方法名             | 参数                                   | 返回值                                         | 说明                                       |
|-----------------|--------------------------------------|---------------------------------------------|------------------------------------------|
| `close`         | `destroy?: boolean`                  | `void`                                      | 关闭当前应用                                   |
| `back`          | -                                    | `void`                                      | 返回上一页                                    |
| `interceptBack` | `callback: (data: any) => boolean`   | `() => void`                                | 设置应用关闭前的回调，返回true可阻止关闭。返回一个可注销监听的函数      |
| `nextZIndex`    | -                                    | `number`                                    | 获取下一个可用的模态框z-index                       |
| `selectUsers`   | `params: SelectUsersParams`          | `Promise<any>`                              | 选择用户，可以传入多种配置来自定义选择器                     |
| `popoutWindow`  | `objects`                            | `void`                                      | 应用窗口独立显示                                 |
| `openWindow`    | `objects`                            | `void`                                      | 打开新窗口（只在 isElectron 环境有效）                |
| `openTabWindow` | `url: string`                        | `void`                                      | 在新标签页打开URL，直接传入URL地址（只在 isElectron 环境有效） |
| `openAppPage`   | `objects`                            | `void`                                      | 打开应用页面（只在 isEEUIApp 环境有效）                |
| `requestAPI`    | `params: requestParams`              | `Promise<responseSuccess \| responseError>` | 请求服务器API                                 |
| `extraCallA`    | `methodName: string, ...args: any[]` | `any`                                       | 调用$A上的额外方法                               |

### 全局函数

| 函数名                  | 参数                                          | 返回值            | 说明               |
|----------------------|---------------------------------------------|----------------|------------------|
| `appReady`           | -                                           | `Promise<any>` | 应用准备就绪的Promise对象 |
| `isMicroApp`         | -                                           | `boolean`      | 检查当前是否在微前端环境中运行  |
| `getAppData`         | `key?: string`                              | `any`          | 获取原始微前端应用数据      |
| `addDataListener`    | `callback: Function, autoTrigger?: boolean` | `void`         | 添加数据监听器          |
| `removeDataListener` | `callback: Function`                        | `void`         | 移除数据监听器          |

### 兼容函数

所有`props`和`methods`中的属性和方法都有对应的函数式调用版本，如：

- `getThemeName()` 对应 `props.themeName`
- `getUserId()` 对应 `props.userId`
- `closeApp()` 对应 `methods.close()`
- 等等...

## 使用示例

### 检测运行环境

```typescript
import {props, appReady, isMicroApp} from '@dootask/tools';

appReady().then(() => {
    console.log('应用已准备就绪');
})

if (isMicroApp()) {
    // 在微前端环境中运行
    if (props.isElectron) {
        console.log('在Electron环境中运行');
        if (props.isMainElectron) {
            console.log('这是主窗口');
        } else if (props.isSubElectron) {
            console.log('这是子窗口');
        }
    }
} else {
    console.log('不在微前端环境中运行');
}
```

### 应用关闭拦截

```typescript
import {methods} from '@dootask/tools';

let hasUnsavedChanges = true;

// 设置应用关闭前的回调
method.interceptBack((data) => {
    if (hasUnsavedChanges) {
        // 如果有未保存的数据，则阻止关闭
        if (confirm('有未保存的数据，确定要关闭吗？')) {
            // 用户确认关闭，可以执行保存操作
            saveData();
            return false; // 允许关闭
        } else {
            return true; // 阻止关闭
        }
    }
    return false; // 没有未保存的数据，允许关闭
});
```

### 选择用户

```typescript
import {methods, selectUsers} from '@dootask/tools';

// 方法一：使用 methods 对象
 methods.selectUsers({
    value: [], // 已选择的值
    projectId: 123, // 指定项目ID
    title: '选择成员', // 弹窗标题
    showSelectAll: true // 显示全选项
}).then(result => {
    console.log('选择的用户：', result);
});

// 方法二：使用兼容函数
selectUsers({
    value: [],
    onlyGroup: true, // 仅显示群组
    showBot: false // 不显示机器人
}).then(result => {
    console.log('选择的群组：', result);
});
```

### 监听数据变化

```typescript
import {addDataListener, removeDataListener} from '@dootask/tools';

// 添加数据监听器
const dataListener = (data) => {
    console.log('收到新数据:', data);
};

// 添加监听，并在初次绑定时触发
addDataListener(dataListener, true);

// 移除监听
// removeDataListener(dataListener);
```

### 弹出窗口和页面

```typescript
import {methods, props} from '@dootask/tools';

// 将当前页面作为独立窗口显示
methods.popoutWindow();

// 将当前页面作为独立窗口显示（自定义窗口信息，信息仅对Electron环境有效）
methods.popoutWindow({
    title: '独立窗口',      // 窗口标题
    width: 1000,           // 窗口宽度
    height: 700,           // 窗口高度
    minWidth: 800          // 窗口最小宽度
});

// 在Electron环境中将当前页面以独立窗口形式显示
if (props.isElectron) {
    // 打开新窗口
    methods.openWindow({
        name: 'my-window-id',  // 窗口唯一标识
        url: 'https://example.com',  // 访问地址
        force: false,  // 是否强制创建新窗口，而不是重用已有窗口
        config: {
            title: '标题',  // 窗口标题
            titleFixed: true,       // 窗口标题是否固定
            width: Math.min(window.screen.availWidth, 1200),  // 窗口宽度
            height: Math.min(window.screen.availHeight, 800),  // 窗口高度
        }
    });

    // 在新标签页打开URL
    methods.openTabWindow('https://example.com');
}

// 在EEUI环境中打开应用页面
if (props.isEEUIApp) {
    methods.openAppPage({
        title: '标题',       // 页面标题
        titleFixed: true,    // 窗口标题是否固定
        url: 'https://example.com',  // 访问地址
    });
}
```

### 请求服务器API

```typescript
import {methods, props} from '@dootask/tools';

// 请求服务器API
methods.requestAPI({
    url: 'users/info',  // 访问接口路径，接口文档请查看 https://你的域名/docs/index.html
}).then((res) => {
    console.log(res);
});
```

## 注意事项

1. 这个库会自动检测是否在微前端环境中运行。如果不在微前端环境中，大部分属性将返回空值，方法将无操作。

2. `getAppData`方法可以获取微前端应用的原始数据，那些未被`props`和`methods`封装的数据也可以通过这个方法获取。

3. 如果你希望调用$A上的方法，可以使用`methods.extraCallA`或`callExtraA`方法。

## 贡献和反馈

如果你在使用中发现任何问题，或者有改进建议，欢迎在GitHub仓库提交Issue或Pull Request。

## 许可证

MIT
