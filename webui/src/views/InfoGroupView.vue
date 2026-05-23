<script>
import {
    setConversationName,
    setConversationPhotoUrl,
    setConversationId,
    setConversationIsGroup
} from "@/state/conversation"
import { getMembers } from "@/utils/conversations"
import { conversation } from "@/state/conversation"
import { Poller } from "@/services/poller"
import { removeUser } from "@/services/groups"
import { handleError } from "@/utils/errors"
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
                this.members = this.members.filter(m => m.userId !== member.userId)
                this.$notifier.success(`User ${member.name} removed`)

            } catch (e) {
                handleError(e)
            }
        },

        async startConversation(member) {
            setConversationName(member.name)
            setConversationPhotoUrl(member.photoUrl)
            setConversationId(member.userId)
            setConversationIsGroup(false)

            this.$router.push("/conversation")
        }
    },

    async mounted() {
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
            <MemberList :members="members" @removeUser="removeUser" @selectUser="startConversation" />
        </div>
        <Bottombar />
    </div>
</template>

<style scoped></style>
