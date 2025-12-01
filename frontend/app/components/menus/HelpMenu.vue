[file name]: menus/HelpMenu.vue
[file content begin]
<template>
  <div class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48">
    <button 
      @click="showKeyboardShortcuts"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      KEYBOARD SHORTCUTS
    </button>
    <button 
      @click="showAbout"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      ABOUT
    </button>
    <div class="border-t border-[#5a548d] my-1"></div>
    <button 
      @click="checkForUpdates"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      CHECK FOR UPDATES
    </button>
    <button 
      @click="viewDocumentation"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      DOCUMENTATION
    </button>
  </div>
</template>

<script setup lang="ts">
const emit = defineEmits(['close'])

const showKeyboardShortcuts = () => {
  alert(`
Keyboard Shortcuts:
Space - Play/Pause
← → - Seek forward/backward
↑ ↓ - Volume up/down
N - Next song
P - Previous song
S - Shuffle
R - Repeat mode
F - Toggle fullscreen
ESC - Close dialogs
  `)
  emit('close')
}

const showAbout = () => {
  alert(`
MSX Player v1.0.0
A retro-futuristic music player
Built with Vue.js, TypeScript, and Go
© 2024 MSX Player Team
  `)
  emit('close')
}

const checkForUpdates = () => {
  alert('Checking for updates...\n(This would check your package.json or a remote API)')
  emit('close')
}

const viewDocumentation = () => {
  alert('Documentation would open in a browser')
  emit('close')
}
</script>
[file content end]