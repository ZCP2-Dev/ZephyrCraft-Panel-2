<template>
  <div class="files-container">
    <h3>文件管理</h3>
      <div class="files-list">
      <div v-for="file in files" :key="file.name" class="file-item">
        <span class="file-icon">{{ file.type === 'folder' ? '📁' : '📄' }}</span>
        <span class="file-name">{{ file.name }}</span>
        <span class="file-size" v-if="file.type === 'file'">{{ file.size }}</span>
        <span class="file-actions">
          <button class="file-btn">下载</button>
          <button class="file-btn">重命名</button>
          <button class="file-btn">删除</button>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const files = [
  { name: 'world', type: 'folder' },
  { name: 'server.properties', type: 'file', size: '2 KB' },
  { name: 'ops.json', type: 'file', size: '1 KB' },
  { name: 'plugins', type: 'folder' },
  { name: 'logs', type: 'folder' },
  { name: 'banned-players.json', type: 'file', size: '0.5 KB' }
];
</script>

<style scoped>
.files-container {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  padding: 2rem 2.5rem;
}
.files-list {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.file-item {
  display: flex;
  align-items: center;
  gap: 1.2rem;
  padding: 0.8rem 1.2rem;
  background: #f5f6fa;
  border-radius: 6px;
}
.file-icon {
  font-size: 1.3rem;
}
.file-name {
  font-weight: bold;
  color: #23272e;
  flex: 1;
}
.file-size {
  color: #888;
  font-size: 0.98rem;
}
.file-actions {
  display: flex;
  gap: 0.8rem;
}
.file-btn {
  background: #88bf64;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 0.3rem 1rem;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background 0.2s;
}
.file-btn:hover {
  background: #539951;
}
</style>