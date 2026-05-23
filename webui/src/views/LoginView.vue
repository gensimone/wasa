<script>
import { setName, setPhotoUrl, setUserId } from "@/state/user"
import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import LoginBox from "@/components/Login/LoginBox.vue"
import { doLogin } from "@/services/session"
import { handleError } from "@/utils/errors"

export default {
    components: {
        Topbar,
        Bottombar,
        LoginBox
    },

    data() {
        return {
            name: null,
            loading: false
        }
    },

    methods: {
        async doLogin() {
            this.loading = true

            try {
                const user = await doLogin(this.name)

                setName(user.name)
                setPhotoUrl(user.photoUrl)
                setUserId(user.userId)

                this.$router.push("/home")

            } catch (e) {
                handleError(e)

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

        <LoginBox :name="name" :loading="loading" @update:name="name = $event" @submit="doLogin" />

        <Bottombar text="Made by Simone Gentili" />
    </div>
</template>
