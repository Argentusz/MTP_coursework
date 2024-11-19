<template>
  <div class="app">
    <v-header class="header" @open="fileOpen" @save="fileSave"/>
    <status-bar :working-state="workingState" :filePath="filePath" @request="request"/>
    <div class="viewers">
      <div class="viewer">
        <div v-if="program" style="padding: 20px"> <textarea v-model="program"/> </div>
      </div>
      <template v-if="state">
        <div class="viewer">
          <div style="padding: 20px">
            <input-viewer :input="state.Input" :nib="state.CPU.RRAM.SYS.NIB"/>
          </div>
        </div>
        <memory-viewer :xmem="state.CPU.XMEM" class="viewer"/>
        <registers-viewer :rram="state.CPU.RRAM" class="viewer"/>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import VHeader from "./components/VHeader.vue";
import consts from "./js/consts";
import StatusBar from "./components/StatusBar.vue";
import MemoryViewer from "./components/MemoryViewer.vue";
import RegistersViewer from "./components/RegistersViewer.vue";
import InputViewer from "./components/InputViewer.vue";

const request = (command) => window.mtpAPI.request(command)

const
  program = ref(null),
  filePath = ref(null)

const fileOpen = () => {
  window.fileAPI.open()
}
window.fileAPI.onOpened((a) => {
  program.value = a.data
  filePath.value = a.filePath
})

const fileSave = () => {
  window.fileAPI.save(filePath.value, program.value)
}
window.fileAPI.onSaved((a) => {

})

const state = ref(null)
window.mtpAPI.onUpdate((a) => {
  try {
    state.value = JSON.parse(a)
  } catch (e) {
    console.error(e)
  }
})
const error = ref("")
window.mtpAPI.onError((a) => {
  error.value = a
})

const workingState = computed(() => {
  if (!state.value) {
    return { color: "gray", text: "Not connected" }
  }

  if (state.value.Running) {
    return { color: "yellow", text: "Running" }
  }

  if (state.value.CPU.OUTP.INTA) {
    const signal = consts.SIGNALS[state.value.CPU.OUTP.INTN]
    if (!signal) {
      return {color: "red", text: `Interrupted: code = ${state.value.CPU.OUTP.INTN}` }
    }

    if (signal === "SIGTRACE") {
      return { color: "aqua", text: "Tracing" }
    }

    return { color: "red", text: `Interrupted (${signal})` }
  }

  return { color: "lime", text: "Waiting" }
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