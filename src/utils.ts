import {Any, Func, MicroAppData, OpenAppPageParams, OpenWindowParams, PopoutWindowParams, SelectUsersParams, requestParams, responseSuccess, responseError, ModalParams} from './types';

// 存储微应用数据
let microAppData: MicroAppData | null = null;

// 微应用是否已准备好
let microAppReady = false;

// 备用z-index值，当无法从主应用获取nextZIndex时使用
let zIndexMissing = 1000;

// 存储主应用方法调用结果
const parentEvents: Record<string, (data: Any, error: Any) => void> = {};

// 调用主应用方法，如果主应用没有该方法，则向主应用发送消息
const methodTryParent = async (method: string, ...args: Any[]): Promise<Any | null> => {
    if (typeof window === 'undefined') {
        return null;
    }
    const methodFunc = await getAppData('methods.' + method);
    if (typeof methodFunc === 'function') {
        return methodFunc(...args);
    }
    return new Promise<Any | null>((resolve, reject) => {
        const id = Math.random().toString(36).substring(2, 15);
        parentEvents[id] = (data: Any, error: Any) => {
            delete parentEvents[id];
            if (error) {
                reject(error);
            } else {
                resolve(data);
            }
        };
        window.parent.postMessage({
            type: 'MICRO_APP_METHOD',
            message: {id, method, args}
        }, '*');
    });
};

if (typeof window !== 'undefined') {
    // 监听主应用注入的 microApp 对象
    window.addEventListener('message', (event) => {
        if (!event.data) {
            return;
        }
        const {type, message} = event.data;
        switch (type) {
            case 'MICRO_APP_INJECT':
                window.microApp = {
                    getData: () => {
                        return {
                            type: message.type,
                            props: message.props,
                        }
                    },
                }
                break;

            case 'MICRO_APP_METHOD_RESULT':
                const {id, result, error} = message;
                if (parentEvents[id]) {
                    parentEvents[id](result, error);
                }
                break;

            default:
                break;
        }
    });
}

/**
 * 检查当前应用是否为微前端应用
 * @returns {Promise<MicroAppData | null>} 返回微应用数据或null
 */
export const appReady = (): Promise<MicroAppData | null> => {
    return new Promise<MicroAppData | null>(async(resolve) => {
        if (typeof window === 'undefined') {
            resolve(null);
            return;
        }
        if (microAppReady) {
            resolve(microAppData);
            return;
        }
        let count = 0;
        while (typeof window.microApp === 'undefined' || typeof window.microApp.getData !== 'function') {
            await new Promise(resolve => setTimeout(resolve, 100));
            count++;
            if (count > 30) {
                resolve(null);
                return;
            }
        }
        microAppReady = true;
        microAppData = window.microApp.getData();
        resolve(microAppData);
    })
};

/**
 * 获取应用数据
 * @param {string | null} key - 可选参数，指定要获取的数据键名
 * @returns {Promise<Any>} 返回应用数据
 */
const getAppData = async (key: string | null = null): Promise<Any> => {
    if (await appReady() === null) {
        return null;
    }

    if (!key) return microAppData;

    return key.split('.').reduce((obj, k) => {
        if (obj && typeof obj === 'object') {
            const arrayIndex = /^\d+$/.test(k) ? parseInt(k) : k;
            return (obj as Record<string | number, Any>)[arrayIndex];
        }
        return null;
    }, microAppData);
};


// **************************************************************************************
// **************************************************************************************
// **************************************************************************************

/**
 * 检查当前应用是否为微前端应用
 * @returns {Promise<boolean>} 返回是否为微前端应用
 */
export const isMicroApp = async (): Promise<boolean> => {
    return await appReady() !== null;
};

/**
 * 检查是否为EEUI应用
 * @returns {Promise<boolean>} 返回是否为EEUI应用
 */
export const isEEUIApp = async (): Promise<boolean> => {
    return await getAppData('props.isEEUIApp');
};

/**
 * 检查是否为Electron应用
 * @returns {Promise<boolean>} 返回是否为Electron应用
 */
