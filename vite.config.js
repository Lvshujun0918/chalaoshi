import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { execSync } from 'child_process'

// 获取当前 git commit 短哈希
let gitHash = 'unknown'
try {
  gitHash = execSync('git rev-parse --short HEAD').toString().trim()
} catch (e) {
  console.warn('无法获取 git commit hash:', e.message)
}

// https://vite.dev/config/
export default defineConfig(({ command }) => {
  // 开发环境下通过代理访问时使用子路径
  const isDev = command === 'serve'
  return {
    plugins: [vue()],
    base: isDev ? '/absproxy/5173/' : '/',  // 生产环境可根据需要调整
    define: {
      __GIT_HASH__: JSON.stringify(gitHash),
    },
    server: {
      allowedHosts: ['code.cszj.wang'],
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
      },
    },
  }
})