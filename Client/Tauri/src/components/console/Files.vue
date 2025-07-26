<template>
  <div class="files-container">
    <h2>文件管理</h2>
    
    <!-- 路径导航栏 -->
    <div class="path-navigation">
      <div class="path-breadcrumb">
        <span
          class="path-segment"
          :class="{ current: currentPath === '' }"
          @click="navigateToPath('')"
        >根目录</span>
        <template v-for="(segment, index) in pathSegments" :key="index">
          <span class="separator">/</span>
          <span
            class="path-segment"
            :class="{ current: index === pathSegments.length - 1 && currentPath !== '' }"
            @click="navigateToPath(getPathUpTo(index))"
          >{{ segment }}</span>
        </template>
      </div>
      <div class="path-actions">
        <button @click="refreshFiles" class="action-btn" :disabled="loading">
          <IconMdiRefresh style="margin-right: 5px;" />刷新
        </button>
        <button @click="createFolder" class="action-btn">
          <IconMdiFolderPlus style="margin-right: 5px;" />新建文件夹
        </button>
        <button @click="uploadFile" class="action-btn">
          <IconMdiUpload style="margin-right: 5px;" />上传文件
        </button>
      </div>
    </div>

    <!-- 文件列表 -->
    <div class="files-section">
      <div v-if="loading" class="loading">
        <div class="spinner"></div>
        <span>加载中...</span>
      </div>
      
      <div v-else-if="error" class="error-message">
        <IconMdiAlertCircle style="margin-right: 8px;" />
        {{ error }}
      </div>
      
      <div v-else class="files-list">
        <div 
          v-for="file in files" 
          :key="file.path"
          class="file-item"
          :class="{ 'is-directory': file.isDir, 'is-editable': file.isEditable }"
          @click="handleFileClick(file)"
          @contextmenu.prevent="showContextMenu($event, file)"
          @dblclick="handleFileDoubleClick(file)"
        >
          <div class="file-icon">
            <IconMdiFolder v-if="file.isDir" />
            <IconMdiFileDocument v-else-if="file.isEditable" />
            <IconMdiFile v-else />
          </div>
          <div class="file-info">
            <div class="file-name">{{ file.name }}</div>
            <div class="file-details">
              <span v-if="!file.isDir">{{ formatFileSize(file.size) }}</span>
              <span>{{ formatDate(file.modTime) }}</span>
            </div>
          </div>
          <div class="file-actions">
            <button 
              v-if="file.isEditable && !file.isDir && file.path" 
              @click.stop="editFile(file)"
              class="edit-btn"
              title="编辑文件"
            >
              <IconMdiPencil />
        </button>
            <button 
              v-if="!file.isDir && file.path" 
              @click.stop="downloadFile(file)"
              class="download-btn"
              title="下载文件"
            >
              <IconMdiDownload />
        </button>
          </div>
        </div>
        
        <div v-if="files.length === 0" class="no-files">
          <IconMdiFolderOpen style="font-size: 3rem; color: #bdc3c7; margin-bottom: 1rem;" />
          <p>当前目录为空</p>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <div 
      v-if="contextMenu.show" 
      class="context-menu"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
      @click.stop
    >
      <div 
        v-if="contextMenu.file.isDir"
        class="context-menu-item"
        @click="navigateToFile(contextMenu.file)"
      >
        <IconMdiFolderOpen style="margin-right: 8px;" />
        打开
      </div>
      <div 
        v-if="contextMenu.file && contextMenu.file.isEditable && !contextMenu.file.isDir && contextMenu.file.path"
        class="context-menu-item"
        @click="editFile(contextMenu.file)"
      >
        <IconMdiPencil style="margin-right: 8px;" />
        编辑
      </div>
      <div 
        class="context-menu-item"
        @click="downloadFile(contextMenu.file)"
      >
        <IconMdiDownload style="margin-right: 8px;" />
        下载
      </div>
      <div class="context-menu-separator"></div>
      <div 
        class="context-menu-item"
        @click="renameFile(contextMenu.file)"
      >
        <IconMdiRename style="margin-right: 8px;" />
        重命名
      </div>
      <div 
        class="context-menu-item danger"
        @click="deleteFile(contextMenu.file)"
      >
        <IconMdiDelete style="margin-right: 8px;" />
        删除
      </div>
    </div>
    
    <!-- 文件编辑器对话框 -->
    <div v-if="editor.show" class="editor-modal">
      <div class="editor-container">
        <div class="editor-header">
          <h3>编辑文件: {{ editor.fileName }}</h3>
          <div class="editor-actions">
            <button @click="saveFile" class="save-btn" :disabled="editor.saving">
              <IconMdiContentSave style="margin-right: 5px;" />
              {{ editor.saving ? '保存中...' : '保存' }}
            </button>
            <button @click="closeEditor" class="close-btn">
              <IconMdiClose />
            </button>
          </div>
        </div>
        <div class="editor-content">
          <textarea
            v-model="editor.content"
            class="file-editor"
            placeholder="文件内容..."
            spellcheck="false"
          ></textarea>
        </div>
      </div>
    </div>

    <!-- 重命名对话框 -->
    <div v-if="renameDialog.show" class="modal-overlay">
      <div class="modal">
        <h3>重命名</h3>
        <input
          v-model="renameDialog.newName"
          type="text"
          class="rename-input"
          placeholder="输入新名称"
          @keyup.enter="confirmRename"
          @keyup.esc="cancelRename"
        />
        <div class="modal-actions">
          <button @click="confirmRename" class="confirm-btn">确认</button>
          <button @click="cancelRename" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>

    <!-- 新建文件夹对话框 -->
    <div v-if="createFolderDialog.show" class="modal-overlay">
      <div class="modal">
        <h3>新建文件夹</h3>
        <input
          v-model="createFolderDialog.name"
          type="text"
          class="folder-input"
          placeholder="输入文件夹名称"
          @keyup.enter="confirmCreateFolder"
          @keyup.esc="cancelCreateFolder"
        />
        <div class="modal-actions">
          <button @click="confirmCreateFolder" class="confirm-btn">创建</button>
          <button @click="cancelCreateFolder" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>

    <!-- 隐藏的文件上传输入 -->
    <input
      ref="fileInput"
      type="file"
      multiple
      style="display: none"
      @change="handleFileUpload"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, onUnmounted, watch } from 'vue';

