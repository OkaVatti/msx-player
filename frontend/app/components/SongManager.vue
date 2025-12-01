<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="border-2 border-lime-400 p-4 bg-black">
      <h2 class="text-xl mb-2 flex items-center gap-2">
        <span class="animate-pulse">&gt;&gt;</span>
        SONG MANAGEMENT SYSTEM
        <span class="text-xs text-white">[EDIT, ORGANIZE, DELETE]</span>
      </h2>
      <div class="text-sm text-lime-300 font-mono">
        LIBRARY STATS: 
        <span class="text-white">{{ songs.length }} SONGS</span> • 
        <span class="text-yellow-400">{{ totalDuration }}</span> • 
        <span class="text-blue-400">{{ uniqueArtists }} ARTISTS</span>
        <span v-if="loading" class="text-yellow-400 animate-pulse ml-2">[LOADING...]</span>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="border-2 border-lime-400 p-4 bg-black">
      <div class="flex flex-wrap gap-2">
        <button
          @click="exportLibrary"
          class="px-4 py-2 border-2 border-lime-400 bg-black text-lime-400 hover:bg-lime-400 hover:text-black transition-all font-mono text-sm"
        >
          [EXPORT LIBRARY]
        </button>
        <button
          @click="scanForDuplicates"
          class="px-4 py-2 border-2 border-yellow-400 bg-black text-yellow-400 hover:bg-yellow-400 hover:text-black transition-all font-mono text-sm"
        >
          [SCAN DUPLICATES]
        </button>
        <button
          @click="refreshLibrary"
          :disabled="loading"
          class="px-4 py-2 border-2 border-blue-400 bg-black text-blue-400 hover:bg-blue-400 hover:text-black transition-all font-mono text-sm disabled:opacity-50 disabled:cursor-not-allowed"
        >
          [{{ loading ? 'REFRESHING...' : 'REFRESH LIBRARY' }}]
        </button>
        <button
          @click="showMetadataEditor = true"
          class="px-4 py-2 border-2 border-purple-400 bg-black text-purple-400 hover:bg-purple-400 hover:text-black transition-all font-mono text-sm"
        >
          [BATCH EDIT]
        </button>
      </div>
    </div>

    <!-- Status Messages -->
    <div v-if="errorMessage" class="border-2 border-red-400 p-4 bg-black">
      <div class="text-red-400 font-mono text-sm">
        ERROR: {{ errorMessage }}
      </div>
    </div>

    <div v-if="successMessage" class="border-2 border-green-400 p-4 bg-black">
      <div class="text-green-400 font-mono text-sm">
        SUCCESS: {{ successMessage }}
      </div>
    </div>

    <!-- Songs Table -->
    <div class="border-2 border-lime-400 bg-black overflow-x-auto">
      <table class="w-full font-mono text-sm">
        <thead class="border-b-2 border-lime-400">
          <tr>
            <th class="p-3 text-left text-lime-400 font-mono">
              <input type="checkbox" v-model="selectAll" />
            </th>
            <th class="p-3 text-left text-lime-400 font-mono">TITLE</th>
            <th class="p-3 text-left text-lime-400 font-mono">ARTIST</th>
            <th class="p-3 text-left text-lime-400 font-mono">ALBUM</th>
            <th class="p-3 text-left text-lime-400 font-mono">DURATION</th>
            <th class="p-3 text-left text-lime-400 font-mono">RATING</th>
            <th class="p-3 text-left text-lime-400 font-mono">PLAYS</th>
            <th class="p-3 text-left text-lime-400 font-mono">ACTIONS</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="song in songs"
            :key="song.id"
            class="border-b border-lime-800 hover:bg-lime-900 hover:bg-opacity-20 transition-all"
          >
            <td class="p-3">
              <input type="checkbox" v-model="selectedSongs" :value="song.id" />
            </td>
            <td class="p-3">
              <div class="text-lime-300">{{ song.title || 'Unknown Title' }}</div>
              <div class="text-lime-500 text-xs">
                <span v-if="song.explicit" class="text-red-400">[EXPLICIT]</span>
                <span v-else-if="song.clean" class="text-blue-400">[CLEAN]</span>
              </div>
            </td>
            <td class="p-3 text-lime-300">{{ song.artist || 'Unknown Artist' }}</td>
            <td class="p-3 text-lime-300">{{ song.album || 'Unknown Album' }}</td>
            <td class="p-3 text-lime-300 font-mono">{{ formatDuration(song.duration) }}</td>
            <td class="p-3">
              <div class="flex items-center gap-1">
                <button
                  v-for="star in 10"
                  :key="star"
                  @click="rateSong(song.id, star)"
                  class="text-xs transition-all hover:scale-125"
                  :class="star <= (song.rating || 0) ? 'text-yellow-400' : 'text-lime-600'"
                >
                  ★
                </button>
              </div>
            </td>
            <td class="p-3 text-lime-300 font-mono">{{ song.playCount || 0 }}</td>
            <td class="p-3">
              <div class="flex gap-2">
                <button
                  @click="editSong(song)"
                  class="text-lime-400 hover:text-yellow-400 transition-colors text-xs"
                  title="Edit"
                >
                  [EDIT]
                </button>
                <button
                  @click="playSong(song)"
                  class="text-lime-400 hover:text-green-400 transition-colors text-xs"
                  title="Play"
                >
                  [PLAY]
                </button>
                <button
                  @click="deleteSong(song.id)"
                  class="text-lime-400 hover:text-red-400 transition-colors text-xs"
                  title="Delete"
                >
                  [DEL]
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Loading State -->
      <div v-if="loading && songs.length === 0" class="text-center py-8">
        <div class="text-lime-400 font-mono animate-pulse">
          LOADING SONGS FROM DATABASE...
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="songs.length === 0" class="text-center py-8">
        <pre class="text-xs text-lime-400">
