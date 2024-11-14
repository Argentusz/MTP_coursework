<template>
  <div>
    <div style="display: flex; flex-direction: column">
      <span v-if="error" style="color:red">{{ error }}</span>
      State: {{ state }}
    </div>
    <div style="display: flex">
      <button @click="request('compile examples/dividers.mtp')">Compile</button>
      <button @click="request('run')">Run</button>
      <button @click="request('run all')">Run All</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";

const state = ref("")
window.mtpAPI.onUpdate((a) => {
  state.value = a
})
const error = ref("")
window.mtpAPI.onError((a) => {
  error.value = a
})


const request = (command) => {
  window.mtpAPI.request(command)
}
</script>