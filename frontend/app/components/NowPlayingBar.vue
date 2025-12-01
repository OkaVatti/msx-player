<template>
  <div class="h-20 border-t border-[#837dbd] bg-black flex items-center px-4 gap-4 flex-shrink-10">
    <!-- Song Info -->
    <div class="flex items-center gap-3 min-w-0 w-64">
      <div v-if="playerStore.currentSong" class="w-14 h-14 bg-linear-to-br from-[#5a548d] to-[#837dbd] flex items-center justify-center flex-shrink-10">
        <span class="text-2xl">♪</span>
      </div>
      <div v-else class="w-14 h-14 bg-[#1a1a1a] flex items-center justify-center flex-shrink-10">
        <span class="text-2xl text-[#5a548d]">♪</span>
      </div>
      <div v-if="playerStore.currentSong" class="min-w-0 flex-1">
        <div class="text-[#d3ceff] text-sm font-mono truncate">
          {{ playerStore.currentSong.title }}
        </div>
        <div class="text-[#837dbd] text-xs font-mono truncate">
          {{ playerStore.currentSong.artist }}
        </div>
      </div>
      <div v-else class="text-[#5a548d] text-sm font-mono">
        No song selected
      </div>
    </div>

    <!-- Player Controls -->
    <div class="flex-1 flex flex-col items-center gap-2 max-w-2xl mx-auto">
      <div class="flex items-center gap-4">
        <!-- Shuffle -->
        <button
          @click="playerStore.setShuffle(!playerStore.shuffle)"
          :class="[
            'text-lg transition-colors',
            playerStore.shuffle ? 'text-[#d3ceff]' : 'text-[#837dbd] hover:text-[#d3ceff]'
          ]"
          title="Shuffle"
        >
          🔀
        </button>

        <!-- Previous -->
        <button
          @click="playerStore.previousSong()"
          :disabled="!playerStore.hasPrevious"
          class="text-xl text-[#837dbd] hover:text-[#d3ceff] disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          title="Previous"
        >
          ⏮
        </button>

        <!-- Play/Pause -->
        <button
          v-if="!playerStore.isPlaying"
          @click="handlePlayPause"
          :disabled="!playerStore.currentSong"
          class="w-10 h-10 bg-[#d3ceff] text-black rounded-full flex items-center justify-center hover:bg-white disabled:opacity-30 disabled:cursor-not-allowed transition-all"
          title="Play"
        >
          <span class="ml-0.5">▶</span>
        </button>
        <button
          v-else
          @click="playerStore.pause()"
          class="w-10 h-10 bg-[#d3ceff] text-black rounded-full flex items-center justify-center hover:bg-white transition-all"
          title="Pause"
        >
          ⏸
        </button>

        <!-- Next -->
        <button
          @click="playerStore.nextSong()"
          :disabled="!playerStore.hasNext"
          class="text-xl text-[#837dbd] hover:text-[#d3ceff] disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
          title="Next"
        >
          ⏭
        </button>

        <!-- Repeat -->
        <button
          @click="cycleRepeatMode"
          :class="[
            'text-lg transition-colors',
            playerStore.repeatMode !== 'none' ? 'text-[#d3ceff]' : 'text-[#837dbd] hover:text-[#d3ceff]'
          ]"
          :title="`Repeat: ${playerStore.repeatMode}`"
        >
          {{ repeatIcon }}
        </button>
      </div>

      <!-- Progress Bar -->
      <div class="w-full flex items-center gap-3">
        <span class="text-xs text-[#5a548d] font-mono w-12 text-right">
          {{ formatTime(playerStore.currentTime) }}
        </span>
        <div
          ref="progressBar"
          class="flex-1 h-1 bg-[#5a548d] rounded-full cursor-pointer group relative"
          @click="handleSeek"
          @mouseenter="showProgressHover = true"
          @mouseleave="showProgressHover = false"
          @mousemove="updateHoverTime"
        >
          <div
            class="h-full bg-[#d3ceff] rounded-full transition-all duration-100"
            :style="{ width: playerStore.progressPercentage + '%' }"
          ></div>
          <div
            v-if="showProgressHover"
            class="absolute -top-8 px-2 py-1 bg-black border border-[#837dbd] text-xs text-[#d3ceff] font-mono pointer-events-none"
            :style="{ left: hoverPosition + 'px' }"
          >
            {{ formatTime(hoverTime) }}
          </div>
        </div>
        <span class="text-xs text-[#5a548d] font-mono w-12">
          {{ formatTime(playerStore.duration) }}
        </span>
      </div>
    </div>

    <!-- Volume Control -->
    <div class="flex items-center gap-3 w-40">
      <button
        @click="toggleMute"
        class="text-[#837dbd] hover:text-[#d3ceff] transition-colors"
      >
        {{ volumeIcon }}
      </button>
      <input
        type="range"
        v-model.number="volume"
        @input="handleVolumeChange"
        min="0"
        max="1"
        step="0.01"
        class="flex-1 accent-[#d3ceff]"
      />
      <span class="text-xs text-[#5a548d] font-mono w-8">
        {{ Math.round((volume as any) * 100) }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { usePlayerStore } from '../../stores/player'

const playerStore = usePlayerStore()

const progressBar = ref<HTMLDivElement>()
const showProgressHover = ref(false)
const hoverPosition = ref(0)
const hoverTime = ref(0)

// volume synced to store via computed getter/setter
const volume = computed<number>({
  get: () => playerStore.volume as unknown as number,
  set: (v: number) => {
    playerStore.setVolume(v)
  }
})

const previousVolume = ref(playerStore.volume as unknown as number)

const repeatIcon = computed(() => {
  switch (playerStore.repeatMode) {
    case 'all': return '🔁'
    case 'one': return '🔂'
    default: return '↻'
  }
})

const volumeIcon = computed(() => {
  if ((volume as any) === 0) return '🔇'
  if ((volume as any) < 0.5) return '🔉'
  return '🔊'
})

const formatTime = (seconds: number) => {
  if (!seconds || isNaN(seconds)) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

const handlePlayPause = () => {
  if (playerStore.currentSong) {
    if (playerStore.isPlaying) {
      playerStore.pause()
    } else {
      // Try to resume via server; playerStore.resume() will fallback to local play if necessary
      playerStore.resume()
    }
  }
}

const handleSeek = (event: MouseEvent) => {
  if (!progressBar.value || !playerStore.duration) return
  
  const rect = progressBar.value.getBoundingClientRect()
  const percent = (event.clientX - rect.left) / rect.width
  const newTime = percent * (playerStore.duration as number)
  playerStore.seek(Math.max(0, Math.min(newTime, playerStore.duration as number)))
}

const updateHoverTime = (event: MouseEvent) => {
  if (!progressBar.value || !playerStore.duration) return
  
  const rect = progressBar.value.getBoundingClientRect()
  const percent = Math.max(0, Math.min(1, (event.clientX - rect.left) / rect.width))
  hoverTime.value = percent * (playerStore.duration as number)
  hoverPosition.value = event.clientX - rect.left - 25
}

const handleVolumeChange = () => {
  // the computed setter already calls store.setVolume
}

const toggleMute = () => {
  if ((volume as any) > 0) {
    previousVolume.value = (volume as any)
    volume.value = 0
  } else {
    volume.value = previousVolume.value || 0.7
  }
  playerStore.setVolume(volume as unknown as number)
}

const cycleRepeatMode = () => {
  const modes: Array<'none' | 'all' | 'one'> = ['none', 'all', 'one']
  const currentIndex = modes.indexOf(playerStore.repeatMode as any)
  const nextIndex = (currentIndex + 1) % modes.length
  playerStore.setRepeatMode(modes[nextIndex])
}
</script>
