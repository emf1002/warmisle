import { ref, watch } from 'vue'
import { defineStore } from 'pinia'

export type ThemeMode = 'light' | 'dark'

export const useThemeStore = defineStore('theme', () => {
  const STORAGE_KEY = 'warmisle-theme'

  // 从 localStorage 读取，默认 light
  function getInitialTheme(): ThemeMode {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved === 'light' || saved === 'dark') return saved
    return 'light'
  }

  const theme = ref<ThemeMode>(getInitialTheme())

  // 切换主题
  function toggleTheme() {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  // 设置指定主题
  function setTheme(mode: ThemeMode) {
    theme.value = mode
  }

  // 同步到 <html data-theme="..."> 和 localStorage
  function applyTheme(mode: ThemeMode) {
    document.documentElement.setAttribute('data-theme', mode)
    localStorage.setItem(STORAGE_KEY, mode)
  }

  // 初始化时立即应用
  applyTheme(theme.value)

  // 监听变化自动应用
  watch(theme, (newTheme) => {
    applyTheme(newTheme)
  })

  return { theme, toggleTheme, setTheme }
})
