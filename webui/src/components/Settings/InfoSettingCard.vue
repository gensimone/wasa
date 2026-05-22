<script>
import PhotoEditor from "@/components/Settings/PhotoEditor.vue"
import StatusMessage from "@/components/Shared/StatusMessage.vue"
import TextEditor from "@/components/Settings/TextEditor.vue"

export default {
    components: {
        PhotoEditor,
        StatusMessage,
        TextEditor
    },

    props: {
        photoUrl: String,
        photoChanged: Boolean,

        title: String,
        text: String,

        loading: Boolean,
        message: String,
        error: Boolean
    },

    emits: [
        "uploadPhoto",
        "deletePhoto",
        "revertPhoto",
        "keyPress",
        "save"
    ]
}
</script>

<template>
    <div class="item-setting-card">

        <!-- PHOTO EDITOR-->
        <PhotoEditor :photoUrl="photoUrl" :photoChanged="photoChanged" :loading="loading"
            @uploadPhoto="$emit('uploadPhoto', $event)" @revertPhoto="$emit('revertPhoto')"
            @deletePhoto="$emit('deletePhoto')" />

        <!-- TEXT ENTRY EDITOR -->
        <TextEditor :title="title" :text="text" @keyPress="$emit('keyPress', $event)" />

        <!-- SAVE BUTTON -->
        <button class="icon-btn" @click="$emit('save')">
            <img src="/icons/save.svg" class="icon-img">
        </button>

        <!-- STATUS MESSAGE -->
        <StatusMessage :message="message" :error="error" />

    </div>
</template>

<style scoped>
.item-setting-card {
    width: min(520px, 92%);
    padding: 34px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 30px;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.06);
    border-radius: 22px;
    backdrop-filter: blur(20px);
    box-shadow: 0 25px 90px rgba(0, 0, 0, 0.75);
}
</style>
