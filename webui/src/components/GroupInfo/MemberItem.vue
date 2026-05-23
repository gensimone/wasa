<script>
import { user } from "@/state/user"
import { group } from "@/state/group"
import { expandUrl } from "@/utils/media"

export default {
    props: {
        member: { type: Object, required: true }
    },

    emits: ["removeUser", "selectUser"],

    computed: {
        isFounder() {
            return this.member.userId === group.founderId
        },

        currentUserIsFounder() {
            return user.userId === group.founderId
        }
    },

    methods: {
        expandUrl,

        onSelect() {
            if (this.member.userId !== user.userId)
                this.$emit("selectUser", this.member)
        },

        onRemove(event) {
            event.stopPropagation()
            this.$emit("removeUser", this.member)
        }
    }
}
</script>

<template>
    <div class="list-item" @click="onSelect">

        <div class="item-photo-wrapper">
            <img :src="expandUrl(member.photoUrl)" class="item-photo" />
        </div>

        <div class="item-info">
            <div class="item-name">
                {{ member.name }}
            </div>
        </div>

        <div class="item-badge" v-if="isFounder">
            Founder
        </div>

        <div class="item-actions" v-if="currentUserIsFounder && !isFounder">
            <button class="icon-btn" @click="onRemove">
                <img src="/icons/minus.svg" class="icon-img" />
            </button>
        </div>

    </div>
</template>

<style scoped>
.icon-btn:hover {
    background: rgba(255, 0, 0, 0.08);
    border-color: rgba(255, 0, 0, 0.2);
}
</style>
