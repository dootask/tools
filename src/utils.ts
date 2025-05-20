import {Any, Func, MicroAppData, OpenAppPageParams, OpenWindowParams, PopoutWindowParams, SelectUsersParams, requestParams, responseSuccess, responseError} from './types';

// 存储微应用数据
let microAppData: MicroAppData | null = null;
let microAppReady = false;

// 备用z-index值，当无法从主应用获取nextZIndex时使用
let zIndexMissing = 1000;

// 在导入时环境允许自动初始化
if (!(typeof window === 'undefined' || typeof window.microApp === 'undefined' || typeof window.microApp.getData !== 'function')) {
    try {
        microAppData = window.microApp.getData();
    } catch (e) {
        console.warn('Failed to initialize DooTask tools:', e);
    }
}

/**
 * 检查当前应用是否为微前端应用
 * @returns {Promise<void>} 返回一个 Promise，当微前端应用准备好时解析
 */
export const appReady = (): Promise<MicroAppData> => {
    return new Promise<MicroAppData>((resolve) => {
        if (typeof window === 'undefined' || typeof window.microApp === 'undefined' || typeof window.microApp.getData !== 'function') {
            return;
        }
        if (!microAppReady) {
            microAppReady = true;
            microAppData = window.microApp.getData();
            resolve(microAppData)
        }
    })
};

/**
 * 检查当前应用是否为微前端应用
 * @returns {boolean} 如果当前应用是微前端应用，返回true；否则返回false
 */
export const isMicroApp = (): boolean => {
    return !(typeof window === 'undefined' || typeof window.microApp === 'undefined' || typeof window.microApp.getData !== 'function');
};

/**
 * 获取应用数据
 * @param {string | null} key - 可选参数，指定要获取的数据键名
 * @returns {Any} 当不传key时返回全部共享数据；传key时返回对应值
 */
export const getAppData = (key: string | null = null): Any => {
    if (!isMicroApp()) {
        return null;
    }

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
            return (obj as Record<string | number, Any>)[arrayIndex];
        }
        return null;
    }, microAppData);
};

// 直接导出属性变量
export const props = {
    /** 当前主题名称 */
    themeName: getAppData('props.themeName') || '',

    /** 当前用户ID */
    userId: Math.max(0, Number(getAppData('props.userId')) || 0),

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
    get: function (key: string, defaultValue: Any = null): Any {
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
 * @returns {number} 当前用户ID
 */
export const getUserId = (): number => {
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
 * @returns {Any} 当前用户信息对象
 */
export const getUserInfo = (): Any => {
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
 * @returns {Any} 系统信息对象
 */
export const getSystemInfo = (): Any => {
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
 * @returns {Any[]} 语言列表
 */
export const getLanguageList = (): Any[] => {
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

    /** 应用关闭前的回调
     * @param callback - 回调函数，返回true则阻止关闭，false则允许关闭
     * @description 用于在应用关闭前执行操作，可以通过返回true来阻止关闭
     * @returns 返回一个函数，执行该函数可以注销监听器
     */
    interceptBack: (callback: (data: Any) => boolean): (() => void) => {
        if (window.microApp?.addDataListener) {
            const interceptListener = (data: Any) => {
                if (data && data.type === 'beforeClose') {
                    return callback(data);
                }
                return false;
            };
            window.microApp.addDataListener(interceptListener, false);

            // 返回注销监听的函数
            return () => {
                if (window.microApp?.removeDataListener) {
                    window.microApp.removeDataListener(interceptListener);
                }
            };
        }
        // 如果没有添加监听，返回空函数
        return () => {
        };
    },

    /** 获取下一个模态框z-index */
    nextZIndex: (): number => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.nextZIndex === 'function') {
            return methodsData.nextZIndex();
        }
        return zIndexMissing++;
    },

    /** 选择用户 */
    selectUsers: async (params: SelectUsersParams): Promise<Any> => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.selectUsers === 'function') {
            return methodsData.selectUsers(params);
        }
        return null;
    },

    /** 应用窗口独立显示（只在 isElectron 环境有效） */
    popoutWindow: (params: PopoutWindowParams): void => {
        const methodsData = getAppData('methods');
        if (methodsData && typeof methodsData.popoutWindow === 'function') {
            methodsData.popoutWindow(params);
        }
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

    /** 请求服务器API */
    requestAPI: async (params: requestParams): Promise<responseSuccess | responseError> => {
        const dispatch = getAppData('instance.store.dispatch');
        if (dispatch && typeof dispatch === 'function') {
            return dispatch("call", params);
        } else {
            throw new Error('requestAPI method not found');
        }
    },

    /** 调用$A上的额外方法 */
    extraCallA: (methodName: string, ...args: Any[]): Any => {
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
 * 设置应用关闭前的回调 (兼容方法)
 * @param callback - 回调函数，返回true则阻止关闭，false则允许关闭
 * @description 用于在应用关闭前执行操作，可以通过返回true来阻止关闭
 * @returns 返回一个函数，执行该函数可以注销监听器
 */
export const interceptBack = (callback: (data: Any) => boolean): (() => void) => {
    return methods.interceptBack(callback);
};

/**
 * 获取下一个可用的模态框 z-index 值 (兼容方法)
 * @returns {number} 返回一个递增的 z-index 值
 */
export const nextZIndex = (): number => {
    return methods.nextZIndex();
};

/**
 * 应用窗口独立显示 (兼容方法)
 * @param params - 窗口参数
 * @description 只在 isElectron 环境有效
 */
export const popoutWindow = (params: PopoutWindowParams): void => {
    methods.popoutWindow(params);
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
 * 请求服务器API (兼容方法)
 * @param params - API请求参数
 * @returns Promise 返回API请求结果
 */
export const requestAPI = async (params: requestParams): Promise<responseSuccess | responseError> => {
    return methods.requestAPI(params);
};

/**
 * 调用$A上的额外方法 (兼容方法)
 * @param methodName - 方法名
 * @param args - 参数列表
 * @returns 方法返回值
 */
export const callExtraA = (methodName: string, ...args: Any[]): Any => {
    return methods.extraCallA(methodName, ...args);
};

/**
 * 选择用户 (兼容方法)
 * @param params - 可以是值或配置对象
 * @returns Promise 返回选择的用户结果
 */
export const selectUsers = async (params: SelectUsersParams): Promise<Any> => {
    return methods.selectUsers(params);
};

/**
 * 添加数据监听器
 * @param callback - 回调函数，当数据发生变化时调用
 * @param autoTrigger - 在初次绑定监听函数时如果有缓存数据，是否需要主动触发一次
 */
export const addDataListener = (callback: Func, autoTrigger = false): void => {
    if (window.microApp?.addDataListener) {
        window.microApp.addDataListener(callback, autoTrigger);
    }
};

/**
 * 移除数据监听器
 * @param callback - 回调函数，之前添加的监听器
 */
export const removeDataListener = (callback: Func): void => {
    if (window.microApp?.removeDataListener) {
        window.microApp.removeDataListener(callback);
    }
};
