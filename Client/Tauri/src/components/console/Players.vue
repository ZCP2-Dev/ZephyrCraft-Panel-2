<template>
  <div class="players-container" :key="refreshKey">
    <h2>玩家管理</h2>
    
    <div class="players-stats">
      <div class="stat-card">
        <div class="stat-title">在线玩家</div>
        <div class="stat-value">{{ players.length }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-title">最大玩家数</div>
        <div class="stat-value">{{ maxPlayers }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-title">服务器状态</div>
        <div class="stat-value" :class="globalIsRunning ? 'status-running' : 'status-stopped'">
          {{ globalIsRunning ? '运行中' : '已停止' }}
        </div>
      </div>
    </div>

    <div class="players-section">
      <h3>在线玩家</h3>
      <div class="players-list">
        <div v-for="player in players" :key="player.name" class="player-item">
          <div class="player-info">
            <div class="player-avatar">
              {{ player.name.charAt(0).toUpperCase() }}
            </div>
            <div class="player-details">
              <div class="player-name">{{ player.name }}</div>
              <div class="player-status">在线</div>
              <div v-if="player.xuid" class="player-xuid">XUID: {{ player.xuid }}</div>
            </div>
          </div>
          <div class="player-actions">
            <button @click="sendMessage(player.name)" class="message-btn" :disabled="!globalIsRunning">
              <IconMdiMessage style="margin-right: 5px;" />消息
            </button>
            <button @click="kickPlayer(player.name)" class="kick-btn" :disabled="!globalIsRunning">
              <IconMdiAccountRemove style="margin-right: 5px;" />踢出
            </button>
          </div>
        </div>
        <div v-if="players.length === 0" class="no-players">
          {{ globalIsRunning ? '暂无在线玩家' : '服务器未运行' }}
        </div>
      </div>
    </div>

    <div class="players-section">
      <h3>玩家操作</h3>
      <div class="action-buttons">
        <button @click="forceReload" class="action-btn" :disabled="!globalIsRunning">
          <IconMdiRefresh style="margin-right: 5px;" />刷新玩家列表
        </button>
        <!-- <button @click="sendCommand('whitelist list')" class="action-btn" :disabled="!globalIsRunning">
          <IconMdiAccountGroup style="margin-right: 5px;" />查看白名单
        </button>
        <button @click="sendCommand('list')" class="action-btn" :disabled="!globalIsRunning">
          <IconMdiAccountGroup style="margin-right: 5px;" />查看在线玩家
        </button> -->
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, watch, onUnmounted, onActivated } from 'vue';

const props = defineProps<{ server?: any }>();
const wsApi = inject('wsApi') as any;
const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);
const isInConsole = inject('isInConsole') as any;

// 玩家列表
const players = ref<Array<{ name: string; xuid?: string }>>([]);
const maxPlayers = ref(20);

// 从全局状态获取服务器运行状态
const globalIsRunning = inject('isRunning') as any;

// 新增：用于强制刷新组件的key
const refreshKey = ref(0);

let refreshTimer: number | null = null;

// 保存监听器引用以便清理
let systemMessageListener: ((data: any) => void) | null = null;
let terminalMessageListener: ((data: any) => void) | null = null;

// 踢出玩家
function kickPlayer(playerName: string) {
  if (wsApi && isConnected.value && globalIsRunning) {
    wsApi.send({ command: 'input', content: `kick ${playerName}` });
    // 踢出玩家后延迟刷新列表
    setTimeout(() => {
      requestPlayers();
    }, 1000);
  }
}

// 发送消息给玩家
function sendMessage(playerName: string) {
  const message = prompt(`发送消息给 ${playerName}:`);
  if (message && wsApi && isConnected.value && globalIsRunning) {
    wsApi.send({ command: 'input', content: `tell ${playerName} ${message}` });
  }
}


