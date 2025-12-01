<!-- App.vue -->
<template>
  <div class="app-container">
    <!-- Audio Player Element (Hidden) -->
    <audio 
      ref="audioElement" 
      preload="auto"
      @timeupdate="onTimeUpdate"
      @ended="onEnded"
      @loadedmetadata="onLoadedMetadata"
      @error="onAudioError"
    ></audio>
    
    <div class="h-screen flex flex-col bg-black text-white">
      <!-- Header -->
      <header class="border-b-2 border-[#837dbd] p-4 bg-black">
        <div class="container mx-auto flex justify-between items-center">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 bg-[#837dbd] rounded-full animate-pulse"></div>
            <h1 class="text-xl font-bold text-[#d3ceff] font-mono">MSX AUDIO PLAYER</h1>
          </div>
          
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2">
              <div class="w-2 h-2 rounded-full animate-pulse" :class="isConnected ? 'bg-green-400' : 'bg-red-400'"></div>
              <span class="text-sm font-mono" :class="isConnected ? 'text-green-400' : 'text-red-400'">
                {{ isConnected ? 'CONNECTED' : 'DISCONNECTED' }}
              </span>
            </div>
            <button
              @click="reconnectAudio"
              class="px-3 py-1 border border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black transition-all text-xs font-mono"
            >
              RECONNECT
            </button>
          </div>
        </div>
      </header>

      <!-- Main Content -->
      <main class="flex-1 overflow-hidden">
        <router-view 
          :audio-element="audioElement"
          :current-song="currentSong"
          :is-playing="isPlaying"
          :volume="volume"
          @play="playSong"
          @pause="pausePlayer"
          @volume-change="setVolume"
          @seek="seekTo"
        />
      </main>

      <!-- Now Playing Bar (Fixed at bottom) -->
      <div v-if="currentSong" class="border-t-2 border-[#837dbd] bg-black p-4">
        <div class="container mx-auto">
          <div class="flex items-center justify-between">
            <!-- Song Info -->
            <div class="flex items-center gap-4 flex-1">
              <div class="w-12 h-12 bg-gradient-to-br from-[#837dbd] to-[#d3ceff] flex items-center justify-center">
                <span class="text-2xl">♪</span>
              </div>
              <div class="flex-1 min-w-0">
                <div class="text-[#d3ceff] font-mono truncate">{{ currentSong.title }}</div>
                <div class="text-[#837dbd] text-sm font-mono truncate">{{ currentSong.artist }} • {{ currentSong.album }}</div>
              </div>
            </div>

            <!-- Playback Controls -->
            <div class="flex items-center gap-6">
              <button @click="previousSong" class="text-2xl text-[#837dbd] hover:text-[#d3ceff]">
                ⏮
              </button>
              <button
                v-if="!isPlaying"
                @click="resumePlayer"
                class="w-12 h-12 bg-[#837dbd] text-black rounded-full flex items-center justify-center hover:bg-[#d3ceff]"
              >
                ▶
              </button>
              <button
                v-else
                @click="pausePlayer"
                class="w-12 h-12 bg-[#837dbd] text-black rounded-full flex items-center justify-center hover:bg-[#d3ceff]"
              >
                ⏸
              </button>
              <button @click="nextSong" class="text-2xl text-[#837dbd] hover:text-[#d3ceff]">
                ⏭
              </button>
            </div>

            <!-- Progress and Volume -->
            <div class="flex-1 max-w-md space-y-2">
              <div class="flex items-center gap-3">
                <span class="text-xs text-[#837dbd] font-mono w-12">{{ formatTime(currentTime) }}</span>
                <input
                  v-model="currentTime"
                  @input="onSeekInput"
                  type="range"
                  min="0"
                  :max="duration"
                  step="1"
                  class="flex-1 h-1 bg-[#5a548d] rounded-lg appearance-none cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-[#d3ceff]"
                />
                <span class="text-xs text-[#837dbd] font-mono w-12">{{ formatTime(duration) }}</span>
              </div>
              
              <div class="flex items-center gap-2">
                <span class="text-xs text-[#837dbd] font-mono">VOL</span>
                <input
                  v-model="volume"
                  @input="setVolume"
                  type="range"
                  min="0"
                  max="1"
                  step="0.01"
                  class="flex-1 h-1 bg-[#5a548d] rounded-lg appearance-none cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-[#d3ceff]"
                />
                <span class="text-xs text-[#837dbd] font-mono w-12">{{ Math.round(volume * 100) }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Connection Error Modal -->
    <div v-if="showConnectionError" class="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center z-50">
      <div class="border-2 border-red-400 bg-black p-8 max-w-md">
        <div class="text-center space-y-4">
          <div class="text-4xl text-red-400">⚠</div>
          <h3 class="text-xl text-red-400 font-mono">AUDIO CONNECTION ERROR</h3>
          <p class="text-[#837dbd]">
            Unable to connect to audio stream. Please check:
          </p>
          <ul class="text-left text-sm text-[#837dbd] space-y-2 list-disc list-inside">
            <li>Server is running on localhost:1323</li>
            <li>CORS is properly configured</li>
            <li>Audio files are accessible</li>
          </ul>
          <div class="flex gap-4 justify-center mt-6">
            <button
              @click="reconnectAudio"
              class="px-4 py-2 bg-[#837dbd] text-black hover:bg-[#d3ceff] font-mono"
            >
              RETRY
            </button>
            <button
              @click="showConnectionError = false"
              class="px-4 py-2 border border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black font-mono"
            >
              CLOSE
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const audioElement = ref<HTMLAudioElement>()
const isConnected = ref(false)
const isPlaying = ref(false)
const currentSong = ref<any>(null)
const currentTime = ref(0)
const duration = ref(0)
const volume = ref(0.7)
const showConnectionError = ref(false)

// Fetch initial player state
const fetchPlayerState = async () => {
  try {
    const response = await fetch('http://localhost:1323/api/player/state')
    if (response.ok) {
      const state = await response.json()
      isConnected.value = state.is_connected
      currentSong.value = state.current_song
      isPlaying.value = state.is_playing
      currentTime.value = state.current_time
      duration.value = state.duration
      volume.value = state.volume
      
      // If there's a current song, load it
      if (currentSong.value && audioElement.value) {
        await loadAudioStream(currentSong.value.id)
      }
    }
  } catch (error) {
    console.error('Failed to fetch player state:', error)
    isConnected.value = false
  }
}

// Load audio stream from API
const loadAudioStream = async (songId: number) => {
  if (!audioElement.value) return
  
  try {
    // Clear previous source
    audioElement.value.pause()
    audioElement.value.src = ''
    
    // Load stream from API endpoint
    const streamUrl = `http://localhost:1323/api/songs/${songId}/stream`
    audioElement.value.src = streamUrl
    audioElement.value.volume = volume.value
    
    // Add timestamp to prevent caching issues
    audioElement.value.src += `?t=${Date.now()}`
    
    await audioElement.value.load()
    isConnected.value = true
    showConnectionError.value = false
  } catch (error) {
    console.error('Failed to load audio stream:', error)
    showConnectionError.value = true
    isConnected.value = false
  }
}

// Play song by ID
const playSong = async (song: any) => {
  if (!song) return
  
  currentSong.value = song
  
  try {
    // First, try to set the song via API
    const response = await fetch(`http://localhost:1323/api/player/play/${song.id}`, {
      method: 'POST'
    })
    
    if (response.ok) {
      // Then load the stream
      await loadAudioStream(song.id)
      
      if (audioElement.value) {
        await audioElement.value.play()
        isPlaying.value = true
      }
    }
  } catch (error) {
    console.error('Failed to play song:', error)
    // Fallback: load stream directly
    await loadAudioStream(song.id)
    if (audioElement.value) {
      await audioElement.value.play()
      isPlaying.value = true
    }
  }
}

// Resume player
const resumePlayer = async () => {
  if (!audioElement.value) return
  
  try {
    if (currentSong.value) {
      // Resume via API
      const response = await fetch('http://localhost:1323/api/player/resume', {
        method: 'POST'
      })
      
      if (response.ok) {
        await audioElement.value.play()
        isPlaying.value = true
      }
    } else {
      // Get first song from library and play it
      const response = await fetch('http://localhost:1323/api/songs?limit=1')
      if (response.ok) {
        const songs = await response.json()
        if (songs.length > 0) {
          await playSong(songs[0])
        }
      }
    }
  } catch (error) {
    console.error('Failed to resume:', error)
    // Fallback: play directly
    if (audioElement.value) {
      await audioElement.value.play()
      isPlaying.value = true
    }
  }
}

// Pause player
const pausePlayer = async () => {
  if (!audioElement.value) return
  
  try {
    // Pause via API
    await fetch('http://localhost:1323/api/player/pause', {
      method: 'POST'
    })
    
    audioElement.value.pause()
    isPlaying.value = false
  } catch (error) {
    console.error('Failed to pause:', error)
    // Fallback: pause directly
    audioElement.value.pause()
    isPlaying.value = false
  }
}

// Previous song
const previousSong = async () => {
  try {
    const response = await fetch('http://localhost:1323/api/player/previous', {
      method: 'POST'
    })
    
    if (response.ok) {
      await fetchPlayerState()
    }
  } catch (error) {
    console.error('Failed to go to previous song:', error)
  }
}

// Next song
const nextSong = async () => {
  try {
    const response = await fetch('http://localhost:1323/api/player/next', {
      method: 'POST'
    })
    
    if (response.ok) {
      await fetchPlayerState()
    }
  } catch (error) {
    console.error('Failed to go to next song:', error)
  }
}

// Set volume
const setVolume = (event?: Event) => {
  if (!audioElement.value) return
  
  const newVolume = event ? parseFloat((event.target as HTMLInputElement).value) : volume.value
  volume.value = newVolume
  audioElement.value.volume = newVolume
  
  // Also update on server
  fetch(`http://localhost:1323/api/player/volume?value=${newVolume}`, {
    method: 'POST'
  }).catch(console.error)
}

// Seek to position
const seekTo = (time: number) => {
  if (!audioElement.value) return
  
  currentTime.value = time
  audioElement.value.currentTime = time
  
  // Update on server
  fetch(`http://localhost:1323/api/player/seek?time=${time}`, {
    method: 'POST'
  }).catch(console.error)
}

const onSeekInput = (event: Event) => {
  const time = parseFloat((event.target as HTMLInputElement).value)
  seekTo(time)
}

// Format time (MM:SS)
const formatTime = (seconds: number) => {
  if (!seconds || isNaN(seconds)) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

// Event listeners
const onTimeUpdate = () => {
  if (audioElement.value) {
    currentTime.value = audioElement.value.currentTime
  }
}

const onEnded = () => {
  isPlaying.value = false
  // Auto-play next song
  nextSong()
}

const onLoadedMetadata = () => {
  if (audioElement.value) {
    duration.value = audioElement.value.duration
  }
}

const onAudioError = (event: Event) => {
  console.error('Audio error:', event)
  showConnectionError.value = true
  isConnected.value = false
}

const reconnectAudio = async () => {
  showConnectionError.value = false
  await fetchPlayerState()
}

// Initialize
onMounted(async () => {
  await fetchPlayerState()
  
  // Set initial volume
  if (audioElement.value) {
    audioElement.value.volume = volume.value
  }
  
  // Poll for updates
  const interval = setInterval(fetchPlayerState, 5000)
  
  onUnmounted(() => {
    clearInterval(interval)
  })
})

// Watch for route changes to update player
watch(() => router.currentRoute.value, fetchPlayerState)
</script>

<style scoped>
/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
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

/* Range slider styling */
input[type="range"] {
  -webkit-appearance: none;
  appearance: none;
  background: transparent;
  cursor: pointer;
}

input[type="range"]::-webkit-slider-track {
  background: #5a548d;
  height: 0.25rem;
  border-radius: 0.25rem;
}

input[type="range"]::-moz-range-track {
  background: #5a548d;
  height: 0.25rem;
  border-radius: 0.25rem;
}

input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  height: 1rem;
  width: 1rem;
  background-color: #d3ceff;
  border-radius: 50%;
  margin-top: -0.375rem;
}

input[type="range"]::-moz-range-thumb {
  border: none;
  border-radius: 50%;
  height: 1rem;
  width: 1rem;
  background-color: #d3ceff;
}
</style>