import { defineStore } from "pinia"

export const useMainStore = defineStore('mainStore', {
    state: () => ({
        isLoading: false
    }),

    getters: {
        loading(state) {
            return state.isLoading
        }
    },

    actions: {
        setLoadingState(status: boolean) {
            this.isLoading = status
        }

    }
})