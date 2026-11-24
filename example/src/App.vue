<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <div class="max-w-5xl mx-auto px-4 py-6 flex flex-col gap-6">
      <header>
        <h1 class="text-2xl font-bold">DooTask Tools 2.x 示例</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400">
          把这个页面配置成应用的窗口入口，在 DooTask 里打开后逐个试接口
        </p>
      </header>

      <section class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div v-for="item in status" :key="item.label" class="rounded-lg bg-white dark:bg-gray-800 p-3">
          <div class="text-xs text-gray-500 dark:text-gray-400">{{ item.label }}</div>
          <div class="mt-1 font-medium break-all">{{ item.value }}</div>
        </div>
      </section>

      <section v-for="group in groups" :key="group.title">
        <h2 class="mb-2 font-semibold">{{ group.title }}</h2>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="action in group.actions"
            :key="action.label"
            class="px-3 py-1.5 rounded-md text-sm bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 hover:border-blue-500 disabled:opacity-40"
            :disabled="!capabilities.includes(action.name)"
            @click="run(action)"
          >
            {{ action.label }}
          </button>
        </div>
      </section>

      <section>
        <h2 class="mb-2 font-semibold">日志</h2>
        <pre class="rounded-lg bg-white dark:bg-gray-800 p-3 text-xs h-64 overflow-auto whitespace-pre-wrap">{{
          logs.join("\n")
        }}</pre>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue"
import * as doo from "../../src/index"

interface Action {
  /** 接口名，当前环境不支持时按钮置灰 */
  name: string
  label: string
  handler: () => Promise<unknown>
}

const env = ref<doo.Env | null>(null)
const theme = ref("")
const language = ref("")
const user = ref("")
const capabilities = ref<string[]>([])
const logs = ref<string[]>([])
const cleanups: Array<() => void> = []
let offBeforeClose: (() => void) | null = null

const status = computed(() => [
  { label: "通道", value: env.value?.channel ?? "不在 DooTask 里" },
  { label: "呈现方式", value: env.value?.presentation ?? "--" },
  { label: "平台", value: env.value ? `${env.value.platform}${env.value.mobile ? "（移动端）" : ""}` : "--" },
  { label: "当前用户", value: user.value || "--" },
  { label: "主题", value: theme.value || "--" },
  { label: "语言", value: language.value || "--" },
  { label: "可用接口", value: env.value ? `${capabilities.value.length} 个` : "--" },
])

function log(label: string, value?: unknown) {
  const time = new Date().toLocaleTimeString()
  const text = value === undefined ? "" : ` ${typeof value === "string" ? value : JSON.stringify(value)}`
  logs.value.unshift(`[${time}] ${label}${text}`)
}

async function run(action: Action) {
  try {
    log(`${action.label} →`, await action.handler())
  } catch (error) {
    if (error instanceof doo.DooTaskError) {
      log(`${action.label} ✗ ${error.code}`, error.message)
    } else {
      log(`${action.label} ✗`, String(error))
    }
  }
}

function applyTheme(value: doo.Theme) {
  theme.value = value
  document.documentElement.classList.toggle("dark", value === "dark")
}

async function toggleBeforeClose() {
  if (offBeforeClose) {
    offBeforeClose()
    offBeforeClose = null
    return "已关闭"
  }
  offBeforeClose = doo.onBeforeClose(() => doo.dialog.confirm({ title: "确定关闭吗？" }))
  return "已开启"
}

const groups: Array<{ title: string; actions: Action[] }> = [
  {
    title: "身份和请求",
    actions: [
      { name: "getAppToken", label: "应用令牌", handler: async () => `${(await doo.getAppToken()).slice(0, 12)}…` },
      { name: "request", label: "请求 users/info", handler: () => doo.request("users/info") },
      { name: "fetchUsers", label: "查用户 1", handler: () => doo.fetchUsers([1]) },
      { name: "getUserInfo", label: "当前用户", handler: () => doo.getUserInfo() },
    ],
  },
  {
    title: "界面",
    actions: [
      { name: "toast.success", label: "成功提示", handler: () => doo.toast.success("操作成功") },
      { name: "toast.error", label: "错误提示", handler: () => doo.toast.error("操作失败") },
      {
        name: "dialog.confirm",
        label: "确认框",
        handler: () => doo.dialog.confirm({ title: "删除", content: "确定删除这条记录吗？" }),
      },
      {
        name: "selectUsers",
        label: "选人（最多 3 个）",
        handler: () => doo.selectUsers({ max: 3, title: "选择成员", beforeSubmit: ids => ids.users.length > 0 }),
      },
      {
        name: "selectImage",
        label: "选图片（512 裁满）",
        handler: () => doo.selectImage({ width: 512, height: 512, crop: true }),
      },
      {
        name: "previewImage",
        label: "预览图片",
        handler: () => doo.previewImage(`${location.origin}/favicon.ico`),
      },
    ],
  },
  {
    title: "操作 DooTask",
    actions: [
      { name: "addTask", label: "新建任务", handler: () => doo.addTask({ name: "来自示例应用" }) },
      { name: "openSearch", label: "搜索", handler: () => doo.openSearch("任务") },
      { name: "openAI", label: "AI 助手", handler: () => doo.openAI({ prompt: "总结一下今天的任务" }) },
      { name: "startMeeting", label: "发起会议", handler: () => doo.startMeeting() },
    ],
  },
  {
    title: "窗口",
    actions: [
      { name: "popout", label: "弹出窗口", handler: () => doo.popout() },
      {
        name: "openWindow",
        label: "打开外部窗口",
        handler: () => doo.openWindow("https://www.dootask.com", { title: "DooTask" }),
      },
      { name: "close", label: "关闭", handler: () => doo.close() },
      { name: "getCapsuleRect", label: "胶囊位置", handler: () => doo.getCapsuleRect() },
      {
        name: "setCapsuleMenus",
        label: "胶囊菜单",
        handler: () => doo.setCapsuleMenus([{ label: "刷新数据", value: "refresh" }]),
      },
      { name: "close", label: "关闭前询问（开/关）", handler: toggleBeforeClose },
    ],
  },
]

onMounted(async () => {
  try {
    env.value = await doo.ready()
  } catch (error) {
    log("不在 DooTask 里", error instanceof doo.DooTaskError ? error.code : String(error))
    return
  }
  capabilities.value = await doo.capabilities()
  applyTheme(await doo.getTheme())
  language.value = await doo.getLanguage()
  doo.getUserInfo().then(info => (user.value = `${info.nickname}（${info.userid}）`))

  cleanups.push(
    doo.onThemeChange(applyTheme),
    doo.onLanguageChange(value => (language.value = value)),
    doo.onPresentationChange(value => {
      if (env.value) {
        env.value.presentation = value
      }
    }),
    doo.onCapsuleMenu(value => log("胶囊菜单", value)),
    doo.onCapsuleChange(rect => log("胶囊位置", rect))
  )
})

onUnmounted(() => {
  cleanups.forEach(off => off())
  offBeforeClose?.()
})
</script>
