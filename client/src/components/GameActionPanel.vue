<script lang="ts">
import {defineComponent} from "vue";
import { store } from "../store";
import LevelPicker from "./LevelPicker.vue";
import ProblemPanel from "./ProblemPanel.vue";

export default defineComponent({
  name: 'GameActionPanel',
  components: { LevelPicker, ProblemPanel },
  data() {
    return {
      store: store,
    };
  },
  computed: {
    isMyTurn(): boolean {
      return (
        !store.gameOver &&
        store.currentTurnPlayer === store.playerName &&
        store.playerName !== ''
      )
    },
  },
});
</script>

<template>
  <!-- Game Action Panel: LevelPicker, ProblemPanel -->
  <div v-if="isMyTurn" class="game-action-container">
    <LevelPicker v-if="!store.questionActive" />
    <ProblemPanel v-else />
  </div>
</template>

<style scoped>
.game-action-container {
  margin: auto;
  width: 100%;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.7);

  position: absolute;
  top: 0;
  left: 0;
  right: 0;

  display: flex;
  justify-content: center;
  align-items: center;
  
}
</style>