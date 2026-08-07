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
      return store.currentTurnPlayer === store.playerName && store.playerName !== ''
    },
  },
});
</script>

<template>
  <!-- Game Action Panel: LevelPicker, ProblemPanel -->
  <div class="game-action-container">
    <LevelPicker v-if="isMyTurn && !store.questionActive" />
    <ProblemPanel v-else-if="store.questionActive" />
  </div>
</template>

<style scoped>
.game-action-container {
  margin: 0.5rem 0;
  width: 100%;
}
</style>