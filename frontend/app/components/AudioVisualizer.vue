<template>
  <div class="space-y-4">
    <div class="border-2 border-lime-400 p-4 bg-black">
      <h2 class="text-xl mb-4 flex items-center gap-2">
        <span class="animate-pulse">&gt;&gt;</span>
        AUDIO VISUALIZER
        <span class="text-xs text-white">[DRAG TO ROTATE | SCROLL TO ZOOM]</span>
      </h2>

      <!-- Mode Selection -->
      <div class="flex gap-2 mb-4">
        <button
          v-for="mode in visualizerModes"
          :key="mode"
          @click="currentMode = mode"
          :class="[
            'px-3 py-1 border-2 transition-all text-sm',
            currentMode === mode
              ? 'bg-lime-400 text-black border-lime-400'
              : 'bg-black text-lime-400 border-lime-400 hover:bg-lime-400 hover:text-black'
          ]"
        >
          [{{ mode.toUpperCase() }}]
        </button>
      </div>
    </div>

    <!-- Visualizer Canvas -->
    <div class="border-2 border-lime-400 bg-black relative">
      <canvas
        ref="canvas"
        class="w-full h-[500px] cursor-grab active:cursor-grabbing"
        @mousedown="startDrag"
        @mousemove="drag"
        @mouseup="stopDrag"
        @mouseleave="stopDrag"
        @wheel="handleWheel"
      />

      <!-- Stats Overlay -->
      <div class="absolute top-2 left-2 text-xs text-lime-400 bg-black border border-lime-400 p-2 font-mono">
        <div>MODE: {{ currentMode.toUpperCase() }}</div>
        <div>ROT_X: {{ rotation.x.toFixed(2) }}°</div>
        <div>ROT_Y: {{ rotation.y.toFixed(2) }}°</div>
        <div>SCALE: {{ scale.toFixed(2) }}x</div>
        <div>FPS: {{ fps }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

const canvas = ref<HTMLCanvasElement | null>(null);
const ctx = ref<CanvasRenderingContext2D | null>(null);
const isDragging = ref(false);
const rotation = ref({ x: 0, y: 0 });
const scale = ref(1);
const currentMode = ref<'bars' | 'wave' | 'particles' | 'circular'>('bars');
const fps = ref(0);

const visualizerModes: Array<'bars' | 'wave' | 'particles' | 'circular'> = ['bars', 'wave', 'particles', 'circular'];

let animationFrame = 0;
let audioContext: AudioContext | null = null;
let analyser: AnalyserNode | null = null;
let dataArray: Uint8Array | null = null; // <--- plain Uint8Array | null
let lastFrameTime = 0;
let frameCount = 0;
let lastMousePos = { x: 0, y: 0 };

const setupAudioContext = async () => {
  if (!audioContext) {
    // create audio context with backwards compatibility
    audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
    analyser = audioContext.createAnalyser();
    analyser.fftSize = 256;
    // create a concrete Uint8Array backed by an ArrayBuffer
    dataArray = new Uint8Array(analyser.frequencyBinCount);
    // ensure the context is running (some browsers start suspended)
    try { await audioContext.resume(); } catch {}
  }
};

// safe connect function — call with an audio element (from your app)
const connectAudioElement = (element: HTMLAudioElement | null) => {
  if (!element) return;
  if (!audioContext || !analyser || !dataArray) {
    // ensure audio context + analyser exist
    setupAudioContext().catch(() => {});
  }
  if (!audioContext || !analyser || !dataArray) return;

  try {
    // createMediaElementSource may throw if element was already connected previously
    const src = audioContext.createMediaElementSource(element);
    src.connect(analyser);
    analyser.connect(audioContext.destination);
    // resume if suspended
    if (audioContext.state === 'suspended') {
      audioContext.resume().catch(() => {});
    }
  } catch (err) {
    // connecting the same element twice can throw — non-fatal
    console.warn('connectAudioElement warning:', err);
  }
};

const draw = (timestamp: number) => {
  const cvs = canvas.value;
  const c = ctx.value;
  if (!cvs || !c) return;

  // FPS
  frameCount++;
  if (timestamp - lastFrameTime >= 1000) {
    fps.value = frameCount;
    frameCount = 0;
    lastFrameTime = timestamp;
  }

  const width = cvs.width;
  const height = cvs.height;

  // Clear with fade
  c.fillStyle = 'rgba(0, 0, 0, 0.1)';
  c.fillRect(0, 0, width, height);

  // Transform
  c.save();
  c.translate(width / 2, height / 2);
  c.rotate(rotation.value.x * Math.PI / 180);
  c.scale(scale.value, scale.value);
  c.translate(-width / 2, -height / 2);

  // obtain audio data: if analyser + dataArray present use it, otherwise generate deterministic synthetic data
  if (analyser && dataArray) {
    // ensure the AudioContext isn't suspended
    if (audioContext && audioContext.state === 'suspended') {
      audioContext.resume().catch(() => {});
    }
    // fill the dataArray with real analyser values
    analyser.getByteFrequencyData(dataArray);
  } else {
    // fallback: create or reuse a concrete Uint8Array and fill it with deterministic waveform-like values
    if (!dataArray) dataArray = new Uint8Array(128);
    for (let i = 0; i < dataArray.length; i++) {
      // produce a smooth pseudo-wave for demo visuals
      const value = (Math.sin(timestamp / 1000 + i / 10) * 0.5 + 0.5) * 255;
      dataArray[i] = Math.max(0, Math.min(255, Math.floor(value)));
    }
  }

  // Draw selected mode. Cast to non-null since we always ensure dataArray exists above.
  const audioData = dataArray as Uint8Array;

  switch (currentMode.value) {
    case 'bars':
      drawBars(c, width, height, audioData);
      break;
    case 'wave':
      drawWave(c, width, height, audioData);
      break;
    case 'particles':
      drawParticles(c, width, height, audioData, timestamp);
      break;
    case 'circular':
      drawCircular(c, width, height, audioData);
      break;
  }

  c.restore();

  animationFrame = requestAnimationFrame(draw);
};

const drawBars = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const barWidth = (width / data.length) * 2.5;
  let x = 0;

  for (let i = 0; i < data.length; i++) {
    const barHeight = ((data?.[i] ?? 0) / 255) * height;

    ctx.strokeStyle = '#84cc16';
    ctx.fillStyle = '#84cc16';
    ctx.globalAlpha = 0.8;

    ctx.fillRect(x, height - barHeight, barWidth - 1, barHeight);

    ctx.fillStyle = '#a3e635';
    ctx.fillRect(x, height - barHeight, barWidth - 1, 2);

    x += barWidth;
  }
  ctx.globalAlpha = 1;
};

