<template>
  <h2>应用设置</h2>
  <div class="app-settings-container">
    <div class="settings-section">
      <h3>外观设置</h3>
      <div class="setting-item">
        <label>主题模式</label>
        <select v-model="currentTheme" @change="handleThemeChange">
          <option value="system">跟随系统</option>
          <option value="light">浅色模式</option>
          <option value="dark">深色模式</option>
        </select>
      </div>
    </div>

    <div class="settings-section">
      <h3>自定义背景</h3>
      <div class="setting-item">
        <label>启用自定义背景</label>
        <input type="checkbox" v-model="backgroundSettings.enabled" @change="handleBackgroundChange">
      </div>
      
      <div v-if="backgroundSettings.enabled" class="background-settings">
        <div class="setting-item">
          <label>背景图片URL</label>
          <input 
            type="text" 
            v-model="backgroundSettings.url" 
            placeholder="https://example.com/image.jpg"
            @input="handleUrlInput"
          >
        </div>
        
        <div class="setting-item">
          <label>本地文件</label>
          <div class="file-input-container">
            <input 
              type="file" 
              ref="fileInput"
              accept="image/*,video/*"
              @change="handleFileSelect"
              class="file-input"
            >
            <button @click="triggerFileSelect" class="file-select-btn">
              <IconMdiFolderOpen style="margin-right: 5px;" />选择本地文件
            </button>
            <span v-if="backgroundSettings.url && backgroundSettings.url.startsWith('blob:')" class="file-info">
              已选择本地文件
            </span>
          </div>
        </div>
        
        <div class="setting-item">
          <label>预设背景</label>
          <div class="preset-backgrounds">
            <button 
              v-for="preset in backgroundPresets" 
              :key="preset.name"
              @click="selectPreset(preset)"
              class="preset-btn"
              :title="preset.name"
            >
              <div 
                class="preset-preview"
                :style="{ backgroundImage: `url(${preset.url})` }"
              ></div>
              <span class="preset-name">{{ preset.name }}</span>
            </button>
          </div>
        </div>
        
        <div class="setting-item">
          <label>背景模糊度</label>
          <div class="slider-container">
            <input 
              type="range" 
              v-model="backgroundSettings.blur" 
              min="0" 
              max="20" 
              step="1"
              @input="handleBackgroundChange"
            >
            <span class="slider-value">{{ backgroundSettings.blur }}px</span>
          </div>
        </div>
        
        <div class="setting-item">
          <label>背景透明度</label>
          <div class="slider-container">
            <input 
              type="range" 
              v-model="backgroundSettings.opacity" 
              min="0.1" 
              max="1" 
              step="0.1"
              @input="handleBackgroundChange"
            >
            <span class="slider-value">{{ Math.round(backgroundSettings.opacity * 100) }}%</span>
          </div>
        </div>
        
        <div class="setting-item">
          <label>背景位置</label>
          <select v-model="backgroundSettings.position" @change="handleBackgroundChange">
            <option value="center">居中</option>
            <option value="top">顶部</option>
            <option value="bottom">底部</option>
            <option value="left">左侧</option>
            <option value="right">右侧</option>
          </select>
        </div>
        
        <div class="setting-item">
          <label>背景尺寸</label>
          <select v-model="backgroundSettings.size" @change="handleBackgroundChange">
            <option value="cover">覆盖</option>
            <option value="contain">包含</option>
            <option value="auto">自动</option>
          </select>
        </div>
        
        <div class="background-preview" v-if="backgroundSettings.url">
          <h4>背景预览</h4>
          <div class="preview-container">
            <div 
              class="preview-image"
              :style="{
                backgroundImage: `url(${backgroundSettings.url})`,
                backgroundPosition: backgroundSettings.position,
                backgroundSize: backgroundSettings.size,
                filter: `blur(${backgroundSettings.blur}px)`,
                opacity: backgroundSettings.opacity
              }"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- <div class="settings-section">
      <h3>常规设置</h3>
      <div class="setting-item">
        <label>自动连接服务器</label>
        <input type="checkbox" v-model="autoConnect">
      </div>
      <div class="setting-item">
        <label>显示通知</label>
        <input type="checkbox" v-model="showNotifications">
      </div>
    </div>
    <div class="settings-section">
      <h3>网络设置</h3>
      <div class="setting-item">
        <label>默认WebSocket地址</label>
        <input type="text" v-model="defaultWsUrl" placeholder="ws://...">
      </div>
    </div>
    <button class="save-btn">保存设置</button> -->
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue';
import { useTheme, type ThemeMode } from '../composables/useTheme';

const { theme, backgroundSettings, saveTheme, saveBackgroundSettings } = useTheme();

const currentTheme = ref<ThemeMode>(theme.value);
const fileInput = ref<HTMLInputElement>();

// 背景预设
const backgroundPresets = ref([
  {
    name: '自然风景',
    url: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=1920&h=1080&fit=crop'
  },
  {
    name: '城市夜景',
    url: 'https://images.unsplash.com/photo-1519501025264-65ba15a82390?w=1920&h=1080&fit=crop'
  },
  {
    name: '抽象几何',
    url: 'https://images.unsplash.com/photo-1557683316-973673baf926?w=1920&h=1080&fit=crop'
  },
  {
    name: '科技感',
    url: 'https://images.unsplash.com/photo-1518709268805-4e9042af2176?w=1920&h=1080&fit=crop'
  },
  {
    name: '简约渐变',
    url: 'https://images.unsplash.com/photo-1557682250-33bd709cbe85?w=1920&h=1080&fit=crop'
  },
  {
    name: '星空',
    url: 'https://images.unsplash.com/photo-1534796636912-3b95b3ab5986?w=1920&h=1080&fit=crop'
  }
]);

