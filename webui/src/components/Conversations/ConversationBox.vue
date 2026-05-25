<script>
import ImageModal from "@/components/Shared/ImageModal.vue"
import MessageList from "@/components/Messages/MessageList.vue"
import ConversationInput from "@/components/Conversations/ConversationInput.vue"
import { userConversations, groupConversations } from "@/state/conversations"
import { expandUrl } from "@/utils/media"
import { getUserById } from "@/services/users"
import { handleError } from "@/utils/errors"

export default {
    components: { MessageList, ConversationInput, ImageModal },

    data() {
        return {
            scrollTick: 0,
            zoomedImage: null,
            showImageModal: false,

            userFetchPromise: null,
            userFetched: null
        }
    },

    props: {
        id: { type: Number, required: true },
        direct: { type: Boolean, required: true }
    },

    computed: {
        conversationData() {
            return this.direct
                ? userConversations.value.get(this.id)
                : groupConversations.value.get(this.id)
        },

        messages() {
            return this.conversationData?.messages || []
        },

        name() {
            console.log(this.conversationData?.name)
            return this.conversationData?.name
                || this.userFetched?.name
                || ""
        },

        photoUrl() {
            return this.conversationData?.photoUrl
                || this.userFetched?.photoUrl
                || ""
        }
    },

    watch: {
        conversationData: {
            immediate: true,
            handler(val) {
                if (!val?.name || !val?.photoUrl) {
                    this.ensureUserData();
                }
            }
        }
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

        async ensureUserData() {
            if (this.userFetched) return

            try {
                this.userFetched = await getUserById(this.id)
            } catch (e) {
                handleError(e)
                this.$router.push('/home')
            }
        }
    }
}
</script>

<template>
    <div class="conversation-box">

        <div class="conversation-box-header">
            <img class="conversation-box-photo" :src="expandUrl(photoUrl)" />
            <div class="conversation-box-name">
                {{ name }}
            </div>
            <div class="conversation-box-info-button">
                <button v-if="!direct" class="conversation-box-info-btn" @click="openGroupInfo">
                    <img src="/icons/info.svg" class="icon-img">
                </button>
            </div>
        </div>

        <MessageList :messages="messages" :scrollTick="scrollTick" @openImage="openImage" />

        <ConversationInput @triggerScrolldown="scrollTick++" :direct="direct" :id="id" />
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
