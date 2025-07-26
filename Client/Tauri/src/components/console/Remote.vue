<template>
  <div class="remote-container">
    <h3>远程连接设置</h3>
    <div class="remote-form">
      <div class="form-group">
        <label>WebSocket 地址</label>
        <input 
          v-model="newWsUrl" 
          type="text" 
          placeholder="ws://localhost:1145"
          class="form-input"
        />
      </div>
      <div class="form-group">
        <label>连接密码</label>
        <input 
          v-model="newPassword" 
          type="password" 
          placeholder="输入连接密码"
          class="form-input"
        />
      </div>
      <div class="form-actions">
        <button @click="saveSettings" class="save-btn" :disabled="!isValid">
          <IconMdiContentSave style="margin-right: 5px;" />保存设置
        </button>
        <button @click="testConnection" class="test-btn" :disabled="!isValid || !isConnected">
          <IconMdiConnection style="margin-right: 5px;" />测试连接
        </button>
      </div>
    </div>
    
    <div class="connection-info">
      <h4>当前连接信息</h4>
      <div class="info-item">
        <span class="info-label">地址：</span>
        <span class="info-value">{{ server?.wsUrl || '未设置' }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">状态：</span>
        <span class="info-value" :class="{ 'connected': isConnected, 'disconnected': !isConnected }">
          {{ isConnected ? '已连接' : '未连接' }}
        </span>
      </div>
      <div class="info-item">
        <span class="info-label">服务器：</span>
        <span class="info-value" :class="{ 'running': isRunning, 'stopped': !isRunning }">
          {{ isRunning ? '运行中' : '已停止' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, onUnmounted } from 'vue';

const props = defineProps<{ server?: any }>();
const wsApi = inject('wsApi') as any;
const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);
const isRunning = ref(false);

const newWsUrl = ref('');
const newPassword = ref('');

// 保存监听器引用以便清理
let terminalMessageListener: ((data: any) => void) | null = null;

const isValid = computed(() => {
  return newWsUrl.value.trim() && newPassword.value.trim();
});

onMounted(() => {
  // 初始化表单数据
  newWsUrl.value = props.server?.wsUrl || '';
  newPassword.value = props.server?.password || '';
  
  // 监听服务器状态
  const bus = (window as any).__TERMINAL_BUS__;
  if (bus && typeof bus.on === 'function') {
    terminalMessageListener = (data: any) => {
      if (data && data.status) {
        isRunning.value = data.status === 'running';
      }
    };
    bus.on('terminal-message', terminalMessageListener);
  }
});

// 组件卸载时清理监听器
onUnmounted(() => {
  console.log('Remote: Component unmounting, cleaning up listeners...');
  
  // 清理终端消息监听器
  if (terminalMessageListener) {
    const bus = (window as any).__TERMINAL_BUS__;
    if (bus && typeof bus.off === 'function') {
      bus.off('terminal-message', terminalMessageListener);
    }
    terminalMessageListener = null;
  }
  
  console.log('Remote: Listeners cleaned up');
});

function saveSettings() {
  if (!isValid.value) return;
  
  // 更新服务器信息
  if (props.server) {
    props.server.wsUrl = newWsUrl.value.trim();
    props.server.password = newPassword.value.trim();
  }
  
  // 重新连接
  wsApi.disconnect();
  setTimeout(() => {
    wsApi.connect();
  }, 500);
}

function testConnection() {
  if (!isValid.value) return;
  
  // 临时连接测试
  const testWs = new WebSocket(newWsUrl.value.trim());
  
  testWs.onopen = () => {
    alert('连接测试成功！');
    testWs.close();
  };
  
  testWs.onerror = () => {
    alert('连接测试失败，请检查地址和网络连接。');
  };
}
</script>

<style scoped>
.remote-container {
  background: #ffffff;
  padding: 1.5rem;
}

.remote-container h3 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.8rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.remote-form {
  margin-bottom: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.6rem;
  font-weight: 600;
  color: #2c3e50;
  font-size: 1rem;
}

.form-input {
  width: 100%;
  padding: 0.8rem 1rem;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: #ffffff;
  color: #2c3e50;
}

.form-input:focus {
  outline: none;
  border-color: #27ae60;
  box-shadow: 0 0 0 3px rgba(39, 174, 96, 0.1);
  background: #ffffff;
}

.form-input::placeholder {
  color: #95a5a6;
}

.form-actions {
  display: flex;
  gap: 1.2rem;
  margin-top: 2rem;
}

.save-btn, .test-btn {
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  position: relative;
  overflow: hidden;
}

.save-btn::before, .test-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.save-btn:hover::before, .test-btn:hover::before {
  left: 100%;
}

.save-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
}

.save-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

.save-btn:disabled {
  background: linear-gradient(135deg, #bdc3c7 0%, #95a5a6 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.test-btn {
  background: linear-gradient(135deg, #f39c12 0%, #e67e22 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(243, 156, 18, 0.3);
}

.test-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(243, 156, 18, 0.4);
}

.test-btn:disabled {
  background: linear-gradient(135deg, #ecf0f1 0%, #bdc3c7 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.connection-info {
  border-top: 2px solid #e9ecef;
  padding-top: 2rem;
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.5rem;
  margin-top: 1.5rem;
}

.connection-info h4 {
  margin-bottom: 1.2rem;
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.3rem;
  text-align: center;
}

.info-item {
  display: flex;
  margin-bottom: 0.8rem;
  padding: 0.7rem 0.8rem;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #e9ecef;
}

.info-label {
  min-width: 100px;
  font-weight: 600;
  color: #7f8c8d;
  font-size: 0.95rem;
}

.info-value {
  color: #2c3e50;
  font-weight: 500;
  flex: 1;
}

.info-value.connected {
  color: #27ae60;
  font-weight: 600;
}

.info-value.disconnected {
  color: #e74c3c;
  font-weight: 600;
}

.info-value.running {
  color: #27ae60;
  font-weight: 600;
}

.info-value.stopped {
  color: #e74c3c;
  font-weight: 600;
}
</style> 