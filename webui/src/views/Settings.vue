<script>
import { setName, setPhotoUrl, user } from "@/state/user"
import backIcon from "@/assets/icons/back.svg"
import saveIcon from "@/assets/icons/save.svg"
import deleteIcon from "@/assets/icons/delete.svg"
import uploadIcon from "@/assets/icons/upload.svg"
import LoggedAs from "@/components/LoggedAs.vue"

export default {
    components: {
        LoggedAs
    },

    data() {
        return {
            name: user.name,          // Name of the user.
            avatarUrl: user.photoUrl, // UI Avatar.
            photo: null,              // Photo to upload.

            message: null,
            error: false,        // Change status area messages color.
            photoChanged: false, // Enable/Disable upload button.
            loading: false,      // Enable/Disable upload button during uploading.

            backIcon,
            saveIcon,
            deleteIcon,
            uploadIcon,
        }
    },

    computed: {
        isDefault() {
            const filename = new URL(this.avatarUrl).pathname.split("/").pop()
            return filename == "default-user-photo.jpg"
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
            }

            this.loading = false
        },

        async setNewName() {
            this.loading = true
            this.message = null

            if (this.name !== user.name) {
                try {
                    await this.$axios.put(`/users/${user.userId}/name`,
                        { name: this.name },
                        { headers: { Authorization: user.userId } }
                    )

                    setName(this.name)
                    this.error = false
                    this.message = "Name changed successfully"
                } catch (e) {
                    this.error = true
                    this.message = e?.response?.data?.error || "Unexpected error"
                }
            }

            this.loading = false
        },

        async deleteMyPhoto() {
            this.loading = true

            // Remove local uploaded photo.
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

            // Remove remote uploaded photo.
            try {
                await this.$axios.delete(`/users/${user.userId}/photo`,
                    { headers: { Authorization: user.userId } }
                )

                try {
                    const response = await this.$axios.get(`/users/${user.userId}`, {
                        headers: { Authorization: user.userId }
                    })

                    setPhotoUrl(response.data.photoUrl)

                    if (this.avatarUrl) {
                        URL.revokeObjectURL(this.avatarUrl)
                    }

                    this.avatarUrl = userState.photoUrl
                    this.photo = null
                } catch (e) {
                    this.error = true
                    this.message = e?.response?.data?.error || "Unexpected error"
                    return
                }

                this.error = false
                this.message = "Photo removed successfully"
            } catch (e) {
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
        <header class="topbar">
            <div class="header-title"> WASAText </div>
            <div class="actions">
                <!-- BACK BUTTON -->
                <button class="icon-btn" @click="$router.back()">
                    <img :src="backIcon" class="icon-img">
                </button>
                <!-- SET NEW NAME -->
                <button class="icon-btn" :disabled="loading" @click="setNewName">
                    <img :src="saveIcon" class="icon-img">
                </button>
            </div>
        </header>
        <form class="settings-page" @submit.prevent="setNewName">
            <div class="settings-container">
                <!-- PROFILE CARD -->
                <div class="profile-card">
                    <div class="avatar-row">
                        <!-- LEFT BUTTON -->
                        <button type="button" class="icon-btn" :disabled="isDefault" @click="deleteMyPhoto">
                            <img :src="deleteIcon" class="icon-img">
                        </button>
                        <!-- AVATAR -->
                        <div class="avatar-wrapper">
                            <label class="avatar-clickable">
                                <img :src="avatarUrl" class="avatar-big" />
                                <input type="file" accept="image/*" @change="uploadPhoto" hidden>
                            </label>
                        </div>
                        <!-- RIGHT BUTTON -->
                        <button type="button" class="icon-btn" :disabled="!photoChanged" @click="setNewPhoto">
                            <img :src="uploadIcon" class="icon-img">
                        </button>
                    </div>
                    <!-- USERNAME -->
                    <div class="username-box">
                        <h2>Username</h2>
                        <input v-model="name" class="prompt" type="text" placeholder="Enter username" />
                    </div>
                </div>
                <!-- STATUS MESSAGES -->
                <div class="status-area">
                    <p v-if="message" :class="error ? 'error' : 'success'">
                        {{ message }}
                    </p>
                </div>
            </div>
        </form>
        <LoggedAs />
    </div>
</template>

<style scoped>
.settings-page,
profile-card {
    position: relative;
    z-index: 1;
}

.settings-page {
    min-height: calc(100vh - 70px);
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 20px;
}

.settings-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
}

.profile-card {
    width: min(520px, 92%);
    padding: 34px;

    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 30px;

    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.06);
    border-radius: 22px;

    backdrop-filter: blur(20px);

    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.75);
}

.avatar-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
}

.avatar-big {
    width: 170px;
    height: 170px;

    border-radius: 50%;
    object-fit: cover;

    border: 2px solid rgba(0, 255, 120, 0.18);
    box-shadow: 0 0 50px rgba(0, 255, 120, 0.08);

    transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.avatar-big:hover {
    transform: scale(1.04);
    box-shadow: 0 0 70px rgba(0, 255, 120, 0.15);
}

/* clickable upload button */
.avatar-edit {
    font-size: 0.85rem;
    letter-spacing: 1px;
    text-transform: uppercase;

    color: rgba(0, 255, 120, 0.7);
    cursor: pointer;

    transition: color 0.2s ease;
}

.avatar-edit:hover {
    color: rgba(0, 255, 120, 1);
}

.username-box {
    width: 100%;
    text-align: center;
}

.username-box h2 {
    margin-bottom: 14px;

    font-size: 1.2rem;
    letter-spacing: 1px;
}

.prompt {
    width: 100%;
    padding: 14px;

    border-radius: 14px;

    background: rgba(0, 0, 0, 0.5);
    border: 1px solid rgba(255, 255, 255, 0.06);

    color: rgba(245, 245, 245, 0.92);
    outline: none;

    transition: all 0.25s ease;
}

.prompt:focus {
    border: 1px solid rgba(0, 255, 120, 0.25);
    box-shadow: 0 0 20px rgba(0, 255, 120, 0.08);
}

.error {
    margin-top: 10px;
    color: rgba(255, 80, 80, 0.85);
    font-size: 0.9rem;
}

.success {
    margin-top: 10px;
    color: rgba(0, 255, 120, 0.7);
    font-size: 0.9rem;
}

.status-area {
    width: 100%;
    max-width: 520px;

    margin-top: 14px;

    display: flex;
    flex-direction: column;
    gap: 8px;

    text-align: center;
}

@media (max-width: 600px) {
    .profile-card {
        padding: 22px;
    }

    .avatar-big {
        width: 140px;
        height: 140px;
    }
}

.avatar-row {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
}

.avatar-actions {
    display: flex;
    flex-direction: column;
    gap: 12px;
}
</style>
