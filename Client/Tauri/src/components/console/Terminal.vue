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
      <button @click="testColors" class="test-btn" style="display: none;">
        <IconMdiPalette style="margin-right: 5px;" />测试颜色
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick, computed, inject } from 'vue';
import { useWebSocket } from '../../useWebSocket';

const props = defineProps<{ server?: any }>();
const commandInput = ref('');
const outputLines = ref<Array<{ text: string }>>([]);
const outputRef = ref<HTMLElement | null>(null);

const wsApi = inject('wsApi') as ReturnType<typeof import('../../useWebSocket').useWebSocket>;
const STORAGE_KEY = computed(() => `terminal_history_${props.server?.wsUrl || ''}`);

// 保存监听器引用以便清理
let terminalMessageListener: ((data: any) => void) | null = null;

function ansiToHtml(str: string): string {
  // 完整的ANSI转义序列解析器，支持Windows 10+ CMD的所有颜色、格式和控制字符
  let result = '';
  let currentIndex = 0;
  let currentStyles: string[] = [];
  
  while (currentIndex < str.length) {
    const char = str[currentIndex];
    
    if (char === '\u001b' && str[currentIndex + 1] === '[') {
      // 找到ANSI转义序列
      let endIndex = str.indexOf('m', currentIndex);
      let commandChar = 'm';
      
      // 检查是否有其他命令字符（如K, C, D等）
      if (endIndex === -1) {
        // 查找其他可能的命令字符
        for (let i = currentIndex + 2; i < str.length; i++) {
          const c = str[i];
          if (c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z') {
            endIndex = i;
            commandChar = c;
            break;
          }
        }
      }
      
      if (endIndex === -1) {
        result += char;
        currentIndex++;
        continue;
      }
      
      const sequence = str.substring(currentIndex + 2, endIndex);
      const codes = sequence.split(';').map(code => parseInt(code) || 0);
      
      // 处理不同的命令字符
      if (commandChar === 'm') {
        // 处理颜色和格式代码
        let newStyles: string[] = [];
        
        for (let i = 0; i < codes.length; i++) {
          const code = codes[i];
          
          switch (code) {
            // 重置所有属性
            case 0:
              newStyles = [];
              break;
              
            // 文本属性
            case 1: // 粗体
              newStyles.push('font-weight: bold');
              break;
            case 2: // 暗淡
              newStyles.push('opacity: 0.6');
              break;
            case 3: // 斜体
              newStyles.push('font-style: italic');
              break;
            case 4: // 下划线
              newStyles.push('text-decoration: underline');
              break;
            case 5: // 闪烁
              newStyles.push('animation: blink 1s infinite');
              break;
            case 7: // 反显
              newStyles.push('background-color: currentColor; color: #2c3e50');
              break;
            case 8: // 隐藏
              newStyles.push('visibility: hidden');
              break;
            case 9: // 删除线
              newStyles.push('text-decoration: line-through');
              break;
              
            // 前景色 (30-37)
            case 30: newStyles.push('color: #000000'); break; // 黑
            case 31: newStyles.push('color: #cd3131'); break; // 红
            case 32: newStyles.push('color: #0dbc79'); break; // 绿
            case 33: newStyles.push('color: #e5e510'); break; // 黄
            case 34: newStyles.push('color: #2472c8'); break; // 蓝
            case 35: newStyles.push('color: #bc3fbc'); break; // 紫
            case 36: newStyles.push('color: #11a8cd'); break; // 青
            case 37: newStyles.push('color: #e5e5e5'); break; // 白
            
            // 前景色高亮 (90-97)
            case 90: newStyles.push('color: #666666'); break; // 亮黑
            case 91: newStyles.push('color: #f14c4c'); break; // 亮红
            case 92: newStyles.push('color: #23d18b'); break; // 亮绿
            case 93: newStyles.push('color: #f5f543'); break; // 亮黄
            case 94: newStyles.push('color: #3b8eea'); break; // 亮蓝
            case 95: newStyles.push('color: #d670d6'); break; // 亮紫
            case 96: newStyles.push('color: #29b8db'); break; // 亮青
            case 97: newStyles.push('color: #ffffff'); break; // 亮白
            
            // 背景色 (40-47)
            case 40: newStyles.push('background-color: #000000'); break; // 黑
            case 41: newStyles.push('background-color: #cd3131'); break; // 红
            case 42: newStyles.push('background-color: #0dbc79'); break; // 绿
            case 43: newStyles.push('background-color: #e5e510'); break; // 黄
            case 44: newStyles.push('background-color: #2472c8'); break; // 蓝
            case 45: newStyles.push('background-color: #bc3fbc'); break; // 紫
            case 46: newStyles.push('background-color: #11a8cd'); break; // 青
            case 47: newStyles.push('background-color: #e5e5e5'); break; // 白
            
            // 背景色高亮 (100-107)
            case 100: newStyles.push('background-color: #666666'); break; // 亮黑
            case 101: newStyles.push('background-color: #f14c4c'); break; // 亮红
            case 102: newStyles.push('background-color: #23d18b'); break; // 亮绿
            case 103: newStyles.push('background-color: #f5f543'); break; // 亮黄
            case 104: newStyles.push('background-color: #3b8eea'); break; // 亮蓝
            case 105: newStyles.push('background-color: #d670d6'); break; // 亮紫
            case 106: newStyles.push('background-color: #29b8db'); break; // 亮青
            case 107: newStyles.push('background-color: #ffffff'); break; // 亮白
          }
        }
        
        // 处理256色模式 (38;5;x 和 48;5;x)
        if (codes.length >= 3 && codes[0] === 38 && codes[1] === 5) {
          const colorIndex = codes[2];
          const color = get256Color(colorIndex);
          newStyles.push(`color: ${color}`);
        } else if (codes.length >= 3 && codes[0] === 48 && codes[1] === 5) {
          const colorIndex = codes[2];
          const color = get256Color(colorIndex);
          newStyles.push(`background-color: ${color}`);
        }
        
        // 处理RGB颜色 (38;2;r;g;b 和 48;2;r;g;b)
        if (codes.length >= 5 && codes[0] === 38 && codes[1] === 2) {
          const r = codes[2];
          const g = codes[3];
          const b = codes[4];
          newStyles.push(`color: rgb(${r}, ${g}, ${b})`);
        } else if (codes.length >= 5 && codes[0] === 48 && codes[1] === 2) {
          const r = codes[2];
          const g = codes[3];
          const b = codes[4];
          newStyles.push(`background-color: rgb(${r}, ${g}, ${b})`);
        }
        
        // 更新当前样式
        if (newStyles.length > 0) {
          currentStyles = newStyles;
        } else {
          currentStyles = [];
        }
      } else if (commandChar === 'K') {
        // 清除行命令 - 在HTML中我们忽略这个，因为每行都是独立的
        // 如果需要实现，可以添加特殊的清除标记
      } else if (commandChar === 'C') {
        // 光标向右移动 - 在HTML中我们忽略这个，因为文本会自动换行
        // 如果需要实现，可以添加空格
        const spaces = codes[0] || 1;
        result += '&nbsp;'.repeat(spaces);
      } else if (commandChar === 'D') {
        // 光标向左移动 - 在HTML中我们忽略这个
      } else if (commandChar === 'A') {
        // 光标向上移动 - 在HTML中我们忽略这个
      } else if (commandChar === 'B') {
        // 光标向下移动 - 在HTML中我们忽略这个
      } else if (commandChar === 'H') {
        // 光标定位 - 在HTML中我们忽略这个
      } else if (commandChar === 'J') {
        // 清除屏幕 - 在HTML中我们忽略这个
      } else {
        // 其他未知的控制字符，忽略它们
        // console.log(`忽略未知ANSI控制字符: ${commandChar} with codes:`, codes);
      }
      
      currentIndex = endIndex + 1;
    } else {
      // 普通字符
      if (currentStyles.length > 0) {
        result += `<span style="${currentStyles.join('; ')}">${escapeHtml(char)}</span>`;
      } else {
        result += escapeHtml(char);
      }
      currentIndex++;
    }
  }
  
  return result;
}

