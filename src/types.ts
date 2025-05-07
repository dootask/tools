/**
 * 微应用数据类型定义
 */

// 微应用实例类型
export interface MicroAppInstance {
    Vue: any;
    store: any;
    components: {
        DialogWrapper: any;
        UserSelect: any;
        DatePicker: any;
        [key: string]: any;
    };
}

// 微应用属性接口
export interface MicroAppProps {
    userId: string | number;
    userToken: string;
    userInfo: any;
    baseUrl: string;
    systemInfo: any;
    isEEUIApp: boolean;
    isElectron: boolean;
    isMainElectron: boolean;
    isSubElectron: boolean;
    languageList: any[];
    languageName: string;
    themeName: string;

    [key: string]: any;
}

// 微应用方法接口
export interface MicroAppMethods {
    close: (destroy?: boolean) => void;
    back: () => void;
    interceptBack: (callback: (data: any) => boolean) => (() => void);
    nextZIndex: () => number;
    openWindow: (params: OpenWindowParams) => void;
    openTabWindow: (url: string) => void;
    openAppPage: (params: OpenAppPageParams) => void;
    extraCallA: (...args: any[]) => any;

    [key: string]: any;
}

// 完整微应用数据接口
export interface MicroAppData {
    type: string;
    instance: MicroAppInstance;
    props: MicroAppProps;
    methods: MicroAppMethods;
}

// 打开窗口参数接口
export interface OpenWindowParams {
    name?: string;       // 窗口唯一标识
    url?: string;        // 访问地址
    force?: boolean;     // 是否强制创建新窗口，而不是重用已有窗口
    config?: {
        title?: string;     // 窗口标题
        titleFixed?: boolean; // 窗口标题是否固定
        width?: number;    // 窗口宽度
        height?: number;   // 窗口高度
        [key: string]: any;
    };

    [key: string]: any;
}

// 打开应用页面参数接口
export interface OpenAppPageParams {
    title?: string;      // 页面标题
    titleFixed?: boolean; // 窗口标题是否固定
    url?: string;        // 访问地址
    [key: string]: any;
}

// 扩展Window全局接口
declare global {
    interface Window {
        microApp?: {
            getData: () => MicroAppData;
            addDataListener: (callback: Function, autoTrigger?: boolean) => void;
            removeDataListener: (callback: Function) => void;
        };
        modalTransferIndex?: number;
        systemInfo?: any;
    }
}
