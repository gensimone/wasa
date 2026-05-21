<script>
import UserItem from "@/components/Users/UserItem.vue"

export default {
    components: {
        UserItem
    },
    props: {
        users: Array,
        loading: Boolean // FIXME: implement loading screen.
    },
    emits: ["select"]
}
</script>

<template>
    <div class="items-list">
        <div class="items-list-header">
            <h2>Users</h2>
        </div>
        <UserItem v-for="u in users" :key="u.userId" :user="u" @select="$emit('select', $event)" />
    </div>
</template>
