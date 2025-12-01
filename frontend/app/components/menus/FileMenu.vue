[file name]: menus/FileMenu.vue
[file content begin]
<template>
  <div class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48">
    <button 
      @click="$emit('upload')"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      UPLOAD SONG...
    </button>
    <button 
      @click="scanLibrary"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      SCAN LIBRARY
    </button>
    <div class="border-t border-[#5a548d] my-1"></div>
    <button 
      @click="exportLibrary"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      EXPORT LIBRARY
    </button>
    <button 
      @click="importLibrary"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      IMPORT LIBRARY
    </button>
  </div>
</template>

<script setup lang="ts">
import { useLibraryStore } from '../../../stores/library'

const libraryStore = useLibraryStore()

const scanLibrary = () => {
  alert('Library scan functionality would be implemented here')
  libraryStore.fetchSongs()
  emit('close')
}

const exportLibrary = () => {
  alert('Export functionality would be implemented here')
  emit('close')
}

const importLibrary = () => {
  alert('Import functionality would be implemented here')
  emit('close')
}

const emit = defineEmits(['close', 'upload'])
</script>
[file content end]