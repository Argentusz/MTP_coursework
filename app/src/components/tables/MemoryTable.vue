<template>
  <div class="memory-table">
    <div class="memory-table__header">
      <div class="memory-table__selectors">
        <div>
          <label>Base:</label>
          <select v-model.number="base" class="memory-table__select">
            <option>2</option>
            <option>8</option>
            <option>10</option>
            <option>16</option>
          </select>
        </div>

        <div>
          <label>Group:</label>
          <select v-model.number="group" class="memory-table__select">
            <option>8</option>
            <option>16</option>
            <option>32</option>
          </select>
        </div>
      </div>
      <div style="margin-bottom: 7px;">
        <label>Address:</label>
        <input :value="fromRenderer" class="memory-table__input" @input="onFromInput"/>
      </div>
    </div>
    <div class="memory-table__content">
      <v-table
        :columns="columns"
        :rows="rows"
        :chars="chars"
        editable
        style="height: 100%"
        @edit="onEdit"
      />
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from "vue";
import VTable from "./VTable.vue";

const ROWS_ON_PAGE = 30

const $emit = defineEmits(["edit"])

const $props = defineProps({
  table: {
    type: Object,
    default: () => {}
  },
  toSet: {
    type: Object,
    default: () => ({ XMEM: {} })
  }
})

const base = ref(2)
const group = ref(8)
const from = ref(0)

const columns = computed(() => {
  const col = ["addr"]
  const n = Math.floor(32 / group.value)
  for (let i = 0; i < n; i++) col.push(`w${i}`)
  return col
})

const chars = computed(() => {
  switch (base.value) {
    case 2: return group.value
    case 8: return Math.ceil(group.value / 3) + 2
    case 10: return Math.max(Math.floor(group.value / 3), 3)
    case 16: return Math.floor(group.value / 4) + 2
  }
})

const rows = computed(() => {
  const res = []

  const maxAddr = ((ROWS_ON_PAGE - 1) * 4 + from.value).toString(16)

  const wordsAmt = Math.floor(32 / group.value)
  for (let i = 0; i < ROWS_ON_PAGE; i++) {
    const addr = i * 4 + from.value
    let addrStr = addr.toString(16)
    addrStr = `0x${"0".repeat(maxAddr.length - addrStr.length)}${addrStr}`

    const obj = { "addr": `${addrStr}` }
    for (let j = 0; j < wordsAmt; j++) {
      const offset = j * Math.floor(group.value / 8)
      obj[`w${j}`] = genCell(addr + offset)
    }

    res.push(obj)
  }
  return res
})

const genCell = (addr) => {
  if (group.value === 8) {
    return formatByBase(($props.table[addr] || 0).toString(base.value))
  }

  const cells = []
  for (let i = 0; i < Math.floor(group.value / 8); i++)
    cells.push($props.toSet[addr+i] || $props.table[addr+i] || 0)

  const str = cells.reduce((acc, cell) => {
    let cellStr = cell.toString(2)
    cellStr = `${"0".repeat(8 - cellStr.length)}${cellStr}`
    return `${acc}${cellStr}`
  }, "")

  if (base.value === 2) {
    return formatByBase(str)
  }

  const num = parseInt(str, 2)
  return formatByBase(num.toString(base.value))
}

const formatByBase = (str) => {
  switch (base.value) {
    case 2: return `${"0".repeat(chars.value - str.length)}${str}`
    case 8: return `0o${str}`
    case 10: return str
    case 16: return `0x${str}`
    default: return str
  }
}

const fromRenderer = computed(() => {
  return `0x${from.value.toString(16)}`
})

const onFromInput = (e) => {
  from.value = parseInt(e.target.value.slice(2) || "0", 16)
}

const onEdit = (i, j, v) => {
  const baseAddr = parseInt(rows.value[i].addr.replace("0x", ""), 16)

  v = v.replace(/0x|0o/, "")
  v = `${"0".repeat(group.value - v.length)}${v}`
  const values = []
  for (let sub = 0; sub < group.value; sub+=8) {
    values.push(parseInt(v.substring(sub, sub + 8), base.value))
  }

  values.forEach((value, indent) => {
    $emit("edit", baseAddr + indent + j - 1, value)
  })
}
</script>

<style scoped>
.memory-table {
  padding: 10px;
}

.memory-table__header {
  display: flex;
  flex-direction: column;
}

.memory-table__selectors {
  display: flex;
  flex-direction: row;
  gap: 10px;
  margin-bottom: 7px;
}

.memory-table__select,
.memory-table__input {
  background: none;
  margin-left: 15px;
  padding: 1px 5px;
  border: 1px solid gray;
  color: inherit;
  border-radius: 5px;
}
</style>