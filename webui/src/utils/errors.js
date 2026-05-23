import { getNotifier } from "@/notifier";
import { logger } from "@/utils/logger";

export function handleError(e) {
    const notifier = getNotifier()

    if (e.response) {
        notifier.error(e.response.data.error)

    } else if (e.request) {
        notifier.error("Network error")

    } else {
        notifier.error("Unexpected error")
    }

    logger.error(e)
}
