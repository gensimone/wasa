<script>
import { setChatName, setChatPhotoUrl, setChatUserId } from "@/state/chat"
import { user } from "@/state/user"
import LoggedAs from "@/components/LoggedAs.vue"
import backIcon from "@/assets/icons/back.svg"
export default {
    components: {
        LoggedAs
    },

    data() {
        return {
            backIcon,
            users: [],
            loading: false
        }
    },

    methods: {
        chatUser(user) {
            setChatName(user.name)
            setChatPhotoUrl(user.photoUrl)
            setChatUserId(user.userId)
            this.$router.push("/chat")
        },

        async fetchUser(userId) {
            const response = await this.$axios.get(`/users/${userId}`,
                { headers: { Authorization: user.userId } }
            )

            return response.data
        },

        async fetchUsers() {
            this.loading = true
            this.users = []

            try {
                const response = await this.$axios.get(`/users`, {
                    headers: { Authorization: user.userId }
                })

                const userPromises = response.data.users
                    .filter(userId => userId !== user.userId)
                    .map(userId =>
                        this.fetchUser(userId)
                    )

                const usersData = await Promise.all(userPromises)

                this.users = usersData
                    .filter(data => data.userId !== user.userId)
                    .map(data => ({
                        userId: data.userId,
                        name: data.name,
                        photoUrl: `${__API_URL__}${data.photoUrl}`
                    }))
            } catch (e) {
                this.message = e?.response?.data?.error || "Unexpected error"
                this.users = []
            }

            this.loading = false
        }
    },

    mounted() {
        this.fetchUsers()
    }
}
</script>

<template>
    <div class="app">
        <header class="topbar">
            <div class="header-title"> WASAText </div>
            <div class="actions">
                <button class="icon-btn" @click="$router.back()">
                    <img :src="backIcon" class="icon-img">
                </button>
            </div>
        </header>
        <div class="content">
            <div class="user-list">
                <div class="user-list-header">
                    <h2>Users</h2>
                </div>
                <div v-for="user in users" :key="user.userId" class="user-item" @click="chatUser(user)">
                    <!-- PHOTO -->
                    <div class="user-photo-wrapper">
                        <img :src="user.photoUrl" class="user-photo" alt="user photo" />
                    </div>
                    <!-- INFO -->
                    <div class="user-info">
                        <div class="user-name">
                            {{ user.name }}
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <LoggedAs />
    </div>
</template>

<style scoped>
.user-list {
    position: relative;
    z-index: 1;
    width: min(720px, 100%);
    border-radius: 22px;
    padding: 20px;

    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(20px);

    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.7);

    overflow: hidden;

    animation: floatCard 5s ease-in-out infinite alternate;
}

/* border glow */
.user-list::before {
    content: "";
    position: absolute;
    inset: -1px;
    border-radius: 22px;

    background: linear-gradient(120deg,
            rgba(0, 255, 120, 0.14),
            rgba(255, 255, 255, 0.05),
            rgba(0, 255, 120, 0.14));

    -webkit-mask: linear-gradient(#000 0 0) content-box,
        linear-gradient(#000 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;

    opacity: 0.25;
    pointer-events: none;
}

/* =========================
   HEADER
   ========================= */
.user-list-header {
    margin-bottom: 18px;
    padding-bottom: 12px;

    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.user-list-header h2 {
    margin: 0;
    font-size: 1.15rem;
    font-weight: 800;
    letter-spacing: 1px;
    color: rgba(245, 245, 245, 0.9);
}

/* =========================
   USER ITEM
   ========================= */
.user-item {
    display: flex;
    align-items: center;
    gap: 14px;

    padding: 14px;
    margin-bottom: 10px;
    border-radius: 18px;

    background: rgba(255, 255, 255, 0.02);
    border: 1px solid rgba(255, 255, 255, 0.06);

    cursor: pointer;
    position: relative;
    overflow: hidden;

    transition: transform 0.25s ease, border 0.25s ease, box-shadow 0.25s ease;

    animation: fadeInUp 0.35s ease both;
}

.user-item:hover {
    transform: translateY(-6px) scale(1.02);
    border: 1px solid rgba(0, 255, 120, 0.18);
    box-shadow: 0 18px 65px rgba(0, 0, 0, 0.75);
}

/* shine effect */
.user-item::before {
    content: "";
    position: absolute;
    top: -40%;
    left: -90%;
    width: 130%;
    height: 200%;

    background: linear-gradient(90deg,
            transparent,
            rgba(0, 255, 120, 0.08),
            rgba(255, 255, 255, 0.04),
            transparent);

    transform: rotate(18deg);
    transition: left 0.65s ease;
}

.user-item:hover::before {
    left: 70%;
}

/* =========================
   USER PHOTO
   ========================= */
.user-photo-wrapper {
    width: 52px;
    height: 52px;
    border-radius: 16px;
    overflow: hidden;

    background: linear-gradient(145deg,
            rgba(255, 255, 255, 0.06),
            rgba(0, 255, 120, 0.12));
}

.user-photo {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

/* =========================
   USER INFO
   ========================= */
.user-info {
    display: flex;
    flex-direction: column;
    gap: 6px;
    min-width: 0;
}

.user-name {
    font-size: 1.05rem;
    font-weight: 800;
    color: rgba(245, 245, 245, 0.92);

    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.user-id {
    font-size: 0.9rem;
    opacity: 0.65;
    color: rgba(200, 200, 200, 0.75);
}

/* =========================
   ANIMATIONS
   ========================= */
@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(12px) scale(0.98);
    }

    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

@keyframes floatCard {
    from {
        transform: translateY(0);
    }

    to {
        transform: translateY(10px);
    }
}
</style>
