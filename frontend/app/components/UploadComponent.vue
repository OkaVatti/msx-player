<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="border-2 border-lime-400 p-4 bg-black">
      <h2 class="text-xl mb-2 flex items-center gap-2">
        <span class="animate-pulse">&gt;&gt;</span>
        FILE UPLOAD SYSTEM
        <span class="text-xs text-white">[DRAG & DROP SUPPORTED]</span>
      </h2>
      <div class="text-sm text-lime-300 font-mono">
        UPLOAD STATUS: <span :class="uploadStatusClass">{{ uploadStatus }}</span>
        <span v-if="debugInfo" class="text-yellow-400 ml-2">[DEBUG: {{ debugInfo }}]</span>
      </div>
    </div>

    <!-- Upload Area -->
    <div 
      class="border-2 border-dashed border-lime-400 p-8 text-center bg-black transition-all duration-300 hover:border-solid hover:bg-lime-400 hover:text-black"
      @drop="handleDrop"
      @dragover="handleDragOver"
      @dragleave="handleDragLeave"
      :class="{ 'bg-lime-400 text-black border-solid': isDragging }"
    >
      <div class="space-y-3">
        <div class="text-4xl">📁</div>
        <h3 class="text-lg font-mono">DROP MP3 FILES HERE</h3>
        <p class="text-sm opacity-80 font-mono">OR CLICK TO BROWSE</p>
        <input
          ref="fileInput"
          type="file"
          multiple
          accept=".mp3,audio/mpeg"
          @change="handleFileSelect"
          class="hidden"
        />
        <button
          @click="triggerFileInput"
          class="px-6 py-3 border-2 border-current bg-transparent hover:bg-black hover:text-lime-400 transition-all font-mono"
        >
          [SELECT FILES]
        </button>
      </div>
    </div>

    <!-- Connection Status -->
    <div v-if="connectionStatus" class="border-2 border-yellow-400 p-4 bg-black">
      <h3 class="text-lg mb-2 flex items-center gap-2 text-yellow-400">
        <span class="animate-pulse">⚠</span>
        CONNECTION STATUS
      </h3>
      <div class="text-sm text-yellow-300 font-mono">
        {{ connectionStatus }}
      </div>
    </div>

    <!-- Upload Queue -->
    <div v-if="uploadQueue.length > 0" class="border-2 border-lime-400 p-4 bg-black">
      <h3 class="text-lg mb-3 flex items-center gap-2">
        <span class="animate-pulse">&gt;</span>
        UPLOAD QUEUE
        <span class="text-xs text-white">[{{ uploadQueue.length }} FILES]</span>
      </h3>
      
      <div class="space-y-2 max-h-64 overflow-y-auto">
        <div
          v-for="(file, index) in uploadQueue"
          :key="file.id"
          class="border border-lime-400 p-3 font-mono text-sm"
        >
          <div class="flex justify-between items-center">
            <div class="flex-1 min-w-0">
              <div class="text-lime-300 truncate">{{ file.name }}</div>
              <div class="text-lime-500 text-xs">
                {{ formatFileSize(file.size) }} • 
                <span :class="getStatusClass(file.status)">{{ file.status.toUpperCase() }}</span>
                <span v-if="file.status === 'uploading'" class="text-yellow-400">
                  • {{ file.progress }}%
                </span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <!-- Progress Bar -->
              <div v-if="file.status === 'uploading'" class="w-20 bg-lime-900 border border-lime-400">
                <div 
                  class="h-2 bg-lime-400 transition-all duration-300"
                  :style="{ width: `${file.progress}%` }"
                ></div>
              </div>
              
              <!-- Status Icon -->
              <span v-if="file.status === 'pending'" class="text-lime-500">[ ]</span>
              <span v-else-if="file.status === 'uploading'" class="text-yellow-400 animate-pulse">[~]</span>
              <span v-else-if="file.status === 'completed'" class="text-green-400">[✓]</span>
              <span v-else-if="file.status === 'error'" class="text-red-400">[!]</span>
              
              <button
                @click="removeFromQueue(index)"
                class="text-lime-500 hover:text-red-400 transition-colors"
                :disabled="file.status === 'uploading'"
                :class="{ 'opacity-50 cursor-not-allowed': file.status === 'uploading' }"
              >
                [X]
              </button>
            </div>
          </div>
          
          <!-- Error Message -->
          <div v-if="file.error" class="text-red-400 text-xs mt-1 font-mono">
            ERROR: {{ file.error }}
          </div>

          <!-- Success Message -->
          <div v-if="file.status === 'completed'" class="text-green-400 text-xs mt-1 font-mono">
            UPLOAD SUCCESSFUL • ADDED TO LIBRARY
          </div>
        </div>
      </div>

      <!-- Queue Controls -->
      <div class="flex gap-2 mt-4 pt-3 border-t border-lime-400">
        <button
          @click="startUpload"
          :disabled="!hasPendingFiles || isUploading"
          class="px-4 py-2 border-2 border-lime-400 bg-lime-400 text-black hover:bg-black hover:text-lime-400 disabled:opacity-50 disabled:cursor-not-allowed transition-all font-mono text-sm"
        >
          [{{ isUploading ? 'UPLOADING...' : 'START UPLOAD' }}]
        </button>
        <button
          @click="clearCompleted"
          class="px-4 py-2 border-2 border-lime-400 bg-black text-lime-400 hover:bg-lime-400 hover:text-black transition-all font-mono text-sm"
        >
          [CLEAR COMPLETED]
        </button>
        <button
          @click="clearQueue"
          class="px-4 py-2 border-2 border-red-400 bg-black text-red-400 hover:bg-red-400 hover:text-black transition-all font-mono text-sm"
        >
          [CLEAR ALL]
        </button>
      </div>
    </div>

    <!-- Upload Stats -->
    <div v-if="uploadStats.total > 0" class="border-2 border-lime-400 p-4 bg-black">
      <h3 class="text-lg mb-3 flex items-center gap-2">
        <span class="animate-pulse">&gt;</span>
        UPLOAD STATISTICS
      </h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm font-mono">
        <div class="text-center p-2 border border-lime-400">
          <div class="text-lime-300">TOTAL</div>
          <div class="text-xl text-white">{{ uploadStats.total }}</div>
        </div>
        <div class="text-center p-2 border border-lime-400">
          <div class="text-lime-300">SUCCESS</div>
          <div class="text-xl text-green-400">{{ uploadStats.completed }}</div>
        </div>
        <div class="text-center p-2 border border-lime-400">
          <div class="text-lime-300">FAILED</div>
          <div class="text-xl text-red-400">{{ uploadStats.failed }}</div>
        </div>
        <div class="text-center p-2 border border-lime-400">
          <div class="text-lime-300">PENDING</div>
          <div class="text-xl text-yellow-400">{{ uploadStats.pending }}</div>
        </div>
      </div>
    </div>

    <!-- Duplicate Protection Info -->
    <div class="border-2 border-blue-400 p-4 bg-black">
      <h3 class="text-lg mb-2 flex items-center gap-2 text-blue-400">
        <span class="animate-pulse">ℹ</span>
        DUPLICATE PROTECTION
      </h3>
      <div class="text-sm text-blue-300 font-mono space-y-1">
        <div>• SYSTEM AUTOMATICALLY DETECTS DUPLICATE SONGS</div>
        <div>• EXPLICIT/CLEAN VERSIONS ARE ALLOWED</div>
        <div>• HIGHEST RATED VERSION IS KEPT</div>
        <div>• OLDEST ID WINS ON TIE</div>
      </div>
    </div>

    <!-- Debug Panel -->
    <div v-if="showDebug" class="border-2 border-red-400 p-4 bg-black">
      <h3 class="text-lg mb-2 flex items-center gap-2 text-red-400">
        <span class="animate-pulse">🐛</span>
        DEBUG PANEL
      </h3>
      <div class="text-sm text-red-300 font-mono space-y-2">
        <div>Backend URL: {{ backendUrl }}</div>
        <div>Connection: {{ connectionTest }}</div>
        <div>Last Error: {{ lastError }}</div>
        <button
          @click="testBackendConnection"
          class="px-4 py-2 border-2 border-red-400 bg-black text-red-400 hover:bg-red-400 hover:text-black transition-all font-mono text-sm"
        >
          [TEST CONNECTION]
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';

