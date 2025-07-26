<template>
  <div class="core-container">
    <h2>核心设置</h2>
    <div class="core-section">
      <h3>服务器状态</h3>
      <div class="core-row">
        <div class="core-label">服务器状态：</div>
        <div class="core-value" :class="isRunning ? 'status-running' : 'status-stopped'">
          {{ isRunning ? '运行中' : '已停止' }}
        </div>
        <button @click="toggleServer" class="core-btn" :disabled="!isConnected">
          {{ isRunning ? '停止服务器' : '启动服务器' }}
        </button>
      </div>
      <div class="core-row">
        <div class="core-label">在线玩家：</div>
        <div class="core-value">{{ serverInfo.playerCount }} / {{ serverInfo.maxPlayers }}</div>
      </div>
      <div class="core-row">
        <div class="core-label">服务器TPS：</div>
        <div class="core-value">{{ serverInfo.tps.toFixed(1) }}</div>
      </div>
      <div class="core-row">
        <div class="core-label">服务器版本：</div>
        <div class="core-value">{{ serverInfo.version }}</div>
      </div>
      <div class="core-row">
        <div class="core-label">启动时间：</div>
        <div class="core-value">{{ serverInfo.startTime || '----' }}</div>
      </div>
    </div>

    <div class="core-section">
      <h3>系统状态</h3>
      <div class="core-row">
        <div class="core-label">CPU使用率：</div>
        <div class="core-value">{{ systemInfo.cpuUsage.toFixed(1) }}%</div>
      </div>
      <div class="core-row">
        <div class="core-label">内存使用：</div>
        <div class="core-value">{{ formatMemory(systemInfo.memoryUsed, systemInfo.memoryTotal) }}</div>
      </div>
      <div class="core-row">
        <div class="core-label">磁盘使用：</div>
        <div class="core-value">{{ formatDisk(systemInfo.diskUsed, systemInfo.diskTotal) }}</div>
      </div>
      <div class="core-row">
        <div class="core-label">系统运行时间：</div>
        <div class="core-value">{{ formatUptime(systemInfo.uptime) }}</div>
      </div>
    </div>

    <div class="core-section">
      <h3>快速操作</h3>
      <div class="button-group">
        <button @click="sendCommand('save')" class="core-btn" :disabled="!isRunning">
          保存世界
        </button>
        <button @click="sendCommand('stop')" class="core-btn" :disabled="!isRunning">
          安全停止
        </button>
        <button @click="sendCommand('reload')" class="core-btn" :disabled="!isRunning">
          重载配置
        </button>
        <button @click="sendCommand('list')" class="core-btn" :disabled="!isRunning">
          查看玩家
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, watch, onUnmounted } from 'vue';

const props = defineProps<{ server?: any }>();
const wsApi = inject('wsApi') as any;
const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);

// 加载状态
const isLoading = ref(true);

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

// 服务器运行状态
const isRunning = ref(false);

let refreshTimer: number | null = null;

// 保存监听器引用以便清理
let systemMessageListener: ((data: any) => void) | null = null;
let terminalMessageListener: ((data: any) => void) | null = null;

// 格式化函数
function formatMemory(used: number, total: number): string {
  if (total === 0) return '0 MB / 0 MB';
  const usedGB = (used / 1024).toFixed(1);
  const totalGB = (total / 1024).toFixed(1);
  return `${usedGB} GB / ${totalGB} GB`;
}

function formatDisk(used: number, total: number): string {
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
    console.log('Core: Detected fake data, using last valid data:', newData);
    // 使用上次的真实数据
    systemInfo.value = { ...lastValidSystemInfo.value };
  } else {
    console.log('Core: Using real data:', newData);
    // 保存真实数据
    lastValidSystemInfo.value = { ...newData };
    systemInfo.value = { ...newData };
  }
}

// 请求系统信息
function requestSystemInfo() {
  if (wsApi && isConnected.value) {
    wsApi.send({ command: 'getSystemInfo' });
  }
}

// 请求服务器信息
function requestServerInfo() {
  if (wsApi && isConnected.value) {
    wsApi.send({ command: 'getServerInfo' });
  }
}

// 请求服务器状态
function requestStatus() {
  if (wsApi && isConnected.value) {
    wsApi.send({ command: 'status' });
  }
}

// 切换服务器状态
function toggleServer() {
  if (!isConnected.value) return;
  
  if (isRunning.value) {
    wsApi.send({ command: 'stop' });
  } else {
    wsApi.send({ command: 'start' });
  }
}

// 发送命令
function sendCommand(command: string) {
  if (wsApi && isConnected.value && isRunning.value) {
    wsApi.send({ command: 'input', content: command });
  }
}

// 开始定时刷新
function startRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer);
  }
  refreshTimer = setInterval(() => {
    requestSystemInfo();
    requestServerInfo();
    requestStatus();
  }, 5000); // 每5秒刷新一次
}

// 停止定时刷新
function stopRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer);
    refreshTimer = null;
  }
}

onMounted(() => {
  // 监听全局系统监控消息总线
  const systemBus = (window as any).__SYSTEM_BUS__;
  const terminalBus = (window as any).__TERMINAL_BUS__;
  
  if (systemBus && typeof systemBus.on === 'function') {
    systemMessageListener = (data: any) => {
      if (data && data.systemInfo) {
        validateAndUpdateSystemInfo(data.systemInfo);
        isLoading.value = false;
      }
      if (data && data.serverInfo) {
        serverInfo.value = data.serverInfo;
        isLoading.value = false;
      }
    };
    systemBus.on('system-message', systemMessageListener);
  }
  
  // 监听终端消息总线获取状态信息
  if (terminalBus && typeof terminalBus.on === 'function') {
    terminalMessageListener = (data: any) => {
      if (data && data.status) {
        isRunning.value = data.status === 'running';
      }
    };
    terminalBus.on('terminal-message', terminalMessageListener);
  }

  // 立即请求一次信息
  requestSystemInfo();
  requestServerInfo();
  requestStatus();
  
  // 开始定时刷新
  startRefresh();
});

onUnmounted(() => {
  console.log('Core: Component unmounting, cleaning up listeners...');
  
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
  
  console.log('Core: Listeners cleaned up');
});

watch(() => props.server, () => {
  // 重新请求信息
  requestSystemInfo();
  requestServerInfo();
  requestStatus();
});

watch(isConnected, (connected) => {
  if (connected) {
    requestSystemInfo();
    requestServerInfo();
    requestStatus();
    startRefresh();
  } else {
    stopRefresh();
  }
});
</script>

<style scoped>
.core-container {
  /* background: var(--bg-primary); */
  padding: 1.5rem;
}

.core-container h2 {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.core-section {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  border: 1px solid var(--border-color);
}

.core-section h3 {
  color: var(--text-primary);
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
  background: var(--bg-primary);
  border-radius: 6px;
  border: 1px solid var(--border-color);
}

.core-row:last-child {
  margin-bottom: 0;
}

.core-label {
  min-width: 120px;
  font-weight: 600;
  color: var(--text-secondary);
  font-size: 0.95rem;
}

.core-value {
  flex: 1;
  color: var(--text-primary);
  font-weight: 500;
}

.status-running {
  color: var(--success-color);
  font-weight: 600;
}

.status-stopped {
  color: var(--error-color);
  font-weight: 600;
}

.button-group {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.core-btn {
  background: var(--accent-gradient);
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
  box-shadow: 0 2px 8px var(--accent-light);
}

.core-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--accent-light);
}

.core-btn:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
  box-shadow: 0 2px 8px var(--shadow-color);
}
</style>