const props = defineProps<{ server?: any }>();
const wsApi = inject('wsApi') as any;
const isConnected = computed(() => wsApi?.isConnected && typeof wsApi.isConnected === 'object' ? wsApi.isConnected.value : wsApi.isConnected);

// 文件列表状态
const files = ref<Array<{
  name: string;
  path: string;
  isDir: boolean;
  size: number;
  modTime: string;
  permissions: string;
  isEditable: boolean;
}>>([]);
const currentPath = ref('');
const loading = ref(false);
const error = ref('');

// 右键菜单状态
const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  file: null as any
});

// 编辑器状态
const editor = ref({
  show: false,
  fileName: '',
  filePath: '',
  content: '',
  saving: false
});

// 重命名对话框状态
const renameDialog = ref({
  show: false,
  file: null as any,
  newName: ''
});

// 新建文件夹对话框状态
const createFolderDialog = ref({
  show: false,
  name: ''
});

// 文件上传输入引用
const fileInput = ref<HTMLInputElement>();

// 计算路径分段
const pathSegments = computed(() => {
  if (!currentPath.value) return [];
  return currentPath.value.split('/').filter(segment => segment);
});

// 格式化文件大小
function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 格式化日期
function formatDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleString('zh-CN');
}

// 获取指定索引的路径
function getPathUpTo(index: number): string {
  return pathSegments.value.slice(0, index + 1).join('/');
}

// 导航到指定路径
function navigateToPath(path: string) {
  currentPath.value = path;
  loadFiles();
}

// 加载文件列表
function loadFiles() {
  if (!wsApi || !isConnected.value) {
    error.value = 'WebSocket连接未建立';
    return;
  }

  loading.value = true;
  error.value = '';
  
  wsApi.send({ command: 'getFileList', content: currentPath.value });
}

// 刷新文件列表
function refreshFiles() {
  loadFiles();
}

// 处理文件点击
function handleFileClick(file: any) {
  if (file.isDir) {
    navigateToFile(file);
  }
}

// 处理文件双击
function handleFileDoubleClick(file: any) {
  if (file.isDir) {
    navigateToFile(file);
  } else if (file.isEditable) {
    editFile(file);
  }
}

// 导航到文件/文件夹
function navigateToFile(file: any) {
  if (file.isDir) {
    const newPath = currentPath.value ? `${currentPath.value}/${file.name}` : file.name;
    navigateToPath(newPath);
  }
}

// 显示右键菜单
function showContextMenu(event: MouseEvent, file: any) {
  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    file
  };
}

// 隐藏右键菜单
function hideContextMenu() {
  contextMenu.value.show = false;
}

// 编辑文件
function editFile(file: any) {

  console.log(file);

  if (!file || !file.isEditable || file.isDir || !file.path) return;
  
  loading.value = true;
  wsApi.send({ command: 'readFile', filePath: file.path });
}

// 保存文件
function saveFile() {
  if (!editor.value.content || editor.value.saving) return;
  
  editor.value.saving = true;
  wsApi.send({
    command: 'writeFile',
    filePath: editor.value.filePath,
    fileContent: editor.value.content
  });
}

// 关闭编辑器
function closeEditor() {
  editor.value.show = false;
  editor.value.fileName = '';
  editor.value.filePath = '';
  editor.value.content = '';
  editor.value.saving = false;
}

