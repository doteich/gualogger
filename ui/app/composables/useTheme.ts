import { ref, onMounted, watch } from 'vue'

export function useTheme() {
  const is_light = ref(true)

  function toggleTheme() {
    is_light.value = !is_light.value
  }

  watch(is_light, (newTheme) => {

    if (!newTheme) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  })

  onMounted(() => {
    if (is_light.value) {
      document.documentElement.classList.remove('dark')
    }
  })

  return {
    is_light,
    toggleTheme
  }
}
