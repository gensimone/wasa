<script>
export default {
  props: {
    name: String,
    loading: Boolean,
  },
  emits: ["update:name", "submit"],
};
</script>

<template>
  <div class="login-box">
    <h3 class="login-box-title">
      Web Application And Software Architecture <br />
      Project
    </h3>

    <div class="login-box-card">
      <h2>Login</h2>

      <form @submit.prevent="$emit('submit')">
        <input
          name="login"
          class="input-bar"
          :value="name"
          @input="$emit('update:name', $event.target.value)"
          type="text"
          placeholder="Username"
          required
        />

        <button class="submit-button" type="submit" :disabled="loading">
          {{ loading ? "Logging in..." : "Login" }}
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-box {
  min-height: calc(100vh - 70px);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  padding: 20px;
}

.login-box-title {
  font-size: 1.5rem;
  color: var(--text);
  text-transform: uppercase;
}

.login-box-card {
  width: min(420px, 90%);
  margin: 60px auto;
  padding: 28px;

  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 20px;

  box-shadow: 0 25px 90px var(--shadow);
}

.login-box-card h2 {
  margin: 0 0 18px 0;
  font-size: 1.4rem;
  font-weight: 800;
  letter-spacing: 1px;
  color: var(--text);
}
</style>
