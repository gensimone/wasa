export class Poller {
    constructor(callback, interval = 3000) {
        this.callback = callback
        this.interval = interval
        this.poller = null
        this.lock = false
        this.error = null
    }

    startPolling() {
        this.stopPolling()

        this.poller = setInterval(async () => {
            if (this.lock) return
            this.lock = true

            try {
                this.callback()
            } catch (e) {
                this.error = e
            } finally {
                this.lock = false
            }
        }, this.interval)
    }

    stopPolling() {
        if (this.poller) {
            clearInterval(this.poller)
            this.poller = null
        }
    }
}