// 获取256色模式的颜色值
function get256Color(index: number): string {
  if (index < 16) {
    // 标准16色
    const colors = [
      '#000000', '#cd3131', '#0dbc79', '#e5e510', '#2472c8', '#bc3fbc', '#11a8cd', '#e5e5e5',
      '#666666', '#f14c4c', '#23d18b', '#f5f543', '#3b8eea', '#d670d6', '#29b8db', '#ffffff'
    ];
    return colors[index] || '#ffffff';
  } else if (index < 232) {
    // 216色立方体
    const cubeIndex = index - 16;
    const r = Math.floor(cubeIndex / 36) * 51;
    const g = Math.floor((cubeIndex % 36) / 6) * 51;
    const b = (cubeIndex % 6) * 51;
    return `rgb(${r}, ${g}, ${b})`;
  } else {
    // 24级灰度
    const gray = (index - 232) * 10 + 8;
    return `rgb(${gray}, ${gray}, ${gray})`;
  }
}

// HTML转义函数
function escapeHtml(text: string): string {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
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

function testColors() {
  appendLine('\n=== ANSI颜色和控制字符测试 ===');
  appendLine('\u001b[31m红色文本\u001b[0m \u001b[32m绿色文本\u001b[0m \u001b[33m黄色文本\u001b[0m \u001b[34m蓝色文本\u001b[0m');
  appendLine('\u001b[35m紫色文本\u001b[0m \u001b[36m青色文本\u001b[0m \u001b[37m白色文本\u001b[0m \u001b[30m黑色文本\u001b[0m');
  appendLine('\u001b[1m\u001b[31m粗体红色\u001b[0m \u001b[4m\u001b[32m下划线绿色\u001b[0m \u001b[3m\u001b[33m斜体黄色\u001b[0m');
  appendLine('\u001b[38;5;9m256色红色\u001b[0m \u001b[48;5;12m256色背景\u001b[0m \u001b[38;2;255;100;100mRGB红色\u001b[0m');
  appendLine('\u001b[7m反显文本\u001b[0m \u001b[9m删除线\u001b[0m \u001b[5m闪烁文本\u001b[0m');
  appendLine('测试控制字符: 文本\u001b[K清除行\u001b[0m');
  appendLine('测试光标移动: 开始\u001b[10C中间\u001b[10C结束');
  appendLine('=== 测试完成 ===\n');
  scrollToBottom();
}

const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);

