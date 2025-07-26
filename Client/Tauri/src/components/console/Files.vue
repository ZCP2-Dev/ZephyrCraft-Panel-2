<template>
  <div class="files-container">
    <h2>文件管理</h2>
    <div class="file-actions">
      <h3>文件操作</h3>
      <div class="action-grid">
        <button @click="downloadFiles" class="file-btn">
          <IconMdiDownload style="margin-right: 5px;" />下载文件
        </button>
        <button @click="uploadFiles" class="file-btn upload-btn">
          <IconMdiUpload style="margin-right: 5px;" />上传文件
        </button>
        <button @click="backupWorld" class="file-btn backup-btn">
          <IconMdiBackupRestore style="margin-right: 5px;" />备份世界
        </button>
        <button @click="restoreWorld" class="file-btn restore-btn">
          <IconMdiRestore style="margin-right: 5px;" />恢复世界
        </button>
      </div>
    </div>
    
    <div class="file-info">
      <h4>文件统计</h4>
      <div class="file-stats">
        <div class="stat-item">
          <div class="stat-value">156</div>
          <div class="stat-label">总文件数</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">2.3 GB</div>
          <div class="stat-label">总大小</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">12</div>
          <div class="stat-label">备份数量</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
function downloadFiles() {
  console.log('下载文件');
}

function uploadFiles() {
  console.log('上传文件');
}

function backupWorld() {
  console.log('备份世界');
}

function restoreWorld() {
  console.log('恢复世界');
}
</script>

<style scoped>
.files-container {
  background: #ffffff;
  padding: 1.5rem;
}

.files-container h2 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.files-container h3 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.4rem;
  margin: 0 0 1.2rem 0;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.file-actions {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  border: 1px solid #e9ecef;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.file-btn {
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 1rem 1.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.3);
  position: relative;
  overflow: hidden;
}

.file-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.file-btn:hover::before {
  left: 100%;
}

.file-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(39, 174, 96, 0.4);
}

.file-btn:disabled {
  background: linear-gradient(135deg, #bdc3c7 0%, #95a5a6 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.upload-btn {
  background: linear-gradient(135deg, #3498db 0%, #2980b9 100%);
  box-shadow: 0 2px 8px rgba(52, 152, 219, 0.3);
}

.upload-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(52, 152, 219, 0.4);
}

.backup-btn {
  background: linear-gradient(135deg, #f39c12 0%, #e67e22 100%);
  box-shadow: 0 2px 8px rgba(243, 156, 18, 0.3);
}

.backup-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(243, 156, 18, 0.4);
}

.restore-btn {
  background: linear-gradient(135deg, #9b59b6 0%, #8e44ad 100%);
  box-shadow: 0 2px 8px rgba(155, 89, 182, 0.3);
}

.restore-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(155, 89, 182, 0.4);
}

.file-info {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 1.5rem;
  border: 1px solid #e9ecef;
}

.file-info h4 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.2rem;
  margin: 0 0 1rem 0;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.file-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

.stat-item {
  background: #ffffff;
  padding: 1rem;
  border-radius: 6px;
  border: 1px solid #e9ecef;
  text-align: center;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: #27ae60;
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 0.9rem;
  color: #7f8c8d;
  font-weight: 500;
}
</style>