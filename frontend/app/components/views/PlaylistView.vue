<template>
  <div v-if="playlist" class="p-4">
    <div class="flex items-center gap-4 mb-4">
      <div class="w-24 h-24 bg-linear-to-br from-[#5a548d] to-[#837dbd] flex items-center justify-center text-4xl">♪</div>
      <div>
        <div class="text-2xl text-[#d3ceff]">{{ playlist.name }}</div>
        <div class="text-sm text-[#837dbd]">{{ playlist.songs.length }} songs</div>
      </div>
    </div>

    <div>
      <SongRow v-for="s in playlist.songs" :key="s.id" :song="s" />
    </div>
  </div>
  <div v-else class="p-4 text-[#837dbd]">Playlist not selected</div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useLibraryStore } from '../../../stores/library'
import SongRow from '../library/SongRow.vue'

const library = useLibraryStore()
const playlist = computed(() => {
  const id = library.selectedPlaylistId
  return library.playlists.find(p => p.id === id) ?? null
})
</script>
