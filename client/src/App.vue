<script lang="ts">
import { defineComponent } from 'vue'
import LobbyView from './components/LobbyView.vue'
import BoardView from './components/BoardView.vue'
import AdminView from './components/AdminView.vue'
import LanguageSwitcher from './components/LanguageSwitcher.vue'
import ThemeSwitcher from './components/ThemeSwitcher.vue'
import { store } from './store'

export default defineComponent({
  name: 'App',
  components: { LobbyView, BoardView, AdminView, LanguageSwitcher, ThemeSwitcher },
  data() {
    return {
      isAdmin: false,
    }
  },
  computed: {
    showBoard(): boolean {
      return store.connected && store.boardCells.length > 0
    },
  },
  created() {
    this.isAdmin = new URLSearchParams(window.location.search).get('admin') === '1'
  },
})
</script>

<template>
  <div class="app-controls">
    <ThemeSwitcher />
    <LanguageSwitcher />
  </div>
  <!-- Admin area mode: open the page with ?admin=1 -->
  <AdminView v-if="isAdmin" />
  <template v-else>
    <LobbyView v-if="!showBoard" />
    <BoardView v-else />
  </template>
</template>


<style scoped>
/* Theme and language controls stay reachable on every screen. */
.app-controls {
  position: fixed;
  top: 12px;
  right: 12px;
  z-index: 10000;
  display: flex;
  align-items: flex-start;
  gap: 8px;
}
</style>
