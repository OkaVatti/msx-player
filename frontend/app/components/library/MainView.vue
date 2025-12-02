<template>
  <div class="p-6">
    <div v-if="view === 'library'">
      <div class="mb-4 flex justify-between items-center">
        <h2 class="text-xl text-[#d3ceff]">Library</h2>
        <div class="flex gap-2">
          <input v-model="filters.search" placeholder="Search" class="bg-black border border-[#837dbd] px-2 py-1" />
          <select v-model="filters.sortBy" class="bg-black border border-[#837dbd] px-2 py-1">
            <option value="title">Title</option>
            <option value="artist">Artist</option>
            <option value="album">Album</option>
            <option value="rating">Rating</option>
            <option value="play_count">Play Count</option>
            <option value="last_played">Last Played</option>
            <option value="duration">Duration</option>
            <option value="year">Year</option>
          </select>
        </div>
      </div>

      <div>
        <div v-for="s in library.filteredSongs" :key="s.id" class="p-3 border-b border-[#1a1a1a] flex items-center justify-between">
          <div>
            <div class="text-[#d3ceff] font-mono">{{ s.title }}</div>
            <div class="text-[#837dbd] text-xs">{{ s.artist }} • {{ s.album }}</div>
          </div>
          <div class="flex items-center gap-3">
            <button @click="play(s)" class="text-[#837dbd]">Play</button>
            <div class="text-xs text-[#5a548d]">{{ formatDuration(s.duration) }}</div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="view === 'visualizer'">
      <Visualizer style="height: 70vh" />
    </div>

    <div v-else-if="view === 'playlists'">
      <h2 class="text-xl text-[#d3ceff] mb-4">Playlists</h2>
      <div v-for="pl in library.playlists" :key="pl.id" class="p-3 border border-[#837dbd] mb-2">
        <div class="flex justify-between items-center">
          <div>
            <div class="text-[#d3ceff]">{{ pl.name }}</div>
            <div class="text-[#837dbd] text-sm">{{ pl.description }}</div>
          </div>
          <div>
            <button @click="open(pl.id)" class="text-[#837dbd]">Open</button>
          </div>
        </div>
        <div class="mt-2 grid grid-cols-2 gap-2">
          <div v-for="s in pl.songs" :key="s.id" class="p-2 border border-[#1a1a1a]">
            <div class="text-[#d3ceff]">{{ s.title }}</div>
            <div class="text-[#837dbd] text-xs">{{ s.artist }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Visualizer from '~/components/player/Visualizer.vue'
import { useLibraryStore } from '@/stores/library'
import { usePlayerStore } from '@/stores/player'
import type { Song } from '@/types'

const props = defineProps<{ view: 'library' | 'playlists' | 'visualizer' }>()
const library = useLibraryStore()
const player = usePlayerStore()

const filters = library.filters

const play = async (s: Song) => {
  await player.playSong(s)
}

const formatDuration = (s:number) => {
  if (!s) return '0:00'
  const m = Math.floor(s/60); const sec = Math.floor(s%60).toString().padStart(2,'0')
  return `${m}:${sec}`
}

const open = (id:number) => {
  library.selectedPlaylistId = id
}
</script>
