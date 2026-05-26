import { getNotifier } from "@/notifier";
import logger from "@/utils/logger";

export function handleError(e) {
    const notifier = getNotifier()

    if (e.request) {
        notifier.error("Network error")

    } else {
        notifier.error("Unexpected error")
    }

    logger.error(e)
}
