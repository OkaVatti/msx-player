<template>
  <div class="space-y-6">
    <!-- Library Header -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <h3 class="text-2xl font-bold text-lime-400">YOUR LIBRARY</h3>
        <p class="text-gray-400">{{ libraryStore.songs.length }} songs, sorted by {{ libraryStore.filters.sortBy }}</p>
      </div>
      <div class="flex gap-3">
        <select
          v-model="libraryStore.filters.sortBy"
          @change="handleFilterChange"
          class="bg-gray-700 border-2 border-lime-400 text-lime-400 px-3 py-2 rounded-lg focus:outline-none focus:bg-lime-400 focus:text-black transition-all"
        >
          <option value="title">Sort by Title</option>
          <option value="artist">Sort by Artist</option>
          <option value="album">Sort by Album</option>
          <option value="rating">Sort by Rating</option>
          <option value="playCount">Sort by Plays</option>
        </select>
        
        <select
          v-model="libraryStore.filters.genre"
          @change="handleFilterChange"
          class="bg-gray-700 border-2 border-lime-400 text-lime-400 px-3 py-2 rounded-lg focus:outline-none focus:bg-lime-400 focus:text-black transition-all"
        >
          <option value="">All Genres</option>
          <option v-for="genre in libraryStore.genres" :key="genre" :value="genre">
            {{ genre }}
          </option>
        </select>
        
        <!-- Play All Button -->
        <button
          @click="playAll"
          :disabled="libraryStore.songs.length === 0"
          class="bg-lime-400 text-black px-4 py-2 rounded-lg border-2 border-lime-400 hover:bg-black hover:text-lime-400 disabled:opacity-50 transition-all font-semibold"
        >
          ▶ PLAY ALL
        </button>
      </div>
    </div>

    <!-- Debug Info -->
    <div v-if="debug" class="bg-yellow-900 border-2 border-yellow-400 p-4 rounded-lg">
      <div class="text-yellow-400 text-sm font-mono">
        <div>🔍 DEBUG INFO</div>
        <div>Songs: {{ libraryStore.songs.length }}</div>
        <div>Current Song: {{ playerStore.currentSong?.title || 'None' }}</div>
        <div>Playing: {{ playerStore.isPlaying }}</div>
        <div>Audio Element: {{ playerStore.audioElement ? 'Loaded' : 'None' }}</div>
      </div>
    </div>

    <!-- Songs Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div
        v-for="(song, index) in libraryStore.filteredSongs"
        :key="song.id"
        class="bg-gray-800 rounded-xl border-2 border-gray-700 hover:border-lime-400 hover:glow-green transition-all duration-300 group cursor-pointer overflow-hidden"
        @click="playSong(song, index)"
      >
        <!-- Album Art Placeholder -->
        <div class="relative">
          <div class="w-full h-48 bg-linear-to-br from-purple-500 to-lime-400 flex items-center justify-center">
            <span class="text-4xl text-white">♪</span>
          </div>
          <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-20 transition-all flex items-center justify-center">
            <button class="opacity-0 group-hover:opacity-100 transform scale-75 group-hover:scale-100 transition-all duration-300 w-12 h-12 bg-lime-400 rounded-full flex items-center justify-center text-black text-xl">
              ▶
            </button>
          </div>
          <div class="absolute top-2 right-2">
            <span v-if="song.explicit" class="bg-red-500 text-white text-xs px-2 py-1 rounded">E</span>
            <span v-else-if="song.clean" class="bg-blue-500 text-white text-xs px-2 py-1 rounded">C</span>
          </div>
          
          <!-- Now Playing Indicator -->
          <div 
            v-if="playerStore.currentSong?.id === song.id"
            class="absolute bottom-2 left-2 bg-lime-400 text-black text-xs px-2 py-1 rounded font-semibold animate-pulse"
          >
            NOW PLAYING
          </div>
        </div>

        <!-- Song Info -->
        <div class="p-4">
          <h4 class="font-bold text-white truncate mb-1">{{ song.title || 'Unknown Title' }}</h4>
          <p class="text-lime-400 text-sm truncate mb-2">{{ song.artist || 'Unknown Artist' }}</p>
          <div class="flex justify-between items-center text-xs text-gray-400">
            <span>{{ song.album || 'Unknown Album' }}</span>
            <span>{{ formatDuration(song.duration) }}</span>
          </div>
          
          <!-- Rating & Plays -->
          <div class="flex justify-between items-center mt-3">
            <div class="flex items-center gap-1">
              <span 
                v-for="star in 5" 
                :key="star"
                class="text-sm"
                :class="star <= Math.floor((song.rating || 0) / 2) ? 'text-yellow-400' : 'text-gray-600'"
              >
                ★
              </span>
            </div>
            <div class="text-xs text-gray-400 flex items-center gap-1">
              <span>▶</span>
              <span>{{ song.play_count || 0 }}</span>
            </div>
          </div>
          
          <!-- File Info -->
          <div class="mt-2 text-xs text-gray-500 truncate">
            {{ getFilename(song.file_path) }}
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-if="libraryStore.filteredSongs.length === 0 && !libraryStore.loading"
      class="text-center py-16"
    >
      <div class="text-6xl mb-4">🎵</div>
      <h3 class="text-xl text-lime-400 mb-2">No songs found</h3>
      <p class="text-gray-400 mb-6">Try adjusting your filters or upload some music</p>
      <button 
        @click="libraryStore.fetchSongs()"
        class="px-6 py-3 border-2 border-lime-400 bg-lime-400 text-black hover:bg-black hover:text-lime-400 transition-all rounded-lg font-semibold"
      >
        Refresh Library
      </button>
    </div>

    <!-- Loading State -->
    <div
      v-if="libraryStore.loading"
      class="text-center py-16"
    >
      <div class="animate-pulse text-lime-400 text-lg">Loading your music library...</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { usePlayerStore } from '../../stores/player';
import { useLibraryStore } from '../../stores/library';
import type { Song } from '../../types';

const playerStore = usePlayerStore();
const libraryStore = useLibraryStore();
const debug = ref(true); // Set to false to hide debug info

onMounted(() => {
  console.log('📚 LibraryComponent mounted');
  libraryStore.fetchSongs();
});

const handleFilterChange = () => {
  libraryStore.fetchSongs();
};

const playSong = async (song: Song, index: number = 0) => {
  console.log('🎵 Playing song:', song);
  
  // Set the queue to current filtered songs starting from this index
  playerStore.setQueue(libraryStore.filteredSongs, index);
  
  // Play the song
  await playerStore.playSong(song);
};

const playAll = () => {
  if (libraryStore.filteredSongs.length > 0) {
    playerStore.setQueue(libraryStore.filteredSongs, 0);
    playerStore.playSong(libraryStore.filteredSongs[0]);
  }
};

const formatDuration = (seconds: number) => {
  if (!seconds || seconds === 0) return '0:00';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
};

const getFilename = (filePath: string) => {
  return filePath.split('/').pop() || filePath;
};
</script>

<style scoped>
.glow-green {
  box-shadow: 0 0 20px rgba(132, 204, 22, 0.3);
}
</style>