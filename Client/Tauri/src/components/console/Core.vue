<template>
  <div class="core-container">
    <h2>核心设置</h2>
    <div class="core-section">
      <h3>服务器状态</h3>
      <div class="core-row">
        <div class="core-label">服务器状态：</div>
        <div class="core-value status-running">运行中</div>
      </div>
      <div class="core-row">
        <div class="core-label">最大玩家数：</div>
        <div class="core-value">20</div>
      </div>
      <div class="core-row">
        <div class="core-label">白名单：</div>
        <div class="core-value">已启用</div>
      </div>
      <div class="core-row">
        <div class="core-label">MOTD：</div>
        <div class="core-value">欢迎来到 ZephyrCraft 服务器！</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 静态假数据
</script>

<style scoped>
.core-container {
  background: #ffffff;
  padding: 1.5rem;
}

.core-container h2 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.core-section {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  border: 1px solid #e9ecef;
}

.core-section h3 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.4rem;
  margin: 0 0 1.2rem 0;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.core-row {
  display: flex;
  align-items: center;
  margin-bottom: 1.2rem;
  padding: 0.8rem 1rem;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.core-row:last-child {
  margin-bottom: 0;
}

.core-label {
  min-width: 120px;
  font-weight: 600;
  color: #7f8c8d;
  font-size: 0.95rem;
}

.core-value {
  flex: 1;
  color: #2c3e50;
  font-weight: 500;
}

.status-running {
  color: #27ae60;
  font-weight: 600;
}

.status-stopped {
  color: #e74c3c;
  font-weight: 600;
}

.core-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.7rem 1.2rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
  margin-left: 1rem;
}

.core-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

.core-btn:disabled {
  background: linear-gradient(135deg, #bdc3c7 0%, #95a5a6 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
</style>