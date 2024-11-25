<template>
  <button
      :class="[
        'v-button',
        `v-button_${theme}`,
        `v-button_${size}`,
        {'v-button_disabled': disabled},
      ]"
      :disabled="disabled"
  >
    <slot/>
  </button>
</template>

<script setup>
const $props = defineProps({
  disabled: {
    type: Boolean,
    default: false
  },
  theme: {
    type: String,
    default: "default"
  },
  size: {
    type: String,
    default: "md"
  }
})
</script>

<style lang="scss" scoped>
.v-button {
  background: none;
  border: none;
  padding: 5px;
  border-radius: 3px;
  min-width: 60px;
  opacity: 0.8;
  &:hover {
    opacity: 1;
  }
}

.v-button_disabled {
  background-color: gray !important;
  color: white;
  opacity: 0.3 !important;
}

.v-button_green {
  background-color: olivedrab;
  color: white;
}

.v-button_red {
  background-color: red;
  color: white;
}

.v-button_pink {
  background-color: #f86ea8;
  color: white;
}

.v-button_glass {
  background-color: rgba(255, 255, 255, 0.1);
  color: white;
  &:hover {
    background-color: rgba(255, 255, 255, 0.3);
  }
}

.v-button_lg {
  font-size: 1rem;
  min-width: 70px;
  font-weight: bold;
}
</style>