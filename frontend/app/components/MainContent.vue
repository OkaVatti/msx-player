<!-- components/MainContent.vue -->
<template>
  <div class="flex-1 flex flex-col">
    <!-- Content Header -->
    <div class="h-16 border-b border-[#837dbd] px-6 flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-[#d3ceff]">{{ viewTitle }}</h2>
        <div class="text-sm text-[#837dbd]">{{ viewDescription }}</div>
      </div>
      
      <!-- Search and Controls -->
      <div class="flex items-center gap-4">
        <!-- Search Bar -->
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="SEARCH LIBRARY..."
            class="bg-black border border-[#837dbd] text-white px-4 py-2 pl-8 w-64 focus:outline-none focus:border-[#d3ceff]"
            @input="onSearch"
          />
          <div class="absolute left-2 top-1/2 transform -translate-y-1/2 text-[#837dbd]">⌕</div>
        </div>
        
        <!-- View Controls -->
        <div class="flex border border-[#837dbd]">
          <button @click="changeViewMode('list')" :class="viewMode === 'list' ? 'bg-[#1a1a1a]' : ''" 
                  class="px-3 py-1 border-r border-[#837dbd] hover:bg-[#1a1a1a]">
            ☰
          </button>
          <button @click="changeViewMode('grid')" :class="viewMode === 'grid' ? 'bg-[#1a1a1a]' : ''" 
                  class="px-3 py-1 hover:bg-[#1a1a1a]">
            ▦
          </button>
        </div>
      </div>
    </div>

    <!-- Content Body -->
    <div class="flex-1 overflow-auto bg-black p-6">
      <component :is="currentComponent" 
                 :view-mode="viewMode"
                 :search-query="searchQuery"
                 @play-song="playSong" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, defineAsyncComponent, markRaw } from 'vue'
import { usePlayerStore } from '../../stores/player'
import { useLibraryStore } from '../../stores/library'
import type { Song } from '../../types'

interface Props {
  currentView: string
}

const props = defineProps<Props>()
const emit = defineEmits(['play-song'])

const playerStore = usePlayerStore()
const libraryStore = useLibraryStore()

const searchQuery = ref('')
const viewMode = ref<'list' | 'grid'>('list')

// Async components for better performance
const LibraryView = defineAsyncComponent(() => import('~/views/LibraryView.vue'))
const SearchView = defineAsyncComponent(() => import('~/views/SearchView.vue'))
const PlaylistsView = defineAsyncComponent(() => import('~//views/PlaylistsView.vue'))
const UploadView = defineAsyncComponent(() => import('~/views/UploadView.vue'))
const ManageView = defineAsyncComponent(() => import('~/views/ManageView.vue'))
const VisualizerView = defineAsyncComponent(() => import('~//views/VisualizerView.vue'))
const ControlsView = defineAsyncComponent(() => import('~/views/ControlsView.vue'))

const viewTitle = computed(() => {
  const titles: Record<string, string> = {
    library: 'Music Library',
    search: 'Search',
    playlists: 'Playlists',
    upload: 'Upload Music',
    manage: 'Song Management',
    visualizer: 'Audio Visualizer',
    controls: 'Control Center',
    recent: 'Recently Added',
    top: 'Top Rated'
  }
  return titles[props.currentView] || 'Music Library'
})

const viewDescription = computed(() => {
  const descriptions: Record<string, string> = {
    library: `${libraryStore.songs.length} songs in your library`,
    search: 'Find songs in your library',
    playlists: 'Manage your playlists',
    upload: 'Add new music to your library',
    manage: 'Song management and organization',
    visualizer: 'Audio visualization',
    controls: 'Player controls and settings',
    recent: 'Music added in the last 30 days',
    top: 'Your highest rated tracks'
  }
  return descriptions[props.currentView] || 'Browse your music collection'
})

const currentComponent = computed(() => {
  const components: Record<string, any> = {
    library: markRaw(LibraryView),
    search: markRaw(SearchView),
    playlists: markRaw(PlaylistsView),
    upload: markRaw(UploadView),
    manage: markRaw(ManageView),
    visualizer: markRaw(VisualizerView),
    controls: markRaw(ControlsView),
    recent: markRaw(LibraryView),
    top: markRaw(LibraryView)
  }
  return components[props.currentView] || markRaw(LibraryView)
})

const onSearch = () => {
  // Debounced search
  libraryStore.searchQuery = searchQuery.value
  if (props.currentView === 'library' || props.currentView === 'search') {
    libraryStore.fetchSongs()
  }
}

const changeViewMode = (mode: 'list' | 'grid') => {
  viewMode.value = mode
}

const playSong = (song: Song) => {
  emit('play-song', song)
}
</script>