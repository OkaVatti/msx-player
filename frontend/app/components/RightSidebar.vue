<template>
  <div class="w-80 bg-black border-l border-[#837dbd] p-4 overflow-y-auto">
    <!-- Now Playing -->
    <div class="mb-6">
      <div class="text-xs text-[#837dbd] uppercase tracking-wider mb-4">NOW PLAYING</div>
      
      <div v-if="playerStore.currentSong" class="space-y-4">
        <!-- Album Art -->
        <div class="border border-[#837dbd] p-4 bg-[#111]">
          <div class="aspect-square bg-linear-to-br from-[#837dbd] to-[#d3ceff] flex items-center justify-center">
            <span class="text-4xl text-black">♪</span>
          </div>
        </div>

        <!-- Song Info -->
        <div>
          <div class="text-lg font-bold text-[#d3ceff] truncate">{{ playerStore.currentSong.title }}</div>
          <div class="text-sm text-[#837dbd]">{{ playerStore.currentSong.artist }}</div>
          <div class="text-xs text-[#837dbd] mt-1">{{ playerStore.currentSong.album }}</div>
        </div>

        <!-- Rating -->
        <div class="space-y-2">
          <div class="text-xs text-[#837dbd]">RATING</div>
          <div class="flex justify-center space-x-1">
            <button
              v-for="star in 10"
              :key="star"
              @click="rateSong(star)"
              class="text-xl transition-all duration-200 hover:scale-125"
              :class="star <= (playerStore.currentSong?.rating || 0) ? 'text-yellow-400' : 'text-[#333]'"
            >
              ★
            </button>
          </div>
        </div>

        <!-- Progress Bar -->
        <div class="space-y-2">
          <div class="flex justify-between text-xs text-[#837dbd]">
            <span>{{ formatTime(playerStore.currentTime) }}</span>
            <span>{{ formatTime(playerStore.duration) }}</span>
          </div>
          <div class="h-1 bg-[#333] cursor-pointer" @click="seekToTime">
            <div class="h-full bg-[#d3ceff]" :style="{ width: playerStore.progressPercentage + '%' }"></div>
          </div>
        </div>
      </div>
      
      <div v-else class="text-center p-8 border border-[#837dbd]">
        <pre class="text-[#837dbd] text-xs">
╔═══════════════════════════╗
║   NO SONG PLAYING         ║
║   SELECT A TRACK TO START ║
╚═══════════════════════════╝</pre>
      </div>
    </div>

    <!-- Up Next -->
    <div class="mb-6">
      <div class="text-xs text-[#837dbd] uppercase tracking-wider mb-2 flex justify-between">
        <span>UP NEXT</span>
        <button @click="shuffleQueue" class="text-xs text-[#d3ceff] hover:text-white">SHUFFLE</button>
      </div>
      <div class="space-y-2">
        <div v-for="song in upNextSongs" :key="song.id" 
             class="flex items-center gap-3 p-2 hover:bg-[#1a1a1a] cursor-pointer"
             @click="playSong(song)">
          <div class="text-xs text-[#837dbd] w-4">{{ song.index }}</div>
          <div class="flex-1 min-w-0">
            <div class="text-sm truncate text-[#d3ceff]">{{ song.title }}</div>
            <div class="text-xs text-[#837dbd] truncate">{{ song.artist }}</div>
          </div>
          <div class="text-xs text-[#837dbd]">{{ formatDuration(song.duration) }}</div>
        </div>
      </div>
    </div>

    <!-- Mini Visualizer -->
    <div class="border border-[#837dbd] p-4">
      <div class="text-xs text-[#837dbd] uppercase tracking-wider mb-2 flex justify-between items-center">
        <span>VISUALIZER</span>
        <button @click="cycleVisualizerMode" class="text-xs text-[#d3ceff] hover:text-white">
          [{{ visualizerMode.toUpperCase() }}]
        </button>
      </div>
      <div class="h-32 bg-black border border-[#837dbd] relative overflow-hidden">
        <canvas ref="visualizerCanvas" class="w-full h-full"></canvas>
        <div class="absolute inset-0 pointer-events-none scanlines"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { usePlayerStore } from '../../stores/player'
