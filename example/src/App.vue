<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- 导航栏 -->
    <nav v-if="showNav" class="sticky top-0 z-50 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-4 py-3">
      <div class="flex items-center justify-between max-w-7xl mx-auto">
        <button @click="handleCloseApp" class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
        <h1 class="text-lg font-semibold text-gray-900 dark:text-white">DooTask Tools</h1>
        <div class="w-9"></div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto px-4 py-6">
      <!-- 页面标题 -->
      <header class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">DooTask Tools</h1>
        <p class="text-gray-600 dark:text-gray-400">现代化的 Vite 开发工具集成示例</p>
      </header>

      <!-- 应用状态信息 -->
      <section class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">应用状态</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600 dark:text-gray-400">是否为微应用</span>
              <span class="font-medium" :class="isMicroAppRef ? 'text-green-600' : 'text-gray-500'">
                {{ isMicroAppRef ? '是' : '否' }}
              </span>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600 dark:text-gray-400">用户ID</span>
              <span class="font-medium text-blue-600">{{ userId }}</span>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600 dark:text-gray-400">主题</span>
              <span class="font-medium text-purple-600">{{ themeName || '--' }}</span>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600 dark:text-gray-400">语言</span>
              <span class="font-medium text-orange-600">{{ languageName || '--' }}</span>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600 dark:text-gray-400">是否为Electron</span>
              <span class="font-medium" :class="isElectronRef ? 'text-green-600' : 'text-gray-500'">
                {{ isElectronRef ? '是' : '否' }}
              </span>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600 dark:text-gray-400">是否为EEUI应用</span>
              <span class="font-medium" :class="isEEUIAppRef ? 'text-green-600' : 'text-gray-500'">
                {{ isEEUIAppRef ? '是' : '否' }}
              </span>
            </div>
          </div>
        </div>
      </section>

      <!-- 功能演示 -->
      <section class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">功能演示</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <button @click="handlePopoutWindow" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-3 rounded-lg font-medium">
            打开独立窗口
          </button>
          <button @click="handleOpenWindow" class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-3 rounded-lg font-medium">
            打开新窗口
          </button>
          <button @click="handleSelectUsers" class="bg-green-600 hover:bg-green-700 text-white px-4 py-3 rounded-lg font-medium">
            选择用户
          </button>
          <button @click="handleRequestAPI" class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-3 rounded-lg font-medium">
            测试API请求
          </button>
          <button @click="handleCloseApp" class="bg-red-600 hover:bg-red-700 text-white px-4 py-3 rounded-lg font-medium">
            关闭应用
          </button>
          <button @click="handleBackApp" class="bg-orange-600 hover:bg-orange-700 text-white px-4 py-3 rounded-lg font-medium">
            返回
          </button>

          <!-- 是否阻止关闭应用 -->
          <button
            @click="preventCloseApp = !preventCloseApp"
            class="bg-gray-600 hover:bg-gray-700 px-4 py-3 rounded-lg font-medium flex items-center justify-center gap-2"
            :aria-pressed="preventCloseApp"
            type="button"
          >
            <span>
              <svg v-if="preventCloseApp" class="h-5 w-5 text-blue-600" fill="currentColor" viewBox="0 0 20 20">
                <rect width="20" height="20" rx="4" fill="currentColor"/>
                <polyline points="5 11 9 15 15 7" fill="none" stroke="#fff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
              <svg v-else class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 20 20">
                <rect x="2" y="2" width="16" height="16" rx="4" stroke="currentColor" stroke-width="2" fill="none"/>
              </svg>
            </span>
            <span class="text-white">阻止返回</span>
          </button>
        </div>
      </section>

      <!-- 提示框演示 -->
      <section class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">提示框演示</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <button @click="handleOpenModal('info')" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-3 rounded-lg font-medium">
            打开默认提示框
          </button>
          <button @click="handleOpenModal('warning')" class="bg-yellow-600 hover:bg-yellow-700 text-white px-4 py-3 rounded-lg font-medium">
            打开警告提示框
          </button>
          <button @click="handleOpenModal('error')" class="bg-red-600 hover:bg-red-700 text-white px-4 py-3 rounded-lg font-medium">
            打开错误提示框
          </button>
          <button @click="handleOpenModal('success')" class="bg-green-600 hover:bg-green-700 text-white px-4 py-3 rounded-lg font-medium">
            打开成功提示框
          </button>
          <button @click="handleOpenModal('confirm')" class="bg-purple-600 hover:bg-purple-700 text-white px-4 py-3 rounded-lg font-medium">
            打开确认提示框
          </button>
          <button @click="handleOpenModal('alert')" class="bg-gray-600 hover:bg-gray-700 text-white px-4 py-3 rounded-lg font-medium">
            打开系统提示框
          </button>
        </div>
      </section>

      <!-- 消息框演示 -->
      <section class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">消息框演示</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <button @click="handleOpenMessage('info')" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-3 rounded-lg font-medium">
            打开默认消息框
          </button>
          <button @click="handleOpenMessage('warning')" class="bg-yellow-600 hover:bg-yellow-700 text-white px-4 py-3 rounded-lg font-medium">
            打开警告消息框
          </button>
          <button @click="handleOpenMessage('error')" class="bg-red-600 hover:bg-red-700 text-white px-4 py-3 rounded-lg font-medium">
            打开错误消息框
          </button>
          <button @click="handleOpenMessage('success')" class="bg-green-600 hover:bg-green-700 text-white px-4 py-3 rounded-lg font-medium">
            打开成功消息框
          </button>
        </div>
      </section>

      <!-- 用户信息 -->
      <section v-if="userInfo" class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">用户信息</h2>
        <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
          <pre class="text-sm text-gray-700 dark:text-gray-300 overflow-x-auto">{{ JSON.stringify(userInfo, null, 2) }}</pre>
        </div>
      </section>

      <!-- 系统信息 -->
      <section v-if="systemInfo" class="mb-8">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">系统信息</h2>
        <div class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
          <pre class="text-sm text-gray-700 dark:text-gray-300 overflow-x-auto">{{ JSON.stringify(systemInfo, null, 2) }}</pre>
        </div>
      </section>
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
const preventCloseApp = ref(true)
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
    if (preventCloseApp.value) {
      modalInfo("阻止返回")
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
  closeApp()
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
