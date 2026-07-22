<script lang="ts">
import { defineComponent } from 'vue'
import LobbyView from './components/LobbyView.vue'
import BoardView from './components/BoardView.vue'
import { store } from './store'

export default defineComponent({
  name: 'App',
  components: { LobbyView, BoardView },
  computed: {
    showBoard(): boolean {
      return store.connected && store.boardCells.length > 0
    },
  },
})
</script>

<template>
  <LobbyView v-if="!showBoard" />
  <BoardView v-else />
</template>
