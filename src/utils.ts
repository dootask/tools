import {MicroAppData, OpenAppPageParams, OpenWindowParams} from './types';

// 存储微应用数据（自动初始化）
let microAppData: MicroAppData | null = null;

// 在导入时自动初始化
if (typeof window !== 'undefined' && window.microApp && window.microApp.getData) {
    try {
        microAppData = window.microApp.getData();
    } catch (e) {
        console.warn('Failed to initialize dootask-tools:', e);
    }
}

/**
 * 检查当前应用是否为微前端应用
 * @returns {boolean} 如果当前应用是微前端应用，返回true；否则返回false
 */
export const isMicroApp = (): boolean => {
    return window.microApp !== undefined && typeof window.microApp.getData === 'function';
};

/**
 * 获取应用数据
 * @param {string | null} key - 可选参数，指定要获取的数据键名
 * @returns {any} 当不传key时返回全部共享数据；传key时返回对应值
 */
export const getAppData = (key: string | null = null): any => {
    if (!microAppData && window.microApp?.getData) {
        microAppData = window.microApp.getData();
    }

    if (!microAppData) return null;
    if (!key) return microAppData;

    return key.split('.').reduce((obj, k) => {
        if (obj && typeof obj === 'object') {
            // 处理数组索引（如 items.0）
            const arrayIndex = /^\d+$/.test(k) ? parseInt(k) : k;
            // 使用类型断言解决动态索引的类型问题
            return (obj as Record<string | number, any>)[arrayIndex];
        }
        return null;
    }, microAppData);
};

// 直接导出属性变量
export const props = {
    /** 当前主题名称 */
    themeName: getAppData('props.themeName') || '',

    /** 当前用户ID */
    userId: getAppData('props.userId') || '',

    /** 当前用户Token */
    userToken: getAppData('props.userToken') || '',

    /** 当前用户信息 */
    userInfo: getAppData('props.userInfo') || null,

    /** 基础URL */
    baseUrl: getAppData('props.baseUrl') || '',

    /** 系统信息 */
    systemInfo: getAppData('props.systemInfo') || null,

    /** 是否为EEUI应用 */
    isEEUIApp: !!getAppData('props.isEEUIApp'),

    /** 是否为Electron应用 */
    isElectron: !!getAppData('props.isElectron'),

    /** 是否为主Electron窗口 */
    isMainElectron: !!getAppData('props.isMainElectron'),

    /** 是否为子Electron窗口 */
    isSubElectron: !!getAppData('props.isSubElectron'),

    /** 语言列表 */
    languageList: getAppData('props.languageList') || [],

    /** 当前语言名称 */
    languageName: getAppData('props.languageName') || '',

    /** 获取原始属性字段 */
    get: function (key: string, defaultValue: any = null): any {
        return getAppData(`props.${key}`) || defaultValue;
    }
};

// 兼容保留原来的方法

/**
 * 获取当前主题名称 (兼容方法)
 * @returns {string} 当前主题名称
 */
export const getThemeName = (): string => {
    return props.themeName;
};

/**
 * 获取当前用户ID (兼容方法)
 * @returns {string | number} 当前用户ID
 */
export const getUserId = (): string | number => {
    return props.userId;
};

/**
 * 获取当前用户Token (兼容方法)
 * @returns {string} 当前用户Token
 */
export const getUserToken = (): string => {
    return props.userToken;
};

/**
 * 获取当前用户信息 (兼容方法)
 * @returns {any} 当前用户信息对象
 */
export const getUserInfo = (): any => {
    return props.userInfo;
};

/**
 * 获取基础URL (兼容方法)
 * @returns {string} 基础URL
 */
export const getBaseUrl = (): string => {
    return props.baseUrl;
};

/**
 * 获取系统信息 (兼容方法)
 * @returns {any} 系统信息对象
 */
export const getSystemInfo = (): any => {
    return props.systemInfo;
};

/**
 * 检查是否为EEUI应用 (兼容方法)
 * @returns {boolean} 是否为EEUI应用
 */
export const isEEUIApp = (): boolean => {
    return props.isEEUIApp;
};

/**
 * 检查是否为Electron应用 (兼容方法)
 * @returns {boolean} 是否为Electron应用
 */
export const isElectron = (): boolean => {
    return props.isElectron;
};

/**
 * 检查是否为主Electron窗口 (兼容方法)
 * @returns {boolean} 是否为主Electron窗口
 */
