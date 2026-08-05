<script lang="ts">
import { defineComponent } from 'vue'
import LobbyView from './components/LobbyView.vue'
import BoardView from './components/BoardView.vue'
import AdminSpectatorView from './components/AdminSpectatorView.vue'
import { store } from './store'

export default defineComponent({
  name: 'App',
  components: { LobbyView, BoardView, AdminSpectatorView },
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
  <!-- Admin spectator mode: open the page with ?admin=1 -->
  <AdminSpectatorView v-if="isAdmin" />
  <template v-else>
    <LobbyView v-if="!showBoard" />
    <BoardView v-else />
  </template>
</template>

