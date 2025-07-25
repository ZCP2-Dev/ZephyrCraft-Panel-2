<template>
  <div class="terminal-container">
    <div class="terminal-output" ref="outputRef">
      <div v-for="(line, index) in outputLines" :key="index" class="output-line">
        <span class="ansi-line" v-html="line.text"></span>
      </div>
    </div>
    <div class="terminal-input-area">
      <input
        v-model="commandInput"
        @keyup.enter="sendCommand"
        placeholder="输入命令..."
        class="command-input-field"
      />
      <button @click="sendCommand" class="send-btn">发送</button>
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

const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);
// 监听 wsApi 消息
onMounted(() => {
  // 载入历史
  const history = localStorage.getItem(STORAGE_KEY.value);
  if (history) {
    try {
      const arr = JSON.parse(history);
      outputLines.value = arr.map((text: string) => ({ text }));
    } catch {}
  }
  // 监听全局终端消息
  const bus = (window as any).__TERMINAL_BUS__;
  if (bus && typeof bus.on === 'function') {
    bus.on('terminal-message', (data: any) => {
      if (typeof data === 'string') {
        appendLine(data);
        return;
      }
      if (data.output) {
        appendLine(data.output);
      }
      if (data.error) {
        appendLine(data.error);
      }
      if (data.status) {
        appendLine(data.status);
      }
      scrollToBottom();
    });
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
    } catch {}
  } else {
    outputLines.value = [];
  }
});

const sendCommand = () => {
  if (!commandInput.value.trim() || !isConnected.value) return;
  const input = commandInput.value.trim();
  wsApi.send({ command: 'input', content: input });
  commandInput.value = '';
  scrollToBottom();
};

function scrollToBottom() {
  nextTick(() => {
    if (outputRef.value) {
      outputRef.value.scrollTop = outputRef.value.scrollHeight;
    }
  });
}
</script>

<style scoped>
.terminal-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}
.terminal-output {
  background-color: #1e1e1e;
  color: #ffffff;
  border-radius: 10px 10px 0 0;
  padding: 1.2rem 1.5rem;
  height: inherit;
  overflow-y: auto;
  font-family: 'Courier New', monospace;
  font-size: 0.98rem;
}
.output-line {
  margin-bottom: 0.5rem;
  line-height: 1.4;
}
.ansi-line span {
  font-family: inherit;
  font-size: inherit;
}
.terminal-input-area {
  display: flex;
  border-top: 1px solid #eee;
  background: #fafbfc;
  border-radius: 0 0 10px 10px;
  padding: 1rem 1.5rem;
}
.command-input-field {
  flex: 1;
  padding: 0.8rem;
  border: none;
  border-radius: 6px 0 0 6px;
  font-family: 'Courier New', monospace;
  font-size: 1rem;
  background-color: #333;
  color: white;
}
.send-btn {
  background: #88bf64;
  color: #fff;
  border: none;
  border-radius: 0 6px 6px 0;
  padding: 0 1.5rem;
  font-size: 1rem;
  cursor: pointer;
  transition: background 0.2s;
}
.send-btn:hover {
  background: #539951;
}
</style>