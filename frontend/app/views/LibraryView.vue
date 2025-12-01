<!-- views/LibraryView.vue -->
<template>
  <div class="h-full flex flex-col p-4 space-y-4">
    <!-- Library Header -->
    <div class="border-2 border-[#837dbd] p-4 bg-black">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-xl text-[#d3ceff] font-bold tracking-wide flex items-center gap-2">
            <span class="animate-pulse">&gt;&gt;</span>
            MUSIC LIBRARY
          </h1>
          <p class="text-sm text-[#837dbd] font-mono mt-1">
            {{ filteredSongs.length }} of {{ songs.length }} tracks
            <span v-if="searchQuery">• searching: "{{ searchQuery }}"</span>
          </p>
        </div>
        
        <div class="flex items-center gap-4">
          <!-- Search -->
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="SEARCH LIBRARY..."
              class="bg-black border-2 border-[#837dbd] text-[#d3ceff] px-4 py-2 font-mono focus:outline-none focus:border-[#d3ceff] w-64"
            />
            <span class="absolute right-3 top-2.5 text-[#837dbd]">⌕</span>
          </div>
          
          <!-- View Toggle -->
          <div class="flex border-2 border-[#837dbd]">
            <button
              @click="viewMode = 'grid'"
              :class="[
                'px-3 py-2 font-mono text-sm transition-all',
                viewMode === 'grid' 
                  ? 'bg-[#837dbd] text-black' 
                  : 'bg-black text-[#837dbd] hover:bg-[#837dbd] hover:text-black'
              ]"
            >
              GRID
            </button>
            <button
              @click="viewMode = 'list'"
              :class="[
                'px-3 py-2 font-mono text-sm transition-all',
                viewMode === 'list' 
                  ? 'bg-[#837dbd] text-black' 
                  : 'bg-black text-[#837dbd] hover:bg-[#837dbd] hover:text-black'
              ]"
            >
              LIST
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Stats -->
    <div class="grid grid-cols-4 gap-4">
      <div class="border-2 border-[#837dbd] p-4 bg-black">
        <div class="text-xs text-[#837dbd] font-mono">TOTAL TRACKS</div>
        <div class="text-2xl text-[#d3ceff] font-mono">{{ songs.length }}</div>
      </div>
      <div class="border-2 border-[#837dbd] p-4 bg-black">
        <div class="text-xs text-[#837dbd] font-mono">TOTAL TIME</div>
        <div class="text-2xl text-[#d3ceff] font-mono">{{ formatTotalTime() }}</div>
      </div>
      <div class="border-2 border-[#837dbd] p-4 bg-black">
        <div class="text-xs text-[#837dbd] font-mono">ARTISTS</div>
        <div class="text-2xl text-[#d3ceff] font-mono">{{ uniqueArtists }}</div>
      </div>
      <div class="border-2 border-[#837dbd] p-4 bg-black">
        <div class="text-xs text-[#837dbd] font-mono">ALBUMS</div>
        <div class="text-2xl text-[#d3ceff] font-mono">{{ uniqueAlbums }}</div>
      </div>
    </div>

    <!-- Controls -->
    <div class="flex gap-3">
      <button
        @click="playAll"
        :disabled="filteredSongs.length === 0"
        class="px-4 py-2 bg-[#837dbd] text-black hover:bg-[#d3ceff] disabled:opacity-50 disabled:cursor-not-allowed transition-all font-mono flex items-center gap-2"
      >
        <span>▶</span>
        <span>PLAY ALL</span>
      </button>
      <button
        @click="shuffleAll"
        :disabled="filteredSongs.length === 0"
        class="px-4 py-2 border-2 border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black disabled:opacity-50 disabled:cursor-not-allowed transition-all font-mono flex items-center gap-2"
      >
        <span>🔀</span>
        <span>SHUFFLE ALL</span>
      </button>
      
      <!-- Sort Dropdown -->
      <select
        v-model="sortBy"
        class="bg-black border-2 border-[#837dbd] text-[#d3ceff] px-3 py-2 font-mono focus:outline-none focus:border-[#d3ceff]"
      >
        <option value="title">SORT: TITLE</option>
        <option value="artist">SORT: ARTIST</option>
        <option value="album">SORT: ALBUM</option>
        <option value="rating">SORT: RATING</option>
        <option value="duration">SORT: DURATION</option>
      </select>
      
      <!-- Filter Dropdown -->
      <select
        v-model="genreFilter"
        class="bg-black border-2 border-[#837dbd] text-[#d3ceff] px-3 py-2 font-mono focus:outline-none focus:border-[#d3ceff]"
      >
        <option value="">GENRE: ALL</option>
        <option v-for="genre in uniqueGenres" :key="genre" :value="genre">
          GENRE: {{ genre.toUpperCase() }}
        </option>
      </select>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-4">
        <div class="text-4xl text-[#837dbd] animate-pulse">♫</div>
        <div class="text-[#d3ceff] font-mono">LOADING MUSIC LIBRARY...</div>
        <div class="text-sm text-[#837dbd] font-mono">please wait</div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-4">
        <div class="text-4xl text-red-400">⚠</div>
        <div class="text-xl text-red-400 font-mono">FAILED TO LOAD LIBRARY</div>
        <div class="text-[#837dbd]">{{ error }}</div>
        <button
          @click="fetchSongs"
          class="mt-4 px-4 py-2 border-2 border-red-400 text-red-400 hover:bg-red-400 hover:text-black transition-all font-mono"
        >
          [RETRY]
        </button>
      </div>
    </div>

    <!-- Grid View -->
    <div v-else-if="viewMode === 'grid' && filteredSongs.length > 0" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 flex-1 overflow-y-auto">
      <div
        v-for="song in filteredSongs"
        :key="song.id"
        class="border-2 border-[#837dbd] bg-black hover:border-[#d3ceff] transition-all group cursor-pointer"
        @click="playSong(song)"
        @dblclick="queueSong(song)"
      >
        <!-- Album Art -->
        <div class="relative aspect-square bg-gradient-to-br from-[#5a548d] to-[#837dbd]">
          <div class="absolute inset-0 flex items-center justify-center">
            <span class="text-6xl text-black/50">♪</span>
          </div>
          
          <!-- Overlay with play button -->
          <div class="absolute inset-0 bg-black/70 opacity-0 group-hover:opacity-100 transition-all flex items-center justify-center">
            <div class="w-16 h-16 bg-[#d3ceff] rounded-full flex items-center justify-center text-black text-2xl transform scale-75 group-hover:scale-100 transition-all">
              ▶
            </div>
          </div>
          
          <!-- Current playing indicator -->
          <div v-if="currentSong && currentSong.id === song.id" class="absolute top-2 right-2">
            <div class="w-3 h-3 bg-green-400 animate-pulse rounded-full"></div>
          </div>
        </div>
        
        <!-- Song Info -->
        <div class="p-4">
          <div class="text-[#d3ceff] font-mono truncate" :title="song.title">{{ song.title }}</div>
          <div class="text-sm text-[#837dbd] font-mono truncate" :title="song.artist">{{ song.artist }}</div>
          <div class="text-xs text-[#5a548d] font-mono truncate" :title="song.album">{{ song.album }}</div>
          
          <div class="flex justify-between items-center mt-3">
            <div class="text-xs text-[#837dbd] font-mono">{{ formatDuration(song.duration) }}</div>
            <div class="flex items-center gap-1">
              <span
                v-for="star in 5"
                :key="star"
                class="text-sm"
                :class="star <= Math.floor((song.rating || 0) / 2) ? 'text-yellow-400' : 'text-[#5a548d]'"
              >
                ★
              </span>
            </div>
          </div>
          
          <!-- Quick Actions -->
          <div class="flex gap-2 mt-3 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              @click.stop="queueSong(song)"
              class="flex-1 px-2 py-1 border border-[#837dbd] text-[#837dbd] text-xs font-mono hover:bg-[#837dbd] hover:text-black transition-all"
            >
              QUEUE
            </button>
            <button
              @click.stop="rateSong(song.id, 10)"
              class="px-2 py-1 border border-yellow-400 text-yellow-400 text-xs font-mono hover:bg-yellow-400 hover:text-black transition-all"
              title="Rate 10/10"
            >
              ★
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- List View -->
    <div v-else-if="viewMode === 'list' && filteredSongs.length > 0" class="border-2 border-[#837dbd] bg-black flex-1 overflow-y-auto">
      <div class="grid grid-cols-12 border-b-2 border-[#5a548d] text-xs text-[#837dbd] font-mono uppercase p-3 sticky top-0 bg-black z-10">
        <div class="col-span-1 text-center">#</div>
        <div class="col-span-4">TITLE</div>
        <div class="col-span-3">ARTIST</div>
        <div class="col-span-2">ALBUM</div>
        <div class="col-span-1 text-center">TIME</div>
        <div class="col-span-1 text-center">RATING</div>
      </div>
      
      <div
        v-for="(song, index) in filteredSongs"
        :key="song.id"
        :class="[
          'grid grid-cols-12 p-3 border-b border-[#5a548d] hover:bg-[#5a548d]/20 cursor-pointer transition-all group',
          currentSong && currentSong.id === song.id ? 'bg-[#837dbd]/10' : ''
        ]"
        @click="playSong(song)"
        @dblclick="queueSong(song)"
      >
        <!-- Index -->
        <div class="col-span-1 flex items-center justify-center">
          <div v-if="currentSong && currentSong.id === song.id" class="w-6 h-6 flex items-center justify-center">
            <div class="w-2 h-2 bg-green-400 animate-pulse rounded-full"></div>
          </div>
          <span v-else class="text-sm text-[#837dbd]">{{ index + 1 }}</span>
        </div>
        
        <!-- Title -->
        <div class="col-span-4 flex items-center gap-3">
          <div class="w-10 h-10 bg-gradient-to-br from-[#5a548d] to-[#837dbd] flex items-center justify-center flex-shrink-0">
            <span class="text-lg">♪</span>
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-[#d3ceff] font-mono truncate">{{ song.title }}</div>
            <div v-if="song.track_number" class="text-xs text-[#5a548d] font-mono">
              Track {{ song.track_number }}
            </div>
          </div>
        </div>
        
        <!-- Artist -->
        <div class="col-span-3 flex items-center">
          <span class="text-[#837dbd] font-mono truncate">{{ song.artist }}</span>
        </div>
        
        <!-- Album -->
        <div class="col-span-2 flex items-center">
          <span class="text-[#837dbd] font-mono truncate">{{ song.album }}</span>
        </div>
        
        <!-- Duration -->
        <div class="col-span-1 flex items-center justify-center">
          <span class="text-sm text-[#837dbd] font-mono">{{ formatDuration(song.duration) }}</span>
        </div>
        
        <!-- Rating -->
        <div class="col-span-1 flex items-center justify-center">
          <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              v-for="star in 5"
              :key="star"
              @click.stop="rateSong(song.id, star * 2)"
              class="text-sm hover:scale-125 transition-transform"
              :class="star <= Math.floor((song.rating || 0) / 2) ? 'text-yellow-400' : 'text-[#5a548d]'"
            >
              ★
            </button>
          </div>
          <span v-if="!song.rating" class="text-xs text-[#5a548d]">-</span>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredSongs.length === 0" class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-4">
        <div class="text-6xl text-[#5a548d]">♫</div>
        <div class="text-xl text-[#d3ceff] font-mono">NO SONGS FOUND</div>
        <div v-if="searchQuery" class="text-[#837dbd]">
          No results for "{{ searchQuery }}"
        </div>
        <div v-else class="text-[#837dbd]">
          Upload some music to get started
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'

interface Song {
  id: number
  title: string
  artist: string
  album: string
  duration: number
  genre: string
  rating: number
  track_number: number
  year: number
  file_path: string
  play_count: number
  created_at: string
}

// Props
interface Props {
  audioElement?: HTMLAudioElement
  currentSong?: Song
  isPlaying?: boolean
}

const props = defineProps<Props>()

// Emits
const emit = defineEmits(['play', 'pause', 'volume-change', 'seek'])

// State
const songs = ref<Song[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const searchQuery = ref('')
const sortBy = ref('title')
const genreFilter = ref('')
const viewMode = ref<'grid' | 'list'>('grid')

// Computed
const uniqueArtists = computed(() => {
  const artists = new Set(songs.value.map(song => song.artist))
  return artists.size
})

const uniqueAlbums = computed(() => {
  const albums = new Set(songs.value.map(song => song.album))
  return albums.size
})

const uniqueGenres = computed(() => {
  const genres = new Set(songs.value.map(song => song.genre).filter(Boolean))
  return Array.from(genres).sort()
})

const filteredSongs = computed(() => {
  let filtered = songs.value

  // Apply search
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(song =>
      song.title.toLowerCase().includes(query) ||
      song.artist.toLowerCase().includes(query) ||
      song.album.toLowerCase().includes(query) ||
      (song.genre && song.genre.toLowerCase().includes(query))
    )
  }

  // Apply genre filter
  if (genreFilter.value) {
    filtered = filtered.filter(song => song.genre === genreFilter.value)
  }

  // Apply sorting
  filtered = [...filtered].sort((a, b) => {
    switch (sortBy.value) {
      case 'title':
        return a.title.localeCompare(b.title)
      case 'artist':
        return a.artist.localeCompare(b.artist)
      case 'album':
        return a.album.localeCompare(b.album)
      case 'rating':
        return (b.rating || 0) - (a.rating || 0)
      case 'duration':
        return b.duration - a.duration
      default:
        return 0
    }
  })

  return filtered
})

// Methods
const fetchSongs = async () => {
  loading.value = true
  error.value = null
  
  try {
    const response = await fetch('http://localhost:1323/api/songs')
    if (!response.ok) throw new Error(`HTTP ${response.status}`)
    
    const data = await response.json()
    songs.value = data.map((song: any) => ({
      ...song,
      // Ensure file_path points to stream endpoint
      file_path: song.file_path.startsWith('http') 
        ? song.file_path 
        : `http://localhost:1323/api/songs/${song.id}/stream`
    }))
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load songs'
    console.error('Error fetching songs:', err)
  } finally {
    loading.value = false
  }
}

const formatDuration = (seconds: number) => {
  if (!seconds || isNaN(seconds)) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const formatTotalTime = () => {
  const totalSeconds = songs.value.reduce((sum, song) => sum + (song.duration || 0), 0)
  const hours = Math.floor(totalSeconds / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  
  if (hours > 0) {
    return `${hours}h ${minutes}m`
  } else {
    return `${minutes}m`
  }
}

const playSong = async (song: Song) => {
  // Construct proper stream URL
  const streamUrl = `http://localhost:1323/api/songs/${song.id}/stream?t=${Date.now()}`
  
  emit('play', {
    ...song,
    file_path: streamUrl
  })
}

const queueSong = async (song: Song) => {
  try {
    const response = await fetch('http://localhost:1323/api/player/queue', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ songId: song.id })
    })
    
    if (response.ok) {
      console.log('Song queued:', song.title)
    }
  } catch (err) {
    console.error('Failed to queue song:', err)
  }
}

