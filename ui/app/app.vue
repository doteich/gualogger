<script setup lang="ts">
import { useTheme } from './composables/useTheme';
import { useMainStore } from './stores/mainStore';
import loadingScreen from './components/loadingScreen.vue';
import navbar from '@/components/navbar.vue';

useTheme();

const store = useMainStore()

</script>

<template>
  <div id="app">

    <navbar id="nav"></navbar>
    <section class="main-content">
      <loading-screen v-if="store.loading"></loading-screen>
      <NuxtPage ></NuxtPage>
    </section>

  </div>
</template>

<style>
#app {
  display: flex;
  height: 100vh;
}

.main-content {
  position: relative;
  flex-grow: 1;
  overflow-y: auto;
}
</style>
