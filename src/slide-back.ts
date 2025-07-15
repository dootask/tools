/**
 * 滑动返回功能模块
 * 监听触摸事件，实现从屏幕左边向右滑动触发返回操作
 * 
 * @example
 * ```typescript
 * import { SlideBack, createSlideBack } from 'dootask-tools';
 * 
 * // 方式一：使用便捷函数（推荐）
 * const slideBack = createSlideBack({
 *     onBack: () => {
 *         console.log('执行返回操作');
 *         // 执行返回逻辑
 *     }
 * });
 * 
 * // 方式二：手动创建实例
 * const slideBack = new SlideBack({
 *     edgeWidth: 30,           // 左边缘检测区域宽度（默认30px）
 *     minSwipeDistance: 100,   // 水平滑动最小距离（默认100px）
 *     maxVerticalOffset: 100,  // 最大垂直偏移距离（默认100px）
 *     maxSwipeTime: 2000,      // 最大滑动时间（默认2000ms）
 *     onBack: async () => {
 *         // 自定义返回逻辑
 *         await someAsyncOperation();
 *     }
 * });
 * 
 * // 启用滑动返回
 * slideBack.enable();
 * 
 * // 禁用滑动返回
 * slideBack.disable();
 * 
 * // 检查是否已启用
 * if (slideBack.enabled) {
 *     console.log('滑动返回已启用');
 * }
 * 
 * // 运行时更新配置
 * slideBack.updateOptions({
 *     edgeWidth: 50,
 *     onBack: newBackHandler
 * });
 * 
 * // 销毁实例（会自动禁用）
 * slideBack.destroy();
 * ```
 * 
 * @example
 * ```typescript
 * // 与微前端应用集成
 * import { backApp, createSlideBack } from 'dootask-tools';
 * 
 * const slideBack = createSlideBack({
 *     onBack: async () => {
 *         try {
 *             await backApp(); // 调用微前端返回方法
 *         } catch (error) {
 *             console.warn('返回失败:', error);
 *         }
 *     }
 * });
 * 
 * // 在组件卸载时清理
 * onUnmounted(() => {
 *     slideBack.destroy();
 * });
 * ```
 * 
 * @example
 * ```typescript
 * // 自定义配置示例
 * const slideBack = new SlideBack({
 *     edgeWidth: 40,          // 更大的边缘检测区域
 *     minSwipeDistance: 80,   // 更短的滑动距离要求
 *     maxVerticalOffset: 120, // 允许更大的垂直偏移
 *     maxSwipeTime: 1500,     // 更短的时间限制
 *     onBack: () => {
 *         // 可以添加振动反馈
 *         if (navigator.vibrate) {
 *             navigator.vibrate(50);
 *         }
 *         
 *         // 执行返回逻辑
 *         history.back();
 *     }
 * });
 * 
 * slideBack.enable();
 * ```
 * 
 * @description
 * 滑动返回的触发条件：
 * 1. 触摸起始点在屏幕左边缘指定区域内
 * 2. 水平向右滑动距离达到最小要求
 * 3. 垂直偏移不超过最大限制
 * 4. 滑动时间在允许范围内
 * 
 * @author DooTask Team
 * @version 1.0.0
 */

export interface SlideBackOptions {
    /** 左边缘检测区域宽度，默认30px */
    edgeWidth?: number;
    /** 水平滑动最小距离，默认100px */
    minSwipeDistance?: number;
    /** 最大垂直偏移距离，默认100px */
    maxVerticalOffset?: number;
    /** 最大滑动时间，默认2000ms */
    maxSwipeTime?: number;
    /** 返回回调函数 */
    onBack?: () => Promise<void> | void;
}

export class SlideBack {
    private touchStartX = 0;
    private touchStartY = 0;
    private touchStartTime = 0;
    private isSwiping = false;
    private isEnabled = false;
    
    private options: Required<SlideBackOptions>;
    
    /** 事件处理函数，需要保存引用以便移除监听 */
    private handleTouchStart = (e: TouchEvent) => {
        const touch = e.touches[0];
        this.touchStartX = touch.clientX;
        this.touchStartY = touch.clientY;
        this.touchStartTime = Date.now();
        
        /** 只有在屏幕左边缘开始触摸才启用滑动检测 */
        if (this.touchStartX <= this.options.edgeWidth) {
            this.isSwiping = true;
        } else {
            this.isSwiping = false;
        }
    };
    
    private handleTouchMove = (e: TouchEvent) => {
        if (!this.isSwiping) return;
        
        const touch = e.touches[0];
        const deltaX = touch.clientX - this.touchStartX;
        const deltaY = Math.abs(touch.clientY - this.touchStartY);
        
        /** 如果垂直偏移太大，取消滑动检测 */
        if (deltaY > this.options.maxVerticalOffset) {
            this.isSwiping = false;
        }
        
        /** 如果向右滑动距离超过一半最小距离且垂直偏移不大，可以考虑预处理 */
        if (this.isSwiping && deltaX > this.options.minSwipeDistance / 2 && deltaY < this.options.maxVerticalOffset) {
            /** 可以在这里添加视觉反馈，比如显示返回指示器 */
        }
    };
    
    private handleTouchEnd = (e: TouchEvent) => {
        if (!this.isSwiping) return;
        
        const touch = e.changedTouches[0];
        const deltaX = touch.clientX - this.touchStartX;
        const deltaY = Math.abs(touch.clientY - this.touchStartY);
        const deltaTime = Date.now() - this.touchStartTime;
        
        // 判断是否满足向右滑动的条件
        if (deltaX > this.options.minSwipeDistance && 
            deltaY < this.options.maxVerticalOffset && 
            deltaTime < this.options.maxSwipeTime) {
            try {
                this.options.onBack();
            } catch (error) {
                console.warn('滑动返回失败:', error);
            }
        }
        
        this.isSwiping = false;
    };
    
    constructor(options: SlideBackOptions = {}) {
        this.options = {
            edgeWidth: 30,
            minSwipeDistance: 100,
            maxVerticalOffset: 100,
            maxSwipeTime: 2000,
            onBack: () => {},
            ...options
        };
    }
    
    /**
     * 启用滑动返回功能
     */
    enable(): void {
        if (this.isEnabled || typeof document === 'undefined') return;
        
        document.addEventListener('touchstart', this.handleTouchStart);
        document.addEventListener('touchmove', this.handleTouchMove);
        document.addEventListener('touchend', this.handleTouchEnd);
        
        this.isEnabled = true;
    }
    
    /**
     * 禁用滑动返回功能
     */
    disable(): void {
        if (!this.isEnabled || typeof document === 'undefined') return;
        
        document.removeEventListener('touchstart', this.handleTouchStart);
        document.removeEventListener('touchmove', this.handleTouchMove);
        document.removeEventListener('touchend', this.handleTouchEnd);
        
        this.isEnabled = false;
    }
    
    /**
     * 更新配置选项
     */
    updateOptions(options: Partial<SlideBackOptions>): void {
        this.options = { ...this.options, ...options };
    }
    
    /**
     * 销毁实例
     */
    destroy(): void {
        this.disable();
    }
    
    /**
     * 检查是否已启用
     */
    get enabled(): boolean {
        return this.isEnabled;
    }
}

/**
 * 创建并启用滑动返回功能的便捷函数
 * @param options 配置选项
 * @returns SlideBack 实例
 */
export const createSlideBack = (options: SlideBackOptions = {}): SlideBack => {
    const slideBack = new SlideBack(options);
    slideBack.enable();
    return slideBack;
};

/**
 * 默认导出，用于快速创建实例
 */
export default SlideBack;
