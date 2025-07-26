<template>
  <div class="files-container">
    <h2>文件管理</h2>

    <!-- 路径导航栏 -->
    <div class="path-navigation">
      <div class="path-breadcrumb">
        <span class="path-segment" :class="{ current: currentPath === '' }" @click="navigateToPath('')">根目录</span>
        <template v-for="(segment, index) in pathSegments" :key="index">
          <span class="separator">/</span>
          <span class="path-segment" :class="{ current: index === pathSegments.length - 1 && currentPath !== '' }"
            @click="navigateToPath(getPathUpTo(index))">{{ segment }}</span>
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
        <!-- 多选模式按钮 -->
        <button @click="toggleMultiSelectMode" class="action-btn" :class="{ 'active': isMultiSelectMode }">
          <IconMdiCheckboxMultipleMarkedOutline v-if="isMultiSelectMode" style="margin-right: 5px;" />
          <IconMdiCheckboxMultipleBlankOutline v-else style="margin-right: 5px;" />
          {{ isMultiSelectMode ? '退出多选' : '多选模式' }}
        </button>
        <!-- 压缩按钮 -->
        <button v-if="isMultiSelectMode && selectedFiles.size > 0" @click="showZipDialog"
          class="action-btn compress-btn">
          <IconMdiArchive style="margin-right: 5px;" />压缩选中文件
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
        <div v-for="file in files" :key="file.path" class="file-item" :class="{
          'is-directory': file.isDir,
          'is-editable': file.isEditable,
          'is-selected': isFileSelected(file),
          'multi-select-mode': isMultiSelectMode
        }" @click="isMultiSelectMode ? toggleFileSelection(file) : handleFileClick(file)"
          @contextmenu.prevent="showContextMenu($event, file)" @dblclick="handleFileDoubleClick(file)">
          <!-- 多选复选框 -->
          <div v-if="isMultiSelectMode" class="file-checkbox" @click.stop="toggleFileSelection(file)">
            <IconMdiCheckboxMarked v-if="isFileSelected(file)" class="checkbox-checked" />
            <IconMdiCheckboxBlankOutline v-else class="checkbox-unchecked" />
          </div>

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
            <button v-if="file.isEditable && !file.isDir && file.path" @click.stop="editFile(file)" class="edit-btn"
              title="编辑文件">
              <IconMdiPencil />
            </button>
            <button v-if="!file.isDir && file.path" @click.stop="downloadFile(file)" class="download-btn" title="下载文件">
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

    <!-- 压缩对话框 -->
    <div v-if="zipDialog.show" class="dialog-overlay" @click="cancelZip">
      <div class="dialog-content" @click.stop>
        <div class="dialog-header">
          <h3>创建压缩文件</h3>
          <button @click="cancelZip" class="close-btn">
            <IconMdiClose />
          </button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>压缩文件名：</label>
            <input v-model="zipDialog.fileName" type="text" placeholder="请输入压缩文件名" class="form-input" />
            <small>将自动添加 .zip 扩展名</small>
          </div>
          <div class="selected-files-info">
            <h4>选中的文件 ({{ selectedFiles.size }} 个)：</h4>
            <div class="selected-files-list">
              <div v-for="file in getSelectedFiles()" :key="file.path || file.name" class="selected-file-item">
                <IconMdiFolder v-if="file.isDir" style="margin-right: 5px;" />
                <IconMdiFile v-else style="margin-right: 5px;" />
                {{ file.name }}
              </div>
            </div>
          </div>
        </div>
        <div class="dialog-footer">
          <button @click="cancelZip" class="btn btn-secondary">取消</button>
          <button @click="createZip" class="btn btn-primary">创建压缩文件</button>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <div v-if="contextMenu.show" class="context-menu" :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
      @click.stop>
      <!-- 多选模式下的右键菜单 -->
      <template v-if="isMultiSelectMode">
        <div class="context-menu-item" @click="selectAllFiles">
          <IconMdiCheckboxMultipleMarked style="margin-right: 8px;" />
          全选
        </div>
        <div class="context-menu-item" @click="deselectAllFiles">
          <IconMdiCheckboxMultipleBlankOutline style="margin-right: 8px;" />
          取消全选
        </div>
        <div class="context-menu-separator"></div>
        <div class="context-menu-item" @click="showZipDialog">
          <IconMdiArchive style="margin-right: 8px;" />
          压缩选中文件
        </div>
      </template>

      <!-- 普通模式下的右键菜单 -->
      <template v-else>
        <div v-if="contextMenu.file && contextMenu.file.isEditable && !contextMenu.file.isDir && contextMenu.file.path"
          class="context-menu-item" @click="editFile(contextMenu.file)">
          <IconMdiPencil style="margin-right: 8px;" />
          编辑
        </div>
        <div v-if="contextMenu.file && !contextMenu.file.isDir && contextMenu.file.path" class="context-menu-item"
          @click="downloadFile(contextMenu.file)">
          <IconMdiDownload style="margin-right: 8px;" />
          下载
        </div>
        <div class="context-menu-separator"></div>
        <div class="context-menu-item" @click="renameFile(contextMenu.file)">
          <IconMdiRename style="margin-right: 8px;" />
          重命名
        </div>
        <div class="context-menu-item delete-item" @click="deleteFile(contextMenu.file)">
          <IconMdiDelete style="margin-right: 8px;" />
          删除
        </div>
      </template>
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
          <textarea v-model="editor.content" class="file-editor" placeholder="文件内容..." spellcheck="false"></textarea>
        </div>
      </div>
    </div>

    <!-- 重命名对话框 -->
    <div v-if="renameDialog.show" class="modal-overlay">
      <div class="modal">
        <h3>重命名</h3>
        <input v-model="renameDialog.newName" type="text" class="rename-input" placeholder="输入新名称"
          @keyup.enter="confirmRename" @keyup.esc="cancelRename" />
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
        <input v-model="createFolderDialog.name" type="text" class="folder-input" placeholder="输入文件夹名称"
          @keyup.enter="confirmCreateFolder" @keyup.esc="cancelCreateFolder" />
        <div class="modal-actions">
          <button @click="confirmCreateFolder" class="confirm-btn">创建</button>
          <button @click="cancelCreateFolder" class="cancel-btn">取消</button>
        </div>
      </div>
    </div>

    <!-- 隐藏的文件上传输入 -->
    <input ref="fileInput" type="file" multiple style="display: none" @change="handleFileUpload" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted, onUnmounted, watch} from 'vue';

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