interface UploadFile {
  id: string;
  file: File;
  name: string;
  size: number;
  status: 'pending' | 'uploading' | 'completed' | 'error';
  progress: number;
  error?: string;
}

const fileInput = ref<HTMLInputElement>();
const isDragging = ref(false);
const uploadQueue = ref<UploadFile[]>([]);
const uploadStatus = ref('READY');
const isUploading = ref(false);
const activeUploads = ref(0);
const connectionStatus = ref('');
const debugInfo = ref('');
const showDebug = ref(true); // Set to true for debugging
const backendUrl = ref('http://localhost:1323/api/songs');
const connectionTest = ref('NOT TESTED');
const lastError = ref('');

const uploadStats = ref({
  total: 0,
  completed: 0,
  failed: 0,
  pending: 0
});

const uploadStatusClass = computed(() => {
  switch (uploadStatus.value) {
    case 'READY': return 'text-lime-400';
    case 'UPLOADING': return 'text-yellow-400 animate-pulse';
    case 'COMPLETED': return 'text-green-400';
    case 'ERROR': return 'text-red-400';
    default: return 'text-lime-400';
  }
});

const hasPendingFiles = computed(() => {
  return uploadQueue.value.some(file => file.status === 'pending');
});

const getStatusClass = (status: string) => {
  switch (status) {
    case 'pending': return 'text-lime-400';
    case 'uploading': return 'text-yellow-400';
    case 'completed': return 'text-green-400';
    case 'error': return 'text-red-400';
    default: return 'text-lime-400';
  }
};

