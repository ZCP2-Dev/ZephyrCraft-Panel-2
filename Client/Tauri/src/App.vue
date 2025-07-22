<script setup lang="ts">
import { ref } from "vue";

const wsUrl = ref("");
const ws = ref<WebSocket | null>(null);
const receivedData = ref("");
const command = ref("");

interface Message {
  command: string;
  content: string;
  output: string;
  error: string;
  status: string;
}

function connect() {
  if (ws.value) {
    ws.value.close();
  }
  ws.value = new WebSocket(wsUrl.value);

  ws.value.onopen = () => {
    receivedData.value += "连接成功\n";
  };

  ws.value.onmessage = (event) => {
    try {
      const data: Message = JSON.parse(event.data);
      if (data.output) {
        receivedData.value += `输出: ${data.output}\n`;
      }
      if (data.error) {
        receivedData.value += `错误: ${data.error}\n`;
      }
      if (data.status) {
        receivedData.value += `状态: ${data.status}\n`;
      }
    } catch (err) {
      receivedData.value += `解析数据出错: ${err}\n收到原始数据: ${event.data}\n`;
    }
  };

  ws.value.onclose = () => {
    receivedData.value += "连接已断开\n";
  };

  ws.value.onerror = (error) => {
    receivedData.value += `发生错误: ${error}\n`;
  };
}

function sendMessage() {
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    const message = JSON.stringify({
      command: "input",
      content: command.value
    });
    ws.value.send(message);
    receivedData.value += `已发送命令: ${command.value}\n`;
    command.value = "";
  }
}

function disconnect() {
  if (ws.value) {
    ws.value.close();
  }
}

function startProcess() {
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    const message = JSON.stringify({
      command: "start",
      content: ""
    });
    ws.value.send(message);
    receivedData.value += `已发送启动命令\n`;
  }
}

function stopProcess() {
  if (ws.value && ws.value.readyState === WebSocket.OPEN) {
    const message = JSON.stringify({
      command: "stop",
      content: ""
    });
    ws.value.send(message);
    receivedData.value += `已发送停止命令\n`;
  }
}
</script>

<template>
  <main class="container">
    <h1>WebSocket 测试工具</h1>
    
    <div class="connection-controls">
      <input 
        v-model="wsUrl" 
        placeholder="输入 WebSocket 地址 (ws:// 或 wss://)" 
        class="input-field"
      />
      <button @click="connect" class="btn">连接</button>
      <button @click="disconnect" class="btn">断开连接</button>
      <button @click="startProcess" class="btn">启动进程</button>
      <button @click="stopProcess" class="btn">停止进程</button>
    </div>

    <div class="message-receiver">
      <h2>接收数据</h2>
      <textarea 
        v-model="receivedData" 
        readonly 
        rows="10" 
        class="data-display"
      ></textarea>
    </div>

    <div class="message-sender">
      <h2>发送数据</h2>
      <div class="input-group">
        <input 
          v-model="command" 
          placeholder="输入命令" 
          class="input-field"
          style="width: 415px;"
        />
        <button @click="sendMessage" class="btn">发送</button>
      </div>
    </div>
  </main>
</template>

<style scoped>
.connection-controls, .message-receiver, .message-sender {
  margin-bottom: 20px;
  width: 80%;
  max-width: 800px;
  margin-left: auto;
  margin-right: auto;
}

.input-field {
  margin-right: 10px;
  margin-bottom: 10px;
  width: 200px;
}

.btn {
  margin-right: 10px;
}

.data-display {
  width: 100%;
  padding: 10px;
  border-radius: 8px;
}

.input-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
</style>
<style>
.container {
  padding-top: 2vh;
}

:root {
  font-family: Inter, Avenir, Helvetica, Arial, sans-serif;
  font-size: 16px;
  line-height: 24px;
  font-weight: 400;

  color: #0f0f0f;
  background-color: #f6f6f6;

  font-synthesis: none;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  -webkit-text-size-adjust: 100%;
}

.container {
  margin: 0;
  padding-top: 10vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  text-align: center;
}

.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: 0.75s;
}

.logo.tauri:hover {
  filter: drop-shadow(0 0 2em #24c8db);
}

.row {
  display: flex;
  justify-content: center;
}

a {
  font-weight: 500;
  color: #646cff;
  text-decoration: inherit;
}

a:hover {
  color: #535bf2;
}

h1 {
  text-align: center;
}

input,
button {
  border-radius: 8px;
  border: 1px solid transparent;
  padding: 0.6em 1.2em;
  font-size: 1em;
  font-weight: 500;
  font-family: inherit;
  color: #0f0f0f;
  background-color: #ffffff;
  transition: border-color 0.25s;
  box-shadow: 0 2px 2px rgba(0, 0, 0, 0.2);
}

button {
  cursor: pointer;
}

button:hover {
  border-color: #396cd8;
}
button:active {
  border-color: #396cd8;
  background-color: #e8e8e8;
}

input,
button {
  outline: none;
}

@media (prefers-color-scheme: dark) {
  :root {
    color: #f6f6f6;
    background-color: #2f2f2f;
  }

  a:hover {
    color: #24c8db;
  }

  input,
  button {
    color: #ffffff;
    background-color: #0f0f0f98;
  }
  button:active {
    background-color: #0f0f0f69;
  }
}
</style>