import type { Song } from '../../types'

const playerStore = usePlayerStore()
const visualizerCanvas = ref<HTMLCanvasElement>()

const visualizerMode = ref<'wave' | 'bars' | 'circular'>('wave')

const upNextSongs = computed(() => {
  if (!playerStore.currentSong || !playerStore.queue.length) return []
  
  const currentIndex = playerStore.queue.findIndex(song => song.id === playerStore.currentSong?.id)
  return playerStore.queue.slice(currentIndex + 1, currentIndex + 6).map((song, i) => ({
    ...song,
    index: i + 1
  }))
})

const formatTime = (seconds: number) => {
  if (!seconds || isNaN(seconds)) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const formatDuration = (seconds: number) => {
  if (!seconds) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const seekToTime = (event: MouseEvent) => {
  const element = event.currentTarget as HTMLElement
  const rect = element.getBoundingClientRect()
  const percent = (event.clientX - rect.left) / rect.width
  const newTime = percent * (playerStore.duration as number)
  playerStore.seek(newTime)
}

const rateSong = async (rating: number) => {
  if (playerStore.currentSong) {
    await playerStore.rateSong(playerStore.currentSong.id, rating)
  }
}

const playSong = async (song: Song) => {
  await playerStore.playSong(song)
}

const shuffleQueue = () => {
  playerStore.shuffleQueue()
}

const cycleVisualizerMode = () => {
  const modes: Array<'wave' | 'bars' | 'circular'> = ['wave', 'bars', 'circular']
  const currentIndex = modes.indexOf(visualizerMode.value)
  const nextIndex = (currentIndex + 1) % modes.length
  visualizerMode.value = modes[nextIndex]
}

let animationFrame: number
let audioContext: AudioContext | null = null
let analyser: AnalyserNode | null = null
let dataArray: Uint8Array | null = null

const setupAudioContext = async () => {
  if (!audioContext) {
    audioContext = new (window.AudioContext || (window as any).webkitAudioContext)()
    analyser = audioContext.createAnalyser()
    analyser.fftSize = 128
    dataArray = new Uint8Array(analyser.frequencyBinCount)
    // Some browsers require user gesture: resume when created
    await audioContext.resume().catch(() => {})
  }
}

const connectAudioElement = async (audioElement: HTMLAudioElement | null) => {
  if (!audioElement) return
  await setupAudioContext()
  if (!audioContext || !analyser) return

  try {
    const source = audioContext.createMediaElementSource(audioElement)
    source.connect(analyser)
    analyser.connect(audioContext.destination)
  } catch (err) {
    // connecting twice throws in some cases; ignore non-fatal
    console.warn("Visualizer connectAudioElement:", err)
  }
}

const drawVisualizer = () => {
  if (!visualizerCanvas.value) return
  
  const canvas = visualizerCanvas.value
  const ctx = canvas.getContext('2d')
  if (!ctx) return
  
  // Set canvas dimensions
  canvas.width = canvas.offsetWidth
  canvas.height = canvas.offsetHeight
  
  const width = canvas.width
  const height = canvas.height
  
  // Clear with fade effect
  ctx.fillStyle = 'rgba(0, 0, 0, 0.1)'
  ctx.fillRect(0, 0, width, height)
  
  // Get or generate audio data (ensure ArrayBuffer-backed Uint8Array)
  let audioData: Uint8Array = dataArray ?? new Uint8Array(new ArrayBuffer(64));
  
  if (analyser) {
    // Create a view that TypeScript will accept as ArrayBuffer-backed.
    // If audioData.buffer is already an ArrayBuffer this is zero-copy.
    const safeArray = new Uint8Array(
      audioData.buffer as ArrayBuffer,
      audioData.byteOffset,
      audioData.length
    );
  
    analyser.getByteFrequencyData(safeArray);
  
    // use safeArray for downstream calculations/rendering
    audioData = safeArray;
  } else {
    // Demo data — create a concrete ArrayBuffer-backed array to avoid typing issues
    audioData = new Uint8Array(new ArrayBuffer(64));
    const time = Date.now() / 1000;
    for (let i = 0; i < audioData.length; i++) {
      // keep values inside [0,255]
      audioData[i] = Math.floor(Math.abs(Math.sin(time * 2 + i * 0.1)) * 255);
    }
  }

  
  // Draw based on mode
  switch (visualizerMode.value) {
    case 'wave':
      drawWave(ctx, width, height, audioData)
      break
    case 'bars':
      drawBars(ctx, width, height, audioData)
      break
    case 'circular':
      drawCircular(ctx, width, height, audioData)
      break
  }
  
  animationFrame = requestAnimationFrame(drawVisualizer)
}

const drawWave = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  ctx.beginPath()
  ctx.lineWidth = 2
  ctx.strokeStyle = '#d3ceff'
  
  const sliceWidth = width / data.length
  let x = 0

  for (let i = 0; i < data.length; i++) {
    const v = data[i] / 255
    const y = (1 - v) * height

    if (i === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }

    x += sliceWidth
  }

  ctx.stroke()
}

const drawBars = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const barCount = 16
  const barWidth = width / barCount
  
  for (let i = 0; i < barCount; i++) {
    const dataIndex = Math.floor((i / barCount) * data.length)
    const amplitude = data[dataIndex] / 255
    const barHeight = amplitude * height * 0.8
    
    const x = i * barWidth
    const y = height - barHeight
    
    // Gradient effect
    const gradient = ctx.createLinearGradient(x, y, x, y + barHeight)
    gradient.addColorStop(0, '#d3ceff')
    gradient.addColorStop(0.5, '#837dbd')
    gradient.addColorStop(1, '#5a5586')
    
    ctx.fillStyle = gradient
    ctx.fillRect(x + 1, y, barWidth - 2, barHeight)
  }
}