// Test backend connection on component mount
onMounted(async () => {
  await testBackendConnection();
});

const testBackendConnection = async () => {
  try {
    connectionStatus.value = 'Testing connection...';
    const response = await fetch('http://localhost:1323/api/songs');
    
    if (response.ok) {
      connectionStatus.value = '✅ Backend connected successfully';
      connectionTest.value = 'SUCCESS';
      debugInfo.value = 'Backend OK';
    } else {
      connectionStatus.value = '❌ Backend responded with error';
      connectionTest.value = `FAILED: ${response.status}`;
      debugInfo.value = `HTTP ${response.status}`;
    }
  } catch (error) {
    connectionStatus.value = '❌ Cannot connect to backend';
    connectionTest.value = 'FAILED: Network error';
    debugInfo.value = 'Network error';
    lastError.value = error instanceof Error ? error.message : 'Unknown error';
    console.error('Backend connection test failed:', error);
  }
};

const triggerFileInput = () => {
  fileInput.value?.click();
};

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files.length > 0) {
    addFilesToQueue(Array.from(target.files));
    // Reset the input to allow uploading the same file again
    target.value = '';
  }
};

const handleDragOver = (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = true;
};

const handleDragLeave = (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = false;
};

const handleDrop = (event: DragEvent) => {
  event.preventDefault();
  isDragging.value = false;
  
  if (event.dataTransfer?.files && event.dataTransfer.files.length > 0) {
    addFilesToQueue(Array.from(event.dataTransfer.files));
  }
};

const addFilesToQueue = (files: File[]) => {
  const mp3Files = files.filter(file => 
    file.type === 'audio/mpeg' || 
    file.name.toLowerCase().endsWith('.mp3')
  );

  // Filter out duplicates by filename
  const existingFilenames = new Set(uploadQueue.value.map(f => f.name));
  const newFiles = mp3Files.filter(file => !existingFilenames.has(file.name));

  if (newFiles.length === 0) {
    connectionStatus.value = '⚠ Duplicate files detected';
    return;
  }

  newFiles.forEach(file => {
    uploadQueue.value.push({
      id: Math.random().toString(36).substr(2, 9),
      file,
      name: file.name,
      size: file.size,
      status: 'pending',
      progress: 0
    });
  });

  updateStats();
  connectionStatus.value = `📥 Added ${newFiles.length} files to queue`;
};

