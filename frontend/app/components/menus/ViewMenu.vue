[file name]: menus/ViewMenu.vue
[file content begin]
<template>
  <div class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48">
    <button 
      @click="toggleDarkMode"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      {{ darkMode ? 'LIGHT MODE' : 'DARK MODE' }}
    </button>
    <button 
      @click="toggleCompactView"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      {{ compactView ? 'EXPAND VIEW' : 'COMPACT VIEW' }}
    </button>
    <div class="border-t border-[#5a548d] my-1"></div>
    <button 
      @click="changeView('library')"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      LIBRARY VIEW
    </button>
    <button 
      @click="changeView('playlists')"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      PLAYLISTS VIEW
    </button>
    <button 
      @click="changeView('visualizer')"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      VISUALIZER VIEW
    </button>
    <button 
      @click="changeView('search')"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      SEARCH VIEW
    </button>
    <button 
      @click="changeView('settings')"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      SETTINGS VIEW
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const darkMode = ref(true)
const compactView = ref(false)
const emit = defineEmits(['close', 'change-view'])

const toggleDarkMode = () => {
  darkMode.value = !darkMode.value
  if (darkMode.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
  emit('close')
}

const toggleCompactView = () => {
  compactView.value = !compactView.value
  // This would adjust CSS classes or store preference
  emit('close')
}

const changeView = (view: string) => {
  emit('change-view', view)
  emit('close')
}
</script>
[file content end]