/** 滑动返回 */
export class SlideBack {
  /** 容器元素 */
  private slideBackContainer: HTMLDivElement | null = null
  /** 回调函数 */
  private slideBackCallback: (() => void) | null = null

  /** 窗口高度 */
  private windowHeight = 0
  /** 窗口滚动位置 */
  private windowScrollY = 0

  /** 触摸位置 */
  private touchX = 0
  private touchY = 0

  /** 触摸开始位置 */
  private startX = 0
  private startY = 0

  /** 是否显示滑动返回容器 */
  private isVisible = false
  /** 是否触摸 */
  private isTouched = false
  /** 是否滚动 */
  private isScrolling: boolean | undefined = undefined

  /** 初始化样式 */
  private initSlideStyle = () => {
    if (document.querySelector("style.slide-back-style")) {
      return
    }
    const style = document.createElement("style")
    style.id = "slide-back-style"
    style.textContent = `
      .slide-back-container {
        position: fixed;
        top: 200px;
        left: -50px;
        width: 500px;
        height: 500px;
        background-color: rgba(0, 0, 0, 0.1);
        z-index: 9999;
        border-radius: 50%;
        transform: translate(-460px, -50%);
        transition: left 0.2s ease;
      }

      .slide-back-container.visible {
        left: 0;
      }
      `
    document.head.appendChild(style)
  }

  /** 初始化容器 */
  private initContainer = () => {
    if (this.slideBackContainer) {
      return
    }
    this.slideBackContainer = document.createElement("div")
    this.slideBackContainer.classList.add("slide-back-container")
    document.body.appendChild(this.slideBackContainer)
  }

  /** 初始化滑动返回 */
  private initSlideBack = () => {
    document.addEventListener("touchstart", this.touchstart)
    document.addEventListener("touchmove", this.touchmove, { passive: false })
    document.addEventListener("touchend", this.touchend)
  }

  /** 获取触摸位置 */
  private getXY = (event: TouchEvent) => {
    const touch = event.touches[0]
    this.touchX = touch.clientX
    this.touchY = touch.clientY
    this.updateContainer()
  }

  /** 更新滑动返回容器 */
  private updateContainer = () => {
    if (!this.slideBackContainer) {
      return
    }

    const isVisible = this.isVisibleSlideBack()
    if (isVisible) {
      const offset = 135
      const top = Math.max(offset, this.touchY) + this.windowScrollY
      const maxTop = this.windowHeight - offset
      this.slideBackContainer.style.top = `${Math.min(top, maxTop) + "px"}`
    }
    this.slideBackContainer.classList.toggle("visible", isVisible)
  }

  /** 是否显示滑动返回容器 */
  private isVisibleSlideBack = () => {
    return this.isVisible && this.touchX > 0
  }

  /** 触摸开始 */
  private touchstart = (e: TouchEvent) => {
    this.getXY(e)

    this.isTouched = this.touchX < 20
    this.isScrolling = undefined

    this.startX = e.targetTouches[0].pageX
    this.startY = e.targetTouches[0].pageY
  }

  /** 触摸移动 */
  private touchmove = (e: TouchEvent) => {
    if (!this.isTouched) {
      return
    }
    const pageX = e.targetTouches[0].pageX
    const pageY = e.targetTouches[0].pageY
    if (typeof this.isScrolling === "undefined") {
      const verticalMove = Math.abs(pageY - this.startY)
      const horizontalMove = Math.abs(pageX - this.startX) * 1.5
      this.isScrolling = verticalMove > horizontalMove
    }
    if (this.isScrolling) {
      this.isTouched = false
      return
    }
    this.isVisible = true
    this.getXY(e)
    e.preventDefault()
  }

  /** 触摸结束 */
  private touchend = (e: TouchEvent) => {
    // 判断停止时的位置偏移
    if (this.touchX > 90 && this.isVisible) {
      this.slideBackCallback?.()
    }
    this.touchX = 0
    this.isVisible = false
    this.updateContainer()
  }

  /** 滚动事件 */
  private scroll = () => {
    this.windowScrollY = window.scrollY
  }

  /** 构造函数 */
  constructor(callback: () => void) {
    this.slideBackCallback = callback
    this.initSlideStyle()
    this.initContainer()
    this.initSlideBack()

    this.windowHeight = window.innerHeight
    window.addEventListener("scroll", this.scroll)
  }

  /** 销毁 */
  public destroy = () => {
    document.removeEventListener("touchstart", this.touchstart)
    document.removeEventListener("touchmove", this.touchmove)
    document.removeEventListener("touchend", this.touchend)
    window.removeEventListener("scroll", this.scroll)
  }
}
/**
 * 自动初始化滑动返回功能
 */
export const initSlideBack = (callback: () => void): SlideBack => {
  return new SlideBack(callback)
}
