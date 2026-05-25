<script>
import MemberList from "@/components/Groups/MemberList.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import Topbar from "@/components/Shared/Topbar.vue"

export default {
    components: { Bottombar, Topbar, MemberList },

    data() {
        return {
            members: []
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
