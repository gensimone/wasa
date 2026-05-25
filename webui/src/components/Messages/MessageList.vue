<script>
import MessageItem from "@/components/Messages/MessageItem.vue"

import { user } from "@/state/user"
import { setMessageStatusAsRead } from "@/services/messages"
import { handleError } from "@/utils/errors"

export default {
    components: { MessageItem },

    data() {
        return {
            user
        }
    },

    props: {
        messages: { type: Array, required: true },
        scrollTick: { type: Number, required: true }
    },

    watch: {
        scrollTick() {
            this.scrollToBottomIfNeeded()
        },

        async messages(newValue, oldValue) {
            if (this.isNearBottom()) {
                this.scrollToBottomIfNeeded()
            }

            const equal = JSON.stringify(newValue) === JSON.stringify(oldValue)

            if (equal) return

            try {
                const requests = this.messages
                    .filter(m => m.senderId !== user.userId)
                    .map(m => setMessageStatusAsRead(m.messageId))

                await Promise.all(requests)
            } catch (e) {
                handleError(e)
            }
        }
    },

    methods: {
        isNearBottom() {
            const el = this.$refs.container
            return (
                el.scrollHeight - el.scrollTop - el.clientHeight < 80
            )
        },

        scrollToBottomIfNeeded() {
            this.$nextTick(() => {
                const bottom = this.$refs.bottom
                if (!bottom) return

                bottom.scrollIntoView({
                    behavior: "smooth",
                    block: "end"
                })
            })
        }
    },

    mounted() {
        this.scrollToBottomIfNeeded()
    },

    emits: ["openImage"]
}
</script>

<template>
    <div class="message-list" ref="container">
        <MessageItem v-for="message in messages" :key="message.messageId" :message="message"
            :isMine="message.senderId === user.userId" @openImage="$emit('openImage', $event)" />

        <div ref="bottom"></div>
    </div>
</template>

<style scoped>
.message-list {
    height: 100vh;
    overflow-y: auto;
    overflow-x: hidden;
    display: flex;
    flex-direction: column;
    position: relative;
}
</style>
