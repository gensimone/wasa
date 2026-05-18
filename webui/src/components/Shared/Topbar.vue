<script>
export default {
    props: {
        title: {
            String,
            default: "Wasatext"
        },
        actions: {
            type: Array,
            default: () => []
        },
        links: {
            type: Array,
            default: () => []
        }
    }
}
</script>

<template>
    <header class="topbar">
        <div class="header-title">{{ title }}</div>
        <div class="actions">
            <a v-for="(l, i) in links" :key="i" class="icon-btn" :href="l.link" target="_blank">
                <img :src="l.icon" class="icon-img">
            </a>
            <button v-for="(a, i) in actions" :key="i" class="icon-btn" @click="a.onClick">
                <img :src="a.icon" class="icon-img" />
            </button>
        </div>
    </header>
</template>

<style scoped>
.topbar {
    left: 0;
    position: fixed;
    width: 100%;
    top: 0px;
    height: 65px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 22px;
    z-index: 9999;
    background: rgba(0, 0, 0, 0.35);
    backdrop-filter: blur(10px);
    border-bottom: 1px solid rgba(255, 255, 255, 0.08);
    box-shadow: 0 18px 70px rgba(0, 0, 0, 0.7);

}

.header-title {
    font-size: 1.45rem;
    font-weight: 900;
    letter-spacing: 2px;
    text-transform: uppercase;
    color: rgba(245, 245, 245, 0.92);
    text-shadow: 0 0 12px rgba(0, 255, 120, 0.18);
}

.actions {
    display: flex;
    gap: 12px;
}
</style>
