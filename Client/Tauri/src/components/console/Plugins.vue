<template>
  <div class="plugins-container">
    <h3>插件管理</h3>
    <div class="plugins-list">
      <div v-for="plugin in plugins" :key="plugin.name" class="plugin-item">
        <div class="plugin-info">
          <span class="plugin-name">{{ plugin.name }}</span>
          <span class="plugin-version">v{{ plugin.version }}</span>
          <span class="plugin-status" :class="plugin.enabled ? 'enabled' : 'disabled'">
            {{ plugin.enabled ? '已启用' : '已禁用' }}
          </span>
        </div>
        <div class="plugin-actions">
          <button class="plugin-btn">重载</button>
          <button class="plugin-btn">{{ plugin.enabled ? '禁用' : '启用' }}</button>
          <button class="plugin-btn">删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const plugins = [
  { name: 'EssentialsX', version: '2.19.7', enabled: true },
  { name: 'LuckPerms', version: '5.4.120', enabled: true },
  { name: 'Dynmap', version: '3.4-beta-5', enabled: false },
  { name: 'WorldEdit', version: '7.2.12', enabled: true },
  { name: 'Vault', version: '1.7.3', enabled: true }
];
</script>

<style scoped>
.plugins-container {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  padding: 2rem 2.5rem;
}
.plugins-list {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.plugin-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.8rem 1.2rem;
  background: #f5f6fa;
  border-radius: 6px;
}
.plugin-info {
  display: flex;
  align-items: center;
  gap: 1.2rem;
}
.plugin-name {
  font-weight: bold;
  color: #23272e;
}
.plugin-version {
  color: #888;
  font-size: 0.98rem;
}
.plugin-status {
  font-size: 0.95rem;
  font-weight: bold;
}
.enabled {
  color: #88bf64;
}
.disabled {
  color: #ff6b6b;
}
.plugin-actions {
  display: flex;
  gap: 0.8rem;
}
.plugin-btn {
  background: #88bf64;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 0.3rem 1rem;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background 0.2s;
}
.plugin-btn:hover {
  background: #539951;
}
</style>