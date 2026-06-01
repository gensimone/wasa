<script>
import { user, updateUserState } from "@/state/user";
import { defaultUserPhotoUrl } from "@/assets/default";
import { setMyUserName, setMyPhoto, deleteMyPhoto } from "@/services/users";
import { handleError } from "@/utils/errors";
import { usePhotoManager } from "@/composables/usePhotoManager";
import { useSettingsForm } from "@/composables/useSettingsForm";
import { getIcon } from "@/state/theme";
import SettingsCard from "@/components/Settings/SettingsCard.vue";
import Topbar from "@/components/Shared/Topbar.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";

export default {
  components: { SettingsCard, Topbar, Bottombar },

  data() {
    const photo = usePhotoManager(user.photoUrl, defaultUserPhotoUrl);

    const form = useSettingsForm(user.name, user.name);

    return {
      user,

      ...photo,
      ...form,
    };
  },

  watch: {
    "user.name"(v) {
      this.text = v;
      this.placeholder = v;
    },

    "user.photoUrl"(v) {
      this.photoUrl = v;
      this.photo = null;
      this.photoChanged = false;
    },
  },

  beforeUnmount() {
    this.revoke();
  },

  methods: {
    getIcon,

    async updateProfile() {
      try {
        await this.submit(async (name) => {
          let changed = false;

          if (name !== user.name) {
            const updatedName = await setMyUserName(name);

            updateUserState({ name: updatedName });
            changed = true;
          }

          if (this.photo) {
            const url = await setMyPhoto(this.photo);
            updateUserState({ photoUrl: url });
            changed = true;
          } else if (this.photoChanged) {
            await deleteMyPhoto();
            updateUserState({ photoUrl: defaultUserPhotoUrl });
            changed = true;
          }

          if (!changed) throw new Error("NO_CHANGE");
        });

        this.$notifier.success("Profile updated successfully");
      } catch (e) {
        if (e.message === "EMPTY_NAME")
          this.$notifier.error("Invalid user name");
        else if (e.message === "NO_CHANGE")
          this.$notifier.warning("Nothing to do..");
        else {
          handleError(e);
        }
      }
    },
  },
};
</script>

<template>
  <div class="app">
    <Topbar :actions="[{ icon: 'back', onClick: () => $router.back() }]" />
    <div class="content-center">
      <div class="user-settings">
        <SettingsCard
          :photoUrl="photoUrl"
          :enableEditing="true"
          :photoChanged="photoChanged"
          :text="text"
          title="Username"
          submitButtonText="Update"
          :loading="loading"
          @uploadPhoto="uploadPhoto"
          @revertPhoto="revertPhoto"
          @deletePhoto="deletePhoto"
          @keyPress="setText"
          @submit="updateProfile"
        />
      </div>
    </div>
    <Bottombar />
  </div>
</template>

<style scoped>
.content-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.user-settings {
  width: min(600px, 50%);
}
</style>
