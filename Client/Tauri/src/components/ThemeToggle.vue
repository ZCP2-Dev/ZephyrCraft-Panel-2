<template>
  <button class="theme-toggle" @click="toggleTheme" :title="themeTitle">
    <svg v-if="isDark" class="theme-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="12" r="5"/>
      <line x1="12" y1="1" x2="12" y2="3"/>
      <line x1="12" y1="21" x2="12" y2="23"/>
      <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
      <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
      <line x1="1" y1="12" x2="3" y2="12"/>
      <line x1="21" y1="12" x2="23" y2="12"/>
      <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
      <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
    </svg>
    <svg v-else class="theme-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
    </svg>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useTheme } from '../composables/useTheme'

const { theme, isDark, saveTheme } = useTheme()

const themeTitle = computed(() => {
  switch (theme.value) {
    case 'system':
      return '当前：跟随系统'
    case 'light':
      return '当前：浅色模式'
    case 'dark':
      return '当前：深色模式'
    // @ts-expect-error 忽略类型检查以支持透明主题
    case 'light-transparent':
      return '当前：浅色透明模式'
    // @ts-expect-error 忽略类型检查以支持透明主题
    case 'dark-transparent':
      return '当前：深色透明模式'
    default:
      return '切换主题'
  }
})

const toggleTheme = () => {
  const themes: Array<'system' | 'light' | 'dark' | 'light-transparent' | 'dark-transparent'> = [
    'system', 
    'light', 
    'dark', 
    'light-transparent', 
    'dark-transparent'
  ]
  // 由于 saveTheme 只接受 ThemeMode 类型（不包括 transparent 主题），需要类型断言或类型兼容处理
  const currentIndex = themes.indexOf(theme.value as any)
  const nextIndex = (currentIndex + 1) % themes.length
  // @ts-expect-error 忽略类型检查以支持透明主题
  saveTheme(themes[nextIndex])
}
</script>

<style scoped>
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.3s ease;
  padding: 0;
}

.theme-toggle:hover {
  background: var(--bg-tertiary);
  border-color: var(--accent-color);
  box-shadow: 0 2px 8px var(--shadow-color);
}

.theme-icon {
  width: 20px;
  height: 20px;
  transition: transform 0.3s ease;
}

.theme-toggle:hover .theme-icon {
  transform: rotate(15deg);
}
</style> 