╔═══════════════════════════╗
║   NO SONGS IN LIBRARY     ║
║   UPLOAD SOME MUSIC       ║
╚═══════════════════════════╝
        </pre>
      </div>
    </div>

    <!-- Bulk Actions -->
    <div v-if="selectedSongs.length > 0" class="border-2 border-yellow-400 p-4 bg-black">
      <h3 class="text-lg mb-3 flex items-center gap-2 text-yellow-400">
        <span class="animate-pulse">!</span>
        BULK ACTIONS ({{ selectedSongs.length }} SONGS SELECTED)
      </h3>
      <div class="flex flex-wrap gap-2">
        <button
          @click="addToPlaylistBulk"
          class="px-4 py-2 border-2 border-lime-400 bg-black text-lime-400 hover:bg-lime-400 hover:text-black transition-all font-mono text-sm"
        >
          [ADD TO PLAYLIST]
        </button>
        <button
          @click="rateBulk"
          class="px-4 py-2 border-2 border-yellow-400 bg-black text-yellow-400 hover:bg-yellow-400 hover:text-black transition-all font-mono text-sm"
        >
          [SET RATING]
        </button>
        <button
          @click="deleteBulk"
          class="px-4 py-2 border-2 border-red-400 bg-black text-red-400 hover:bg-red-400 hover:text-black transition-all font-mono text-sm"
        >
          [DELETE SELECTED]
        </button>
        <button
          @click="clearSelection"
          class="px-4 py-2 border-2 border-gray-400 bg-black text-gray-400 hover:bg-gray-400 hover:text-black transition-all font-mono text-sm"
        >
          [CLEAR SELECTION]
        </button>
      </div>
    </div>

    <!-- Edit Song Modal -->
    <div v-if="editingSong" class="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50 p-4">
      <div class="border-2 border-lime-400 bg-black w-full max-w-2xl">
        <div class="border-b-2 border-lime-400 p-4">
          <div class="flex justify-between items-center">
            <h3 class="text-lg text-lime-400 font-mono">EDIT SONG METADATA</h3>
            <button
              @click="editingSong = null"
              class="text-lime-400 hover:text-red-400 transition-colors"
            >
              [CLOSE]
            </button>
          </div>
        </div>

        <div class="p-4 space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm text-lime-400 mb-2 font-mono">TITLE</label>
              <input
                v-model="editingSong.title"
                type="text"
                class="w-full px-3 py-2 bg-black border-2 border-lime-400 text-lime-400 font-mono focus:outline-none focus:border-lime-300"
              />
            </div>
            <div>
              <label class="block text-sm text-lime-400 mb-2 font-mono">ARTIST</label>
              <input
                v-model="editingSong.artist"
                type="text"
                class="w-full px-3 py-2 bg-black border-2 border-lime-400 text-lime-400 font-mono focus:outline-none focus:border-lime-300"
              />
            </div>
            <div>
              <label class="block text-sm text-lime-400 mb-2 font-mono">ALBUM</label>
              <input
                v-model="editingSong.album"
                type="text"
                class="w-full px-3 py-2 bg-black border-2 border-lime-400 text-lime-400 font-mono focus:outline-none focus:border-lime-300"
              />
            </div>
            <div>
              <label class="block text-sm text-lime-400 mb-2 font-mono">GENRE</label>
              <input
                v-model="editingSong.genre"
                type="text"
                class="w-full px-3 py-2 bg-black border-2 border-lime-400 text-lime-400 font-mono focus:outline-none focus:border-lime-300"
              />
            </div>
            <div>
              <label class="block text-sm text-lime-400 mb-2 font-mono">YEAR</label>
              <input
                v-model="editingSong.year"
                type="number"
                min="1900"
                max="2030"
                class="w-full px-3 py-2 bg-black border-2 border-lime-400 text-lime-400 font-mono focus:outline-none focus:border-lime-300"
              />
            </div>
            <div>
              <label class="block text-sm text-lime-400 mb-2 font-mono">TRACK NUMBER</label>
              <input
                v-model="editingSong.trackNumber"
                type="number"
                min="1"
                class="w-full px-3 py-2 bg-black border-2 border-lime-400 text-lime-400 font-mono focus:outline-none focus:border-lime-300"
              />
            </div>
          </div>

          <div class="flex gap-4">
            <div class="flex-1">
              <label class="block text-sm text-lime-400 mb-2 font-mono">CONTENT TYPE</label>
              <div class="flex gap-4">
                <label class="flex items-center gap-2 text-lime-400">
                  <input type="radio" v-model="contentType" value="explicit" />
                  EXPLICIT
                </label>
                <label class="flex items-center gap-2 text-lime-400">
                  <input type="radio" v-model="contentType" value="clean" />
                  CLEAN
                </label>
                <label class="flex items-center gap-2 text-lime-400">
                  <input type="radio" v-model="contentType" value="neutral" />
                  NEUTRAL
                </label>
              </div>
            </div>
          </div>

          <div class="flex gap-2 pt-4 border-t border-lime-400">
            <button
              @click="saveSongEdit"
              :disabled="saving"
              class="px-6 py-2 border-2 border-lime-400 bg-lime-400 text-black hover:bg-black hover:text-lime-400 transition-all font-mono disabled:opacity-50 disabled:cursor-not-allowed"
            >
              [{{ saving ? 'SAVING...' : 'SAVE CHANGES' }}]
            </button>
            <button
              @click="editingSong = null"
              class="px-6 py-2 border-2 border-lime-400 bg-black text-lime-400 hover:bg-lime-400 hover:text-black transition-all font-mono"
            >
              [CANCEL]
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useLibraryStore } from '../../stores/library';
import { usePlayerStore } from '../../stores/player';
import type { Song } from '../../types';

