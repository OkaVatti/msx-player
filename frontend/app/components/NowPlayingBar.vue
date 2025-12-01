<template>
  <div class="h-12 border-t border-[#837dbd] bg-black flex items-center px-4">
    <div class="flex items-center gap-4 flex-1">
      <!-- Song Info -->
      <div v-if="playerStore.currentSong" class="flex items-center gap-3 min-w-0 flex-1">
        <div class="w-8 h-8 bg-[#5a548d] flex items-center justify-center flex-shrink-0">
          <span class="text-[#d3ceff]">♪</span>
        </div>
        <div class="min-w-0 flex-1">
          <div class="text-[#d3ceff] text-sm font-mono truncate">{{ playerStore.currentSong.title }}</div>
          <div class="text-[#837dbd] text-xs font-mono truncate">{{ playerStore.currentSong.artist }}</div>
        </div>
      </div>
      <div v-else class="text-[#5a548d] text-sm font-mono">No song selected</div>
      
      <!-- Mini Player Controls -->
      <div class="flex items-center gap-3">
        <button
          @click="playerStore.previousSong()"
          class="text-[#837dbd] hover:text-[#d3ceff] transition-colors"
          title="Previous"
        >
          ⏮
        </button>
        <button
          @click="togglePlay"
          class="text-[#d3ceff] hover:text-white transition-colors"
          title="Play/Pause"
        >
          {{ playerStore.isPlaying ? '⏸' : '▶' }}
        </button>
        <button
          @click="playerStore.nextSong()"
          class="text-[#837dbd] hover:text-[#d3ceff] transition-colors"
          title="Next"
        >
          ⏭
        </button>
      </div>
      
      <!-- Progress Bar -->
      <div v-if="playerStore.currentSong" class="flex-1 max-w-md">
        <div class="h-1 bg-[#5a548d] rounded-full overflow-hidden cursor-pointer" @click="seekToTime">
          <div 
            class="h-full bg-[#d3ceff] transition-all duration-300"
            :style="{ width: playerStore.progressPercentage + '%' }"
          ></div>
        </div>
        <div class="flex justify-between text-[#5a548d] text-xs mt-1">
          <span>{{ formatTime(playerStore.currentTime) }}</span>
          <span>{{ formatTime(playerStore.duration) }}</span>
        </div>
      </div>
    </div>
    
    <!-- Volume Control -->
    <div class="flex items-center gap-2 w-32">
      <span class="text-[#837dbd] text-xs font-mono">VOL</span>
      <input
        type="range"
        v-model="playerStore.volume"
        min="0"
        max="1"
        step="0.01"
        @change="playerStore.setVolume(playerStore.volume)"
        class="flex-1"
      />
      <span class="text-[#5a548d] text-xs font-mono w-8">
        {{ Math.round(playerStore.volume * 100) }}%
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { usePlayerStore } from '../../stores/player'

const playerStore = usePlayerStore()

const togglePlay = () => {
  if (playerStore.currentSong) {
    if (playerStore.isPlaying) {
      playerStore.pause()
    } else {
      playerStore.playSong(playerStore.currentSong)
    }
  }
}

const formatTime = (seconds: number) => {
  if (!seconds || isNaN(seconds)) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const seekToTime = (event: MouseEvent) => {
  if (!playerStore.currentSong) return
  
  const element = event.currentTarget as HTMLElement
  const rect = element.getBoundingClientRect()
  const percent = (event.clientX - rect.left) / rect.width
  const newTime = percent * playerStore.duration
  playerStore.seek(newTime)
}
</script>