<template>
  <div class="max-w-6xl mx-auto p-5 bg-background text-foreground">
    <nav class="sticky top-2.5 z-50 grid grid-cols-[60px_1fr_60px] items-center p-2.5 mb-5 bg-card border border-border rounded-lg shadow-sm" v-if="showNav">
      <div class="flex justify-start items-center cursor-pointer p-2 rounded-sm transition-colors hover:bg-accent" @click="handleCloseApp">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
          stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
          class="lucide lucide-x-icon lucide-x">
          <path d="M18 6 6 18" />
          <path d="m6 6 12 12" />
        </svg>
      </div>
      <div class="flex justify-center items-center text-xl font-semibold text-foreground">DooTask Tools</div>
      <div class="flex justify-end items-center"></div>
    </nav>

    <header class="text-center mb-10 p-5 bg-primary text-primary-foreground rounded-lg shadow-md">
      <h1 class="text-4xl font-bold mb-2.5">DooTask Tools - Vite 示例</h1>
      <p class="text-lg opacity-90">展示如何在Vite项目中使用dootask-tools</p>
    </header>

    <main>
      <!-- 应用状态信息 -->
      <section class="mb-10 p-6 bg-card border border-border rounded-lg shadow-sm">
        <h2 class="text-2xl font-semibold mb-5 text-card-foreground">应用状态</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
          <div class="flex justify-between items-center p-3 px-4 bg-muted rounded-sm border-l-4 border-primary">
            <span class="font-medium text-muted-foreground">是否为微应用:</span>
            <span class="font-semibold" :class="isMicroAppRef ? 'text-green-600' : 'text-destructive'">
              {{ isMicroAppRef ? '是' : '否' }}
            </span>
          </div>
          <div class="flex justify-between items-center p-3 px-4 bg-muted rounded-sm border-l-4 border-primary">
            <span class="font-medium text-muted-foreground">用户ID:</span>
            <span class="font-semibold text-foreground">{{ userId }}</span>
          </div>
          <div class="flex justify-between items-center p-3 px-4 bg-muted rounded-sm border-l-4 border-primary">
            <span class="font-medium text-muted-foreground">主题:</span>
            <span class="font-semibold text-foreground">{{ themeName || '--' }}</span>
          </div>
          <div class="flex justify-between items-center p-3 px-4 bg-muted rounded-sm border-l-4 border-primary">
            <span class="font-medium text-muted-foreground">语言:</span>
            <span class="font-semibold text-foreground">{{ languageName || '--' }}</span>
          </div>
          <div class="flex justify-between items-center p-3 px-4 bg-muted rounded-sm border-l-4 border-primary">
            <span class="font-medium text-muted-foreground">是否为Electron:</span>
            <span class="font-semibold" :class="isElectronRef ? 'text-green-600' : 'text-destructive'">
              {{ isElectronRef ? '是' : '否' }}
            </span>
          </div>
          <div class="flex justify-between items-center p-3 px-4 bg-muted rounded-sm border-l-4 border-primary">
            <span class="font-medium text-muted-foreground">是否为EEUI应用:</span>
            <span class="font-semibold" :class="isEEUIAppRef ? 'text-green-600' : 'text-destructive'">
              {{ isEEUIAppRef ? '是' : '否' }}
            </span>
          </div>
        </div>
      </section>

      <!-- 功能演示 -->
      <section class="mb-10 p-6 bg-card border border-border rounded-lg shadow-sm">
        <h2 class="text-2xl font-semibold mb-5 text-card-foreground">功能演示</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <button @click="handlePopoutWindow" class="px-6 py-3 bg-primary text-primary-foreground rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开独立窗口
          </button>
          <button @click="handleOpenWindow" class="px-6 py-3 bg-secondary text-secondary-foreground rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开新窗口
          </button>
          <button @click="handleSelectUsers" class="px-6 py-3 bg-green-600 text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            选择用户
          </button>
          <button @click="handleRequestAPI" class="px-6 py-3 bg-blue-600 text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            测试API请求
          </button>
          <button @click="handleCloseApp" class="px-6 py-3 bg-destructive text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            关闭应用
          </button>
          <button @click="handleBackApp" class="px-6 py-3 bg-orange-500 text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            返回
          </button>
        </div>
      </section>

      <!-- 提示框演示 -->
      <section class="mb-10 p-6 bg-card border border-border rounded-lg shadow-sm">
        <h2 class="text-2xl font-semibold mb-5 text-card-foreground">提示框演示</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <button @click="handleOpenModal('info')" class="px-6 py-3 bg-primary text-primary-foreground rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开默认提示框
          </button>
          <button @click="handleOpenModal('warning')" class="px-6 py-3 bg-orange-500 text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开警告提示框
          </button>
          <button @click="handleOpenModal('error')" class="px-6 py-3 bg-destructive text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开错误提示框
          </button>
          <button @click="handleOpenModal('success')" class="px-6 py-3 bg-green-600 text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开成功提示框
          </button>
          <button @click="handleOpenModal('alert')" class="px-6 py-3 bg-blue-600 text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开系统提示框
          </button>
        </div>
      </section>

      <!-- 消息框演示 -->
      <section class="mb-10 p-6 bg-card border border-border rounded-lg shadow-sm">
        <h2 class="text-2xl font-semibold mb-5 text-card-foreground">消息框演示</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <button @click="handleOpenMessage('info')" class="px-6 py-3 bg-primary text-primary-foreground rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开默认消息框
          </button>
          <button @click="handleOpenMessage('warning')" class="px-6 py-3 bg-orange-500 text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开警告消息框
          </button>
          <button @click="handleOpenMessage('error')" class="px-6 py-3 bg-destructive text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开错误消息框
          </button>
          <button @click="handleOpenMessage('success')" class="px-6 py-3 bg-green-600 text-white rounded-sm font-medium uppercase tracking-wide shadow-sm hover:-translate-y-0.5 hover:shadow-md transition-all duration-200">
            打开成功消息框
          </button>
        </div>
      </section>

      <!-- 用户信息 -->
      <section class="mb-10 p-6 bg-card border border-border rounded-lg shadow-sm" v-if="userInfo">
        <h2 class="text-2xl font-semibold mb-5 text-card-foreground">用户信息</h2>
        <pre class="bg-muted border border-border rounded-sm p-4 overflow-x-auto font-mono text-sm leading-6 text-muted-foreground">{{ JSON.stringify(userInfo, null, 2) }}</pre>
      </section>

      <!-- 系统信息 -->
      <section class="mb-10 p-6 bg-card border border-border rounded-lg shadow-sm" v-if="systemInfo">
        <h2 class="text-2xl font-semibold mb-5 text-card-foreground">系统信息</h2>
        <pre class="bg-muted border border-border rounded-sm p-4 overflow-x-auto font-mono text-sm leading-6 text-muted-foreground">{{ JSON.stringify(systemInfo, null, 2) }}</pre>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from "vue"
