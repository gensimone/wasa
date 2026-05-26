<script>
import UserItem from "@/components/Users/UserItem.vue"
import Poller from "@/services/poller";
import { getUsers } from "@/services/users";
import { user } from "@/state/user";

export default {
    components: { UserItem },

    data() {
        return {
            query: "",
            users: [],
            poller: null
        }
    },

    computed: {
        usersToShow() {
            if (!this.query.trim()) return this.users

            return this.users.filter(u =>
                u.name.toLowerCase().includes(this.query.toLowerCase())
            )
        }
    },

    emits: ["select"],

    async mounted() {
        this.poller = new Poller(async () => {
            const users = await getUsers()
            this.users = users
                .filter(u => u.userId != user.userId)
        })

        this.poller.startPolling()
    },

    beforeUnmount() {
        this.poller?.stopPolling()
    }
}
</script>

<template>
    <div class="users-list">
        <input class="input-bar" placeholder="Search.." @input="query = $event.target.value" />
        <UserItem v-for="u in usersToShow" :key="u.userId" :user="u" @select="$emit('select', $event)" />
    </div>
</template>

<style scoped>
.users-list {
    width: min(720px, 100%);
    padding: 20px;
    border-radius: 22px;

    background: var(--surface);
    border: 1px solid var(--border);
    box-shadow: 0 25px 90px var(--shadow);

    backdrop-filter: blur(20px);
}
</style>
