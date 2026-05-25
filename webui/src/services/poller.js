import { handleError } from "@/utils/errors";

export class Poller {
    constructor(
        callback,
        interval = 3000,
        maxInterval = 10000,
        factor = 1.5
    ) {
        this.callback = callback
        this.running = false

        this.timer = null
        this.baseInterval = interval
        this.currentInterval = interval
        this.maxInterval = maxInterval
        this.factor = factor
    }

    async loop() {
        if (!this.running) return

        if (document.visibilityState === "visible") {
            try {
                await this.callback()
                this.currentInterval = this.baseInterval

            } catch (e) {
                handleError(e)

                this.currentInterval = Math.min(
                    this.currentInterval * this.factor,
                    this.maxInterval
                )
            }
        }

        this.timer = setTimeout(() => this.loop(), this.currentInterval)
    }

    startPolling() {
        this.stopPolling()
        this.running = true
        this.loop()
    }

    stopPolling() {
        this.running = false
        if (this.timer) {
            clearTimeout(this.timer)
            this.timer = null
        }
    }
}
