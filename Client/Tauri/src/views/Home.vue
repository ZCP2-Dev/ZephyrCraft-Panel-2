<template>
  <div class="layout-root">
    <div class="sidebar">
      <div class="sidebar-header">
        <img src="../assets/logo.png" alt="logo" style="width: 60px; height: 60px;" />
        <h2 style="margin-top: 5px;">ZephyrCraft Panel</h2>
        <div class="header-actions">
          <!-- <ThemeToggle /> 颜色模式 -->
        </div>
      </div>
      <div class="sidebar-menu">
  <template v-if="!isInConsole">
          <button @click="showServerList" class="menu-item" :class="{ 'active': currentView === 'serverList' }">
            <IconMdiHome style="margin-right: 10px; font-size: 1.2em;" />服务器列表
          </button>
          <button @click="showAppSettings" class="menu-item" :class="{ 'active': currentView === 'appSettings' }">
            <IconMdiCog style="margin-right: 10px; font-size: 1.2em;" />应用设置
          </button>
  </template>
  <template v-else>
          <button @click="showServerList" class="menu-item return-btn">
            <IconMdiHome style="margin-right: 10px; font-size: 1.2em;" />返回服务器列表
          </button>
    <div class="console-menu-header">控制台</div>
          <div class="console-menu-items">
          <button v-for="item in consoleMenuItems" :key="item.key" @click="changeConsoleSection(item.key)" class="menu-item console-item" :class="{ 'active': consoleSection === item.key }">
              <IconMdiViewDashboard v-if="item.key === 'overview'" style="margin-right: 10px; font-size: 1.2em;" />
              <IconMdiConsole v-else-if="item.key === 'terminal'" style="margin-right: 10px; font-size: 1.2em;" />
              <IconMdiAccountGroup v-else-if="item.key === 'players'" style="margin-right: 10px; font-size: 1.2em;" />
              <IconMdiChip v-else-if="item.key === 'core'" style="margin-right: 10px; font-size: 1.2em;" />
              <IconMdiPuzzle v-else-if="item.key === 'plugins'" style="margin-right: 10px; font-size: 1.2em;" />
              <IconMdiFileDocument v-else-if="item.key === 'files'" style="margin-right: 10px; font-size: 1.2em;" />
              <IconMdiCog v-else-if="item.key === 'remote'" style="margin-right: 10px; font-size: 1.2em;" />
              <IconMdiDotsHorizontal v-else-if="item.key === 'other'" style="margin-right: 10px; font-size: 1.2em;" />
              {{ item.label }}
            </button>
          </div>
          
          <!-- 连接状态和服务器控制 -->
          <div class="connection-status">
            <div class="status-indicator" :class="{ 
              'connected': connectionStatus === 'connected', 
              'connecting': connectionStatus === 'connecting',
              'failed': connectionStatus === 'failed',
              'disconnected': connectionStatus === 'disconnected' 
            }">
              <IconMdiCircle v-if="connectionStatus === 'connected'" style="color: #27ae60; margin-right: 5px;" />
              <IconMdiLoading v-else-if="connectionStatus === 'connecting'" class="spin" style="color: #f39c12; margin-right: 5px;" />
              <IconMdiCircle v-else-if="connectionStatus === 'failed'" style="color: #e74c3c; margin-right: 5px;" />
              <IconMdiCircle v-else style="color: #95a5a6; margin-right: 5px;" />
              {{ getStatusText() }}
            </div>
            
            <!-- 连接失败或断开时显示重连控件 -->
            <div v-if="(connectionStatus === 'failed' || connectionStatus === 'disconnected')" class="reconnect-controls">
              <div class="error-message">{{ lastError || '连接已断开' }}</div>
              <div class="reconnect-buttons">
                <button @click="manualReconnect" class="reconnect-btn" :disabled="isConnecting">
                  <IconMdiRefresh v-if="!isConnecting" style="margin-right: 5px;" />
                  <IconMdiLoading v-else class="spin" style="margin-right: 5px;" />
                  {{ isConnecting ? '重连中...' : '重连' }}
                </button>
                <!-- <button @click="manualConnect" class="connect-btn" :disabled="isConnecting || !currentServer">
                  <IconMdiConnection v-if="!isConnecting" style="margin-right: 5px;" />
                  <IconMdiLoading v-else class="spin" style="margin-right: 5px;" />
                  {{ isConnecting ? '连接中...' : '连接' }}
                </button> -->
                <button @click="resetConnection" class="reset-btn">
                  <IconMdiClose style="margin-right: 5px;" />重置
                </button>
              </div>
              <!-- 错误四次后显示重连链接 -->
              <div v-if="reconnectAttempts >= 4" class="reconnect-link">
                <a href="#" @click.prevent="manualReconnect" class="reconnect-a">
                  <IconMdiRefresh v-if="!isConnecting" style="margin-right: 5px;" />
                  <IconMdiLoading v-else class="spin" style="margin-right: 5px;" />
                  {{ isConnecting ? '重连中...' : '点击重新连接' }}
                </a>
              </div>
            </div>
            
            <div class="server-controls">
              <button v-if="!isRunning" class="control-btn start" :disabled="connectionStatus !== 'connected'" @click="startServer">
                <IconMdiPlay style="margin-right: 5px;" />启动
              </button>
              <button v-else class="control-btn stop" :disabled="connectionStatus !== 'connected'" @click="stopServer">
                <IconMdiStop style="margin-right: 5px;" />停止
              </button>
              <button class="control-btn restart" :disabled="connectionStatus !== 'connected' || !isRunning" @click="restartServer">
                <IconMdiRestart style="margin-right: 5px;" />重启
              </button>
            </div>
            
            <!-- 调试信息 -->
            <div v-if="isInConsole" style="font-size: 12px; color: #666; margin-top: 10px;display: none;">
              调试: 连接状态={{ connectionStatus }}, 运行状态={{ isRunning }}, 按钮可用={{ connectionStatus === 'connected' }}
              <button @click="refreshServerStatus" style="margin-left: 10px; padding: 2px 8px; font-size: 10px; background: #f0f0f0; border: 1px solid #ccc; border-radius: 4px; cursor: pointer;">
                刷新状态
              </button>
            </div>
          </div>
  </template>
