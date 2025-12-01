<!-- components/LeftSidebar.vue -->
<template>
  <div class="w-64 bg-black border-r border-[#837dbd] p-4 overflow-y-auto">
    <!-- Library Section -->
    <div class="mb-6">
      <div class="text-xs text-[#837dbd] uppercase tracking-wider mb-2">LIBRARY</div>
      <div class="space-y-1">
        <button @click="emit('change-view', 'library')" :class="sidebarButtonClass('library')">
          <span class="text-[#d3ceff]">[</span>Music<span class="text-[#d3ceff]">]</span>
        </button>
        <button @click="emit('change-view', 'playlists')" :class="sidebarButtonClass('playlists')">
          <span class="text-[#d3ceff]">[</span>Playlists<span class="text-[#d3ceff]">]</span>
        </button>
        <button @click="emit('change-view', 'recent')" :class="sidebarButtonClass('recent')">
          <span class="text-[#d3ceff]">[</span>Recently Added<span class="text-[#d3ceff]">]</span>
        </button>
        <button @click="emit('change-view', 'top')" :class="sidebarButtonClass('top')">
          <span class="text-[#d3ceff]">[</span>Top Rated<span class="text-[#d3ceff]">]</span>
        </button>
        <button @click="emit('change-view', 'search')" :class="sidebarButtonClass('search')">
          <span class="text-[#d3ceff]">[</span>Search<span class="text-[#d3ceff]">]</span>
        </button>
        <button @click="emit('change-view', 'upload')" :class="sidebarButtonClass('upload')">
          <span class="text-[#d3ceff]">[</span>Upload<span class="text-[#d3ceff]">]</span>
        </button>
        <button @click="emit('change-view', 'manage')" :class="sidebarButtonClass('manage')">
          <span class="text-[#d3ceff]">[</span>Manage<span class="text-[#d3ceff]">]</span>
        </button>
        <button @click="emit('change-view', 'visualizer')" :class="sidebarButtonClass('visualizer')">
          <span class="text-[#d3ceff">[</span>Visualizer<span class="text-[#d3ceff">]</span>
        </button>
        <button @click="emit('change-view', 'controls')" :class="sidebarButtonClass('controls')">
          <span class="text-[#d3ceff">[</span>Controls<span class="text-[#d3ceff">]</span>
        </button>
      </div>
    </div>

    <!-- Playlists Section -->
    <div class="mb-6">
      <div class="text-xs text-[#837dbd] uppercase tracking-wider mb-2 flex justify-between">
        <span>PLAYLISTS</span>
        <button @click="emit('create-playlist')" class="text-[#d3ceff] hover:text-white text-xs">[+]</button>
      </div>
      <div class="space-y-1">
        <div v-for="playlist in playlists" :key="playlist.id" 
             @click="emit('select-playlist', playlist.id)"
             class="flex justify-between items-center hover:bg-[#1a1a1a] px-2 py-1 rounded cursor-pointer group">
          <span class="text-sm truncate text-[#d3ceff]">{{ playlist.name }}</span>
          <span class="text-xs text-[#837dbd] group-hover:opacity-100 opacity-0">{{ playlist.songs.entries || 0 }}</span>
        </div>
      </div>
    </div>

    <!-- Now Playing -->
    <div class="mb-6 border border-[#837dbd] p-3 bg-[#111]">
      <div class="text-xs text-[#837dbd] uppercase tracking-wider mb-2">NOW PLAYING</div>
      <div v-if="currentSong" class="space-y-2">
        <div class="text-sm font-bold text-[#d3ceff] truncate">{{ currentSong.title }}</div>
        <div class="text-xs text-[#837dbd] truncate">{{ currentSong.artist }}</div>
        <div class="w-full bg-[#333] h-1 rounded">
          <div class="bg-[#d3ceff] h-1 rounded" :style="{ width: progressPercentage + '%' }"></div>
        </div>
      </div>
      <div v-else class="text-xs text-[#837dbd] italic">
        NO TRACK SELECTED
      </div>
    </div>

    <!-- ASCII Art Section -->
    <div class="mt-8 pt-6 border-t border-[#837dbd]">
      <pre class="text-[#d3ceff] text-xs leading-3">
╔══════════════════╗
║   MSX PLAYER     ║
║   v1.3.37        ║
║   [ACTIVE]       ║
╚══════════════════╝</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerStore } from '../../stores/player'
import type { Playlist } from '../../types'

interface Props {
  currentView: string
  playlists: Playlist[]
}

const props = defineProps<Props>()
const emit = defineEmits(['change-view', 'create-playlist', 'select-playlist'])

const playerStore = usePlayerStore()

const currentSong = computed(() => playerStore.currentSong)
const progressPercentage = computed(() => playerStore.progressPercentage)

const sidebarButtonClass = (view: string) => {
  const base = 'w-full text-left px-3 py-2 text-sm rounded hover:bg-[#1a1a1a] transition-colors'
  return props.currentView === view 
    ? `${base} bg-[#1a1a1a] text-[#d3ceff] border-l-2 border-[#d3ceff]`
    : `${base} text-[#837dbd]`
}
</script>