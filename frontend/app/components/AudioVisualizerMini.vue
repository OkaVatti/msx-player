<template>
  <div class="w-full h-full relative bg-black border-2 border-lime-400 overflow-hidden">
    <!-- Header -->
    <div class="absolute top-0 left-0 right-0 bg-black border-b-2 border-lime-400 p-1 z-10">
      <div class="flex justify-between items-center text-xs">
        <span class="text-lime-400 font-mono">[MINI-VIS]</span>
        <span class="text-white/60">v1.337</span>
        <button 
          @click="$emit('close')"
          class="text-lime-400 hover:text-white transition-colors"
        >
          [X]
        </button>
      </div>
    </div>

    <!-- Canvas -->
    <canvas
      ref="canvas"
      class="w-full h-full cursor-crosshair"
      @mousedown="startDrag"
      @mousemove="drag"
      @mouseup="stopDrag"
      @mouseleave="stopDrag"
      @wheel="handleWheel"
    />

    <!-- Mode Indicator -->
    <div class="absolute bottom-2 left-2 bg-black border border-lime-400 px-2 py-1">
      <span class="text-xs text-lime-400 font-mono">MODE:{{ currentMode.toUpperCase() }}</span>
    </div>

    <!-- Stats -->
    <div class="absolute top-8 left-2 bg-black border border-lime-400 px-2 py-1 text-xs font-mono">
      <div class="text-lime-400">ROT:{{ rotation.x.toFixed(0) }},{{ rotation.y.toFixed(0) }}</div>
      <div class="text-green-400">SCL:{{ scale.toFixed(1) }}x</div>
    </div>

    <!-- Mode Switch Button -->
    <button
      @click="cycleMode"
      class="absolute bottom-2 right-2 bg-black border border-lime-400 px-2 py-1 hover:bg-lime-400 hover:text-black transition-all"
    >
      <span class="text-xs text-lime-400 hover:text-black font-mono">[SWITCH]</span>
    </button>

    <!-- Audio Input Indicator -->
    <div class="absolute top-2 right-2">
      <div class="flex items-center gap-1">
        <div 
          v-for="n in 3" 
          :key="n"
          class="w-1 bg-lime-400 transition-all duration-100"
          :style="{ height: `${Math.random() * 8 + 2}px` }"
          :class="{ 'animate-pulse': isAudioPlaying }"
        ></div>
      </div>
    </div>

    <!-- Scanlines Overlay -->
    <div class="absolute inset-0 pointer-events-none scanlines"></div>

    <!-- CRT Glow -->
    <div class="absolute inset-0 pointer-events-none crt-glow"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

const emit = defineEmits(['close']);

const canvas = ref<HTMLCanvasElement>();
const ctx = ref<CanvasRenderingContext2D>();
const isDragging = ref(false);
const rotation = ref({ x: 0, y: 0 });
const scale = ref(1);
const currentMode = ref<'bars' | 'wave' | 'particles' | 'circular' | 'matrix'>('bars');
const isAudioPlaying = ref(false);

let animationFrame: number;
let audioContext: AudioContext;
let analyser: AnalyserNode;
let dataArray: Uint8Array<ArrayBuffer>;
let lastMousePos = { x: 0, y: 0 };

const modes: typeof currentMode.value[] = ['bars', 'wave', 'particles', 'circular', 'matrix'];

const setupAudioContext = () => {
  if (!audioContext) {
    audioContext = new AudioContext();
    analyser = audioContext.createAnalyser();
    analyser.fftSize = 128;
    dataArray = new Uint8Array(analyser.frequencyBinCount) as Uint8Array<ArrayBuffer>;
  }
};

const connectAudioElement = (element: HTMLAudioElement) => {
  setupAudioContext();
  const source = audioContext.createMediaElementSource(element);
  source.connect(analyser);
  analyser.connect(audioContext.destination);
  isAudioPlaying.value = true;
};

