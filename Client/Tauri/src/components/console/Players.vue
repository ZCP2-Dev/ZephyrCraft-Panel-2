<template>
  <div class="players-container">
    <h3>在线玩家</h3>
    <div class="players-list">
      <div v-for="player in players" :key="player.name" class="player-item">
        <span class="player-name">{{ player.name }}</span>
        <span class="player-actions">
          <button class="action-btn">踢出</button>
          <button class="action-btn">OP</button>
          <button class="action-btn">禁言</button>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const players = [
  { name: 'Steve' },
  { name: 'Alex' },
  { name: 'Notch' },
  { name: 'Herobrine' },
  { name: 'Zombie' },
  { name: 'Skeleton' }
];
</script>

<style scoped>
.players-container {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  padding: 2rem 2.5rem;
}
.players-list {
  margin-top: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.player-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.8rem 1.2rem;
  background: #f5f6fa;
  border-radius: 6px;
}
.player-name {
  font-weight: bold;
  color: #23272e;
}
.player-actions {
  display: flex;
  gap: 0.8rem;
}
.action-btn {
  background: #88bf64;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 0.3rem 1rem;
  font-size: 0.95rem;
  cursor: pointer;
  transition: background 0.2s;
}
.action-btn:hover {
  background: #539951;
}
</style>