<script>
import { group } from "@/state/group"
import { expandUrl } from "@/utils/media"
export default {
    props: {
        member: { type: Object, required: true }
    },

    computed: {
        isFounder() {
            return this.member.userId === group.founderId
        }
    },

    emits: ["select"],

    methods: {
        expandUrl
    }
}
</script>

<template>
    <div class="list-item" @click="$emit('select', member)">
        <div class="item-photo-wrapper">
            <img :src="expandUrl(member.photoUrl)" class="item-photo" />
        </div>
        <div class="item-info">
            <div class="item-name">
                {{ member.name }}
            </div>
            <div v-if="isFounder">
                Founder
            </div>
        </div>
    </div>
</template>
