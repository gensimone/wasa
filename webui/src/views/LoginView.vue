<script>
import { setUserId, setName, setPhotoUrl } from "@/state/user"
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
            this.error = null
            this.loading = true

            try {
                let response = await this.$axios.post("/session", {
                    name: this.name
                })

                const userId = response.data.userId
                response = await this.$axios.get(`/users/${userId}`, {
                    headers: { Authorization: userId }
                })

                const data = response.data
                setUserId(data.userId)
                setName(data.name)
                setPhotoUrl(data.photoUrl)

                this.$router.push("/home")

            } catch (e) {
                this.error = e?.response?.data?.error || "Unexpected error"
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