const handleThemeChange = () => {
  saveTheme(currentTheme.value);
};

const handleBackgroundChange = () => {
  saveBackgroundSettings(backgroundSettings.value);
};

const selectPreset = (preset: { name: string; url: string }) => {
  // 清理之前的 blob URL
  if (backgroundSettings.value.url && backgroundSettings.value.url.startsWith('blob:')) {
    URL.revokeObjectURL(backgroundSettings.value.url);
  }
  
  backgroundSettings.value.url = preset.url;
  backgroundSettings.value.isVideo = false; // 预设背景都是图片
  saveBackgroundSettings(backgroundSettings.value);
};

const triggerFileSelect = () => {
  fileInput.value?.click();
};

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  
  if (file) {
    // 清理之前的 blob URL
    if (backgroundSettings.value.url && backgroundSettings.value.url.startsWith('blob:')) {
      URL.revokeObjectURL(backgroundSettings.value.url);
    }
    
    // 创建 blob URL
    const blobUrl = URL.createObjectURL(file);
    
    // 检测文件类型
    const isVideo = file.type.startsWith('video/');
    
    // 更新背景设置
    backgroundSettings.value.url = blobUrl;
    backgroundSettings.value.isVideo = isVideo;
    saveBackgroundSettings(backgroundSettings.value);
  }
};

const handleUrlInput = () => {
  // 手动输入的 URL 被标记为图片类型
  backgroundSettings.value.isVideo = false;
};

// 组件卸载时清理 blob URL
onUnmounted(() => {
  if (backgroundSettings.value.url && backgroundSettings.value.url.startsWith('blob:')) {
    URL.revokeObjectURL(backgroundSettings.value.url);
  }
});
</script>

<style scoped>
.app-settings-container {
  padding: 2rem 2rem 2rem 0;
}

.settings-section {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 2px 4px var(--shadow-color);
  background: var(--bg-secondary);
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border-light);
}

.setting-item label {
  color: var(--text-primary);
  font-weight: 500;
  min-width: 120px;
}

input[type="text"],
select {
  padding: 0.5rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  width: 250px;
  background: var(--bg-primary);
  color: var(--text-primary);
}

input[type="text"]:focus,
select:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: 0 0 0 2px var(--accent-light);
}

input[type="checkbox"] {
  width: auto;
  margin: 0;
}

.file-input-container {
  display: flex;
  align-items: center;
  gap: 1rem;
  width: 100%;
  max-width: 500px;
  flex-direction: row-reverse;
}

.file-input {
  display: none;
}

.file-select-btn {
  padding: 0.5rem 1rem;
  border: 2px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-primary);
  color: var(--text-primary);
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
}

.file-select-btn:hover {
  border-color: var(--accent-color);
  background: var(--accent-light);
}

.file-info {
  color: var(--success-color);
  font-size: 0.9rem;
  font-weight: 500;
}

.slider-container {
  display: flex;
  align-items: center;
  gap: 1rem;
  width: 250px;
}

input[type="range"] {
  flex: 1;
  height: 6px;
  border-radius: 3px;
  background: var(--bg-tertiary);
  outline: none;
  -webkit-appearance: none;
}

input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--accent-color);
  cursor: pointer;
  box-shadow: 0 2px 4px var(--shadow-color);
}

input[type="range"]::-moz-range-thumb {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--accent-color);
  cursor: pointer;
  border: none;
  box-shadow: 0 2px 4px var(--shadow-color);
}

.slider-value {
  min-width: 40px;
  text-align: right;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.background-settings {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-light);
}

.background-preview {
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-light);
}

.background-preview h4 {
  color: var(--text-primary);
  margin-bottom: 1rem;
  font-size: 1rem;
}

.preview-container {
  width: 100%;
  height: 120px;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid var(--border-color);
  position: relative;
}

.preview-image {
  width: 100%;
  height: 100%;
  background-repeat: no-repeat;
  background-attachment: fixed;
}

.preset-backgrounds {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 0.8rem;
  width: 100%;
  max-width: 500px;
}

.preset-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 0.8rem;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  cursor: pointer;
  transition: all 0.3s ease;
  min-height: 80px;
}

.preset-btn:hover {
  border-color: var(--accent-color);
  box-shadow: 0 2px 8px var(--accent-light);
}

.preset-preview {
  width: 60px;
  height: 40px;
  border-radius: 4px;
  background-size: cover;
  background-position: center;
  border: 1px solid var(--border-light);
}

.preset-name {
  font-size: 0.8rem;
  color: var(--text-secondary);
  text-align: center;
  font-weight: 500;
}

.save-btn {
  background: var(--accent-gradient);
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.6rem 1.2rem;
  cursor: pointer;
  margin-top: 1rem;
  font-weight: 600;
  transition: all 0.3s ease;
}

.save-btn:hover {
  box-shadow: 0 4px 12px var(--accent-light);
}
</style>