<script>
import { expandUrl } from "@/utils/media"
export default {
    props: {
        user: { type: Object, required: true }
    },

    emits: ["select"],

    methods: {
        expandUrl
    }
}
</script>

<template>
    <div class="list-item" @click="$emit('select', user)">
        <div class="item-photo-wrapper">
            <img :src="expandUrl(user.photoUrl)" class="item-photo" />
        </div>
        <div class="item-info">
            <div class="item-name">
                {{ user.name }}
            </div>
        </div>
    </div>
</template>
