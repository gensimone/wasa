<script>
import MessageItem from "@/components/Messages/MessageItem.vue"

import { user } from "@/state/user"

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
        }
    },

    methods: {
        isNearBottom() {
            const el = this.$refs.container
            return (
                el.scrollHeight - el.scrollTop - el.clientHeight < 80
            )
        },

        // FIXME: It doesn't work.
        scrollToBottomIfNeeded() {
            this.$nextTick(() => {
                const bottom = this.$refs.bottom
                if (!bottom) return

                bottom.scrollIntoView({
                    behavior: "smooth",
                    block: "end"
                })
            })
        },

        scrollToBottomInstant() {
            const bottom = this.$refs.bottom
            if (!bottom) return

            bottom.scrollIntoView({
                behavior: "auto",
                block: "end"
            })
        }
    },

    mounted() {
        this.scrollToBottomInstant()
    }
}
</script>

<template>
    <div class="message-list" ref="container">
        <MessageItem v-for="m in messages" :key="m.messageId" :message="m" />
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
