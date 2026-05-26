<script>
import MemberList from "@/components/Groups/MemberList.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";
import Topbar from "@/components/Shared/Topbar.vue";
import Poller from "@/services/poller";
import SettingsCard from "@/components/Settings/SettingsCard.vue";
import { usePhotoManager } from "@/composables/usePhotoManager";
import { useSettingsForm } from "@/composables/useSettingsForm";
import { defaultGroupPhotoUrl } from "@/assets/default";
import { handleError } from "@/utils/errors";
import { groupConversations } from "@/state/conversations";
import { getMemberIds, removeUser } from "@/services/groups";
import { getUserById } from "@/services/users";
import { defaultUserPhotoUrl } from "@/assets/default";
import { user } from "@/state/user";

export default {
  components: { Bottombar, Topbar, SettingsCard, MemberList },

  computed: {
    groupId() {
      return Number(this.$route.params.id);
    },

    founderId() {
      return groupConversations.value.get(this.groupId)?.founderId;
    },

    isFounder() {
      return this.founderId === user.userId;
    },

    name() {
      return groupConversations.value.get(this.groupId)?.name;
    },

    groupPhotoUrl() {
      const photo = groupConversations.value.get(this.groupId)?.photoUrl;
      console.log("photoURL:", photo);
      return photo;
    },
  },

  data() {
    const photo = usePhotoManager(this.groupPhotoUrl, defaultGroupPhotoUrl);

    const form = useSettingsForm(this.name, this.name);

    return {
      ...photo,
      ...form,

      poller: null,
      members: [],
    };
  },

  methods: {
    async removeUser(member) {
      try {
        await removeUser(this.groupId, member.userId);
        this.members = this.members.filter((m) => m.userId !== member.userId);
        this.$notifier.success(`User "${member.name}" removed`);
      } catch (e) {
        handleError(e);
      }
    },

    async startConversation(member) {
      this.$router.push({
        name: "conversation",
        params: { id: member.userId },
        query: { direct: true },
      });
    },

    updateGroup() {},
  },

  async mounted() {
    this.poller = new Poller(async () => {
      const memberIds = await getMemberIds(this.groupId);
      this.members = await Promise.all(
        memberIds.map(async (userId) => {
          const user = await getUserById(userId);
          if (!user.photoUrl) {
            user.photoUrl = defaultUserPhotoUrl;
          }

          return user;
        }),
      );
    }, 10000);

    this.poller.startPolling();
  },

  unmounted() {
    this.poller?.stopPolling();
  },
};
</script>

<template>
  <div class="app">
    <Topbar :actions="[{ icon: 'back', onClick: () => $router.back() }]" />
    <div class="content-column">
      <SettingsCard
        :enableEditing="isFounder"
        :photoUrl="photoUrl"
        :photoChanged="photoChanged"
        :text="text"
        title="Group name"
        submitButtonText="Update"
        :loading="loading"
        @uploadPhoto="uploadPhoto"
        @revertPhoto="revertPhoto"
        @deletePhoto="deletePhoto"
        @keyPress="setText"
        @submit="updateGroup"
      />
      <MemberList
        :members="members"
        @removeUser="removeUser"
        :founderId="founderId"
        @selectUser="startConversation"
      />
    </div>
    <Bottombar />
  </div>
</template>
