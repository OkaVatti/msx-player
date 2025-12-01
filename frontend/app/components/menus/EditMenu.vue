[file name]: menus/EditMenu.vue
[file content begin]
<template>
  <div class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48">
    <button 
      @click="findDuplicates"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      FIND DUPLICATES
    </button>
    <button 
      @click="cleanLibrary"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      CLEAN LIBRARY
    </button>
    <button 
      @click="updateMetadata"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      UPDATE METADATA
    </button>
    <div class="border-t border-[#5a548d] my-1"></div>
    <button 
      @click="refreshLibrary"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      REFRESH LIBRARY
    </button>
  </div>
</template>

<script setup lang="ts">
import { useLibraryStore } from '../../../stores/library'

const libraryStore = useLibraryStore()
const emit = defineEmits(['close'])

const findDuplicates = () => {
  const duplicates = libraryStore.scanForDuplicates()
  if (duplicates.length > 0) {
    alert(`Found ${duplicates.length} duplicate songs. Check the console for details.`)
    console.log('Duplicates found:', duplicates)
  } else {
    alert('No duplicates found!')
  }
  emit('close')
}

const cleanLibrary = () => {
  alert('Clean library functionality would be implemented here')
  emit('close')
}

const updateMetadata = () => {
  alert('Update metadata functionality would be implemented here')
  emit('close')
}

const refreshLibrary = () => {
  libraryStore.fetchSongs()
  libraryStore.fetchPlaylists()
  emit('close')
}
</script>
[file content end]