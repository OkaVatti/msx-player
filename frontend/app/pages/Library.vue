<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-2xl text-[#d3ceff]">Library</h1>
        <div class="text-xs text-[#837dbd]">Your collection</div>
      </div>
      <div class="flex items-center gap-2">
        <SearchBar />
        <select v-model="library.filters.sortBy" class="bg-black border border-[#837dbd] p-2">
          <option value="title">Title</option>
          <option value="artist">Artist</option>
          <option value="album">Album</option>
          <option value="rating">Rating</option>
          <option value="play_count">Play Count</option>
          <option value="last_played">Last Played</option>
        </select>
      </div>
    </div>

    <div class="bg-[#0b0b0b] border border-[#837dbd] rounded">
      <SongRow v-for="s in library.filteredSongs" :key="s.id" :song="s" />
    </div>
  </div>
</template>

<script setup lang="ts">
import SearchBar from '~/components/library/SearchBar.vue'
import SongRow from '~/components/library/SongRow.vue'
import { useLibraryStore } from '../../stores/library'

const library = useLibraryStore()
library.fetchSongs()
</script>
