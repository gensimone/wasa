<script setup>
    import { userState } from "@/state/user"
</script>

<template>
  <footer class="footer">
    Logged as {{ userState.name }}
  </footer>
</template>

<style>
.footer {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    padding: 5px;
    text-align: center;
    font-size: 12px;
    color: rgba(0, 255, 65, 0.6);
    text-shadow: 0 0 6px rgba(0, 255, 65, 0.3);
}
</style>
