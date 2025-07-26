<template>
  <div class="overview-container">
    <h2>服务器概览</h2>
    
    <!-- 假数据警告 -->
    <div v-if="showFakeDataWarning" class="fake-data-warning">
      <div class="warning-icon">⚠️</div>
      <div class="warning-text">
        <strong>数据更新中</strong><br>
        正在更新数据，请稍后...
      </div>
    </div>
    <div class="card-row">
      <div class="info-card">
        <div class="card-title">服务器名称</div>
        <div class="card-value">{{ server?.name || 'BDS服务器' }}</div>
      </div>
      <div class="info-card">
        <div class="card-title">在线玩家</div>
        <div class="card-value">{{ serverInfo.playerCount }} / {{ serverInfo.maxPlayers }}</div>
      </div>
      <div class="info-card">
        <div class="card-title">CPU 占用</div>
        <div class="card-value">
          <div v-if="isLoading" class="loading-spinner"></div>
          <span v-else>{{ formatCPUUsage(systemInfo.cpuUsage) }}</span>
        </div>
      </div>
      <div class="info-card">
        <div class="card-title">内存占用</div>
        <div class="card-value">
          <div v-if="isLoading" class="loading-spinner"></div>
          <span v-else>{{ formatMemoryUsage(systemInfo.memoryUsed, systemInfo.memoryTotal) }}</span>
        </div>
      </div>
    </div>
    
    <div class="card-row">
      <div class="info-card">
        <div class="card-title">磁盘使用</div>
        <div class="card-value">
          <div v-if="isLoading" class="loading-spinner"></div>
          <span v-else>{{ formatDiskUsage(systemInfo.diskUsed, systemInfo.diskTotal) }}</span>
        </div>
      </div>
      <div class="info-card">
        <div class="card-title">服务器TPS</div>
        <div class="card-value">
          <div v-if="isLoading" class="loading-spinner"></div>
          <span v-else>{{ serverInfo.tps.toFixed(1) }}</span>
        </div>
      </div>
      <div class="info-card">
        <div class="card-title">系统运行时间</div>
        <div class="card-value">
          <div v-if="isLoading" class="loading-spinner"></div>
          <span v-else>{{ formatUptime(systemInfo.uptime) }}</span>
        </div>
      </div>
      <div class="info-card">
        <div class="card-title">连接状态</div>
        <div class="card-value" :class="isConnected ? 'status-running' : 'status-stopped'">
          {{ isConnected ? '已连接' : '未连接' }}
        </div>
      </div>
    </div>
    
    <div class="overview-section">
      <h3>服务器状态</h3>
      <p>版本：{{ serverInfo.version }}</p>
      <p>启动时间：{{ serverInfo.startTime || '----' }}</p>
      <p>服务器地址：{{ server?.wsUrl || 'ws://127.0.0.1:19132' }}</p>
    </div>
    
    <div class="overview-section">
      <h3>系统状态</h3>
      <div class="progress-container">
        <div class="progress-item">
          <div class="progress-label">CPU使用率</div>
          <div class="progress-bar">
            <div v-if="isLoading" class="progress-loading"></div>
            <div v-else class="progress-fill" :style="{ width: systemInfo.cpuUsage + '%' }"></div>
          </div>
          <div class="progress-text">
            <div v-if="isLoading" class="loading-spinner"></div>
            <span v-else>{{ systemInfo.cpuUsage.toFixed(1) }}%</span>
          </div>
        </div>
        <div class="progress-item">
          <div class="progress-label">内存使用率</div>
          <div class="progress-bar">
            <div v-if="isLoading" class="progress-loading"></div>
            <div v-else class="progress-fill" :style="{ width: systemInfo.memoryUsage + '%' }"></div>
          </div>
          <div class="progress-text">
            <div v-if="isLoading" class="loading-spinner"></div>
            <span v-else>{{ systemInfo.memoryUsage.toFixed(1) }}%</span>
          </div>
        </div>
        <div class="progress-item">
          <div class="progress-label">磁盘使用率</div>
          <div class="progress-bar">
            <div v-if="isLoading" class="progress-loading"></div>
            <div v-else class="progress-fill" :style="{ width: systemInfo.diskUsage + '%' }"></div>
          </div>
          <div class="progress-text">
            <div v-if="isLoading" class="loading-spinner"></div>
            <span v-else>{{ systemInfo.diskUsage.toFixed(1) }}%</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, watch, onUnmounted } from 'vue';

const props = defineProps<{ server?: any }>();
const version = ref('');
const startTime = ref('');
const wsApi = inject('wsApi') as any;
const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);

// 加载状态
const isLoading = ref(true);

