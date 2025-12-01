<template>
  <div class="space-y-4">
    <div class="border-2 border-lime-400 p-4 bg-black">
      <h2 class="text-xl mb-4 flex items-center gap-2">
        <span class="animate-pulse">&gt;&gt;</span>
        SEARCH SYSTEM
        <span class="text-xs text-white">[REAL-TIME FILTERING]</span>
      </h2>

      <div class="relative">
        <input
          v-model="libraryStore.searchQuery"
          type="text"
          placeholder="ENTER SEARCH QUERY..."
          class="w-full bg-black border-2 border-lime-400 text-lime-400 p-3 focus:outline-none focus:bg-lime-400 focus:text-black font-mono"
          @input="handleSearch"
        />
        <div
          v-if="libraryStore.searchQuery && libraryStore.loading"
          class="absolute right-2 top-1/2 transform -translate-y-1/2 text-lime-400 animate-pulse font-mono"
        >
          [SEARCHING...]
        </div>
        <div
          v-else-if="libraryStore.searchQuery"
          class="absolute right-2 top-1/2 transform -translate-y-1/2 text-white font-mono"
        >
          [{{ searchResults.length }}]
        </div>
      </div>

      <div class="mt-4 text-sm text-lime-400 font-mono">
        <span>RESULTS: {{ searchResults.length }}</span>
        <span class="mx-2">|</span>
        <span>QUERY: "{{ libraryStore.searchQuery || 'NONE' }}"</span>
      </div>
    </div>

    <!-- Search Results -->
    <div class="border-2 border-lime-400 p-4 bg-black max-h-[500px] overflow-y-auto">
      <div
        v-for="song in searchResults"
        :key="song.id"
        class="mb-2 border border-lime-400 p-3 hover:bg-lime-400 hover:text-black transition-all cursor-pointer group"
        @click="playSong(song)"
      >
        <div class="flex justify-between items-center">
          <div class="flex-1 min-w-0">
            <div class="font-bold mb-1 flex items-center gap-2">
              {{ song.title }}
              <span v-if="song.explicit" class="text-red-500 text-xs border border-red-500 px-1">
                EXPLICIT
              </span>
              <span v-else-if="song.clean" class="text-blue-500 text-xs border border-blue-500 px-1">
                CLEAN
              </span>
            </div>
            <div class="text-sm opacity-80">
              {{ song.artist }} • {{ song.album }}
              <span v-if="song.year && song.year !== 2025">({{ song.year }})</span>
            </div>
            <div class="text-xs opacity-60 mt-1">
              {{ song.genre || 'Unknown' }} | {{ formatDuration(song.duration) }} | 
              PLAYED: {{ song.playCount || 0 }}x
            </div>
          </div>
          <div class="flex items-center gap-3">
            <!-- Rating -->
            <div class="text-yellow-400 text-sm">
              {{ '★'.repeat(Math.floor((song.rating || 0) / 2)) }}{{ '☆'.repeat(5 - Math.floor((song.rating || 0) / 2)) }}
            </div>
            <!-- Play Button -->
            <button
              @click.stop="playSong(song)"
              class="opacity-0 group-hover:opacity-100 transition-opacity text-lime-400 hover:text-white"
              title="Play"
            >
              ▶
            </button>
          </div>
        </div>
      </div>

      <div
        v-if="libraryStore.searchQuery && searchResults.length === 0 && !libraryStore.loading"
        class="text-center py-8 text-lime-400"
      >
        <pre class="text-xs">
  ╔═══════════════════════════╗
  ║   SEARCH RETURNED 0       ║
  ║   TRY DIFFERENT QUERY     ║
  ╚═══════════════════════════╝
        </pre>
      </div>

      <div v-if="!libraryStore.searchQuery" class="text-center py-8 text-lime-400">
        <pre class="text-xs">
  ╔═══════════════════════════╗
  ║   ENTER SEARCH QUERY      ║
  ║   ABOVE TO BEGIN          ║
  ╚═══════════════════════════╝
        </pre>
      </div>

      <div v-if="libraryStore.loading" class="text-center py-8 text-lime-400 animate-pulse">
        <pre class="text-xs">
  ╔═══════════════════════════╗
  ║       SEARCHING...        ║
  ║   PLEASE WAIT             ║
  ╚═══════════════════════════╝
        </pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useLibraryStore } from '../../stores/library';
import { usePlayerStore } from '../../stores/player';
import type { Song } from '../../types';

const libraryStore = useLibraryStore();
const playerStore = usePlayerStore();

const searchResults = computed(() => libraryStore.filteredSongs);

// Debounced search
let searchTimeout: NodeJS.Timeout;
const handleSearch = () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    libraryStore.fetchSongs();
  }, 500);
};

onMounted(() => {
  libraryStore.fetchSongs();
});

const playSong = (song: Song) => {
  playerStore.playSong(song);
};

const formatDuration = (seconds: number) => {
  if (!seconds || seconds === 0) return '0:00';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
};
</script>