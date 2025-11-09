// https://nuxt.com/docs/api/configuration/nuxt-config

import Aura from '@primeuix/themes/aura'
import { option } from '@primeuix/themes/aura/autocomplete'
import ConfirmPopup from 'primevue/confirmpopup'

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxt/eslint', '@nuxt/icon', '@nuxt/fonts', "@primevue/nuxt-module", '@pinia/nuxt'],
  css: ['~/assets/main.css', "bootstrap-icons/font/bootstrap-icons.css"],

  primevue: {
    directives:{
      include: ["ConfirmPopup"]
    },
    options: {
      theme: {
        preset: Aura,
        options: {
          darkModeSelector: ".dark"
        }
      },

    }
  },
  runtimeConfig: {
    public: {
      clientId: "geist",
      authority: "http://localhost:8080/realms/master",
      redirectUrl: "http://localhost:3000/auth/callback",
      postLogoutUrl: "http://localhost:3000",
      scopes: "openid profile email microprofile-jwt"
    }
  },
  imports: {
    dirs: ["~/composables"]
  }
})