[file name]: views/SettingsView.vue
[file content begin]
<template>
  <div class="h-full flex flex-col p-4 space-y-4">
    <!-- Settings Header -->
    <div class="border border-[#837dbd] p-4 bg-black">
      <h2 class="text-lg text-[#d3ceff] font-bold">SETTINGS</h2>
      <p class="text-[#837dbd] text-sm mt-1">Configure your MSX Player</p>
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Player Settings -->
      <div class="border border-[#837dbd] p-4 bg-black">
        <h3 class="text-[#d3ceff] text-sm font-bold mb-4">PLAYER SETTINGS</h3>
        <div class="space-y-4">
          <div>
            <label class="text-[#837dbd] text-xs block mb-2">VOLUME</label>
            <div class="flex items-center gap-3">
              <input
                type="range"
                v-model="playerStore.volume"
                min="0"
                max="1"
                step="0.01"
                @change="playerStore.setVolume(playerStore.volume)"
                class="flex-1"
              />
              <span class="text-[#5a548d] text-xs font-mono w-12">
                {{ Math.round(playerStore.volume * 100) }}%
              </span>
            </div>
          </div>
          
          <div>
            <label class="text-[#837dbd] text-xs block mb-2">PLAYBACK SPEED</label>
            <div class="flex items-center gap-3">
              <input
                type="range"
                v-model="playerStore.speed"
                min="0.5"
                max="2"
                step="0.1"
                @change="playerStore.setSpeed(playerStore.speed)"
                class="flex-1"
              />
              <span class="text-[#5a548d] text-xs font-mono w-12">
                {{ playerStore.speed.toFixed(1) }}x
              </span>
            </div>
          </div>
          
          <div class="flex items-center justify-between">
            <span class="text-[#837dbd] text-xs">SHUFFLE</span>
            <button
              @click="playerStore.setShuffle(!playerStore.shuffle)"
              :class="[
                'w-8 h-4 rounded-full transition-all relative',
                playerStore.shuffle ? 'bg-[#d3ceff]' : 'bg-[#5a548d]'
              ]"
            >
              <div
                :class="[
                  'absolute top-0.5 w-3 h-3 rounded-full transition-all',
                  playerStore.shuffle ? 'right-0.5 bg-black' : 'left-0.5 bg-white'
                ]"
              ></div>
            </button>
          </div>
        </div>
      </div>
      
      <!-- Library Settings -->
      <div class="border border-[#837dbd] p-4 bg-black">
        <h3 class="text-[#d3ceff] text-sm font-bold mb-4">LIBRARY SETTINGS</h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <span class="text-[#837dbd] text-xs">AUTO-SCAN FOR NEW FILES</span>
            <button
              @click="toggleAutoScan"
              :class="[
                'w-8 h-4 rounded-full transition-all relative',
                autoScan ? 'bg-[#d3ceff]' : 'bg-[#5a548d]'
              ]"
            >
              <div
                :class="[
                  'absolute top-0.5 w-3 h-3 rounded-full transition-all',
                  autoScan ? 'right-0.5 bg-black' : 'left-0.5 bg-white'
                ]"
              ></div>
            </button>
          </div>
          
          <div>
            <label class="text-[#837dbd] text-xs block mb-2">DEFAULT SORT ORDER</label>
            <select
              v-model="libraryStore.filters.sortBy"
              class="w-full px-3 py-1 bg-black border border-[#837dbd] text-[#837dbd] text-xs focus:outline-none focus:border-[#d3ceff] font-mono"
              @change="saveSettings"
            >
              <option value="title">Title</option>
              <option value="artist">Artist</option>
              <option value="album">Album</option>
              <option value="rating">Rating</option>
            </select>
          </div>
          
          <div>
            <label class="text-[#837dbd] text-xs block mb-2">LIBRARY PATH</label>
            <div class="flex gap-2">
              <input
                type="text"
                :value="libraryPath"
                readonly
                class="flex-1 px-3 py-1 bg-black border border-[#5a548d] text-[#837dbd] text-xs font-mono"
              />
              <button
                @click="changeLibraryPath"
                class="px-3 py-1 border border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black transition-all text-xs font-mono"
              >
                BROWSE
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Display Settings -->
      <div class="border border-[#837dbd] p-4 bg-black">
        <h3 class="text-[#d3ceff] text-sm font-bold mb-4">DISPLAY SETTINGS</h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <span class="text-[#837dbd] text-xs">DARK MODE</span>
            <button
              @click="toggleDarkMode"
              :class="[
                'w-8 h-4 rounded-full transition-all relative',
                darkMode ? 'bg-[#d3ceff]' : 'bg-[#5a548d]'
              ]"
            >
              <div
                :class="[
                  'absolute top-0.5 w-3 h-3 rounded-full transition-all',
                  darkMode ? 'right-0.5 bg-black' : 'left-0.5 bg-white'
                ]"
              ></div>
            </button>
          </div>
          
          <div class="flex items-center justify-between">
            <span class="text-[#837dbd] text-xs">COMPACT VIEW</span>
            <button
              @click="toggleCompactView"
              :class="[
                'w-8 h-4 rounded-full transition-all relative',
                compactView ? 'bg-[#d3ceff]' : 'bg-[#5a548d]'
              ]"
            >
              <div
                :class="[
                  'absolute top-0.5 w-3 h-3 rounded-full transition-all',
                  compactView ? 'right-0.5 bg-black' : 'left-0.5 bg-white'
                ]"
              ></div>
            </button>
          </div>
          
          <div class="flex items-center justify-between">
            <span class="text-[#837dbd] text-xs">SCANLINES EFFECT</span>
            <button
              @click="toggleScanlines"
              :class="[
                'w-8 h-4 rounded-full transition-all relative',
                scanlines ? 'bg-[#d3ceff]' : 'bg-[#5a548d]'
              ]"
            >
              <div
                :class="[
                  'absolute top-0.5 w-3 h-3 rounded-full transition-all',
                  scanlines ? 'right-0.5 bg-black' : 'left-0.5 bg-white'
                ]"
              ></div>
            </button>
          </div>
        </div>
      </div>
      
      <!-- Advanced Settings -->
      <div class="border border-[#837dbd] p-4 bg-black">
        <h3 class="text-[#d3ceff] text-sm font-bold mb-4">ADVANCED SETTINGS</h3>
        <div class="space-y-4">
          <button
            @click="clearCache"
            class="w-full px-3 py-2 border border-red-400 text-red-400 hover:bg-red-400 hover:text-black transition-all text-xs font-mono"
          >
            CLEAR CACHE
          </button>
          
          <button
            @click="resetSettings"
            class="w-full px-3 py-2 border border-yellow-400 text-yellow-400 hover:bg-yellow-400 hover:text-black transition-all text-xs font-mono"
          >
            RESET TO DEFAULTS
          </button>
          
          <button
            @click="exportSettings"
            class="w-full px-3 py-2 border border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black transition-all text-xs font-mono"
          >
            EXPORT SETTINGS
          </button>
          
          <button
            @click="importSettings"
            class="w-full px-3 py-2 border border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black transition-all text-xs font-mono"
          >
            IMPORT SETTINGS
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { usePlayerStore } from '../../stores/player'
import { useLibraryStore } from '../../stores/library'

