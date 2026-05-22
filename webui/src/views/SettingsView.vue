<script>
import { user, setName, setPhotoUrl } from "@/state/user"
import { setMyUserName, setMyPhoto, deleteMyPhoto } from "@/services/users"
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
            name: user.name,

            oldPhotoUrl: null,
            photoUrl: expandUrl(user.photoUrl),
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
            this.oldPhotoUrl = expandUrl(user.photoUrl)
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
            this.photoUrl = expandUrl(user.defaultUserPhoto)
            this.photo = null
            this.photoChanged = true

            this.error = false
            this.message = null
        },

        async submit() {
            this.name = this.name.trim()
            if (this.name === "") {
                this.error = true
                this.message = "Provide a name"
                return
            }

            this.loading = true

            let profileChanged = false

            this.name = this.name?.trim()
            if (this.name !== user.name) {
                try {
                    const name = await setMyUserName(this.name)
                    setName(name)

                    this.error = false
                    profileChanged = true

                } catch (e) {
                    this.loading = false
                    this.error = true
                    this.message = e?.response?.data?.error || "Unexpected error"
                    return
                }
            }

            if (this.photo !== null) {
                try {
                    const photoUrl = await setMyPhoto(this.photo)
                    setPhotoUrl(photoUrl)
                    this.revokePhotoUrl()
                    this.photoChanged = false
                    this.error = false
                    profileChanged = true

                } catch (e) {
                    this.loading = false
                    this.error = true
                    this.message = e?.response?.data?.error || "Unexpected error"
                    return
                }
            } else if (this.photoChanged) {
                try {
                    await deleteMyPhoto()
                    setPhotoUrl(user.defaultUserPhoto)
                    this.photoChanged = false
                    this.error = false
                    profileChanged = true

                } catch (e) {
                    this.loading = false
                    this.error = true
                    this.message = e?.response?.data?.error || "Unexpected error"
                    return
                }
            }

            if (profileChanged) {
                this.error = false
                this.message = "Profile updated successfully"
            }

            this.loading = false
        },

        handleKeyPress(name) {
            this.name = name === "" ? user.name : name
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
            <InfoSettingCard :photoUrl="photoUrl" :photoChanged="photoChanged" :text="name" title="Username"
                submitButtonText="Update" :loading="loading" :message="message" :error="error"
                @revertPhoto="revertPhoto" @deletePhoto="deletePhoto" @uploadPhoto="uploadPhoto"
                @keyPress="handleKeyPress" @submit="submit" />
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
