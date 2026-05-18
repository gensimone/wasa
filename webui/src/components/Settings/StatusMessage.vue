<script>
export default {
    props: {
        message: String,
        error: Boolean
    }
}
</script>

<template>
    <div v-if="message" class="status-area">
        <p :class="error ? 'error' : 'success'">
            {{ message }}
        </p>
    </div>
</template>

<style scoped>
.status-area {
  width: 100%;
  max-width: 520px;
  margin-top: 14px;

  display: flex;
  flex-direction: column;
  text-align: center;
}

.error {
  margin-top: 10px;
  color: rgba(255, 80, 80, 0.85);
  font-size: 0.9rem;
}

.success {
  margin-top: 10px;
  color: rgba(0, 255, 120, 0.7);
  font-size: 0.9rem;
}
</style>
