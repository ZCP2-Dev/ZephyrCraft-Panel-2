<template>
  <div class="core-container">
    <h3>核心设置</h3>
    <div class="core-row">
      <div class="core-label">服务器状态：</div>
      <div class="core-value status-running">运行中</div>
      <button class="core-btn" :disabled="!isConnected" @click="restartServer">重启</button>
      <button class="core-btn" :disabled="!isConnected" @click="stopServer">关闭</button>
        </div>
    <div class="core-row">
      <div class="core-label">白名单：</div>
      <div class="core-value">已开启</div>
      <button class="core-btn">管理</button>
    </div>
    <div class="core-row">
      <div class="core-label">最大人数：</div>
      <div class="core-value">20</div>
      <button class="core-btn">修改</button>
    </div>
    <div class="core-row">
      <div class="core-label">服务器 MOTD：</div>
      <div class="core-value">欢迎来到 ZephyrCraft！</div>
      <button class="core-btn">修改</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject, computed } from 'vue';
const props = defineProps<{ server?: any }>();
const wsApi = inject('wsApi') as any;
const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);

function restartServer() {
  wsApi.send({ command: 'restart' });
}
function stopServer() {
  wsApi.send({ command: 'stop' });
}
</script>

<style scoped>
.core-container {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  padding: 2rem 2.5rem;
}
.core-row {
  display: flex;
  align-items: center;
  gap: 1.2rem;
  margin-bottom: 1.2rem;
}
.core-label {
  min-width: 100px;
  color: #888;
  font-size: 1rem;
}
.core-value {
  font-weight: bold;
  color: #23272e;
}
.status-running {
  color: #88bf64;
}
.core-btn {
  background: #88bf64;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 0.3rem 1.2rem;
  font-size: 0.98rem;
  cursor: pointer;
  transition: background 0.2s;
}
.core-btn:hover {
  background: #539951;
}
</style>