// 请求玩家列表
function requestPlayers() {
  requestServerInfo();
  console.log('Players: Requesting players list, wsApi:', !!wsApi, 'isConnected:', isConnected.value, 'isRunning:', globalIsRunning);
  if (wsApi && isConnected.value && globalIsRunning) {
    console.log('Players: Sending getPlayers command');
    wsApi.send({ command: 'getPlayers' });
    // 显示刷新状态
    console.log('Players: Refresh request sent successfully');
  } else {
    console.log('Players: Cannot send getPlayers - wsApi:', !!wsApi, 'isConnected:', isConnected.value, 'isRunning:', globalIsRunning);
  }
}

function requestServerInfo() {
  if (wsApi && isConnected.value) {
    console.log('Players: Sending getServerInfo command');
    wsApi.send({ command: 'getServerInfo' });
  } else {
    console.log('Players: Cannot send getServerInfo - wsApi:', !!wsApi, 'isConnected:', isConnected.value);
  }
}

// 修改刷新按钮逻辑
function forceReload() {
  refreshKey.value++;
}


// 组件激活时刷新
function handleActivated() {
  console.log('Players: Component activated, refreshing...');
  if (globalIsRunning) {
    requestPlayers();
  }
}

// 开始定时刷新
function startRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer);
  }
  refreshTimer = setInterval(() => {
    if (globalIsRunning) {
      requestPlayers();
    }
  }, 10000); // 每10秒刷新一次
}

// 停止定时刷新
function stopRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer);
    refreshTimer = null;
  }
}

onMounted(() => {
  console.log('Players: Component mounted, setting up listeners...');
  
  // 清理之前的监听器（如果存在）
  if (systemMessageListener) {
    const systemBus = (window as any).__SYSTEM_BUS__;
    if (systemBus && typeof systemBus.off === 'function') {
      systemBus.off('system-message', systemMessageListener);
    }
    systemMessageListener = null;
  }
  
  if (terminalMessageListener) {
    const terminalBus = (window as any).__TERMINAL_BUS__;
    if (terminalBus && typeof terminalBus.off === 'function') {
      terminalBus.off('terminal-message', terminalMessageListener);
    }
    terminalMessageListener = null;
  }
  
  // 重置玩家列表状态
  players.value = [];
  
  // 监听全局系统监控消息总线和终端消息总线
  const systemBus = (window as any).__SYSTEM_BUS__;
  const terminalBus = (window as any).__TERMINAL_BUS__;
  
  console.log('Players: SystemBus available:', !!systemBus);
  console.log('Players: TerminalBus available:', !!terminalBus);
  
  if (systemBus && typeof systemBus.on === 'function') {
    systemMessageListener = (data: any) => {
      console.log('Players: Received system message:', data);
      if (data && data.serverInfo) {
        maxPlayers.value = data.serverInfo.maxPlayers || 20;
      }
      if (data && data.players) {
        console.log('Players: Received players list:', data.players);
        const oldCount = players.value.length;
        players.value = data.players;
        const newCount = players.value.length;
        console.log(`Players: Updated players list, count: ${oldCount} -> ${newCount}`);
        
        // 如果玩家数量减少，记录详细信息
        if (newCount < oldCount) {
          console.log('Players: Player count decreased, this might indicate a player left');
        }
      }
    };
    systemBus.on('system-message', systemMessageListener);
    console.log('Players: System message listener registered');
  } else {
    console.error('Players: SystemBus not available or does not have on method');
  }
  
  if (terminalBus && typeof terminalBus.on === 'function') {
    terminalMessageListener = (data: any) => {
      if (data && data.status) {
        // 服务器状态变化时清空玩家列表
        if (data.status !== 'running') {
          players.value = [];
        }
      }
      // 监听命令输出，如果是玩家相关命令，延迟刷新列表
      if (data && data.output) {
        const output = data.output.toLowerCase();
        if (output.includes('whitelist') || output.includes('list') || output.includes('kick') || output.includes('player')) {
          console.log('Players: Detected player-related command output, refreshing list...');
          setTimeout(() => {
            requestPlayers();
          }, 500);
        }
        
        // 特别监听玩家退出事件
        if (output.includes('left the game') || output.includes('disconnected') || output.includes('left')) {
          console.log('Players: Detected player leave/disconnect event, refreshing list...');
          setTimeout(() => {
            requestPlayers();
          }, 300);
        }
      }
    };
    terminalBus.on('terminal-message', terminalMessageListener);
  }

  // 立即请求一次玩家列表
  if (globalIsRunning) {
    console.log('Players: Server is running, requesting initial players list...');
    requestPlayers();
  } else {
    console.log('Players: Server is not running, skipping initial request');
  }
  
  // 开始定时刷新
  startRefresh();
  
  // 添加调试：检查当前玩家列表状态
  console.log('Players: Current players list:', players.value);
  console.log('Players: Current globalIsRunning:', globalIsRunning);
});

