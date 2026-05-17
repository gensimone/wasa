<script>
import backIcon from "@/assets/icons/back.svg"
export default {
    data() {
        return {
            backIcon
        }
    }
}
</script>

<template>
    <div class="app">
        <header class="topbar">
            <div class="header-title"> WASAText </div>
            <div class="actions">
                <button class="icon-btn" @click="$router.back()">
                    <img :src="backIcon" class="icon-img">
                </button>
            </div>
        </header>
        <LoggedAs />
    </div>
</template>

<style scoped></style>
