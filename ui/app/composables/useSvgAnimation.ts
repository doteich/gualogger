import { ref, onMounted, watch } from 'vue';

export function useSvgAnimation() {
  const pathEl = ref<SVGPathElement | null>(null);
  const pathLength = ref(0);

  onMounted(() => {
    if (pathEl.value) {
      pathLength.value = pathEl.value.getTotalLength();
      // Set CSS variables for animation
      pathEl.value.style.setProperty('--path-length', `${pathLength.value}`);
    }
  });

  // Watch for changes in pathLength in case the element is rendered later or changes
  watch(pathLength, (newLength) => {
    if (pathEl.value) {
      pathEl.value.style.setProperty('--path-length', `${newLength}`);
    }
  });

  return {
    pathEl,
    pathLength, // Expose pathLength if needed for other purposes
  };
}