const libraryStore = useLibraryStore();
const playerStore = usePlayerStore();

// Reactive data
const selectedSongs = ref<number[]>([]);
const editingSong = ref<Song | null>(null);
const showMetadataEditor = ref(false);
const loading = ref(false);
const saving = ref(false);
const errorMessage = ref('');
const successMessage = ref('');

// Computed properties
const songs = computed(() => libraryStore.songs);

const selectAll = computed({
  get: () => selectedSongs.value.length === songs.value.length && songs.value.length > 0,
  set: (value: boolean) => {
    selectedSongs.value = value ? songs.value.map(s => s.id) : [];
  }
});

const totalDuration = computed(() => {
  const totalSeconds = songs.value.reduce((sum, song) => sum + (song.duration || 0), 0);
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  return `${hours}h ${minutes}m`;
});

const uniqueArtists = computed(() => {
  const artists = songs.value.map(s => s.artist || 'Unknown Artist').filter(artist => artist !== 'Unknown Artist');
  return new Set(artists).size;
});

const contentType = computed({
  get: () => {
    if (!editingSong.value) return 'neutral';
    if (editingSong.value.explicit) return 'explicit';
    if (editingSong.value.clean) return 'clean';
    return 'neutral';
  },
  set: (value: string) => {
    if (!editingSong.value) return;
    editingSong.value.explicit = value === 'explicit';
    editingSong.value.clean = value === 'clean';
  }
});

