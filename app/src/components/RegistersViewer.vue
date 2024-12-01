<template>
  <div class="register-viewer">
    <div class="register-viewer__tabs">
      <tab-button text="General Purpose" :active="!isSystemRegistersShown" @click="isSystemRegistersShown = false"/>
      <tab-button text="System" :active="isSystemRegistersShown" @click="isSystemRegistersShown = true"/>
    </div>
    <div class="register-viewer__table">
      <template v-if="isSystemRegistersShown">
        <sys-registers-table :sys="rram.SYS"/>
      </template>
      <template v-else>
        <registers-table :registers="rram.SGPRs" :size="8" baseName="rb" @edit="(i, v) => onEdit(i, v)"/>
        <registers-table :registers="rram.GPRs" :size="32" baseName="rw" @edit="(i, v) => onEdit(i, v, 8)"/>
        <registers-table :registers="rram.XGPRs" :size="32" baseName="rx" @edit="(i, v) => onEdit(i, v, 8 + 32)"/>
        <registers-table :registers="rram.HGPRs" :size="16" baseName="rh" @edit="(i, v) => onEdit(i, v, 8 + 32 + 8)"/>
        <registers-table :registers="rram.LGPRs" :size="16" baseName="rl" @edit="(i, v) => onEdit(i, v, 8 + 32 + 8 + 8)"/>
      </template>
    </div>
  </div>
</template>

<script setup>
import TabButton from "./buttons/TabButton.vue";
import { ref } from "vue";
import RegistersTable from "./tables/RegistersTable.vue";
import SysRegistersTable from "./tables/SysRegistersTable.vue";

const isSystemRegistersShown = ref(false)

const toEdit = defineModel({ type: Object, default: () => ({ RRAM: {} }) })

defineProps({
  rram: {
    type: Object,
    default: () => {}
  }
})

const onEdit = (i, v, offset = 0) => {
  toEdit.value.RRAM[offset + i] = v
}
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