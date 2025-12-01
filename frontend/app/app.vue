<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-900 to-black text-white font-mono">
    <!-- Top Navigation Bar -->
    <div class="bg-black border-b-2 border-lime-400 p-4">
      <div class="flex justify-between items-center">
        <!-- Logo -->
        <div class="flex items-center gap-4">
          <div class="text-2xl text-lime-400 animate-pulse">🎵</div>
          <h1 class="text-xl font-bold text-lime-400">
            MSX<span class="text-white">PLAYER</span>
          </h1>
          <div class="text-xs text-gray-400">v1.3.37</div>
        </div>

        <!-- Top Controls -->
        <div class="flex items-center gap-3">
          <button class="p-2 hover:bg-lime-400 hover:text-black transition-all rounded">
            <span class="text-xs">[MINIMIZE]</span>
          </button>
          <button class="p-2 hover:bg-lime-400 hover:text-black transition-all rounded">
            <span class="text-xs">[MAXIMIZE]</span>
          </button>
          <button class="p-2 hover:bg-red-500 hover:text-white transition-all rounded">
            <span class="text-xs">[CLOSE]</span>
          </button>
        </div>
      </div>
    </div>

    <div class="flex h-[calc(100vh-80px)]">
      <!-- Sidebar -->
      <div class="w-64 bg-gray-900 border-r-2 border-lime-400 p-4">
        <!-- Navigation -->
        <nav class="space-y-2 mb-8">
          <div class="text-xs text-gray-400 uppercase tracking-wider mb-2">NAVIGATION</div>
          <button
            v-for="page in pages"
            :key="page.id"
            @click="currentPage = page.id"
            :class="[
              'w-full text-left p-3 rounded-lg border-2 transition-all duration-200 flex items-center gap-3',
              currentPage === page.id
                ? 'bg-lime-400 text-black border-lime-400 shadow-lg shadow-lime-400/25'
                : 'bg-gray-800 text-lime-400 border-gray-600 hover:bg-lime-400 hover:text-black'
            ]"
          >
            <span class="text-lg">{{ page.icon }}</span>
            <span class="font-semibold">{{ page.name }}</span>
          </button>
        </nav>

        <!-- Now Playing -->
        <div class="bg-gray-800 rounded-lg p-4 border-2 border-lime-400 mb-4">
          <div class="text-xs text-gray-400 mb-2">NOW PLAYING</div>
          <div v-if="playerStore.currentSong" class="space-y-2">
            <div class="text-sm font-bold text-white truncate">{{ playerStore.currentSong.title }}</div>
            <div class="text-xs text-lime-400 truncate">{{ playerStore.currentSong.artist }}</div>
            <div class="w-full bg-gray-700 rounded-full h-1">
              <div 
                class="bg-lime-400 h-1 rounded-full transition-all duration-500"
                :style="{ width: playerStore.progressPercentage + '%' }"
              ></div>
            </div>
          </div>
          <div v-else class="text-xs text-gray-500 italic">
            NO TRACK SELECTED
          </div>
        </div>

        <!-- Quick Stats -->
        <div class="bg-gray-800 rounded-lg p-4 border-2 border-purple-400">
          <div class="text-xs text-gray-400 mb-2">LIBRARY STATS</div>
          <div class="space-y-1 text-xs">
            <div class="flex justify-between">
              <span class="text-gray-400">Tracks:</span>
              <span class="text-lime-400">{{ libraryStore.songs.length }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-400">Playlists:</span>
              <span class="text-purple-400">12</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-400">Storage:</span>
              <span class="text-blue-400">2.3 GB</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content -->
      <div class="flex-1 flex flex-col">
        <!-- Content Header -->
        <div class="bg-gradient-to-r from-gray-800 to-gray-900 border-b-2 border-lime-400 p-6">
          <div class="flex justify-between items-center">
            <div>
              <h2 class="text-2xl font-bold text-lime-400 mb-2">
                {{ currentPageName }}
              </h2>
              <div class="text-sm text-gray-400">
                {{ currentPageDescription }}
              </div>
            </div>
            <div class="flex items-center gap-4">
              <!-- Search Bar -->
              <div class="relative">
                <input
                  type="text"
                  placeholder="SEARCH LIBRARY..."
                  class="bg-black border-2 border-lime-400 text-lime-400 px-4 py-2 rounded-lg w-64 focus:outline-none focus:bg-lime-400 focus:text-black transition-all font-mono"
                />
                <span class="absolute right-3 top-2 text-gray-400">⌕</span>
              </div>
              <!-- Visualizer Toggle -->
              <button
                @click="toggleMiniVisualizer"
                class="p-2 border-2 border-lime-400 bg-black text-lime-400 hover:bg-lime-400 hover:text-black transition-all rounded-lg"
              >
                {{ showMiniVisualizer ? '[- VIS]' : '[+ VIS]' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Page Content -->
        <div class="flex-1 p-6 overflow-auto bg-gradient-to-b from-gray-800 to-gray-900">
          <Transition name="page-fade" mode="out-in">
            <component :is="currentComponent" />
          </Transition>
        </div>
      </div>
    </div>

    <!-- Now Playing Bar -->
    <div class="fixed bottom-0 left-0 right-0 bg-black border-t-2 border-lime-400 p-4">
      <div class="flex items-center justify-between">
        <!-- Song Info -->
        <div class="flex items-center gap-4 w-64">
          <div v-if="playerStore.currentSong" class="flex items-center gap-3">
            <div class="w-12 h-12 bg-gradient-to-br from-lime-400 to-green-500 rounded-lg flex items-center justify-center">
              <span class="text-white text-lg">♪</span>
            </div>
            <div>
              <div class="text-sm font-semibold text-white">{{ playerStore.currentSong.title }}</div>
              <div class="text-xs text-lime-400">{{ playerStore.currentSong.artist }}</div>
            </div>
          </div>
          <div v-else class="text-xs text-gray-500">
            NO SONG PLAYING
          </div>
        </div>

        <!-- Player Controls -->
        <div class="flex-1 max-w-2xl">
          <div class="flex flex-col items-center gap-3">
            <!-- Control Buttons -->
            <div class="flex items-center gap-6">
              <button @click="playerStore.previousSong" class="text-2xl text-lime-400 hover:text-white transition-all">
                ⏮
              </button>
              <button 
                v-if="!playerStore.isPlaying"
                @click="resumePlayer" 
                class="w-12 h-12 bg-lime-400 text-black rounded-full flex items-center justify-center hover:bg-lime-300 transition-all"
              >
                ▶
              </button>
              <button 
                v-else
                @click="playerStore.pause" 
                class="w-12 h-12 bg-yellow-400 text-black rounded-full flex items-center justify-center hover:bg-yellow-300 transition-all"
              >
                ⏸
              </button>
              <button @click="playerStore.nextSong" class="text-2xl text-lime-400 hover:text-white transition-all">
                ⏭
              </button>
            </div>
            
            <!-- Progress Bar -->
            <div class="w-full flex items-center gap-3">
              <span class="text-xs text-gray-400">{{ formatTime(playerStore.currentTime) }}</span>
              <div 
                class="flex-1 bg-gray-700 rounded-full h-1 cursor-pointer"
                @click="seekToTime"
                ref="progressBar"
              >
                <div 
                  class="bg-gradient-to-r from-lime-400 to-green-500 h-1 rounded-full transition-all duration-500"
                  :style="{ width: playerStore.progressPercentage + '%' }"
                ></div>
              </div>
              <span class="text-xs text-gray-400">{{ formatTime(playerStore.duration) }}</span>
            </div>
          </div>
        </div>

        <!-- Volume & Extra Controls -->
        <div class="w-64 flex items-center justify-end gap-4">
          <button class="text-lime-400 hover:text-white transition-all">
            ♡
          </button>
          <div class="flex items-center gap-2">
            <span class="text-xs text-gray-400">🔊</span>
            <input 
              type="range" 
              v-model="playerStore.volume"
              @change="playerStore.setVolume(playerStore.volume)"
              class="w-20 accent-lime-400"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Mini Visualizer -->
    <div
      v-if="showMiniVisualizer"
      class="fixed bottom-24 right-6 w-80 h-48 border-2 border-lime-400 bg-black rounded-lg overflow-hidden z-50"
    >
      <AudioVisualizerMini @close="showMiniVisualizer = false" />
    </div>

    <!-- Retro Scanlines Overlay -->
    <div class="fixed inset-0 pointer-events-none scanlines opacity-20 z-40"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { storeToRefs } from 'pinia';
import { usePlayerStore } from '../stores/player';
import { useLibraryStore } from '../stores/library';

// Components
import SearchComponent from './components/SearchComponent.vue';
import LibraryComponent from './components/LibraryComponent.vue';
import AudioVisualizer from './components/AudioVisualizer.vue';
import ControlCenter from './components/ControlCenter.vue';
import UploadComponent from './components/UploadComponent.vue';
import PlaylistManager from './components/PlaylistManager.vue';
import SongManager from './components/SongManager.vue';
import AudioVisualizerMini from './components/AudioVisualizerMini.vue';

const playerStore = usePlayerStore();
const libraryStore = useLibraryStore();
const { currentSong, currentTime, duration, isPlaying, volume } = storeToRefs(playerStore);

const currentPage = ref('library');
const showMiniVisualizer = ref(false);
const progressBar = ref<HTMLDivElement>();

const pages = [
  { id: 'library', name: 'Library', icon: '📚', description: 'Browse your music collection' },
  { id: 'search', name: 'Search', icon: '🔍', description: 'Find songs in your library' },
  { id: 'playlists', name: 'Playlists', icon: '🎧', description: 'Manage your playlists' },
  { id: 'upload', name: 'Upload', icon: '📁', description: 'Add new music to your library' },
  { id: 'manage', name: 'Manage', icon: '⚙️', description: 'Song management and organization' },
  { id: 'visualizer', name: 'Visualizer', icon: '🌌', description: 'Audio visualization' },
  { id: 'controls', name: 'Controls', icon: '🎛️', description: 'Player controls and settings' },
];

const currentPageName = computed(() => {
  return pages.find(p => p.id === currentPage.value)?.name || 'Library';
});

const currentPageDescription = computed(() => {
  return pages.find(p => p.id === currentPage.value)?.description || 'Browse your music collection';
});

const currentComponent = computed(() => {
  const components = {
    search: SearchComponent,
    library: LibraryComponent,
    visualizer: AudioVisualizer,
    controls: ControlCenter,
    upload: UploadComponent,
    playlists: PlaylistManager,
    manage: SongManager,
  };
  return components[currentPage.value as keyof typeof components] || LibraryComponent;
});

const toggleMiniVisualizer = () => {
  showMiniVisualizer.value = !showMiniVisualizer.value;
};

const formatTime = (seconds: number) => {
  if (!seconds || isNaN(seconds)) return '0:00';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
};

const seekToTime = (event: MouseEvent) => {
  if (!progressBar.value || !duration.value) return;
  
  const rect = progressBar.value.getBoundingClientRect();
  const percent = (event.clientX - rect.left) / rect.width;
  const newTime = percent * duration.value;
  
  playerStore.seek(newTime);
};

const resumePlayer = () => {
  if (currentSong.value) {
    playerStore.playSong(currentSong.value);
  }
};

onMounted(() => {
  console.log('🚀 MSX Player initialized');
  libraryStore.fetchSongs();
});
</script>

<style>
@keyframes page-fade {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.page-fade-enter-active,
.page-fade-leave-active {
  transition: all 0.3s ease;
}

.page-fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.page-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* Retro scanlines effect */
.scanlines {
  background: linear-gradient(
    to bottom,
    transparent 50%,
    rgba(0, 0, 0, 0.1) 51%
  );
  background-size: 100% 4px;
  animation: scanline 6s linear infinite;
}

@keyframes scanline {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 0 100%;
  }
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #1f2937;
}

::-webkit-scrollbar-thumb {
  background: #84cc16;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a3e635;
}

/* Glow effects */
.glow-green {
  box-shadow: 0 0 20px rgba(132, 204, 22, 0.3);
}

/* Retro CRT curvature effect */
.crt-curve {
  transform: perspective(1000px) rotateX(2deg);
}
</style>