</div>
    </div>
    <div class="main-content">
      <component :is="currentComponent" v-bind="consoleProps" :key="`${currentView}-${consoleSection}`" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, defineAsyncComponent, watch, computed, provide, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useWebSocket } from '../useWebSocket';

const ServerList = defineAsyncComponent(() => import('../components/ServerList.vue'));
const AppSettings = defineAsyncComponent(() => import('../components/AppSettings.vue'));
const Overview = defineAsyncComponent(() => import('../components/console/Overview.vue'));
const Terminal = defineAsyncComponent(() => import('../components/console/Terminal.vue'));
const Players = defineAsyncComponent(() => import('../components/console/Players.vue'));
const Core = defineAsyncComponent(() => import('../components/console/Core.vue'));
const Plugins = defineAsyncComponent(() => import('../components/console/Plugins.vue'));
const Files = defineAsyncComponent(() => import('../components/console/Files.vue'));
const Other = defineAsyncComponent(() => import('../components/console/Other.vue'));
const Remote = defineAsyncComponent(() => import('../components/console/Remote.vue'));

const route = useRoute();
const router = useRouter();
const currentView = ref('serverList');
const currentComponent = ref(ServerList);
const consoleSection = ref('overview');
const isInConsole = ref(false);
const currentServer = ref<any>(null);

const wsUrl = ref('');
const wsPassword = ref('');
const isRunning = ref(false);
const reconnectAttempts = ref(0);
const maxReconnectAttempts = 3;
const reconnectDelay = 2000; // 2秒

const wsApi = useWebSocket({
  get url() { return wsUrl.value; },
  get password() { return wsPassword.value; },
});
provide('wsApi', wsApi);
provide('isRunning', isRunning);
provide('isInConsole', isInConsole);

