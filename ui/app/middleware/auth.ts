import { useOidcAuth } from "~/composables/useOidcAuth"

export default defineNuxtRouteMiddleware(() => {
    const { isLoggedIn } = useOidcAuth();

    if (!isLoggedIn()) {
        return navigateTo('/');
    }
})