// 在script setup部分添加多选相关的响应式变量
const selectedFiles = ref<Set<string>>(new Set());
const isMultiSelectMode = ref(false);
const zipDialog = ref({
  show: false,
  fileName: ''
});

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
  // 获取视口尺寸
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;

  // 预估菜单尺寸
  const menuWidth = 200;
  const menuHeight = 150;

  // 计算菜单位置
  let x = event.clientX;
  let y = event.clientY;

  // 检查右边界
  if (x + menuWidth > viewportWidth) {
    x = viewportWidth - menuWidth - 10;
  }

  // 检查下边界
  if (y + menuHeight > viewportHeight) {
    y = viewportHeight - menuHeight - 10;
  }

  // 确保不会超出左边界和上边界
  x = Math.max(10, x);
  y = Math.max(10, y);

  contextMenu.value = {
    show: true,
    x: x,
    y: y,
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
  reader.onload = function (e) {
    const base64 = (e.target?.result as string).split(',')[1];
    wsApi.send({ command: 'uploadFile', filePath: (currentPath.value ? currentPath.value + '/' : '') + file.name, fileContent: base64 });
  };
  reader.readAsDataURL(file);
  target.value = '';
}

// 多选相关方法
function toggleMultiSelectMode() {
  isMultiSelectMode.value = !isMultiSelectMode.value;
  if (!isMultiSelectMode.value) {
    selectedFiles.value.clear();
  }
}

function toggleFileSelection(file: any) {
  if (!isMultiSelectMode.value) return;

  const fileKey = file.path || file.name;
  if (selectedFiles.value.has(fileKey)) {
    selectedFiles.value.delete(fileKey);
  } else {
    selectedFiles.value.add(fileKey);
  }
}

function isFileSelected(file: any) {
  const fileKey = file.path || file.name;
  return selectedFiles.value.has(fileKey);
}

function selectAllFiles() {
  if (!isMultiSelectMode.value) return;
  files.value.forEach(file => {
    const fileKey = file.path || file.name;
    selectedFiles.value.add(fileKey);
  });
}

function deselectAllFiles() {
  selectedFiles.value.clear();
}

function getSelectedFiles() {
  return files.value.filter(file => {
    const fileKey = file.path || file.name;
    return selectedFiles.value.has(fileKey);
  });
}

function showZipDialog() {
  if (selectedFiles.value.size === 0) {
    alert('请先选择要压缩的文件');
    return;
  }
  zipDialog.value.show = true;
  zipDialog.value.fileName = `压缩文件_${new Date().toISOString().slice(0, 10)}`;
}

