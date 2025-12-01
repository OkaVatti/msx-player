<template>
  <div class="h-full flex flex-col p-4 space-y-4">
    <!-- Header -->
    <div class="border-2 border-[#837dbd] p-4 bg-black">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-xl text-[#d3ceff] font-bold tracking-wide flex items-center gap-2">
            <span class="animate-pulse">&gt;&gt;</span>
            PLAYLISTS
          </h1>
          <p class="text-sm text-[#837dbd] font-mono mt-1">
            {{ playlists.length }} playlists • {{ totalSongs }} tracks
          </p>
        </div>
        <button
          @click="showCreateModal = true"
          class="px-4 py-2 bg-[#837dbd] text-black hover:bg-[#d3ceff] transition-all font-mono text-sm"
        >
          [NEW PLAYLIST]
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="text-[#837dbd] text-xl font-mono animate-pulse">LOADING PLAYLISTS...</div>
        <div class="text-[#837dbd] text-sm mt-2">please wait</div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="text-red-400 text-xl font-mono">ERROR LOADING PLAYLISTS</div>
        <div class="text-[#837dbd] text-sm mt-2">{{ error }}</div>
        <button
          @click="fetchPlaylists"
          class="mt-4 px-4 py-2 border-2 border-red-400 text-red-400 hover:bg-red-400 hover:text-black transition-all font-mono"
        >
          [RETRY]
        </button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="playlists.length === 0" class="flex-1 flex items-center justify-center">
      <div class="text-center">
        <div class="text-[#837dbd] text-xl font-mono">NO PLAYLISTS FOUND</div>
        <div class="text-[#837dbd] text-sm mt-2">create your first playlist to get started</div>
        <button
          @click="showCreateModal = true"
          class="mt-4 px-4 py-2 bg-[#837dbd] text-black hover:bg-[#d3ceff] transition-all font-mono"
        >
          [CREATE PLAYLIST]
        </button>
      </div>
    </div>

    <!-- Playlists Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 flex-1 overflow-y-auto">
      <div
        v-for="playlist in playlists"
        :key="playlist.id"
        class="border-2 border-[#837dbd] p-4 bg-black hover:border-[#d3ceff] transition-all group"
        @dblclick="openPlaylist(playlist)"
      >
        <!-- Playlist Header -->
        <div class="flex justify-between items-start mb-3">
          <div class="flex-1 min-w-0">
            <h3 class="text-lg text-[#d3ceff] font-mono truncate">{{ playlist.name }}</h3>
            <div class="text-sm text-[#837dbd] font-mono mt-1">
              {{ playlist.songs?.length || 0 }} tracks
            </div>
          </div>
          <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              @click.stop="editPlaylist(playlist)"
              class="text-[#837dbd] hover:text-[#d3ceff] transition-colors text-sm"
              title="Edit"
            >
              [E]
            </button>
            <button
              @click.stop="deletePlaylist(playlist.id)"
              class="text-[#837dbd] hover:text-red-400 transition-colors text-sm"
              title="Delete"
            >
              [X]
            </button>
          </div>
        </div>

        <!-- Description -->
        <div v-if="playlist.description" class="text-sm text-[#837dbd] font-mono mb-4 line-clamp-2">
          {{ playlist.description }}
        </div>

        <!-- Created Date -->
        <div class="text-xs text-[#5a548d] font-mono mb-3">
          created {{ formatDate(playlist.created_at) }}
        </div>

        <!-- Quick Actions -->
        <div class="flex gap-2">
          <button
            @click.stop="playPlaylist(playlist)"
            class="flex-1 px-3 py-2 border border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black transition-all text-sm font-mono flex items-center justify-center gap-2"
          >
            <span>▶</span>
            <span>PLAY</span>
          </button>
          <button
            @click.stop="managePlaylist(playlist)"
            class="flex-1 px-3 py-2 border border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black transition-all text-sm font-mono flex items-center justify-center gap-2"
          >
            <span>⚙</span>
            <span>MANAGE</span>
          </button>
        </div>

        <!-- Recent Songs Preview -->
        <div class="mt-4 space-y-1 max-h-24 overflow-y-auto border-t border-[#5a548d] pt-3">
          <div
            v-for="song in playlist.songs?.slice(0, 3)"
            :key="song.id"
            class="text-xs text-[#837dbd] font-mono truncate flex items-center gap-2"
          >
            <span class="text-[#5a548d]">›</span>
            <span class="truncate">{{ song.title }}</span>
          </div>
          <div v-if="!playlist.songs || playlist.songs.length === 0" class="text-xs text-[#5a548d] font-mono italic">
            No songs in playlist
          </div>
        </div>
      </div>
    </div>

    <!-- Create Playlist Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50 p-4">
      <div class="border-2 border-[#837dbd] bg-black w-full max-w-md">
        <!-- Modal Header -->
        <div class="border-b-2 border-[#837dbd] p-4 bg-black">
          <div class="flex justify-between items-center">
            <h3 class="text-lg text-[#d3ceff] font-mono">CREATE PLAYLIST</h3>
            <button
              @click="showCreateModal = false"
              class="text-[#837dbd] hover:text-red-400 transition-colors"
            >
              [X]
            </button>
          </div>
        </div>

        <!-- Form -->
        <div class="p-4">
          <div class="space-y-4">
            <div>
              <label class="block text-sm text-[#837dbd] mb-2 font-mono">NAME</label>
              <input
                v-model="newPlaylist.name"
                type="text"
                placeholder="Enter playlist name..."
                class="w-full px-3 py-2 bg-black border-2 border-[#837dbd] text-[#d3ceff] font-mono focus:outline-none focus:border-[#d3ceff]"
                @keyup.enter="createPlaylist"
              />
            </div>
            <div>
              <label class="block text-sm text-[#837dbd] mb-2 font-mono">DESCRIPTION (OPTIONAL)</label>
              <textarea
                v-model="newPlaylist.description"
                placeholder="Enter description..."
                rows="3"
                class="w-full px-3 py-2 bg-black border-2 border-[#837dbd] text-[#d3ceff] font-mono focus:outline-none focus:border-[#d3ceff] resize-none"
              ></textarea>
            </div>
          </div>

          <div class="flex justify-end gap-2 mt-6">
            <button
              @click="showCreateModal = false"
              class="px-4 py-2 border-2 border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black transition-all font-mono"
            >
              [CANCEL]
            </button>
            <button
              @click="createPlaylist"
              :disabled="!newPlaylist.name.trim()"
              class="px-4 py-2 bg-[#837dbd] text-black hover:bg-[#d3ceff] disabled:opacity-50 disabled:cursor-not-allowed transition-all font-mono"
            >
              [CREATE]
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Manage Playlist Modal -->
    <div v-if="selectedPlaylist" class="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50 p-4">
      <div class="border-2 border-[#837dbd] bg-black w-full max-w-4xl max-h-[80vh] overflow-hidden flex flex-col">
        <!-- Modal Header -->
        <div class="border-b-2 border-[#837dbd] p-4 bg-black flex-shrink-10">
          <div class="flex justify-between items-center">
            <h3 class="text-lg text-[#d3ceff] font-mono">
              MANAGE: {{ selectedPlaylist.name }}
            </h3>
            <button
              @click="selectedPlaylist = null"
              class="text-[#837dbd] hover:text-red-400 transition-colors"
            >
              [CLOSE]
            </button>
          </div>
        </div>

        <!-- Content -->
        <div class="flex-1 overflow-hidden flex">
          <!-- Available Songs -->
          <div class="w-1/2 border-r border-[#5a548d] flex flex-col">
            <div class="p-4 border-b border-[#5a548d]">
              <h4 class="text-sm text-[#d3ceff] mb-3 font-mono">AVAILABLE SONGS</h4>
              <input
                v-model="songSearch"
                type="text"
                placeholder="Search songs..."
                class="w-full px-3 py-2 bg-black border-2 border-[#837dbd] text-[#d3ceff] font-mono text-sm focus:outline-none focus:border-[#d3ceff]"
              />
            </div>
            <div class="flex-1 overflow-y-auto p-4">
              <div class="space-y-2">
                <div
                  v-for="song in filteredSongs"
                  :key="song.id"
                  class="flex justify-between items-center p-3 border border-[#5a548d] hover:border-[#837dbd] transition-all"
                >
                  <div class="flex-1 min-w-0">
                    <div class="text-[#d3ceff] text-sm font-mono truncate">{{ song.title }}</div>
                    <div class="text-[#837dbd] text-xs font-mono">{{ song.artist }} • {{ song.album }}</div>
                  </div>
                  <button
                    @click="addSongToPlaylist(song)"
                    :disabled="isSongInPlaylist(song.id)"
                    class="px-3 py-1 border border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black disabled:opacity-50 disabled:cursor-not-allowed transition-all text-xs font-mono"
                  >
                    {{ isSongInPlaylist(song.id) ? 'ADDED' : 'ADD' }}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Playlist Songs -->
          <div class="w-1/2 flex flex-col">
            <div class="p-4 border-b border-[#5a548d]">
              <h4 class="text-sm text-[#d3ceff] mb-3 font-mono">
                PLAYLIST SONGS ({{ selectedPlaylist.songs?.length || 0 }})
              </h4>
            </div>
            <div class="flex-1 overflow-y-auto p-4">
              <div class="space-y-2">
                <div
                  v-for="(song, index) in selectedPlaylist.songs"
                  :key="song.id"
                  class="flex justify-between items-center p-3 border border-[#5a548d] bg-[#837dbd] bg-opacity-10"
                >
                  <div class="flex items-center gap-3 flex-1 min-w-0">
                    <span class="text-[#837dbd] text-xs font-mono w-4">{{ index + 1 }}</span>
                    <div class="flex-1 min-w-0">
                      <div class="text-[#d3ceff] text-sm font-mono truncate">{{ song.title }}</div>
                      <div class="text-[#837dbd] text-xs font-mono">{{ song.artist }}</div>
                    </div>
                  </div>
                  <div class="flex gap-2">
                    <button
                      @click="moveSong(index, -1)"
                      :disabled="index === 0"
                      class="text-[#837dbd] hover:text-[#d3ceff] disabled:opacity-30 transition-colors text-sm"
                      title="Move up"
                    >
                      ↑
                    </button>
                    <button
                      @click="moveSong(index, 1)"
                      :disabled="index === selectedPlaylist.songs.length - 1"
                      class="text-[#837dbd] hover:text-[#d3ceff] disabled:opacity-30 transition-colors text-sm"
                      title="Move down"
                    >
                      ↓
                    </button>
                    <button
                      @click="removeSongFromPlaylist(song.id)"
                      class="text-[#837dbd] hover:text-red-400 transition-colors text-sm"
                      title="Remove"
                    >
                      [X]
                    </button>
                  </div>
                </div>
                <div v-if="!selectedPlaylist.songs || selectedPlaylist.songs.length === 0" class="text-center text-[#837dbd] text-sm font-mono py-4">
                  No songs in this playlist
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="border-t-2 border-[#837dbd] p-4 bg-black flex-shrink-10">
          <div class="flex justify-between items-center">
            <div class="text-sm text-[#837dbd] font-mono">
              Drag to reorder • Double-click to play
            </div>
            <button
              @click="savePlaylistChanges"
              class="px-4 py-2 bg-[#837dbd] text-black hover:bg-[#d3ceff] transition-all font-mono"
            >
              [SAVE CHANGES]
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface Song {
  id: number
  title: string
  artist: string
  album: string
  duration: number
  file_path: string
  rating: number
  play_count: number
}

interface Playlist {
  id: number
  name: string
  description: string
  songs: Song[]
  created_at: string
}

const playlists = ref<Playlist[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const showCreateModal = ref(false)
const selectedPlaylist = ref<Playlist | null>(null)
const songSearch = ref('')
const allSongs = ref<Song[]>([])

const newPlaylist = ref({
  name: '',
  description: ''
})

const totalSongs = computed(() => {
  return playlists.value.reduce((total, playlist) => total + (playlist.songs?.length || 0), 0)
})

const filteredSongs = computed(() => {
  if (!songSearch.value.trim()) return allSongs.value
  
  const search = songSearch.value.toLowerCase()
  return allSongs.value.filter(song =>
    song.title.toLowerCase().includes(search) ||
    song.artist.toLowerCase().includes(search) ||
    song.album.toLowerCase().includes(search)
  )
})

const fetchPlaylists = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await fetch('http://localhost:1323/api/playlists')
    if (!response.ok) throw new Error('Failed to fetch playlists')
    
    const data = await response.json()
    playlists.value = data
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Unknown error'
    console.error('Error fetching playlists:', err)
  } finally {
    loading.value = false
  }
}

