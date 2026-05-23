<script>
import { user, setName, setPhotoUrl, defaultUserPhotoUrl } from "@/state/user"
import { setMyUserName, setMyPhoto, deleteMyPhoto } from "@/services/users"
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
            expandUrl(user.photoUrl),
            expandUrl(defaultUserPhotoUrl)
        )

        const form = useSettingsForm(user.name, user.name)

        return {
            user,

            ...photo,
            ...form
        }
    },

    watch: {
        "user.name"(v) {
            this.text = v
            this.placeholder = v
        },

        "user.photoUrl"(v) {
            this.photoUrl = expandUrl(v)
            this.photo = null
            this.photoChanged = false
        }
    },

    methods: {
        expandUrl,

        async updateProfile() {
            try {
                await this.submit(async (name) => {
                    let changed = false

                    if (name !== user.name) {
                        const updatedName = await setMyUserName(name)
                        setName(updatedName)
                        changed = true
                    }

                    if (this.photo) {
                        const url = await setMyPhoto(this.photo)
                        setPhotoUrl(url)
                        changed = true
                    } else if (this.photoChanged) {
                        await deleteMyPhoto()
                        setPhotoUrl(defaultUserPhotoUrl)
                        changed = true
                    }

                    if (!changed) throw new Error("NO_CHANGE")
                })

                this.$notifier.success("Profile updated successfully")

            } catch (e) {
                if (e.message === "EMPTY_NAME")
                    this.$notifier.error("Invalid user name")
                else if (e.message === "NO_CHANGE")
                    this.$notifier.error("Nothing to do..")
                else {
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
            <InfoSettingCard :photoUrl="photoUrl" :photoChanged="photoChanged" :text="text" title="Username"
                submitButtonText="Update" :loading="loading" @uploadPhoto="uploadPhoto" @revertPhoto="revertPhoto"
                @deletePhoto="deletePhoto(defaultUserPhotoUrl)" @keyPress="setText" @submit="updateProfile" />
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
