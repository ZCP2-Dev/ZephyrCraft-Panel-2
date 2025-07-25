<template>
  <div class="layout-root">
    <div class="sidebar">
      <div class="sidebar-header">
        <img src="../assets/logo.png" alt="logo" style="width: 60px; height: 60px;" />
        <h2 style="margin-top: 5px;">ZephyrCraft Panel</h2>
      </div>
      <div class="sidebar-menu">
  <template v-if="!isInConsole">
          <button @click="showServerList" class="menu-item" :class="{ 'active': currentView === 'serverList' }">
            <IconMdiHome style="margin-right: 10px; font-size: 1.2em;" />服务器列表
          </button>
          <button @click="showAppSettings" class="menu-item" :class="{ 'active': currentView === 'appSettings' }">
            <IconMdiCog style="margin-right: 10px; font-size: 1.2em;" />应用设置
          </button>
  </template>
  <template v-else>
          <button @click="showServerList" class="menu-item return-btn">
            <IconMdiHome style="margin-right: 10px; font-size: 1.2em;" />返回服务器列表
          </button>
    <div class="console-menu-header">控制台</div>
          <button v-for="item in consoleMenuItems" :key="item.key" @click="changeConsoleSection(item.key)" class="menu-item console-item" :class="{ 'active': consoleSection === item.key }">
            <component :is="iconMap[item.key]" style="margin-right: 10px; font-size: 1.2em;" />{{ item.label }}
          </button>
  </template>
</div>
    </div>
    <div class="main-content">
      <component :is="currentComponent" v-bind="consoleProps" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, defineAsyncComponent, watch, computed, provide } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useWebSocket } from '../useWebSocket';

const ServerList = defineAsyncComponent(() => import('../components/ServerList.vue'));
const AppSettings = defineAsyncComponent(() => import('../components/AppSettings.vue'));
const Overview = defineAsyncComponent(() => import('../components/console/Overview.vue'));
const Terminal = defineAsyncComponent(() => import('../components/console/Terminal.vue'));
const Players = defineAsyncComponent(() => import('../components/console/Players.vue'));
const Core = defineAsyncComponent(() => import('../components/console/Core.vue'));
const Plugins = defineAsyncComponent(() => import('../components/console/Plugins.vue'));
const Files = defineAsyncComponent(() => import('../components/console/Files.vue'));
const Other = defineAsyncComponent(() => import('../components/console/Other.vue'));

const route = useRoute();
const router = useRouter();
const currentView = ref('serverList');
const currentComponent = ref(ServerList);
const consoleSection = ref('overview');
const isInConsole = ref(false);
const currentServer = ref<any>(null);

const wsUrl = ref('');
const wsPassword = ref('');
const wsApi = useWebSocket({
  get url() { return wsUrl.value; },
  get password() { return wsPassword.value; },
});
provide('wsApi', wsApi);

const consoleMenuItems = [
  { key: 'overview', label: '概览' },
  { key: 'terminal', label: '终端' },
  { key: 'players', label: '玩家' },
  { key: 'core', label: '核心' },
  { key: 'plugins', label: '插件' },
  { key: 'files', label: '文件' },
  { key: 'other', label: '其它' }
];
const iconMap: Record<string, string> = {
  overview: 'IconMdiViewDashboard',
  terminal: 'IconMdiConsole',
  players: 'IconMdiAccountGroup',
  core: 'IconMdiChip',
  plugins: 'IconMdiPuzzle',
  files: 'IconMdiFileDocument',
  other: 'IconMdiDotsHorizontal'
};

watch(
  () => route.query,
  (newQuery) => {
    if (route.path === '/dashboard' && newQuery?.wsUrl) {
      // 从 query 组装 server
      const server = {
        name: newQuery.name as string || '',
        wsUrl: newQuery.wsUrl as string || '',
        password: newQuery.password as string || ''
      };
      enterConsole(server);
    } else if (route.path !== '/dashboard') {
      isInConsole.value = false;
      wsApi.disconnect();
    }
  },
  { immediate: true }
);

// wsApi.onMessage 全局分发终端消息
wsApi.onMessage = (data: any) => {
  if (window && (window as any).__TERMINAL_BUS__) {
    (window as any).__TERMINAL_BUS__.emit('terminal-message', data);
  }
};

function showServerList() {
  currentView.value = 'serverList';
  currentComponent.value = ServerList;
  isInConsole.value = false;
  if (route.path !== '/') {
    router.push('/');
  }
}

function showAppSettings() {
  currentView.value = 'appSettings';
  currentComponent.value = AppSettings;
  isInConsole.value = false;
}

function enterConsole(server: any) {
  currentServer.value = server;
  isInConsole.value = true;
  consoleSection.value = 'overview';
  currentComponent.value = Overview as unknown as any;
  wsUrl.value = server.wsUrl;
  wsPassword.value = server.password;
  wsApi.disconnect();
  // 重置控制台相关状态
  if (window && (window as any).__TERMINAL_BUS__) {
    (window as any).__TERMINAL_BUS__.emit('terminal-message', { status: 'stopped' });
  }
  wsApi.connect();
}

function changeConsoleSection(section: string) {
  consoleSection.value = section;
  switch(section) {
    case 'overview':
      currentComponent.value = Overview as unknown as any;
      break;
    case 'terminal':
      currentComponent.value = Terminal as unknown as any;
      break;
    case 'players':
      currentComponent.value = Players as unknown as any;
      break;
    case 'core':
      currentComponent.value = Core as unknown as any;
      break;
    case 'plugins':
      currentComponent.value = Plugins as unknown as any;
      break;
    case 'files':
      currentComponent.value = Files as unknown as any;
      break;
    case 'other':
      currentComponent.value = Other as unknown as any;
      break;
  }
}

const consoleProps = computed(() => {
  if (!isInConsole.value) return {};
  return { server: currentServer.value };
});
</script>

<style scoped>
.layout-root {
  display: flex;
  height: 100vh;
}
.sidebar {
  width: 220px;
  background: #fff;
  color: #23272e;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  border-radius: 18px;
  margin: 1.5rem 0 1.5rem 1.5rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  min-height: 0;
  padding: 10px;
}
.sidebar-header {
  padding: 2rem 1rem 1rem 1rem;
  text-align: center;
  border-bottom: 1px solid #eee;
}
.sidebar-menu {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 1rem 0;
}
.menu-item {
  background: none;
  border: none;
  color: #23272e;
  text-align: left;
  padding: 0.85rem 2.2rem;
  font-size: 1.08rem;
  font-weight: 600;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.18s, color 0.18s, box-shadow 0.18s;
  outline: none;
  box-shadow: none;
  display: flex;
  align-items: center;
  gap: 0.5em;
  margin-bottom: 0.5em;
}
.menu-item.active, .menu-item:hover {
  background: #eaf7ef;
  color: #3bb06c;
  box-shadow: 0 2px 8px rgba(56,191,100,0.08);
}
.return-btn {
  color: #3bb06c;
  font-weight: bold;
  background: #eaf7ef;
  border-radius: 12px;
  margin-bottom: 1.2rem;
}
.console-menu-header {
  margin: 1.2rem 0 0.5rem 2rem;
  font-size: 1.1rem;
  color: #88bf64;
  font-weight: bold;
}
.main-content {
  flex: 1;
  background: #f5f6fa;
  overflow: auto;
  padding: 2rem 2.5rem;
  border-radius: 18px;
  margin: 1.5rem 1.5rem 1.5rem 0;
}
</style>