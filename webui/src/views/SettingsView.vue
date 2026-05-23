<script>
import { user, setName, setPhotoUrl } from "@/state/user"
import { setMyUserName, setMyPhoto, deleteMyPhoto } from "@/services/users"
import { expandUrl } from "@/utils/media"
import Bottombar from "@/components/Shared/Bottombar.vue"
import Topbar from "@/components/Shared/Topbar.vue"
import InfoSettingCard from "@/components/Settings/InfoSettingCard.vue"
import { handleError } from "@/utils/errors"

export default {
    components: {
        Bottombar,
        Topbar,
        InfoSettingCard
    },

    data() {
        return {
            name: user.name,
            user,

            oldPhotoUrl: null,
            photoUrl: expandUrl(user.photoUrl),
            photo: null,
            photoChanged: false,
            loading: false,
        }
    },

    watch: {
        "user.photoUrl"(newPhotoUrl, _) {
            this.photoUrl = expandUrl(newPhotoUrl)
        },

        "user.name"(newName, _) {
            this.name = newName
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

            event.target.value = ""
        },

        revertPhoto() {
            this.revokePhotoUrl()
            this.photoUrl = this.oldPhotoUrl
            this.photo = null
            this.photoChanged = false
        },

        deletePhoto() {
            this.revokePhotoUrl()
            this.oldPhotoUrl = this.photoUrl
            this.photoUrl = expandUrl(user.defaultUserPhoto)
            this.photo = null
            this.photoChanged = true
        },

        async submit() {
            console.log("submit called")
            this.name = this.name.trim()
            if (!this.name) {
                this.$notifier.error("Provide a name")
                return
            }

            this.loading = true
            let profileChanged = false

            this.name = this.name?.trim()
            if (this.name !== user.name) {
                try {
                    const name = await setMyUserName(this.name)
                    setName(name)

                    profileChanged = true

                } catch (e) {
                    handleError(e)
                    this.loading = false
                    return
                }
            }

            if (this.photo !== null) {
                try {
                    const photoUrl = await setMyPhoto(this.photo)
                    setPhotoUrl(photoUrl)
                    this.revokePhotoUrl()
                    this.photoChanged = false
                    profileChanged = true

                } catch (e) {
                    handleError(e)
                    this.loading = false
                    return
                }

            } else if (this.photoChanged) {
                try {
                    await deleteMyPhoto()
                    setPhotoUrl(user.defaultUserPhoto)
                    this.photoChanged = false
                    profileChanged = true

                } catch (e) {
                    handleError(e)
                    this.loading = false
                    return
                }
            }

            if (profileChanged) {
                this.$notifier.success("Profile updated successfully")
            } else {
                this.$notifier.error("Nothing to do..")
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
        <Topbar :actions="[
            { icon: '/icons/back.svg', onClick: () => $router.back() }
        ]" />
        <div class="settings-page">
            <InfoSettingCard :photoUrl="photoUrl" :photoChanged="photoChanged" :text="name" title="Username"
                submitButtonText="Update" :loading="loading" @revertPhoto="revertPhoto" @deletePhoto="deletePhoto"
                @uploadPhoto="uploadPhoto" @keyPress="handleKeyPress" @submit="submit" />
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

.settings-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
}
</style>
