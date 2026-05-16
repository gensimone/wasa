<script>
import { clearUserState } from "@/state/user"
import logoutIcon from "@/assets/icons/logout.svg";
import settingsIcon from "@/assets/icons/settings.svg";
import LoggedAs from "@/components/LoggedAs.vue";

export default {
    components: {
        LoggedAs
    },
    data() {
        return {
            chats: [
                { id: 1, name: "Neo", lastMessage: "The Matrix has you..." },
                { id: 2, name: "Morpheus", lastMessage: "Follow the white rabbit." },
                { id: 4, name: "Trinity", lastMessage: "Are you awake?" },
                { id: 5, name: "Trinity", lastMessage: "Are you awake?" },
                { id: 6, name: "Trinity", lastMessage: "Are you awake?" },
            ],
            groups: [
                { id: 1, name: "Zion", lastMessage: "People of zion..." },
                { id: 2, name: "Spleeping People", lastMessage: "Wake up now!" },
            ],
            logoutIcon,
            settingsIcon
        };
    },
    methods: {
        logout() {
            clearUserState();
            this.$router.replace("/");
        }
    }
};
</script>

<template>
    <div class="app">
        <header class="topbar">
            <div class="header-title"> WASAText </div>
            <div class="actions">
                <button class="icon-btn" @click="$router.push('/settings')">
                    <img :src="settingsIcon" class="icon-img">
                </button>
                <button class="icon-btn" @click="logout">
                    <img :src="logoutIcon" class="icon-img">
                </button>
            </div>
        </header>
        <main class="content">
            <div class="chat-list">
                <div class="chat-header">
                    <h2> Chats </h2>
                    <button class="matrix-button" @click="$router.push('/users')">
                        <svg viewBox="0 0 24 24" class="plus-icon">
                            <path d="M12 5v14M5 12h14" />
                        </svg>
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
        <main class="content">
            <div class="chat-list">
                <div class="chat-header">
                    <h2> Groups </h2>
                    <button class="matrix-button" @click="$router.push('/create/group')">
                        <svg viewBox="0 0 24 24" class="plus-icon">
                            <path d="M12 5v14M5 12h14" />
                        </svg>
                    </button>
                </div>
                <div v-for="group in groups" :key="group.id" class="chat-item">
                    <div class="chat-photo-preview" />
                    <div class="info">
                        <div class="chat-name"> {{ group.name }} </div>
                        <div class="chat-last-message"> {{ group.lastMessage }} </div>
                    </div>
                </div>
            </div>
        </main>
        <LoggedAs />
    </div>
</template>

<style>
/* keep everything above background */
.topbar,
.content,
.chat-list {
    position: relative;
    z-index: 1;
}

/* =========================
   CONTENT SECTION
   ========================= */
.content {
    padding: 18px;
    display: flex;
    justify-content: center;
}

/* =========================
   CHAT LIST CARD
   ========================= */
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

/* subtle green border shimmer */
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

/* =========================
   CHAT HEADER
   ========================= */
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

/* =========================
   MATRIX BUTTON
   ========================= */
.matrix-button {
    width: 48px;
    height: 48px;
    border-radius: 14px;

    border: 1px solid rgba(0, 255, 120, 0.18);
    background: rgba(255, 255, 255, 0.02);

    display: flex;
    justify-content: center;
    align-items: center;

    cursor: pointer;
    position: relative;
    overflow: hidden;

    transition: transform 0.25s ease, box-shadow 0.25s ease, border 0.25s ease;
}

/* matrix scan effect */
.matrix-button::before {
    content: "";
    position: absolute;
    inset: 0;

    background: repeating-linear-gradient(to bottom,
            rgba(0, 255, 120, 0.18),
            rgba(0, 255, 120, 0.18) 2px,
            transparent 2px,
            transparent 8px);

    opacity: 0;
    transform: translateY(-120%);
}

.matrix-button:hover::before {
    opacity: 0.55;
    animation: matrixScan 0.75s linear infinite;
}

@keyframes matrixScan {
    from {
        transform: translateY(-120%);
    }

    to {
        transform: translateY(120%);
    }
}

.matrix-button:hover {
    transform: rotate(-2deg) scale(1.08);
    border: 1px solid rgba(0, 255, 120, 0.3);
    box-shadow: 0 0 35px rgba(0, 255, 120, 0.12);
}

.matrix-button:active {
    transform: scale(0.95);
}

.plus-icon {
    width: 22px;
    height: 22px;
    stroke: rgba(0, 255, 120, 0.9);
    stroke-width: 2.6;
    fill: none;
    filter: drop-shadow(0 0 10px rgba(0, 255, 120, 0.2));
}

/* =========================
   CHAT ITEM
   ========================= */
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

/* hover streak */
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