// Methods
const formatDuration = (seconds: number) => {
  if (!seconds || seconds === 0) return '0:00';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
};

const showMessage = (message: string, type: 'success' | 'error' = 'success') => {
  if (type === 'success') {
    successMessage.value = message;
    errorMessage.value = '';
  } else {
    errorMessage.value = message;
    successMessage.value = '';
  }
  
  setTimeout(() => {
    successMessage.value = '';
    errorMessage.value = '';
  }, 5000);
};

const loadSongs = async () => {
  try {
    loading.value = true;
    errorMessage.value = '';
    await libraryStore.fetchSongs();
    console.log('✅ Songs loaded:', songs.value.length);
  } catch (error) {
    console.error('❌ Failed to load songs:', error);
    showMessage('Failed to load songs from database', 'error');
  } finally {
    loading.value = false;
  }
};

const editSong = (song: Song) => {
  editingSong.value = { ...song };
};

const saveSongEdit = async () => {
  if (!editingSong.value) return;

  try {
    saving.value = true;
    errorMessage.value = '';

    const updates = {
      title: editingSong.value.title,
      artist: editingSong.value.artist,
      album: editingSong.value.album,
      genre: editingSong.value.genre,
      year: editingSong.value.year,
      trackNumber: editingSong.value.trackNumber,
      explicit: editingSong.value.explicit,
      clean: editingSong.value.clean,
    };

    await libraryStore.updateSong(editingSong.value.id, updates);
    showMessage('Song updated successfully');
    editingSong.value = null;
  } catch (error) {
    console.error('❌ Failed to update song:', error);
    showMessage('Failed to update song', 'error');
  } finally {
    saving.value = false;
  }
};

const playSong = (song: Song) => {
  console.log('🎵 Playing song:', song.title);
  playerStore.playSong(song);
};

const deleteSong = async (id: number) => {
  if (!confirm('Are you sure you want to delete this song from your library?')) {
    return;
  }

  try {
    errorMessage.value = '';
    await libraryStore.deleteSong(id);
    showMessage('Song deleted successfully');
    // Remove from selected songs if it was selected
    selectedSongs.value = selectedSongs.value.filter(songId => songId !== id);
  } catch (error) {
    console.error('❌ Failed to delete song:', error);
    showMessage('Failed to delete song', 'error');
  }
};

const rateSong = async (songId: number, rating: number) => {
  try {
    errorMessage.value = '';
    await playerStore.rateSong(songId, rating);
    // Refresh the songs to show updated ratings
    await loadSongs();
  } catch (error) {
    console.error('❌ Failed to rate song:', error);
    showMessage('Failed to rate song', 'error');
  }
};