const draw = () => {
  if (!canvas.value || !ctx.value) return;

  const width = canvas.value.width;
  const height = canvas.value.height;

  // Clear with retro fade effect
  ctx.value.fillStyle = 'rgba(0, 0, 0, 0.15)';
  ctx.value.fillRect(0, 0, width, height);

  // Apply transformations
  ctx.value.save();
  ctx.value.translate(width / 2, height / 2);
  ctx.value.rotate(rotation.value.x * Math.PI / 180);
  ctx.value.scale(scale.value, scale.value);
  ctx.value.translate(-width / 2, -height / 2);

  // Get or generate audio data
  let audioData = dataArray;
  if (!audioData) {
    audioData = new Uint8Array(64);
  }

  if (analyser && audioData) {
    analyser.getByteFrequencyData(audioData);
  } else {
    // Demo data
    for (let i = 0; i < audioData.length; i++) {
      audioData[i] = Math.random() * 255 * Math.sin(Date.now() / 1000 + i / 10);
    }
  }

  switch (currentMode.value) {
    case 'bars':
      drawMiniBars(ctx.value, width, height, audioData);
      break;
    case 'wave':
      drawMiniWave(ctx.value, width, height, audioData);
      break;
    case 'particles':
      drawMiniParticles(ctx.value, width, height, audioData);
      break;
    case 'circular':
      drawMiniCircular(ctx.value, width, height, audioData);
      break;
    case 'matrix':
      drawMatrix(ctx.value, width, height);
      break;
  }

  ctx.value.restore();
  animationFrame = requestAnimationFrame(draw);
};

const drawMiniBars = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  if (!data || data.length === 0) return;
  
  const barCount = 16;
  const barWidth = width / barCount;
  
  for (let i = 0; i < barCount; i++) {
    const dataIndex = Math.floor((i / barCount) * data.length);
    const amplitude = data[dataIndex]! / 255;
    const barHeight = amplitude * height * 0.8;
    
    const x = i * barWidth;
    const y = height - barHeight;
    
    // Retro green gradient
    const gradient = ctx.createLinearGradient(x, y, x, y + barHeight);
    gradient.addColorStop(0, '#a3e635');
    gradient.addColorStop(0.5, '#84cc16');
    gradient.addColorStop(1, '#65a30d');
    
    ctx.fillStyle = gradient;
    ctx.fillRect(x + 1, y, barWidth - 2, barHeight);
    
    // Scan line effect
    ctx.fillStyle = 'rgba(163, 230, 53, 0.6)';
    ctx.fillRect(x + 1, y, barWidth - 2, 1);
  }
};

const drawMiniWave = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  ctx.beginPath();
  ctx.lineWidth = 2;
  ctx.strokeStyle = '#84cc16';
  
  const sliceWidth = width / data.length;
  let x = 0;

  for (let i = 0; i < data.length; i++) {
    const v = (data[i] ?? 0) / 255;
    const y = (1 - v) * height;

    if (i === 0) {
      ctx.moveTo(x, y);
    } else {
      ctx.lineTo(x, y);
    }

    x += sliceWidth;
  }

  ctx.stroke();
  
  // Glow effect
  ctx.strokeStyle = 'rgba(163, 230, 53, 0.3)';
  ctx.lineWidth = 4;
  ctx.stroke();
};

const drawMiniParticles = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const time = Date.now() / 1000;

  for (let i = 0; i < 20; i++) {
    const dataIndex = i % data.length;
    const amplitude = (data[dataIndex] ?? 0) / 255;
    const angle = (i / 20) * Math.PI * 2 + time;
    const radius = 20 + amplitude * 40;
    
    const x = centerX + Math.cos(angle) * radius;
    const y = centerY + Math.sin(angle) * radius;
    
    const size = 1 + amplitude * 3;
    
    ctx.beginPath();
    ctx.arc(x, y, size, 0, Math.PI * 2);
    ctx.fillStyle = `hsl(${i * 18}, 100%, 60%)`;
    ctx.fill();
    
    // Trail effect
    ctx.globalAlpha = 0.3;
    ctx.beginPath();
    ctx.arc(x, y, size * 2, 0, Math.PI * 2);
    ctx.fillStyle = '#84cc16';
    ctx.fill();
    ctx.globalAlpha = 1;
  }
};

