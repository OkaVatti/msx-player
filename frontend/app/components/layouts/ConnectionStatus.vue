<template>
  <div class="flex items-center gap-2 font-mono text-sm">
    <div :class="['w-2 h-2 rounded-full', connected ? 'bg-green-400 animate-pulse' : 'bg-gray-700']"></div>
    <div :class="connected ? 'text-green-400' : 'text-[#837dbd]'">{{ connected ? 'ONLINE' : 'OFFLINE' }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const connected = ref(false)
let timer: number | null = null

const check = async () => {
  try {
    const r = await fetch('/api/player/state', { cache: 'no-store' })
    connected.value = r.ok
  } catch {
    connected.value = false
  }
}

onMounted(() => {
  check()
  timer = window.setInterval(check, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
</style>