// 连接状态相关的计算属性和方法
const connectionStatus = computed(() => wsApi.connectionStatus.value);
const isConnecting = computed(() => wsApi.isConnecting.value);
const lastError = computed(() => wsApi.lastError.value);

function getStatusText() {
  switch (connectionStatus.value) {
    case 'connected':
      return '已连接';
    case 'connecting':
      return '正在连接...';
    case 'failed':
      return '连接失败';
    case 'disconnected':
    default:
      return '未连接';
  }
}

function manualReconnect() {
  console.log('Manual reconnect triggered');
  reconnectAttempts.value = 0;
  wsApi.connect();
}

function resetConnection() {
  reconnectAttempts.value = 0;
  wsApi.disconnect();
  // 重置错误状态
  setTimeout(() => {
    wsApi.connect();
  }, 1000);
}

const consoleMenuItems = [
  { key: 'overview', label: '概览' },
  { key: 'terminal', label: '终端' },
  { key: 'players', label: '玩家' },
  { key: 'core', label: '核心' },
  { key: 'plugins', label: '插件' },
  { key: 'files', label: '文件' },
  { key: 'remote', label: '远程' },
  { key: 'other', label: '其它' }
];

// 服务器控制方法
function startServer() {
  console.log('Sending start command'); // 调试日志
  if (connectionStatus.value === 'connected') {
    wsApi.send({ command: 'start' });
    // 延迟查询状态
    setTimeout(() => {
      if (connectionStatus.value === 'connected') {
        console.log('Querying status after start command'); // 调试日志
        wsApi.send({ command: 'status' });
      }
    }, 2000);
  } else {
    console.error('Cannot start server: not connected');
  }
}

function stopServer() {
  console.log('Sending stop command'); // 调试日志
  if (connectionStatus.value === 'connected') {
    wsApi.send({ command: 'stop' });
    // 延迟查询状态
    setTimeout(() => {
      if (connectionStatus.value === 'connected') {
        console.log('Querying status after stop command'); // 调试日志
        wsApi.send({ command: 'status' });
      }
    }, 2000);
  } else {
    console.error('Cannot stop server: not connected');
  }
}

function restartServer() {
  console.log('Sending restart command'); // 调试日志
  if (connectionStatus.value === 'connected') {
    wsApi.send({ command: 'restart' });
    // 延迟查询状态
    setTimeout(() => {
      if (connectionStatus.value === 'connected') {
        console.log('Querying status after restart command'); // 调试日志
        wsApi.send({ command: 'status' });
      }
    }, 3000);
  } else {
    console.error('Cannot restart server: not connected');
  }
}

// 手动刷新服务器状态
function refreshServerStatus() {
  console.log('Manually refreshing server status'); // 调试日志
  if (connectionStatus.value === 'connected') {
    wsApi.send({ command: 'status' });
  } else {
    console.error('Cannot refresh status: not connected');
  }
}

// 自动重连逻辑
function attemptReconnect() {
  console.log('Attempting reconnect, attempts:', reconnectAttempts.value, 'max:', maxReconnectAttempts, 'status:', connectionStatus.value, 'isEnteringConsole:', isEnteringConsole);
  
  if (reconnectAttempts.value < maxReconnectAttempts && 
      (connectionStatus.value === 'failed' || connectionStatus.value === 'disconnected') && 
      !isEnteringConsole) {
    reconnectAttempts.value++;
    console.log('Starting reconnect attempt', reconnectAttempts.value);
    setTimeout(() => {
      wsApi.connect();
    }, reconnectDelay);
  } else {
    console.log('Reconnect conditions not met, skipping reconnect');
  }
}

