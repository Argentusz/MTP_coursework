<template>
<table class="v-table">
  <tr v-for="(r, i) of rows" :key="i" class="v-table__row">
    <th v-for="(c, j) of columns" :key="j" class="v-table__cell">
      <span v-if="!editable || j === 0">{{ r[c] }}</span>
      <input v-else :value="r[c]" class="v-table__input" @input="$emit('edit', i, j, $event.target.value)">
    </th>
  </tr>
</table>
</template>

<script setup>
import { computed } from "vue";

const $props = defineProps({
  columns: {
    type: Array,
    required: true
  },
  rows: {
    type: Array,
    default: () => []
  },
  editable: {
    type: Boolean,
    default: false
  },
  chars: {
    type: Number,
    default: 8,
  }
})

const $emit = defineEmits(["edit"])

const inputWidth = computed(() => {
  return `${Math.floor($props.chars) / 2}rem`
})
</script>

<style scoped>
.v-table {
  width: 100%;
  border-spacing: 0;
  font-weight: normal;
}

.v-table__cell {
  border-right: 1px solid gray;
  font-weight: normal;
}

.v-table__cell:first-child {
  border-left: 1px solid gray;
  width: 75px;
  font-weight: bold;
  font-size: 0.8rem;
}

.v-table__row .v-table__cell {
  border-bottom: 1px solid gray;
}

.v-table__row:first-child .v-table__cell {
  border-top: 1px solid gray;
}

.v-table__input {
  background: none;
  color: inherit;
  border: none;
  font-family: "Helvetica Neue", monospace;
  width: v-bind(inputWidth);
  text-align: center;
}
</style>