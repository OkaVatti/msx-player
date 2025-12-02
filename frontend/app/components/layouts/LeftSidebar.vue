<template>
  <div class="p-4 h-full">
    <div class="text-xs text-[#837dbd] uppercase mb-2">Playlists</div>
    <div class="space-y-2">
      <div v-for="pl in library.playlists" :key="pl.id" class="p-2 border border-[#837dbd] rounded cursor-pointer" @click="select(pl.id)">
        <div class="text-[#d3ceff] font-mono">{{ pl.name }}</div>
        <div class="text-xs text-[#837dbd]">{{ pl.songs.length }} songs</div>
      </div>
    </div>
    <div class="mt-4">
      <button @click="create" class="bg-[#837dbd] text-black px-3 py-1 rounded">New Playlist</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLibraryStore } from '@/stores/library'
const library = useLibraryStore()
const select = (id:number) => library.selectedPlaylistId = id
const create = async () => {
  const name = prompt("Playlist name")
  if (name) {
    await library.createPlaylist(name)
    await library.fetchPlaylists()
  }
}
</script>
