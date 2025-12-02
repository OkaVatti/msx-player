<template>
  <div class="w-full h-full bg-black border border-[#837dbd]">
    <canvas ref="canvas" class="w-full h-full"></canvas>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { usePlayerStore } from '@/stores/player'

const player = usePlayerStore()
const canvas = ref<HTMLCanvasElement | null>(null)
let ctx: CanvasRenderingContext2D | null = null
let audioCtx: AudioContext | null = null
let analyser: AnalyserNode | null = null
let dataArray: Uint8Array | null = null
let source: MediaElementAudioSourceNode | null = null
let raf = 0

function setup() {
  if (!player.audioElement) return
  audioCtx = new (window.AudioContext || (window as any).webkitAudioContext)()
  analyser = audioCtx.createAnalyser()
  analyser.fftSize = 256
  const count = analyser.frequencyBinCount
  // explicit ArrayBuffer-backed typed array
  dataArray = new Uint8Array(new ArrayBuffer(count))
  try {
    source = audioCtx.createMediaElementSource(player.audioElement!)
    source.connect(analyser)
    analyser.connect(audioCtx.destination)
  } catch (err) {
    console.warn('media element source may already be created', err)
  }
}

function draw(time: number) {
  if (!canvas.value || !ctx) return
  const w = canvas.value.width = canvas.value.offsetWidth
  const h = canvas.value.height = canvas.value.offsetHeight
  ctx.fillStyle = 'rgba(0,0,0,0.12)'
  ctx.fillRect(0,0,w,h)
  if (analyser && dataArray) {
    // create a view backed by ArrayBuffer to satisfy TS overloads
    const safe = new Uint8Array(dataArray.buffer as ArrayBuffer, dataArray.byteOffset, dataArray.length)
    analyser.getByteFrequencyData(safe)
    const barCount = 32
    const barWidth = w / barCount
    for (let i=0;i<barCount;i++){
      const idx = Math.floor((i / barCount) * safe.length)
      const amp = safe[idx] / 255
      const bh = amp * h * 0.85
      const x = i*barWidth
      const y = h - bh
      // gradient
      const g = ctx.createLinearGradient(x,y,x,y+bh)
      g.addColorStop(0,'#d3ceff'); g.addColorStop(1,'#5a548d')
      ctx.fillStyle = g
      ctx.fillRect(x+1,y, barWidth-2, bh)
    }
  } else {
    // fallback demo
    ctx.fillStyle = '#837dbd'
    ctx.fillRect(0,h/2, w, 1)
  }
  raf = requestAnimationFrame(draw)
}

onMounted(() => {
  ctx = canvas.value?.getContext('2d') ?? null
  setup()
  raf = requestAnimationFrame(draw)
})

onUnmounted(() => {
  cancelAnimationFrame(raf)
  if (audioCtx) audioCtx.close().catch(()=>{})
})
</script>

<style scoped>
canvas { display:block; width:100%; height:100%; }
</style>
