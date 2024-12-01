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
      <v-button theme="green" :disabled="running || !connected" @click="$emit('send')">Send</v-button>
      <v-button theme="red" :disabled="running || !connected" @click="$emit('drop')">Drop</v-button>
      <v-led color="gray" size="xs"/>
      <switch-button :disabled="running || !connected" :on="intr" @click="request(`intr ${intr ? 'off' : 'on'}`)">Interrupts</switch-button>
      <switch-button :disabled="running || !connected" :on="trace" @click="request(`trace ${trace ? 'off' : 'on'}`)">Trace</switch-button>
      <switch-button :disabled="running || !connected" :on="sudo" @click="request(`sudo ${sudo ? 'off' : 'on'}`)">Sudo</switch-button>
      <v-led color="gray" size="xs"/>
      <v-button :disabled="running || !connected || !filePath" theme="pink" @click="compile">Compile</v-button>
      <v-button :disabled="running || !connected" theme="pink" @click="run(false)">Run</v-button>
      <v-button :disabled="running || !connected" theme="pink" @click="run(true)">Run All</v-button>
      <v-led color="gray" size="xs"/>
      <v-button :disabled="running || !connected" theme="red" @click="request('reset')">Reset</v-button>
      <v-button :disabled="!running || !connected" theme="red" @click="request('term')">Terminate</v-button>
      <v-led color="gray" size="xs"/>
      <v-button v-if="!connected" theme="green" @click="$emit('connect')">Connect</v-button>
      <v-button v-else theme="red" @click="$emit('disconnect')">Disconnect</v-button>
    </div>
  </div>

</template>

<script setup>
import VButton from "./buttons/VButton.vue";
import SwitchButton from "./buttons/SwitchButton.vue";
import VLed from "./ui/VLed.vue";
import { computed, ref } from "vue";
import consts from "../js/consts";

const $log = () => console.log($props.state)

const $props = defineProps({
  state: {
    type: [Object, null],
    default: null,
  },
  filePath: {
    type: [String, null],
    default: null,
  },
})

const $emit = defineEmits(["request", "compile", "connect", "disconnect", "run", "send", "drop"])

const request = (cmd) => $emit("request", cmd)
const compile = () => $emit("compile")
const run = (all = false) => $emit("run", all)

const running = computed(() => $props.state?.Running || false)
const connected = computed(() => !!$props.state)

const intr = computed(() => connected.value && isFlagOn(0b10000))
const trace = computed(() => connected.value && isFlagOn(0b100000))
const sudo = computed(() => connected.value && isFlagOn(0b1000000))
const isFlagOn = mask => (($props.state.CPU.RRAM.SYS.FLG || $props.state.CPU.RRAM.SYS.FLB) & mask) === mask

const workingState = computed(() => {
  if (!$props.state)
    return { color: "gray", text: "Not connected" }
  if ($props.state.Running)
    return { color: "yellow", text: "Running" }
  if (!$props.state.CPU.OUTP.INTA)
    return { color: "lime", text: "Waiting" }

  const signal = consts.SIGNALS[$props.state.CPU.OUTP.INTN]
  if (signal === "SIGTRACE")
    return { color: "aqua", text: "Tracing" }

  return { color: "red", text: `Interrupted (${signal || $props.state.CPU.OUP.INTN})` }
})

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