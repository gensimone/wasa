<script>
import {
    setConversationName,
    setConversationPhotoUrl,
    setConversationId,
    setConversationIsGroup
} from "@/state/conversation"

import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import UserList from "@/components/Users/UserList.vue"
import { getUsers } from "@/services/users";
import { user } from "@/state/user";
import { Poller } from "@/services/poller";

export default {
    components: {
        Topbar,
        Bottombar,
        UserList
    },

    data() {
        return {
            users: [],
            poller: null,
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
        }
    },

    async mounted() {
        this.poller = new Poller(async () => {
            const users = await getUsers()
            this.users = users.filter(u => u.userId != user.userId)
        })

        this.poller.startPolling()
    },

    beforeUnmount() {
        this.poller?.stopPolling()
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
                <UserList :users="users" @select="startConversation" />
            </div>
        </div>
        <Bottombar />
    </div>
</template>