export const isElectron = async (): Promise<boolean> => {
    return await getAppData('props.isElectron');
};

/**
 * 检查是否为主Electron窗口
 * @returns {Promise<boolean>} 返回是否为主Electron窗口
 */
export const isMainElectron = async (): Promise<boolean> => {
    return await getAppData('props.isMainElectron');
};

/**
 * 检查是否为子Electron窗口
 * @returns {Promise<boolean>} 返回是否为子Electron窗口
 */
export const isSubElectron = async (): Promise<boolean> => {
    return await getAppData('props.isSubElectron');
};


// **************************************************************************************
// **************************************************************************************
// **************************************************************************************

/**
 * 获取当前主题名称
 * @returns {Promise<string>} 返回当前主题名称
 */
export const getThemeName = async (): Promise<string> => {
    return await getAppData('props.themeName');
};

/**
 * 获取当前用户ID
 * @returns {Promise<number>} 返回当前用户ID
 */
export const getUserId = async (): Promise<number> => {
    return await getAppData('props.userId');
};

/**
 * 获取当前用户Token
 * @returns {Promise<string>} 返回当前用户Token
 */
export const getUserToken = async (): Promise<string> => {
    return await getAppData('props.userToken');
};

/**
 * 获取当前用户信息
 * @returns {Promise<Any>} 返回当前用户信息对象
 */
export const getUserInfo = async (): Promise<Any> => {
    return await getAppData('props.userInfo');
};

/**
 * 获取基础URL
 * @returns {Promise<string>} 返回基础URL
 */
export const getBaseUrl = async (): Promise<string> => {
    return await getAppData('props.baseUrl');
};

/**
 * 获取系统信息
 * @returns {Promise<Any>} 返回系统信息对象
 */
export const getSystemInfo = async (): Promise<Any> => {
    return await getAppData('props.systemInfo');
};

/**
 * 获取页面类型
 * @returns {Promise<string>} 返回页面类型，可能的值为 'popout' 或 'embed'
 */
export const getWindowType = async (): Promise<string> => {
    return await getAppData('props.windowType');
}

/**
 * 获取语言列表
 * @returns {Promise<Any[]>} 返回语言列表
 */
export const getLanguageList = async (): Promise<Any[]> => {
    return await getAppData('props.languageList');
};

/**
 * 获取当前语言名称
 * @returns {Promise<string>} 返回当前语言名称
 */
export const getLanguageName = async (): Promise<string> => {
    return await getAppData('props.languageName');
};

// **************************************************************************************
// **************************************************************************************
// **************************************************************************************

/**
 * 关闭微前端应用
 * @param destroy - 可选参数，布尔值，表示是否销毁应用。默认为false。
 */
export const closeApp = async (destroy = false): Promise<void> => {
    await methodTryParent('close', destroy);
};

/**
 * 逐步返回上一个页面
 * @description 类似于浏览器的后退按钮，返回到最后一个页面时会关闭应用。
 */
export const backApp = async (): Promise<void> => {
    await methodTryParent('back');
};

/**
 * 应用窗口独立显示
 * @param params - 窗口参数
 */
export const popoutWindow = async (params?: PopoutWindowParams): Promise<void> => {
    await methodTryParent('popoutWindow', params);
};

/**
 * 打开新窗口
 * @param params - 窗口参数
 * @description 只在 isElectron 环境有效
 */
export const openWindow = async (params: OpenWindowParams): Promise<void> => {
    await methodTryParent('openWindow', params);
};

/**
 * 在新标签页打开URL
 * @param url - 要打开的URL
 * @description 只在 isElectron 环境有效
 */
export const openTabWindow = async (url: string): Promise<void> => {
    await methodTryParent('openTabWindow', url);
};

/**
 * 打开应用页面
 * @param params - 应用页面参数
 * @description 只在 isEEUIApp 环境有效
 */
export const openAppPage = async (params: OpenAppPageParams): Promise<void> => {
    await methodTryParent('openAppPage', params);
};

/**
 * 请求服务器API
 * @param params - API请求参数
 * @returns Promise 返回API请求结果
 */
