<script>
import { getIcon } from "@/state/theme";
import { expandUrl } from "@/utils/media";
import { defaultGroupPhotoUrl, defaultUserPhotoUrl } from "@/assets/default";
import { setImageModal } from "@/state/imageModal";

export default {
  props: {
    photoUrl: { type: String, required: true },
    photoChanged: { type: Boolean, required: true },
    loading: { type: Boolean, required: true },
    enableEditing: { type: Boolean, required: true },
  },

  emits: ["revertPhoto", "deletePhoto", "uploadPhoto"],

  methods: {
    expandUrl,
    getIcon,
    setImageModal,
  },

  computed: {
    isDefault() {
      return (
        this.photoUrl === defaultUserPhotoUrl ||
        this.photoUrl === defaultGroupPhotoUrl
      );
    },
  },
};
</script>

<template>
  <div class="photo-editor">
    <div v-if="enableEditing">
      <button
        v-if="photoChanged"
        class="icon-btn"
        :disabled="loading"
        @click="$emit('revertPhoto')"
      >
        <img :src="getIcon('revert')" class="icon-img" />
      </button>
      <button
        v-else-if="!isDefault"
        class="icon-btn"
        :disabled="loading"
        @click="$emit('deletePhoto')"
      >
        <img :src="getIcon('trash')" class="icon-img" />
      </button>
      <button v-else class="icon-btn invisible-placeholder" disabled></button>
    </div>

    <label class="photo-editor-photo-wrapper">
      <img
        :src="expandUrl(photoUrl)"
        class="avatar-big"
        @click="setImageModal(photoUrl)"
      />
    </label>

    <label v-if="enableEditing" class="icon-btn">
      <img :src="getIcon('plus')" class="icon-img" />

      <input
        type="file"
        accept="image/*"
        hidden
        @change="$emit('uploadPhoto', $event)"
        :disabled="!enableEditing"
      />
    </label>
  </div>
</template>

<style scoped>
.photo-editor {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
}

.photo-editor-photo-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.invisible-placeholder {
  visibility: hidden;
}

.avatar-big {
  width: 170px;
  height: 170px;

  border-radius: 50%;
  object-fit: cover;

  border: 2px solid var(--accent);
}
</style>