const removeFromQueue = (index: number) => {
  if (index >= 0 && index < uploadQueue.value.length) {
    const file = uploadQueue.value[index];
    if (file.status !== 'uploading') {
      uploadQueue.value.splice(index, 1);
      updateStats();
    }
  }
};

const clearCompleted = () => {
  uploadQueue.value = uploadQueue.value.filter(file => 
    file.status !== 'completed' && file.status !== 'error'
  );
  updateStats();
};

const clearQueue = () => {
  // Only clear non-uploading files
  uploadQueue.value = uploadQueue.value.filter(file => file.status === 'uploading');
  updateStats();
};

const startUpload = async () => {
  if (isUploading.value) return;
  
  isUploading.value = true;
  uploadStatus.value = 'UPLOADING';
  activeUploads.value = uploadQueue.value.filter(f => f.status === 'pending').length;

  // Test connection before starting upload
  await testBackendConnection();
  
  if (connectionTest.value !== 'SUCCESS') {
    uploadStatus.value = 'ERROR';
    connectionStatus.value = '❌ Cannot upload: Backend not connected';
    isUploading.value = false;
    return;
  }

  // Process files sequentially to avoid overwhelming the server
  for (const queueItem of uploadQueue.value.filter(f => f.status === 'pending')) {
    await uploadFile(queueItem);
    activeUploads.value--;
  }

  isUploading.value = false;
  uploadStatus.value = 'COMPLETED';
};

const uploadFile = async (queueItem: UploadFile) => {
  queueItem.status = 'uploading';
  
  try {
    const formData = new FormData();
    formData.append('file', queueItem.file);

    // Use fetch with progress tracking
    const xhr = new XMLHttpRequest();
    
    // Track upload progress
    xhr.upload.addEventListener('progress', (event) => {
      if (event.lengthComputable) {
        const progress = (event.loaded / event.total) * 100;
        queueItem.progress = Math.round(progress);
      }
    });

    // Wait for the upload to complete
    await new Promise((resolve, reject) => {
      xhr.onreadystatechange = () => {
        if (xhr.readyState === 4) {
          if (xhr.status === 200 || xhr.status === 201) {
            try {
              const response = JSON.parse(xhr.responseText);
              console.log('Upload successful:', response);
              resolve(response);
            } catch (e) {
              resolve(xhr.responseText);
            }
          } else {
            let errorMsg = `Upload failed: ${xhr.status} ${xhr.statusText}`;
            try {
              const errorResponse = JSON.parse(xhr.responseText);
              errorMsg = errorResponse.error || errorMsg;
            } catch (e) {
              // Ignore JSON parse error
            }
            reject(new Error(errorMsg));
          }
        }
      };

      xhr.onerror = () => reject(new Error('Network error during upload'));
      xhr.ontimeout = () => reject(new Error('Upload timeout'));
      
      xhr.open('POST', 'http://localhost:1323/api/songs');
      xhr.timeout = 30000; // 30 second timeout
      xhr.send(formData);
    });

    // Success
    queueItem.status = 'completed';
    queueItem.progress = 100;
    connectionStatus.value = `✅ Uploaded: ${queueItem.name}`;
    
  } catch (error) {
    console.error('Upload error:', error);
    queueItem.status = 'error';
    queueItem.error = error instanceof Error ? error.message : 'Upload failed';
    lastError.value = queueItem.error;
    connectionStatus.value = `❌ Upload failed: ${queueItem.name}`;
  }
  
  updateStats();
};

const updateStats = () => {
  const stats = {
    total: uploadQueue.value.length,
    completed: uploadQueue.value.filter(f => f.status === 'completed').length,
    failed: uploadQueue.value.filter(f => f.status === 'error').length,
    pending: uploadQueue.value.filter(f => f.status === 'pending').length
  };
  
  uploadStats.value = stats;
};

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};
</script>

<style scoped>
/* Custom scrollbar for upload queue */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #000;
  border: 1px solid #84cc16;
}

::-webkit-scrollbar-thumb {
  background: #84cc16;
}

::-webkit-scrollbar-thumb:hover {
  background: #a3e635;
}
</style>