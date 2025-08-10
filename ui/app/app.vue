<script setup lang="ts">
import { useTheme } from './composables/useTheme';
import navbar from '@/components/navbar.vue';

useTheme();
</script>

<template>
  <div id="app">

    <navbar id="nav"></navbar>
    <section>
      <NuxtPage></NuxtPage>
    </section>

  </div>
</template>