<script>
import { RouterView } from 'vue-router'
import NotificationContainer from "./notifier/NotificationContainer.vue"

export default {
    components: { NotificationContainer }
}
</script>

<template>
    <div class="container-fluid">
        <RouterView />
    </div>
    <NotificationContainer />
</template>
