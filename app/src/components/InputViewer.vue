<template>
  <div class="input-viewer">
    <textarea :value="program" disabled/>
  </div>
</template>

<script setup>

import { computed } from "vue";

const $props = defineProps({
  input: {
    type: Array,
    default: () => []
  },
  nib: {
    type: [Number, null],
    default: null
  }
})

const program = computed(() => {
  if (!$props.input) return ""
  // if (!$props.nib) return $props.input.join("\n")
  //
  const activeLine = $props.nib ? Math.floor($props.nib / 4) : null
  // $props.input[line] = `[*] ${$props.input[line]}`
  // return $props.input.join("\n")
  return $props.input.reduce((prg, line, i) => {
    if (i !== activeLine) return `${prg}\n    ${line}`
    return `${prg}\n[*] ${line}`
  }, "")
})
</script>

<style scoped>

</style>