<script>
import { createGroup } from "@/services/groups"
import { expandUrl } from "@/utils/media"
import { handleError } from "@/utils/errors"

import { usePhotoManager } from "@/composables/usePhotoManager"
import { useSettingsForm } from "@/composables/useSettingsForm"

import SettingsCard from "@/components/Settings/SettingsCard.vue"
import Topbar from "@/components/Shared/Topbar.vue"
import Bottombar from "@/components/Shared/Bottombar.vue"

export default {
    components: { SettingsCard, Topbar, Bottombar },

    data() {
        const photo = usePhotoManager(
            expandUrl(defaultGroupPhotoUrl),
            expandUrl(defaultGroupPhotoUrl)
        )

        const form = useSettingsForm("...", "...")

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
                    this.$router.push(`/conversation/${group.conversationId}`)
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
        <div class="content-center">
            <SettingsCard :photoUrl="photoUrl" :photoChanged="photoChanged" :text="text" title="Group name"
                submitButtonText="Create group" :loading="loading" @uploadPhoto="uploadPhoto" @revertPhoto="revertPhoto"
                @deletePhoto="deletePhoto" @keyPress="setText" @submit="createGroup" />

        </div>
        <Bottombar />
    </div>
</template>
