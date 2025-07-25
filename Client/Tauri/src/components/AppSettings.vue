<template>
  <h2>应用设置</h2>
  <div class="app-settings-container">
    <div class="settings-section">
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
    <button class="save-btn">保存设置</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const autoConnect = ref(true);
const showNotifications = ref(true);
const defaultWsUrl = ref('ws://localhost:8080');
</script>

<style scoped>
.app-settings-container {
  padding: 2rem 2rem 2rem 0;
}

.settings-section {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid #f0f0f0;
}

input[type="text"] {
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  width: 250px;
}

save-btn {
  background-color: #88bf64;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.6rem 1.2rem;
  cursor: pointer;
  margin-top: 1rem;
}

save-btn:hover {
  background-color: #539951;
}
</style>