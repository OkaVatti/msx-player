<template>
  <div class="p-4 h-full">
    <div class="text-xs text-[#837dbd] uppercase mb-3">Now Playing</div>
    <div v-if="playerStore.currentSong" class="space-y-3">
      <div class="text-[#d3ceff] font-bold">{{ playerStore.currentSong.title }}</div>
      <div class="text-[#837dbd]">{{ playerStore.currentSong.artist }}</div>
      <div class="text-xs text-[#837dbd] mt-2">RATING</div>
      <div class="flex gap-1">
        <button v-for="i in 10" :key="i" @click="rate(i)" :class="i <= (playerStore.currentSong?.rating || 0) ? 'text-yellow-400' : 'text-[#333]'">★</button>
      </div>
    </div>
    <div v-else class="text-[#837dbd]">No song playing</div>

    <div class="mt-6">
      <div class="text-xs text-[#837dbd] uppercase mb-2">Up Next</div>
      <div v-for="s in upNext" :key="s.id" class="flex items-center justify-between py-1">
        <div>
          <div class="text-sm text-[#d3ceff] truncate">{{ s.title }}</div>
          <div class="text-xs text-[#837dbd] truncate">{{ s.artist }}</div>
        </div>
        <button @click="play(s)" class="text-[#837dbd]">▶</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerStore } from '@/stores/player'
import { useLibraryStore } from '@/stores/library'
const playerStore = usePlayerStore()
const lib = useLibraryStore()

const upNext = computed(() => {
  if (!playerStore.currentSong) return lib.songs.slice(0,6)
  const idx = lib.songs.findIndex(s=>s.id===playerStore.currentSong?.id)
  if (idx<0) return lib.songs.slice(0,6)
  return lib.songs.slice(idx+1, idx+7)
})

const play = async (s:any) => await playerStore.playSong(s)
const rate = async (r:number) => {
  if (!playerStore.currentSong) return
  await playerStore.rateSong(playerStore.currentSong.id, r)
}
</script>
