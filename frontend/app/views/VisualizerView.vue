[file name]: VisualizerView.vue
[file content begin]
<template>
  <div class="h-full flex flex-col p-4 space-y-4">
    <!-- Header -->
    <div class="border-2 border-[#837dbd] p-4 bg-black">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-xl text-[#d3ceff] font-bold tracking-wide flex items-center gap-2">
            <span class="animate-pulse">&gt;&gt;</span>
            AUDIO VISUALIZER
          </h1>
          <p class="text-sm text-[#837dbd] font-mono mt-1">
            Real-time audio visualization • {{ currentMode.toUpperCase() }} mode
          </p>
        </div>
        <button
          v-if="!isConnected"
          @click="connectToCurrentAudio"
          class="px-4 py-2 bg-[#837dbd] text-black hover:bg-[#d3ceff] transition-all font-mono text-sm flex items-center gap-2"
        >
          <span>[CONNECT AUDIO]</span>
        </button>
        <div v-else class="flex items-center gap-3">
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 bg-green-400 animate-pulse rounded-full"></div>
            <span class="text-sm text-green-400 font-mono">CONNECTED</span>
          </div>
          <button
            @click="disconnectAudio"
            class="px-3 py-1 border-2 border-red-400 text-red-400 hover:bg-red-400 hover:text-black transition-all font-mono text-xs"
          >
            [DISCONNECT]
          </button>
        </div>
      </div>
    </div>

    <!-- Mode Selector -->
    <div class="border-2 border-[#837dbd] p-4 bg-black">
      <div class="flex flex-wrap gap-2">
        <button
          v-for="mode in visualizerModes"
          :key="mode"
          @click="currentMode = mode"
          :class="[
            'px-4 py-2 border-2 transition-all font-mono text-sm',
            currentMode === mode
              ? 'bg-[#837dbd] text-black border-[#837dbd]'
              : 'bg-black text-[#837dbd] border-[#837dbd] hover:bg-[#837dbd] hover:text-black'
          ]"
        >
          {{ mode.toUpperCase() }}
        </button>
      </div>
    </div>

    <!-- Main Visualizer Canvas -->
    <div class="flex-1 border-2 border-[#837dbd] bg-black relative overflow-hidden">
      <canvas
        ref="canvas"
        class="w-full h-full cursor-grab active:cursor-grabbing"
        @mousedown="startDrag"
        @mousemove="drag"
        @mouseup="stopDrag"
        @mouseleave="stopDrag"
        @wheel="handleWheel"
      />

      <!-- Visualizer Overlay Controls -->
      <div class="absolute top-4 left-4 flex flex-col gap-3">
        <div class="bg-black border border-[#837dbd] p-3 font-mono">
          <div class="text-xs text-[#837dbd]">MODE</div>
          <div class="text-[#d3ceff] text-sm">{{ currentMode.toUpperCase() }}</div>
        </div>
        <div class="bg-black border border-[#837dbd] p-3 font-mono">
          <div class="text-xs text-[#837dbd]">ROTATION</div>
          <div class="text-[#d3ceff] text-sm">X: {{ rotation.x.toFixed(0) }}°</div>
          <div class="text-[#d3ceff] text-sm">Y: {{ rotation.y.toFixed(0) }}°</div>
        </div>
        <div class="bg-black border border-[#837dbd] p-3 font-mono">
          <div class="text-xs text-[#837dbd]">SCALE</div>
          <div class="text-[#d3ceff] text-sm">{{ scale.toFixed(1) }}x</div>
        </div>
      </div>

      <!-- Bottom Controls -->
      <div class="absolute bottom-4 left-4 right-4 flex items-center justify-between">
        <div class="bg-black border border-[#837dbd] p-3 font-mono">
          <div class="text-xs text-[#837dbd]">FPS</div>
          <div class="text-[#d3ceff] text-sm">{{ fps }}</div>
        </div>
        
        <div class="flex items-center gap-3">
          <button
            @click="toggleVisualizer"
            class="px-4 py-2 border-2 border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black transition-all font-mono text-sm"
          >
            {{ isRunning ? '[PAUSE]' : '[PLAY]' }}
          </button>
          <button
            @click="resetView"
            class="px-4 py-2 border-2 border-[#837dbd] text-[#837dbd] hover:bg-[#837dbd] hover:text-black transition-all font-mono text-sm"
          >
            [RESET VIEW]
          </button>
        </div>
      </div>

      <!-- Audio Level Indicator -->
      <div class="absolute top-4 right-4">
        <div class="bg-black border border-[#837dbd] p-3">
          <div class="text-xs text-[#837dbd] font-mono mb-2">AUDIO LEVEL</div>
          <div class="flex items-end gap-1 h-8">
            <div
              v-for="n in 10"
              :key="n"
              class="w-2 bg-[#837dbd] transition-all duration-100"
              :style="{ height: `${(audioLevel / 255) * (n * 2) + 2}px` }"
              :class="{ 'animate-pulse': isConnected }"
            ></div>
          </div>
        </div>
      </div>

      <!-- Scanlines Overlay -->
      <div class="absolute inset-0 pointer-events-none scanlines"></div>
    </div>

    <!-- Instructions -->
    <div class="border-2 border-[#837dbd] p-4 bg-black">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm font-mono">
        <div class="text-center">
          <div class="text-[#837dbd]">DRAG</div>
          <div class="text-[#d3ceff]">Rotate visualizer</div>
        </div>
        <div class="text-center">
          <div class="text-[#837dbd]">SCROLL</div>
          <div class="text-[#d3ceff]">Zoom in/out</div>
        </div>
        <div class="text-center">
          <div class="text-[#837dbd]">CLICK MODES</div>
          <div class="text-[#d3ceff]">Switch visualization</div>
        </div>
      </div>
    </div>

    <!-- Audio Source Info -->
    <div v-if="currentSong" class="border-2 border-[#837dbd] p-4 bg-black">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 bg-[#5a548d] flex items-center justify-center">
            <span class="text-[#d3ceff] text-xl">♪</span>
          </div>
          <div>
            <div class="text-[#d3ceff] font-mono">{{ currentSong.title }}</div>
            <div class="text-[#837dbd] text-sm font-mono">{{ currentSong.artist }}</div>
          </div>
        </div>
        <div class="text-right">
          <div class="text-[#837dbd] text-sm font-mono">NOW PLAYING</div>
          <div class="text-[#d3ceff] font-mono">{{ formatTime(currentTime) }} / {{ formatTime(currentSong.duration) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'

interface Song {
  id: number
  title: string
  artist: string
  duration: number
  file_path: string
}

const canvas = ref<HTMLCanvasElement>()
const ctx = ref<CanvasRenderingContext2D>()
const isDragging = ref(false)
const rotation = ref({ x: 0, y: 0 })
const scale = ref(1)
const currentMode = ref<'bars' | 'wave' | 'particles' | 'circular' | 'matrix'>('bars')
const fps = ref(0)
const isRunning = ref(true)
const isConnected = ref(false)
const audioLevel = ref(0)
const currentSong = ref<Song | null>(null)
const currentTime = ref(0)

const visualizerModes: typeof currentMode.value[] = ['bars', 'wave', 'particles', 'circular', 'matrix']

let animationFrame: number
let audioContext: AudioContext | null = null
let analyser: AnalyserNode | null = null
let dataArray: Uint8Array | null = null
let lastFrameTime = 0
let frameCount = 0
let sourceNode: MediaElementAudioSourceNode | null = null
let audioElement: HTMLAudioElement | null = null
let lastMousePos = { x: 0, y: 0 }

const setupAudioContext = () => {
  if (!audioContext) {
    audioContext = new (window.AudioContext || (window as any).webkitAudioContext)()
    analyser = audioContext.createAnalyser()
    analyser.fftSize = 256
    dataArray = new Uint8Array(analyser.frequencyBinCount)
  }
}

const connectToCurrentAudio = () => {
  // Try to find the audio element from the player
  const audioEl = document.querySelector('audio') as HTMLAudioElement
  if (!audioEl) {
    alert('No audio element found. Please play a song first.')
    return
  }

  setupAudioContext()
  
  if (audioContext && analyser) {
    sourceNode = audioContext.createMediaElementSource(audioEl)
    sourceNode.connect(analyser)
    analyser.connect(audioContext.destination)
    isConnected.value = true
    audioElement = audioEl
    
    // Listen for time updates
    audioEl.addEventListener('timeupdate', updateCurrentTime)
    
    // Get current song info from store or API
    fetchCurrentSong()
  }
}

const disconnectAudio = () => {
  if (sourceNode) {
    sourceNode.disconnect()
    sourceNode = null
  }
  
  if (audioElement) {
    audioElement.removeEventListener('timeupdate', updateCurrentTime)
    audioElement = null
  }
  
  isConnected.value = false
  currentSong.value = null
  currentTime.value = 0
}

const fetchCurrentSong = async () => {
  try {
    const response = await fetch('http://localhost:1323/api/player/state')
    if (!response.ok) throw new Error('Failed to fetch player state')
    
    const state = await response.json()
    if (state.current_song) {
      currentSong.value = state.current_song
      currentTime.value = state.current_time
    }
  } catch (err) {
    console.error('Error fetching current song:', err)
  }
}

const updateCurrentTime = () => {
  if (audioElement) {
    currentTime.value = audioElement.currentTime
  }
}

const formatTime = (seconds: number) => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const draw = (timestamp: number) => {
  if (!canvas.value || !ctx.value) return

  // Calculate FPS
  frameCount++
  if (timestamp - lastFrameTime >= 1000) {
    fps.value = frameCount
    frameCount = 0
    lastFrameTime = timestamp
  }

  const width = canvas.value.width
  const height = canvas.value.height

  // Clear with fade effect
  ctx.value.fillStyle = 'rgba(0, 0, 0, 0.1)'
  ctx.value.fillRect(0, 0, width, height)

  // Apply transformations
  ctx.value.save()
  ctx.value.translate(width / 2, height / 2)
  ctx.value.rotate(rotation.value.x * Math.PI / 180)
  ctx.value.scale(scale.value, scale.value)
  ctx.value.translate(-width / 2, -height / 2)

  // Get audio data
  let audioData = dataArray
  if (analyser && audioData) {
    analyser.getByteFrequencyData(audioData)
    
    // Calculate average audio level
    let sum = 0
    for (let i = 0; i < audioData.length; i++) {
      sum += audioData[i]
    }
    audioLevel.value = sum / audioData.length
  } else {
    // Generate demo data
    if (!audioData) {
      audioData = new Uint8Array(128)
    }
    for (let i = 0; i < audioData.length; i++) {
      audioData[i] = Math.random() * 255 * Math.sin(timestamp / 1000 + i / 10)
    }
    audioLevel.value = Math.random() * 255
  }

  // Draw based on current mode
  if (audioData) {
    switch (currentMode.value) {
      case 'bars':
        drawBars(ctx.value, width, height, audioData)
        break
      case 'wave':
        drawWave(ctx.value, width, height, audioData)
        break
      case 'particles':
        drawParticles(ctx.value, width, height, audioData, timestamp)
        break
      case 'circular':
        drawCircular(ctx.value, width, height, audioData)
        break
      case 'matrix':
        drawMatrix(ctx.value, width, height, timestamp)
        break
    }
  }

  ctx.value.restore()
  
  if (isRunning.value) {
    animationFrame = requestAnimationFrame(draw)
  }
}

const drawBars = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const barCount = 32
  const barWidth = width / barCount
  
  for (let i = 0; i < barCount; i++) {
    const dataIndex = Math.floor((i / barCount) * data.length)
    const amplitude = data[dataIndex] / 255
    const barHeight = amplitude * height * 0.8
    
    const x = i * barWidth
    const y = height - barHeight
    
    // Gradient from purple to cyan
    const gradient = ctx.createLinearGradient(x, y, x, y + barHeight)
    gradient.addColorStop(0, '#d3ceff')
    gradient.addColorStop(0.5, '#837dbd')
    gradient.addColorStop(1, '#5a548d')
    
    ctx.fillStyle = gradient
    ctx.fillRect(x + 1, y, barWidth - 2, barHeight)
    
    // Top highlight
    ctx.fillStyle = 'rgba(255, 255, 255, 0.3)'
    ctx.fillRect(x + 1, y, barWidth - 2, 1)
  }
}

const drawWave = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  ctx.beginPath()
  ctx.lineWidth = 2
  ctx.strokeStyle = '#d3ceff'
  
  const sliceWidth = width / data.length
  let x = 0

  for (let i = 0; i < data.length; i++) {
    const v = data[i] / 255
    const y = v * height / 2 + height / 4

    if (i === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }

    x += sliceWidth
  }

  ctx.stroke()
  
  // Glow effect
  ctx.strokeStyle = 'rgba(131, 125, 189, 0.4)'
  ctx.lineWidth = 4
  ctx.stroke()
}

const drawParticles = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array, timestamp: number) => {
  const centerX = width / 2
  const centerY = height / 2
  const time = timestamp / 1000

  for (let i = 0; i < 50; i++) {
    const dataIndex = i % data.length
    const amplitude = data[dataIndex] / 255
    const angle = (i / 50) * Math.PI * 2 + time
    const radius = 50 + amplitude * 200
    
    const x = centerX + Math.cos(angle) * radius
    const y = centerY + Math.sin(angle) * radius
    
    const size = 1 + amplitude * 4
    
    // Particle
    ctx.beginPath()
    ctx.arc(x, y, size, 0, Math.PI * 2)
    ctx.fillStyle = `hsl(${i * 7.2}, 100%, 70%)`
    ctx.fill()
    
    // Trail
    ctx.globalAlpha = 0.3
    ctx.beginPath()
    ctx.arc(x, y, size * 2, 0, Math.PI * 2)
    ctx.fillStyle = '#837dbd'
    ctx.fill()
    ctx.globalAlpha = 1
  }
}

