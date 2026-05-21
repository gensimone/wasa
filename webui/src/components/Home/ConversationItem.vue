<script>
import { expandURL } from "@/utils/media"
export default {
    props: {
        conversation: Object
    },

    emits: ["select"],

    methods: {
        expandURL
    }
}
</script>

<template>
    <div class="list-item" @click="$emit('select', conversation)">
        <div class="item-photo-wrapper">
            <img :src="expandURL(conversation.photoUrl)" class="item-photo" />
        </div>
        <div class="item-info">
            <div class="item-name">
                {{ conversation.name }}
            </div>
            <div v-if="conversation.lastMessage" class="last-message">
                {{ conversation.lastMessage.text }}
            </div>
        </div>
    </div>
</template>
