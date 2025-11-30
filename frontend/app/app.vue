<template>
  <div class="min-h-screen bg-linear-to-br from-gray-900 to-black text-white">
    <div class="container mx-auto px-4 py-8">
      <header class="mb-8">
        <h1 class="text-4xl font-bold text-center mb-2 bg-linear-to-r from-purple-500 to-pink-500 bg-clip-text text-transparent">
          Nexus Player
        </h1>
        <p class="text-center text-gray-400">The ultimate cross-platform music experience</p>
      </header>

      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Left Sidebar -->
        <div class="lg:col-span-1 space-y-6">
          <SearchComponent />
          <LibraryComponent />
        </div>

        <!-- Main Content -->
        <div class="lg:col-span-3 space-y-6">
          <AudioVisualizer />
          <ControlCenter />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import SearchComponent from './components/SearchComponent.vue'
import LibraryComponent from './components/LibraryComponent.vue'
import AudioVisualizer from './components/AudioVisualizer.vue'
import ControlCenter from './components/ControlCenter.vue'
</script>