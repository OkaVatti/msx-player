<template>
  <div class="flex items-center gap-4 text-xs">
    <div 
      :class="[
        'flex items-center gap-2',
        isConnected ? 'text-green-400' : 'text-red-400'
      ]"
    >
      <div 
        :class="[
          'w-2 h-2 rounded-full animate-pulse',
          isConnected ? 'bg-green-400' : 'bg-red-400'
        ]"
      ></div>
      <span>{{ isConnected ? '[CONNECTED]' : '[DISCONNECTED]' }}</span>
    </div>
    <div class="text-[#837dbd]">
      {{ songCount }} SONGS
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { usePlayerStore } from '../../stores/player'
import { useLibraryStore } from '../../stores/library'

const playerStore = usePlayerStore()
const libraryStore = useLibraryStore()

const isConnected = computed(() => playerStore.isConnected)
const songCount = computed(() => libraryStore.songs.length)
</script>