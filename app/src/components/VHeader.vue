<template>
  <div class="header__container">
    <div class="header__front">
      <img src="../../assets/img/mtp_logo.png" height="50">
      <v-button disabled theme="glass" size="lg" @click="$emit('new')">New</v-button>
      <v-button theme="glass" size="lg" @click="$emit('open')">Open</v-button>
      <v-button theme="glass" size="lg" :disabled="!filePath" @click="$emit('save')">Save</v-button>
      {{ fileName }}
    </div>
    <div class="header__glass"></div>
    <div class="header__gradient">
      <img src="../../assets/img/gradient.png" height="75" width="100%">
    </div>

  </div>
</template>

<script setup>
import VButton from "./buttons/VButton.vue";
import { computed } from "vue";

const $emit = defineEmits(["open", "save", "new"])

const $props = defineProps({
  filePath: {
    type: [String, null],
    default: false,
  }
})

const fileName = computed(() => {
  if (!$props.filePath) return ""
  const split = $props.filePath.split("/")
  if (!split || split.length === 0) return ""
  return split[split.length - 1]
})
</script>

<style scoped>
.header__container {
  height: 75px;
}

.header__front {
  position: absolute;
  z-index: 30;
  height: 75px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 30px;
}

.header__glass {
  position: absolute;
  z-index: 20;
  width: 100%;
  height: 75px;
  background: rgba(50, 50, 50, 0.2);
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

.header__gradient {
  position: absolute;
  z-index: 10;
  width: 100%;
  height: 75px;
}
</style>