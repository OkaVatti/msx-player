<template>
  <div class="flex gap-4 relative">
    <!-- File Menu -->
    <div class="relative">
      <button 
        @click.stop="toggleFileMenu"
        :class="[
          'text-[#837dbd] hover:text-[#d3ceff] transition-colors',
          showFileMenu ? 'text-[#d3ceff]' : ''
        ]"
      >
        FILE
      </button>
      <FileMenu 
        v-if="showFileMenu" 
        @upload="upload" 
        @close="closeAllMenus"
        @click-outside="showFileMenu = false"
      />
    </div>
    
    <!-- Edit Menu -->
    <div class="relative">
      <button 
        @click.stop="toggleEditMenu"
        :class="[
          'text-[#837dbd] hover:text-[#d3ceff] transition-colors',
          showEditMenu ? 'text-[#d3ceff]' : ''
        ]"
      >
        EDIT
      </button>
      <EditMenu 
        v-if="showEditMenu" 
        @close="closeAllMenus"
        @click-outside="showEditMenu = false"
      />
    </div>
    
    <!-- View Menu -->
    <div class="relative">
      <button 
        @click.stop="toggleViewMenu"
        :class="[
          'text-[#837dbd] hover:text-[#d3ceff] transition-colors',
          showViewMenu ? 'text-[#d3ceff]' : ''
        ]"
      >
        VIEW
      </button>
      <ViewMenu 
        v-if="showViewMenu" 
        @change-view="changeView"
        @close="closeAllMenus"
        @click-outside="showViewMenu = false"
      />
    </div>
    
    <!-- Controls Menu -->
    <div class="relative">
      <button 
        @click.stop="toggleControlsMenu"
        :class="[
          'text-[#837dbd] hover:text-[#d3ceff] transition-colors',
          showControlsMenu ? 'text-[#d3ceff]' : ''
        ]"
      >
        CONTROLS
      </button>
      <ControlsMenu 
        v-if="showControlsMenu" 
        @close="closeAllMenus"
        @click-outside="showControlsMenu = false"
      />
    </div>
    
    <!-- Help Menu -->
    <div class="relative">
      <button 
        @click.stop="toggleHelpMenu"
        :class="[
          'text-[#837dbd] hover:text-[#d3ceff] transition-colors',
          showHelpMenu ? 'text-[#d3ceff]' : ''
        ]"
      >
        HELP
      </button>
      <HelpMenu 
        v-if="showHelpMenu" 
        @close="closeAllMenus"
        @click-outside="showHelpMenu = false"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import FileMenu from '~/components/menus/FileMenu.vue'
import EditMenu from '~/components/menus/EditMenu.vue'
import ViewMenu from '~/components/menus/ViewMenu.vue'
import ControlsMenu from '~/components/menus/ControlsMenu.vue'
import HelpMenu from '~/components/menus/HelpMenu.vue'

const emit = defineEmits(['upload', 'change-view'])

const showFileMenu = ref(false)
const showEditMenu = ref(false)
const showViewMenu = ref(false)
const showControlsMenu = ref(false)
const showHelpMenu = ref(false)

const toggleFileMenu = () => {
  showFileMenu.value = !showFileMenu.value
  if (showFileMenu.value) {
    showEditMenu.value = false
    showViewMenu.value = false
    showControlsMenu.value = false
    showHelpMenu.value = false
  }
}

const toggleEditMenu = () => {
  showEditMenu.value = !showEditMenu.value
  if (showEditMenu.value) {
    showFileMenu.value = false
    showViewMenu.value = false
    showControlsMenu.value = false
    showHelpMenu.value = false
  }
}

const toggleViewMenu = () => {
  showViewMenu.value = !showViewMenu.value
  if (showViewMenu.value) {
    showFileMenu.value = false
    showEditMenu.value = false
    showControlsMenu.value = false
    showHelpMenu.value = false
  }
}

const toggleControlsMenu = () => {
  showControlsMenu.value = !showControlsMenu.value
  if (showControlsMenu.value) {
    showFileMenu.value = false
    showEditMenu.value = false
    showViewMenu.value = false
    showHelpMenu.value = false
  }
}

const toggleHelpMenu = () => {
  showHelpMenu.value = !showHelpMenu.value
  if (showHelpMenu.value) {
    showFileMenu.value = false
    showEditMenu.value = false
    showViewMenu.value = false
    showControlsMenu.value = false
  }
}

const closeAllMenus = () => {
  showFileMenu.value = false
  showEditMenu.value = false
  showViewMenu.value = false
  showControlsMenu.value = false
  showHelpMenu.value = false
}

const upload = () => {
  emit('upload')
  closeAllMenus()
}

const changeView = (view: string) => {
  emit('change-view', view)
  closeAllMenus()
}

// Close menus when clicking outside
const handleClickOutside = () => {
  closeAllMenus()
}

// Listen for escape key
document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') {
    closeAllMenus()
  }
})
</script>