const drawCircular = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const centerX = width / 2
  const centerY = height / 2
  const baseRadius = Math.min(width, height) / 4

  // Base circle
  ctx.beginPath()
  ctx.arc(centerX, centerY, baseRadius, 0, Math.PI * 2)
  ctx.strokeStyle = '#5a548d'
  ctx.lineWidth = 1
  ctx.stroke()

  // Frequency response
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

const drawMatrix = (ctx: CanvasRenderingContext2D, width: number, height: number, timestamp: number) => {
  const fontSize = 12
  const columns = Math.floor(width / fontSize)
  
  // Initialize drops if needed
  if (!(window as any).matrixDrops) {
    (window as any).matrixDrops = Array(columns).fill(1)
  }
  const drops = (window as any).matrixDrops as number[]
  
  // Fade effect
  ctx.fillStyle = 'rgba(0, 0, 0, 0.05)'
  ctx.fillRect(0, 0, width, height)
  
  // Draw characters
  ctx.fillStyle = '#837dbd'
  ctx.font = `${fontSize}px monospace`
  
  for (let i = 0; i < drops.length; i++) {
    const char = String.fromCharCode(0x30A0 + Math.random() * 96)
    const x = i * fontSize
    const y = drops[i] * fontSize
    
    // Random green tint
    ctx.fillStyle = Math.random() > 0.7 ? '#d3ceff' : '#837dbd'
    ctx.fillText(char, x, y)
    
    // Reset drop if it reaches bottom
    if (y > height && Math.random() > 0.975) {
      drops[i] = 0
    }
    
    drops[i] = drops[i] + 1
  }
}

