<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-2xl text-[#d3ceff]">Playlists</h1>
        <div class="text-xs text-[#837dbd]">Create and manage playlists</div>
      </div>
      <div>
        <button @click="showCreate = true" class="bg-[#837dbd] text-black px-3 py-1 rounded">New Playlist</button>
      </div>
    </div>

    <div class="grid grid-cols-3 gap-4">
      <div v-for="pl in library.playlists" :key="pl.id" class="p-4 border border-[#837dbd] rounded">
        <div class="text-[#d3ceff] font-medium">{{ pl.name }}</div>
        <div class="text-xs text-[#837dbd]">{{ pl.songs.length }} songs</div>
        <div class="mt-2 flex gap-2">
          <button @click="open(pl.id)" class="text-[#837dbd]">Open</button>
        </div>
      </div>
    </div>

    <CreatePlaylistModal v-if="showCreate" @close="showCreate=false" @created="onCreated" />
    <PlaylistView v-if="activePlaylist" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import CreatePlaylistModal from '../components/modals/CreatePlaylistModal.vue'
import PlaylistView from '../components/views/PlaylistView.vue'
import { useLibraryStore } from '../../stores/library'

const library = useLibraryStore()
library.fetchPlaylists()
const showCreate = ref(false)

const activePlaylist = computed(() => {
  return library.selectedPlaylistId ? library.playlists.find(p => p.id === library.selectedPlaylistId) : null
})

const onCreated = async () => {
  await library.fetchPlaylists()
}

const open = (id:number) => {
  library.selectedPlaylistId = id
}
</script>