const playAll = async () => {
  if (filteredSongs.value.length === 0) return
  
  try {
    const response = await fetch('http://localhost:1323/api/player/playlist', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        songs: filteredSongs.value.map(s => s.id),
        startIndex: 0 
      })
    })
    
    if (response.ok && filteredSongs.value[0]) {
      await playSong(filteredSongs.value[0])
    }
  } catch (err) {
    console.error('Failed to play all:', err)
    // Fallback: play first song
    if (filteredSongs.value[0]) {
      await playSong(filteredSongs.value[0])
    }
  }
}

const shuffleAll = async () => {
  if (filteredSongs.value.length === 0) return
  
  try {
    const shuffled = [...filteredSongs.value].sort(() => Math.random() - 0.5)
    
    const response = await fetch('http://localhost:1323/api/player/playlist', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ 
        songs: shuffled.map(s => s.id),
        startIndex: 0 
      })
    })
    
    if (response.ok && shuffled[0]) {
      await playSong(shuffled[0])
    }
  } catch (err) {
    console.error('Failed to shuffle all:', err)
  }
}

const rateSong = async (songId: number, rating: number) => {
  try {
    const response = await fetch(`http://localhost:1323/api/songs/${songId}/rate`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ rating })
    })
    
    if (response.ok) {
      // Update local state
      const songIndex = songs.value.findIndex(s => s.id === songId)
      if (songIndex !== -1) {
        songs.value[songIndex].rating = rating
      }
    }
  } catch (err) {
    console.error('Failed to rate song:', err)
  }
}

// Lifecycle
onMounted(() => {
  fetchSongs()
  
  // Refresh every minute
  const interval = setInterval(fetchSongs, 60000)
  
  return () => clearInterval(interval)
})

// Watch for search changes
watch(searchQuery, () => {
  // Debounce would be good here
})
</script>

<style scoped>
/* Custom scrollbar */
.border-\[#5a548d\]::-webkit-scrollbar {
  width: 6px;
}

.border-\[#5a548d\]::-webkit-scrollbar-track {
  background: #000;
  border: 1px solid #5a548d;
}

.border-\[#5a548d\]::-webkit-scrollbar-thumb {
  background: #837dbd;
}

.border-\[#5a548d\]::-webkit-scrollbar-thumb:hover {
  background: #d3ceff;
}
</style>