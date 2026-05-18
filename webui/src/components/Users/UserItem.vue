<script>
export default {
  props: {
    user: Object
  },
  emits: ["select"]
}
</script>

<template>
  <div class="user-item" @click="$emit('select', user)">

    <div class="user-photo-wrapper">
      <img :src="user.photoUrl" class="user-photo" />
    </div>

    <div class="user-info">
      <div class="user-name">
        {{ user.name }}
      </div>
    </div>

  </div>
</template>

<style scoped>
.user-item {
  display: flex;
  align-items: center;
  gap: 14px;

  padding: 14px;
  margin-bottom: 10px;
  border-radius: 18px;

  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);

  cursor: pointer;
  overflow: hidden;

  transition: transform 0.25s ease, box-shadow 0.25s ease;
  animation: fadeInUp 0.35s ease both;
}

.user-item:hover {
  transform: translateY(-6px) scale(1.02);
  border: 1px solid rgba(0, 255, 120, 0.18);
  box-shadow: 0 18px 65px rgba(0, 0, 0, 0.75);
}

.user-photo-wrapper {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  overflow: hidden;
}

.user-photo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-name {
  font-size: 1.05rem;
  font-weight: 800;
  color: rgba(245, 245, 245, 0.92);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(12px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
</style>
