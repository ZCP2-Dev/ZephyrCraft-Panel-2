<template>
  <div class="other-container">
    <h2>其他设置</h2>
    <div class="system-info">
      <h3>系统信息</h3>
      <div class="info-grid">
        <div class="info-item">
          <div class="info-icon">V</div>
          <div class="info-content">
            <div class="info-label">服务器版本</div>
            <div class="info-value">1.20.1</div>
          </div>
        </div>
        <div class="info-item">
          <div class="info-icon">T</div>
          <div class="info-content">
            <div class="info-label">运行时间</div>
            <div class="info-value">3天 12小时 45分钟</div>
          </div>
        </div>
        <div class="info-item">
          <div class="info-icon">M</div>
          <div class="info-content">
            <div class="info-label">最大内存</div>
            <div class="info-value">4 GB</div>
          </div>
        </div>
        <div class="info-item">
          <div class="info-icon">U</div>
          <div class="info-content">
            <div class="info-label">已用内存</div>
            <div class="info-value">1.2 GB</div>
          </div>
        </div>
      </div>
    </div>
    
    <div class="quick-actions">
      <h3>快速操作</h3>
      <div class="action-grid">
        <button @click="saveWorld" class="action-btn save-btn">
          <IconMdiContentSave style="margin-right: 5px;" />保存世界
        </button>
        <button @click="clearDrops" class="action-btn clear-btn">
          <IconMdiDelete style="margin-right: 5px;" />清理掉落物
        </button>
        <button @click="setWeather" class="action-btn weather-btn">
          <IconMdiWeatherSunny style="margin-right: 5px;" />天气控制
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
function saveWorld() {
  console.log('保存世界');
}

function clearDrops() {
  console.log('清理掉落物');
}

function setWeather() {
  console.log('设置天气');
}
</script>

<style scoped>
.other-container {
  /* background: var(--bg-primary); */
  padding: 1.5rem;
}

.other-container h2 {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.other-container h3 {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 1.4rem;
  margin: 0 0 1.2rem 0;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.system-info {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  border: 1px solid var(--border-color);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.info-item {
  background: var(--bg-primary);
  padding: 1rem 1.2rem;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.info-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: var(--accent-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 1rem;
}

.info-content {
  display: flex;
  flex-direction: column;
}

.info-label {
  font-size: 0.9rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.info-value {
  font-size: 1.1rem;
  color: var(--text-primary);
  font-weight: 600;
}

.quick-actions {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 1.5rem;
  border: 1px solid var(--border-color);
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1rem;
}

.action-btn {
  background: var(--accent-gradient);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 1rem 1.2rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px var(--accent-light);
  position: relative;
  overflow: hidden;
}

.action-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.action-btn:hover::before {
  left: 100%;
}

.action-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--accent-light);
}

.action-btn:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.save-btn {
  background: var(--accent-gradient);
  box-shadow: 0 2px 8px var(--accent-light);
}

.save-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--accent-light);
}

.clear-btn {
  background: var(--error-color);
  box-shadow: 0 2px 8px var(--error-color);
}

.clear-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--error-color);
}

.weather-btn {
  background: var(--info-color);
  box-shadow: 0 2px 8px var(--info-color);
}

.weather-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--info-color);
}
</style>