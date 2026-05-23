import { createApp } from "vue"
import NotificationContainer from "./NotificationContainer.vue"

let appInstance = null

export default {
    install(app) {
        appInstance = app
        let notifier = null

        app.config.globalProperties.$notifier = {
            success(message, duration = 5000) {
                notifier?.({
                    message,
                    error: false,
                    duration
                })
            },

            error(message, duration = 5000) {
                notifier?.({
                    message,
                    error: true,
                    duration
                })
            },

            show(options) {
                addToast?.(options)
            }
        }

        const container = document.createElement("div")
        document.body.appendChild(container)

        const notifierApp = createApp(NotificationContainer)
        const vm = notifierApp.mount(container)

        notifier = vm.add
    }
}

export function getNotifier() {
    return appInstance.config.globalProperties.$notifier
}
