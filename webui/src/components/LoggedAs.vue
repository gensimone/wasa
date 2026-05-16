<script setup>
import { userState } from "@/state/user"
</script>

<template>
    <footer class="footer">
        Logged as {{ userState.name }}
    </footer>
</template>
