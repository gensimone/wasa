<script>
import { clearUserState } from "@/state/user"
import { clearChatState } from "@/state/chat"
import Bottombar from "@/components/Shared/Bottombar.vue";
import Topbar from "@/components/Shared/Topbar.vue";

export default {
    components: {
        Bottombar,
        Topbar
    },
    data() {
        return {
            chats: [
            ],
        };
    },
    methods: {
        logout() {
            clearUserState()
            clearChatState()
            localStorage.clear()
            this.$router.replace("/")
        }
    }
};
</script>

<template>
    <div class="app">
        <Topbar :actions="[
            { icon: '/icons/settings.svg', onClick: () => $router.push('/settings') },
            { icon: '/icons/logout.svg', onClick: () => logout() }
        ]" />
        <main class="content">
            <div class="chat-list">
                <div class="chat-header">
                    <h2> Chats </h2>
                    <button class="icon-btn" @click="$router.push('/users')">
                        <img src="/icons/plus.svg" class="icon-img">
                    </button>
                </div>
                <div v-for="chat in chats" :key="chat.id" class="chat-item">
                    <div class="chat-photo-preview" />
                    <div class="info">
                        <div class="chat-name"> {{ chat.name }} </div>
                        <div class="chat-last-message"> {{ chat.lastMessage }} </div>
                    </div>
                </div>
            </div>
        </main>
        <Bottombar />
    </div>
</template>

<style scoped>
.chat-list {}

.chat-list {
    width: min(720px, 100%);
    border-radius: 22px;
    padding: 20px;

    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(20px);

    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.7);

    position: relative;
    overflow: hidden;

    animation: floatCard 5s ease-in-out infinite alternate;
}

@keyframes floatCard {
    0% {
        transform: translateY(0px);
    }

    100% {
        transform: translateY(10px);
    }
}

.chat-list::before {
    content: "";
    position: absolute;
    inset: -1px;
    border-radius: 22px;
    padding: 1px;

    background: linear-gradient(120deg,
            rgba(0, 255, 120, 0.14),
            rgba(255, 255, 255, 0.05),
            rgba(0, 255, 120, 0.14));

    -webkit-mask: linear-gradient(#000 0 0) content-box, linear-gradient(#000 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;

    opacity: 0.25;
    pointer-events: none;

    animation: borderGlow 4s ease-in-out infinite alternate;
}

@keyframes borderGlow {
    0% {
        opacity: 0.18;
    }

    100% {
        opacity: 0.35;
    }
}

.chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    margin-bottom: 18px;
    padding-bottom: 12px;

    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.chat-header h2 {
    margin: 0;
    font-size: 1.15rem;
    font-weight: 800;
    letter-spacing: 1px;
    color: rgba(245, 245, 245, 0.9);
}

.chat-item {
    display: flex;
    align-items: center;
    gap: 14px;

    padding: 14px 14px;
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

.chat-item::before {
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

.chat-item:hover::before {
    left: 70%;
}

.chat-item:hover {
    transform: translateY(-6px) scale(1.02);
    border: 1px solid rgba(0, 255, 120, 0.18);
    box-shadow: 0 18px 65px rgba(0, 0, 0, 0.75);
}

.chat-item:active {
    transform: scale(0.98);
}

/* =========================
   CHAT PHOTO
   ========================= */
.chat-photo-preview {
    width: 52px;
    height: 52px;
    border-radius: 16px;

    background: linear-gradient(145deg,
            rgba(255, 255, 255, 0.06),
            rgba(0, 255, 120, 0.12));

    box-shadow: 0 0 20px rgba(0, 0, 0, 0.55);
    position: relative;
    overflow: hidden;
}

/* subtle shimmer */
.chat-photo-preview::before {
    content: "";
    position: absolute;
    inset: 0;
    background: linear-gradient(120deg,
            transparent,
            rgba(255, 255, 255, 0.18),
            transparent);
    transform: translateX(-150%);
    animation: avatarSweep 3s ease-in-out infinite;
}

@keyframes avatarSweep {
    0% {
        transform: translateX(-150%);
    }

    50% {
        transform: translateX(150%);
    }

    100% {
        transform: translateX(150%);
    }
}

/* =========================
   CHAT INFO
   ========================= */
.info {
    display: flex;
    flex-direction: column;
    gap: 6px;
    min-width: 0;
}

.chat-name {
    font-size: 1.05rem;
    font-weight: 800;
    letter-spacing: 0.4px;

    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    color: rgba(245, 245, 245, 0.92);
}

.chat-last-message {
    font-size: 0.9rem;
    opacity: 0.65;

    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    color: rgba(200, 200, 200, 0.75);
}

/* =========================
   ENTRY ANIMATION
   ========================= */
@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(12px) scale(0.98);
    }

    to {
        opacity: 1;
        transform: translateY(0px) scale(1);
    }
}

/* =========================
   RESPONSIVE
   ========================= */
@media (max-width: 600px) {
    .topbar {
        padding: 0 14px;
    }

    .header-title {
        font-size: 1.2rem;
    }

    .chat-list {
        padding: 16px;
    }

    .chat-photo-preview {
        width: 46px;
        height: 46px;
    }
}
</style>
