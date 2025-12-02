<template>
  <div class="flex items-center gap-2">
    <input v-model="query" @input="onInput" placeholder="Search songs, artists, albums" class="bg-black border border-[#837dbd] p-2 w-64" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useLibraryStore } from '../../../stores/library'
const library = useLibraryStore()
const query = ref(library.filters.search)

const onInput = () => {
  library.filters.search = query.value
}
</script>