const exportLibrary = () => {
  try {
    const dataStr = JSON.stringify(songs.value, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    
    const link = document.createElement('a');
    link.href = URL.createObjectURL(dataBlob);
    link.download = `msx-library-export-${new Date().toISOString().split('T')[0]}.json`;
    link.click();
    
    showMessage('Library exported successfully');
  } catch (error) {
    console.error('❌ Failed to export library:', error);
    showMessage('Failed to export library', 'error');
  }
};

const scanForDuplicates = () => {
  // Simple duplicate detection by title and artist
  const seen = new Map();
  const duplicates = [];

  songs.value.forEach(song => {
    const key = `${song.title}-${song.artist}`.toLowerCase();
    if (seen.has(key)) {
      duplicates.push({ song, duplicateOf: seen.get(key) });
    } else {
      seen.set(key, song);
    }
  });

  if (duplicates.length > 0) {
    alert(`Found ${duplicates.length} potential duplicates. Check console for details.`);
    console.log('🔍 Duplicates found:', duplicates);
    showMessage(`Found ${duplicates.length} potential duplicates`, 'error');
  } else {
    showMessage('No duplicates found');
  }
};

const refreshLibrary = async () => {
  await loadSongs();
};

const addToPlaylistBulk = () => {
  if (selectedSongs.value.length === 0) return;
  
  // For now, just show a message - you can implement playlist functionality later
  showMessage(`Added ${selectedSongs.value.length} songs to playlist queue`);
  console.log('Adding selected songs to playlist:', selectedSongs.value);
};

const rateBulk = async () => {
  if (selectedSongs.value.length === 0) return;

  const ratingStr = prompt(`Enter rating (1-10) for ${selectedSongs.value.length} songs:`);
  if (!ratingStr) return;

  const rating = parseInt(ratingStr);
  if (isNaN(rating) || rating < 1 || rating > 10) {
    showMessage('Rating must be a number between 1 and 10', 'error');
    return;
  }

  try {
    errorMessage.value = '';
    
    // Rate each selected song
    for (const songId of selectedSongs.value) {
      await playerStore.rateSong(songId, rating);
    }
    
    // Refresh to show updated ratings
    await loadSongs();
    showMessage(`Rated ${selectedSongs.value.length} songs as ${rating}/10`);
  } catch (error) {
    console.error('❌ Failed to bulk rate songs:', error);
    showMessage('Failed to rate songs', 'error');
  }
};

const deleteBulk = async () => {
  if (selectedSongs.value.length === 0) return;

  if (!confirm(`Are you sure you want to delete ${selectedSongs.value.length} songs? This cannot be undone.`)) {
    return;
  }

  try {
    errorMessage.value = '';
    
    // Delete each selected song
    for (const songId of selectedSongs.value) {
      await libraryStore.deleteSong(songId);
    }
    
    showMessage(`Deleted ${selectedSongs.value.length} songs`);
    selectedSongs.value = [];
  } catch (error) {
    console.error('❌ Failed to bulk delete songs:', error);
    showMessage('Failed to delete songs', 'error');
  }
};

const clearSelection = () => {
  selectedSongs.value = [];
};

// Lifecycle
onMounted(() => {
  console.log('🔄 SongManager mounted - loading songs...');
  loadSongs();
});

// Watch for store changes
watch(() => libraryStore.songs, (newSongs) => {
  console.log('📊 Songs updated in SongManager:', newSongs.length);
});
</script>

<style scoped>
/* Custom scrollbar for table */
::-webkit-scrollbar {
  height: 8px;
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

/* Table styling */
table {
  border-collapse: collapse;
}

th, td {
  border-right: 1px solid #84cc16;
}

th:last-child, td:last-child {
  border-right: none;
}

/* Loading animation */
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>