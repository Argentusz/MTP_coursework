<template>
  <div class="app">
    <v-header
        :file-path="filePath"
        class="header"
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
        @run="e => run(e)"
    />

    <div class="viewers">
      <div class="viewer" style="width: 35%">
        <code-editor v-model="program" :file-path="filePath"/>
      </div>
      <template v-if="state">
        <input-viewer :input="state.Input" :nib="state.CPU.RRAM.SYS.NIB" class="viewer" style="width: 15%"/>
        <memory-viewer :xmem="state.CPU.XMEM" class="viewer"/>
        <registers-viewer :rram="state.CPU.RRAM" class="viewer"/>
      </template>
    </div>

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

const
  program = ref(null),
  filePath = ref(null)

const connect = () => window.mtpAPI.connect(["-intr", "-sudo"])
const disconnect = window.mtpAPI.disconnect

const fileOpen = () => {
  window.fileAPI.open()
}
window.fileAPI.onOpened((a) => {
  program.value = a.data
  filePath.value = a.filePath
})

const fileSave = () => window.fileAPI.save(filePath.value, program.value)
window.fileAPI.onSaved(() => {})

const state = ref(null)
const request = (command) => window.mtpAPI.request(command)
const compile = async () => {
  await fileSave()
  request(`compile ${filePath.value}`)
}
const run = (all = false) => {
  state.value.Running = true
  request(all ? "run all" : "run")
}

window.mtpAPI.onUpdate((a) => {
  try {
    state.value = JSON.parse(a)
  } catch (e) {
    console.error(e, a)
  }
})
onMounted(window.mtpAPI.ping)

const error = ref("")
window.mtpAPI.onError((a) => {
  error.value = a
})
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