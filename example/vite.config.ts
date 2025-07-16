import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import tailwindcss from "@tailwindcss/vite"

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss()
  ],
  server: {
    host: "0.0.0.0",
  },
  build: {
    outDir: "dist",
    sourcemap: true,
  },
})
