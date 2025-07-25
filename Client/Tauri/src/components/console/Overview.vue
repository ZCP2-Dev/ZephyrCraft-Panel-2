<template>
  <div class="overview-container">
    <div class="card-row">
      <div class="info-card">
        <div class="card-title">服务器名称</div>
        <div class="card-value">{{ server?.name || 'BDS服务器' }}</div>
      </div>
      <div class="info-card">
        <div class="card-title">在线玩家</div>
        <div class="card-value">12 / 20</div>
      </div>
      <div class="info-card">
        <div class="card-title">CPU 占用</div>
        <div class="card-value">23%</div>
      </div>
      <div class="info-card">
        <div class="card-title">内存占用</div>
        <div class="card-value">1.2 GB / 4 GB</div>
      </div>
    </div>
    <div class="overview-section">
      <h3>服务器状态</h3>
      <p>状态：<span :class="isRunning ? 'status-running' : 'status-stopped'">{{ isRunning ? '运行中' : '已停止' }}</span></p>
      <p>版本：{{ version || '未知' }}</p>
      <p>启动时间：{{ startTime || '----' }}</p>
      <p>服务器地址：{{ server?.wsUrl || 'ws://127.0.0.1:19132' }}</p>
      <div class="overview-actions">
        <button v-if="!isRunning" class="overview-btn" :disabled="!isConnected" @click="startServer">启动服务器</button>
        <button v-else class="overview-btn stop" :disabled="!isConnected" @click="stopServer">停止服务器</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, watch } from 'vue';
const props = defineProps<{ server?: any }>();
const isRunning = ref(false);
const version = ref('');
const startTime = ref('');
const wsApi = inject('wsApi') as any;
const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);

function startServer() {
  wsApi.send({ command: 'start' });
}
function stopServer() {
  wsApi.send({ command: 'stop' });
}

onMounted(() => {
  // 监听全局 ws 消息，自动同步 isRunning、版本、启动时间
  const bus = (window as any).__TERMINAL_BUS__;
  if (bus && typeof bus.on === 'function') {
    bus.on('terminal-message', (data: any) => {
      if (data && data.status) {
        isRunning.value = data.status === 'running';
      }
      if (data && data.version) {
        version.value = data.version;
      }
      if (data && data.startTime) {
        startTime.value = data.startTime;
      }
      // 断开时重置
      if (data && data.status === 'stopped') {
        startTime.value = '';
      }
    });
  }
});

watch(() => props.server, () => {
  isRunning.value = false;
  version.value = '';
  startTime.value = '';
});
</script>

<style scoped>
.overview-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}
.card-row {
  display: flex;
  gap: 2rem;
}
.info-card {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  padding: 1.5rem 2.5rem;
  min-width: 160px;
  text-align: center;
}
.card-title {
  color: #888;
  font-size: 1rem;
  margin-bottom: 0.5rem;
}
.card-value {
  font-size: 1.6rem;
  font-weight: bold;
  color: #23272e;
}
.status-running {
  color: #88bf64;
  font-weight: bold;
}
.status-stopped {
  color: #ff6b6b;
  font-weight: bold;
}
.overview-section {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  padding: 2rem 2.5rem;
}
.overview-actions {
  margin-top: 1.5rem;
}
.overview-btn {
  background: #88bf64;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 0.6rem 2.2rem;
  font-size: 1.1rem;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s;
}
.overview-btn:disabled {
  background: #b7d7a8;
  cursor: not-allowed;
}
.overview-btn.stop {
  background: #ff6b6b;
}
.overview-btn.stop:disabled {
  background: #ffb3b3;
}
</style>