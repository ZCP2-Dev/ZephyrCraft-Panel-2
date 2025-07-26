<template>
  <div id="app">
    <router-view class="view-content" />
  </div>
</template>

<script setup lang="ts">
import { useTheme } from './composables/useTheme'

// 初始化主题
useTheme()
</script>

<style scoped>
.navbar {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background-color: #f0f0f0;
  border-bottom: 1px solid #e0e0e0;
}

.nav-link {
  padding: 0.5rem 1rem;
  text-decoration: none;
  color: #333;
  border-radius: 4px;
}

.nav-link:hover {
  background-color: #e0e0e0;
}

.nav-link.router-link-exact-active {
  color: #88bf64;
  font-weight: bold;
  border-bottom: 2px solid #88bf64;
}

.view-content {
  box-sizing: border-box;
}
</style>

<style>
/* CSS变量定义 */
:root {
  /* 浅色主题变量 */
  --bg-primary: #ffffff;
  --bg-secondary: #f8f9fa;
  --bg-tertiary: #f0f0f0;
  --bg-warning: #f8a9295d;
  --bg-terminal: #2c3e50;
  --text-primary: #2c3e50;
  --text-secondary: #6c757d;
  --text-muted: #95a5a6;
  --border-color: #e0e0e0;
  --border-light: #f0f0f0;
  --shadow-color: rgba(0, 0, 0, 0.1);
  --shadow-color-light: rgba(0, 0, 0, 0.05);
  --accent-color: #27ae60;
  --accent-hover: #229954;
  --accent-light: rgba(39, 174, 96, 0.1);
  --accent-gradient: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  --accent-gradient-hover: linear-gradient(135deg, #229954 0%, #27ae60 100%);
  --success-color: #27ae60;
  --warning-color: #ff8800;
  --error-color: #e74c3c;
  --info-color: #3498db;
}

[data-theme="dark"] {
  /* 深色主题变量 */
  --bg-primary: #1a1a1a;
  --bg-secondary: #2d2d2d;
  --bg-tertiary: #3d3d3d;
  --bg-warning: #f39d1269;
  --bg-terminal: #2d2d2d;
  --text-primary: #ffffff;
  --text-secondary: #b0b0b0;
  --text-muted: #888888;
  --border-color: #404040;
  --border-light: #505050;
  --shadow-color: rgba(0, 0, 0, 0.3);
  --shadow-color-light: rgba(0, 0, 0, 0.2);
  --accent-color: #4ade80;
  --accent-hover: #22c55e;
  --accent-light: rgba(74, 222, 128, 0.1);
  --accent-gradient: linear-gradient(135deg, #4ade80 0%, #22c55e 100%);
  --accent-gradient-hover: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  --success-color: #4ade80;
  --warning-color: #ff7d32;
  --error-color: #f87171;
  --info-color: #60a5fa;
}

[data-theme="light-transparent"] {
  /* 浅色透明主题变量 */
  --bg-primary: rgba(255, 255, 255, 0.65);
  --bg-secondary: rgba(248, 249, 250, 0.45);
  --bg-tertiary: rgba(240, 240, 240, 0.3);
  --bg-warning: rgba(248, 169, 41, 0.18);
  --bg-terminal: rgba(44, 62, 80, 0.7);
  --text-primary: #222;
  --text-secondary: #444;
  --text-muted: #888;
  --border-color: rgba(224, 224, 224, 0.5);
  --border-light: rgba(240, 240, 240, 0.3);
  --shadow-color: rgba(0, 0, 0, 0.08);
  --shadow-color-light: rgba(0, 0, 0, 0.04);
  --accent-color: #27ae60;
  --accent-hover: #229954;
  --accent-light: rgba(39, 174, 96, 0.12);
  --accent-gradient: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  --accent-gradient-hover: linear-gradient(135deg, #229954 0%, #27ae60 100%);
  --success-color: #27ae60;
  --warning-color: #ff8800;
  --error-color: #e74c3c;
  --info-color: #3498db;
}

[data-theme="dark-transparent"] {
  /* 深色透明主题变量 */
  --bg-primary: rgba(26, 26, 26, 0.65);
  --bg-secondary: rgba(45, 45, 45, 0.45);
  --bg-tertiary: rgba(61, 61, 61, 0.3);
  --bg-warning: rgba(243, 157, 18, 0.18);
  --bg-terminal: rgba(45, 45, 45, 0.7);
  --text-primary: #f3f3f3;
  --text-secondary: #b0b0b0;
  --text-muted: #888888;
  --border-color: rgba(64, 64, 64, 0.4);
  --border-light: rgba(80, 80, 80, 0.2);
  --shadow-color: rgba(0, 0, 0, 0.18);
  --shadow-color-light: rgba(0, 0, 0, 0.08);
  --accent-color: #4ade80;
  --accent-hover: #22c55e;
  --accent-light: rgba(74, 222, 128, 0.12);
  --accent-gradient: linear-gradient(135deg, #4ade80 0%, #22c55e 100%);
  --accent-gradient-hover: linear-gradient(135deg, #22c55e 0%, #16a34a 100%);
  --success-color: #4ade80;
  --warning-color: #ff7d32;
  --error-color: #f87171;
  --info-color: #60a5fa;
}

/* 保留全局样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
  line-height: 1.6;
  color: var(--text-primary);
  background: var(--bg-primary);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* 动画支持 */
@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

@-webkit-keyframes spin {
  0% {
    -webkit-transform: rotate(0deg);
  }
  100% {
    -webkit-transform: rotate(360deg);
  }
}

.spin {
  animation: spin 1s linear infinite !important;
  -webkit-animation: spin 1s linear infinite !important;
  display: inline-block;
}

:root {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
  font-size: 16px;
  line-height: 1.6;
  font-weight: 400;
  color: var(--text-primary);
  background: var(--bg-primary);
  font-synthesis: none;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  -webkit-text-size-adjust: 100%;
}

body {
  margin: 0;
  padding: 0;
  min-height: 100vh;
  background: var(--bg-primary);
}

h1, h2, h3, h4, h5, h6 {
  color: var(--text-primary);
  font-weight: 700;
  line-height: 1.3;
  margin: 0 0 1rem 0;
}

h1 {
  text-align: center;
  font-size: 2.5rem;
  font-weight: 800;
}

input,
button {
  border-radius: 8px;
  border: 2px solid transparent;
  padding: 0.8rem 1.2rem;
  font-size: 1rem;
  font-weight: 500;
  font-family: inherit;
  color: var(--text-primary);
  background: var(--bg-secondary);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px var(--shadow-color);
}

button {
  cursor: pointer;
  font-weight: 600;
  background: var(--accent-gradient);
  color: white;
  box-shadow: 0 4px 12px var(--accent-light);
}

button:hover {
  box-shadow: 0 6px 20px var(--accent-light);
}

button:active {
  box-shadow: 0 2px 8px var(--accent-light);
}

input:focus,
button:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: 0 0 0 3px var(--accent-light);
}

input::placeholder {
  color: var(--text-muted);
  opacity: 0.8;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: var(--bg-tertiary);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: var(--accent-gradient);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--accent-gradient-hover);
}

/* 选择文本样式 */
::selection {
  background: var(--accent-light);
  color: var(--text-primary);
}

/* 链接样式 */
a {
  color: var(--accent-color);
  text-decoration: none;
  transition: color 0.3s ease;
}

a:hover {
  color: var(--accent-hover);
}

/* 代码样式 */
code {
  background: var(--accent-light);
  color: var(--accent-color);
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
  font-size: 0.9em;
}

/* 表格样式 */
table {
  border-collapse: collapse;
  width: 100%;
  background: var(--bg-secondary);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px var(--shadow-color);
}

th, td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid var(--accent-light);
}

th {
  background: var(--accent-light);
  font-weight: 600;
  color: var(--text-primary);
}

tr:hover {
  background: var(--accent-light);
}
</style>