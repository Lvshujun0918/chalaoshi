import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig(({ command }) => {
  // 开发环境下通过代理访问时使用子路径
  const isDev = command === 'serve'
  return {
    plugins: [vue()],
    base: isDev ? '/absproxy/5173/' : '/',  // 生产环境可根据需要调整
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