<script>
import PhotoEditor from "@/components/Settings/PhotoEditor.vue"

export default {
    components: { PhotoEditor },

    props: {
        photoUrl: { type: String, required: true },
        photoChanged: { type: Boolean, required: true },
        submitButtonText: { type: String, required: true },
        title: { type: String, required: true },
        text: { type: String, required: true },
        loading: { type: Boolean, required: true }
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
            <PhotoEditor :photoUrl="photoUrl" :photoChanged="photoChanged" :loading="loading"
                @uploadPhoto="$emit('uploadPhoto', $event)" @revertPhoto="$emit('revertPhoto')"
                @deletePhoto="$emit('deletePhoto')" />

            <div class="text-editor-box">
                <h2> {{ title }} </h2>
                <div class="input-bar">
                    <input class="prompt" :placeholder="text" @input="$emit('keyPress', $event.target.value)" />
                </div>
            </div>

            <button class="submit-button" :disabled="loading">
                {{ submitButtonText }}
            </button>
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
