<template>
  <div class="flex items-center justify-between p-2 border-b border-[#111] hover:bg-[#0e0e0e]">
    <div class="flex items-center gap-3">
      <div class="w-12 h-12 bg-[#111] flex items-center justify-center rounded">
        <span class="text-lg">♪</span>
      </div>
      <div class="min-w-0">
        <div class="text-[#d3ceff] truncate">{{ song.title }}</div>
        <div class="text-xs text-[#837dbd] truncate">{{ song.artist }} • {{ song.album || '—' }}</div>
      </div>
    </div>

    <div class="flex items-center gap-3">
      <div class="text-xs text-[#5a548d]">{{ formatDuration(song.duration) }}</div>
      <button @click="play" class="text-[#837dbd]">Play</button>
      <button @click="addToPlaylist" class="text-[#837dbd]">＋</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Song } from '../../../types'
import { defineProps } from 'vue'
import { usePlayerStore } from '../../../stores/player'
import { useLibraryStore } from '../../../stores/library'

const props = defineProps<{ song: Song }>()
const player = usePlayerStore()
const lib = useLibraryStore()

const play = async () => {
  await player.playSong(props.song)
}

const addToPlaylist = async () => {
  const idStr = prompt('Add to playlist id (use UI to manage playlists). Enter playlist id:')
  if (!idStr) return
  const id = parseInt(idStr, 10)
  if (Number.isNaN(id)) return alert('invalid id')
  // naive API: in future add proper endpoint to add song to playlist
  alert('Server-side add-to-playlist endpoint not implemented in sample. You can extend server to support it.')
}

const formatDuration = (s?: number) => {
  if (!s) return '0:00'
  const m = Math.floor(s/60); const sec = Math.floor(s%60).toString().padStart(2,'0')
  return `${m}:${sec}`
}
</script>
