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
  background: #ffffff;
  padding: 1.5rem;
}

.plugins-container h2 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.plugins-container h3 {
  color: #2c3e50;
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
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.plugin-item:hover {
  background: #ffffff;
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.1);
  border-color: rgba(39, 174, 96, 0.2);
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
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
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
  color: #2c3e50;
  font-size: 1.1rem;
}

.plugin-version {
  font-size: 0.9rem;
  color: #7f8c8d;
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
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.plugin-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.plugin-btn:disabled {
  background: linear-gradient(135deg, #ecf0f1 0%, #bdc3c7 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.enable-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
}

.enable-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

.disable-btn {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(231, 76, 60, 0.3);
}

.disable-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(231, 76, 60, 0.4);
}

.reload-btn {
  background: linear-gradient(135deg, #f39c12 0%, #e67e22 100%);
  color: #fff;
  box-shadow: 0 2px 8px rgba(243, 156, 18, 0.3);
}

.reload-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(243, 156, 18, 0.4);
}

.no-plugins {
  text-align: center;
  padding: 3rem 2rem;
  color: #7f8c8d;
  font-size: 1.1rem;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}
</style>