const drawCircular = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const centerX = width / 2
  const centerY = height / 2
  const baseRadius = Math.min(width, height) / 6

  // Draw base circle
  ctx.beginPath()
  ctx.arc(centerX, centerY, baseRadius, 0, Math.PI * 2)
  ctx.strokeStyle = '#837dbd'
  ctx.lineWidth = 1
  ctx.stroke()

  // Draw frequency response
  ctx.beginPath()
  for (let i = 0; i < data.length; i++) {
    const amplitude = data[i] / 255
    const angle = (i / data.length) * Math.PI * 2
    const radius = baseRadius + amplitude * baseRadius
    
    const x = centerX + Math.cos(angle) * radius
    const y = centerY + Math.sin(angle) * radius
    
    if (i === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  }
  ctx.closePath()
  
  ctx.fillStyle = 'rgba(131, 125, 189, 0.2)'
  ctx.fill()
  ctx.strokeStyle = '#d3ceff'
  ctx.lineWidth = 1
  ctx.stroke()
}

onMounted(async () => {
  drawVisualizer()

  // Try to obtain audio element from store first, otherwise DOM query
  let audioEl: HTMLAudioElement | null = null
  try {
    // If the store stores an element reference directly, use it
    const maybe = (playerStore as any).audioElement
    if (maybe && typeof (maybe as any).currentTime === "number") {
      audioEl = maybe as unknown as HTMLAudioElement
    } else {
      audioEl = document.querySelector('audio') as HTMLAudioElement | null
    }
  } catch {
    audioEl = document.querySelector('audio') as HTMLAudioElement | null
  }

  if (audioEl) {
    await connectAudioElement(audioEl)
  }
})

onUnmounted(() => {
  if (animationFrame) {
    cancelAnimationFrame(animationFrame)
  }
  if (audioContext) {
    audioContext.close().catch(() => {})
  }
})

// Watch for audio element changes (store may update later)
watch(() => (playerStore as any).audioElement, (el) => {
  const audioEl = (el && typeof (el as any).currentTime === "number") ? (el as unknown as HTMLAudioElement) : document.querySelector('audio') as HTMLAudioElement | null
  if (audioEl) connectAudioElement(audioEl)
})
</script>
