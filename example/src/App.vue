<template>
  <div class="min-h-screen bg-gradient-to-br from-background via-background to-muted/20">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 py-4 sm:py-8">
      <nav class="sticky top-2 sm:top-4 z-50 grid grid-cols-[50px_1fr_50px] sm:grid-cols-[60px_1fr_60px] items-center p-3 sm:p-4 mb-6 sm:mb-10 bg-card/70 backdrop-blur-xl border border-border/50 rounded-xl sm:rounded-2xl shadow-lg shadow-black/5" v-if="showNav">
        <div class="flex justify-start items-center cursor-pointer p-1.5 sm:p-2 rounded-lg sm:rounded-xl transition-all hover:bg-accent" @click="handleCloseApp">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
            class="lucide lucide-x-icon lucide-x sm:w-[18px] sm:h-[18px]">
            <path d="M18 6 6 18" />
            <path d="m6 6 12 12" />
          </svg>
        </div>
        <div class="flex justify-center items-center text-base sm:text-lg font-semibold text-foreground">DooTask Tools</div>
        <div class="flex justify-end items-center"></div>
      </nav>

      <header class="text-center mb-12 sm:mb-16 p-6 sm:p-12 bg-gradient-to-br from-primary/10 via-primary/5 to-transparent border border-primary/20 rounded-3xl backdrop-blur-sm shadow-lg shadow-primary/5">
        <h1 class="text-3xl sm:text-4xl md:text-5xl font-bold mb-4 text-foreground tracking-tight">DooTask Tools</h1>
        <p class="text-lg sm:text-xl text-muted-foreground font-medium">现代化的 Vite 开发工具集成示例</p>
      </header>

          <main class="relative">
            <!-- 应用状态信息 -->
      <section class="mb-8 sm:mb-12 p-4 sm:p-8 bg-gradient-to-br from-card/80 to-card/40 backdrop-blur-lg border border-border/50 rounded-2xl sm:rounded-3xl shadow-xl shadow-black/5 hover:shadow-lg hover:shadow-black/10 transition-shadow">
        <h2 class="text-xl sm:text-2xl font-bold mb-6 sm:mb-8 text-card-foreground flex items-center gap-3">
          <div class="w-2 h-6 sm:h-8 bg-primary rounded-full"></div>
          应用状态
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4 sm:gap-6">
          <div class="group flex justify-between items-center p-4 sm:p-6 bg-gradient-to-r from-emerald-50 to-emerald-100 dark:from-emerald-900/20 dark:to-emerald-800/30 rounded-xl sm:rounded-2xl border border-emerald-200 dark:border-emerald-700/50 hover:shadow-lg hover:shadow-emerald-500/10 transition-shadow">
            <span class="font-medium text-sm sm:text-base text-emerald-800 dark:text-emerald-200">是否为微应用:</span>
            <span class="font-bold text-base sm:text-lg" :class="isMicroAppRef ? 'text-emerald-600 dark:text-emerald-400' : 'text-slate-500 dark:text-slate-400'">
              {{ isMicroAppRef ? '是' : '否' }}
            </span>
          </div>
          <div class="group flex justify-between items-center p-4 sm:p-6 bg-gradient-to-r from-blue-50 to-blue-100 dark:from-blue-900/20 dark:to-blue-800/30 rounded-xl sm:rounded-2xl border border-blue-200 dark:border-blue-700/50 hover:shadow-lg hover:shadow-blue-500/10 transition-shadow">
            <span class="font-medium text-sm sm:text-base text-blue-800 dark:text-blue-200">用户ID:</span>
            <span class="font-bold text-base sm:text-lg text-blue-600 dark:text-blue-400">{{ userId }}</span>
          </div>
          <div class="group flex justify-between items-center p-4 sm:p-6 bg-gradient-to-r from-purple-50 to-purple-100 dark:from-purple-900/20 dark:to-purple-800/30 rounded-xl sm:rounded-2xl border border-purple-200 dark:border-purple-700/50 hover:shadow-lg hover:shadow-purple-500/10 transition-shadow">
            <span class="font-medium text-sm sm:text-base text-purple-800 dark:text-purple-200">主题:</span>
            <span class="font-bold text-base sm:text-lg text-purple-600 dark:text-purple-400">{{ themeName || '--' }}</span>
          </div>
          <div class="group flex justify-between items-center p-4 sm:p-6 bg-gradient-to-r from-orange-50 to-orange-100 dark:from-orange-900/20 dark:to-orange-800/30 rounded-xl sm:rounded-2xl border border-orange-200 dark:border-orange-700/50 hover:shadow-lg hover:shadow-orange-500/10 transition-shadow">
            <span class="font-medium text-sm sm:text-base text-orange-800 dark:text-orange-200">语言:</span>
            <span class="font-bold text-base sm:text-lg text-orange-600 dark:text-orange-400">{{ languageName || '--' }}</span>
          </div>
          <div class="group flex justify-between items-center p-4 sm:p-6 bg-gradient-to-r from-cyan-50 to-cyan-100 dark:from-cyan-900/20 dark:to-cyan-800/30 rounded-xl sm:rounded-2xl border border-cyan-200 dark:border-cyan-700/50 hover:shadow-lg hover:shadow-cyan-500/10 transition-shadow">
            <span class="font-medium text-sm sm:text-base text-cyan-800 dark:text-cyan-200">是否为Electron:</span>
            <span class="font-bold text-base sm:text-lg" :class="isElectronRef ? 'text-cyan-600 dark:text-cyan-400' : 'text-slate-500 dark:text-slate-400'">
              {{ isElectronRef ? '是' : '否' }}
            </span>
          </div>
          <div class="group flex justify-between items-center p-4 sm:p-6 bg-gradient-to-r from-pink-50 to-pink-100 dark:from-pink-900/20 dark:to-pink-800/30 rounded-xl sm:rounded-2xl border border-pink-200 dark:border-pink-700/50 hover:shadow-lg hover:shadow-pink-500/10 transition-shadow">
            <span class="font-medium text-sm sm:text-base text-pink-800 dark:text-pink-200">是否为EEUI应用:</span>
            <span class="font-bold text-base sm:text-lg" :class="isEEUIAppRef ? 'text-pink-600 dark:text-pink-400' : 'text-slate-500 dark:text-slate-400'">
              {{ isEEUIAppRef ? '是' : '否' }}
            </span>
          </div>
        </div>
      </section>

      <!-- 功能演示 -->
      <section class="mb-8 sm:mb-12 p-4 sm:p-8 bg-gradient-to-br from-card/80 to-card/40 backdrop-blur-lg border border-border/50 rounded-2xl sm:rounded-3xl shadow-xl shadow-black/5 hover:shadow-lg hover:shadow-black/10 transition-shadow">
        <h2 class="text-xl sm:text-2xl font-bold mb-6 sm:mb-8 text-card-foreground flex items-center gap-3">
          <div class="w-2 h-6 sm:h-8 bg-primary rounded-full"></div>
          功能演示
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
          <button @click="handlePopoutWindow" class="group relative px-4 sm:px-6 py-3 sm:py-4 bg-gradient-to-r from-primary to-primary/80 text-primary-foreground rounded-xl sm:rounded-2xl font-semibold shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-xl sm:rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开独立窗口</span>
          </button>
          <button @click="handleOpenWindow" class="group relative px-4 sm:px-6 py-3 sm:py-4 bg-gradient-to-r from-slate-50 to-slate-100 dark:from-slate-800 dark:to-slate-700 text-slate-700 dark:text-slate-300 rounded-xl sm:rounded-2xl font-semibold border border-slate-200 dark:border-slate-600 shadow-lg shadow-slate-500/10 hover:shadow-xl hover:shadow-slate-500/20 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-xl sm:rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开新窗口</span>
          </button>
          <button @click="handleSelectUsers" class="group relative px-6 py-4 bg-gradient-to-r from-emerald-500 to-emerald-600 text-white rounded-2xl font-semibold shadow-lg shadow-emerald-500/25 hover:shadow-xl hover:shadow-emerald-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">选择用户</span>
          </button>
          <button @click="handleRequestAPI" class="group relative px-6 py-4 bg-gradient-to-r from-blue-500 to-blue-600 text-white rounded-2xl font-semibold shadow-lg shadow-blue-500/25 hover:shadow-xl hover:shadow-blue-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">测试API请求</span>
          </button>
          <button @click="handleCloseApp" class="group relative px-6 py-4 bg-gradient-to-r from-red-500 to-red-600 text-white rounded-2xl font-semibold shadow-lg shadow-red-500/25 hover:shadow-xl hover:shadow-red-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">关闭应用</span>
          </button>
          <button @click="handleBackApp" class="group relative px-6 py-4 bg-gradient-to-r from-orange-500 to-orange-600 text-white rounded-2xl font-semibold shadow-lg shadow-orange-500/25 hover:shadow-xl hover:shadow-orange-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">返回</span>
          </button>
        </div>
      </section>

      <!-- 提示框演示 -->
      <section class="mb-8 sm:mb-12 p-4 sm:p-8 bg-gradient-to-br from-card/80 to-card/40 backdrop-blur-lg border border-border/50 rounded-2xl sm:rounded-3xl shadow-xl shadow-black/5 hover:shadow-lg hover:shadow-black/10 transition-shadow">
        <h2 class="text-xl sm:text-2xl font-bold mb-6 sm:mb-8 text-card-foreground flex items-center gap-3">
          <div class="w-2 h-6 sm:h-8 bg-primary rounded-full"></div>
          提示框演示
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
          <button @click="handleOpenModal('info')" class="group relative px-6 py-4 bg-gradient-to-r from-indigo-500 to-indigo-600 text-white rounded-2xl font-semibold shadow-lg shadow-indigo-500/25 hover:shadow-xl hover:shadow-indigo-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开默认提示框</span>
          </button>
          <button @click="handleOpenModal('warning')" class="group relative px-6 py-4 bg-gradient-to-r from-amber-400 to-amber-500 text-white rounded-2xl font-semibold shadow-lg shadow-amber-400/25 hover:shadow-xl hover:shadow-amber-400/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开警告提示框</span>
          </button>
          <button @click="handleOpenModal('error')" class="group relative px-6 py-4 bg-gradient-to-r from-rose-500 to-rose-600 text-white rounded-2xl font-semibold shadow-lg shadow-rose-500/25 hover:shadow-xl hover:shadow-rose-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开错误提示框</span>
          </button>
          <button @click="handleOpenModal('success')" class="group relative px-6 py-4 bg-gradient-to-r from-teal-500 to-teal-600 text-white rounded-2xl font-semibold shadow-lg shadow-teal-500/25 hover:shadow-xl hover:shadow-teal-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开成功提示框</span>
          </button>
          <button @click="handleOpenModal('confirm')" class="group relative px-6 py-4 bg-gradient-to-r from-sky-500 to-sky-600 text-white rounded-2xl font-semibold shadow-lg shadow-sky-500/25 hover:shadow-xl hover:shadow-sky-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开确认提示框</span>
          </button>
          <button @click="handleOpenModal('alert')" class="group relative px-6 py-4 bg-gradient-to-r from-gray-400 to-gray-500 text-white rounded-2xl font-semibold shadow-lg shadow-gray-400/25 hover:shadow-xl hover:shadow-gray-400/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开系统提示框</span>
          </button>
        </div>
      </section>

      <!-- 消息框演示 -->
      <section class="mb-8 sm:mb-12 p-4 sm:p-8 bg-gradient-to-br from-card/80 to-card/40 backdrop-blur-lg border border-border/50 rounded-2xl sm:rounded-3xl shadow-xl shadow-black/5 hover:shadow-lg hover:shadow-black/10 transition-shadow">
        <h2 class="text-xl sm:text-2xl font-bold mb-6 sm:mb-8 text-card-foreground flex items-center gap-3">
          <div class="w-2 h-6 sm:h-8 bg-primary rounded-full"></div>
          消息框演示
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
          <button @click="handleOpenMessage('info')" class="group relative px-6 py-4 bg-gradient-to-r from-violet-500 to-violet-600 text-white rounded-2xl font-semibold shadow-lg shadow-violet-500/25 hover:shadow-xl hover:shadow-violet-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开默认消息框</span>
          </button>
          <button @click="handleOpenMessage('warning')" class="group relative px-6 py-4 bg-gradient-to-r from-yellow-500 to-yellow-600 text-white rounded-2xl font-semibold shadow-lg shadow-yellow-500/25 hover:shadow-xl hover:shadow-yellow-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开警告消息框</span>
          </button>
          <button @click="handleOpenMessage('error')" class="group relative px-6 py-4 bg-gradient-to-r from-red-600 to-red-700 text-white rounded-2xl font-semibold shadow-lg shadow-red-600/25 hover:shadow-xl hover:shadow-red-600/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开错误消息框</span>
          </button>
          <button @click="handleOpenMessage('success')" class="group relative px-6 py-4 bg-gradient-to-r from-green-500 to-green-600 text-white rounded-2xl font-semibold shadow-lg shadow-green-500/25 hover:shadow-xl hover:shadow-green-500/30 transition-shadow">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span class="relative">打开成功消息框</span>
          </button>
        </div>
      </section>

      <!-- 用户信息 -->
      <section class="mb-8 sm:mb-12 p-4 sm:p-8 bg-gradient-to-br from-card/80 to-card/40 backdrop-blur-lg border border-border/50 rounded-2xl sm:rounded-3xl shadow-xl shadow-black/5 hover:shadow-lg hover:shadow-black/10 transition-shadow" v-if="userInfo">
        <h2 class="text-xl sm:text-2xl font-bold mb-6 sm:mb-8 text-card-foreground flex items-center gap-3">
          <div class="w-2 h-6 sm:h-8 bg-primary rounded-full"></div>
          用户信息
        </h2>
        <div class="bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl sm:rounded-2xl p-4 sm:p-6 shadow-inner">
          <pre class="overflow-x-auto font-mono text-xs sm:text-sm leading-6 text-slate-700 dark:text-slate-300">{{ JSON.stringify(userInfo, null, 2) }}</pre>
        </div>
      </section>

      <!-- 系统信息 -->
      <section class="mb-8 sm:mb-12 p-4 sm:p-8 bg-gradient-to-br from-card/80 to-card/40 backdrop-blur-lg border border-border/50 rounded-2xl sm:rounded-3xl shadow-xl shadow-black/5 hover:shadow-lg hover:shadow-black/10 transition-shadow" v-if="systemInfo">
        <h2 class="text-xl sm:text-2xl font-bold mb-6 sm:mb-8 text-card-foreground flex items-center gap-3">
          <div class="w-2 h-6 sm:h-8 bg-primary rounded-full"></div>
          系统信息
        </h2>
        <div class="bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl sm:rounded-2xl p-4 sm:p-6 shadow-inner">
          <pre class="overflow-x-auto font-mono text-xs sm:text-sm leading-6 text-slate-700 dark:text-slate-300">{{ JSON.stringify(systemInfo, null, 2) }}</pre>
        </div>
      </section>
      
      <!-- 装饰性元素 -->
      <div class="fixed top-10 right-10 w-20 sm:w-32 h-20 sm:h-32 bg-primary/10 rounded-full blur-2xl animate-pulse pointer-events-none"></div>
      <div class="fixed bottom-10 left-10 w-16 sm:w-20 h-16 sm:h-20 bg-secondary/10 rounded-full blur-xl animate-pulse pointer-events-none" style="animation-delay: 1s;"></div>
    </main>
    </div>
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
  modalConfirm,
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
    case 'confirm':
      modalConfirm({
        title: '确认提示框',
        content: '确认提示框内容',
      }).then(res => {
        if (res) {
          messageSuccess('确认')
        } else {
          messageInfo('取消')
        }
      })
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
