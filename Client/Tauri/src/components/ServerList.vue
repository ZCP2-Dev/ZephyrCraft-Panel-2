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
  background: #ffffff;
  padding: 1.5rem;
  box-sizing: border-box;
}

.server-list-container h2 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  box-sizing: border-box;
}

.add-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.8rem 1.5rem;
  cursor: pointer;
  font-weight: 600;
  font-size: 1rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
  position: relative;
  overflow: hidden;
  box-sizing: border-box;
}

.add-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.add-btn:hover::before {
  left: 100%;
}

.add-btn:hover {
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

.cards-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
  box-sizing: border-box;
}

.cards-container p {
  margin: 0;
}

.server-card {
  border: 1px solid #e9ecef;
  border-radius: 10px;
  padding: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  width: 100%;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: #ffffff;
  position: relative;
  overflow: hidden;
  box-sizing: border-box;
}

.server-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.server-card:hover::before {
  transform: scaleX(1);
}

.server-card:hover {
  box-shadow: 0 6px 16px rgba(39, 174, 96, 0.15);
  border-color: rgba(39, 174, 96, 0.3);
}

.server-card h3 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.3rem;
  margin: 0 0 1rem 0;
}

.server-card p {
  color: #7f8c8d;
  font-size: 0.95rem;
  margin: 0.5rem 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 9999;
  box-sizing: border-box;
  padding: 1rem;
}

.modal {
  background: #ffffff;
  padding: 2rem;
  border-radius: 12px;
  width: 450px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  border: 1px solid #e9ecef;
  box-sizing: border-box;
  max-width: calc(100vw - 2rem);
  max-height: calc(100vh - 2rem);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.modal h3 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.5rem;
  margin: 0 0 2rem 0;
  text-align: center;
  flex-shrink: 0;
}

.form-group {
  margin-bottom: 1.5rem;
  box-sizing: border-box;
  flex-shrink: 0;
}

.form-group label {
  display: block;
  margin-bottom: 0.8rem;
  color: #2c3e50;
  font-weight: 600;
  font-size: 0.95rem;
}

.form-group input {
  width: 100%;
  padding: 1rem;
  border: 2px solid #e9ecef;
  border-radius: 8px;
  font-size: 1rem;
  transition: all 0.3s ease;
  background: #ffffff;
  box-sizing: border-box;
  min-width: 0;
}

.form-group input:focus {
  outline: none;
  border-color: #27ae60;
  box-shadow: 0 0 0 3px rgba(39, 174, 96, 0.1);
  background: #ffffff;
}

.modal-buttons {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2.5rem;
  box-sizing: border-box;
  flex-shrink: 0;
}

.cancel-btn {
  background: #bdc3c7;
  color: #2c3e50;
  border: none;
  border-radius: 8px;
  padding: 0.8rem 1.5rem;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s ease;
  box-sizing: border-box;
  white-space: nowrap;
}

.cancel-btn:hover {
  background: #95a5a6;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.confirm-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.8rem 1.5rem;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
  box-sizing: border-box;
  white-space: nowrap;
}

.confirm-btn:hover {
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

/* 响应式设计 */
@media (max-width: 480px) {
  .modal {
    padding: 1.5rem;
    width: calc(100vw - 2rem);
  }
  
  .modal h3 {
    font-size: 1.3rem;
    margin-bottom: 1.5rem;
  }
  
  .form-group {
    margin-bottom: 1.2rem;
  }
  
  .form-group input {
    padding: 0.8rem;
    font-size: 0.95rem;
  }
  
  .modal-buttons {
    flex-direction: column;
    gap: 0.8rem;
  }
  
  .cancel-btn, .confirm-btn {
    width: 100%;
    padding: 0.8rem;
  }
}

@media (max-height: 600px) {
  .modal {
    max-height: calc(100vh - 1rem);
    padding: 1.5rem;
  }
  
  .modal h3 {
    margin-bottom: 1.5rem;
  }
  
  .form-group {
    margin-bottom: 1rem;
  }
  
  .modal-buttons {
    margin-top: 1.5rem;
  }
}

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  background: #ffffff;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  z-index: 1001;
  min-width: 140px;
  overflow: hidden;
  box-sizing: border-box;
}

.context-menu .menu-item {
  padding: 0.8rem 1.2rem;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #2c3e50;
  font-weight: 500;
  box-sizing: border-box;
}

.context-menu .menu-item:hover {
  background: linear-gradient(135deg, rgba(39, 174, 96, 0.1) 0%, rgba(46, 204, 113, 0.1) 100%);
  color: #27ae60;
}
</style>