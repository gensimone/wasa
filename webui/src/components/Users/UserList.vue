<script>
import UserItem from "./UserItem.vue"

export default {
    components: {
        UserItem
    },

    props: {
        users: Array
    },

    emits: ["select"]
}
</script>

<template>
    <div class="user-list">

        <div class="user-list-header">
            <h2>Users</h2>
        </div>

        <UserItem v-for="u in users" :key="u.userId" :user="u" @select="$emit('select', $event)" />

    </div>
</template>

<style scoped>
.user-list {
    position: relative;
    z-index: 1;
    width: min(720px, 100%);
    border-radius: 22px;
    padding: 20px;

    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(20px);

    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.7);

    overflow: hidden;
}

.user-list-header {
    margin-bottom: 18px;
    padding-bottom: 12px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

h2 {
    margin: 0;
    font-size: 1.15rem;
    font-weight: 800;
}
</style>
