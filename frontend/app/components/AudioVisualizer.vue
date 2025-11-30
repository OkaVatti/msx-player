<template>
  <div class="bg-gray-800 rounded-lg p-6">
    <h2 class="text-xl font-semibold mb-4">Audio Visualizer</h2>
    
    <div class="relative">
      <canvas
        ref="canvas"
        class="w-full h-64 bg-gray-900 rounded-lg cursor-grab active:cursor-grabbing"
        @mousedown="startDrag"
        @mousemove="drag"
        @mouseup="stopDrag"
        @mouseleave="stopDrag"
        @wheel="handleWheel"
      />
      
      <div class="absolute top-4 right-4 flex space-x-2">
        <button
          v-for="mode in visualizerModes"
          :key="mode"
          @click="currentMode = mode as typeof currentMode"
            class="px-3 py-1 bg-gray-700 rounded-lg text-sm hover:bg-gray-600 transition-colors"
            :class="{ 'bg-purple-500': currentMode === mode }
          ">
          {{ mode }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const canvas = ref<HTMLCanvasElement>()
const ctx = ref<CanvasRenderingContext2D>()
const isDragging = ref(false)
const rotation = ref({ x: 0, y: 0 })
const scale = ref(1)
const currentMode = ref<'bars' | 'wave' | 'particles' | 'circular'>('bars')

const visualizerModes = ['bars', 'wave', 'particles', 'circular']

let animationFrame: number
let audioContext: AudioContext
let analyser: AnalyserNode
let dataArray: Uint8Array<ArrayBuffer>
let source: MediaElementAudioSourceNode

const setupAudioContext = () => {
  if (!audioContext) {
    audioContext = new AudioContext()
    analyser = audioContext.createAnalyser()
    analyser.fftSize = 256
    dataArray = new Uint8Array(analyser.frequencyBinCount)
  }
}

const connectAudioElement = (audioElement: HTMLAudioElement) => {
  setupAudioContext()
  if (source) {
    source.disconnect()
  }
  source = audioContext.createMediaElementSource(audioElement)
  source.connect(analyser)
  analyser.connect(audioContext.destination)
}

const draw = () => {
  if (!canvas.value || !ctx.value) return

  const width = canvas.value.width
  const height = canvas.value.height

  ctx.value.clearRect(0, 0, width, height)

  // Apply transformations
  ctx.value.save()
  ctx.value.translate(width / 2, height / 2)
  ctx.value.rotate(rotation.value.x * Math.PI / 180)
  ctx.value.scale(scale.value, scale.value)
  ctx.value.translate(-width / 2, -height / 2)

  if (analyser && dataArray) {
    analyser.getByteFrequencyData(dataArray)

    switch (currentMode.value) {
      case 'bars':
        drawBars(ctx.value, width, height, dataArray)
        break
      case 'wave':
        drawWave(ctx.value, width, height, dataArray)
        break
      case 'particles':
        drawParticles(ctx.value, width, height, dataArray)
        break
      case 'circular':
        drawCircular(ctx.value, width, height, dataArray)
        break
    }
  }

  ctx.value.restore()
  animationFrame = requestAnimationFrame(draw)
}

const drawBars = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const barWidth = (width / data.length) * 2.5
  let x = 0

  for (let i = 0; i < data.length; i++) {
    const barHeight = (data[i] / 255) * height
    
    const gradient = ctx.createLinearGradient(0, height, 0, height - barHeight)
    gradient.addColorStop(0, '#8B5CF6')
    gradient.addColorStop(1, '#EC4899')
    
    ctx.fillStyle = gradient
    ctx.fillRect(x, height - barHeight, barWidth, barHeight)
    
    x += barWidth + 1
  }
}

const drawWave = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  ctx.beginPath()
  ctx.lineWidth = 2
  ctx.strokeStyle = '#8B5CF6'
  
  const sliceWidth = width / data.length
  let x = 0

  for (let i = 0; i < data.length; i++) {
    const v = data[i] / 255
    const y = v * height

    if (i === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }

    x += sliceWidth
  }

  ctx.stroke()
}

const drawParticles = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const centerX = width / 2
  const centerY = height / 2

  for (let i = 0; i < data.length; i++) {
    const amplitude = data[i] / 255
    const angle = (i / data.length) * Math.PI * 2
    const radius = 50 + amplitude * 100
    
    const x = centerX + Math.cos(angle) * radius
    const y = centerY + Math.sin(angle) * radius
    
    const size = 2 + amplitude * 8
    
    ctx.beginPath()
    ctx.arc(x, y, size, 0, Math.PI * 2)
    ctx.fillStyle = `hsl(${i * 360 / data.length}, 100%, 50%)`
    ctx.fill()
  }
}

const drawCircular = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const centerX = width / 2
  const centerY = height / 2
  const radius = Math.min(width, height) / 4

  ctx.beginPath()
  ctx.arc(centerX, centerY, radius, 0, Math.PI * 2)
  ctx.strokeStyle = '#374151'
  ctx.lineWidth = 2
  ctx.stroke()

  ctx.beginPath()
  for (let i = 0; i < data.length; i++) {
    const amplitude = data[i] / 255
    const angle = (i / data.length) * Math.PI * 2
    const pointRadius = radius + amplitude * 100
    
    const x = centerX + Math.cos(angle) * pointRadius
    const y = centerY + Math.sin(angle) * pointRadius
    
    if (i === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  }
  ctx.closePath()
  
  const gradient = ctx.createRadialGradient(centerX, centerY, radius, centerX, centerY, radius + 100)
  gradient.addColorStop(0, '#8B5CF6')
  gradient.addColorStop(1, '#EC4899')
  
  ctx.fillStyle = gradient
  ctx.fill()
  ctx.strokeStyle = 'white'
  ctx.lineWidth = 1
  ctx.stroke()
}

const startDrag = (event: MouseEvent) => {
  isDragging.value = true
}

const drag = (event: MouseEvent) => {
  if (!isDragging.value) return
  
  rotation.value.x += event.movementX * 0.5
  rotation.value.y += event.movementY * 0.5
}

const stopDrag = () => {
  isDragging.value = false
}

const handleWheel = (event: WheelEvent) => {
  event.preventDefault()
  scale.value += event.deltaY * -0.001
  scale.value = Math.min(Math.max(0.1, scale.value), 3)
}

onMounted(() => {
  if (canvas.value) {
    canvas.value.width = canvas.value.offsetWidth
    canvas.value.height = canvas.value.offsetHeight
    ctx.value = canvas.value.getContext('2d')!
    
    // Simulate audio data for demo
    setupAudioContext()
    draw()
  }
})

onUnmounted(() => {
  if (animationFrame) {
    cancelAnimationFrame(animationFrame)
  }
  if (audioContext) {
    audioContext.close()
  }
})

// In a real implementation, you would connect this to the actual audio element
defineExpose({
  connectAudioElement
})
</script>