<script>
import MemberList from "@/components/Groups/MemberList.vue";
import Bottombar from "@/components/Shared/Bottombar.vue";
import Topbar from "@/components/Shared/Topbar.vue";
import Poller from "@/services/poller";
import SettingsCard from "@/components/Settings/SettingsCard.vue";
import ItemsList from "@/components/Users/ItemsList.vue";
import { usePhotoManager } from "@/composables/usePhotoManager";
import { useSettingsForm } from "@/composables/useSettingsForm";
import { defaultGroupPhotoUrl } from "@/assets/default";
import { handleError } from "@/utils/errors";
import { users } from "@/state/users";
import { user } from "@/state/user";
import { groups } from "@/state/conversations";
import {
  addToGroup,
  getMemberIds,
  removeUser,
  setGroupName,
  setGroupPhoto,
  deleteGroupPhoto,
  leaveGroup,
  deleteGroup,
} from "@/services/groups";
import { getIcon } from "@/state/theme";

export default {
  components: { Bottombar, Topbar, SettingsCard, MemberList, ItemsList },

  computed: {
    groups() {
      return groups.value;
    },

    groupId() {
      return Number(this.$route.params.id);
    },

    founderId() {
      return this.groups.get(this.groupId)?.founderId;
    },

    isFounder() {
      return this.founderId === user.userId;
    },

    groupName() {
      return this.groups.get(this.groupId)?.name;
    },

    groupPhotoUrl() {
      return this.groups.get(this.groupId)?.photoUrl;
    },
  },

  watch: {
    groupPhotoUrl: {
      immediate: true,
      handler(newUrl) {
        this.photoUrl = newUrl;
        this.photo = null;
        this.photoChanged = false;
      },
    },

    groupName: {
      immediate: true,
      handler(newName) {
        this.text = newName || "";
      },
    },

    groups: {
      handler(groups) {
        if (!groups.has(this.groupId)) {
          this.$router.push("/home");
        }
      },
      deep: true,
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

    async deleteGroup() {
      try {
        await deleteGroup(this.groupId);
        this.$notifier.success(`Group "${this.groupName}" deleted`);
        this.$router.push("/home");
      } catch (e) {
        handleError(e);
      }
    },

    async leaveGroup() {
      try {
        await leaveGroup(this.groupId);
        this.$notifier.success(`Group "${this.groupName}" left`);
        this.$router.push("/home");
      } catch (e) {
        handleError(e);
      }
    },

    async addMemberToGroup(user) {
      try {
        const newMember = await addToGroup(this.groupId, user.id);
        this.$notifier.success(`User ${newMember.name} added`);
        return user;
      } catch (e) {
        handleError(e);
      }
    },

    async addUsersGroup(users) {
      // Avoid UI flickering by sorting members by ID.
      const addedMembers = await Promise.all(users.map(this.addMemberToGroup));

      const newMembersList = [...addedMembers, ...this.members];
      newMembersList.sort((a, b) => a.userId - b.userId);

      this.members = newMembersList;
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

    async updateGroup() {
      try {
        await this.submit(async (name) => {
          let changed = false;

          if (name !== this.groupName) {
            const updatedName = await setGroupName(this.groupId, name);
            this.placeholder = updatedName;
            changed = true;
          }

          if (this.photo) {
            const url = await setGroupPhoto(this.groupId, this.photo);
            this.photUrl = url;
            this.photo = null;
            this.photoChanged = false;

            changed = true;
          } else if (this.photoChanged) {
            await deleteGroupPhoto(this.groupId);
            this.photoUrl = defaultGroupPhotoUrl;
            this.photo = null;
            this.photoChanged = false;

            changed = true;
          }

          if (!changed) throw new Error("NO_CHANGE");
        });

        this.$notifier.success("Group updated successfully");
      } catch (e) {
        if (e.message === "EMPTY_NAME")
          this.$notifier.error("Invalid group name");
        else if (e.message === "NO_CHANGE")
          this.$notifier.warning("Nothing to do..");
        else {
          handleError(e);
        }
      }
    },
  },

  async mounted() {
    this.poller = new Poller(async () => {
      const memberIds = await getMemberIds(this.groupId);
      this.members = users.value
        .values()
        .filter((u) => memberIds.some((m) => u.userId === m))
        .toArray();
    }, 5000);

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
      <div class="group-settings">
        <SettingsCard
          style="display: flex"
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

        <button class="submit-button" @click="leaveGroup">
          <img class="icon-img" :src="getIcon('leave')" />
          Leave
        </button>

        <button v-if="isFounder" class="submit-button" @click="deleteGroup">
          <img class="icon-img" :src="getIcon('trash')" />
          Delete
        </button>
      </div>

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
          <ItemsList
            @select="addUsersGroup"
            :excludedUsers="members"
            :includeUsers="true"
            :includeGroups="false"
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
  gap: 10px;
  align-items: flex-start;
  padding: 20px;

  max-width: 1300px;
  margin: 100px auto;
}

.group-settings {
  width: min(600px, 50%);
}

.items-list {
  width: min(720px, 50%);
  padding: 10px;
  border-radius: 22px;

  border: 1px solid var(--border);
  box-shadow: 0 25px 90px var(--shadow);
}

.items-list {
  flex: 1;
  min-width: 300px;
}

SettingsCard {
  flex: 0 0 360px;
}
</style>
