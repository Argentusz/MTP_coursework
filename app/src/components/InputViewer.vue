<template>
  <div id="input-viewer" style="height:100%; width:100%"/>
</template>

<script setup>

import { computed, onMounted, watch } from "vue";
import { editor } from "monaco-editor";

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
  const activeLine = $props.nib ? Math.floor($props.nib / 4) : null
  return $props.input.reduce((prg, line, i) => {
    return i === activeLine
      ? `${prg}[*] ${line}\n`:
        `${prg}    ${line}\n`
  }, "").trimEnd()
})

let monEditor = null
onMounted(() => {
  monEditor = editor.create(document.getElementById("input-viewer"), {
    value: program.value,
    readOnly: true,
    scrollBeyondLastLine: false,
    theme: "mtpTheme",
    language: "mtp",
    lineNumbers: "off",
    glyphMargin: false,
    folding: false,
    lineDecorationsWidth: 0,
    lineNumbersMinChars: 0,
    minimap:{enabled:false},
  })
})

watch(program, () => monEditor.setValue(program.value))
</script>

<style scoped>

</style>