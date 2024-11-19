<template>
  <div class="register-viewer">
    <div class="register-viewer__tabs">
      <tab-button text="General Purpose" :active="!isSystemRegistersShown" @click="isSystemRegistersShown = false"/>
      <tab-button text="System" :active="isSystemRegistersShown" @click="isSystemRegistersShown = true"/>
    </div>
    <div class="register-viewer__table">
      <template v-if="isSystemRegistersShown">

      </template>
      <template v-else>
        <registers-table :registers="rram.SGPRs" :size="8" baseName="rb"/>
        <registers-table :registers="rram.GPRs" :size="32" baseName="rw"/>
        <registers-table :registers="rram.XGPRs" :size="32" baseName="rx"/>
        <registers-table :registers="rram.HGPRs" :size="16" baseName="rh"/>
        <registers-table :registers="rram.LGPRs" :size="16" baseName="rl"/>
      </template>
    </div>
  </div>
</template>

<script setup>
import TabButton from "./buttons/TabButton.vue";
import { ref } from "vue";
import RegistersTable from "./tables/RegistersTable.vue";

const isSystemRegistersShown = ref(false)

defineProps({
  rram: {
    type: Object,
    default: () => {}
  }
})
</script>

<style scoped>
.register-viewer__tabs {
  display: flex;
  flex-direction: row;
  gap: 5px;
  justify-content: space-around;
  align-items: center;
  height: 75px;
  background: #242021;
}
</style>