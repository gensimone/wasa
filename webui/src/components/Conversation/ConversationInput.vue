<script>
export default {
    props: {
        text: { type: String, default: "" },
        attachment: { type: File, default: null },
        attachmentUrl: { type: String, default: null },
        sending: { type: Boolean, required: true }
    },

    emits: [
        "update:text",
        "send",
        "addAttachment",
        "removeAttachment"
    ],
}
</script>

<template>
    <div>
        <div v-if="attachment" class="attachment-preview">
            <img class="preview-img" :src="attachmentUrl" />
            <button class="icon-btn" @click="$emit('removeAttachment')">
                <img src="/icons/remove.svg" class="icon-img">
            </button>
        </div>

        <div class="input-bar">
            <button class="icon-btn" @click="$refs.fileInput.click()">
                <img src="/icons/upload.svg" class="icon-img">
            </button>

            <input ref="fileInput" type="file" accept="image/*" style="display: none"
                @change="$emit('addAttachment', $event)" />

            <input class="prompt" :value="text" :disabled="sending" @input="$emit('update:text', $event.target.value)"
                placeholder="Type a message..." @keyup.enter="$emit('send')" />

            <button class="icon-btn" :disable="sending" @click="$emit('send')">
                <img src="/icons/send.svg" class="icon-img">
            </button>
        </div>
    </div>
</template>

<style scoped>
.input-bar {
    width: min(720px, 100%);
    margin: 0px auto 0;
    display: flex;
    gap: 10px;
    padding: 10px;
    border-radius: 0px;
    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.08);
}

.prompt {
    width: 100%;
    padding: 14px;
    border-radius: 14px;
    background: rgba(0, 0, 0, 0.5);
    border: 1px solid rgba(255, 255, 255, 0.06);
    color: rgba(245, 245, 245, 0.92);
    outline: none;
    transition: all 0.25s ease;
}

.prompt:focus {
    border: 1px solid rgba(0, 255, 120, 0.25);
    box-shadow: 0 0 20px rgba(0, 255, 120, 0.08);
}

.attachment-preview {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin: 0px auto 0;
    padding: 8px 10px;
    width: min(720px, 100%);
    border-radius: 0px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
    backdrop-filter: blur(10px);
}

.preview-img {
    width: 50px;
    height: 50px;
    border-radius: 12px;
    object-fit: cover;
    border: 1px solid rgba(255, 255, 255, 0.12);
}
</style>