const drawWave = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  ctx.beginPath();
  ctx.lineWidth = 3;
  ctx.strokeStyle = '#84cc16';

  const sliceWidth = width / data.length;
  let x = 0;

  for (let i = 0; i < data.length; i++) {
    const v = (data?.[i] ?? 0) / 255;
    const y = v * height;

    if (i === 0) ctx.moveTo(x, y); else ctx.lineTo(x, y);

    x += sliceWidth;
  }
  ctx.stroke();

  ctx.strokeStyle = '#65a30d';
  ctx.globalAlpha = 0.3;
  ctx.lineWidth = 6;
  ctx.stroke();
  ctx.globalAlpha = 1;
};

const drawParticles = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array, timestamp: number) => {
  const centerX = width / 2;
  const centerY = height / 2;

  for (let i = 0; i < data.length; i += Math.max(1, Math.floor(data.length / 40))) {
    const amplitude = (data?.[i] ?? 0) / 255;
    const angle = (i / data.length) * Math.PI * 2 + timestamp / 1000;
    const radius = 50 + amplitude * 150;

    const x = centerX + Math.cos(angle) * radius;
    const y = centerY + Math.sin(angle) * radius;

    const size = 2 + amplitude * 6;

    ctx.beginPath();
    ctx.arc(x, y, size, 0, Math.PI * 2);
    ctx.fillStyle = '#84cc16';
    ctx.fill();

    ctx.globalAlpha = 0.3;
    ctx.beginPath();
    ctx.arc(x, y, size * 2, 0, Math.PI * 2);
    ctx.fillStyle = '#a3e635';
    ctx.fill();
    ctx.globalAlpha = 1;
  }
};

const drawCircular = (ctx: CanvasRenderingContext2D, width: number, height: number, data: Uint8Array) => {
  const centerX = width / 2;
  const centerY = height / 2;
  const radius = Math.min(width, height) / 4;

  ctx.beginPath();
  ctx.arc(centerX, centerY, radius, 0, Math.PI * 2);
  ctx.strokeStyle = '#374151';
  ctx.lineWidth = 2;
  ctx.stroke();

  ctx.beginPath();
  for (let i = 0; i < data.length; i++) {
    const amplitude = (data?.[i] ?? 0) / 255;
    const angle = (i / data.length) * Math.PI * 2;
    const pointRadius = radius + amplitude * 100;

    const x = centerX + Math.cos(angle) * pointRadius;
    const y = centerY + Math.sin(angle) * pointRadius;

    if (i === 0) ctx.moveTo(x, y); else ctx.lineTo(x, y);
  }
  ctx.closePath();

  ctx.fillStyle = 'rgba(132, 204, 22, 0.3)';
  ctx.fill();
  ctx.strokeStyle = '#84cc16';
  ctx.lineWidth = 2;
  ctx.stroke();
};

const startDrag = (event: MouseEvent) => {
  isDragging.value = true;
  lastMousePos = { x: event.clientX, y: event.clientY };
};

const drag = (event: MouseEvent) => {
  if (!isDragging.value) return;
  rotation.value.x += event.movementY * 0.5;
  rotation.value.y += event.movementX * 0.5;
};

const stopDrag = () => { isDragging.value = false; };

const handleWheel = (event: WheelEvent) => {
  event.preventDefault();
  scale.value += event.deltaY * -0.001;
  scale.value = Math.min(Math.max(0.1, scale.value), 3);
};

onMounted(() => {
  if (!canvas.value) return;
  const rect = canvas.value.getBoundingClientRect();
  canvas.value.width = rect.width;
  canvas.value.height = 500;
  ctx.value = canvas.value.getContext('2d')!;
  // create audio context asap (optional)
  setupAudioContext().catch(() => {});
  animationFrame = requestAnimationFrame(draw);
});

onUnmounted(() => {
  if (animationFrame) cancelAnimationFrame(animationFrame);
  if (audioContext) audioContext.close().catch(() => {});
});
</script>

<style scoped>
/* keep existing styles */
</style>
