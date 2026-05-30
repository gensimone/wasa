<script>
import { createGroup } from "@/services/groups";
import { handleError } from "@/utils/errors";
import { usePhotoManager } from "@/composables/usePhotoManager";
import { useSettingsForm } from "@/composables/useSettingsForm";
import { defaultGroupPhotoUrl } from "@/assets/default";
import { groups } from "@/state/conversations";
import { user } from "@/state/user";
import SettingsCard from "@/components/Settings/SettingsCard.vue";
import Topbar from "@/components/Shared/Topbar.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";

export default {
  components: { SettingsCard, Topbar, Bottombar },

  data() {
    const photo = usePhotoManager(defaultGroupPhotoUrl, defaultGroupPhotoUrl);

    const form = useSettingsForm("...", "...");

    return {
      ...photo,
      ...form,
    };
  },

  methods: {
    async createGroup() {
      try {
        await this.submit(async (name) => {
          const isNotValid = Array.from(groups.value.values()).some((g) => {
            return g.founderId === user.userId && g.name === name;
          });

          if (isNotValid) {
            this.$notifier.error("Invalid group name.");
            return;
          }

          const group = await createGroup(name, this.photo);

          groups.value.set(group.conversationId, group);

          this.$router.push({
            name: "conversation",
            params: { id: group.conversationId },
            query: { direct: false },
          });
        });
      } catch (e) {
        if (e.message === "EMPTY_NAME") {
          this.$notifier.error("Invalid group name");
        } else {
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
      <SettingsCard
        :photoUrl="photoUrl"
        :enableEditing="true"
        :photoChanged="photoChanged"
        :text="text"
        title="Group name"
        submitButtonText="Create group"
        :loading="loading"
        @uploadPhoto="uploadPhoto"
        @revertPhoto="revertPhoto"
        @deletePhoto="deletePhoto"
        @keyPress="setText"
        @submit="createGroup"
      />
    </div>
    <Bottombar />
  </div>
</template>
