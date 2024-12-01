<template>
  <transition-group name="list">
    <div v-if="error" class="error-flash">
      <span>{{ error }}</span>
      <img width="15" src="../../assets/svg/close.svg" class="error-close" @click="error = null"/>
    </div>
  </transition-group>
</template>

<script setup>
const error = defineModel({ type: [String, null], default: null })
</script>

<style scoped>
.error-flash {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  position: absolute;
  bottom: 0;
  right: 0;
  background-color: #ff033e;
  opacity: 0.95;
  color: #191516;
  padding: 10px 20px;
  gap: 20px;
  width: 33vw;
}

.error-close {
  cursor: pointer;
}

.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
}
</style>