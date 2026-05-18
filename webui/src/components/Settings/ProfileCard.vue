<script>
import AvatarEditor from "./AvatarEditor.vue"
import UsernameEditor from "./UsernameEditor.vue"
import StatusMessage from "./StatusMessage.vue"

export default {
    components: {
        AvatarEditor,
        UsernameEditor,
        StatusMessage
    },

    props: {
        avatarUrl: String,
        name: String,
        isDefault: Boolean,
        photoChanged: Boolean,
        loading: Boolean,
        message: String,
        error: Boolean
    },

    emits: [
        "uploadPhoto",
        "removePhoto",
        "savePhoto",
        "updateName",
        "saveName"
    ]
}
</script>

<template>
    <div class="profile-card">

        <AvatarEditor :avatarUrl="avatarUrl" :isDefault="isDefault" :photoChanged="photoChanged" :loading="loading"
            @upload="$emit('uploadPhoto', $event)" @remove="$emit('removePhoto')" @save="$emit('savePhoto')" />

        <UsernameEditor :modelValue="name" @update:modelValue="$emit('updateName', $event)" @save="$emit('saveName')" />

        <StatusMessage :message="message" :error="error" />

    </div>
</template>

<style scoped>
.profile-card {
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
