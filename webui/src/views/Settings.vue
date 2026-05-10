<script>
    import Footer from "@/components/Footer.vue";
    export default {
        data() {
            return {
                username: localStorage.getItem("userName"),
                imageUrl: null,
                imageFile: null,
                error: null
            };
        },
        methods: {
            onFileChange(e) {
                const file = e.target.files[0];
                if (!file) return;
                this.imageFile = file;
                this.imageUrl = URL.createObjectURL(file);
            },

            apply() {
            },
        },
        components: {
            Footer
        }
    };
</script>

<template>
    <div class="app">
        <header class="topbar">
            <div class="header-title"> Settings </div>
            <div class="actions">
                <button @click="$router.back()"> Back </button>
            </div>
        </header>

            <div class="settings-page">
                <h2> Settings </h2>
                <div class="section">
                    <label> Photo </label>
                    <div class="avatar-preview">
                        <img v-if="imageUrl" :src="imageUrl" />
                        <div v-else class="placeholder">No Image</div>
                    </div>
                    <input type="file" @change="onFileChange" accept="image/*" />
                </div>
                <div class="section">
                    <label> Username </label>
                    <input
                            v-model="username"
                            type="text"
                            placeholder="Username"
                            />
                </div>
                <button @click="apply"> Save </button>
            </div>
            <Footer />
    </div>
</template>

<style>
.settings-page {
  min-height: 100vh;
  color: #00ff41;
  font-family: "Courier New", monospace;
  padding: 40px;
}

.content {
    display: flex;
    justify-content: center;
    padding-top: 40px;
}

h2 {
  margin-bottom: 30px;
  text-shadow: 0 0 10px rgba(0, 255, 65, 0.5);
}

.section {
  margin-bottom: 25px;
}

label {
  display: block;
  margin-bottom: 8px;
  opacity: 0.8;
}

input[type="text"] {
  width: 300px;
  padding: 10px;
  background: rgba(0,0,0,0.8);
  border: 1px solid rgba(0, 255, 65, 0.4);
  color: #00ff41;
  outline: none;
}

input[type="file"] {
  margin-top: 10px;
  color: #00ff41;
}

.avatar-preview {
  width: 120px;
  height: 120px;
  margin-bottom: 10px;

  border: 1px solid rgba(0, 255, 65, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;

  overflow: hidden;
}

.avatar-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.placeholder {
  font-size: 12px;
  opacity: 0.5;
}

button {
  padding: 10px 18px;
  background: transparent;
  border: 1px solid #00ff41;
  color: #00ff41;
  cursor: pointer;
}

button:hover {
  background: rgba(0, 255, 65, 0.1);
}
</style>
