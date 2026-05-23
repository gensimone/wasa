<script>
import { getMembers } from "@/utils/conversations"
import { conversation } from "@/state/conversation"
import { group } from "@/state/group"
import { user } from "@/state/user"
import { Poller } from "@/services/poller"
import { removeUser } from "@/services/groups"
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
            poller: { default: null, type: Poller }
        }
    },

    methods: {
        async removeUser(member) {
            try {
                await removeUser(conversation.id, member.userId)
                this.$notifier.success(`User ${member.name} removed`)

            } catch (e) {
                this.$notifier.error(e?.response?.data?.error || "Unexpected error")
            }
        }
    },

    async mounted() {
        this.members = await getMembers(conversation.id)
        this.poller = new Poller(async () => {
            this.members = await getMembers(conversation.id)
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
            <MemberList :members="members" @removeUser="removeUser" />
        </div>
        <Bottombar />
    </div>
</template>

<style scoped></style>