import {
  isMicroApp,
  getUserId,
  getThemeName,
  getLanguageName,
  isElectron,
  isEEUIApp,
  getSystemInfo,
  getUserInfo,
  appReady,
  popoutWindow,
  openWindow,
  selectUsers,
  requestAPI,
  closeApp,
  backApp,
  modalInfo,
  modalWarning,
  modalError,
  modalSuccess,
  modalAlert,
  messageInfo,
  messageWarning,
  messageError,
  messageSuccess,
  DooTaskUserInfo,
  DooTaskSystemInfo,
  isFullScreen,
  getWindowType,
  interceptBack
} from '../../src/index'

// 响应式数据
const isMicroAppRef = ref(false)
const userId = ref(0)
const themeName = ref('')
const languageName = ref('')
const isElectronRef = ref(false)
const isEEUIAppRef = ref(false)
const showNav = ref(false)
const preventClose = ref(true)
const userInfo = ref<DooTaskUserInfo | null>(null)
const systemInfo = ref<DooTaskSystemInfo | null>(null)

// 初始化应用
onMounted(async () => {
  try {
    // 等待应用准备就绪
    const appData = await appReady()
    console.log('应用已准备就绪:', appData)

    // 更新状态
    isMicroAppRef.value = await isMicroApp()
    userId.value = await getUserId()
    themeName.value = await getThemeName()
    languageName.value = await getLanguageName()
    isElectronRef.value = await isElectron()
    isEEUIAppRef.value = await isEEUIApp()
    userInfo.value = await getUserInfo()
    systemInfo.value = await getSystemInfo()
  } catch (error) {
    console.error('应用初始化失败:', error)
  }

  // 阻止返回
  interceptBack(() => {
    if (preventClose.value) {
      preventClose.value = false
      modalInfo("阻止返回，再次点击将关闭应用")
      return true
    } else {
      return false
    }
  })

  // 检查尺寸发生变化
  handleResize()
  window.addEventListener('resize', handleResize)
})

// 卸载时移除事件监听
onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
})

// 监听主题变化
watch(themeName, (newTheme: string) => {
  document.documentElement.classList.toggle('dark', newTheme === 'dark')
})

// 检查尺寸发生变化
const handleResize = async () => {
  // 如果当前是嵌入式窗口，并且是满屏，则显示导航栏
  showNav.value = (await getWindowType()) === 'embed' && (await isFullScreen())
}

// 处理函数
const handlePopoutWindow = () => {
  popoutWindow({
    title: '独立窗口示例',
    width: 800,
    height: 600,
    url: window.location.href
  })
}

const handleOpenWindow = () => {
  openWindow({
    name: 'example-window',
    url: 'https://example.com',
    config: {
      title: '新窗口示例',
      width: 800,
      height: 600
    }
  })
}

const handleSelectUsers = async () => {
  try {
    const result = await selectUsers({
      title: '选择用户',
      placeholder: '搜索用户...',
      multipleMax: 5
    })
    console.log('选择的用户:', result)
    modalAlert(`选择了 ${result.length} 个用户`)
  } catch (error) {
    console.error('选择用户失败:', error)
  }
}

const handleRequestAPI = async () => {
  try {
    const result = await requestAPI({
      url: 'users/info',
      method: 'GET',
      spinner: true
    })
    console.log('API请求结果:', result)
    modalInfo({
      title: 'API请求成功',
      content: '<pre style="white-space:pre-wrap">' + JSON.stringify(result.data, null, 2) + '</pre>',
    })
  } catch (error) {
    console.error('API请求失败:', error)
  }
}

const handleCloseApp = () => {
  closeApp(true)
}

const handleBackApp = () => {
  backApp()
}

const handleOpenModal = (type: string) => {
  switch (type) {
    case 'info':
      modalInfo('info')
      break;
    case 'warning':
      modalWarning('warning')
      break;
    case 'error':
      modalError('error')
      break;
    case 'success':
      modalSuccess('success')
      break;
    case 'alert':
      modalAlert('alert')
      break;
  }
}

const handleOpenMessage = (type: string) => {
  switch (type) {
    case 'info':
      messageInfo('info')
      break;
    case 'warning':
      messageWarning('warning')
      break;
    case 'error':
      messageError('error')
      break;
    case 'success':
      messageSuccess('success')
      break;
  }
}
</script>
