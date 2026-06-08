<script>
import PhotoEditor from "@/components/Settings/PhotoEditor.vue";

export default {
  components: { PhotoEditor },

  props: {
    photoUrl: { type: String, required: false },
    photoChanged: { type: Boolean, required: true },
    title: { type: String, required: true },
    loading: { type: Boolean, required: true },
    text: { type: String, required: true },
    submitButtonText: { type: String, required: true },
    enableEditing: { type: Boolean, required: true },
  },

  emits: ["submit", "uploadPhoto", "revertPhoto", "deletePhoto", "keyPress"],
};
</script>

<template>
  <div class="setting-card">
    <PhotoEditor
      :photo-url="photoUrl"
      :enable-editing="enableEditing"
      :photo-changed="photoChanged"
      :loading="loading"
      @upload-photo="$emit('uploadPhoto', $event)"
      @revert-photo="$emit('revertPhoto')"
      @delete-photo="$emit('deletePhoto')"
    />

    <div class="setting-card-input-box">
      <h2>{{ title }}</h2>

      <input
        name="settings-card-input-bar"
        class="input-bar"
        :placeholder="text"
        :disabled="!enableEditing"
        @input="$emit('keyPress', $event.target.value)"
      >
    </div>

    <button
      v-if="enableEditing"
      class="submit-button"
      :disabled="loading"
      @click="$emit('submit')"
    >
      {{ submitButtonText }}
    </button>
  </div>
</template>

<style scoped>
.setting-card {
  width: min(600px, 100%);
  padding: 34px;
  display: flex;
  flex-direction: column;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 22px;
}

.setting-card-input-box {
  width: 100%;
  text-align: center;
}

.setting-card-input-box h2 {
  margin-bottom: 10px;
  font-size: 1.2rem;
  letter-spacing: 1px;
}
</style>
