<template>
  <div class="other-container">
    <h3>其它功能</h3>
    <div class="other-section">
      <h4>服务器公告</h4>
      <div class="notice-card">欢迎来到 ZephyrCraft 面板！服务器将在每周一凌晨自动重启。</div>
        </div>
    <div class="other-section">
      <h4>操作日志</h4>
      <ul class="log-list">
        <li>2024-06-01 12:00:01 管理员 Steve 重启了服务器</li>
        <li>2024-06-01 11:45:23 玩家 Alex 被踢出</li>
        <li>2024-06-01 11:30:10 插件 LuckPerms 已更新</li>
      </ul>
        </div>
    <div class="other-section">
      <h4>常见问题</h4>
      <ul class="faq-list">
        <li>如何添加插件？—— 请在“插件”页面上传 jar 文件。</li>
        <li>如何修改服务器配置？—— 请在“文件”页面编辑 server.properties。</li>
        <li>如何联系技术支持？—— QQ群：123456789</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
// 假数据
</script>

<style scoped>
.other-container {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  padding: 2rem 2.5rem;
}
.other-section {
  margin-bottom: 2rem;
}
.notice-card {
  background: #f5f6fa;
  border-radius: 6px;
  padding: 1rem 1.5rem;
  color: #23272e;
  font-size: 1.1rem;
  margin-bottom: 1rem;
}
.log-list, .faq-list {
  margin: 0;
  padding-left: 1.2rem;
  color: #555;
  font-size: 1rem;
}
.log-list li, .faq-list li {
  margin-bottom: 0.5rem;
}
</style>