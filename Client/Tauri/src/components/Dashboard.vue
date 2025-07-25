<template>
  <div class="dashboard-container">
    <div class="server-header">
      <h2>{{ currentServer?.name || '服务器控制台' }}</h2>
      <div class="server-status">{{ currentServer?.status || '未连接' }}</div>
    </div>
    <component :is="currentConsoleComponent" :server="currentServer"></component>
  </div>
</template>

<script setup lang="ts">
import { ref, defineProps } from 'vue';
import Overview from './console/Overview.vue';
import Terminal from './console/Terminal.vue';
import Players from './console/Players.vue';
import Core from './console/Core.vue';
import Plugins from './console/Plugins.vue';
import Files from './console/Files.vue';
import Other from './console/Other.vue';

// 定义props
const props = defineProps<{
  server?: any,
  section?: string
}>();

// 当前控制台组件
const currentConsoleComponent = ref<any>(Overview);
const currentServer = ref<any>(props.server);

// 根据传入的section切换组件
if (props.section) {
  switch(props.section) {
    case 'overview':
      currentConsoleComponent.value = Overview;
      break;
    case 'terminal':
      currentConsoleComponent.value = Terminal;
      break;
    case 'players':
      currentConsoleComponent.value = Players;
      break;
    case 'core':
      currentConsoleComponent.value = Core;
      break;
    case 'plugins':
      currentConsoleComponent.value = Plugins;
      break;
    case 'files':
      currentConsoleComponent.value = Files;
      break;
    case 'other':
      currentConsoleComponent.value = Other;
      break;
  }
}
</script>

<style scoped>
.dashboard-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.server-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
  margin-bottom: 1rem;
  border-bottom: 1px solid #e0e0e0;
}

.server-status {
  background-color: #88bf64;
  color: white;
  padding: 0.3rem 0.8rem;
  border-radius: 4px;
  font-size: 0.9rem;
}
</style>