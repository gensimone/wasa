let appInstance = null;
let notifier = null;

export default {
  install(app) {
    appInstance = app;

    const api = {
      success(text, duration = 3000) {
        notifier?.({ text, type: "success", duration });
      },

      error(text, duration = 10000) {
        notifier?.({ text, type: "error", duration });
      },

      warning(text, duration = 7500) {
        notifier?.({ text, type: "warning", duration });
      },

      info(text, duration = 3000) {
        notifier?.({ text, type: "info", duration });
      },

      message(opts) {
        notifier?.({ ...opts, type: "message" });
      },
    };

    app.config.globalProperties.$notifier = api;
  },
};

export function getNotifier() {
  return appInstance?.config?.globalProperties?.$notifier;
}

export function setNotifier(fn) {
  notifier = fn;
}
