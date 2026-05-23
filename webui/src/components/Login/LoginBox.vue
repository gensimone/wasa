<script>
export default {
    props: {
        name: String,
        loading: Boolean
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
                <button class="submit-button" type="submit" :disabled="loading">
                    {{ loading ? "Logging in..." : "Login" }}
                </button>
            </form>
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

@media (max-width: 500px) {
    .login-container {
        margin: 40px auto;
        padding: 22px;
    }
}
</style>
