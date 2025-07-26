import { ref, watch, onMounted } from 'vue'

export type ThemeMode = 'system' | 'light' | 'dark'

export interface BackgroundSettings {
  enabled: boolean
  url: string
  blur: number
  opacity: number
  position: 'center' | 'top' | 'bottom' | 'left' | 'right'
  size: 'cover' | 'contain' | 'auto'
  isVideo: boolean
}

export function useTheme() {
  const theme = ref<ThemeMode>('system')
  const isDark = ref(false)
  const isTransparent = ref(false)
  
  // 背景设置
  const backgroundSettings = ref<BackgroundSettings>({
    enabled: false,
    url: '',
    blur: 0,
    opacity: 0.3,
    position: 'center',
    size: 'cover',
    isVideo: false
  })

  // 记录当前真实主题（不含透明）
  const realTheme = ref<ThemeMode>('system')

  // 从localStorage加载主题设置
  const loadTheme = () => {
    const savedTheme = localStorage.getItem('theme') as ThemeMode
    if (savedTheme && ['system', 'light', 'dark'].includes(savedTheme)) {
      theme.value = savedTheme
      realTheme.value = savedTheme
    }
    // 加载背景设置
    const savedBackground = localStorage.getItem('backgroundSettings')
    if (savedBackground) {
      try {
        backgroundSettings.value = { ...backgroundSettings.value, ...JSON.parse(savedBackground) }
      } catch (e) {
        console.warn('Failed to parse background settings:', e)
      }
    }
  }

  // 保存主题设置到localStorage
  const saveTheme = (newTheme: ThemeMode) => {
    localStorage.setItem('theme', newTheme)
    theme.value = newTheme
    realTheme.value = newTheme
    applyTheme()
  }

  // 保存背景设置到localStorage
  const saveBackgroundSettings = (settings: Partial<BackgroundSettings>) => {
    const newSettings = { ...backgroundSettings.value, ...settings }
    backgroundSettings.value = newSettings
    localStorage.setItem('backgroundSettings', JSON.stringify(newSettings))
    applyTheme()
    applyBackground()
  }

  // 检测系统主题
  const getSystemTheme = (): boolean => {
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }

  // 应用主题
  const applyTheme = () => {
    let shouldBeDark = false
    let shouldBeTransparent = false
    let themeClass = ''
    // 只允许三种主题
    switch (realTheme.value) {
      case 'system':
        shouldBeDark = getSystemTheme()
        break
      case 'dark':
        shouldBeDark = true
        break
      case 'light':
        shouldBeDark = false
        break
    }
    // 启用自定义背景时自动进入透明模式
    if (backgroundSettings.value.enabled && backgroundSettings.value.url) {
      shouldBeTransparent = true
      themeClass = shouldBeDark ? 'dark-transparent' : 'light-transparent'
    } else {
      shouldBeTransparent = false
      themeClass = shouldBeDark ? 'dark' : 'light'
    }
    isDark.value = shouldBeDark
    isTransparent.value = shouldBeTransparent
    document.documentElement.setAttribute('data-theme', themeClass)
    applyBackground()
  }

  // 应用背景
  const applyBackground = () => {
    const root = document.documentElement
    const body = document.body
    
    // 移除之前的背景容器
    const oldBgContainer = document.getElementById('custom-background')
    if (oldBgContainer) {
      oldBgContainer.remove()
    }
    
    if (backgroundSettings.value.enabled && backgroundSettings.value.url) {
      // 创建背景容器
      const bgContainer = document.createElement('div')
      bgContainer.id = 'custom-background'
      bgContainer.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: -1;
        filter: blur(${backgroundSettings.value.blur}px);
        opacity: ${backgroundSettings.value.opacity};
        pointer-events: none;
        overflow: hidden;
      `
      
      if (backgroundSettings.value.isVideo) {
        // 视频背景
        const video = document.createElement('video')
        video.src = backgroundSettings.value.url
        video.style.cssText = `
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          object-fit: ${backgroundSettings.value.size};
          object-position: ${backgroundSettings.value.position};
        `
        video.autoplay = true
        video.muted = true
        video.loop = true
        video.playsInline = true
        
        bgContainer.appendChild(video)
      } else {
        // 图片背景
        bgContainer.style.background = `
          url('${backgroundSettings.value.url}') 
          ${backgroundSettings.value.position} / ${backgroundSettings.value.size} 
          no-repeat fixed
        `
      }
      
      document.body.appendChild(bgContainer)
    }
    
    // 透明模式下设置前景色
    if (isTransparent.value) {
      body.style.background = 'transparent'
      if (isDark.value) {
        // 深色透明
        root.style.setProperty('--bg-primary', 'rgba(26, 26, 26, 0.65)')
        root.style.setProperty('--bg-secondary', 'rgba(45, 45, 45, 0.45)')
        root.style.setProperty('--bg-tertiary', 'rgba(61, 61, 61, 0.3)')
      } else {
        // 浅色透明
        root.style.setProperty('--bg-primary', 'rgba(255, 255, 255, 0.65)')
        root.style.setProperty('--bg-secondary', 'rgba(248, 249, 250, 0.45)')
        root.style.setProperty('--bg-tertiary', 'rgba(240, 240, 240, 0.3)')
      }
    } else {
      body.style.background = ''
      root.style.removeProperty('--bg-primary')
      root.style.removeProperty('--bg-secondary')
      root.style.removeProperty('--bg-tertiary')
    }
  }

  // 监听系统主题变化
  const setupSystemThemeListener = () => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
    mediaQuery.addEventListener('change', () => {
      if (realTheme.value === 'system') {
        applyTheme()
      }
    })
  }

  // 监听主题变化
  watch(theme, () => {
    realTheme.value = theme.value
    applyTheme()
  })

  // 监听背景设置变化
  watch(backgroundSettings, () => {
    applyTheme()
  }, { deep: true })

  onMounted(() => {
    loadTheme()
    applyTheme()
    setupSystemThemeListener()
  })

  return {
    theme,
    isDark,
    isTransparent,
    backgroundSettings,
    saveTheme,
    saveBackgroundSettings
  }
} 