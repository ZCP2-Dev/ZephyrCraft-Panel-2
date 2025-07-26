<template>
  <div class="plugins-container">
    <h2>插件管理</h2>
    <div class="plugins-list">
      <h3>已安装插件</h3>
      <div v-for="plugin in plugins" :key="plugin.name" class="plugin-item">
        <div class="plugin-info">
          <div class="plugin-icon">{{ plugin.name.charAt(0) }}</div>
          <div class="plugin-details">
            <span class="plugin-name">{{ plugin.name }}</span>
            <span class="plugin-version">v{{ plugin.version }}</span>
          </div>
        </div>
        <div class="plugin-actions">
          <button v-if="plugin.enabled" @click="disablePlugin(plugin.name)" class="plugin-btn disable-btn">
            <IconMdiPowerOff style="margin-right: 5px;" />禁用
          </button>
          <button v-else @click="enablePlugin(plugin.name)" class="plugin-btn enable-btn">
            <IconMdiPowerOn style="margin-right: 5px;" />启用
          </button>
          <button @click="reloadPlugin(plugin.name)" class="plugin-btn reload-btn">
            <IconMdiReload style="margin-right: 5px;" />重载
          </button>
        </div>
      </div>
      <div v-if="plugins.length === 0" class="no-plugins">
        暂无已安装插件
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const plugins = ref([
  { name: 'Essentials', version: '2.19.0', enabled: true },
  { name: 'WorldEdit', version: '7.2.5', enabled: true },
  { name: 'WorldGuard', version: '7.0.7', enabled: false },
  { name: 'LuckPerms', version: '5.3.0', enabled: true }
]);

function enablePlugin(pluginName: string) {
  console.log('启用插件:', pluginName);
}

function disablePlugin(pluginName: string) {
  console.log('禁用插件:', pluginName);
}

function reloadPlugin(pluginName: string) {
  console.log('重载插件:', pluginName);
}
</script>

<style scoped>
.plugins-container {
  /* background: var(--bg-primary); */
  padding: 1.5rem;
}

.plugins-container h2 {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.plugins-container h3 {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 1.4rem;
  margin: 0 0 1.2rem 0;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.plugins-list {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.plugin-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.2rem;
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.plugin-item:hover {
  background: var(--bg-primary);
  box-shadow: 0 2px 8px var(--accent-light);
  border-color: var(--accent-color);
}

.plugin-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.plugin-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: var(--accent-gradient);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 1.2rem;
}

.plugin-details {
  display: flex;
  flex-direction: column;
}

.plugin-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 1.1rem;
}

.plugin-version {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.plugin-actions {
  display: flex;
  gap: 0.6rem;
}

.plugin-btn {
  padding: 0.6rem 1rem;
  border: none;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.plugin-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--shadow-color-light);
}

.plugin-btn:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
  box-shadow: 0 2px 8px var(--shadow-color);
}

.enable-btn {
  background: var(--accent-gradient);
  color: #fff;
  box-shadow: 0 2px 8px var(--accent-light);
}

.enable-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--accent-light);
}

.disable-btn {
  background: var(--error-color);
  color: #fff;
  box-shadow: 0 2px 8px var(--error-color);
}

.disable-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--error-color);
}

.reload-btn {
  background: var(--warning-color);
  color: #fff;
  box-shadow: 0 2px 8px var(--warning-color);
}

.reload-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px var(--warning-color);
}

.no-plugins {
  text-align: center;
  padding: 3rem 2rem;
  color: var(--text-secondary);
  font-size: 1.1rem;
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}
</style>