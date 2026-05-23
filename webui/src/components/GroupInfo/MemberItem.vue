<script>
import { user } from "@/state/user"
import { group } from "@/state/group"
import { expandUrl } from "@/utils/media"

export default {
    props: {
        member: { type: Object, required: true }
    },

    computed: {
        memberIsFounder() {
            return this.member.userId === group.founderId
        },

        isFounder() {
            return user.userId === group.founderId
        }
    },

    emits: [
        "removeUser",
        "selectUser"
    ],

    methods: {
        expandUrl
    }
}
</script>

<template>
    <div class="list-item" @click="$emit('selectUser', member)">
        <div class="item-photo-wrapper">
            <img :src="expandUrl(member.photoUrl)" class="item-photo" />
        </div>
        <div class="item-info">
            <div class="item-name">
                {{ member.name }}
            </div>
        </div>
        <div v-if="memberIsFounder" class="is-founder">
            Founder
        </div>
        <div v-else-if="isFounder" class="remove-button">
            <button class="icon-btn" @click="$emit('removeUser', member)">
                <img src="/icons/remove2.svg" class="icon-img" />
            </button>
        </div>
    </div>
</template>

<style scoped>
.remove-button {
    position: absolute;
    right: 20px;
}

.is-founder {
    opacity: 50%;
    position: absolute;
    right: 20px;
}
</style>
