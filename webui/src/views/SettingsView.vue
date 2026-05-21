<script>
import { user, setName, setPhotoUrl } from "@/state/user"
import { updateName, updatePhoto, deletePhoto } from "@/services/userService"
import Bottombar from "@/components/Shared/Bottombar.vue"
import Topbar from "@/components/Shared/Topbar.vue"
import ProfileCard from "@/components/Settings/ProfileCard.vue"

export default {
    components: {
        Bottombar,
        Topbar,
        ProfileCard
    },

    data() {
        return {
            name: user.name,
            avatarUrl: user.photoUrl,
            photo: null,
            message: null,
            error: false,
            photoChanged: false,
            loading: false,
        }
    },

    computed: {
        isDefault() {
            return this.avatarUrl ?
                this.avatarUrl.split("/").pop() === "default-user-photo.jpg"
                : true
        },
    },

    beforeUnmount() {
        this.revokeAvatarUrl()
    },

    methods: {
        revokeAvatarUrl() {
            if (this.avatarUrl) {
                URL.revokeObjectURL(this.avatarUrl)
            }
        },

        uploadPhoto(event) {
            const file = event.target.files[0]
            if (!file) return

            this.revokeAvatarUrl()
            this.avatarUrl = URL.createObjectURL(file)

            this.photoChanged = true

            this.error = false
            this.message = null

            event.target.value = ""
        },

        async setNewPhoto() {
            this.loading = true
            this.message = null

            try {
                const newPhotoUrl = await updatePhoto(this.photo)
                setPhotoUrl(newPhotoUrl)

                this.revokeAvatarUrl()
                this.avatarUrl = user.photoUrl

                this.photoChanged = false

                this.error = false
                this.message = "Profile photo updated successfully"

            } catch (e) {
                this.error = true
                this.message = e?.message || "Unexpected error"

            } finally {
                this.loading = false
            }
        },

        async setNewName() {
            if (this.name === user.name) return

            this.loading = true
            this.message = null

            try {
                const newName = await updateName(this.name)
                setName(newName)

                this.error = false
                this.message = "Name changed successfully"

            } catch (e) {
                this.error = true
                this.message = e?.message || "Unexpected error"

            } finally {
                this.loading = false
            }
        },

        async deleteMyPhoto() {
            this.loading = true
            this.message = null

            try {
                const defaultPhotoUrl = await deletePhoto()
                setPhotoUrl(defaultPhotoUrl)

                this.revokeAvatarUrl()
                this.avatarUrl = user.photoUrl

                this.error = false
                this.message = "Photo removed successfully"

            } catch (e) {
                this.error = true
                this.message = e?.message || "Unexpected error"

            } finally {
                this.loading = false
            }
        },
    }
}
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: '/icons/back.svg', onClick: () => $router.back() },
            { icon: '/icons/save.svg', onClick: () => setNewName() }
        ]" />
        <form class="settings-page" @submit.prevent="setNewName">
            <ProfileCard :avatarUrl="avatarUrl" :name="name" :isDefault="isDefault" :photoChanged="photoChanged"
                :loading="loading" :message="message" :error="error" @uploadPhoto="uploadPhoto"
                @removePhoto="deleteMyPhoto" @savePhoto="setNewPhoto" @updateName="name = $event"
                @saveName="setNewName" />
        </form>
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