// 监听连接状态变化
const connectionStatusWatcher = watch(connectionStatus, (status, oldStatus) => {
  console.log('Connection status changed:', status, 'oldStatus:', oldStatus, 'isInConsole:', isInConsole.value, 'reconnectAttempts:', reconnectAttempts.value, 'isEnteringConsole:', isEnteringConsole);
  
  if ((status === 'failed' || status === 'disconnected') && isInConsole.value && reconnectAttempts.value < maxReconnectAttempts && !isEnteringConsole) {
    console.log('Attempting automatic reconnect...');
    attemptReconnect();
  } else if (status === 'connected') {
    reconnectAttempts.value = 0;
    console.log('Connection established successfully, will query status in 2 seconds');
    // 连接成功后查询服务器状态，增加延迟确保连接稳定
    setTimeout(() => {
      if (connectionStatus.value === 'connected') {
        console.log('Connection established, querying server status (from Home.vue)'); // 调试日志
        wsApi.send({ command: 'status' });
      } else {
        console.log('Connection status changed before status query, skipping');
      }
    }, 2000); // 增加到2秒，确保连接稳定
  }
});

// 监听控制台状态变化，确保在进入控制台时能正确连接
const consoleStatusWatcher = watch(isInConsole, (inConsole, oldInConsole) => {
  console.log('Console status changed:', inConsole, 'oldInConsole:', oldInConsole, 'connectionStatus:', connectionStatus.value, 'isEnteringConsole:', isEnteringConsole);
  
  // 只有在从false变为true，且连接状态为disconnected，且没有重连尝试，且不是通过enterConsole函数进入时才自动连接
  if (inConsole && !oldInConsole && connectionStatus.value === 'disconnected' && reconnectAttempts.value === 0 && !isEnteringConsole) {
    // 首次进入控制台，尝试连接
    console.log('First time entering console, attempting connection...');
    setTimeout(() => {
      wsApi.connect();
    }, 100);
  }
});

// 标记是否已经初始化过，避免重复连接
let isInitialized = false;
// 标记是否正在通过enterConsole函数进入控制台，避免重复连接
let isEnteringConsole = false;

watch(
  () => route.query,
  (newQuery) => {
    if (route.path === '/dashboard' && newQuery?.wsUrl) {
      // 从 query 组装 server
      const server = {
        name: newQuery.name as string || '',
        wsUrl: newQuery.wsUrl as string || '',
        password: newQuery.password as string || ''
      };
      
      // 只有在不是初始化时才自动进入控制台
      if (isInitialized) {
        console.log('Route query changed, entering console for server:', server);
        enterConsole(server);
      } else {
        console.log('Initial route query detected, but skipping auto-connect to avoid duplicate connections');
        isInitialized = true;
      }
    } else if (route.path !== '/dashboard') {
      isInConsole.value = false;
      wsApi.disconnect();
    }
  },
  { immediate: true }
);

// 组件卸载时清理监听器
onUnmounted(() => {
  console.log('Home component unmounting, cleaning up watchers...');
  if (connectionStatusWatcher) {
    connectionStatusWatcher();
  }
  if (consoleStatusWatcher) {
    consoleStatusWatcher();
  }
  // 断开WebSocket连接
  wsApi.disconnect();
});

