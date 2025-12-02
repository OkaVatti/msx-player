<template>
  <div class="min-h-screen flex flex-col bg-black text-white">
    <div class="h-16 flex items-center px-4 border-b border-[#837dbd]">
      <div class="flex items-center gap-6">
        <div class="text-[#d3ceff] font-bold text-lg">MSX PLAYER</div>
        <nav class="text-sm text-[#837dbd] font-mono">
          <button @click="setView('library')" :class="btnClass('library')">Library</button>
          <button @click="setView('playlists')" :class="btnClass('playlists')">Playlists</button>
          <button @click="setView('visualizer')" :class="btnClass('visualizer')">Visualizer</button>
        </nav>
      </div>
      <div class="ml-auto flex items-center gap-4">
        <UploadModal v-if="showUpload" @close="showUpload=false" @uploaded="onUploaded" />
        <button @click="showUpload=true" class="bg-[#837dbd] text-black px-3 py-1 rounded text-sm">Upload</button>
      </div>
    </div>

    <div class="flex flex-1 overflow-hidden">
      <LeftSidebar class="w-80 border-r border-[#837dbd]" />
      <main class="flex-1 overflow-auto">
        <MainView :view="view" />
      </main>
      <RightSidebar class="w-96 border-l border-[#837dbd]" />
    </div>

    <NowPlayingBar />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import LeftSidebar from '~/components/layouts/LeftSidebar.vue'
import RightSidebar from '~/components/layouts/RightSidebar.vue'
import MainView from '~/components/library/MainView.vue'
import NowPlayingBar from '~/components/layouts/NowPlayingBar.vue'
import UploadModal from '@/components/UploadModal.vue'
import { useLibraryStore } from '@/stores/library'

const view = ref<'library' | 'playlists' | 'visualizer'>('library')
const showUpload = ref(false)

const store = useLibraryStore()
store.fetchSongs()
store.fetchPlaylists()

const setView = (v: typeof view.value) => {
  view.value = v
}
const onUploaded = () => {
  showUpload.value = false
  store.fetchSongs()
}

const btnClass = (v: string) => {
  return [
    'px-2 py-1 rounded font-mono',
    view.value === v ? 'bg-[#837dbd] text-black' : 'text-[#837dbd] hover:bg-[#837dbd] hover:text-black'
  ]
}
</script>

<style>
/* small global accent touches kept here */
</style>
