<template>
  <div class="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
    <div class="bg-black border border-[#837dbd] p-6 w-96">
      <h3 class="text-[#d3ceff] mb-2">Upload Audio</h3>
      <input type="file" ref="fileInput" accept="audio/*" @change="onFile" />
      <div class="mt-4 flex justify-end gap-2">
        <button @click="$emit('close')" class="px-3 py-1 border border-[#837dbd]">Cancel</button>
        <button @click="upload" :disabled="!file" class="bg-[#837dbd] text-black px-3 py-1">Upload</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useLibraryStore } from '@/stores/library'
const fileInput = ref<HTMLInputElement|null>(null)
const file = ref<File|null>(null)
const lib = useLibraryStore()

const onFile = () => {
  const f = fileInput.value?.files?.[0] ?? null
  if (f) file.value = f
}

const upload = async () => {
  if (!file.value) return
  try {
    await lib.uploadSong(file.value)
    $emit('uploaded')
  } catch (err) {
    alert('Upload failed: '+err)
  }
}
</script>
