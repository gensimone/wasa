<script>
    import { setUserId, setName, setPhotoUrl, userState } from "@/state/user"
    export default {
        data() {
            return {
                name: null,
                error: null,
                loading: false
            };
        },
        methods: {
        async login() {
            this.error = null;
            this.loading = true;

            try {
                // Get the user id associated with username.
                let response = await this.$axios.post('/session', {
                    name: this.name
                });

                // Get the user informations (username and photo)
                // associated with the user id.
                let userId = response.data.userId;

                try {
                    response = await this.$axios.get(`/users/${userId}`, {
                        headers: { Authorization: userId }
                    });

                    let data = response.data;
                    setUserId(data.userId);
                    setName(data.name);
                    setPhotoUrl(data.photoUrl);

                    console.log(userState.photoUrl);

                    this.error = null;
                    this.$router.push("/home");
                } catch (e) {
                    this.error = e.response.data.error
                }
            } catch (e) {
                this.error = e.response.data.error
            }
            this.loading = false;
        }
    }
  };
</script>

<template>
  <h1 class="title"> WASAText </h1>
  <div class="login-container">
    <h2> Login </h2>
    <form @submit.prevent="login">
      <input
        v-model="name"
        type="text"
        placeholder="Username"
        required
      >
      <button type="submit" :disabled="loading">
        {{ loading ? "Logging in..." : "Login" }}
      </button>
    </form>
    <p v-if="error" class="error"> {{ error }} </p>
  </div>
</template>

<style scoped>
.title {
  text-align: center;
  font-size: 42px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #00ff41;
  text-shadow: 0 0 12px rgba(0, 255, 65, 0.8);
}

.login-container {
  max-width: 340px;
  padding: 30px;
  background: rgba(0, 0, 0, 0.85);
  border: 1px solid rgba(0, 255, 65, 0.5);
  box-shadow: 0 0 25px rgba(0, 255, 65, 0.15);
  color: #00ff41;
}

h2 {
  text-align: center;
  margin-bottom: 20px;
  color: #00ff41;
  font-weight: normal;
  letter-spacing: 2px;
  text-shadow: 0 0 8px rgba(0, 255, 65, 0.6);
}

input {
  width: 100%;
  padding: 12px;
  margin-top: 12px;
  background: rgba(0, 0, 0, 0.9);
  border: 1px solid rgba(0, 255, 65, 0.5);
  color: #00ff41;
  font-size: 14px;
  outline: none;
}

input:focus {
  border: 1px solid #00ff41;
  box-shadow: 0 0 10px rgba(0, 255, 65, 0.5);
}

input::placeholder {
  color: rgba(0, 255, 65, 0.5);
}

button {
  width: 100%;
  margin-top: 18px;
  padding: 12px;
  border: 1px solid rgba(0, 255, 65, 0.6);
  background: rgba(0, 255, 65, 0.12);
  color: #00ff41;
  font-size: 15px;
  letter-spacing: 2px;
  cursor: pointer;
  transition: 0.2s;
}

button:hover {
  background: rgba(0, 255, 65, 0.2);
  box-shadow: 0 0 12px rgba(0, 255, 65, 0.4);
}

button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.error {
  margin-top: 18px;
  padding: 10px;
  background: rgba(255, 0, 0, 0.1);
  border: 1px solid rgba(255, 0, 0, 0.4);
  color: #ff4d4d;
  font-size: 13px;
}
</style>
