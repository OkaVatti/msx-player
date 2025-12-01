<template>
  <div class="space-y-6">
    <!-- Control Center Header -->
    <div class="bg-linear-to-r from-gray-800 to-gray-900 rounded-xl p-6 border-2 border-lime-400">
      <h3 class="text-2xl font-bold text-lime-400 mb-2">CONTROL CENTER</h3>
      <p class="text-gray-400">Advanced playback controls and settings</p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Playback Controls -->
      <div class="bg-gray-800 rounded-xl p-6 border-2 border-lime-400">
        <h4 class="text-lg font-semibold text-lime-400 mb-4">PLAYBACK CONTROLS</h4>
        
        <!-- Control Grid -->
        <div class="grid grid-cols-3 gap-4 mb-6">
          <button
            @click="playerStore.previousSong"
            :disabled="!playerStore.currentSong"
            class="p-4 bg-gray-700 rounded-lg border-2 border-lime-400 hover:bg-lime-400 hover:text-black disabled:opacity-50 transition-all flex flex-col items-center justify-center"
          >
            <span class="text-2xl mb-2">⏮</span>
            <span class="text-xs">PREVIOUS</span>
          </button>

          <button
            v-if="!playerStore.isPlaying"
            @click="resumePlayer"
            :disabled="!playerStore.currentSong"
            class="p-4 bg-lime-400 text-black rounded-lg border-2 border-lime-400 hover:bg-lime-300 disabled:opacity-50 transition-all flex flex-col items-center justify-center"
          >
            <span class="text-2xl mb-2">▶</span>
            <span class="text-xs">PLAY</span>
          </button>
          <button
            v-else
            @click="playerStore.pause"
            class="p-4 bg-yellow-400 text-black rounded-lg border-2 border-yellow-400 hover:bg-yellow-300 transition-all flex flex-col items-center justify-center"
          >
            <span class="text-2xl mb-2">⏸</span>
            <span class="text-xs">PAUSE</span>
          </button>

          <button
            @click="playerStore.nextSong"
            :disabled="!playerStore.currentSong"
            class="p-4 bg-gray-700 rounded-lg border-2 border-lime-400 hover:bg-lime-400 hover:text-black disabled:opacity-50 transition-all flex flex-col items-center justify-center"
          >
            <span class="text-2xl mb-2">⏭</span>
            <span class="text-xs">NEXT</span>
          </button>
        </div>

        <!-- Progress -->
        <div class="space-y-3">
          <div class="flex justify-between text-sm text-lime-400">
            <span>{{ formatTime(playerStore.currentTime) }}</span>
            <span>{{ formatTime(playerStore.duration) }}</span>
          </div>
          <div
            class="w-full h-3 bg-gray-700 rounded-full cursor-pointer"
            @click="seekToTime"
            ref="progressBar"
          >
            <div
              class="h-full bg-linear-to-r from-lime-400 to-green-500 rounded-full transition-all duration-500"
              :style="{ width: playerStore.progressPercentage + '%' }"
            ></div>
          </div>
        </div>
      </div>

      <!-- Audio Settings -->
      <div class="bg-gray-800 rounded-xl p-6 border-2 border-purple-400">
        <h4 class="text-lg font-semibold text-purple-400 mb-4">AUDIO SETTINGS</h4>
        
        <!-- Volume -->
        <div class="space-y-3 mb-6">
          <div class="flex justify-between text-sm">
            <span class="text-lime-400">VOLUME</span>
            <span class="text-white">{{ Math.round(playerStore.volume * 100) }}%</span>
          </div>
          <input
            v-model="playerStore.volume"
            @input="setVolume"
            type="range"
            min="0"
            max="1"
            step="0.01"
            class="w-full h-2 bg-gray-700 rounded-lg appearance-none cursor-pointer [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-lime-400"
          />
        </div>

        <!-- Speed -->
        <div class="space-y-3">
          <div class="flex justify-between text-sm">
            <span class="text-lime-400">PLAYBACK SPEED</span>
            <span class="text-white">{{ playerStore.speed }}x</span>
          </div>
          <select
            v-model="playerStore.speed"
            @change="setSpeed"
            class="w-full bg-gray-700 border-2 border-lime-400 text-lime-400 p-3 rounded-lg focus:outline-none focus:bg-lime-400 focus:text-black transition-all"
          >
            <option value="0.5">0.5x (Slow)</option>
            <option value="0.75">0.75x</option>
            <option value="1">1.0x (Normal)</option>
            <option value="1.25">1.25x</option>
            <option value="1.5">1.5x</option>
            <option value="2">2.0x (Fast)</option>
          </select>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="bg-gray-800 rounded-xl p-6 border-2 border-blue-400">
        <h4 class="text-lg font-semibold text-blue-400 mb-4">QUICK ACTIONS</h4>
        <div class="grid grid-cols-2 gap-3">
          <button
            @click="shuffleAll"
            class="p-4 bg-gray-700 rounded-lg border-2 border-blue-400 hover:bg-blue-400 hover:text-black transition-all flex items-center gap-3"
          >
            <span class="text-xl">🔀</span>
            <span class="text-sm">SHUFFLE ALL</span>
          </button>
          <button
            @click="refreshLibrary"
            class="p-4 bg-gray-700 rounded-lg border-2 border-green-400 hover:bg-green-400 hover:text-black transition-all flex items-center gap-3"
          >
            <span class="text-xl">🔄</span>
            <span class="text-sm">REFRESH</span>
          </button>
        </div>
      </div>

      <!-- Current Track Rating -->
      <div v-if="playerStore.currentSong" class="bg-gray-800 rounded-xl p-6 border-2 border-yellow-400">
        <h4 class="text-lg font-semibold text-yellow-400 mb-4">RATE THIS TRACK</h4>
        <div class="space-y-4">
          <div class="flex justify-center space-x-1">
            <button
              v-for="star in 10"
              :key="star"
              @click="rateCurrentSong(star)"
              class="text-3xl transition-all duration-200 hover:scale-125"
              :class="star <= (playerStore.currentSong?.rating || 0) ? 'text-yellow-400' : 'text-gray-600'"
            >
              ★
            </button>
          </div>
          <div class="text-center text-lime-400">
            Current rating: {{ playerStore.currentSong?.rating || 0 }}/10
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { usePlayerStore } from '../../stores/player';
import { useLibraryStore } from '../../stores/library';