// wsApi.onMessage 全局分发终端消息
wsApi.onMessage = (data: any) => {
  console.log('WebSocket received message:', data); // 调试日志
  
  // 同步服务器运行状态
  if (data && data.status) {
    console.log('Updating server status:', data.status); // 调试日志
    isRunning.value = data.status === 'running';
  }
  
  // 处理命令响应
  if (data && data.command) {
    console.log('Received command response:', data.command); // 调试日志
    // 如果是状态查询响应，更新运行状态
    if (data.command === 'status' && data.status) {
      console.log('Status query response:', data.status); // 调试日志
      isRunning.value = data.status === 'running';
    }
  }
  
  // 处理错误消息
  if (data && data.error) {
    console.error('Server error:', data.error); // 调试日志
  }
  
  // 分发到全局消息总线
  if (window && (window as any).__TERMINAL_BUS__) {
    // 过滤掉系统监控、服务器信息、玩家列表和文件管理消息，这些不应该在终端中显示
    if (data && typeof data === 'object' && (data.systemInfo || data.serverInfo || data.players || 
        data.fileList || data.fileContent !== undefined || data.filePath || data.oldPath || data.newPath ||
        // 新增：文件操作成功响应也应该发送到file-message
        (data.status === 'success' && (data.command === 'renameFile' || data.command === 'deleteFile' || data.command === 'createDirectory' || data.command === 'uploadFile' || data.command === 'writeFile' || data.command === 'createZip')))) {
      // 这些是系统监控、玩家管理和文件管理消息，只发送给相应的组件，不发送到终端
      console.log('System/Player/File monitoring message, not sending to terminal:', data);
      // 发送到专门的消息总线，而不是终端总线
      if ((window as any).__SYSTEM_BUS__) {
        console.log('Home: Emitting to system bus:', data);
        (window as any).__SYSTEM_BUS__.emit('system-message', data);
        // 同时发送文件相关消息到专门的文件消息总线
        if (data.fileList || data.fileContent !== undefined || data.filePath || data.oldPath || data.newPath ||
            // 新增：文件操作成功响应也应该发送到file-message
            (data.status === 'success' && (data.command === 'renameFile' || data.command === 'deleteFile' || data.command === 'createDirectory' || data.command === 'uploadFile' || data.command === 'writeFile' || data.command === 'createZip'))) {
          (window as any).__SYSTEM_BUS__.emit('file-message', data);
        }
      } else {
        console.error('Home: SystemBus not available');
      }
    } else {
      console.log('Emitting to terminal bus:', data); // 调试日志
      (window as any).__TERMINAL_BUS__.emit('terminal-message', data);
    }
  } else {
    console.error('TerminalBus not available for message:', data);
  }
};

function showServerList() {
  currentView.value = 'serverList';
  currentComponent.value = ServerList;
  isInConsole.value = false;
  if (route.path !== '/') {
    router.push('/');
  }
}

function showAppSettings() {
  currentView.value = 'appSettings';
  currentComponent.value = AppSettings;
  isInConsole.value = false;
}

function enterConsole(server: any) {
  console.log('Entering console for server:', server);
  
  // 标记正在通过enterConsole函数进入控制台
  isEnteringConsole = true;
  
  // 先断开当前连接
  wsApi.disconnect();
  
  // 重置重连状态和服务器状态
  reconnectAttempts.value = 0;
  isRunning.value = false; // 重置服务器运行状态
  
  // 更新服务器信息
  currentServer.value = server;
  wsUrl.value = server.wsUrl;
  wsPassword.value = server.password;
  
  // 重置控制台相关状态
  if (window && (window as any).__TERMINAL_BUS__) {
    (window as any).__TERMINAL_BUS__.emit('terminal-message', { status: 'stopped' });
  }
  
  // 设置控制台状态（这会触发consoleStatusWatcher，但我们已经修复了重复连接问题）
  isInConsole.value = true;
  consoleSection.value = 'overview';
  currentComponent.value = Overview as unknown as any;
  
  // 延迟连接，确保状态重置完成
  setTimeout(() => {
    console.log('Connecting to WebSocket:', wsUrl.value); // 调试日志
    wsApi.connect();
    
    // 发送测试消息
    setTimeout(() => {
      if (window && (window as any).__TERMINAL_BUS__) {
        console.log('Sending test message to terminal'); // 调试日志
        (window as any).__TERMINAL_BUS__.emit('terminal-message', '正在连接到服务器...');
      }
    }, 500);
    
    // 连接完成后重置标记
    setTimeout(() => {
      isEnteringConsole = false;
    }, 500);
  }, 300); // 增加延迟时间，确保状态变化完成
}

function changeConsoleSection(section: string) {
  consoleSection.value = section;
  switch(section) {
    case 'overview':
      currentComponent.value = Overview as unknown as any;
      break;
    case 'terminal':
      currentComponent.value = Terminal as unknown as any;
      break;
    case 'players':
      currentComponent.value = Players as unknown as any;
      break;
    case 'core':
      currentComponent.value = Core as unknown as any;
      break;
    case 'plugins':
      currentComponent.value = Plugins as unknown as any;
      break;
    case 'files':
      currentComponent.value = Files as unknown as any;
      break;
    case 'remote':
      currentComponent.value = Remote as unknown as any;
      break;
    case 'other':
      currentComponent.value = Other as unknown as any;
      break;
  }
}