// 下载文件
function downloadFile(file: any) {
  if (file.isDir || !file.path) return;
  wsApi.send({ command: 'downloadFile', filePath: file.path });
}

// 重命名文件
function renameFile(file: any) {
  renameDialog.value = {
    show: true,
    file,
    newName: file.name
  };
  hideContextMenu();
}

// 确认重命名
function confirmRename() {
  const { file, newName } = renameDialog.value;
  if (!newName.trim()) return;
  
  const oldPath = file.path;
  const newPath = file.path.replace(file.name, newName);
  
  wsApi.send({
    command: 'renameFile',
    oldPath,
    newPath
  });
  
  cancelRename();
}

// 取消重命名
function cancelRename() {
  renameDialog.value.show = false;
  renameDialog.value.file = null;
  renameDialog.value.newName = '';
}

// 删除文件
function deleteFile(file: any) {
  if (!confirm(`确定要删除 ${file.name} 吗？此操作不可恢复！`)) return;
  
  wsApi.send({
    command: 'deleteFile',
    content: file.path
  });
  
  hideContextMenu();
}

// 创建文件夹
function createFolder() {
  createFolderDialog.value.show = true;
  createFolderDialog.value.name = '';
}

// 确认创建文件夹
function confirmCreateFolder() {
  const name = createFolderDialog.value.name.trim();
  if (!name) return;
  
  const folderPath = currentPath.value ? `${currentPath.value}/${name}` : name;
  wsApi.send({
    command: 'createDirectory',
    content: folderPath
  });
  
  cancelCreateFolder();
}

// 取消创建文件夹
function cancelCreateFolder() {
  createFolderDialog.value.show = false;
  createFolderDialog.value.name = '';
}

// 上传文件
function uploadFile() {
  fileInput.value?.click();
}

// 处理文件上传
function handleFileUpload(event: Event) {
  const target = event.target as HTMLInputElement;
  const files = target.files;
  if (!files || files.length === 0) return;
  const file = files[0];
  const reader = new FileReader();
  reader.onload = function(e) {
    const base64 = (e.target?.result as string).split(',')[1];
    wsApi.send({ command: 'uploadFile', filePath: (currentPath.value ? currentPath.value + '/' : '') + file.name, fileContent: base64 });
  };
  reader.readAsDataURL(file);
  target.value = '';
}

// 监听WebSocket消息，支持下载
function handleWebSocketMessage(data: any) {
  console.log('Files received WebSocket message:', data); // 调试日志
  
  // 处理文件列表响应
  if (data.fileList) {
    files.value = data.fileList;
    loading.value = false;
    error.value = '';
  } 
  // 处理文件内容读取响应
  else if (data.fileContent !== undefined) {
    editor.value = {
      show: true,
      fileName: data.filePath ? data.filePath.split('/').pop() : '',
      filePath: data.filePath || '',
      content: data.fileContent,
      saving: false
    };
    loading.value = false;
    error.value = '';
  } 
  // 处理文件写入响应
  else if (data.status === 'success' && data.filePath) {
    // 这是文件写入成功的响应
    closeEditor();
    loadFiles(); // 刷新文件列表
  }
  // 处理其他文件操作响应（重命名、删除、创建文件夹等）
  else if (data.status === 'success') {
    // 操作成功，刷新文件列表
    loadFiles();
  } 
  // 处理错误响应
  else if (data.error) {
    error.value = data.error;
    loading.value = false;
    editor.value.saving = false;
  }
  // 下载文件响应
  if (data.status === 'download' && data.fileContent && data.filePath) {
    // base64转Blob并下载
    const bstr = atob(data.fileContent);
    const u8arr = new Uint8Array(bstr.length);
    for (let i = 0; i < bstr.length; ++i) u8arr[i] = bstr.charCodeAt(i);
    const blob = new Blob([u8arr]);
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = data.filePath.split('/').pop() || 'download';
    document.body.appendChild(a);
    a.click();
    setTimeout(() => {
      document.body.removeChild(a);
      URL.revokeObjectURL(a.href);
    }, 100);
    return;
  }
}

// 监听点击事件，隐藏右键菜单
function handleGlobalClick() {
  hideContextMenu();
}

onMounted(() => {
  // 添加全局点击监听
  document.addEventListener('click', handleGlobalClick);
  
  // 加载初始文件列表
  loadFiles();
  
  // 监听全局消息总线中的文件相关消息
  if (window && (window as any).__SYSTEM_BUS__) {
    (window as any).__SYSTEM_BUS__.on('file-message', handleWebSocketMessage);
  }
});

onUnmounted(() => {
  document.removeEventListener('click', handleGlobalClick);
  
  // 清理消息监听
  if (window && (window as any).__SYSTEM_BUS__) {
    (window as any).__SYSTEM_BUS__.off('file-message', handleWebSocketMessage);
  }
});

