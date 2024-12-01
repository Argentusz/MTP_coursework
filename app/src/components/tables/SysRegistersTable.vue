<template>
  <div class="registers-table">
    <div class="registers-table__content">
      <v-table
          :columns="columns"
          :rows="rows"
          :chars="32"
          style="height: 100%"
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import VTable from "./VTable.vue";

const $props = defineProps({
  sys: {
    type: [Object, null],
    default: null
  }
})

const columns = ["name", "val"]

const rows = computed(() => {
  if (!$props.sys) return []

  return [
    { name: "IR",  val: `0x${$props.sys.IR.toString(16)}`},
    { name: "NIR", val: `0x${$props.sys.NIR.toString(16)}`},
    { name: "NIB", val: `0x${$props.sys.NIB.toString(16)}`},
    { name: "CSP", val: `0x${$props.sys.CSP.toString(16)}`},
    { name: "MBR", val: $props.sys.MBR.toString()},
    { name: "TMP", val: $props.sys.TMP.toString()},
    { name: "FLB", val: `0b${$props.sys.FLB.toString(2)}`},
    { name: "FLG", val: `0b${$props.sys.FLG.toString(2)}`},
  ]
})
</script>

<style scoped>
.registers-table {
  padding: 10px;
}
</style>