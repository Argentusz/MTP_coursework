<template>
  <button
      :class="['tab-button', { 'tab-button_active': active }]"
      @click="$emit('click')"
  >
    {{ text }}
  </button>
</template>

<script setup>
defineProps({
  text: {
    type: String,
    default: "Click"
  },
  active: {
    type: Boolean,
    default: false
  }
})

const $emit = defineEmits(["click"])
</script>

<style scoped>
.tab-button {
  background: none;
  color: inherit;
  border: none;
  padding: 0;
  font: inherit;
  cursor: pointer;
  outline: inherit;
}

.tab-button_active {
  color: #f86ea8;
}
</style>