const startDrag = (event: MouseEvent) => {
  isDragging.value = true
  lastMousePos = { x: event.clientX, y: event.clientY }
}

const drag = (event: MouseEvent) => {
  if (!isDragging.value) return
  
  const deltaX = event.clientX - lastMousePos.x
  const deltaY = event.clientY - lastMousePos.y
  
  rotation.value.y += deltaX * 0.5
  rotation.value.x += deltaY * 0.5
  
  lastMousePos = { x: event.clientX, y: event.clientY }
}

const stopDrag = () => {
  isDragging.value = false
}

const handleWheel = (event: WheelEvent) => {
  event.preventDefault()
  scale.value += event.deltaY * -0.001
  scale.value = Math.min(Math.max(0.3, scale.value), 3)
}

const toggleVisualizer = () => {
  isRunning.value = !isRunning.value
  if (isRunning.value) {
    draw(performance.now())
  }
}

const resetView = () => {
  rotation.value = { x: 0, y: 0 }
  scale.value = 1
}

const resizeCanvas = () => {
  if (canvas.value) {
    canvas.value.width = canvas.value.offsetWidth
    canvas.value.height = canvas.value.offsetHeight
  }
}

onMounted(() => {
  if (canvas.value) {
    resizeCanvas()
    ctx.value = canvas.value.getContext('2d')!
    draw(performance.now())
  }
  
  window.addEventListener('resize', resizeCanvas)
})

onUnmounted(() => {
  if (animationFrame) {
    cancelAnimationFrame(animationFrame)
  }
  
  if (audioContext) {
    audioContext.close()
  }
  
  disconnectAudio()
  window.removeEventListener('resize', resizeCanvas)
})

// Watch for mode changes to restart animation if paused
watch(currentMode, () => {
  if (!isRunning.value) {
    isRunning.value = true
    draw(performance.now())
  }
})
</script>

<style scoped>
.scanlines {
  background: linear-gradient(
    to bottom,
    transparent 50%,
    rgba(211, 206, 255, 0.05) 51%
  );
  background-size: 100% 4px;
  animation: scanline 8s linear infinite;
}

@keyframes scanline {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 0 100%;
  }
}

::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #000;
  border: 1px solid #5a548d;
}

::-webkit-scrollbar-thumb {
  background: #837dbd;
}

::-webkit-scrollbar-thumb:hover {
  background: #d3ceff;
}
</style>
[file content end]