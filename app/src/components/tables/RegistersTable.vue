<template>
  <div class="registers-table">
    <div class="registers-table__header">
      <div class="registers-table__selectors">
        <div>
          <label>Base:</label>
          <select v-model.number="base" class="register-table__select">
            <option>2</option>
            <option>8</option>
            <option>10</option>
            <option>16</option>
          </select>
        </div>
      </div>
    </div>
    <div class="registers-table__content">
      <v-table
          :columns="columns"
          :rows="rows"
          :chars="size"
          editable
          style="height: 100%"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import VTable from "./VTable.vue";

const base = ref(2)

const columns = ['name', 'value']

const $props = defineProps({
  registers: {
    type: Array,
    default: () => []
  },
  size: {
    type: Number,
    default: () => 8
  },
  baseName: {
    type: String,
    default: "rb"
  }
})

const rows = computed(() => {
  return $props.registers.map((rg, i) => {
    return { name: `${$props.baseName}${i}`, value: formatByBase(rg.toString(base.value)) }
  })
})

const formatByBase = (str) => {
  switch (base.value) {
    case 2: return `${"0".repeat($props.size - str.length)}${str}`
    case 8: return `0o${str}`
    case 10: return str
    case 16: return `0x${str}`
    default: return str
  }
}
</script>

<style scoped>
.registers-table {
  padding: 10px;
}

.register-table__select {
  background: none;
  margin-left: 15px;
  padding: 1px 5px;
  border: 1px solid gray;
  color: inherit;
  border-radius: 5px;
}
</style>