<template>
  <div class="overview-container">
    <h2>服务器概览</h2>
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
      <p>版本：{{ version || '未知' }}</p>
      <p>启动时间：{{ startTime || '----' }}</p>
      <p>服务器地址：{{ server?.wsUrl || 'ws://127.0.0.1:19132' }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, watch } from 'vue';
const props = defineProps<{ server?: any }>();
const version = ref('');
const startTime = ref('');
const wsApi = inject('wsApi') as any;
const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);

onMounted(() => {
  // 监听全局 ws 消息，自动同步版本、启动时间
  const bus = (window as any).__TERMINAL_BUS__;
  if (bus && typeof bus.on === 'function') {
    bus.on('terminal-message', (data: any) => {
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
  version.value = '';
  startTime.value = '';
});
</script>

<style scoped>
.overview-container {
  background: #ffffff;
  padding: 1.5rem;
}

.overview-container h2 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.card-row {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
}

.info-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.2rem 1.5rem;
  min-width: 160px;
  text-align: center;
  border: 1px solid #e9ecef;
  flex: 1;
  min-width: 200px;
}

.card-title {
  color: #7f8c8d;
  font-size: 1rem;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.card-value {
  font-size: 1.6rem;
  font-weight: 700;
  color: #2c3e50;
}

.overview-section {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  border: 1px solid #e9ecef;
}

.overview-section h3 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.4rem;
  margin: 0 0 1.2rem 0;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.overview-section p {
  color: #34495e;
  font-size: 1rem;
  margin: 0.8rem 0;
  padding: 0.7rem 0.8rem;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.overview-section p:first-of-type {
  margin-top: 0;
}

.overview-section p:last-of-type {
  margin-bottom: 0;
}

.status-running {
  color: #27ae60;
  font-weight: 600;
}

.status-stopped {
  color: #e74c3c;
  font-weight: 600;
}
</style>