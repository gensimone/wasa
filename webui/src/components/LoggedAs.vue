<script setup>
import { user } from "@/state/user"
</script>

<template>
    <footer class="footer">
        Logged as {{ user.name }}
    </footer>
</template>
