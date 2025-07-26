<template>
  <div class="players-container">
    <h2>玩家管理</h2>
    <div class="players-list">
      <h3>在线玩家</h3>
      <div v-for="player in players" :key="player.name" class="player-item">
        <span class="player-name">{{ player.name }}</span>
        <button @click="kickPlayer(player.name)" class="kick-btn">
          <IconMdiAccountRemove style="margin-right: 5px;" />踢出
        </button>
      </div>
      <div v-if="players.length === 0" class="no-players">
        暂无在线玩家
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const players = ref([
  { name: 'Player1' },
  { name: 'Player2' },
  { name: 'Player3' }
]);

function kickPlayer(playerName: string) {
  // 踢出玩家逻辑
  console.log('踢出玩家:', playerName);
}
</script>

<style scoped>
.players-container {
  background: #ffffff;
  padding: 1.5rem;
}

.players-container h2 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 2rem;
  margin: 0 0 1.5rem 0;
  text-align: center;
}

.players-container h3 {
  color: #2c3e50;
  font-weight: 700;
  font-size: 1.4rem;
  margin: 0 0 1.2rem 0;
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.players-list {
  display: flex;
  flex-direction: column;
  gap: 0.8rem;
}

.player-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem 1.2rem;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.player-item:hover {
  background: #ffffff;
  box-shadow: 0 2px 8px rgba(39, 174, 96, 0.1);
  border-color: rgba(39, 174, 96, 0.2);
}

.player-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.player-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #27ae60 0%, #2ecc71 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  font-size: 1.2rem;
}

.player-details {
  display: flex;
  flex-direction: column;
}

.player-name {
  font-weight: 600;
  color: #2c3e50;
  font-size: 1.1rem;
}

.player-status {
  font-size: 0.9rem;
  color: #7f8c8d;
}

.player-actions {
  display: flex;
  gap: 0.8rem;
}

.kick-btn {
  background: linear-gradient(135deg, #e74c3c 0%, #c0392b 100%);
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(231, 76, 60, 0.3);
}

.kick-btn:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(231, 76, 60, 0.4);
}

.kick-btn:disabled {
  background: linear-gradient(135deg, #ecf0f1 0%, #bdc3c7 100%);
  cursor: not-allowed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.no-players {
  text-align: center;
  padding: 3rem 2rem;
  color: #7f8c8d;
  font-size: 1.1rem;
  background: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e9ecef;
}
</style>