<template>
  <div class="flex items-center gap-3">
    <button @click="player.previousSong()" :disabled="!player.hasPrevious" class="text-[#837dbd]">⏮</button>
    <button v-if="!player.isPlaying" @click="play" class="w-10 h-10 bg-[#d3ceff] text-black rounded-full">▶</button>
    <button v-else @click="player.pause" class="w-10 h-10 bg-[#d3ceff] text-black rounded-full">⏸</button>
    <button @click="player.nextSong()" :disabled="!player.hasNext" class="text-[#837dbd]">⏭</button>
  </div>
</template>

<script setup lang="ts">
import { usePlayerStore } from '../../../stores/player'
const player = usePlayerStore()
const play = async () => {
  if (player.currentSong) await player.playSong(player.currentSong)
}
</script>