// 假数据警告状态
const showFakeDataWarning = ref(false);

// 系统信息状态
const systemInfo = ref({
  cpuUsage: 0,
  memoryUsage: 0,
  memoryTotal: 0,
  memoryUsed: 0,
  diskUsage: 0,
  diskTotal: 0,
  diskUsed: 0,
  uptime: 0
});

// 上次获取的真实系统信息
const lastValidSystemInfo = ref({
  cpuUsage: 0,
  memoryUsage: 0,
  memoryTotal: 0,
  memoryUsed: 0,
  diskUsage: 0,
  diskTotal: 0,
  diskUsed: 0,
  uptime: 0
});

// 服务器信息状态
const serverInfo = ref({
  version: '未知',
  startTime: '',
  playerCount: 0,
  maxPlayers: 20,
  tps: 20.0,
  uptime: 0
});

let refreshTimer: number | null = null;

// 格式化函数
function formatCPUUsage(usage: number): string {
  return usage.toFixed(1) + '%';
}

function formatMemoryUsage(used: number, total: number): string {
  if (total === 0) return '0 MB / 0 MB';
  const usedGB = (used / 1024).toFixed(1);
  const totalGB = (total / 1024).toFixed(1);
  return `${usedGB} GB / ${totalGB} GB`;
}

function formatDiskUsage(used: number, total: number): string {
  if (total === 0) return '0 GB / 0 GB';
  const usedGB = (used / 1024).toFixed(1);
  const totalGB = (total / 1024).toFixed(1);
  return `${usedGB} GB / ${totalGB} GB`;
}

function formatUptime(uptime: number): string {
  if (uptime === 0) return '0天';
  const days = Math.floor(uptime / 86400);
  const hours = Math.floor((uptime % 86400) / 3600);
  return `${days}天${hours}小时`;
}

// 判断是否为假数据
function isFakeData(data: any): boolean {
  // 检查内存数据是否为假数据（1000MB = 1GB，明显不是真实系统）
  if (data.memoryTotal === 1000 || data.memoryTotal === 1024) {
    return true;
  }
  
  // 检查磁盘数据是否为假数据（1000MB = 1GB，明显不是真实系统）
  if (data.diskTotal === 1000 || data.diskTotal === 1024) {
    return true;
  }
  
  // 检查其他明显不合理的值
  if (data.memoryTotal < 2000 || data.diskTotal < 10000) { // 小于2GB内存或10GB磁盘
    return true;
  }
  
  return false;
}

// 验证并更新系统信息
function validateAndUpdateSystemInfo(newData: any) {
  if (isFakeData(newData)) {
    console.log('Overview: Detected fake data, using last valid data:', newData);
    // 使用上次的真实数据
    systemInfo.value = { ...lastValidSystemInfo.value };
    // 显示假数据警告
    showFakeDataWarning.value = true;
  } else {
    console.log('Overview: Using real data:', newData);
    // 保存真实数据
    lastValidSystemInfo.value = { ...newData };
    systemInfo.value = { ...newData };
    // 隐藏假数据警告
    showFakeDataWarning.value = false;
  }
}

// 请求系统信息
function requestSystemInfo() {
  console.log('Overview: Requesting system info, wsApi:', !!wsApi, 'isConnected:', isConnected.value);
  if (wsApi && isConnected.value) {
    console.log('Overview: Sending getSystemInfo command');
    wsApi.send({ command: 'getSystemInfo' });
  } else {
    console.log('Overview: Cannot send getSystemInfo - wsApi:', !!wsApi, 'isConnected:', isConnected.value);
  }
}

// 请求服务器信息
function requestServerInfo() {
  console.log('Overview: Requesting server info, wsApi:', !!wsApi, 'isConnected:', isConnected.value);
  if (wsApi && isConnected.value) {
    console.log('Overview: Sending getServerInfo command');
    wsApi.send({ command: 'getServerInfo' });
  } else {
    console.log('Overview: Cannot send getServerInfo - wsApi:', !!wsApi, 'isConnected:', isConnected.value);
  }
}

// 开始定时刷新
function startRefresh() {
  console.log('Overview: Starting refresh timer');
  if (refreshTimer) {
    clearInterval(refreshTimer);
  }
  refreshTimer = setInterval(() => {
    console.log('Overview: Refresh timer triggered');
    requestSystemInfo();
    requestServerInfo();
  }, 5000); // 每5秒刷新一次
}

// 停止定时刷新
function stopRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer);
    refreshTimer = null;
  }
}

// 保存监听器引用以便清理
let systemMessageListener: ((data: any) => void) | null = null;
let terminalMessageListener: ((data: any) => void) | null = null;

