// Self-hosted fonts (replaces Google Fonts CDN)
import '@fontsource/plus-jakarta-sans/400.css'
import '@fontsource/plus-jakarta-sans/500.css'
import '@fontsource/plus-jakarta-sans/600.css'
import '@fontsource/plus-jakarta-sans/700.css'
import '@fontsource-variable/noto-sans-sc'

import '@/styles/themes.css'
import '@/styles/global.css'
import '@/styles/components.css'
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'
import App from './App.vue'
import router from './router'

const app = createApp(App)

// 全局 Vue 错误处理器
app.config.errorHandler = (err, instance, info) => {
  console.error('[Vue Error]', err)
  console.error('[Vue Error Info]', info)
  console.error('[Vue Error Component]', instance)
}

app.use(createPinia())
app.use(router)
app.use(Antd)
app.mount('#app')
