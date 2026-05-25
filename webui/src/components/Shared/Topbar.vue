<script>
export default {
    props: {
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
        <button class="topbar-title" @click="$router.push('/home')">
            WASAText
        </button>
        <div class="topbar-actions">
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
    position: fixed;
    top: 0;
    left: 0;

    width: 100%;
    height: 65px;

    display: flex;
    justify-content: space-between;
    align-items: center;

    padding: 0 22px;

    background: var(--surface);
    border-bottom: 1px solid var(--border);
    box-shadow: 0 18px 70px var(--shadow);

    backdrop-filter: blur(10px);

    z-index: 2;
}

.topbar-actions {
    display: flex;
    gap: 12px;
}

.topbar-title {
    font-size: 1.45rem;
    font-weight: 0;
    border: 0px;
    background: transparent;
    color: var(--text);
}
</style>
