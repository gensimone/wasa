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
    <div class="member-item" @click="onSelect">

        <div class="member-item-photo-wrapper">
            <img :src="expandUrl(member.photoUrl)" class="member-item-photo" />
        </div>

        <div class="member-item-info">
            <div class="member-item-name">
                {{ member.name }}
            </div>
        </div>

        <div class="member-item-badge" v-if="isFounder">
            Founder
        </div>

        <div class="member-item-actions" v-if="currentUserIsFounder && !isFounder">
            <button class="icon-btn" @click="onRemove">
                <img src="/icons/minus.svg" class="icon-img" />
            </button>
        </div>

    </div>
</template>

<style scoped>
.member-item {
    display: flex;
    align-items: center;
    gap: 14px;

    padding: 14px;
    margin-bottom: 10px;

    border-radius: 18px;

    background: var(--surface-2);
    border: 1px solid var(--border);

    cursor: pointer;
    overflow: hidden;

    animation: fadeInUp 0.35s ease both;
}

.member-item:hover {
    transform: translateY(-6px) scale(1.02);
    border: 1px solid var(--accent);
}

.member-item-photo-wrapper {
    width: 75px;
    height: 75px;
    border-radius: 16px;
    overflow: hidden;
    flex-shrink: 0;
}

.member-item-photo {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.member-item-info {
    display: flex;
    flex-direction: column;
    flex: 1;
    min-width: 0;
}

.member-item-name {
    font-size: 1.05rem;
    font-weight: 800;

    color: var(--text);

    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.member-item-badge {
    position: absolute;
    right: 16px;

    font-size: 0.85rem;
    font-weight: 700;

    opacity: 0.6;
    color: var(--text);
}

.member-item-actions {
    position: absolute;
    right: 16px;

    display: flex;
    align-items: center;
    justify-content: center;
}


.icon-btn:hover {
    background: rgba(255, 0, 0, 0.08);
    border-color: rgba(255, 0, 0, 0.2);
}
</style>
