<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="border-2 border-lime-400 p-4 bg-black">
      <h2 class="text-xl mb-2 flex items-center gap-2">
        <span class="animate-pulse">&gt;&gt;</span>
        PLAYLIST MANAGEMENT SYSTEM
        <span class="text-xs text-white">[CREATE, EDIT, ORGANIZE]</span>
      </h2>
      <div class="text-sm text-lime-300 font-mono">
        TOTAL PLAYLISTS: <span class="text-white">{{ playlists.length }}</span>
      </div>
    </div>

    <!-- Create New Playlist -->
    <div class="border-2 border-lime-400 p-4 bg-black">
      <h3 class="text-lg mb-3 flex items-center gap-2">
        <span class="animate-pulse">+</span>
        CREATE NEW PLAYLIST
      </h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm text-lime-400 mb-2 font-mono">PLAYLIST NAME</label>
          <input
            v-model="newPlaylist.name"
            type="text"
            placeholder="ENTER PLAYLIST NAME..."
            class="w-full px-3 py-2 bg-black border-2 border-lime-400 text-lime-400 font-mono focus:outline-none focus:border-lime-300"
          />
        </div>
        <div>
          <label class="block text-sm text-lime-400 mb-2 font-mono">DESCRIPTION</label>
          <input
            v-model="newPlaylist.description"
            type="text"
            placeholder="OPTIONAL DESCRIPTION..."
            class="w-full px-3 py-2 bg-black border-2 border-lime-400 text-lime-400 font-mono focus:outline-none focus:border-lime-300"
          />
        </div>
      </div>
      
      <button
        @click="createPlaylist"
        :disabled="!newPlaylist.name"
        class="mt-4 px-6 py-2 border-2 border-lime-400 bg-lime-400 text-black hover:bg-black hover:text-lime-400 disabled:opacity-50 disabled:cursor-not-allowed transition-all font-mono"
      >
        [CREATE PLAYLIST]
      </button>
    </div>

    <!-- Playlists Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="playlist in playlists"
        :key="playlist.id"
        class="border-2 border-lime-400 p-4 bg-black hover:border-lime-300 transition-all group"
      >
        <!-- Playlist Header -->
        <div class="flex justify-between items-start mb-3">
          <h3 class="text-lg text-lime-300 font-mono truncate flex-1">{{ playlist.name }}</h3>
          <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              @click="editPlaylist(playlist)"
              class="text-lime-400 hover:text-yellow-400 transition-colors text-sm"
              title="Edit"
            >
              [E]
            </button>
            <button
              @click="deletePlaylist(playlist.id)"
              class="text-lime-400 hover:text-red-400 transition-colors text-sm"
              title="Delete"
            >
              [X]
            </button>
          </div>
        </div>

        <!-- Playlist Info -->
        <div class="text-sm text-lime-500 font-mono mb-3">
          {{ playlist.songs.length }} TRACKS
          <span v-if="playlist.description" class="block text-xs mt-1">
            {{ playlist.description }}
          </span>
        </div>

        <!-- Quick Actions -->
        <div class="flex gap-2 mb-3">
          <button
            @click="playPlaylist(playlist)"
            class="flex-1 px-3 py-1 border border-lime-400 text-lime-400 hover:bg-lime-400 hover:text-black transition-all text-sm font-mono"
          >
            PLAY
          </button>
          <button
            @click="manageSongs(playlist)"
            class="flex-1 px-3 py-1 border border-lime-400 text-lime-400 hover:bg-lime-400 hover:text-black transition-all text-sm font-mono"
          >
            MANAGE
          </button>
        </div>

        <!-- Recent Songs Preview -->
        <div class="space-y-1 max-h-32 overflow-y-auto">
          <div
            v-for="song in playlist.songs.slice(0, 3)"
            :key="song.id"
            class="text-xs text-lime-600 font-mono truncate border-b border-lime-800 pb-1 last:border-b-0"
          >
            {{ song.title }} - {{ song.artist }}
          </div>
          <div v-if="playlist.songs.length === 0" class="text-xs text-lime-700 font-mono italic">
            NO SONGS IN PLAYLIST
          </div>
        </div>
      </div>
    </div>

    <!-- Song Management Modal -->
    <div v-if="selectedPlaylist" class="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50 p-4">
      <div class="border-2 border-lime-400 bg-black w-full max-w-4xl max-h-[80vh] overflow-hidden">
        <!-- Modal Header -->
        <div class="border-b-2 border-lime-400 p-4 bg-black">
          <div class="flex justify-between items-center">
            <h3 class="text-lg text-lime-400 font-mono">
              MANAGE SONGS: {{ selectedPlaylist.name }}
            </h3>
            <button
              @click="selectedPlaylist = null"
              class="text-lime-400 hover:text-red-400 transition-colors"
            >
              [CLOSE]
            </button>
          </div>
        </div>

        <!-- Available Songs -->
        <div class="p-4 max-h-96 overflow-y-auto">
          <h4 class="text-sm text-lime-400 mb-3 font-mono">AVAILABLE SONGS</h4>
          <div class="space-y-2">
            <div
              v-for="song in availableSongs"
              :key="song.id"
              class="flex justify-between items-center p-2 border border-lime-800 hover:border-lime-400 transition-all"
            >
              <div class="flex-1 min-w-0">
                <div class="text-lime-300 text-sm font-mono truncate">{{ song.title }}</div>
                <div class="text-lime-500 text-xs font-mono">{{ song.artist }} • {{ song.album }}</div>
              </div>
              <button
                @click="addSongToPlaylist(song)"
                :disabled="isSongInPlaylist(song.id)"
                class="px-3 py-1 border border-lime-400 text-lime-400 hover:bg-lime-400 hover:text-black disabled:opacity-50 disabled:cursor-not-allowed transition-all text-xs font-mono"
              >
                {{ isSongInPlaylist(song.id) ? 'ADDED' : 'ADD' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Current Playlist Songs -->
        <div class="p-4 border-t-2 border-lime-400 max-h-96 overflow-y-auto">
          <h4 class="text-sm text-lime-400 mb-3 font-mono">
            PLAYLIST SONGS ({{ selectedPlaylist.songs.length }})
          </h4>
          <div class="space-y-2">
            <div
              v-for="(song, index) in selectedPlaylist.songs"
              :key="song.id"
              class="flex justify-between items-center p-2 border border-lime-800 bg-lime-900 bg-opacity-20"
            >
              <div class="flex items-center gap-3 flex-1 min-w-0">
                <span class="text-lime-500 text-xs font-mono w-4">{{ index + 1 }}</span>
                <div class="flex-1 min-w-0">
                  <div class="text-lime-300 text-sm font-mono truncate">{{ song.title }}</div>
                  <div class="text-lime-500 text-xs font-mono">{{ song.artist }}</div>
                </div>
              </div>
              <div class="flex gap-2">
                <button
                  @click="moveSong(index, -1)"
                  :disabled="index === 0"
                  class="text-lime-400 hover:text-yellow-400 disabled:opacity-30 transition-colors text-xs"
                >
                  ↑
                </button>
                <button
                  @click="moveSong(index, 1)"
                  :disabled="index === selectedPlaylist.songs.length - 1"
                  class="text-lime-400 hover:text-yellow-400 disabled:opacity-30 transition-colors text-xs"
                >
                  ↓
                </button>
                <button
                  @click="removeSongFromPlaylist(song.id)"
                  class="text-lime-400 hover:text-red-400 transition-colors text-xs"
                >
                  [X]
                </button>
              </div>
            </div>
            <div v-if="selectedPlaylist.songs.length === 0" class="text-center text-lime-600 text-sm font-mono py-4">
              NO SONGS IN THIS PLAYLIST
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

interface Song {
  id: number;
  title: string;
  artist: string;
  album: string;
  duration: number;
}

interface Playlist {
  id: number;
  name: string;
  description: string;
  songs: Song[];
  createdAt: string;
}

// Initialize with proper type safety
const playlists = ref<Playlist[]>([
  {
    id: 1,
    name: 'CHILL VIBES',
    description: 'Relaxing tunes for coding',
    songs: [],
    createdAt: '2024-01-15'
  },
  {
    id: 2,
    name: 'HYPE MODE',
    description: 'Energy boost tracks',
    songs: [],
    createdAt: '2024-01-20'
  },
  {
    id: 3,
    name: 'RETRO WAVES',
    description: '80s and 90s classics',
    songs: [],
    createdAt: '2024-01-25'
  }
]);

const newPlaylist = ref({
  name: '',
  description: ''
});

const selectedPlaylist = ref<Playlist | null>(null);

// Mock data with proper typing
const availableSongs = ref<Song[]>([
  { id: 1, title: 'Digital Dreams', artist: 'Synthwave', album: 'Retro Future', duration: 245 },
  { id: 2, title: 'Neon Nights', artist: 'Cyberpunk', album: 'City Lights', duration: 198 },
  { id: 3, title: 'Code Flow', artist: 'Dev Beats', album: 'Programming', duration: 312 },
  { id: 4, title: 'Binary Love', artist: 'Tech Romance', album: 'Algorithms', duration: 267 },
  { id: 5, title: 'Server Down', artist: 'Error 500', album: 'Debugging', duration: 189 }
]);

const createPlaylist = (): void => {
  if (!newPlaylist.value.name.trim()) return;

  const playlist: Playlist = {
    id: Date.now(),
    name: newPlaylist.value.name,
    description: newPlaylist.value.description || '', // Ensure string type
    songs: [],
    createdAt: new Date().toISOString().split('T')[0] || new Date().toISOString() // Fallback
  };

  playlists.value.push(playlist);
  newPlaylist.value = { name: '', description: '' };
};

const editPlaylist = (playlist: Playlist): void => {
  const newName = prompt('Enter new playlist name:', playlist.name);
  if (newName) {
    playlist.name = newName;
  }
};

const deletePlaylist = (id: number): void => {
  if (confirm('Are you sure you want to delete this playlist?')) {
    playlists.value = playlists.value.filter(p => p.id !== id);
  }
};

const playPlaylist = (playlist: Playlist): void => {
  console.log('Playing playlist:', playlist.name);
};

const manageSongs = (playlist: Playlist): void => {
  selectedPlaylist.value = { ...playlist };
};

const isSongInPlaylist = (songId: number): boolean => {
  return selectedPlaylist.value?.songs.some(s => s.id === songId) || false;
};

const addSongToPlaylist = (song: Song): void => {
  if (selectedPlaylist.value && !isSongInPlaylist(song.id)) {
    selectedPlaylist.value.songs.push({ ...song });
  }
};

const removeSongFromPlaylist = (songId: number): void => {
  if (selectedPlaylist.value) {
    selectedPlaylist.value.songs = selectedPlaylist.value.songs.filter(s => s.id !== songId);
  }
};

const moveSong = (index: number, direction: number): void => {
  if (!selectedPlaylist.value) return;

  const newIndex = index + direction;
  if (newIndex >= 0 && newIndex < selectedPlaylist.value.songs.length) {
    const songs = selectedPlaylist.value.songs;
    [songs[index], songs[newIndex]] = [songs[newIndex], songs[index]];
  }
};
</script>

<style scoped>
/* Custom scrollbars */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #000;
  border: 1px solid #84cc16;
}

::-webkit-scrollbar-thumb {
  background: #84cc16;
}

::-webkit-scrollbar-thumb:hover {
  background: #a3e635;
}
</style>