<template>
  <div class="app">
    <header class="header">
      <h1>DooTask Tools - Vite 示例</h1>
      <p>展示如何在Vite项目中使用dootask-tools</p>
    </header>

    <main class="main">
      <!-- 应用状态信息 -->
      <section class="section">
        <h2>应用状态</h2>
        <div class="status-grid">
          <div class="status-item">
            <span class="label">是否为微应用:</span>
            <span class="value" :class="{ 'true': isMicroAppRef, 'false': !isMicroAppRef }">
              {{ isMicroAppRef ? '是' : '否' }}
            </span>
          </div>
          <div class="status-item">
            <span class="label">用户ID:</span>
            <span class="value">{{ userId }}</span>
          </div>
          <div class="status-item">
            <span class="label">主题:</span>
            <span class="value">{{ themeName || '--' }}</span>
          </div>
          <div class="status-item">
            <span class="label">语言:</span>
            <span class="value">{{ languageName || '--' }}</span>
          </div>
          <div class="status-item">
            <span class="label">是否为Electron:</span>
            <span class="value" :class="{ 'true': isElectronRef, 'false': !isElectronRef }">
              {{ isElectronRef ? '是' : '否' }}
            </span>
          </div>
          <div class="status-item">
            <span class="label">是否为EEUI应用:</span>
            <span class="value" :class="{ 'true': isEEUIAppRef, 'false': !isEEUIAppRef }">
              {{ isEEUIAppRef ? '是' : '否' }}
            </span>
          </div>
        </div>
      </section>

      <!-- 功能演示 -->
      <section class="section">
        <h2>功能演示</h2>
        <div class="button-grid">
          <button @click="handlePopoutWindow" class="btn btn-primary">
            打开独立窗口
          </button>
          <button @click="handleOpenWindow" class="btn btn-secondary">
            打开新窗口
          </button>
          <button @click="handleSelectUsers" class="btn btn-success">
            选择用户
          </button>
          <button @click="handleRequestAPI" class="btn btn-info">
            测试API请求
          </button>
          <button @click="handleCloseApp" class="btn btn-danger">
            关闭应用
          </button>
          <button @click="handleBackApp" class="btn btn-warning">
            返回
          </button>
        </div>
      </section>

      <!-- 提示框演示 -->
      <section class="section">
        <h2>提示框演示</h2>
        <div class="button-grid">
          <button @click="handleOpenModal('info')" class="btn btn-primary">
            打开默认提示框
          </button>
          <button @click="handleOpenModal('warning')" class="btn btn-warning">
            打开警告提示框
          </button>
          <button @click="handleOpenModal('error')" class="btn btn-danger">
            打开错误提示框
          </button>
          <button @click="handleOpenModal('success')" class="btn btn-success">
            打开成功提示框
          </button>
          <button @click="handleOpenModal('alert')" class="btn btn-info">
            打开系统提示框
          </button>
        </div>
      </section>

      <!-- 消息框演示 -->
      <section class="section">
        <h2>消息框演示</h2>
        <div class="button-grid">
          <button @click="handleOpenMessage('info')" class="btn btn-primary">
            打开默认消息框
          </button>
          <button @click="handleOpenMessage('warning')" class="btn btn-warning">
            打开警告消息框
          </button>
          <button @click="handleOpenMessage('error')" class="btn btn-danger">
            打开错误消息框
          </button>
          <button @click="handleOpenMessage('success')" class="btn btn-success">
            打开成功消息框
          </button>
        </div>
      </section>

      <!-- 用户信息 -->
      <section class="section" v-if="userInfo">
        <h2>用户信息</h2>
        <pre class="code-block">{{ JSON.stringify(userInfo, null, 2) }}</pre>
      </section>

      <!-- 系统信息 -->
      <section class="section" v-if="systemInfo">
        <h2>系统信息</h2>
        <pre class="code-block">{{ JSON.stringify(systemInfo, null, 2) }}</pre>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
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
  messageSuccess
} from '../../src/index'

// 响应式数据
const isMicroAppRef = ref(false)
const userId = ref(0)
const themeName = ref('')
const languageName = ref('')
const isElectronRef = ref(false)
const isEEUIAppRef = ref(false)
const userInfo = ref(null)
const systemInfo = ref(null)

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
})

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

<style scoped>
.app {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.header {
  text-align: center;
  margin-bottom: 40px;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 12px;
}

.header h1 {
  margin: 0 0 10px 0;
  font-size: 2.5rem;
  font-weight: 700;
}

.header p {
  margin: 0;
  font-size: 1.1rem;
  opacity: 0.9;
}

.section {
  margin-bottom: 40px;
  padding: 24px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.section h2 {
  margin: 0 0 20px 0;
  color: #333;
  font-size: 1.5rem;
  font-weight: 600;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 8px;
  border-left: 4px solid #007bff;
}

.label {
  font-weight: 500;
  color: #495057;
}

.value {
  font-weight: 600;
  color: #212529;
}

.value.true {
  color: #28a745;
}

.value.false {
  color: #dc3545;
}

.button-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-success {
  background: #28a745;
  color: white;
}

.btn-info {
  background: #17a2b8;
  color: white;
}

.btn-danger {
  background: #dc3545;
  color: white;
}

.btn-warning {
  background: #ffc107;
  color: #212529;
}

.code-block {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.9rem;
  line-height: 1.5;
  color: #495057;
}

@media (max-width: 768px) {
  .app {
    padding: 10px;
  }
  
  .header h1 {
    font-size: 2rem;
  }
  
  .status-grid {
    grid-template-columns: 1fr;
  }
  
  .button-grid {
    grid-template-columns: 1fr;
  }
}
</style> 