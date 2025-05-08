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
    userId: number;
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
    selectUsers: (params: SelectUsersParams) => Promise<any>;
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

// 窗口配置接口
export interface WindowConfig {
    title?: string;     // 窗口标题
    titleFixed?: boolean; // 窗口标题是否固定
    width?: number;    // 窗口宽度
    height?: number;   // 窗口高度
    minWidth?: number; // 窗口最小宽度
    [key: string]: any;
    // 更多配置项参考 https://www.electronjs.org/docs/latest/api/structures/base-window-options
}

// 打开独立窗口参数接口
export interface PopoutWindowParams extends WindowConfig {
    url?: string;   // 自定义访问地址，如果为空则打开当前页面
    [key: string]: any;
}

// 打开窗口参数接口
export interface OpenWindowParams {
    name?: string;       // 窗口唯一标识
    url?: string;        // 访问地址
    force?: boolean;     // 是否强制创建新窗口，而不是重用已有窗口
    config?: WindowConfig;

    [key: string]: any;
}

// 打开应用页面参数接口
export interface OpenAppPageParams {
    title?: string;      // 页面标题
    titleFixed?: boolean; // 窗口标题是否固定
    url?: string;        // 访问地址
    [key: string]: any;
}

// 选择用户参数接口
export interface SelectUsersParams {
    value?: string | number | Array<any>;     // 已选择的值，默认值: []
    uncancelable?: Array<any>;                // 不允许取消的列表，默认值: []
    disabledChoice?: Array<any>;              // 禁止选择的列表，默认值: []
    projectId?: number;                       // 指定项目ID，默认值: 0
    noProjectId?: number;                     // 指定非项目ID，默认值: 0
    dialogId?: number;                        // 指定会话ID，默认值: 0
    showBot?: boolean;                        // 是否显示机器人，默认值: false
    showDisable?: boolean;                    // 是否显示禁用的，默认值: false
    multipleMax?: number;                     // 最大选择数量
    title?: string;                           // 弹窗标题
    placeholder?: string;                     // 搜索提示
    showSelectAll?: boolean;                  // 显示全选项，默认值: true
    showDialog?: boolean;                     // 是否显示会话，默认值: false
    onlyGroup?: boolean;                      // 仅显示群组，默认值: false
    beforeSubmit?: Function;                  // 提交前的回调
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