const playerStore = usePlayerStore()
const libraryStore = useLibraryStore()

const darkMode = ref(true)
const compactView = ref(false)
const scanlines = ref(true)
const autoScan = ref(true)
const libraryPath = ref('./music')

onMounted(() => {
  // Load settings from localStorage
  const savedAutoScan = localStorage.getItem('autoScan')
  if (savedAutoScan !== null) {
    autoScan.value = savedAutoScan === 'true'
  }
  
  const savedDarkMode = localStorage.getItem('darkMode')
  if (savedDarkMode !== null) {
    darkMode.value = savedDarkMode === 'true'
  }
  
  const savedScanlines = localStorage.getItem('scanlines')
  if (savedScanlines !== null) {
    scanlines.value = savedScanlines === 'true'
  }
})

const toggleAutoScan = () => {
  autoScan.value = !autoScan.value
  localStorage.setItem('autoScan', autoScan.value.toString())
}

const toggleDarkMode = () => {
  darkMode.value = !darkMode.value
  localStorage.setItem('darkMode', darkMode.value.toString())
  if (darkMode.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}

const toggleCompactView = () => {
  compactView.value = !compactView.value
  localStorage.setItem('compactView', compactView.value.toString())
}

const toggleScanlines = () => {
  scanlines.value = !scanlines.value
  localStorage.setItem('scanlines', scanlines.value.toString())
}

const saveSettings = () => {
  localStorage.setItem('sortBy', libraryStore.filters.sortBy)
}

const changeLibraryPath = () => {
  alert('This would open a directory browser')
}

const clearCache = () => {
  if (confirm('Are you sure you want to clear the cache?')) {
    localStorage.removeItem('libraryCache')
    alert('Cache cleared successfully')
  }
}

const resetSettings = () => {
  if (confirm('Are you sure you want to reset all settings to defaults?')) {
    localStorage.clear()
    window.location.reload()
  }
}

const exportSettings = () => {
  alert('Settings export functionality would be implemented here')
}

const importSettings = () => {
  alert('Settings import functionality would be implemented here')
}
</script>
[file content end]