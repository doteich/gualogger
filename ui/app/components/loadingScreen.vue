<script setup lang="ts">
import { useSvgAnimation } from '~/composables/useSvgAnimation';

const { pathEl } = useSvgAnimation();
</script>

<template>
    <div>
        <svg height="100px" width="100px" viewBox="0 0 3000 3000" style="transform: scale(-1, 1) rotate(135deg)">
            <path ref="pathEl"
                d="M61 733q0 367 214 554.5t622 187.5h778v-267h-755q-268 0 -396.5 -112.5t-128.5 -362.5q0 -238 113 -351.5t330 -123.5v569h917v-837h-899q-375 0 -585 187.5t-210 555.5zM1145 256h282v305h-282v-305z"
                stroke="currentColor" stroke-width="30" fill="transparent" class="endless-animated-path" />
        </svg>
    </div>
</template>

<style scoped>
div {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--p-background);
    opacity: 0.95;
    z-index: 10;
}

svg {
    margin-top: 1rem;
    color: var( --a-color-1);
}

.endless-animated-path {
    stroke-dasharray: var(--path-length);
    stroke-dashoffset: var(--path-length);
    animation: draw-endless 3s linear infinite;
}

@keyframes draw-endless {
    0% {
        stroke-dashoffset: var(--path-length);
        fill: transparent;
    }

    40% {
        /* Draw in */
        stroke-dashoffset: 0;
        fill: currentColor;
    }

    60% {
        /* Stay filled */
        stroke-dashoffset: 0;
        fill: currentColor;
    }

    100% {
        /* Reset for next cycle */
        stroke-dashoffset: var(--path-length);
        fill: transparent;
    }
}
</style>