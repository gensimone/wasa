<script>
import { user, setName, setPhotoUrl } from "@/state/user"
import { updateName, updatePhoto, deletePhoto } from "@/services/userService"
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
            photoUrl: expandUrl("/media/default-group-photo.jpg"),
            photo: null,
            message: null,
            error: false,
            loading: false
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
            this.photoUrl = URL.createObjectURL(file)

            this.photo = file

            this.error = false
            this.message = null

            event.target.value = ""
        },

        deletePhoto() {
            this.revokePhotoUrl()
            this.photoUrl = expandUrl("/media/default-group-photo.jpg")
        },

        createGroup() {
            console.log("Creating group with:")
            console.log(this.name)
            console.log(this.photoUrl)
        }
    }
}
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: '/icons/back.svg', onClick: () => $router.back() }
        ]" />
        <form class="settings-page" @submit.prevent="updateName">
            <InfoSettingCard :photoUrl="photoUrl" :uploadButton="false" :text="name" title="Group name"
                :loading="loading" :message="message" :error="error" @deletePhoto="deletePhoto"
                @uploadPhoto="uploadPhoto" @keyPress="name = $event" @updateText="createGroup" />
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
