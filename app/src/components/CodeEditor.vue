<template>
<!--  <textarea v-model="program" style="resize: none; height: 100%; width: 100%"/>-->
  <div id="code-editor" style="height:100%; width:100%"/>
</template>

<script setup>
import { onMounted, watch } from "vue";
import { editor, languages } from "monaco-editor";

const program = defineModel({ type: [String, null], default: null })
const $props = defineProps({ filePath: { type: [String, null], default: null } })

let monEditor = null
onMounted(() => {
  monEditor = editor.create(document.getElementById("code-editor"), {
    value: program.value,
    scrollBeyondLastLine: false,
    theme: "mtpTheme",
    language: "mtp"
  })
  monEditor.onDidChangeModelContent(() => {
    program.value = monEditor.getValue()
  })
})

watch(() => $props.filePath, () => monEditor.setValue(program.value))
</script>

<style scoped>

</style>