export const requestAPI = async (params: requestParams): Promise<responseSuccess | responseError> => {
    return await methodTryParent('requestAPI', params);
};

/**
 * 选择用户
 * @param params - 可以是值或配置对象
 * @returns Promise 返回选择的用户结果
 */
export const selectUsers = async (params: SelectUsersParams): Promise<Any> => {
    return await methodTryParent('selectUsers', params);
};

/**
 * 调用$A上的额外方法
 * @param methodName - 方法名
 * @param args - 参数列表
 * @returns 方法返回值
 */
export const callExtraA = async (methodName: string, ...args: Any[]): Promise<Any> => {
    return await methodTryParent('extraCallA', methodName, ...args);
};

// **************************************************************************************
// **************************************************************************************
// **************************************************************************************

/**
 * 弹出成功提示框
 * @param message - 提示框内容
 * @returns Promise 返回提示框结果
 */
export const modalSuccess = async (message: string | ModalParams): Promise<Any> => {
    return await methodTryParent('extraCallA', 'modalSuccess', message);
};

/**
 * 弹出错误提示框
 * @param message - 提示框内容
 * @returns Promise 返回提示框结果
 */
export const modalError = async (message: string | ModalParams): Promise<Any> => {
    return await methodTryParent('extraCallA', 'modalError', message);
};

/**
 * 弹出警告提示框
 * @param message - 提示框内容
 * @returns Promise 返回提示框结果
 */
export const modalWarning = async (message: string | ModalParams): Promise<Any> => {
    return await methodTryParent('extraCallA', 'modalWarning', message);
};

/**
 * 弹出信息提示框
 * @param message - 提示框内容
 * @returns Promise 返回提示框结果
 */
export const modalInfo = async (message: string | ModalParams): Promise<Any> => {
    return await methodTryParent('extraCallA', 'modalInfo', message);
};

/**
 * 弹出系统提示框
 * @param message - 提示框内容
 * @returns Promise 返回提示框结果
 */
export const modalAlert = async (message: string): Promise<Any> => {
    return await methodTryParent('extraCallA', 'modalAlert', message);
};

// **************************************************************************************
// **************************************************************************************
// **************************************************************************************

/**
 * 弹出成功消息
 * @param message - 消息内容
 * @returns Promise 返回消息结果
 */
export const messageSuccess = async (message: string): Promise<Any> => {
    return await methodTryParent('extraCallA', 'messageSuccess', message);
};

/**
 * 弹出错误消息
 * @param message - 消息内容
 * @returns Promise 返回消息结果
 */
export const messageError = async (message: string): Promise<Any> => {
    return await methodTryParent('extraCallA', 'messageError', message);
};

/**
 * 弹出警告消息
 * @param message - 消息内容
 * @returns Promise 返回消息结果
 */
export const messageWarning = async (message: string): Promise<Any> => {
    return await methodTryParent('extraCallA', 'messageWarning', message);
};

/**
 * 弹出信息消息
 * @param message - 消息内容
 * @returns Promise 返回消息结果
 */
export const messageInfo = async (message: string): Promise<Any> => {
    return await methodTryParent('extraCallA', 'messageInfo', message);
};

// **************************************************************************************
// **************************************************************************************
// **************************************************************************************

/**
 * 获取下一个可用的模态框 z-index 值
 * @returns {number} 返回一个递增的 z-index 值
 */
export const nextZIndex = async (): Promise<number> => {
    const func = await getAppData('methods.nextZIndex');
    if (typeof func === 'function') {
        return func();
    }
    return zIndexMissing++;
};

/**
 * 设置应用关闭前的回调
 * @param callback - 回调函数，返回true则阻止关闭，false则允许关闭
 * @description 用于在应用关闭前执行操作，可以通过返回true来阻止关闭
 * @returns 返回一个函数，执行该函数可以注销监听器
 */
export const interceptBack = (callback: (data: Any) => boolean): () => void => {
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
    return () => { };
};

// **************************************************************************************
// **************************************************************************************
// **************************************************************************************

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
