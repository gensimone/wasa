<script>
  import { setName, setPhotoUrl, userState } from "@/state/user"
  import backIcon from "@/assets/icons/back.svg"
  import saveIcon from "@/assets/icons/save.svg"
  import LoggedAs from "@/components/LoggedAs.vue"

  export default {
      components: {
          LoggedAs
      },
      data() {
          return {
              name: userState.name,
              avatarURL: userState.photoUrl,
              photo: null,

              photoStatusMsg: null,
              nameStatusMsg: null,
              errName: false,
              errPhoto: false,

              photoChanged: false,
              loading: false,

              backIcon,
              saveIcon,
          }
      },
      beforeUnmount() {
          if (this.avatarUrl) {
              URL.revokeObjectURL(this.avatarUrl)
          }
      },
      methods: {
          uploadPhoto(e) {
              const file = e.target.files[0]
              if (!file) return

              if (this.avatarURL) {
                  URL.revokeObjectURL(this.avatarURL)
              }

              this.photo = file
              this.photoChanged = true
              this.avatarURL = URL.createObjectURL(file)
          },
          async apply() {
              this.loading = true
              this.nameStatusMsg = null
              this.photoStatusMsg = null

              const userId = userState.userId

              if (this.name !== userState.name ) {
                  try {
                      await this.$axios.put(`/users/${userId}/name`,
                          { name: this.name },
                          { headers: { Authorization: userId } }
                      )

                      setName(this.name);
                      this.errName = false
                      this.nameStatusMsg = "Success!"
                  } catch (e) {
                      this.errName = true
                      this.nameStatusMsg = e.response.data.error
                  }
              }

              if (this.photoChanged) {
                  try {
                      const formData = new FormData()
                      formData.append("photo", this.photo)

                      const res = await this.$axios.put(
                          `/users/${userId}/photo`,
                          formData, {
                              headers: {
                                  Authorization: userId,
                                  "Content-Type": "multipart/form-data"
                              }
                          }
                      )
                      setPhotoUrl(res.data.photoUrl);

                      this.errPhoto = false
                      this.photoStatusMsg = "Success!"
                      this.photoChanged = false
                  } catch (e) {
                      this.errPhoto = true
                      this.photoStatusMsg = e.response.data.error
                  }
              }
              this.loading = false
          }
      }
  }
</script>

<template>
  <div class="app">

    <!-- Header -->
    <header class="topbar">
      <div class="header-title"> Settings </div>
      <div class="actions">
        <button class="icon-btn" @click="$router.back()">
          <img :src="backIcon" class="icon-img">
        </button>
        <button class="icon-btn" :disabled="loading" @click="apply">
          <img :src="saveIcon" class="icon-img">
        </button>
      </div>
    </header>

    <div class="settings-page">

      <!-- Username section -->
      <div class="section">
        <h2> Username </h2>
        <input
          v-model="name"
          class="prompt"
          type="text"
          placeholder=">_"
        >
      </div>
      <p v-if="nameStatusMsg" :class="errName ? 'error' : 'success'"> {{ nameStatusMsg }} </p>

      <!-- Photo section -->
      <div class="section">
        <h2> Photo </h2>
        <div class="avatar-preview">
          <img :src="avatarURL" />
        </div>
        <input type="file" accept="image/*" @change="uploadPhoto">
      </div>
      <p v-if="photoStatusMsg" :class="errPhoto ? 'error' : 'success'"> {{ photoStatusMsg }} </p>

    <!-- Footer -->
    </div>
    <LoggedAs />

  </div>
</template>

<style>
.settings-page {
  min-height: 100vh;
  color: #00ff41;
  padding: 40px;
}

.content {
    display: flex;
    justify-content: center;
    padding-top: 40px;
}

.section {
  margin-bottom: 25px;
}

.prompt {
  width: 100%;
  padding: 12px;
  margin-top: 12px;
  background: rgba(0, 0, 0, 0.9);
  border: 1px solid rgba(0, 255, 65, 0.5);
  color: #00ff41;
  font-size: 14px;
  outline: none;
}

.prompt:focus {
  border: 1px solid #00ff41;
  box-shadow: 0 0 10px rgba(0, 255, 65, 0.5);
}

.prompt::placeholder {
  color: rgba(0, 255, 65, 0.5);
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
  font-size: 0px;
}

.success {
  margin-top: 18px;
  padding: 10px;
  background: rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(0, 255, 65, 0.4);
  color: #00ff41;
  font-size: 13px;
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
