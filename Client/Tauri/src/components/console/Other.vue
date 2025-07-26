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
  background: #ffffff;
  padding: 1.5rem;
}

.other-container h2 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.other-container h3 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.4rem;
  margin: 0 0 1.2rem 0;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.system-info {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  border: 1px solid #e9ecef;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.info-item {
  background: #ffffff;
  padding: 1rem 1.2rem;
  border-radius: 6px;
  border: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.info-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
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
  color: #7f8c8d;
  font-weight: 500;
}

.info-value {
  font-size: 1.1rem;
  color: #2c3e50;
  font-weight: 600;
}

.quick-actions {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.5rem;
  border: 1px solid #e9ecef;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1rem;
}

.action-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
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
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
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
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

.action-btn:disabled {
  background: linear-gradient(135deg, #bdc3c7 0%, #95a5a6 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.save-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
}

.save-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

.clear-btn {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
  box-shadow: 0 2px 8px rgba(231, 76, 60, 0.3);
}

.clear-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(231, 76, 60, 0.4);
}

.weather-btn {
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  box-shadow: 0 2px 8px rgba(52, 152, 219, 0.3);
}

.weather-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.4);
}
</style>