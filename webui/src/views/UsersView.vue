<script>
import { user } from "@/state/user"
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

        async fetchUser(userId) {
            const response = await this.$axios.get(`/users/${userId}`, {
                headers: { Authorization: user.userId }
            })
            return response.data
        },

        async fetchUsers() {
            this.loading = true
            this.users = []

            try {
                const response = await this.$axios.get(`/users`, {
                    headers: { Authorization: user.userId }
                })

                const userPromises = response.data.users
                    .filter(id => id !== user.userId)
                    .map(id => this.fetchUser(id))

                const usersData = await Promise.all(userPromises)

                this.users = usersData
                    .map(u => ({
                        userId: u.userId,
                        name: u.name,
                        photoUrl: `${__API_URL__}${u.photoUrl}`
                    }))
            } catch (e) {
                this.users = []
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
