<script>
import MemberItem from "@/components/GroupInfo/MemberItem.vue"

export default {
    components: {
        MemberItem
    },

    props: {
        members: { type: Array, required: true }
    },

    emits: [
        "removeUser",
        "addToGroup",
        "selectUser"
    ],
}
</script>

<template>
    <div class="items-list">

        <div class="items-list-header">
            <h2>Members</h2>

            <button class="icon-btn" @click="$router.push('/conversation/add')">
                <img src="/icons/plus.svg" class="icon-img" />
            </button>
        </div>

        <MemberItem v-for="m in members" :key="m.userId" :member="m" @removeUser="$emit('removeUser', $event)"
            @selectUser="$emit('selectUser', $event)" />

    </div>
</template>

<style scoped>
.items-list {
    position: relative;
    z-index: 1;

    width: min(720px, 100%);
    padding: 20px;

    border-radius: 22px;

    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(20px);

    border: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.7);

    overflow: hidden;
}

.items-list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    margin-bottom: 18px;
    padding-bottom: 12px;

    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.items-list-header h2 {
    margin: 0;

    font-size: 1.15rem;
    font-weight: 800;
    letter-spacing: 1px;

    color: rgba(245, 245, 245, 0.9);
}

.items-list>*:not(.items-list-header) {
    margin-top: 10px;
}
</style>