const drawMiniCircular = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const baseRadius = Math.min(width, height) / 6;

  ctx.beginPath();
  ctx.arc(centerX, centerY, baseRadius, 0, Math.PI * 2);
  ctx.strokeStyle = '#374151';
  ctx.lineWidth = 1;
  ctx.stroke();

  // Frequency response
  ctx.beginPath();
  for (let i = 0; i < data.length; i++) {
    const amplitude = (data[i] ?? 0) / 255;
    const angle = (i / data.length) * Math.PI * 2;
    const radius = baseRadius + amplitude * baseRadius;
    
    const x = centerX + Math.cos(angle) * radius;
    const y = centerY + Math.sin(angle) * radius;
    
    if (i === 0) {
      ctx.moveTo(x, y);
    } else {
      ctx.lineTo(x, y);
    }
  }
  ctx.closePath();
  
  ctx.fillStyle = 'rgba(132, 204, 22, 0.2)';
  ctx.fill();
  ctx.strokeStyle = '#84cc16';
  ctx.lineWidth = 1;
  ctx.stroke();
};

const drawMatrix = (ctx: CanvasRenderingContext2D, width: number, height: number) => {
  // Create matrix rain effect
  const fontSize = 10;
  const columns = Math.floor(width / fontSize);
  const drops: number[] = Array(columns).fill(1);
  
  ctx.fillStyle = 'rgba(0, 0, 0, 0.04)';
  ctx.fillRect(0, 0, width, height);
  
  ctx.fillStyle = '#84cc16';
  ctx.font = `${fontSize}px monospace`;
  
  for (let i = 0; i < drops.length; i++) {
    const text = String.fromCharCode(0x30A0 + Math.random() * 96);
    const x = i * fontSize;
    const y = (drops[i] ?? 1) * fontSize;
    
    ctx.fillText(text, x, y);
    
    if (y > height && Math.random() > 0.975) {
      drops[i] = 0;
    }
    if (typeof drops[i] === 'number') {
      drops[i] = ((drops[i] ?? 0) + 1);
    }
  }
};

const startDrag = (event: MouseEvent) => {
  isDragging.value = true;
  lastMousePos = { x: event.clientX, y: event.clientY };
};

const drag = (event: MouseEvent) => {
  if (!isDragging.value) return;
  
  const deltaX = event.clientX - lastMousePos.x;
  const deltaY = event.clientY - lastMousePos.y;
  
  rotation.value.y += deltaX * 0.5;
  rotation.value.x += deltaY * 0.5;
  
  lastMousePos = { x: event.clientX, y: event.clientY };
};

const stopDrag = () => {
  isDragging.value = false;
};

const handleWheel = (event: WheelEvent) => {
  event.preventDefault();
  scale.value += event.deltaY * -0.002;
  scale.value = Math.min(Math.max(0.3, scale.value), 2);
};

const cycleMode = () => {
  const currentIndex = modes.indexOf(currentMode.value);
  const nextIndex = (currentIndex + 1) % modes.length;
  currentMode.value = modes[nextIndex] ?? 'bars';
};

onMounted(() => {
  if (canvas.value) {
    const rect = canvas.value.getBoundingClientRect();
    canvas.value.width = rect.width;
    canvas.value.height = rect.height;
    ctx.value = canvas.value.getContext('2d')!;
    
    setupAudioContext();
    draw();
  }
});

onUnmounted(() => {
  if (animationFrame) {
    cancelAnimationFrame(animationFrame);
  }
  if (audioContext) {
    audioContext.close();
  }
});

// Expose method to connect to audio element
defineExpose({
  connectAudioElement
});
</script>

<style scoped>
.scanlines {
  background: linear-gradient(
    to bottom,
    transparent 50%,
    rgba(0, 0, 0, 0.1) 51%
  );
  background-size: 100% 4px;
  animation: scanline 8s linear infinite;
}

.crt-glow {
  box-shadow: 
    inset 0 0 20px rgba(132, 204, 22, 0.1),
    0 0 30px rgba(132, 204, 22, 0.1);
  animation: flicker 0.15s infinite alternate;
}

@keyframes scanline {
  0% {
    background-position: 0 0;
  }
  100% {
    background-position: 0 100%;
  }
}

@keyframes flicker {
  0%, 19%, 21%, 23%, 25%, 54%, 56%, 100% {
    opacity: 0.1;
  }
  20%, 24%, 55% {
    opacity: 0.2;
  }
}
</style>