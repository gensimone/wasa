<script>
import { fetchUsers } from "@/services/usersService"
import {
    setConversationName,
    setConversationPhotoUrl,
    setConversationId,
    setConversationIsGroup
} from "@/state/conversation"

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
        startConversation(user) {
            setConversationName(user.name)
            setConversationPhotoUrl(user.photoUrl)
            setConversationId(user.userId)
            setConversationIsGroup(false)

            this.$router.push("/conversation")
        },

        async fetchUsers() {
            this.loading = true

            try {
                this.users = await fetchUsers()
            } catch (e) {
                this.error = e?.response?.data?.error || "Unexpected error"
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
            <div class="items-list">
                <div class="list-item" @click="$router.push('/group/create')">
                    <img src="/icons/plus.svg" class="icon-img" />
                    <div class="item-info">
                        <div class="item-name">
                            Create a new group
                        </div>
                    </div>
                </div>
                <UserList :users="users" :loading="loading" @select="startConversation" />
            </div>
        </div>
        <Bottombar />
    </div>
</template>
