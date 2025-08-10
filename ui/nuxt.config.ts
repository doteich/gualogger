// https://nuxt.com/docs/api/configuration/nuxt-config

import Aura from '@primeuix/themes/aura'
import { option } from '@primeuix/themes/aura/autocomplete'

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxt/eslint', '@nuxt/icon', '@nuxt/fonts', "@primevue/nuxt-module"],
  css: ['~/assets/main.css', "bootstrap-icons/font/bootstrap-icons.css"],
  primevue: {
    options: {
      theme: {
        preset: Aura,
        options: {
          darkModeSelector: ".dark"
        }
      },

    }
  }
}) 