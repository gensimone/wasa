<script>
import { getMembers } from "@/utils/conversations"
import { conversation } from "@/state/conversation"
import { group } from "@/state/group"
import { user } from "@/state/user"
import { Poller } from "@/services/poller"
import MemberList from "@/components/GroupInfo/MemberList.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import Topbar from "@/components/Shared/Topbar.vue"

export default {
    components: {
        Bottombar,
        Topbar,
        MemberList
    },

    data() {
        return {
            members: [],
            poller: null
        }
    },

    computed: {
        isFounder() {
            return user.userId === group.founderId
        }
    },

    async mounted() {
        this.members = await getMembers(conversation.id)
        this.poller = new Poller(async () => {
            this.members = await getMembers(conversation.id)
        })
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
            <MemberList :members="members" />
        </div>
        <Bottombar />
    </div>
</template>

<style scoped></style>