const fetchAllSongs = async () => {
  try {
    const response = await fetch('http://localhost:1323/api/songs')
    if (!response.ok) throw new Error('Failed to fetch songs')
    
    const data = await response.json()
    allSongs.value = data
  } catch (err) {
    console.error('Error fetching songs:', err)
  }
}

const createPlaylist = async () => {
  if (!newPlaylist.value.name.trim()) return
  
  try {
    const response = await fetch('http://localhost:1323/api/playlists', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: newPlaylist.value.name,
        description: newPlaylist.value.description
      })
    })
    
    if (!response.ok) throw new Error('Failed to create playlist')
    
    const newPlaylistData = await response.json()
    playlists.value.push(newPlaylistData)
    showCreateModal.value = false
    newPlaylist.value = { name: '', description: '' }
  } catch (err) {
    console.error('Error creating playlist:', err)
    alert('Failed to create playlist')
  }
}

const editPlaylist = async (playlist: Playlist) => {
  const newName = prompt('Enter new playlist name:', playlist.name)
  if (!newName || newName === playlist.name) return
  
  try {
    const response = await fetch(`http://localhost:1323/api/playlists/${playlist.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ name: newName })
    })
    
    if (!response.ok) throw new Error('Failed to update playlist')
    
    const updatedPlaylist = await response.json()
    const index = playlists.value.findIndex(p => p.id === playlist.id)
    if (index !== -1) {
      playlists.value[index] = updatedPlaylist
    }
  } catch (err) {
    console.error('Error updating playlist:', err)
    alert('Failed to update playlist')
  }
}

const deletePlaylist = async (id: number) => {
  if (!confirm('Are you sure you want to delete this playlist?')) return
  
  try {
    const response = await fetch(`http://localhost:1323/api/playlists/${id}`, {
      method: 'DELETE'
    })
    
    if (!response.ok) throw new Error('Failed to delete playlist')
    
    playlists.value = playlists.value.filter(p => p.id !== id)
  } catch (err) {
    console.error('Error deleting playlist:', err)
    alert('Failed to delete playlist')
  }
}