// 组件激活时强制刷新
onActivated(() => {
  handleActivated();
});

onUnmounted(() => {
  console.log('Players: Component unmounting, cleaning up listeners...');
  
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
  
  console.log('Players: Listeners cleaned up');
});

watch(() => props.server, () => {
  // 重新请求信息
  if (globalIsRunning) {
    requestPlayers();
  }
});

watch(isConnected, (connected) => {
  if (connected) {
    if (globalIsRunning) {
      requestPlayers();
    }
    startRefresh();
  } else {
    stopRefresh();
  }
});

// 监听服务器运行状态变化
watch(globalIsRunning, (running) => {
  if (running) {
    // 服务器启动时请求玩家列表
    requestPlayers();
    startRefresh();
  } else {
    // 服务器停止时清空玩家列表并停止刷新
    players.value = [];
    stopRefresh();
  }
});

// 监听组件激活状态，确保在组件激活时重新获取最新数据
watch(() => isInConsole, (inConsole) => {
  if (inConsole && globalIsRunning) {
    console.log('Players: Console activated, refreshing players list...');
    // 延迟一点时间确保组件完全激活
    setTimeout(() => {
      requestPlayers();
    }, 200);
  }
}, { immediate: true });


</script>

<style scoped>
.players-container {
  /* background: var(--bg-primary); */
  padding: 1.5rem;
}

.players-container h2 {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.players-stats {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}

.stat-card {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 1.2rem 1.5rem;
  text-align: center;
  border: 1px solid var(--border-color);
  flex: 1;
  min-width: 150px;
}

.stat-title {
  color: var(--text-secondary);
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
  font-weight: 600;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: 700;
  color: var(--text-primary);
}

.players-section {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  border: 1px solid var(--border-color);
}

.players-section h3 {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 1.4rem;
  margin: 0 0 1.2rem 0;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.players-list {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.player-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.2rem;
  background: var(--bg-primary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.player-item:hover {
  box-shadow: 0 2px 8px var(--accent-light);
  border-color: var(--accent-color);
}

.player-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.player-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--accent-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 1.2rem;
}

.player-details {
  display: flex;
  flex-direction: column;
}

.player-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 1.1rem;
}

.player-status {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.player-xuid {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-top: 0.2rem;
}

.player-actions {
  display: flex;
  gap: 0.8rem;
}

.message-btn {
  background: var(--info-color);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px var(--info-color);
}

.message-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--info-color);
}

.message-btn:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.kick-btn {
  background: var(--error-color);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px var(--error-color);
}

.kick-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--error-color);
}

.kick-btn:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.action-buttons {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.action-btn {
  background: var(--accent-gradient);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.8rem 1.2rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px var(--accent-light);
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

.no-players {
  text-align: center;
  padding: 3rem 2rem;
  color: var(--text-secondary);
  font-size: 1.1rem;
  background: var(--bg-primary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.status-running {
  color: var(--success-color);
  font-weight: 600;
}

.status-stopped {
  color: var(--error-color);
  font-weight: 600;
}
</style>