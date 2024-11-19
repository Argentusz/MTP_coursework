<template>
  <div class="status-bar"
    :style="{ borderBottomColor: workingState.color }"
  >
    <div class="status-bar__side">
      <div class="status-bar__led"
           :style="{backgroundColor: workingState.color}"
      />
      {{ workingState.text }}
    </div>
    <div class="status-bar__side">
      <button :disabled="!filePath" @click="request(`compile ${filePath}`)">Compile</button>
      <button @click="request('run')">Run</button>
      <button @click="request('run all')">Run All</button>
      <button @click="request('term')">Terminate</button>
      <button @click="request('reset')">Reset</button>
    </div>
  </div>

</template>

<script setup>
defineProps({
  workingState: {
    type: Object,
    default: () => {}
  },
  filePath: {
    type: [String, null],
    default: null
  }
})

const $emit = defineEmits(["request"])

const request = (cmd) => $emit("request", cmd)
</script>

<style scoped>
.status-bar {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  height: 75px;
  padding: 0 10px;
  background-color: #242021;
  border-bottom: 2px solid;
}

.status-bar__side {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 10px;
}

.status-bar__led {
  width: 10px;
  height: 10px;
  border-radius: 10px;
}
</style>