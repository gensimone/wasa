<script>
import {
    groupConversations,
    userConversations,
} from "@/state/conversations"

import ConversationItem from "@/components/Conversations/ConversationItem.vue"

export default {
    components: { ConversationItem },

    emits: ["select"],

    computed: {
        conversationsList() {
            return [
                ...Array.from(conversations.groupConversations.value.values()),
                ...Array.from(conversations.userConversations.value.values())
            ]
        }
    }
}
</script>

<template>
    <div class="conversations-list">

        <div class="conversations-list-header">
            <h2> Conversations </h2>

            <button class="icon-btn" @click="$router.push('/conversation/add')">
                <img src="/icons/plus.svg" class="icon-img">
            </button>
        </div>

        <ConversationItem v-for="c in conversationsList" :key="`${c.isGroup ? 'g' : 'u'}-${c.id}`" :id="c.id"
            :isGroup="c.isGroup" @select="$emit('select', c)" />

    </div>
</template>

<style scoped>
.conversations-list {
    width: min(720px, 100%);
    padding: 20px;
    border-radius: 22px;

    background: var(--surface);
    border: 1px solid var(--border);
    box-shadow: 0 25px 90px var(--shadow);

    backdrop-filter: blur(20px);
}

.conversations-list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 18px;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--border);
}

.conversations-list-header h2 {
    margin: 0;
    font-size: 1.15rem;
    font-weight: 800;
    color: var(--text);
}
</style>
