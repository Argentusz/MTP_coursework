<template>
  <div class="app">
    <v-header
        :file-path="filePath"
        class="header"
        @new="fileNew"
        @open="fileOpen"
        @save="fileSave"
    />

    <status-bar
        :state="state"
        :file-path="filePath"
        @request="request"
        @connect="connect"
        @disconnect="disconnect"
        @compile="compile"
        @run="run"
        @send="onSend"
        @drop="onDrop"
    />

    <div class="viewers">
      <div class="viewer" style="width: 33%">
        <code-editor v-model="program" :file-path="filePath"/>
      </div>
      <template v-if="state">
        <input-viewer :input="state.Input" :nib="state.CPU.RRAM.SYS.NIB" :lbl="state.CPU.XMEM.Segments[3].Table" class="viewer" style="width: 17%"/>
        <memory-viewer v-model="toSet" :xmem="state.CPU.XMEM" class="viewer"/>
        <registers-viewer v-model="toSet" :rram="state.CPU.RRAM" class="viewer"/>
      </template>
    </div>

    <error-flash v-model="error"/>

  </div>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import VHeader from "./components/VHeader.vue";
import consts from "./js/consts";
import StatusBar from "./components/StatusBar.vue";
import MemoryViewer from "./components/MemoryViewer.vue";
import RegistersViewer from "./components/RegistersViewer.vue";
import InputViewer from "./components/InputViewer.vue";
import CodeEditor from "./components/CodeEditor.vue";
import ErrorFlash from "./components/ErrorFlash.vue";

const
  program = ref(null),
  filePath = ref(null)

const connect = () => window.mtpAPI.connect(["-intr", "-sudo"])
const disconnect = window.mtpAPI.disconnect

const fileNew = () => window.fileAPI.new()
const fileOpen = () => window.fileAPI.open()
const fileSave = () => window.fileAPI.save(filePath.value, program.value)

window.fileAPI.onOpened(a => { program.value = a.data; filePath.value = a.filePath })
window.fileAPI.onSaved(() => {})

const state = ref(null)
const request = (command) => {
  console.log("[LOG] Request", command)
  window.mtpAPI.request(command)
}
const compile = async () => {
  await fileSave()
  request(`compile ${filePath.value}`)
}
const run = (all = false) => {
  state.value.Running = true
  request(all ? "run all" : "run")
}

window.mtpAPI.onUpdate(a => {
  if (!a) {
    state.value = null
    return
  }

  try {
    state.value = JSON.parse(a?.trim())
    console.log("[LOG] Update", state.value)
  } catch (e) {
    console.error(e, a)
  }
})
onMounted(window.mtpAPI.ping)

const error = ref("")
window.mtpAPI.onError((a) => {
  error.value = a
  console.error(a)
})

const toSet = ref({
  RRAM: {},
  XMEM: {},
})
const onSend = () => {
  request(`set ${JSON.stringify(toSet.value)}`)
}
const onDrop = () => {
  toSet.value = { RRAM: {}, XMEM: {}, }
}
</script>

<style scoped>
.app {
  background-color: #191516;
  color: white;
  height: 100%;
  font-family: "Helvetica Neue", serif;
}

.viewers {
  display: flex;
  flex-direction: row;
  height: calc(100% - 153px);
}

.viewer {
  width: 25%;
  border-right: 1px solid gray;
  max-height: 100%;
  overflow-y: scroll;
}
</style>