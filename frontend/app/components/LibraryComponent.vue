<template>
  <div class="bg-gray-800 rounded-lg p-6 max-h-[600px] overflow-hidden flex flex-col">
    <h2 class="text-xl font-semibold mb-4">Your Library</h2>
    
    <div class="flex-1 overflow-y-auto">
      <div class="space-y-2">
        <div
          v-for="song in libraryStore.filteredSongs"
          :key="song.id"
          class="flex items-center space-x-4 p-3 bg-gray-700 rounded-lg hover:bg-gray-600 transition-colors cursor-pointer group"
          @click="playSong(song)"
        >
          <div class="flex-shrink-10 w-12 h-12 bg-purple-500 rounded-lg flex items-center justify-center">
            <span class="text-lg">🎵</span>
          </div>
          
          <div class="flex-1 min-w-0">
            <h3 class="font-medium truncate" :class="{ 'text-purple-400': playerStore.currentSong?.id === song.id }">
              {{ song.title }}
            </h3>
            <p class="text-sm text-gray-400 truncate">{{ song.artist }} • {{ song.album }}</p>
          </div>
          
          <div class="flex items-center space-x-3 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              v-for="rating in [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]"
              :key="rating"
              @click.stop="rateSong(song.id, rating)"
              class="text-xs hover:text-yellow-400 transition-colors"
              :class="rating <= song.rating ? 'text-yellow-400' : 'text-gray-500'"
            >
              ★
            </button>
            
            <button
              @click.stop="deleteSong(song.id)"
              class="text-red-400 hover:text-red-300 transition-colors"
            >
              🗑️
            </button>
          </div>
          
          <div class="text-right text-sm text-gray-400">
            <div class="flex items-center space-x-1">
              <span>★</span>
              <span>{{ song.rating }}/10</span>
            </div>
            <div class="text-xs">{{ formatDuration(song.duration) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { usePlayerStore } from '../../stores/player'
import { useLibraryStore } from '../../stores/library'
import { onMounted } from 'vue'

const playerStore = usePlayerStore()
const libraryStore = useLibraryStore()

onMounted(() => {
  libraryStore.fetchSongs()
})

const playSong = (song: any) => {
  playerStore.playSong(song)
}

const rateSong = async (songId: number, rating: number) => {
  await playerStore.rateSong(songId, rating)
  await libraryStore.fetchSongs()
}

const deleteSong = async (songId: number) => {
  if (confirm('Are you sure you want to delete this song from your library?')) {
    await libraryStore.deleteSong(songId)
  }
}

const formatDuration = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}
</script>