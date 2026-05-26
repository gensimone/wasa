const DEV = import.meta.env.DEV

function wrap(fn) {
    return (...args) => {
        if (DEV) fn(...args)
    }
}

const logger = {
    log: wrap(console.log),
    error: wrap(console.error),
    warn: wrap(console.warn),
    info: wrap(console.info),
    debug: wrap(console.debug),
    table: wrap(console.table),
    trace: wrap(console.trace),
    assert: (cond, ...args) => {
        if (DEV) console.assert(cond, ...args)
    }
}

export default logger