export const isMainElectron = (): boolean => {
    return props.isMainElectron;
};

/**
 * 检查是否为子Electron窗口 (兼容方法)
 * @returns {boolean} 是否为子Electron窗口
 */
export const isSubElectron = (): boolean => {
    return props.isSubElectron;
};

/**
 * 获取语言列表 (兼容方法)
 * @returns {any[]} 语言列表
 */
export const getLanguageList = (): any[] => {
    return props.languageList;
};

/**
 * 获取当前语言名称 (兼容方法)
 * @returns {string} 当前语言名称
 */
export const getLanguageName = (): string => {
    return props.languageName;
};

// 直接导出方法变量
export const methods = {
    /** 关闭当前应用 */
    close: (destroy = false): void => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.close === 'function') {
            methodsData.close(destroy);
        }
    },

    /** 返回上一页 */
    back: (): void => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.back === 'function') {
            methodsData.back();
        }
    },

    /** 获取下一个模态框z-index */
    nextZIndex: (): number => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.nextZIndex === 'function') {
            return methodsData.nextZIndex();
        }
        return 1000;
    },

    /** 打开新窗口（只在 isElectron 环境有效） */
    openWindow: (params: OpenWindowParams): void => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.openWindow === 'function') {
            methodsData.openWindow(params);
        }
    },

    /** 在新标签页打开URL（只在 isElectron 环境有效） */
    openTabWindow: (url: string): void => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.openTabWindow === 'function') {
            methodsData.openTabWindow(url);
        }
    },

    /** 打开应用页面（只在 isEEUIApp 环境有效） */
    openAppPage: (params: OpenAppPageParams): void => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.openAppPage === 'function') {
            methodsData.openAppPage(params);
        }
    },

    /** 调用$A上的额外方法 */
    extraCallA: (methodName: string, ...args: any[]): any => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.extraCallA === 'function') {
            return methodsData.extraCallA(methodName, ...args);
        }
        return null;
    }
};

// 兼容保留原来的方法

/**
 * 关闭微前端应用 (兼容方法)
 * @param destroy - 可选参数，布尔值，表示是否销毁应用。默认为false。
 */
export const closeApp = (destroy = false): void => {
    methods.close(destroy);
};

/**
 * 逐步返回上一个页面 (兼容方法)
 * @description 类似于浏览器的后退按钮，返回到最后一个页面时会关闭应用。
 */
export const backApp = (): void => {
    methods.back();
};

/**
 * 获取下一个可用的模态框 z-index 值 (兼容方法)
 * @returns {number} 返回一个递增的 z-index 值
 */
export const nextZIndex = (): number => {
    return methods.nextZIndex();
};

/**
 * 打开新窗口 (兼容方法)
 * @param params - 窗口参数
 * @description 只在 isElectron 环境有效
 */
export const openWindow = (params: OpenWindowParams): void => {
    methods.openWindow(params);
};

/**
 * 在新标签页打开URL (兼容方法)
 * @param url - 要打开的URL
 * @description 只在 isElectron 环境有效
 */
export const openTabWindow = (url: string): void => {
    methods.openTabWindow(url);
};

/**
 * 打开应用页面 (兼容方法)
 * @param params - 应用页面参数
 * @description 只在 isEEUIApp 环境有效
 */
export const openAppPage = (params: OpenAppPageParams): void => {
    methods.openAppPage(params);
};

/**
 * 调用$A上的额外方法 (兼容方法)
 * @param methodName - 方法名
 * @param args - 参数列表
 * @returns 方法返回值
 */
export const callExtraA = (methodName: string, ...args: any[]): any => {
    return methods.extraCallA(methodName, ...args);
};

/**
 * 添加数据监听器
 * @param callback - 回调函数，当数据发生变化时调用
 * @param autoTrigger - 在初次绑定监听函数时如果有缓存数据，是否需要主动触发一次
 */
export const addDataListener = (callback: Function, autoTrigger = false): void => {
    if (window.microApp?.addDataListener) {
        window.microApp.addDataListener(callback, autoTrigger);
    }
};

/**
 * 移除数据监听器
 * @param callback - 回调函数，之前添加的监听器
 */
export const removeDataListener = (callback: Function): void => {
    if (window.microApp?.removeDataListener) {
        window.microApp.removeDataListener(callback);
    }
};
