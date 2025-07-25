<template>
  <div class="server-list-container">
    <div class="header-actions">
      <h2>服务器列表</h2>
      <button class="add-btn" @click="showAddServerModal = true">添加服务器</button>
    </div>

    <!-- 服务器卡片列表 -->
    <div class="cards-container">
      <div
        v-for="(server, index) in servers"
        :key="index"
        class="server-card"
        @click="goToConsole(server)"
        @contextmenu.prevent="showContextMenu($event, server, index)"
      >
        <h3>{{ server.name || '未命名服务器' }}</h3>
        <p>地址: {{ server.wsUrl }}</p>
        <p>状态: {{ server.status || '未知' }}</p>
        <p>CPU: {{ server.cpu || '0%' }}</p>
        <p>内存: {{ server.memory || '0%' }}</p>
      </div>
    </div>

    <!-- 右键菜单 -->
    <div
  v-if="showContextMenuStatus"
  ref="contextMenuRef"
  class="context-menu"
  :style="{ top: contextMenuY + 'px', left: contextMenuX + 'px' }"
  @click="hideContextMenu"
>
      <div v-if="contextMenuServer" @click.stop="goToConsole(contextMenuServer!)" class="menu-item">打开</div>
      <div @click.stop="deleteServer(contextMenuIndex)" class="menu-item">删除</div>
    </div>
  </div>

  <!-- 添加服务器弹窗 -->
  <div v-if="showAddServerModal" class="modal-overlay">
    <div class="modal">
      <h3>添加新服务器</h3>
      <div class="form-group">
        <label>服务器名称</label>
        <input v-model="newServer.name" type="text" placeholder="输入服务器名称">
      </div>
      <div class="form-group">
        <label>WS地址</label>
        <input v-model="newServer.wsUrl" type="text" placeholder="ws://example.com">
      </div>
      <div class="form-group">
        <label>连接密码</label>
        <input v-model="newServer.password" type="password" placeholder="输入连接密码">
      </div>
      <div class="modal-buttons">
        <button @click="showAddServerModal = false" class="cancel-btn">取消</button>
        <button @click="addServer" class="confirm-btn">添加</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';

interface Server {
  name: string;
  wsUrl: string;
  password: string;
  status?: string;
  cpu?: string;
  memory?: string;
}

// 路由导航
const router = useRouter();

// 服务器数据
const servers = ref<Server[]>([]);
const newServer = ref<Server>({ name: '', wsUrl: '', password: '' });

// 模态框状态
const showAddServerModal = ref(false);



// 右键菜单状态
const showContextMenuStatus = ref(false);
const contextMenuX = ref(0);
const contextMenuY = ref(0);
const contextMenuServer = ref<Server | null>(null);
const contextMenuIndex = ref(-1);
const contextMenuRef = ref<HTMLDivElement | null>(null);

// 从本地存储加载服务器
// 保存服务器到本地存储
const saveServersToLocalStorage = () => {
  localStorage.setItem('servers', JSON.stringify(servers.value));
};

// 添加服务器
const addServer = () => {
  if (!newServer.value.wsUrl || !newServer.value.password) return;
  
  // 简单验证WS地址格式
  if (!newServer.value.wsUrl.startsWith('ws://') && !newServer.value.wsUrl.startsWith('wss://')) {
    alert('请输入有效的WS地址（以ws://或wss://开头）');
    return;
  }

  servers.value.push({
    ...newServer.value,
    status: '未连接',
    cpu: '0%',
    memory: '0%'
  });

  saveServersToLocalStorage();
  newServer.value = { name: '', wsUrl: '', password: '' };
  showAddServerModal.value = false;
};

// 删除服务器
const deleteServer = (index: number) => {
  servers.value.splice(index, 1);
  saveServersToLocalStorage();
  hideContextMenu();
};

// 进入控制台
const goToConsole = (server: Server) => {
  router.push({
    path: '/dashboard',
    query: {
      name: server.name,
      wsUrl: server.wsUrl,
      password: server.password
    }
  });
};

// 显示右键菜单
const showContextMenu = (e: MouseEvent, server: Server, index: number) => {
  e.preventDefault();
  contextMenuX.value = e.clientX;
  contextMenuY.value = e.clientY;
  contextMenuServer.value = server;
  contextMenuIndex.value = index;
  showContextMenuStatus.value = true;
};

// 隐藏右键菜单
const hideContextMenu = () => {
  showContextMenuStatus.value = false;
};

const handleClickOutside = (event: MouseEvent) => {
  if (contextMenuRef.value && event.target instanceof Node && !contextMenuRef.value.contains(event.target)) {
    hideContextMenu();
  }
};

onMounted(() => {
  const savedServers = localStorage.getItem('servers');
  if (savedServers) {
    servers.value = JSON.parse(savedServers);
  }
  document.addEventListener('mousedown', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside);
});
</script>

<style scoped>
.server-list-container {
  padding: 0;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.add-btn {
  background-color: #88bf64;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.5rem 1rem;
  cursor: pointer;
}

.add-btn:hover {
  background-color: #539951;
}

.cards-container {
  display: flex;
  flex-wrap: wrap;
  gap: 1.5rem;
}

.cards-container p {
  margin: 0;
}

.server-card {
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  width: 300px;
  margin-bottom: 1rem;
  cursor: pointer;
  transition: transform 0.2s;
}

.server-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
}

.modal {
  background-color: white;
  padding: 2rem;
  border-radius: 8px;
  width: 400px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
}

.form-group input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.modal-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.cancel-btn {
  background-color: #ddd;
  border: none;
  border-radius: 4px;
  padding: 0.5rem 1rem;
  cursor: pointer;
}

.confirm-btn {
  background-color: #88bf64;
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.5rem 1rem;
  cursor: pointer;
}

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  /* 显示状态由v-if控制 */
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 1001;
  min-width: 120px;
}

.context-menu .menu-item {
  padding: 0.5rem 1rem;
  cursor: pointer;
}

.context-menu .menu-item:hover {
  background-color: #f5f5f5;
}
</style>