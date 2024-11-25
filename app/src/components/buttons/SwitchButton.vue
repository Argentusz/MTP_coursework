<template>
  <v-button :disabled="disabled" theme="glass" @click="$emit('click')">
    <div class="switch-btn__content">
      <slot/>
      <v-led
        :color="on ? '#f86ea8' : '#000'"
        size="sm"
      />
    </div>
  </v-button>
</template>

<script setup>
import VButton from "./VButton.vue";
import VLed from "../ui/VLed.vue";

const $emit = defineEmits(["click"])

defineProps({
  on: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  }
})
</script>

<style scoped>
.switch-btn__content {
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
  align-items: center;
  justify-content: space-between;
  gap: 5px;
}
</style>