function createZip() {
  const selectedFileList = getSelectedFiles();
  if (selectedFileList.length === 0) {
    alert('请先选择要压缩的文件');
    return;
  }

  const fileName = zipDialog.value.fileName.trim();
  if (!fileName) {
    alert('请输入压缩文件名');
    return;
  }

  // 构建文件路径列表
  const filesToZip = selectedFileList.map(file => file.path || file.name);

  // 发送压缩请求
  wsApi.send({
    command: 'createZip',
    filesToZip: filesToZip,
    zipFileName: fileName
  });

  // 关闭对话框并清除选择
  zipDialog.value.show = false;
  zipDialog.value.fileName = '';
  selectedFiles.value.clear();
  isMultiSelectMode.value = false;
}

function cancelZip() {
  zipDialog.value.show = false;
  zipDialog.value.fileName = '';
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
  else if (data.status !== 'download' && data.fileContent !== undefined) {
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
  /* background: var(--bg-primary); */
  padding: 1.5rem;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.files-container h2 {
  color: var(--text-primary);
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.path-navigation {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
  border: 1px solid var(--border-color);
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
  color: var(--info-color);
  cursor: pointer;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  transition: all 0.2s;
  font-weight: 500;
}

.path-segment:hover {
  background: var(--accent-light);
}

.path-segment.current {
  color: var(--text-primary);
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
  padding: 0.5rem 1rem;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--text-primary);
  transition: all 0.2s;
  display: flex;
  align-items: center;
  box-shadow: none;
}

.action-btn:hover {
  background: var(--border-light);
  border-color: var(--accent-color);
}

.action-btn.active {
  background: var(--info-color);
  color: #fff;
  border-color: var(--info-color);
}

.action-btn.compress-btn {
  background: var(--warning-color);
  color: #fff;
  border-color: var(--warning-color);
  animation: pulse 2s infinite;
}

.action-btn.compress-btn:hover {
  background: var(--warning-color);
  box-shadow: 0 2px 8px var(--warning-color);
}

@keyframes pulse {
  0% { box-shadow: 0 0 0 0 var(--warning-color); }
  70% { box-shadow: 0 0 0 10px rgba(230, 126, 34, 0); }
  100% { box-shadow: 0 0 0 0 rgba(230, 126, 34, 0); }
}

.files-section {
  flex: 1;
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
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
  color: var(--text-secondary);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--border-light);
  border-top: 4px solid var(--accent-color);
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
  color: var(--error-color);
  background: var(--bg-primary);
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
  background: var(--bg-primary);
  border-radius: 6px;
  border: 1px solid var(--border-color);
  margin-bottom: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.file-item:hover {
  box-shadow: 0 2px 8px var(--accent-light);
  border-color: var(--accent-color);
}

.file-item.is-directory {
  background: var(--bg-secondary);
}

.file-item.is-editable {
  border-left: 3px solid var(--info-color);
}

.file-item.is-selected {
  background: var(--accent-light);
  border-color: var(--info-color);
  box-shadow: 0 2px 8px var(--info-color);
}

.file-item.multi-select-mode {
  padding-left: 2.5rem;
}

.file-checkbox {
  position: absolute;
  left: 0.5rem;
  top: 50%;
  transform: translateY(-50%);
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1;
}

.checkbox-checked {
  color: var(--info-color);
}

.checkbox-unchecked {
  color: var(--text-secondary);
}

.file-icon {
  margin-right: 1rem;
  font-size: 1.5rem;
  color: var(--text-secondary);
  width: 24px;
  text-align: center;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 0.2rem;
  word-break: break-all;
}

.file-details {
  font-size: 0.8rem;
  color: var(--text-secondary);
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

.edit-btn,
.download-btn {
  background: none;
  border: none;
  padding: 0.3rem;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--text-secondary);
}

.edit-btn:hover {
  background: var(--accent-light);
  color: var(--info-color);
}

.download-btn:hover {
  background: var(--accent-light);
  color: var(--accent-color);
}

.no-files {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  color: var(--text-secondary);
  text-align: center;
}

.context-menu {
  position: fixed;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  box-shadow: 0 4px 12px var(--shadow-color);
  z-index: 1000;
  min-width: 160px;
  max-width: 250px;
  overflow: hidden;
  animation: contextMenuFadeIn 0.15s ease-out;
  max-height: calc(100vh - 20px);
  overflow-y: auto;
}

@keyframes contextMenuFadeIn {
  from { opacity: 0; transform: scale(0.95) translateY(-5px); }
  to { opacity: 1; transform: scale(1) translateY(0); }
}

.context-menu-item {
  padding: 0.75rem 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  font-size: 0.9rem;
  color: var(--text-primary);
  transition: background-color 0.2s;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.context-menu-item:hover {
  background: var(--bg-secondary);
}

.context-menu-item:active {
  background: var(--border-light);
}

.context-menu-item.danger:hover {
  background: var(--error-color);
  color: #fff;
}

.context-menu-separator {
  height: 1px;
  background: var(--border-color);
  margin: 0.5rem 0;
}

.context-menu-item.delete-item {
  color: var(--error-color);
}

.context-menu-item.delete-item:hover {
  background: var(--error-color);
  color: #fff;
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
  background: var(--bg-primary);
  border-radius: 8px;
  box-shadow: 0 8px 32px var(--shadow-color);
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
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
  border-radius: 8px 8px 0 0;
}

.editor-header h3 {
  margin: 0;
  color: var(--text-primary);
  font-weight: 600;
}

.editor-actions {
  display: flex;
  gap: 0.8rem;
}

.save-btn,
.close-btn {
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
  background: var(--accent-gradient);
  color: #fff;
}

.save-btn:hover:not(:disabled) {
  box-shadow: 0 2px 8px var(--accent-light);
}

.save-btn:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
}

.close-btn {
  background: var(--error-color);
  color: #fff;
}

.close-btn:hover {
  background: var(--error-color);
  opacity: 0.8;
}

.editor-content {
  flex: 1;
  padding: 1rem;
}

.file-editor {
  width: 100%;
  height: 100%;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 1rem;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
  line-height: 1.5;
  resize: none;
  outline: none;
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.file-editor:focus {
  border-color: var(--info-color);
  box-shadow: 0 0 0 2px var(--info-color);
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
  background: var(--bg-primary);
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 8px 32px var(--shadow-color);
  min-width: 400px;
}

.modal h3 {
  margin: 0 0 1.5rem 0;
  color: var(--text-primary);
  font-weight: 600;
}

.rename-input,
.folder-input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 1rem;
  margin-bottom: 1.5rem;
  outline: none;
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.rename-input:focus,
.folder-input:focus {
  border-color: var(--info-color);
  box-shadow: 0 0 0 2px var(--info-color);
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.confirm-btn,
.cancel-btn {
  padding: 0.6rem 1.2rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
}

.confirm-btn {
  background: var(--accent-gradient);
  color: #fff;
}

.confirm-btn:hover {
  box-shadow: 0 2px 8px var(--accent-light);
}

.cancel-btn {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.cancel-btn:hover {
  background: var(--border-light);
}

.dialog-overlay {
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

.dialog-content {
  background: var(--bg-primary);
  border-radius: 8px;
  box-shadow: 0 8px 32px var(--shadow-color);
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
  border-radius: 8px 8px 0 0;
}

.dialog-header h3 {
  margin: 0;
  color: var(--text-primary);
  font-weight: 600;
}

.dialog-body {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: var(--text-primary);
  font-size: 0.9rem;
}

.form-input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 1rem;
  outline: none;
  background: var(--bg-secondary);
  color: var(--text-primary);
}

.form-input:focus {
  border-color: var(--info-color);
  box-shadow: 0 0 0 2px var(--info-color);
}

.selected-files-info {
  margin-top: 1.5rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-color);
}

.selected-files-info h4 {
  margin: 0 0 0.8rem 0;
  color: var(--text-primary);
  font-weight: 600;
}

.selected-files-list {
  max-height: 150px;
  overflow-y: auto;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  padding: 0.5rem;
}

.selected-file-item {
  display: flex;
  align-items: center;
  padding: 0.5rem 0.8rem;
  background: var(--bg-secondary);
  border-radius: 3px;
  margin-bottom: 0.3rem;
  font-size: 0.9rem;
  color: var(--text-primary);
}

.selected-file-item:last-child {
  margin-bottom: 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--border-color);
  background: var(--bg-secondary);
  border-radius: 0 0 8px 8px;
}

.btn {
  padding: 0.6rem 1.2rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 600;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-primary {
  background: var(--accent-gradient);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  box-shadow: 0 2px 8px var(--accent-light);
}

.btn-primary:disabled {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  margin-right: 10px;
  box-shadow: none;
}

.btn-secondary:hover {
  background: var(--border-light);
}
</style>