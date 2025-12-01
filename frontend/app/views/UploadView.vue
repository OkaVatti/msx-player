<!-- views/UploadView.vue -->
<template>
  <div class="space-y-6">
    <!-- Upload Header -->
    <div class="border border-[#837dbd] p-4 bg-black">
      <h2 class="text-xl mb-2 flex items-center gap-2">
        <span class="text-[#d3ceff]">></span>
        <span class="font-bold">UPLOAD SYSTEM</span>
        <span class="text-xs text-[#837dbd] ml-2">[DRAG & DROP SUPPORTED]</span>
      </h2>
      
      <div class="text-sm text-[#d3ceff]">
        STATUS: <span :class="uploadStatusClass">{{ uploadStatus }}</span>
      </div>
    </div>

    <!-- Upload Area -->
    <div class="border-2 border-dashed border-[#837dbd] p-8 text-center bg-black transition-all hover:border-[#d3ceff]"
         @dragover.prevent="isDragging = true"
         @dragleave="isDragging = false"
         @drop.prevent="handleDrop"
         :class="{ 'border-[#d3ceff]': isDragging }">
      <div class="space-y-3">
        <div class="text-4xl text-[#837dbd]">📁</div>
        <h3 class="text-lg text-[#d3ceff]">DROP MP3 FILES HERE</h3>
        <p class="text-sm text-[#837dbd]">OR CLICK TO BROWSE</p>
        <input ref="fileInput" type="file" multiple accept=".mp3,audio/mpeg" @change="handleFileSelect" class="hidden">
        <button @click="triggerFileInput" 
                class="px-6 py-3 border border-[#837dbd] text-[#837dbd] hover:border-[#d3ceff] hover:text-[#d3ceff] transition-colors">
          [SELECT FILES]
        </button>
      </div>
    </div>

    <!-- Upload Queue -->
    <div v-if="uploadQueue.length > 0" class="border border-[#837dbd] p-4">
      <h3 class="text-lg mb-3 text-[#d3ceff]">UPLOAD QUEUE ({{ uploadQueue.length }})</h3>
      
      <div class="space-y-2 max-h-64 overflow-y-auto">
        <div v-for="(item, index) in uploadQueue" :key="index" 
             class="border border-[#837dbd] p-3">
          <div class="flex justify-between items-center">
            <div class="flex-1 min-w-0">
              <div class="text-[#d3ceff] truncate">{{ item.file.name }}</div>
              <div class="text-xs text-[#837dbd]">
                {{ formatFileSize(item.file.size) }}
                <span v-if="item.status === 'uploading'" class="text-yellow-400">
                  • {{ item.progress }}%
                </span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <div v-if="item.status === 'uploading'" class="w-20 bg-[#333]">
                <div class="h-2 bg-[#d3ceff]" :style="{ width: `${item.progress}%` }"></div>
              </div>
              <span v-if="item.status === 'pending'" class="text-[#837dbd]">[ ]</span>
              <span v-else-if="item.status === 'uploading'" class="text-yellow-400 animate-pulse">[~]</span>
              <span v-else-if="item.status === 'completed'" class="text-green-400">[✓]</span>
              <span v-else-if="item.status === 'error'" class="text-red-400">[!]</span>
              <button @click="removeFromQueue(index)" class="text-[#837dbd] hover:text-red-400">[X]</button>
            </div>
          </div>
        </div>
      </div>

      <div class="flex gap-2 mt-4">
        <button @click="startUpload" :disabled="!hasPendingFiles || isUploading"
                class="px-4 py-2 bg-[#d3ceff] text-black hover:bg-[#837dbd] disabled:opacity-50">
          [START UPLOAD]
        </button>
        <button @click="clearCompleted" class="px-4 py-2 border border-[#837dbd] text-[#837dbd] hover:border-[#d3ceff]">
          [CLEAR COMPLETED]
        </button>
      </div>
    </div>

    <!-- Upload Stats -->
    <div v-if="stats.total > 0" class="grid grid-cols-4 gap-4">
      <div class="border border-[#837dbd] p-4 text-center">
        <div class="text-xs text-[#837dbd]">TOTAL</div>
        <div class="text-2xl text-[#d3ceff]">{{ stats.total }}</div>
      </div>
      <div class="border border-[#837dbd] p-4 text-center">
        <div class="text-xs text-[#837dbd]">SUCCESS</div>
        <div class="text-2xl text-green-400">{{ stats.completed }}</div>
      </div>
      <div class="border border-[#837dbd] p-4 text-center">
        <div class="text-xs text-[#837dbd]">FAILED</div>
        <div class="text-2xl text-red-400">{{ stats.failed }}</div>
      </div>
      <div class="border border-[#837dbd] p-4 text-center">
        <div class="text-xs text-[#837dbd]">PENDING</div>
        <div class="text-2xl text-yellow-400">{{ stats.pending }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useLibraryStore } from '../../stores/library'

