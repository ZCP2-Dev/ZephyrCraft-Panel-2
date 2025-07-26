<template>
  <div class="settings-container">
    <AppSettings />
  </div>
</template>

<script setup lang="ts">
import AppSettings from '../components/AppSettings.vue';
</script>

<style scoped>
.settings-container {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}
</style>