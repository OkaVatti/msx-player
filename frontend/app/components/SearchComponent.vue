<template>
  <div class="bg-gray-800 rounded-lg p-6">
    <h2 class="text-xl font-semibold mb-4">Search</h2>
    
    <div class="space-y-4">
      <input
        v-model="libraryStore.searchQuery"
        type="text"
        placeholder="Search songs, artists, albums..."
        class="w-full px-4 py-2 bg-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
        @input="debouncedSearch"
      >

      <div class="space-y-2">
        <select
          v-model="libraryStore.filters.genre"
          class="w-full px-4 py-2 bg-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
        >
          <option value="">All Genres</option>
          <option v-for="genre in libraryStore.genres" :key="genre" :value="genre">
            {{ genre }}
          </option>
        </select>

        <select
          v-model="libraryStore.filters.rating"
          class="w-full px-4 py-2 bg-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
        >
          <option value="0">All Ratings</option>
          <option value="7">7+ Stars</option>
          <option value="8">8+ Stars</option>
          <option value="9">9+ Stars</option>
          <option value="10">10 Stars</option>
        </select>

        <select
          v-model="libraryStore.filters.sortBy"
          class="w-full px-4 py-2 bg-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
        >
          <option value="title">Title</option>
          <option value="artist">Artist</option>
          <option value="album">Album</option>
          <option value="rating">Rating</option>
          <option value="play_count">Play Count</option>
          <option value="last_played">Last Played</option>
        </select>
      </div>

      <div class="border-t border-gray-700 pt-4">
        <h3 class="text-lg font-medium mb-3">Upload Music</h3>
        <input
          type="file"
          accept=".mp3"
          @change="handleFileUpload"
          class="w-full text-sm text-gray-400 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-purple-500 file:text-white hover:file:bg-purple-600"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLibraryStore } from '../../stores/library'
import { debounce } from '../../utils/debounce'

const libraryStore = useLibraryStore()

const debouncedSearch = debounce(() => {
  libraryStore.fetchSongs()
}, 300)

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    libraryStore.uploadSong(target.files[0])
  }
}
</script>