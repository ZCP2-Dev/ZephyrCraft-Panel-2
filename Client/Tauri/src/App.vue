<template>
  <div id="app">
    
    <router-view class="view-content" />
  </div>
</template>

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
/* 保留全局样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
  line-height: 1.6;
  color: #2c3e50;
  background: #ffffff;
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
  color: #2c3e50;
  background: #ffffff;
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
  background: #ffffff;
}

h1, h2, h3, h4, h5, h6 {
  color: #2c3e50;
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
  color: #2c3e50;
  background: rgba(255, 255, 255, 0.9);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

button {
  cursor: pointer;
  font-weight: 600;
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.3);
}

button:hover {
  box-shadow: 0 6px 20px rgba(39, 174, 96, 0.4);
}

button:active {
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
}

input:focus,
button:focus {
  outline: none;
  border-color: #27ae60;
  box-shadow: 0 0 0 3px rgba(39, 174, 96, 0.1);
}

input::placeholder {
  color: #95a5a6;
  opacity: 0.8;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #229954 0%, #27ae60 100%);
}

/* 选择文本样式 */
::selection {
  background: rgba(39, 174, 96, 0.3);
  color: #2c3e50;
}

/* 链接样式 */
a {
  color: #27ae60;
  text-decoration: none;
  transition: color 0.3s ease;
}

a:hover {
  color: #229954;
}

/* 代码样式 */
code {
  background: rgba(39, 174, 96, 0.1);
  color: #27ae60;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-family: 'Fira Code', 'Monaco', 'Consolas', monospace;
  font-size: 0.9em;
}

/* 表格样式 */
table {
  border-collapse: collapse;
  width: 100%;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

th, td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid rgba(39, 174, 96, 0.1);
}

th {
  background: linear-gradient(135deg, rgba(39, 174, 96, 0.1) 0%, rgba(46, 204, 113, 0.1) 100%);
  font-weight: 600;
  color: #2c3e50;
}

tr:hover {
  background: rgba(39, 174, 96, 0.05);
}
</style>