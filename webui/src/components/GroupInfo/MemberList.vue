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
        "removeUser", // Remove user from group (founder only).
        "addToGroup"  // Add user to group      (any member).
    ],
}
</script>

<template>
    <div class="items-list">
        <div class="item-list-header">
            <h2> Members </h2>
            <button class="icon-btn" @click="$router.push('/conversation/add')">
                <img src="/icons/plus.svg" class="icon-img">
            </button>
        </div>
        <MemberItem v-for="m in members" :key="m.userId" :member="m" @removeUser="$emit('removeUser', $event)" />
    </div>
</template>
