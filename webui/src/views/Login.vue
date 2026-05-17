<script>
import { setUserId, setName, setPhotoUrl } from "@/state/user"
import githubIcon from "@/assets/icons/github.svg"
export default {
    data() {
        return {
            name: null,
            error: null,
            loading: false,
            githubIcon
        }
    },
    methods: {
        async login() {
            this.error = null
            this.loading = true

            try {
                // Get the user id associated with username.
                let response = await this.$axios.post('/session', {
                    name: this.name
                })

                // Get the user informations (username and photo)
                // associated with the user id.
                const userId = response.data.userId

                try {
                    response = await this.$axios.get(`/users/${userId}`, {
                        headers: { Authorization: userId }
                    })

                    const data = response.data
                    setUserId(data.userId)
                    setName(data.name)
                    setPhotoUrl(data.photoUrl)

                    this.error = null
                    this.$router.push("/home")
                } catch (e) {
                    this.error = e?.response?.data?.error || "Unexpected error"
                }
            } catch (e) {
                this.error = e?.response?.data?.error || "Unexpected error"
            }
            this.loading = false
        }
    }
};
</script>

<template>
    <div class="app">
        <header class="topbar">
            <div class="header-title"> WASAText </div>
            <div class="actions">
                <a class="icon-btn" href="https://github.com/gensimone" target="_blank">
                    <img :src="githubIcon" alt="logo" class="icon-img">
                </a>
            </div>
        </header>
        <div class="login-page">
            <h3 class="subtitle">
                Web Application And Software Architecture <br>
                Project
            </h3>
            <div class="login-container">
                <h2> Login </h2>
                <form @submit.prevent="login">
                    <input v-model="name" type="text" placeholder="Username" required>
                    <button type="submit" :disabled="loading">
                        {{ loading ? "Logging in..." : "Login" }}
                    </button>
                </form>
                <p v-if="error" class="error"> {{ error }} </p>
            </div>
        </div>
        <footer class="footer">
            Made by Simone Gentili
        </footer>
    </div>
</template>

<style scoped>
.login-page {
    min-height: calc(100vh - 70px);
    /* accounts for topbar */
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;

    text-align: center;
    padding: 20px;
}

.subtitle {
    margin-bottom: 18px;

    font-size: 0.95rem;
    letter-spacing: 2px;
    text-transform: uppercase;

    color: rgba(200, 200, 200, 0.65);

    text-shadow: 0 0 10px rgba(0, 255, 120, 0.08);
}

/* keep UI above background */
.login-container {
    position: relative;
    z-index: 1;
}

.login-container {
    justify-content: center;
    width: min(420px, 90%);
    margin: 60px auto;
    padding: 28px;

    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(22px);

    border-radius: 20px;
    border: 1px solid rgba(255, 255, 255, 0.08);

    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.75);

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

/* =========================
   LOGIN TITLE
   ========================= */

.login-container h2 {
    margin: 0 0 18px 0;
    font-size: 1.4rem;
    font-weight: 800;
    letter-spacing: 1px;

    color: rgba(245, 245, 245, 0.9);
}

/* =========================
   INPUT
   ========================= */

input {
    width: 100%;
    padding: 14px 14px;
    margin-bottom: 14px;

    border-radius: 14px;

    background: rgba(0, 0, 0, 0.45);
    border: 1px solid rgba(255, 255, 255, 0.06);

    color: rgba(245, 245, 245, 0.92);
    outline: none;

    transition: all 0.25s ease;
}

/* subtle green focus glow */
input:focus {
    border: 1px solid rgba(0, 255, 120, 0.25);
    box-shadow: 0 0 20px rgba(0, 255, 120, 0.08);
}

/* =========================
   BUTTON
   ========================= */

button {
    width: 100%;
    padding: 14px;

    border-radius: 14px;
    border: 1px solid rgba(0, 255, 120, 0.18);

    background: rgba(255, 255, 255, 0.02);
    color: rgba(245, 245, 245, 0.92);

    font-weight: 700;
    letter-spacing: 1px;

    cursor: pointer;
    position: relative;
    overflow: hidden;

    transition: all 0.25s ease;
}

/* matrix scan effect */
button::before {
    content: "";
    position: absolute;
    inset: 0;

    background: linear-gradient(120deg,
            transparent,
            rgba(0, 255, 120, 0.12),
            transparent);

    transform: translateX(-140%);
    transition: transform 0.7s ease;
}

button:hover::before {
    transform: translateX(140%);
}

button:hover {
    transform: translateY(-2px);
    border: 1px solid rgba(0, 255, 120, 0.3);
    box-shadow: 0 0 25px rgba(0, 255, 120, 0.1);
}

button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
}

/* =========================
   ERROR TEXT
   ========================= */

.error {
    margin-top: 12px;
    font-size: 0.9rem;
    color: rgba(255, 80, 80, 0.85);
    text-shadow: 0 0 10px rgba(255, 0, 0, 0.08);
}

/* =========================
   RESPONSIVE
   ========================= */

@media (max-width: 500px) {
    .login-container {
        margin: 40px auto;
        padding: 22px;
    }

    .title {
        font-size: 2rem;
    }
}
</style>
