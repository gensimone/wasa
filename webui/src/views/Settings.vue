<script>
import { setName, setPhotoUrl, userState } from "@/state/user"
import backIcon from "@/assets/icons/back.svg"
import saveIcon from "@/assets/icons/save.svg"
import LoggedAs from "@/components/LoggedAs.vue"

export default {
    components: {
        LoggedAs
    },
    data() {
        return {
            name: userState.name,
            avatarURL: userState.photoUrl,
            photo: null,

            photoStatusMsg: null,
            nameStatusMsg: null,
            errName: false,
            errPhoto: false,

            photoChanged: false,
            loading: false,

            backIcon,
            saveIcon,
        }
    },
    beforeUnmount() {
        if (this.avatarUrl) {
            URL.revokeObjectURL(this.avatarUrl)
        }
    },
    methods: {
        uploadPhoto(e) {
            const file = e.target.files[0]
            if (!file) return

            if (this.avatarURL) {
                URL.revokeObjectURL(this.avatarURL)
            }

            this.photo = file
            this.photoChanged = true
            this.avatarURL = URL.createObjectURL(file)
        },
        async apply() {
            this.loading = true
            this.nameStatusMsg = null
            this.photoStatusMsg = null

            const userId = userState.userId

            if (this.name !== userState.name) {
                try {
                    await this.$axios.put(`/users/${userId}/name`,
                        { name: this.name },
                        { headers: { Authorization: userId } }
                    )

                    setName(this.name);
                    this.errName = false
                    this.nameStatusMsg = "Success!"
                } catch (e) {
                    this.errName = true
                    this.nameStatusMsg = e.response.data.error
                }
            }

            if (this.photoChanged) {
                try {
                    const formData = new FormData()
                    formData.append("photo", this.photo)

                    const res = await this.$axios.put(
                        `/users/${userId}/photo`,
                        formData, {
                        headers: {
                            Authorization: userId,
                            "Content-Type": "multipart/form-data"
                        }
                    }
                    )
                    setPhotoUrl(res.data.photoUrl);

                    this.errPhoto = false
                    this.photoStatusMsg = "Success!"
                    this.photoChanged = false
                } catch (e) {
                    this.errPhoto = true
                    this.photoStatusMsg = e.response.data.error
                }
            }
            this.loading = false
        }
    }
}
</script>

<template>
    <div class="app">
        <header class="topbar">
            <div class="header-title"> Settings </div>
            <div class="actions">
                <button class="icon-btn" @click="$router.back()">
                    <img :src="backIcon" class="icon-img">
                </button>
                <button class="icon-btn" :disabled="loading" @click="apply">
                    <img :src="saveIcon" class="icon-img">
                </button>
            </div>
        </header>
        <div class="settings-page">
            <div class="section">
                <h2> Username </h2>
                <input v-model="name" class="prompt" type="text" placeholder=">_">
            </div>
            <p v-if="nameStatusMsg" :class="errName ? 'error' : 'success'"> {{ nameStatusMsg }} </p>
            <div class="section">
                <h2> Photo </h2>
                <div class="avatar-preview">
                    <img :src="avatarURL" />
                </div>
                <input type="file" accept="image/*" @change="uploadPhoto">
            </div>
            <p v-if="photoStatusMsg" :class="errPhoto ? 'error' : 'success'"> {{ photoStatusMsg }} </p>
        </div>
        <LoggedAs />
    </div>
</template>

<style scoped>
/* keep UI above background */
.topbar,
.settings-page {
    position: relative;
    z-index: 1;
}

/* =========================
   SETTINGS LAYOUT
   ========================= */

.settings-page {
    width: min(720px, 92%);
    margin: 40px auto;

    display: flex;
    flex-direction: column;
    gap: 26px;
}

/* =========================
   SECTION CARD
   ========================= */

.section {
    padding: 20px;

    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.06);

    border-radius: 18px;
    backdrop-filter: blur(18px);

    box-shadow: 0 20px 80px rgba(0, 0, 0, 0.75);

    transition: transform 0.25s ease, border 0.25s ease;
}

.section:hover {
    transform: translateY(-3px);
    border: 1px solid rgba(0, 255, 120, 0.15);
}

.section h2 {
    margin: 0 0 14px 0;
    font-size: 1.1rem;
    font-weight: 800;
    letter-spacing: 1px;
}

/* =========================
   INPUT (cyber minimal)
   ========================= */

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

/* =========================
   AVATAR
   ========================= */

.avatar-preview {
    width: 90px;
    height: 90px;
    border-radius: 18px;

    margin-bottom: 12px;
    overflow: hidden;

    border: 1px solid rgba(255, 255, 255, 0.08);
    background: rgba(255, 255, 255, 0.02);

    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.6);
    position: relative;
}

/* subtle scan shine */
.avatar-preview::before {
    content: "";
    position: absolute;
    inset: 0;

    background: linear-gradient(120deg,
            transparent,
            rgba(255, 255, 255, 0.12),
            transparent);

    transform: translateX(-140%);
    animation: avatarScan 3s ease-in-out infinite;
}

@keyframes avatarScan {
    0% {
        transform: translateX(-140%);
    }

    60% {
        transform: translateX(140%);
    }

    100% {
        transform: translateX(140%);
    }
}

.avatar-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

/* =========================
   FILE INPUT
   ========================= */

input[type="file"] {
    color: rgba(200, 200, 200, 0.8);
    font-size: 0.9rem;
}

/* =========================
   STATUS MESSAGES
   ========================= */

.error {
    color: rgba(255, 80, 80, 0.85);
    font-size: 0.9rem;
}

.success {
    color: rgba(0, 255, 120, 0.7);
    font-size: 0.9rem;
}

/* =========================
   RESPONSIVE
   ========================= */

@media (max-width: 600px) {
    .settings-page {
        margin: 20px auto;
    }

    .section {
        padding: 16px;
    }
}
</style>
