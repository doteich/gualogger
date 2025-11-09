<template>
  <section class="callback">
    <loading-screen></loading-screen>
  </section>
</template>

<script setup lang="ts">

import { useOidcAuth } from "~/composables/useOidcAuth";
import loadingScreen from "~/components/loadingScreen.vue";


const { manager } = useOidcAuth();
const router = useRouter();

onMounted(() => {
  manager.signinRedirectCallback()
    .then((user) => {
      // Login was successful, redirect to the home page
      router.push('/');
      location.reload()
    })
    .catch((err) => {
      console.error(err);
      // Handle login error, maybe redirect to an error page
      router.push('/');
    });
});
</script>

<style>
.callback {
  width: 100%;
  height: 100%;
  align-items: center;
  justify-content: center;
  display: flex;
  font-size: 2em;

}
</style>