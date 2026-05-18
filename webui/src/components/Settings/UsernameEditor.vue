<script>
export default {
    props: {
        modelValue: String
    },
    emits: ["update:modelValue", "save"]
}
</script>

<template>
    <div class="username-box">
        <h2>Username</h2>
        <input class="prompt" :value="modelValue" @input="$emit('update:modelValue', $event.target.value)"
            @keyup.enter="$emit('save')" />
    </div>
</template>

<style scoped>
.username-box {
  width: 100%;
  text-align: center;
}

.username-box h2 {
  margin-bottom: 14px;
  font-size: 1.2rem;
  letter-spacing: 1px;
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
</style>
