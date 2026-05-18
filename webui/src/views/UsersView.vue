<script>
import { fetchUsers } from "@/services/usersService"
import {
    setChatName,
    setChatPhotoUrl,
    setChatUserId
} from "@/state/chat"

import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import UserList from "@/components/Users/UserList.vue"

export default {
    components: {
        Topbar,
        Bottombar,
        UserList
    },

    data() {
        return {
            users: [],
            error: null,
            loading: false
        }
    },

    methods: {
        chatUser(user) {
            setChatName(user.name)
            setChatPhotoUrl(user.photoUrl)
            setChatUserId(user.userId)
            this.$router.push("/chat")
        },

        async fetchUsers() {
            this.loading = true

            try {
                this.users = await fetchUsers()
            } catch (e) {
                this.error = e?.response?.data?.error || e.message || "Unexpected error"
            } finally {
                this.loading = false
            }
        }
    },

    mounted() {
        this.fetchUsers()
    }
}
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: '/icons/back.svg', onClick: () => $router.back() }
        ]" />
        <div class="content">
            <UserList :users="users" :loading="loading" @select="chatUser" />
        </div>
        <Bottombar />
    </div>
</template>
