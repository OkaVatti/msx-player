<template>
  <div class="bg-gray-800 rounded-lg p-6">
    <h2 class="text-xl font-semibold mb-4">Control Center</h2>
    
    <!-- Current Song Info -->
    <div v-if="playerStore.currentSong" class="mb-6 p-4 bg-gray-700 rounded-lg">
      <div class="flex items-center space-x-4">
        <div class="flex-shrink-10 w-16 h-16 bg-purple-500 rounded-lg flex items-center justify-center">
          <span class="text-2xl">🎵</span>
        </div>
        <div class="flex-1">
          <h3 class="font-semibold text-lg">{{ playerStore.currentSong.title }}</h3>
          <p class="text-gray-400">{{ playerStore.currentSong.artist }} • {{ playerStore.currentSong.album }}</p>
        </div>
        <div class="text-right">
          <div class="text-2xl font-bold text-yellow-400">
            {{ playerStore.currentSong.rating }}/10
          </div>
          <div class="text-sm text-gray-400">Rating</div>
        </div>
      </div>
    </div>

    <!-- Progress Bar -->
    <div class="mb-6">
      <div class="flex justify-between text-sm text-gray-400 mb-2">
        <span>{{ formatTime(playerStore.currentTime) }}</span>
        <span>{{ formatTime(playerStore.duration) }}</span>
      </div>
      <div
        class="h-2 bg-gray-700 rounded-full cursor-pointer"
        @click="seekToTime"
        ref="progressBar"
      >
        <div
          class="h-full bg-linear-to-r from-purple-500 to-pink-500 rounded-full transition-all duration-100"
          :style="{ width: progressPercentage + '%' }"
        />
      </div>
    </div>

    <!-- Control Buttons -->
    <div class="flex items-center justify-center space-x-6 mb-6">
      <button
        @click="previousSong"
        class="p-3 rounded-full bg-gray-700 hover:bg-gray-600 transition-colors"
        title="Previous"
      >
        ⏮️
      </button>
      
      <button
        v-if="!playerStore.isPlaying"
        @click="playCurrent"
        class="p-4 rounded-full bg-purple-500 hover:bg-purple-600 transition-colors"
        title="Play"
      >
        ▶️
      </button>
      <button
        v-else
        @click="playerStore.pause()"
        class="p-4 rounded-full bg-pink-500 hover:bg-pink-600 transition-colors"
        title="Pause"
      >
        ⏸️
      </button>
      
      <button
        @click="nextSong"
        class="p-3 rounded-full bg-gray-700 hover:bg-gray-600 transition-colors"
        title="Next"
      >
        ⏭️
      </button>
      
      <button
        @click="playerStore.stop()"
        class="p-3 rounded-full bg-gray-700 hover:bg-gray-600 transition-colors"
        title="Stop"
      >
        ⏹️
      </button>
    </div>

    <!-- Additional Controls -->
    <div class="grid grid-cols-2 gap-4">
      <div class="space-y-2">
        <label class="text-sm text-gray-400">Volume</label>
        <input
          v-model="playerStore.volume"
          type="range"
          min="0"
          max="1"
          step="0.01"
          @input="playerStore.setVolume(playerStore.volume)"
          class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-purple-500"
        />
      </div>
      
      <div class="space-y-2">
        <label class="text-sm text-gray-400">Playback Speed</label>
        <select
          v-model="playerStore.speed"
          @change="playerStore.setSpeed(playerStore.speed)"
          class="w-full px-3 py-2 bg-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
        >
          <option value="0.5">0.5x</option>
          <option value="0.75">0.75x</option>
          <option value="1">1x</option>
          <option value="1.25">1.25x</option>
          <option value="1.5">1.5x</option>
          <option value="2">2x</option>
        </select>
      </div>
    </div>

    <!-- Rating Controls -->
    <div v-if="playerStore.currentSong" class="mt-6 p-4 bg-gray-700 rounded-lg">
      <label class="block text-sm text-gray-400 mb-3">Rate this song</label>
      <div class="flex justify-center space-x-1">
        <button
          v-for="star in 10"
          :key="star"
          @click="rateCurrentSong(star)"
          class="text-2xl transition-transform hover:scale-125"
          :class="star <= (playerStore.currentSong?.rating || 0) ? 'text-yellow-400' : 'text-gray-500'"
        >
          ★
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { usePlayerStore } from '../../stores/player'
import { useLibraryStore } from '../../stores/library'
import { ref, computed } from 'vue'

const playerStore = usePlayerStore()
const libraryStore = useLibraryStore()
const progressBar = ref<HTMLDivElement>()

const progressPercentage = computed(() => {
  if (playerStore.duration === 0) return 0
  return (playerStore.currentTime / playerStore.duration) * 100
})

const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const seekToTime = (event: MouseEvent) => {
  if (!progressBar.value) return
  
  const rect = progressBar.value.getBoundingClientRect()
  const percent = (event.clientX - rect.left) / rect.width
  const newTime = percent * playerStore.duration
  
  playerStore.seek(newTime)
}

const playCurrent = () => {
  if (playerStore.currentSong) {
    playerStore.playSong(playerStore.currentSong)
  }
}

const previousSong = () => {
  // Implementation would depend on your playlist logic
  console.log('Previous song')
}

const nextSong = () => {
  // Implementation would depend on your playlist logic
  console.log('Next song')
}

const rateCurrentSong = async (rating: number) => {
  if (playerStore.currentSong) {
    await playerStore.rateSong(playerStore.currentSong.id, rating)
    await libraryStore.fetchSongs()
  }
}
</script>