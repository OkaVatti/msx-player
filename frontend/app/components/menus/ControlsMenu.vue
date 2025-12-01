[file name]: menus/ControlsMenu.vue
[file content begin]
<template>
  <div class="absolute top-full left-0 mt-1 bg-black border border-[#837dbd] p-2 z-50 min-w-48">
    <button 
      @click="playerStore.shuffleQueue()"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      SHUFFLE QUEUE
    </button>
    <button 
      @click="clearQueue"
      class="block w-full text-left px-3 py-2 text-sm text-[#837dbd] hover:text-[#d3ceff] hover:bg-[#837dbd]/10 transition-colors font-mono"
    >
      CLEAR QUEUE
    </button>
    <div class="border-t border-[#5a548d] my-1"></div>
    <div class="px-3 py-2 text-xs text-[#5a548d] font-mono">REPEAT MODE</div>
    <div class="flex gap-1 px-3 py-2">
      <button 
        @click="setRepeatMode('none')"
        :class="[
          'px-2 py-1 text-xs border transition-colors',
          playerStore.repeatMode === 'none' 
            ? 'border-[#d3ceff] text-[#d3ceff]' 
            : 'border-[#5a548d] text-[#5a548d] hover:border-[#837dbd] hover:text-[#837dbd]'
        ]"
      >
        OFF
      </button>
      <button 
        @click="setRepeatMode('all')"
        :class="[
          'px-2 py-1 text-xs border transition-colors',
          playerStore.repeatMode === 'all' 
            ? 'border-[#d3ceff] text-[#d3ceff]' 
            : 'border-[#5a548d] text-[#5a548d] hover:border-[#837dbd] hover:text-[#837dbd]'
        ]"
      >
        ALL
      </button>
      <button 
        @click="setRepeatMode('one')"
        :class="[
          'px-2 py-1 text-xs border transition-colors',
          playerStore.repeatMode === 'one' 
            ? 'border-[#d3ceff] text-[#d3ceff]' 
            : 'border-[#5a548d] text-[#5a548d] hover:border-[#837dbd] hover:text-[#837dbd]'
        ]"
      >
        ONE
      </button>
    </div>
    <div class="border-t border-[#5a548d] my-1"></div>
    <div class="flex items-center justify-between px-3 py-2">
      <span class="text-xs text-[#837dbd]">SHUFFLE</span>
      <button
        @click="toggleShuffle"
        :class="[
          'w-8 h-4 rounded-full transition-all relative',
          playerStore.shuffle ? 'bg-[#d3ceff]' : 'bg-[#5a548d]'
        ]"
      >
        <div
          :class="[
            'absolute top-0.5 w-3 h-3 rounded-full transition-all',
            playerStore.shuffle ? 'right-0.5 bg-black' : 'left-0.5 bg-white'
          ]"
        ></div>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { usePlayerStore } from '../../../stores/player'

const playerStore = usePlayerStore()
const emit = defineEmits(['close'])

const clearQueue = () => {
  playerStore.setQueue([])
  emit('close')
}

const setRepeatMode = (mode: 'none' | 'all' | 'one') => {
  playerStore.setRepeatMode(mode)
  emit('close')
}

const toggleShuffle = () => {
  playerStore.setShuffle(!playerStore.shuffle)
  emit('close')
}
</script>
[file content end]