const playerStore = usePlayerStore();
const libraryStore = useLibraryStore();
const progressBar = ref<HTMLDivElement>();

const formatTime = (seconds: number) => {
  if (!seconds || isNaN(seconds)) return '0:00';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
};

const seekToTime = (event: MouseEvent) => {
  if (!progressBar.value || !playerStore.duration) return;
  
  const rect = progressBar.value.getBoundingClientRect();
  const percent = (event.clientX - rect.left) / rect.width;
  const newTime = percent * playerStore.duration;
  
  playerStore.seek(newTime);
};

const resumePlayer = () => {
  if (playerStore.currentSong) {
    playerStore.playSong(playerStore.currentSong);
  }
};

const setVolume = () => {
  playerStore.setVolume(playerStore.volume);
};

const setSpeed = () => {
  playerStore.setSpeed(playerStore.speed);
};

const rateCurrentSong = async (rating: number) => {
  if (playerStore.currentSong) {
    await playerStore.rateSong(playerStore.currentSong.id, rating);
    await libraryStore.fetchSongs();
  }
};

const shuffleAll = async () => {
  try {
    await fetch('http://localhost:1323/api/player/shuffle', { method: 'POST' });
    if (libraryStore.songs.length > 0) {
      await playerStore.playSong(libraryStore.songs[0]);
    }
  } catch (error) {
    console.error('Error shuffling songs:', error);
  }
};

const refreshLibrary = () => {
  libraryStore.fetchSongs();
};
</script>