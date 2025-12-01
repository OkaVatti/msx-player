<!-- app.vue -->
<template>
  <div class="h-screen flex flex-col bg-black text-white font-mono overflow-hidden">
    <!-- Hidden Audio Element -->
    <audio
      ref="audioElement"
      @timeupdate="handleTimeUpdate"
      @loadedmetadata="handleLoadedMetadata"
      @ended="handleEnded"
      @error="handleError"
    ></audio>

    <!-- Top Menu Bar -->
    <div class="h-12 bg-black border-b border-[#837dbd] flex items-center px-4 flex-shrink-10">
      <div class="flex items-center gap-6 text-sm">
        <div class="text-[#d3ceff] font-bold tracking-wide">MSX PLAYER</div>
        <NavMenu @upload="showUploadModal = true" @change-view="changeView" />
      </div>
      
      <div class="ml-auto">
        <ConnectionStatus />
      </div>
    </div>

    <!-- Main Layout -->
    <div class="flex flex-1 overflow-hidden">
      <!-- Left Sidebar -->
      <LeftSidebar
        :current-view="currentView"
        :playlists="libraryStore.playlists"
        @change-view="changeView"
        @create-playlist="showCreatePlaylistModal = true"
        @select-playlist="selectPlaylist"
      />

      <!-- Main Content Area -->
      <div class="flex-1 flex flex-col overflow-hidden">
        <MainView :current-view="currentView" @play-song="handlePlaySong" />
      </div>

      <!-- Right Sidebar (Now Playing) -->
      <RightSidebar />
    </div>

    <!-- Bottom Now Playing Bar -->
    <NowPlayingBar />

    <!-- Upload Modal -->
    <UploadModal v-if="showUploadModal" @close="showUploadModal = false" />

    <!-- Create Playlist Modal -->
    <CreatePlaylistModal
      v-if="showCreatePlaylistModal"
      @close="showCreatePlaylistModal = false"
      @created="handlePlaylistCreated"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { usePlayerStore } from '../stores/player'
import { useLibraryStore } from '../stores/library'
import type { Song } from '../types'

const playerStore = usePlayerStore()
const libraryStore = useLibraryStore()

const audioElement = ref<HTMLAudioElement>()
const currentView = ref('library')
const showUploadModal = ref(false)
const showCreatePlaylistModal = ref(false)

// Initialize audio element in store
watch(audioElement, (el) => {
  if (el) {
    playerStore.setAudioElement(el)
  }
})

// Audio Event Handlers
const handleTimeUpdate = () => {
  if (audioElement.value) {
    playerStore.currentTime = audioElement.value.currentTime
  }
}

const handleLoadedMetadata = () => {
  if (audioElement.value) {
    playerStore.duration = audioElement.value.duration
  }
}

const handleEnded = () => {
  playerStore.handleSongEnded()
}

const handleError = (event: Event) => {
  console.error('Audio error:', event)
  playerStore.isPlaying = false
}

// Navigation
const changeView = (view: string) => {
  currentView.value = view
}

const selectPlaylist = (playlistId: number) => {
  currentView.value = 'playlist'
  // Store selected playlist ID for playlist view
  libraryStore.selectedPlaylistId = playlistId
}

// Play Song Handler
const handlePlaySong = async (song: Song) => {
  await playerStore.playSong(song)
}

const handlePlaylistCreated = () => {
  libraryStore.fetchPlaylists()
  showCreatePlaylistModal.value = false
}

// Keyboard Shortcuts
const handleKeydown = (event: KeyboardEvent) => {
  if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
    return
  }

  switch (event.key) {
    case ' ':
      event.preventDefault()
      if (playerStore.currentSong) {
        if (playerStore.isPlaying) {
          playerStore.pause()
        } else {
          playerStore.playSong(playerStore.currentSong)
        }
      }
      break
    case 'ArrowRight':
      event.preventDefault()
      if (playerStore.currentSong) {
        const newTime = Math.min(playerStore.currentTime + 10, playerStore.duration)
        playerStore.seek(newTime)
      }
      break
    case 'ArrowLeft':
      event.preventDefault()
      if (playerStore.currentSong) {
        const newTime = Math.max(playerStore.currentTime - 10, 0)
        playerStore.seek(newTime)
      }
      break
    case 'ArrowUp':
      event.preventDefault()
      playerStore.setVolume(Math.min(playerStore.volume + 0.1, 1))
      break
    case 'ArrowDown':
      event.preventDefault()
      playerStore.setVolume(Math.max(playerStore.volume - 0.1, 0))
      break
    case 'n':
    case 'N':
      event.preventDefault()
      playerStore.nextSong()
      break
    case 'p':
    case 'P':
      event.preventDefault()
      playerStore.previousSong()
      break
  }
}

// Lifecycle
onMounted(() => {
  libraryStore.fetchSongs()
  libraryStore.fetchPlaylists()
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style>
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #000;
  border: 1px solid #837dbd;
}

::-webkit-scrollbar-thumb {
  background: #837dbd;
}

::-webkit-scrollbar-thumb:hover {
  background: #d3ceff;
}

::selection {
  background: rgba(211, 206, 255, 0.3);
}

input[type="range"] {
  appearance: none;
  height: 2px;
  background: #5a548d;
}

input[type="range"]::-webkit-slider-thumb {
  appearance: none;
  width: 12px;
  height: 12px;
  background: #d3ceff;
  border-radius: 50%;
  cursor: pointer;
}

input[type="range"]::-moz-range-thumb {
  width: 12px;
  height: 12px;
  background: #d3ceff;
  border-radius: 50%;
  cursor: pointer;
  border: none;
}
</style>