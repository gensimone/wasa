<script>
import PhotoEditor from "@/components/Settings/PhotoEditor.vue"

export default {
    components: { PhotoEditor },

    props: {
        photoUrl: { type: String, required: true },
        photoChanged: { type: Boolean, required: true },
        text: { type: String, required: true },
        title: { type: String, required: true },
        submitButtonText: { type: String, required: true },
        loading: { type: Boolean, required: true }
    },

    emits: [
        "submit",
        "uploadPhoto",
        "revertPhoto",
        "deletePhoto",
        "keyPress"
    ]
}
</script>

<template>
    <div class="setting-card">
        <form @submit.prevent="$emit('submit')">

            <PhotoEditor :photoUrl="photoUrl" :photoChanged="photoChanged" :loading="loading"
                @uploadPhoto="$emit('uploadPhoto', $event)" @revertPhoto="$emit('revertPhoto')"
                @deletePhoto="$emit('deletePhoto')" />

            <div class="setting-card-input-box">
                <h2>{{ title }}</h2>

                <input class="input-bar" :placeholder="text" @input="$emit('keyPress', $event.target.value)" />
            </div>

            <button class="submit-button" :disabled="loading">
                {{ submitButtonText }}
            </button>

        </form>
    </div>
</template>

<style scoped>
.setting-card {
    width: min(520px, 92%);
    padding: 34px;
    display: flex;
    flex-direction: column;
    gap: 30px;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.06);
    border-radius: 22px;
}

.setting-card-input-box {
    width: 100%;
    text-align: center;
}

.setting-card-input-box h2 {
    margin-top: 35px;
    margin-bottom: 14px;
    font-size: 1.2rem;
    letter-spacing: 1px;
}
</style>
