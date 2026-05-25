<script>
import { conversations } from "@/state/conversations"

import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import ConversationBox from "@/components/Conversations/ConversationBox.vue"

export default {
    components: { Topbar, Bottombar, ConversationBox },

}
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: '/icons/home.svg', onClick: () => $router.push('/home') }
        ]" />
        <div class="content">
            <ConversationBox :id="id" :routeType="routeType" @reportConversationId="reportConversationId" />
        </div>
        <Bottombar />
    </div>
</template>
