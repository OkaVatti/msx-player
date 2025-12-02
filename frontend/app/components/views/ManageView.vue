<!-- views/ManageView.vue -->
<template>
  <div class="space-y-6">
    <!-- Manage Header -->
    <div class="border border-[#837dbd] p-4 bg-black">
      <h2 class="text-xl mb-2 flex items-center gap-2">
        <span class="text-[#d3ceff]">></span>
        <span class="font-bold">SONG MANAGEMENT</span>
        <span class="text-xs text-[#837dbd] ml-2">[EDIT, ORGANIZE, DELETE]</span>
      </h2>
      
      <div class="text-sm text-[#d3ceff]">
        LIBRARY: {{ libraryStore.songs.length }} SONGS • {{ uniqueArtists }} ARTISTS
      </div>
    </div>

    <!-- Bulk Actions -->
    <div class="border border-[#837dbd] p-4">
      <div class="flex gap-2 mb-4">
        <button @click="selectAll" class="px-4 py-2 border border-[#837dbd] text-[#837dbd] hover:border-[#d3ceff]">
          [SELECT ALL]
        </button>
        <button @click="deselectAll" class="px-4 py-2 border border-[#837dbd] text-[#837dbd] hover:border-[#d3ceff]">
          [DESELECT ALL]
        </button>
        <button @click="exportLibrary" class="px-4 py-2 border border-[#837dbd] text-[#837dbd] hover:border-[#d3ceff]">
          [EXPORT]
        </button>
        <button @click="scanDuplicates" class="px-4 py-2 border border-[#837dbd] text-[#837dbd] hover:border-[#d3ceff]">
          [SCAN DUPLICATES]
        </button>
      </div>

      <div v-if="selectedSongs.length > 0" class="flex gap-2">
        <button @click="rateSelected(5)" class="px-4 py-2 bg-[#d3ceff] text-black hover:bg-[#837dbd]">
          [RATE 5/10]
        </button>
        <button @click="rateSelected(10)" class="px-4 py-2 bg-[#d3ceff] text-black hover:bg-[#837dbd]">
          [RATE 10/10]
        </button>
        <button @click="deleteSelected" class="px-4 py-2 bg-red-500 text-white hover:bg-red-600">
          [DELETE {{ selectedSongs.length }}]
        </button>
      </div>
    </div>

    <!-- Songs Table -->
    <div class="border border-[#837dbd]">
      <div class="grid grid-cols-12 border-b border-[#837dbd] text-xs text-[#837dbd] uppercase tracking-wider">
        <div class="col-span-1 p-3 text-center">
          <input type="checkbox" @change="toggleAll" :checked="isAllSelected">
        </div>
        <div class="col-span-4 p-3">TRACK</div>
        <div class="col-span-3 p-3">ARTIST</div>
        <div class="col-span-2 p-3">ALBUM</div>
        <div class="col-span-1 p-3 text-right">RATING</div>
        <div class="col-span-1 p-3 text-right">ACTIONS</div>
      </div>

      <div v-for="song in libraryStore.songs" :key="song.id"
           :class="['grid grid-cols-12 border-b border-[#222] hover:bg-[#1a1a1a]', 
                   selectedSongs.includes(song.id) ? 'bg-[#1a1a1a]' : '']">
        <div class="col-span-1 p-3 text-center">
          <input type="checkbox" v-model="selectedSongs" :value="song.id">
        </div>
        <div class="col-span-4 p-3">
          <div class="text-[#d3ceff]">{{ song.title }}</div>
          <div class="text-xs text-[#837dbd]">{{ formatDuration(song.duration) }} • {{ song.playCount }} plays</div>
        </div>
        <div class="col-span-3 p-3 text-[#837dbd]">{{ song.artist }}</div>
        <div class="col-span-2 p-3 text-[#837dbd]">{{ song.album }}</div>
        <div class="col-span-1 p-3 text-right">
          <div class="flex justify-end space-x-1">
            <button v-for="star in 5" :key="star"
                    @click="rateSong(song.id, star * 2)"
                    class="text-xs hover:scale-125 transition-transform"
                    :class="star <= Math.floor((song.rating || 0) / 2) ? 'text-yellow-400' : 'text-[#333]'">
              ★
            </button>
          </div>
        </div>
        <div class="col-span-1 p-3 text-right">
          <button @click="editSong(song)" class="text-[#837dbd] hover:text-[#d3ceff] mr-2">[E]</button>
          <button @click="deleteSong(song.id)" class="text-[#837dbd] hover:text-red-400">[X]</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useLibraryStore } from '../../stores/library'

const libraryStore = useLibraryStore()
const selectedSongs = ref<number[]>([])

const uniqueArtists = computed(() => {
  const artists = new Set(libraryStore.songs.map(song => song.artist))
  return artists.size
})

const isAllSelected = computed(() => {
  return selectedSongs.value.length === libraryStore.songs.length && libraryStore.songs.length > 0
})

const formatDuration = (seconds: number) => {
  if (!seconds || seconds === 0) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const selectAll = () => {
  selectedSongs.value = libraryStore.songs.map(song => song.id)
}

const deselectAll = () => {
  selectedSongs.value = []
}

const toggleAll = () => {
  if (isAllSelected.value) {
    deselectAll()
  } else {
    selectAll()
  }
}

const exportLibrary = () => {
  const data = JSON.stringify(libraryStore.songs, null, 2)
  const blob = new Blob([data], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `msx-library-${new Date().toISOString().split('T')[0]}.json`
  a.click()
  URL.revokeObjectURL(url)
}

const scanDuplicates = () => {
  const duplicates = libraryStore.scanForDuplicates()
  if (duplicates.length > 0) {
    alert(`Found ${duplicates.length} duplicate song groups. Check console for details.`)
    console.log('Duplicates:', duplicates)
  } else {
    alert('No duplicates found!')
  }
}

const rateSelected = async (rating: number) => {
  for (const songId of selectedSongs.value) {
    // In a real app, you would call an API to rate the song
    console.log(`Rating song ${songId} as ${rating}/10`)
  }
  alert(`Rated ${selectedSongs.value.length} songs as ${rating}/10`)
}

const deleteSelected = async () => {
  if (!confirm(`Delete ${selectedSongs.value.length} songs? This cannot be undone.`)) {
    return
  }
  
  for (const songId of selectedSongs.value) {
    try {
      await libraryStore.deleteSong(songId)
    } catch (error) {
      console.error(`Failed to delete song ${songId}:`, error)
    }
  }
  
  selectedSongs.value = []
  alert(`Deleted ${selectedSongs.value.length} songs`)
}

const editSong = (song: any) => {
  const newTitle = prompt('Edit title:', song.title)
  if (newTitle) {
    libraryStore.updateSong(song.id, { title: newTitle })
  }
}

const deleteSong = async (songId: number) => {
  if (!confirm('Delete this song?')) {
    return
  }
  
  try {
    await libraryStore.deleteSong(songId)
  } catch (error) {
    console.error('Failed to delete song:', error)
  }
}

const rateSong = async (songId: number, rating: number) => {
  try {
    await libraryStore.updateSong(songId, { rating })
  } catch (error) {
    console.error('Failed to rate song:', error)
  }
}
</script>