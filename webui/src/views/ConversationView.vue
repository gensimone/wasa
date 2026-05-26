<script>
import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import ConversationBox from "@/components/Conversations/ConversationBox.vue"

export default {
    components: { Topbar, Bottombar, ConversationBox },

    computed: {
        id() {
            return Number(this.$route.params.id);
        },

        direct() {
            return this.$route.query.direct === "true"
        }
    }
}
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: 'back', onClick: () => $router.back() }
        ]" />
        <div class="content">
            <ConversationBox :id="id" :direct="direct" />
        </div>
        <Bottombar />
    </div>
</template>
