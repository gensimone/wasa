<script>
import { user } from "@/state/user";

export default {
  computed: {
    footerText() {
      if (user.userId) {
        return `Logged as ${user.name}`;
      }

      return "Made by Simone Gentili";
    },
  },
};
</script>

<template>
  <footer class="bottombar">
    {{ footerText }}
  </footer>
</template>

<style scoped>
.bottombar {
  position: fixed;
  bottom: 0;
  left: 0;

  width: 100%;
  padding: 5px;

  text-align: center;
  font-size: 12px;

  color: var(--text-muted);

  background: var(--surface);
  border-top: 1px solid var(--border);
  backdrop-filter: blur(10px);

  z-index: 1;
}
</style>