// 监听连接状态变化
watch(isConnected, (connected) => {
  if (connected) {
    loadFiles();
  }
});
</script>

<style scoped>
.files-container {
  background: #ffffff;
  padding: 1.5rem;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.files-container h2 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.path-navigation {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
  border: 1px solid #e9ecef;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.path-breadcrumb {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.path-segment {
  color: #3498db;
  cursor: pointer;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  transition: all 0.2s;
  font-weight: 500;
}

.path-segment:hover {
  background: rgba(52, 152, 219, 0.1);
}

.path-segment.current {
  color: #2c3e50;
  font-weight: 600;
  cursor: default;
}

.path-segment.current:hover {
  background: none;
}

.path-actions {
  display: flex;
  gap: 0.8rem;
}

.action-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
}

.action-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

.action-btn:disabled {
  background: linear-gradient(135deg, #ecf0f1 0%, #bdc3c7 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.files-section {
  flex: 1;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  color: #7f8c8d;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #27ae60;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  color: #e74c3c;
  background: #fdf2f2;
  border-radius: 8px;
  margin: 1rem;
}

.files-list {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 0.8rem 1rem;
  background: #ffffff;
  border-radius: 6px;
  border: 1px solid #e9ecef;
  margin-bottom: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.file-item:hover {
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.1);
  border-color: rgba(39, 174, 96, 0.2);
}

.file-item.is-directory {
  background: linear-gradient(135deg, #f8f9fa 0%, #e9ecef 100%);
}

.file-item.is-editable {
  border-left: 3px solid #3498db;
}

.file-icon {
  margin-right: 1rem;
  font-size: 1.5rem;
  color: #7f8c8d;
  width: 24px;
  text-align: center;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 0.2rem;
  word-break: break-all;
}

.file-details {
  font-size: 0.8rem;
  color: #7f8c8d;
  display: flex;
  gap: 1rem;
}

.file-actions {
  display: flex;
  gap: 0.5rem;
  opacity: 0;
  transition: opacity 0.2s;
}

.file-item:hover .file-actions {
  opacity: 1;
}

.edit-btn, .download-btn {
  background: none;
  border: none;
  padding: 0.3rem;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  color: #7f8c8d;
}

.edit-btn:hover {
  background: rgba(52, 152, 219, 0.1);
  color: #3498db;
}

.download-btn:hover {
  background: rgba(39, 174, 96, 0.1);
  color: #27ae60;
}

.no-files {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  color: #7f8c8d;
  text-align: center;
}

.context-menu {
  position: fixed;
  background: #ffffff;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  min-width: 150px;
  padding: 0.5rem 0;
}

.context-menu-item {
  padding: 0.6rem 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: background 0.2s;
  font-size: 0.9rem;
}

.context-menu-item:hover {
  background: #f8f9fa;
}

.context-menu-item.danger:hover {
  background: #fdf2f2;
  color: #e74c3c;
}

.context-menu-separator {
  height: 1px;
  background: #e9ecef;
  margin: 0.5rem 0;
}

.editor-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.editor-container {
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  width: 100%;
  max-width: 800px;
  height: 80vh;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e9ecef;
  background: #f8f9fa;
  border-radius: 8px 8px 0 0;
}

.editor-header h3 {
  margin: 0;
  color: #2c3e50;
  font-weight: 600;
}

.editor-actions {
  display: flex;
  gap: 0.8rem;
}

.save-btn, .close-btn {
  border: none;
  border-radius: 4px;
  padding: 0.5rem 1rem;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
  display: flex;
  align-items: center;
}

.save-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: #fff;
}

.save-btn:hover:not(:disabled) {
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
}

.save-btn:disabled {
  background: #bdc3c7;
  cursor: not-allowed;
}

.close-btn {
  background: #e74c3c;
  color: #fff;
}

.close-btn:hover {
  background: #c0392b;
}

.editor-content {
  flex: 1;
  padding: 1rem;
}

.file-editor {
  width: 100%;
  height: 100%;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  padding: 1rem;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  resize: none;
  outline: none;
}

.file-editor:focus {
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1500;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: #ffffff;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  min-width: 400px;
}

.modal h3 {
  margin: 0 0 1.5rem 0;
  color: #2c3e50;
  font-weight: 600;
}

.rename-input, .folder-input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  font-size: 1rem;
  margin-bottom: 1.5rem;
  outline: none;
}

.rename-input:focus, .folder-input:focus {
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.confirm-btn, .cancel-btn {
  padding: 0.6rem 1.2rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
}

.confirm-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: #fff;
}

.confirm-btn:hover {
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
}

.cancel-btn {
  background: #ecf0f1;
  color: #2c3e50;
}

.cancel-btn:hover {
  background: #bdc3c7;
}
</style>