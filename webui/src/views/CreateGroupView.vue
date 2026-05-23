<script>
import {
    setConversationName,
    setConversationPhotoUrl,
    setConversationId,
    setConversationIsGroup
} from "@/state/conversation"

import {
    setGroupFounderId,
    setGroupCreatedAt,
    defaultGroupPhotoUrl
} from "@/state/group"

import { createGroup } from "@/services/groups"
import { expandUrl } from "@/utils/media"
import { handleError } from "@/utils/errors"

import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"
import InfoSettingCard from "@/components/Settings/InfoSettingCard.vue"

import { usePhotoManager } from "@/composables/usePhotoManager"
import { useSettingsForm } from "@/composables/useSettingsForm"

export default {
    components: { Topbar, Bottombar, InfoSettingCard },

    data() {
        const photo = usePhotoManager(
            expandUrl(defaultGroupPhotoUrl),
            expandUrl(defaultGroupPhotoUrl)
        )

        const form = useSettingsForm("", "...")

        return {
            ...photo,
            ...form
        }
    },

    methods: {
        expandUrl,

        async createGroup() {
            try {
                await this.submit(async (name) => {
                    const group = await createGroup(name, this.photo)

                    setConversationName(group.name)
                    setConversationPhotoUrl(group.photoUrl)
                    setConversationId(group.conversationId)
                    setConversationIsGroup(true)
                    setGroupFounderId(group.founderid)
                    setGroupCreatedAt(group.createdAt)

                    this.$router.push("/conversation")
                })

            } catch (e) {
                if (e.message === "EMPTY_NAME") {
                    this.$notifier.error("Invalid group name")
                } else {
                    handleError(e)
                }
            }
        }
    }
}
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: '/icons/back.svg', onClick: () => $router.back() }
        ]" />

        <div class="settings-page">
            <InfoSettingCard :photoUrl="photoUrl" :photoChanged="photoChanged" :text="text" title="Group name"
                submitButtonText="Create group" :loading="loading" @uploadPhoto="uploadPhoto" @revertPhoto="revertPhoto"
                @deletePhoto="deletePhoto" @keyPress="setText" @submit="createGroup" />
        </div>

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
</style>
