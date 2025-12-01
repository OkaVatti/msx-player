<template>
  <div class="p-4 border-2 border-red-400">
    <h2 class="text-xl mb-4">TEST COMPONENT - IT WORKS!</h2>
    <p>Current time: {{ new Date().toLocaleTimeString() }}</p>
    <p>Store test: {{ storeTest }}</p>
    <button @click="testClick" class="px-4 py-2 border-2 border-lime-400 mt-4">
      CLICK ME
    </button>
    <p v-if="clicked" class="text-green-400 mt-2">✅ Button was clicked!</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const clicked = ref(false);
const storeTest = ref('Store not tested (but component works!)');

const testClick = () => {
  clicked.value = true;
  console.log('✅ Test component button clicked!');
};
</script>