const consoleProps = computed(() => {
  if (!isInConsole.value) return {};
  return { server: currentServer.value };
});
</script>

<style scoped>
.layout-root {
  display: flex;
  height: 100vh;
  background: var(--bg-primary);
}

.sidebar {
  width: 280px;
  background: var(--bg-primary);
  color: var(--text-primary);
  display: flex;
  flex-direction: column;
  align-items: stretch;
  border-right: 1px solid var(--border-color);
  margin: 0;
  min-height: 0;
  padding: 20px;
  position: relative;
  z-index: 10;
  height: 100vh;
  overflow: hidden;
  box-sizing: border-box;
}

.sidebar-header {
  padding: 2rem 1rem 1.5rem 1rem;
  text-align: center;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 1rem;
  flex-shrink: 0;
}

.sidebar-header h2 {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 1.4rem;
  margin: 0;
}

.header-actions {
  display: flex;
  justify-content: center;
  margin-top: 1rem;
}

.sidebar-menu {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 1rem 0;
  overflow-y: auto;
  min-height: 0;
  max-height: calc(100vh - 200px);
}

.sidebar-menu::-webkit-scrollbar {
  width: 4px;
}

.sidebar-menu::-webkit-scrollbar-track {
  background: transparent;
}

.sidebar-menu::-webkit-scrollbar-thumb {
  background: rgba(39, 174, 96, 0.2);
  border-radius: 2px;
}

.sidebar-menu::-webkit-scrollbar-thumb:hover {
  background: rgba(39, 174, 96, 0.4);
}

.menu-item {
  background: none;
  border: none;
  color: var(--text-secondary);
  text-align: left;
  padding: 0.8rem 1.2rem;
  font-size: 1rem;
  font-weight: 500;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  outline: none;
  box-shadow: none;
  display: flex;
  align-items: center;
  gap: 0.8em;
  margin-bottom: 0.4rem;
  position: relative;
  overflow: hidden;
  flex-shrink: 0;
}

.menu-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--accent-gradient);
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: -1;
}

.menu-item:hover::before {
  opacity: 0.1;
}

.menu-item.active, .menu-item:hover {
  background: var(--accent-light);
  color: var(--accent-color);
  box-shadow: 0 2px 8px var(--accent-light);
}

.return-btn {
  color: var(--accent-color);
  font-weight: 600;
  background: var(--accent-light);
  border-radius: 8px;
  margin-bottom: 1.2rem;
  border: 1px solid var(--accent-color);
  flex-shrink: 0;
}

.console-menu-header {
  margin: 1.2rem 0 0.8rem 1.5rem;
  font-size: 1.1rem;
  color: var(--accent-color);
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  opacity: 0.8;
  flex-shrink: 0;
}

.console-menu-items {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  min-height: 0;
  max-height: calc(100vh - 400px);
  margin-bottom: 0.8rem;
}

.console-menu-items::-webkit-scrollbar {
  width: 4px;
}

.console-menu-items::-webkit-scrollbar-track {
  background: transparent;
}

.console-menu-items::-webkit-scrollbar-thumb {
  background: rgba(39, 174, 96, 0.2);
  border-radius: 2px;
}

.console-menu-items::-webkit-scrollbar-thumb:hover {
  background: rgba(39, 174, 96, 0.4);
}

.console-item {
  margin-bottom: 0.2rem;
}

.connection-status {
  margin-top: auto;
  padding: 1.2rem;
  background: var(--bg-secondary);
  border-radius: 12px;
  border: 1px solid var(--border-color);
  flex-shrink: 0;
  margin-bottom: 1rem;
  position: sticky;
  bottom: 0;
}

.status-indicator {
  display: flex;
  align-items: center;
  font-size: 0.9rem;
  margin-bottom: 1rem;
  font-weight: 600;
  padding: 0.4rem 0;
}

