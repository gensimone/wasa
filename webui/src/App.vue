<script>
import { RouterView } from "vue-router"
import NotificationContainer from "@/notifier/NotificationContainer.vue"
import ImageModal from "@/components/Shared/ImageModal.vue";

export default {
    components: { NotificationContainer, ImageModal }
}
</script>

<template>
    <div class="container-fluid">
        <RouterView />
    </div>
    <NotificationContainer />
    <ImageModal />
</template>