const playPlaylist = async (playlist: Playlist) => {
  if (!playlist.songs || playlist.songs.length === 0) {
    alert('Playlist is empty')
    return
  }
  
  // In a real app, you would play the first song and queue the rest
  console.log('Playing playlist:', playlist.name)
  // Emit event to play first song
}

const managePlaylist = (playlist: Playlist) => {
  selectedPlaylist.value = { ...playlist, songs: [...playlist.songs] }
}

const isSongInPlaylist = (songId: number) => {
  return selectedPlaylist.value?.songs.some(s => s.id === songId) || false
}

const addSongToPlaylist = (song: Song) => {
  if (selectedPlaylist.value && !isSongInPlaylist(song.id)) {
    selectedPlaylist.value.songs.push({ ...song })
  }
}

const removeSongFromPlaylist = (songId: number) => {
  if (selectedPlaylist.value) {
    selectedPlaylist.value.songs = selectedPlaylist.value.songs.filter(s => s.id !== songId)
  }
}

const moveSong = (index: number, direction: number) => {
  if (!selectedPlaylist.value) return
  
  const newIndex = index + direction
  if (newIndex >= 0 && newIndex < selectedPlaylist.value.songs.length) {
    const songs = selectedPlaylist.value.songs
    const [song] = songs.splice(index, 1)
    songs.splice(newIndex, 0, song)
  }
}

const savePlaylistChanges = async () => {
  if (!selectedPlaylist.value) return
  
  try {
    // First update the playlist metadata
    const response = await fetch(`http://localhost:1323/api/playlists/${selectedPlaylist.value.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: selectedPlaylist.value.name,
        description: selectedPlaylist.value.description
      })
    })
    
    if (!response.ok) throw new Error('Failed to update playlist')
    
    // Update the original playlist in the list
    const index = playlists.value.findIndex(p => p.id === selectedPlaylist.value!.id)
    if (index !== -1) {
      playlists.value[index] = { ...selectedPlaylist.value }
    }
    
    selectedPlaylist.value = null
  } catch (err) {
    console.error('Error saving playlist changes:', err)
    alert('Failed to save changes')
  }
}

const openPlaylist = (playlist: Playlist) => {
  // Navigate to playlist detail view
  console.log('Opening playlist:', playlist.name)
}

const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

onMounted(() => {
  fetchPlaylists()
  fetchAllSongs()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #000;
  border: 1px solid #5a548d;
}

::-webkit-scrollbar-thumb {
  background: #837dbd;
}

::-webkit-scrollbar-thumb:hover {
  background: #d3ceff;
}
</style>