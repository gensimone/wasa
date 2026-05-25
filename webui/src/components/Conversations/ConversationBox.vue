<script>
import ImageModal from "@/components/Shared/ImageModal.vue"
import MessageList from "@/components/Messages/MessageList.vue"
import ConversationInput from "@/components/Conversations/ConversationInput.vue"
import { conversations } from "@/state/conversations"
import { getUserById } from "@/services/users"
import { handleError } from "@/utils/errors"
import { expandUrl } from "@/utils/media"

export default {
    components: { MessageList, ConversationInput, ImageModal },

    data() {
        return {
            scrollTick: 0,
            photoUrl: null,
            name: null,
            zoomedImage: null,
            showImageModal: false,
            isGroup: false
        }
    },

    props: {
        id: { type: Number, required: true },
        routeType: {
            type: String,
            required: true,
            validator: (v) =>
                ["user", "conversation"].includes(v)
        }
    },

    watch: {
        async id() {
            await this.configureConversationData()
        }
    },

    computed: {
        messages() {
            if (this.routeType == "conversation") {
                return conversations.conversationsMap.value.get(this.id)?.messages || []
            } else {
                return []
            }
        },
    },

    methods: {
        expandUrl,

        openImage(url) {
            this.zoomedImage = url
            this.showImageModal = true
        },

        closeImage() {
            this.showImageModal = false
            this.zoomedImage = null
        },

        pushMessage(message) {
            console.log(message)
            this.messages.push(message)
        },

        async configureConversationData() {
            if (this.routeType == "conversation") {
                const conversation = conversations.conversationsMap.value.get(this.id)
                this.photoUrl = conversation?.photoUrl
                this.name = conversation?.name
                this.isGroup = conversation?.isGroup

            } else {
                try {
                    const user = await getUserById(this.id)
                    this.photoUrl = user.photoUrl
                    this.name = user.name
                    this.isGroup = false
                } catch (e) {
                    handleError(e)
                    this.$router.push('/home')
                }
            }
        },

        openGroupInfo() {
            this.$router.push(`/group/${this.id}/info`)
        }
    },

    async mounted() {
        await this.configureConversationData()
    },

    emits: ["reportConversationId"]
}
</script>

<template>
    <div class="conversation-box">

        <div class="conversation-box-header">

            <img class="conversation-box-photo" :src="expandUrl(photoUrl)" @click="openImage(photoUrl)" />

            <div class="conversation-box-name">
                {{ name }}
            </div>

            <div class="conversation-box-info-button">
                <button v-if="isGroup" class="conversation-box-info-btn" @click="openGroupInfo">
                    <img src="/icons/info.svg" class="icon-img">
                </button>
            </div>
        </div>

        <MessageList :messages="messages" :scrollTick="scrollTick" @openImage="openImage" />

        <ConversationInput @reportConversationId="$emit('reportConversationId', $event)" @pushMessage="pushMessage"
            @triggerScrolldown="scrollTick++" :id="id" :routeType="routeType" />

    </div>

    <ImageModal :visible="showImageModal" :imageUrl="zoomedImage" @close="closeImage" />
</template>

<style scoped>
.conversation-box {
    width: min(800px, 75%);
    height: calc(100vh - 150px);
    display: flex;
    flex-direction: column;
    border-radius: 22px;

    background: var(--surface);
    border: 1px solid var(--border);
    box-shadow: 0 25px 90px var(--shadow);
}

.conversation-box-header {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;

    padding: 15px 16px;
    border-bottom: 1px solid var(--border);
    background: transparent;
}

.conversation-box-photo {
    position: absolute;
    left: 16px;

    width: 45px;
    height: 45px;
    border-radius: 50%;
    object-fit: cover;

    border: 1px solid var(--border);
}

.conversation-box-name {
    font-size: 1rem;
    font-weight: 500;
    color: var(--text);
}

.conversation-box-info-button {
    position: absolute;
    right: 16px;
}

.conversation-box-info-btn {
    width: 46px;
    height: 46px;
    border-radius: 100px;
    border: 0px;
    background: rgba(0, 0, 0, 0.0);

    display: flex;
    justify-content: center;
    align-items: center;

    cursor: pointer;
    position: relative;
    overflow: hidden;

    transition: transform 0.25s ease, border 0.25s ease, box-shadow 0.25s ease;
}

.conversation-box-info-btn:hover::before {
    transform: translateX(140%);
}

.conversation-box-info-btn:hover {
    transform: translateY(-4px) scale(1.05);
}

.conversation-box-info-btn:active {
    transform: scale(0.95);
}
</style>
