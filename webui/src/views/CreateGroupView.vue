<script>
import {
    conversation,
    setConversationName,
    setConversationPhotoUrl,
    setConversationId,
    setConversationIsGroup
} from "@/state/conversation"

import {
    setGroupFounderId,
    setGroupCreatedAt
} from "@/state/group"

import { createGroup } from "@/services/groups"
import { expandUrl } from "@/utils/media"
import Bottombar from "@/components/Shared/Bottombar.vue"
import Topbar from "@/components/Shared/Topbar.vue"
import InfoSettingCard from "@/components/Settings/InfoSettingCard.vue"

export default {
    components: {
        Bottombar,
        Topbar,
        InfoSettingCard
    },

    data() {
        return {
            name: null,

            oldPhotoUrl: null,
            photoUrl: expandUrl(conversation.defaultGroupPhoto),
            photo: null,
            photoChanged: false,

            message: null,
            error: false,
            loading: false,
        }
    },

    beforeUnmount() {
        this.revokePhotoUrl()
    },

    methods: {
        expandUrl,

        revokePhotoUrl() {
            if (this.photoUrl) {
                URL.revokeObjectURL(this.photoUrl)
            }
        },

        uploadPhoto(event) {
            const file = event.target.files[0]
            if (!file) return

            this.revokePhotoUrl()
            this.oldPhotoUrl = expandUrl(conversation.defaultGroupPhoto)
            this.photoUrl = URL.createObjectURL(file)

            this.photo = file
            this.photoChanged = true

            this.error = false
            this.message = null

            event.target.value = ""
        },

        revertPhoto() {
            this.revokePhotoUrl()
            this.photoUrl = this.oldPhotoUrl
            this.photo = null
            this.photoChanged = false

            this.error = false
            this.message = null
        },

        deletePhoto() {
            this.revokePhotoUrl()
            this.oldPhotoUrl = this.photoUrl
            this.photoUrl = expandUrl(conversation.defaultGroupPhoto)
            this.photo = null
            this.photoChanged = true

            this.error = false
            this.message = null
        },

        async submit() {
            this.loading = true

            try {
                const group = await createGroup(this.name, this.photo)

                setConversationName(group.name)
                setConversationPhotoUrl(group.photoUrl)
                setConversationId(group.conversationId)
                setConversationIsGroup(true)
                setGroupFounderId(group.founderid)
                setGroupCreatedAt(group.createdAt)

                this.$router.push("/conversation")

            } catch (e) {
                console.log(e)
                this.error = true
                this.message = e?.response?.data?.error || "Unexpected error"
            }

            this.loading = false
        }
    }
}
</script>

<template>
    <div class="app">
        <!-- TOPBAR -->
        <Topbar :actions="[
            { icon: '/icons/back.svg', onClick: () => $router.back() }
        ]" />

        <!-- INFO SETTING CARD -->
        <div class="settings-page">
            <InfoSettingCard :photoUrl="photoUrl" :photoChanged="photoChanged" submitButtonText="Create group"
                :text="name" title="Group name" :loading="loading" :message="message" :error="error"
                @revertPhoto="revertPhoto" @deletePhoto="deletePhoto" @uploadPhoto="uploadPhoto"
                @keyPress="name = $event" @submit="submit" />
        </div>

        <!-- BOTTOMBAR -->
        <Bottombar />

    </div>
</template>

<style scoped>
.settings-page {
    min-height: calc(100vh - 70px);
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 20px;
    position: relative;
    z-index: 1;
}

.settings-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
}
</style>