onMounted(() => {
  console.log('Overview: Component mounted');
  // 监听全局系统监控消息总线
  const systemBus = (window as any).__SYSTEM_BUS__;
  const terminalBus = (window as any).__TERMINAL_BUS__;
  
  if (systemBus && typeof systemBus.on === 'function') {
    console.log('Overview: Setting up system bus listener');
    systemMessageListener = (data: any) => {
      console.log('Overview: Received system message:', data);
      if (data && data.systemInfo) {
        console.log('Overview: Updating systemInfo:', data.systemInfo);
        validateAndUpdateSystemInfo(data.systemInfo);
        // 收到系统信息后关闭加载状态
        isLoading.value = false;
      }
      if (data && data.serverInfo) {
        console.log('Overview: Updating serverInfo:', data.serverInfo);
        serverInfo.value = data.serverInfo;
        // 收到服务器信息后关闭加载状态
        isLoading.value = false;
      }
    };
    systemBus.on('system-message', systemMessageListener);
  } else {
    console.error('Overview: SystemBus not found or invalid');
  }
  
  // 监听终端消息总线获取其他信息
  if (terminalBus && typeof terminalBus.on === 'function') {
    console.log('Overview: Setting up terminal bus listener');
    terminalMessageListener = (data: any) => {
      console.log('Overview: Received terminal message:', data);
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
    };
    terminalBus.on('terminal-message', terminalMessageListener);
  } else {
    console.error('Overview: TerminalBus not found or invalid');
  }

  // 立即请求一次信息
  console.log('Overview: Initial requests');
  requestSystemInfo();
  requestServerInfo();
  
  // 开始定时刷新
  startRefresh();
  
  // 测试消息总线是否工作（仅在开发环境）
  if (import.meta.env.DEV) {
    setTimeout(() => {
      console.log('Overview: Testing message bus');
      if (window && (window as any).__SYSTEM_BUS__) {
        (window as any).__SYSTEM_BUS__.emit('system-message', {
          systemInfo: { cpuUsage: 50, memoryUsage: 60, memoryTotal: 1000, memoryUsed: 600, diskUsage: 70, diskTotal: 1000, diskUsed: 700, uptime: 3600 },
          serverInfo: { version: '测试版本', startTime: '2024-01-01', playerCount: 5, maxPlayers: 20, tps: 18.5, uptime: 1800 }
        });
      }
    }, 2000);
  }
});

// 组件卸载时清理监听器
onUnmounted(() => {
  console.log('Overview: Component unmounting, cleaning up listeners...');
  
  // 停止定时刷新
  stopRefresh();
  
  // 清理系统消息监听器
  if (systemMessageListener) {
    const systemBus = (window as any).__SYSTEM_BUS__;
    if (systemBus && typeof systemBus.off === 'function') {
      systemBus.off('system-message', systemMessageListener);
    }
    systemMessageListener = null;
  }
  
  // 清理终端消息监听器
  if (terminalMessageListener) {
    const terminalBus = (window as any).__TERMINAL_BUS__;
    if (terminalBus && typeof terminalBus.off === 'function') {
      terminalBus.off('terminal-message', terminalMessageListener);
    }
    terminalMessageListener = null;
  }
  
  console.log('Overview: Listeners cleaned up');
});

watch(() => props.server, () => {
  version.value = '';
  startTime.value = '';
  // 重新设置加载状态
  isLoading.value = true;
  // 重新请求信息
  requestSystemInfo();
  requestServerInfo();
});

watch(isConnected, (connected) => {
  if (connected) {
    // 重新设置加载状态
    isLoading.value = true;
    requestSystemInfo();
    requestServerInfo();
    startRefresh();
  } else {
    stopRefresh();
  }
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

.progress-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.progress-item {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.progress-label {
  min-width: 100px;
  font-weight: 600;
  color: #2c3e50;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: #e9ecef;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #27ae60, #2ecc71);
  transition: width 0.3s ease;
}

.progress-text {
  min-width: 60px;
  text-align: right;
  font-weight: 600;
  color: #2c3e50;
}

/* 加载动画 */
.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid #e9ecef;
  border-top: 2px solid #27ae60;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

.progress-loading {
  height: 100%;
  background: linear-gradient(90deg, #e9ecef 25%, #f8f9fa 50%, #e9ecef 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
  border-radius: 4px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}

/* 假数据警告样式 */
.fake-data-warning {
  background: #fff3cd;
  border: 1px solid #ffeaa7;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  animation: fadeIn 0.3s ease-in;
}

.warning-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
}

.warning-text {
  color: #856404;
  font-size: 0.9rem;
  line-height: 1.4;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>