<!-- views/ControlsView.vue -->
<template>
  <div class="space-y-6">
    <!-- Controls Header -->
    <div class="border border-[#837dbd] p-4 bg-black">
      <h2 class="text-xl mb-2 flex items-center gap-2">
        <span class="text-[#d3ceff]">></span>
        <span class="font-bold">CONTROL CENTER</span>
        <span class="text-xs text-[#837dbd] ml-2">[PLAYBACK SETTINGS]</span>
      </h2>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Playback Controls -->
      <div class="border border-[#837dbd] p-6">
        <h3 class="text-lg mb-4 text-[#d3ceff]">PLAYBACK CONTROLS</h3>
        
        <div class="grid grid-cols-3 gap-4 mb-6">
          <button @click="playerStore.previousSong" 
                  :disabled="!playerStore.currentSong || playerStore.queue.length <= 1"
                  class="p-4 border border-[#837dbd] hover:border-[#d3ceff] hover:text-[#d3ceff] disabled:opacity-50">
            ⏮
          </button>
          
          <button v-if="!playerStore.isPlaying" 
                  @click="resumePlayer"
                  :disabled="!playerStore.currentSong"
                  class="p-4 bg-[#d3ceff] text-black hover:bg-[#837dbd] disabled:opacity-50">
            ▶
          </button>
          <button v-else
                  @click="playerStore.pause"
                  class="p-4 bg-[#837dbd] text-black hover:bg-[#d3ceff]">
            ⏸
          </button>
          
          <button @click="playerStore.nextSong"
                  :disabled="!playerStore.currentSong || playerStore.queue.length <= 1"
                  class="p-4 border border-[#837dbd] hover:border-[#d3ceff] hover:text-[#d3ceff] disabled:opacity-50">
            ⏭
          </button>
        </div>

        <!-- Playback Speed -->
        <div class="space-y-3">
          <div class="text-sm text-[#837dbd]">PLAYBACK SPEED</div>
          <select v-model="playerStore.speed" @change="playerStore.setSpeed(playerStore.speed)"
                  class="w-full bg-black border border-[#837dbd] text-white p-2 focus:outline-none focus:border-[#d3ceff]">
            <option value="0.5">0.5x (Slow)</option>
            <option value="0.75">0.75x</option>
            <option value="1">1.0x (Normal)</option>
            <option value="1.25">1.25x</option>
            <option value="1.5">1.5x</option>
            <option value="2">2.0x (Fast)</option>
          </select>
        </div>
      </div>

      <!-- Audio Settings -->
      <div class="border border-[#837dbd] p-6">
        <h3 class="text-lg mb-4 text-[#d3ceff]">AUDIO SETTINGS</h3>
        
        <!-- Volume -->
        <div class="space-y-3 mb-6">
          <div class="text-sm text-[#837dbd]">VOLUME: {{ Math.round(playerStore.volume * 100) }}%</div>
          <input type="range" v-model="playerStore.volume" min="0" max="1" step="0.01"
                 @input="playerStore.setVolume(playerStore.volume)"
                 class="w-full accent-[#d3ceff] bg-[#333] h-1">
        </div>

        <!-- Repeat Mode -->
        <div class="space-y-3 mb-6">
          <div class="text-sm text-[#837dbd]">REPEAT MODE</div>
          <div class="flex gap-2">
            <button @click="playerStore.setRepeatMode('none')"
                    :class="playerStore.repeatMode === 'none' ? 'bg-[#d3ceff] text-black' : 'border border-[#837dbd]'"
                    class="px-3 py-2">
              NONE
            </button>
            <button @click="playerStore.setRepeatMode('all')"
                    :class="playerStore.repeatMode === 'all' ? 'bg-[#d3ceff] text-black' : 'border border-[#837dbd]'"
                    class="px-3 py-2">
              ALL
            </button>
            <button @click="playerStore.setRepeatMode('one')"
                    :class="playerStore.repeatMode === 'one' ? 'bg-[#d3ceff] text-black' : 'border border-[#837dbd]'"
                    class="px-3 py-2">
              ONE
            </button>
          </div>
        </div>

        <!-- Shuffle -->
        <div class="space-y-3">
          <div class="text-sm text-[#837dbd]">SHUFFLE</div>
          <button @click="playerStore.setShuffle(!playerStore.shuffle)"
                  :class="playerStore.shuffle ? 'bg-[#d3ceff] text-black' : 'border border-[#837dbd]'"
                  class="px-6 py-2">
            {{ playerStore.shuffle ? 'ON' : 'OFF' }}
          </button>
        </div>
      </div>

      <!-- Queue Management -->
      <div class="border border-[#837dbd] p-6 lg:col-span-2">
        <h3 class="text-lg mb-4 text-[#d3ceff]">QUEUE MANAGEMENT</h3>
        
        <div class="space-y-4">
          <div class="text-sm text-[#837dbd]">CURRENT QUEUE: {{ playerStore.queue.length }} songs</div>
          
          <div class="grid grid-cols-2 gap-4">
            <button @click="playerStore.shuffleQueue" 
                    class="px-4 py-2 border border-[#837dbd] text-[#837dbd] hover:border-[#d3ceff] hover:text-[#d3ceff]">
              SHUFFLE QUEUE
            </button>
            <button @click="clearQueue"
                    class="px-4 py-2 border border-red-500 text-red-500 hover:bg-red-500 hover:text-white">
              CLEAR QUEUE
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { usePlayerStore } from '../../stores/player'

const playerStore = usePlayerStore()

const resumePlayer = () => {
  if (playerStore.currentSong) {
    playerStore.playSong(playerStore.currentSong)
  }
}

const clearQueue = () => {
  if (confirm('Clear the current queue?')) {
    playerStore.setQueue([], -1)
  }
}
</script>