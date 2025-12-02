<template>
  <div class="fixed inset-0 z-50 bg-black/70 flex items-center justify-center">
    <div class="bg-black border border-[#837dbd] p-6 w-96">
      <h3 class="text-[#d3ceff] mb-2">Create Playlist</h3>
      <input v-model="name" placeholder="Playlist name" class="w-full mb-2 p-2 bg-black border border-[#837dbd]" />
      <textarea v-model="description" placeholder="Description (optional)" class="w-full p-2 bg-black border border-[#837dbd]" rows="3"></textarea>
      <div class="mt-4 flex justify-end gap-2">
        <button @click="$emit('close')" class="px-3 py-1 border border-[#837dbd]">Cancel</button>
        <button @click="create" class="px-3 py-1 bg-[#837dbd] text-black">Create</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useLibraryStore } from '../../../stores/library'

const props = defineProps<{[key: string]: never}>()
const emit = defineEmits(['close','created'])

const name = ref('')
const description = ref('')
const library = useLibraryStore()

const create = async () => {
  if (!name.value.trim()) {
    alert('Name required')
    return
  }
  await library.createPlaylist(name.value.trim(), description.value.trim())
  emit('created')
  emit('close')
}
</script>
