<script>
import { expandUrl } from "@/utils/media"

export default {
    props: {
        conversation: { type: Object, required: true }
    },

    emits: ["select"],

    computed: {
        lastMessage() {
            return this.conversation?.messages?.at?.(-1)
        }
    },

    methods: { expandUrl }
}
</script>

<template>
    <div class="conversation-item" @click="$emit('select', conversation)">

        <div class="conversation-item-photo-wrapper">
            <img :src="expandUrl(conversation.photoUrl)" class="conversation-item-photo" />
        </div>

        <div class="conversation-item-info">
            <div class="conversation-item-name">
                {{ conversation.name }}
            </div>

            <div v-if="lastMessage" class="conversation-item-last-message">
                {{ lastMessage.text }}
            </div>
        </div>

        <div v-if="lastMessage?.attachmentUrl" class="item-photo-wrapper">
            <img :src="expandUrl(lastMessage.attachmentUrl)" class="conversation-item-photo" />
        </div>

    </div>
</template>

<style scoped>
.conversation-item {
    display: flex;
    align-items: center;
    gap: 14px;

    padding: 14px;
    margin-bottom: 10px;

    border-radius: 18px;

    background: var(--surface-2);
    border: 1px solid var(--border);

    cursor: pointer;
    overflow: hidden;

    animation: fadeInUp 0.35s ease both;
}

.conversation-item:hover {
    transform: translateY(-6px) scale(1.02);
    border: 1px solid var(--accent);
}

.conversation-item-photo-wrapper {
    width: 75px;
    height: 75px;
    border-radius: 16px;
    overflow: hidden;
    flex-shrink: 0;
}

.conversation-item-photo {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.conversation-item-info {
    display: flex;
    flex-direction: column;
    flex: 1;
    min-width: 0;
}

.conversation-item-name {
    font-size: 1.05rem;
    font-weight: 800;

    color: var(--text);

    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.conversation-item-last-message {
    margin-top: 4px;
    font-size: 0.92rem;
    color: var(--text-muted);

    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
