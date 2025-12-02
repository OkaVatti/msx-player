<template>
  <div class="h-20 border-t border-[#837dbd] bg-black flex items-center px-4 gap-4">
    <div class="flex items-center gap-3 min-w-0 w-72">
      <div class="w-14 h-14 bg-gradient-to-br from-[#5a548d] to-[#837dbd] flex items-center justify-center">
        <span class="text-2xl">♪</span>
      </div>
      <div class="min-w-0">
        <div class="text-[#d3ceff] truncate">{{ playerStore.currentSong?.title || 'No song selected' }}</div>
        <div class="text-[#837dbd] text-xs truncate">{{ playerStore.currentSong?.artist || '' }}</div>
      </div>
    </div>

    <div class="flex-1 flex flex-col items-center">
      <div class="flex items-center gap-4">
        <button @click="playerStore.previousSong()" :disabled="!playerStore.hasPrevious" class="text-[#837dbd]">⏮</button>
        <button v-if="!playerStore.isPlaying" @click="play" class="w-10 h-10 bg-[#d3ceff] text-black rounded-full">▶</button>
        <button v-else @click="playerStore.pause()" class="w-10 h-10 bg-[#d3ceff] text-black rounded-full">⏸</button>
        <button @click="playerStore.nextSong()" :disabled="!playerStore.hasNext" class="text-[#837dbd]">⏭</button>

        <div class="w-96 flex items-center gap-3">
          <span class="text-xs text-[#5a548d]">{{ formatTime(playerStore.currentTime) }}</span>
          <div class="flex-1 h-1 bg-[#5a548d] rounded relative" @click="seekBar">
            <div class="h-full bg-[#d3ceff] rounded" :style="{ width: playerStore.progressPercentage + '%' }"></div>
          </div>
          <span class="text-xs text-[#5a548d]">{{ formatTime(playerStore.duration) }}</span>
        </div>
      </div>
    </div>

    <div class="flex items-center gap-3 w-44">
      <input type="range" step="0.01" v-model.number="volume" @input="onVolume" min="0" max="1" />
      <div class="text-xs text-[#5a548d]">{{ Math.round(volume*100) }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { usePlayerStore } from '@/stores/player'

const playerStore = usePlayerStore()
const volume = ref(playerStore.volume)
watch(() => playerStore.volume, (v) => volume.value = v)

const onVolume = () => playerStore.setVolume(volume.value)
const formatTime = (s: number) => {
  if (!s || isNaN(s)) return '0:00'
  const m = Math.floor(s/60); const secs = Math.floor(s%60).toString().padStart(2,'0')
  return `${m}:${secs}`
}
const play = async () => {
  if (playerStore.currentSong) await playerStore.playSong(playerStore.currentSong)
}

const seekBar = (e: MouseEvent) => {
  const el = e.currentTarget as HTMLElement
  const rect = el.getBoundingClientRect()
  const percent = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width))
  playerStore.seek(percent * playerStore.duration)
}
</script>

<style scoped>
input[type="range"] { accent-color: #d3ceff; }
</style>