.status-indicator.connected {
  color: var(--success-color);
}

.status-indicator.connecting {
  color: var(--warning-color);
}

.status-indicator.failed {
  color: var(--error-color);
}

.status-indicator.disconnected {
  color: var(--text-muted);
}

.reconnect-controls {
  margin-bottom: 1rem;
  padding: 0.8rem;
  background: var(--bg-warning);
  border: 1px solid var(--warning-color);
  border-radius: 8px;
}

.error-message {
  color: var(--warning-color);
  font-size: 0.85rem;
  margin-bottom: 0.8rem;
  font-weight: 500;
}

.reconnect-buttons {
  display: flex;
  gap: 0.6rem;
}

.reconnect-btn {
  background: var(--warning-color);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px var(--warning-color);
  flex: 1;
}

.reconnect-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--warning-color);
}

.reconnect-btn:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.connect-btn {
  background: var(--info-color);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px var(--info-color);
  flex: 1;
}

.connect-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--info-color);
}

.connect-btn:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.reset-btn {
  background: var(--text-muted);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px var(--text-muted);
}

.reset-btn:hover {
  box-shadow: 0 4px 12px var(--text-muted);
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.spin {
  animation: spin 1s linear infinite !important;
  display: inline-block;
}

.server-controls {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.control-btn {
  background: var(--accent-gradient);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 0.7rem 1rem;
  font-size: 0.9rem;
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

.control-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.control-btn:hover::before {
  left: 100%;
}

.control-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--accent-light);
}

.control-btn:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.control-btn.start {
  background: var(--accent-gradient);
  box-shadow: 0 2px 8px var(--accent-light);
}

.control-btn.start:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--accent-light);
}

.control-btn.stop {
  background: var(--error-color);
  box-shadow: 0 2px 8px var(--error-color);
}

.control-btn.stop:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--error-color);
}

.control-btn.stop:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  box-shadow: 0 2px 8px var(--shadow-color);
}

.control-btn.restart {
  background: var(--warning-color);
  box-shadow: 0 2px 8px var(--warning-color);
}

.control-btn.restart:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--warning-color);
}

.control-btn.restart:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  box-shadow: 0 2px 8px var(--shadow-color);
}

.reconnect-link {
  margin-top: 0.8rem;
  text-align: center;
  padding: 0.6rem;
  background: var(--accent-light);
  border-radius: 6px;
  border: 1px solid var(--accent-color);
}

.reconnect-a {
  color: var(--accent-color);
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 0.5em;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.reconnect-a:hover {
  color: var(--accent-hover);
  text-decoration: underline;
}

.main-content {
  flex: 1;
  background: var(--bg-primary);
  overflow: auto;
  padding: 1.5rem 2.5rem;
}

/* 响应式设计 */
@media (max-height: 600px) {
  .console-menu-items {
    max-height: calc(100vh - 350px);
  }
  
  .connection-status {
    padding: 1rem;
  }
  
  .control-btn {
    padding: 0.6rem 1rem;
    font-size: 0.85rem;
  }
  
  .menu-item {
    padding: 0.6rem 1rem;
    margin-bottom: 0.3rem;
  }
}

@media (max-height: 500px) {
  .console-menu-items {
    max-height: calc(100vh - 300px);
  }
  
  .sidebar-header {
    padding: 1rem 1rem 1rem 1rem;
  }
  
  .console-menu-header {
    margin: 1rem 0 0.5rem 1.5rem;
  }
  
  .connection-status {
    padding: 0.8rem;
  }
  
  .control-btn {
    padding: 0.5rem 0.8rem;
    font-size: 0.8rem;
  }
}

@media (max-width: 768px) {
  .sidebar {
    width: 250px;
    padding: 15px;
  }
  
  .main-content {
    padding: 1rem 1.5rem;
  }
  
  .menu-item {
    padding: 0.7rem 1rem;
    font-size: 0.95rem;
  }
  
  .control-btn {
    padding: 0.6rem 0.8rem;
    font-size: 0.85rem;
  }
}
</style>