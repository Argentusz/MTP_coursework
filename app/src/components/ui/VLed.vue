<template>
  <div
    :class="['v-led', `v-led_${size}`]"
    :style="{backgroundColor: color}"
  />
</template>

<script setup>

defineProps({
  color: {
    type: String,
    default: "#000"
  },
  size: {
    type: String,
    default: "md"
  }
})

</script>

<style scoped>
.v-led {
  width: 10px;
  height: 10px;
  border-radius: 10px;
}

.v-led_sm {
  width: 7px;
  height: 7px;
}

.v-led_xs {
  width: 5px;
  height: 5px;
}
</style>