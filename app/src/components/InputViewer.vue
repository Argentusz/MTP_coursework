<template>
  <div id="input-viewer" style="height:100%; width:100%"/>
</template>

<script setup>

import { computed, onMounted, watch } from "vue";
import * as monaco from "monaco-editor";

const $props = defineProps({
  input: {
    type: Array,
    default: () => []
  },
  nib: {
    type: [Number, null],
    default: null
  },
  lbl: {
    type: [Object, null],
    default: null
  }
})

const program = computed(() => {
  if (!$props.input) return ""

  let maxLen = 0
  $props.input.forEach(line => {
    maxLen = maxLen < line.length ? line.length : maxLen
  })

  const reverseLbl = !$props.lbl
    ? {}
    : Object.keys($props.lbl).reduce((obj, addrStr) => {
        let addr = parseInt(addrStr)
        if (addr % 4 !== 0) {
          addr -= addr %4
          if (obj[addr - (addr % 4)] !== undefined) {
            return obj
          }
        }

        let exeAddrStr = ""
        for (let i = 0; i < 4; i++) {
          const part = $props.lbl[addr+i] || 0
          const partStr = part.toString(2)
          const prefix = "0".repeat(8 - partStr.length)
          exeAddrStr += `${prefix}${partStr}`
        }

        const lineN = Math.floor(parseInt(exeAddrStr, 2) / 4)

        obj[lineN] = `0x${Math.floor(addr/4).toString(16)}`
        return obj
    }, {})

  return $props.input.reduce((prg, line, i) =>
      reverseLbl[i] === undefined
        ? `${prg}${line}\n`
        : `${prg}${line}${" ".repeat(maxLen-line.length)} ; <${reverseLbl[i]}>\n`
  , "").trimEnd()
})

let monEditor = null
onMounted(() => {
  monEditor = monaco.editor.create(document.getElementById("input-viewer"), {
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

let activeLineDecorations = null

watch(program, () => monEditor.setValue(program.value))
watch(() => $props.nib, () => {
  if (activeLineDecorations) activeLineDecorations.clear()
  if (!$props.nib) return

  const activeLine = Math.floor($props.nib / 4) + 1
  activeLineDecorations = monEditor.createDecorationsCollection([{
      range: new monaco.Range(activeLine, 1, activeLine, 1),
      options: {
        isWholeLine: true,
        className: 'activeLineDecoration',
        marginClassName: 'activeLineDecoration'
      }
  }])

  console.log(activeLine, activeLineDecorations)
})
</script>

<style>
.activeLineDecoration {
  background-color: #DBB8FF;
  opacity: 0.5;
}
</style>