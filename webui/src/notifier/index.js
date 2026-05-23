import { createApp } from "vue"
import NotificationContainer from "./NotificationContainer.vue"

let appInstance = null
let notifier = null

export default {
    install(app) {
        appInstance = app

        const api = {
            success(message, duration = 3000) {
                notifier?.({ message, type: "success", duration })
            },

            error(message, duration = 10000) {
                notifier?.({ message, type: "error", duration })
            },

            warning(message, duration = 7500) {
                notifier?.({ message, type: "warning", duration })
            },

            info(message, duration = 3000) {
                notifier?.({ message, type: "info", duration })
            },

            show(options) {
                notifier?.(options)
            }
        }

        app.config.globalProperties.$notifier = api

        const container = document.createElement("div")
        document.body.appendChild(container)

        const notifierApp = createApp(NotificationContainer)
        const vm = notifierApp.mount(container)

        notifier = vm.add
    }
}

export function getNotifier() {
    return appInstance?.config?.globalProperties?.$notifier
}
