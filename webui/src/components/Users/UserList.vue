<script>
import UserItem from "@/components/Users/UserItem.vue"

export default {
    components: {
        UserItem
    },

    props: {
        users: { type: Array, required: true },
        loading: { type: Boolean, required: true }
    },

    emits: ["select"]
}
</script>

<template>
    <div class="items-list">
        <div class="item-list-header">
            <h2>Users</h2>
        </div>
        <UserItem v-for="u in users" :key="u.userId" :user="u" @select="$emit('select', $event)" />
    </div>
</template>
