<script>
import MemberList from "@/components/Groups/MemberList.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";
import Topbar from "@/components/Shared/Topbar.vue";
import Poller from "@/services/poller";
import SettingsCard from "@/components/Settings/SettingsCard.vue";
import UsersList from "@/components/Users/UsersList.vue";
import { usePhotoManager } from "@/composables/usePhotoManager";
import { useSettingsForm } from "@/composables/useSettingsForm";
import { defaultGroupPhotoUrl } from "@/assets/default";
import { handleError } from "@/utils/errors";
import { groupConversations } from "@/state/conversations";
import { getMemberIds, removeUser } from "@/services/groups";
import { getUserById } from "@/services/users";
import { defaultUserPhotoUrl } from "@/assets/default";
import { user } from "@/state/user";
import { getIcon } from "@/state/theme";
import { addToGroup } from "@/services/groups";

export default {
  components: { Bottombar, Topbar, SettingsCard, MemberList, UsersList },

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

    groupName() {
      return groupConversations.value.get(this.groupId)?.name;
    },

    groupPhotoUrl() {
      return (
        groupConversations.value.get(this.groupId)?.photoUrl ||
        defaultGroupPhotoUrl
      );
    },
  },

  watch: {
    groupPhotoUrl: {
      immediate: true,
      handler(newUrl) {
        this.photoUrl = newUrl;
      },
    },

    groupName: {
      immediate: true,
      handler(newName) {
        this.text = newName || "";
      },
    },
  },

  data() {
    const photo = usePhotoManager(defaultGroupPhotoUrl, defaultGroupPhotoUrl);

    const form = useSettingsForm();

    return {
      ...photo,
      ...form,

      poller: null,
      members: [],
      inAddToGroup: false,
    };
  },

  methods: {
    getIcon,

    async addMemberToGroup(user) {
      try {
        const newMember = await addToGroup(this.groupId, user.userId);
        this.members.push(newMember);
        this.$notifier.success(`User ${newMember.name} added`);
      } catch (e) {
        handleError(e);
      }
    },

    addUsersGroup(users) {
      Promise.all(users.map(this.addMemberToGroup));
      this.inAddToGroup = false;
    },

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
      <div class="items-list">
        <div v-if="inAddToGroup">
          <button
            class="submit-button"
            style="margin-bottom: 10px"
            @click="inAddToGroup = false"
          >
            <img class="icon-img" :src="getIcon('back')" />
            Members
          </button>
          <UsersList
            @select="addUsersGroup"
            :excludeUsers="members"
            :canSelectMultiple="true"
          />
        </div>
        <div v-else>
          <MemberList
            :members="members"
            @goInAddToGroup="inAddToGroup = true"
            @removeUser="removeUser"
            :founderId="founderId"
            @selectUser="startConversation"
          />
        </div>
      </div>
    </div>
    <Bottombar />
  </div>
</template>
<style scoped>
.content-column {
  display: flex;
  gap: 20px;
  align-items: flex-start;

  max-width: 1200px; /* controlla la larghezza del layout */
  margin: 100px auto; /* centra orizzontalmente */
  padding: 20px;
}

@media (max-width: 900px) {
  .content-column {
    flex-direction: column;
  }
}
.items-list {
  width: min(720px, 100%);
  padding: 20px;
  border-radius: 22px;

  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: 0 25px 90px var(--shadow);

  backdrop-filter: blur(20px);
}

.items-list {
  flex: 1;
  min-width: 300px;
}

SettingsCard {
  flex: 0 0 360px;
}
</style>
