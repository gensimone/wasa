<script>
export default {
    props: {
        text: String,
        title: String
    },
    emits: ["keyPress"]
}
</script>

<template>
    <div class="text-editor-box">
        <h2> {{ title }} </h2>
        <div class="input-bar">
            <input class="prompt" :value="text" @input="$emit('keyPress', $event.target.value)" />
        </div>
    </div>
</template>

<style scoped>
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
    margin-bottom: 14px;
    font-size: 1.2rem;
    letter-spacing: 1px;
}
</style>
