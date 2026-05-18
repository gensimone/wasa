<script>
import { login } from "@/services/authService"
import Topbar from "@/components/Shared/Topbar.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";
import LoginBox from "@/components/Login/LoginBox.vue"

export default {
    components: {
        Topbar,
        Bottombar,
        LoginBox
    },

    data() {
        return {
            name: null,
            error: null,
            loading: false
        }
    },

    methods: {
        async login() {
            if (this.loading) return

            this.loading = true
            this.error = null

            try {
                await login(this.name)
                this.$router.push("/home")
            } catch (e) {
                this.error = e?.response?.data?.error || e.message || "Unexpected error"
            } finally {
                this.loading = false
            }
        }
    }
}
</script>

<template>
    <div class="app">
        <Topbar :links="[
            { icon: '/icons/github.svg', link: 'https://github.com/gensimone' }
        ]" />

        <LoginBox :name="name" :loading="loading" :error="error" @update:name="name = $event" @submit="login" />

        <Bottombar text="Made by Simone Gentili" />
    </div>
</template>
