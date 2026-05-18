<script>
import { user } from "@/state/user"
export default {
    computed: {
        footerText() {
            if (user.name) {
                return `Logged as ${user.name}`
            }

            return "Made by Simone Gentili"
        }
    }
}
</script>

<template>
    <footer class="footer">
        {{ footerText }}
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

    color: rgba(255, 255, 255, 0.6);
    text-shadow: 0 0 6px rgba(0, 255, 65, 0.3);

    z-index: 9999;

    background: rgba(0, 0, 0, 0.35);
    backdrop-filter: blur(10px);
    border-top: 1px solid rgba(255, 255, 255, 0.06);
}
</style>