interface UploadItem {
  file: File
  status: 'pending' | 'uploading' | 'completed' | 'error'
  progress: number
  error?: string
}

const libraryStore = useLibraryStore()
const fileInput = ref<HTMLInputElement>()
const isDragging = ref(false)
const uploadQueue = ref<UploadItem[]>([])
const uploadStatus = ref('READY')
const isUploading = ref(false)

const stats = computed(() => ({
  total: uploadQueue.value.length,
  completed: uploadQueue.value.filter(item => item.status === 'completed').length,
  failed: uploadQueue.value.filter(item => item.status === 'error').length,
  pending: uploadQueue.value.filter(item => item.status === 'pending').length
}))

const hasPendingFiles = computed(() => {
  return uploadQueue.value.some(item => item.status === 'pending')
})

const uploadStatusClass = computed(() => {
  switch (uploadStatus.value) {
    case 'UPLOADING': return 'text-yellow-400 animate-pulse'
    case 'COMPLETED': return 'text-green-400'
    case 'ERROR': return 'text-red-400'
    default: return 'text-[#837dbd]'
  }
})

const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files) {
    addFilesToQueue(Array.from(target.files))
    target.value = ''
  }
}

const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  if (event.dataTransfer?.files) {
    const files = Array.from(event.dataTransfer.files)
    const mp3Files = files.filter(file => 
      file.type === 'audio/mpeg' || file.name.toLowerCase().endsWith('.mp3')
    )
    addFilesToQueue(mp3Files)
  }
}

const addFilesToQueue = (files: File[]) => {
  files.forEach(file => {
    uploadQueue.value.push({
      file,
      status: 'pending',
      progress: 0
    })
  })
}

const removeFromQueue = (index: number) => {
  uploadQueue.value.splice(index, 1)
}

const clearCompleted = () => {
  uploadQueue.value = uploadQueue.value.filter(item => 
    item.status !== 'completed' && item.status !== 'error'
  )
}

const startUpload = async () => {
  if (isUploading.value) return
  
  isUploading.value = true
  uploadStatus.value = 'UPLOADING'
  
  for (const item of uploadQueue.value.filter(item => item.status === 'pending')) {
    try {
      item.status = 'uploading'
      
      // Simulate progress for demo
      for (let i = 0; i <= 100; i += 10) {
        await new Promise(resolve => setTimeout(resolve, 100))
        item.progress = i
      }
      
      // Actual upload
      await libraryStore.uploadSong(item.file)
      
      item.status = 'completed'
      item.progress = 100
    } catch (error) {
      item.status = 'error'
      item.error = error instanceof Error ? error.message : 'Upload failed'
    }
  }
  
  isUploading.value = false
  uploadStatus.value = uploadQueue.value.some(item => item.status === 'error') ? 'ERROR' : 'COMPLETED'
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
</script>