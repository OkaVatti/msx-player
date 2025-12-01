<template>
  <div class="flex-1 overflow-hidden flex flex-col">
    <!-- Dynamic Content Based on View -->
    <component 
      :is="currentComponent"
      @play-song="$emit('play-song', $event)"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent } from 'vue'
import LibraryView from '~/views/LibraryView.vue'
import PlaylistsView from '~/views/PlaylistsView.vue'
import VisualizerView from '~/views/VisualizerView.vue'
import SettingsView from '~/views/SettingsView.vue'
import SearchView from '~/views/SearchView.vue'

interface Props {
  currentView: string
}

const props = defineProps<Props>()
defineEmits(['play-song'])

// Async components for better performance
const UploadView = defineAsyncComponent(() => import('~/views/UploadView.vue'))
const ManageView = defineAsyncComponent(() => import('~/views/ManageView.vue'))
const ControlCenterView = defineAsyncComponent(() => import('~//views/ControlCenterView.vue'))

const currentComponent = computed(() => {
  const components: Record<string, any> = {
    library: LibraryView,
    playlists: PlaylistsView,
    visualizer: VisualizerView,
    settings: SettingsView,
    search: SearchView,
    upload: UploadView,
    manage: ManageView,
    controls: ControlCenterView,
    recent: LibraryView,
    top: LibraryView
  }
  return components[props.currentView] || LibraryView
})
</script>