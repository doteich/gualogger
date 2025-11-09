<script setup lang="ts">

import { useOidcAuth } from "~/composables/useOidcAuth";

definePageMeta({
  middleware: ['auth']
})


const { user, login, logout } = useOidcAuth();
</script>

<template>
  <div>
    <h1>Welcome to the Application</h1>
    <div v-if="user">
      <p>Hello, <strong>{{ user.profile.name }}</strong>! 👋</p>
      <pre><code>Access Token: {{ user.access_token.substring(0, 30) }}...</code></pre>
      <button @click="logout()">Sign Out</button>
    </div>
    <div v-else>
      <p>You are not logged in.</p>
      <button @click="login()">Sign In with Keycloak</button>
    </div>
  </div>
</template>
