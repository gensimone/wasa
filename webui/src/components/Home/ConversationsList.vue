<script>
import ConversationItem from "@/components/Home/ConversationItem.vue"

export default {
    components: {
        ConversationItem
    },

    props: {
        conversations: { type: Array, required: true },
        loading: { type: Boolean, required: true }
    },

    emits: ["select"]
}
</script>

<template>
    <div class="items-list">
        <div class="items-list-header">
            <h2> Conversations </h2>
            <button class="icon-btn" @click="$router.push('/conversation/add')">
                <img src="/icons/plus.svg" class="icon-img">
            </button>
        </div>
        <ConversationItem v-for="c in conversations" :key="c.conversationId" :conversation="c"
            @select="$emit('select', $event)" />
    </div>
</template>
