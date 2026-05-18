<script>
import { user } from "@/state/user"
import { setName, setPhotoUrl } from "@/state/user"

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
            const filename = new URL(this.avatarUrl).pathname.split("/").pop()
            return filename === "default-user-photo.jpg"
        }
    },

    beforeUnmount() {
        if (this.avatarUrl) {
            URL.revokeObjectURL(this.avatarUrl)
        }
    },

    methods: {

        uploadPhoto(event) {
            const file = event.target.files[0]
            if (!file) return

            if (this.avatarUrl) {
                URL.revokeObjectURL(this.avatarUrl)
            }

            this.photo = file
            this.photoChanged = true
            this.avatarUrl = URL.createObjectURL(file)

            this.error = false
            this.message = null
            event.target.value = ""
        },

        async setNewPhoto() {
            this.loading = true
            this.message = null

            try {
                const formData = new FormData()
                formData.append("photo", this.photo)

                const response = await this.$axios.put(
                    `/users/${user.userId}/photo`,
                    formData,
                    {
                        headers: {
                            Authorization: user.userId,
                            "Content-Type": "multipart/form-data"
                        }
                    }
                )

                setPhotoUrl(response.data.photoUrl)

                this.photoChanged = false
                this.error = false
                this.message = "Profile photo updated successfully"
            } catch (e) {
                this.error = true
                this.message = e?.response?.data?.error || "Unexpected error"
            } finally {
                this.loading = false
            }
        },

        async setNewName() {
            if (this.name === user.name) return

            this.loading = true
            this.message = null

            try {
                await this.$axios.put(
                    `/users/${user.userId}/name`,
                    { name: this.name },
                    { headers: { Authorization: user.userId } }
                )

                setName(this.name)

                this.error = false
                this.message = "Name changed successfully"
            } catch (e) {
                this.error = true
                this.message = e?.response?.data?.error || "Unexpected error"
            } finally {
                this.loading = false
            }
        },

        async deleteMyPhoto() {
            this.loading = true

            if (this.photoChanged) {
                if (this.avatarUrl) {
                    URL.revokeObjectURL(this.avatarUrl)
                }

                this.photo = null
                this.avatarUrl = user.photoUrl
                this.photoChanged = false
                this.loading = false
                this.message = null
                return
            }

            try {
                await this.$axios.delete(`/users/${user.userId}/photo`, {
                    headers: { Authorization: user.userId }
                })

                const response = await this.$axios.get(`/users/${user.userId}`, {
                    headers: { Authorization: user.userId }
                })

                setPhotoUrl(response.data.photoUrl)

                if (this.avatarUrl) {
                    URL.revokeObjectURL(this.avatarUrl)
                }

                this.avatarUrl = user.photoUrl
                this.photo = null

                this.error = false
                this.message = "Photo removed successfully"
            } catch (e) {
                this.error = true
                this.message = e?.response?.data?.error || "Unexpected error"
            } finally {
                this.loading = false
            }
        }
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
