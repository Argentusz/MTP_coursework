<template>
  <div class="memory-viewer">
    <div class="memory-viewer__tabs">
      <tab-button text="User" :active="segment === 2" @click="segment = 2"/>
      <tab-button text="Int" :active="segment === 1" @click="segment = 1"/>
      <tab-button text="Label" :active="segment === 3" @click="segment = 3"/>
      <tab-button text="Calls" :active="segment === 4" @click="segment = 4"/>
      <tab-button text="Exec" :active="segment === 0" @click="segment = 0"/>
    </div>
    <div class="memory-viewer__table">
      <memory-table
        :table="xmem.Segments[segment].Table"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import TabButton from "./buttons/TabButton.vue";
import MemoryTable from "./tables/MemoryTable.vue";

const segment = ref(2)

const $props = defineProps({
  xmem: {
    type: [Object, null],
    default: null
  }
})
</script>

<style scoped>
.memory-viewer__tabs {
  display: flex;
  flex-direction: row;
  gap: 5px;
  justify-content: space-around;
  align-items: center;
  height: 75px;
  background: #242021;
}
</style>