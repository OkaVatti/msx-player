<template>
  <div class="flex items-center gap-2">
    <label class="text-xs text-[#837dbd]">Visualizer</label>
    <select v-model="mode" class="bg-black border border-[#837dbd] p-1">
      <option value="bars">Bars</option>
      <option value="wave">Wave</option>
      <option value="particles">Particles</option>
      <option value="circular">Circular</option>
      <option value="matrix">Matrix</option>
    </select>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const mode = ref('bars')
</script>
