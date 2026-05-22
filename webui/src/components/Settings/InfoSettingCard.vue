<script>
import PhotoEditor from "@/components/Settings/PhotoEditor.vue"
import StatusMessage from "@/components/Shared/StatusMessage.vue"

export default {
    components: {
        PhotoEditor,
        StatusMessage
    },

    props: {
        photoUrl: String,
        photoChanged: Boolean,
        submitButtonText: String,

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
        "submit"
    ]
}
</script>

<template>
    <div class="item-setting-card">
        <form @submit.prevent="$emit('submit')">

            <!-- PHOTO EDITOR-->
            <PhotoEditor :photoUrl="photoUrl" :photoChanged="photoChanged" :loading="loading"
                @uploadPhoto="$emit('uploadPhoto', $event)" @revertPhoto="$emit('revertPhoto')"
                @deletePhoto="$emit('deletePhoto')" />

            <!-- INPUT BOX -->
            <div class="text-editor-box">
                <h2> {{ title }} </h2>
                <div class="input-bar">
                    <input class="prompt" :placeholder="text" @input="$emit('keyPress', $event.target.value)"
                        />
                </div>
            </div>

            <!-- SAVE BUTTON -->
            <button class="submit-button" :disabled="loading">
                {{ submitButtonText }}
            </button>

            <!-- STATUS MESSAGE -->
            <StatusMessage :message="message" :error="error" />
        </form>
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

.input-bar {
    display: flex;
    gap: 10px;
    border-radius: 0px;
}

.text-editor-box {
    width: 100%;
    text-align: center;
}

.text-editor-box h2 {
    margin-top: 35px;
    margin-bottom: 14px;
    font-size: 1.2rem;
    letter-spacing: 1px;
}
</style>
