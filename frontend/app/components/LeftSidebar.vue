<template>
  <div class="w-64 bg-black border-r border-[#837dbd] flex flex-col overflow-hidden flex-shrink-10">
    <!-- Library Section -->
    <div class="p-4 border-b border-[#5a548d]">
      <div class="text-xs text-[#837dbd] uppercase tracking-wider mb-3 font-mono">LIBRARY</div>
      <div class="space-y-1">
        <button
          @click="$emit('change-view', 'library')"
          :class="navButtonClass('library')"
        >
          <span class="text-[#d3ceff]">[</span>Music<span class="text-[#d3ceff]">]</span>
        </button>
        <button
          @click="$emit('change-view', 'playlists')"
          :class="navButtonClass('playlists')"
        >
          <span class="text-[#d3ceff]">[</span>Playlists<span class="text-[#d3ceff]">]</span>
        </button>
        <button
          @click="$emit('change-view', 'search')"
          :class="navButtonClass('search')"
        >
          <span class="text-[#d3ceff]">[</span>Search<span class="text-[#d3ceff]">]</span>
        </button>
      </div>
    </div>

    <!-- Playlists Section -->
    <div class="flex-1 overflow-y-auto p-4">
      <div class="flex justify-between items-center mb-3">
        <div class="text-xs text-[#837dbd] uppercase tracking-wider font-mono">PLAYLISTS</div>
        <button
          @click="$emit('create-playlist')"
          class="text-[#d3ceff] hover:text-white transition-colors text-xs"
          title="Create new playlist"
        >
          [+]
        </button>
      </div>
      <div class="space-y-1">
        <button
          v-for="playlist in playlists"
          :key="playlist.id"
          @click="$emit('select-playlist', playlist.id)"
          class="w-full flex justify-between items-center hover:bg-[#1a1a1a] px-2 py-2 rounded cursor-pointer group transition-colors"
        >
          <span class="text-sm truncate text-[#d3ceff] font-mono">{{ playlist.name }}</span>
          <span class="text-xs text-[#5a548d] group-hover:text-[#837dbd] font-mono">
            {{ playlist.songs?.length || 0 }}
          </span>
        </button>
        <div v-if="playlists.length === 0" class="text-center text-[#5a548d] text-xs py-4 font-mono italic">
          No playlists yet
        </div>
      </div>
    </div>

    <!-- Now Playing Section -->
    <div class="border-t border-[#5a548d] p-4">
      <div class="text-xs text-[#837dbd] uppercase tracking-wider mb-3 font-mono">NOW PLAYING</div>
      <div v-if="playerStore.currentSong" class="space-y-2">
        <div class="text-sm font-bold text-[#d3ceff] truncate font-mono">
          {{ playerStore.currentSong.title }}
        </div>
        <div class="text-xs text-[#837dbd] truncate font-mono">
          {{ playerStore.currentSong.artist }}
        </div>
        <div class="w-full bg-[#333] h-1 rounded overflow-hidden">
          <div
            class="bg-[#d3ceff] h-1 rounded transition-all duration-300"
            :style="{ width: playerStore.progressPercentage + '%' }"
          ></div>
        </div>
      </div>
      <div v-else class="text-xs text-[#5a548d] italic font-mono text-center py-4">
        NO TRACK SELECTED
      </div>
    </div>

    <!-- ASCII Art Footer -->
    <div class="border-t border-[#5a548d] p-4">
      <pre class="text-[#5a548d] text-xs leading-3">
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
defineEmits(['change-view', 'create-playlist', 'select-playlist'])

const playerStore = usePlayerStore()

const navButtonClass = (view: string) => {
  const base = 'w-full text-left px-3 py-2 text-sm rounded hover:bg-[#1a1a1a] transition-colors font-mono'
  return props.currentView === view
    ? `${base} bg-[#1a1a1a] text-[#d3ceff] border-l-2 border-[#d3ceff]`
    : `${base} text-[#837dbd]`
}
</script>