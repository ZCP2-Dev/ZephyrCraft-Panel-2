<template>
  <div class="terminal-container">
    <h2>终端控制台</h2>
    <div class="terminal-output" ref="outputRef">
      <div v-for="(line, index) in outputLines" :key="index" class="output-line">
        <span class="ansi-line" v-html="line.text"></span>
      </div>
    </div>
    <div class="terminal-input-container">
      <input v-model="commandInput" @keyup.enter="sendCommand" placeholder="输入命令..." class="terminal-input" />
      <button @click="sendCommand" class="send-btn" :disabled="!isConnected">
        <IconMdiSend style="margin-right: 5px;" />发送
      </button>
      <button @click="clearOutput" class="clear-btn">
        <IconMdiDelete style="margin-right: 5px;" />清除
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, nextTick, computed, inject } from 'vue';
import { useWebSocket } from '../../useWebSocket';

const props = defineProps<{ server?: any }>();
const commandInput = ref('');
const outputLines = ref<Array<{ text: string }>>([]);
const outputRef = ref<HTMLElement | null>(null);

const wsApi = inject('wsApi') as ReturnType<typeof import('../../useWebSocket').useWebSocket>;
const STORAGE_KEY = computed(() => `terminal_history_${props.server?.wsUrl || ''}`);

function ansiToHtml(str: string): string {
  // 只支持常见的 38;5;X 颜色序列，简化实现
  // 例如：\u001b[38;5;9m红色  \u001b[38;5;15m白色  \u001b[m重置
  const colorMap: Record<string, string> = {
    '9': '#ff6b6b', // 红
    '15': '#ffffff', // 白
    // 可扩展更多
  };
  return str
    .replace(/\u001b\[38;5;(\d+)m/g, (_, color) => `<span style=\"color:${colorMap[color] || '#fff'}\">`)
    .replace(/\u001b\[m/g, '</span>')
    .replace(/\u001b\[[0-9;]*m/g, '</span>'); // 其它重置
}

function appendLine(raw: string) {
  const html = ansiToHtml(raw);
  outputLines.value.push({ text: html });
  // 存储到localStorage，最多100条
  const arr = outputLines.value.slice(-100).map(l => l.text);
  localStorage.setItem(STORAGE_KEY.value, JSON.stringify(arr));
}

function clearOutput() {
  outputLines.value = [];
  localStorage.removeItem(STORAGE_KEY.value);
}

const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);

function sendCommand() {
  if (!commandInput.value.trim() || !isConnected.value) return;

  const command = commandInput.value;
  wsApi.send({ command });
  commandInput.value = '';
}

function scrollToBottom() {
  nextTick(() => {
    if (outputRef.value) {
      outputRef.value.scrollTop = outputRef.value.scrollHeight;
    }
  });
}

// 监听 wsApi 消息
onMounted(() => {
  // 载入历史
  const history = localStorage.getItem(STORAGE_KEY.value);
  if (history) {
    try {
      const arr = JSON.parse(history);
      outputLines.value = arr.map((text: string) => ({ text }));
    } catch { }
  }

  // 添加测试消息
  appendLine('终端已初始化，等待连接...');

  // 监听全局终端消息
  const bus = (window as any).__TERMINAL_BUS__;
  if (bus && typeof bus.on === 'function') {
    console.log('Terminal: TerminalBus found, setting up listener'); // 调试日志
    bus.on('terminal-message', (data: any) => {
      console.log('Terminal received message:', data); // 调试日志

      if (typeof data === 'string') {
        appendLine(data);
        scrollToBottom();
        return;
      }

      if (data && typeof data === 'object') {
        if (data.output) {
          appendLine(data.output);
        }
        if (data.error) {
          appendLine(`错误: ${data.error}`);
        }
        if (data.status) {
          appendLine(`状态: ${data.status}`);
        }
        if (data.content) {
          appendLine(data.content);
        }
        if (data.message) {
          appendLine(data.message);
        }
        // 如果没有特定字段，尝试直接显示整个数据
        if (!data.output && !data.error && !data.status && !data.content && !data.message) {
          appendLine(JSON.stringify(data));
        }
        scrollToBottom();
      }
    });
  } else {
    console.error('TerminalBus not found or invalid');
    appendLine('错误: 终端消息总线未初始化');
  }

  scrollToBottom();
});

watch(() => props.server, (newServer) => {
  // 载入新服务器的历史
  const history = localStorage.getItem(`terminal_history_${newServer?.wsUrl || ''}`);
  if (history) {
    try {
      const arr = JSON.parse(history);
      outputLines.value = arr.map((text: string) => ({ text }));
    } catch { }
  } else {
    outputLines.value = [];
  }
  scrollToBottom();
});
</script>

<style scoped>
.terminal-container {
  background: #ffffff;
  padding: 1.5rem;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.terminal-container h2 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.terminal-output {
  background: #2c3e50;
  color: #ecf0f1;
  padding: 1.2rem;
  border-radius: 8px;
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
  line-height: 1.5;
  height: 400px;
  overflow-y: auto;
  margin-bottom: 1.5rem;
  border: 1px solid #34495e;
  flex: 1;
}

.terminal-output::-webkit-scrollbar {
  width: 8px;
}

.terminal-output::-webkit-scrollbar-track {
  background: #34495e;
  border-radius: 4px;
}

.terminal-output::-webkit-scrollbar-thumb {
  background: #7f8c8d;
  border-radius: 4px;
}

.terminal-output::-webkit-scrollbar-thumb:hover {
  background: #95a5a6;
}

.output-line {
  margin-bottom: 0.5rem;
  line-height: 1.4;
}

.ansi-line span {
  font-family: inherit;
  font-size: inherit;
}

.terminal-input-container {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.terminal-input {
  flex: 1;
  padding: 0.8rem 1rem;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 1rem;
  font-family: 'Courier New', monospace;
  background: #ffffff;
  color: #2c3e50;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.terminal-input:focus {
  outline: none;
  border-color: #27ae60;
  box-shadow: 0 0 0 3px rgba(39, 174, 96, 0.1);
  background: #ffffff;
}

.send-btn,
.clear-btn {
  padding: 0.8rem 1.2rem;
  border: none;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  position: relative;
  overflow: hidden;
}

.send-btn::before,
.clear-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.send-btn:hover::before,
.clear-btn:hover::before {
  left: 100%;
}

.send-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
}

.send-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

.send-btn:disabled {
  background: linear-gradient(135deg, #bdc3c7 0%, #95a5a6 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.clear-btn {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(231, 76, 60, 0.3);
}

.clear-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(231, 76, 60, 0.4);
}

.clear-btn:disabled {
  background: linear-gradient(135deg, #ecf0f1 0%, #bdc3c7 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
</style>