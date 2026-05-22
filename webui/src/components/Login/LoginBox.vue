<script>
export default {
    props: {
        name: String,
        loading: Boolean,
        error: String
    },
    emits: ["update:name", "submit"]
}
</script>

<template>
    <div class="login-page">
        <h3 class="subtitle">
            Web Application And Software Architecture <br />
            Project
        </h3>
        <div class="login-container">
            <h2>Login</h2>
            <form @submit.prevent="$emit('submit')">
                <input class="prompt" :value="name" @input="$emit('update:name', $event.target.value)" type="text"
                    placeholder="Username" required />
                <button type="submit" :disabled="loading">
                    {{ loading ? "Logging in..." : "Login" }}
                </button>
            </form>
            <p v-if="error" class="error">
                {{ error }}
            </p>
        </div>
    </div>
</template>

<style scoped>
.login-page {
    min-height: calc(100vh - 70px);
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

.login-container {
    position: relative;
    z-index: 1;

    justify-content: center;
    width: min(420px, 90%);
    margin: 60px auto;
    padding: 28px;

    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(22px);

    border-radius: 20px;
    border: 1px solid rgba(255, 255, 255, 0.08);

    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.75);
}

.login-container h2 {
    margin: 0 0 18px 0;
    font-size: 1.4rem;
    font-weight: 800;
    letter-spacing: 1px;
    color: rgba(245, 245, 245, 0.9);
}

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

.error {
    margin-top: 12px;
    font-size: 0.9rem;
    color: rgba(255, 80, 80, 0.85);
    text-shadow: 0 0 10px rgba(255, 0, 0, 0.08);
}

@media (max-width: 500px) {
    .login-container {
        margin: 40px auto;
        padding: 22px;
    }
}
</style>
