<script>
    import logoutIcon from "@/assets/icons/logout.svg";
    import settingsIcon from "@/assets/icons/settings.svg";
    import Footer from "@/components/Footer.vue";
    export default {
        data() {
            return {
                chats: [
                    { id: 1, name: "Neo", lastMessage: "The Matrix has you..." },
                    { id: 2, name: "Morpheus", lastMessage: "Follow the white rabbit." },
                    { id: 3, name: "Trinity", lastMessage: "Are you awake?" },
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
                localStorage.clear();
                this.$router.replace("/");
            }
        },
        components: {
            Footer
        }
    };
</script>

<template>
    <div class="app">
        <header class="topbar">
            <div class="header-title"> WASAText </div>
            <div class="actions">
                <button class="icon-btn" @click="logout">
                    <img :src="logoutIcon" class="icon-img" />
                </button>
                <button class="icon-btn" @click="$router.push('/settings')">
                    <img :src="settingsIcon" class="icon-img" />
                </button>
            </div>
        </header>
        <main class="content">
            <div class="chat-list">
                <div class="chat-header">
                    <h2> Chats </h2>
                    <button class="matrix-button" @click="$router.push('/newChat')">
                        <svg viewBox="0 0 24 24" class="plus-icon">
                            <path d="M12 5v14M5 12h14" />
                        </svg>
                    </button>
                </div>
                <div
                        v-for="chat in chats"
                        :key="chat.id"
                        class="chat-item"
                        >
                        <div class="chat-photo-preview"></div>
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
                    <button class="matrix-button" @click="$router.push('/newGroup')">
                        <svg viewBox="0 0 24 24" class="plus-icon">
                            <path d="M12 5v14M5 12h14" />
                        </svg>
                    </button>
                </div>
                <div
                        v-for="group in groups"
                        :key="group.id"
                        class="chat-item"
                        >
                        <div class="chat-photo-preview"></div>
                        <div class="info">
                            <div class="chat-name"> {{ group.name }} </div>
                            <div class="chat-last-message"> {{ group.lastMessage }} </div>
                        </div>
                </div>
            </div>
        </main>
        <Footer />
    </div>
</template>

<style>
.actions {
  display: flex;
  flex-direction: row;
  gap: 10px;
  align-items: center;
}

.content {
    display: flex;
    justify-content: center;
    padding-top: 40px;
}

.icon-img {
    width: 36px;
    height: 36px;
    filter: drop-shadow(0 0 6px rgba(0, 255, 65, 0.8));
    opacity: 0.9;
}

.icon-btn {
    left: 0%;
    width: 10px;
    height: 36px;
    display: flex;
    justify-content: center;
    align-items: center;
    background: rgba(0, 0, 0, 0.6);
    border: 0px;
    border-radius: 50%;
    cursor: pointer;
}
</style>
