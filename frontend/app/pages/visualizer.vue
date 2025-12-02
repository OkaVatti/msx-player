<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-4">
      <div>
        <h1 class="text-2xl text-[#d3ceff]">Visualizer</h1>
        <div class="text-xs text-[#837dbd]">Real-time audio visualization</div>
      </div>
      <div class="flex items-center gap-2">
        <AudioVisualizerSelector />
      </div>
    </div>

    <div class="h-[70vh] border border-[#837dbd]">
      <Visualizer />
    </div>
  </div>
</template>

<script setup lang="ts">
import Visualizer from '~/components/player/Visualizer.vue'
import AudioVisualizerSelector from '@/components/AudioVisualizerSelector.vue'
</script>
