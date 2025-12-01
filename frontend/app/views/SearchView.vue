<!-- views/SearchView.vue -->
<template>
  <div class="space-y-6">
    <!-- Search Header -->
    <div class="border border-[#837dbd] p-4 bg-black">
      <h2 class="text-xl mb-2 flex items-center gap-2">
        <span class="text-[#d3ceff]">></span>
        <span class="font-bold">SEARCH SYSTEM</span>
        <span class="text-xs text-[#837dbd] ml-2">[REAL-TIME FILTERING]</span>
      </h2>
      
      <!-- Search Input -->
      <div class="relative">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="ENTER SEARCH QUERY..."
          class="w-full bg-black border border-[#837dbd] text-white p-3 pl-10 focus:outline-none focus:border-[#d3ceff] font-mono"
          @input="handleSearch"
        />
        <div class="absolute left-3 top-1/2 transform -translate-y-1/2 text-[#837dbd]">⌕</div>
        
        <div v-if="libraryStore.loading" 
             class="absolute right-3 top-1/2 transform -translate-y-1/2 text-[#837dbd] animate-pulse font-mono text-xs">
          [SEARCHING...]
        </div>
        <div v-else-if="searchQuery"
             class="absolute right-3 top-1/2 transform -translate-y-1/2 text-[#d3ceff] font-mono text-xs">
          {{ searchResults.length }}
        </div>
      </div>
    </div>

    <!-- Search Stats -->
    <div class="grid grid-cols-3 gap-4">
      <div class="border border-[#837dbd] p-4 text-center">
        <div class="text-xs text-[#837dbd]">RESULTS</div>
        <div class="text-2xl text-[#d3ceff]">{{ searchResults.length }}</div>
      </div>
      <div class="border border-[#837dbd] p-4 text-center">
        <div class="text-xs text-[#837dbd]">QUERY</div>
        <div class="text-lg text-[#d3ceff] truncate">"{{ searchQuery || 'NONE' }}"</div>
      </div>
      <div class="border border-[#837dbd] p-4 text-center">
        <div class="text-xs text-[#837dbd]">TIME</div>
        <div class="text-xl text-[#d3ceff]">{{ new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }}</div>
      </div>
    </div>

    <!-- Results Table -->
    <div class="border border-[#837dbd]">
      <!-- Table Header -->
      <div class="grid grid-cols-12 border-b border-[#837dbd] text-xs text-[#837dbd] uppercase tracking-wider">
        <div class="col-span-1 p-3 text-center">#</div>
        <div class="col-span-5 p-3">TRACK</div>
        <div class="col-span-3 p-3">ARTIST</div>
        <div class="col-span-2 p-3">ALBUM</div>
        <div class="col-span-1 p-3 text-right">TIME</div>
      </div>

      <!-- Results -->
      <div v-for="(song, index) in searchResults" :key="song.id"
           @click="playSong(song, index)"
           :class="['grid grid-cols-12 border-b border-[#222] hover:bg-[#1a1a1a] cursor-pointer', 
                   playerStore.currentSong?.id === song.id ? 'bg-[#1a1a1a]' : '']">
        <div class="col-span-1 p-3 text-center text-[#837dbd]">
          {{ playerStore.currentSong?.id === song.id ? '▶' : index + 1 }}
        </div>
        <div class="col-span-5 p-3">
          <div class="flex items-center gap-2">
            <span class="text-[#d3ceff]">{{ song.title }}</span>
            <span v-if="song.explicit" class="text-xs text-red-400 border border-red-400 px-1">E</span>
            <span v-else-if="song.clean" class="text-xs text-blue-400 border border-blue-400 px-1">C</span>
          </div>
        </div>
        <div class="col-span-3 p-3 text-[#837dbd]">{{ song.artist }}</div>
        <div class="col-span-2 p-3 text-[#837dbd]">{{ song.album }}</div>
        <div class="col-span-1 p-3 text-right text-[#837dbd]">{{ formatDuration(song.duration) }}</div>
      </div>

      <!-- Empty States -->
      <div v-if="searchQuery && searchResults.length === 0 && !libraryStore.loading" 
           class="p-12 text-center">
        <pre class="text-[#837dbd] text-sm">
╔═══════════════════════════╗
║   NO RESULTS FOUND        ║
║   TRY A DIFFERENT QUERY   ║
╚═══════════════════════════╝</pre>
      </div>

      <div v-if="!searchQuery" class="p-12 text-center">
        <pre class="text-[#837dbd] text-sm">
╔═══════════════════════════╗
║   ENTER SEARCH QUERY      ║
║   ABOVE TO BEGIN          ║
╚═══════════════════════════╝</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { usePlayerStore } from '../../stores/player'
import { useLibraryStore } from '../../stores/library'
import type { Song } from '../../types'

const playerStore = usePlayerStore()
const libraryStore = useLibraryStore()

const searchQuery = ref('')
let searchTimeout: NodeJS.Timeout | null = null

const searchResults = computed(() => libraryStore.filteredSongs)

const handleSearch = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  
  searchTimeout = setTimeout(() => {
    libraryStore.searchQuery = searchQuery.value
    libraryStore.fetchSongs()
  }, 300)
}

const formatDuration = (seconds: number) => {
  if (!seconds || seconds === 0) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const playSong = async (song: Song, index: number) => {
  playerStore.setQueue(searchResults.value, index)
  await playerStore.playSong(song)
}

// Clear timeout on component unmount
import { onUnmounted } from 'vue'
onUnmounted(() => {
  if (searchTimeout) clearTimeout(searchTimeout)
})
</script>../