function sendCommand() {
  if (!commandInput.value.trim() || !isConnected.value) return;

  const message = commandInput.value;
  wsApi.send({ command: 'input', content: message });
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

  // 添加测试消息和颜色示例
  // appendLine('终端已初始化，等待连接...');
  // appendLine('\u001b[32m✓ 彩色文本支持已启用\u001b[0m');
  // appendLine('\u001b[31m红色文本\u001b[0m \u001b[33m黄色文本\u001b[0m \u001b[34m蓝色文本\u001b[0m');
  // appendLine('\u001b[1m\u001b[36m粗体青色文本\u001b[0m \u001b[4m下划线文本\u001b[0m');
  // appendLine('\u001b[38;5;9m256色模式测试\u001b[0m \u001b[48;5;12m背景色测试\u001b[0m');
  // appendLine('控制字符测试: 文本\u001b[K清除行\u001b[0m');
  // appendLine('光标移动测试: 开始\u001b[10C中间\u001b[10C结束');

  // 监听全局终端消息
  const bus = (window as any).__TERMINAL_BUS__;
  if (bus && typeof bus.on === 'function') {
    console.log('Terminal: TerminalBus found, setting up listener'); // 调试日志
    terminalMessageListener = (data: any) => {
      console.log('Terminal received message:', data); // 调试日志

      if (typeof data === 'string') {
        // appendLine(data);
        // scrollToBottom();
        return;
      }

      if (data && typeof data === 'object') {
        // 过滤掉系统监控和服务器信息消息，这些不应该在终端中显示
        if (data.systemInfo || data.serverInfo) {
          // 这些是系统监控消息，不在终端中显示，但允许其他组件接收
          return;
        }
        
        if (data.output) {
          appendLine(data.output);
        }
        if (data.error) {
          appendLine(`错误: ${data.error}`);
        }
        if (data.status) {
          // appendLine(`状态: ${data.status}`);
        }
        if (data.content) {
          appendLine(data.content);
        }
        if (data.message) {
          appendLine(data.message);
        }
        // 如果没有特定字段，尝试直接显示整个数据
        if (!data.output && !data.error && !data.status && !data.content && !data.message) {
          // appendLine(JSON.stringify(data));
        }
        scrollToBottom();
      }
    };
    bus.on('terminal-message', terminalMessageListener);
  } else {
    console.error('TerminalBus not found or invalid');
    appendLine('错误: 终端消息总线未初始化');
  }

  scrollToBottom();
});

// 组件卸载时清理监听器
onUnmounted(() => {
  console.log('Terminal: Component unmounting, cleaning up listeners...');
  
  // 清理终端消息监听器
  if (terminalMessageListener) {
    const bus = (window as any).__TERMINAL_BUS__;
    if (bus && typeof bus.off === 'function') {
      bus.off('terminal-message', terminalMessageListener);
    }
    terminalMessageListener = null;
  }
  
  console.log('Terminal: Listeners cleaned up');
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
/* 闪烁动画 */
@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

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
  display: inline;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 确保彩色文本在深色背景上清晰可见 */
.terminal-output .ansi-line span {
  border-radius: 2px;
  padding: 0 1px;
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

.test-btn {
  background: linear-gradient(135deg, #9b59b6 0%, #8e44ad 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(155, 89, 182, 0.3);
}

.test-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(155, 89, 182, 0.4);
}

.test-btn:disabled {
  background: linear-gradient(135deg, #ecf0f1 0%, #bdc3c7 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}
</style>