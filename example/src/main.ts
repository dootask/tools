import { createApp } from "vue"
import App from "./App.vue"
import "./style.css"

const theme = localStorage.getItem('theme')
if (theme) {
  document.documentElement.classList.add(theme)
}

createApp(App).mount("#app")
