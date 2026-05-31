<script>
import Poller from "@/services/poller";
import { getIcon } from "@/state/theme";
import { getCheckIcon } from "@/utils/messages";
import { getReceipts } from "@/services/messages";

export default {
  props: {
    messageId: { type: Number, required: true },
  },

  data() {
    return {
      poller: null,
      checkIcon: "check-sent",
    };
  },

  methods: {
    getIcon,
  },

  mounted() {
    this.poller = new Poller();

    this.poller.callback = async () => {
      const receipts = await getReceipts(this.messageId);
      const icon = getCheckIcon(receipts);
      if (!icon || icon === "check-read") this.poller.stopPolling();

      this.checkIcon = icon;
    };

    this.poller.startPolling();
  },

  beforeUnmount() {
    this.poller?.stopPolling();
  },
};
</script>

<template>
  <img class="message-item-check-icon" :src="getIcon(checkIcon)" alt="" />
</template>

<style scoped>
.message-item-check-icon {
  width: 16px;
  height: 16px;
  opacity: 0.9;
}
</style>
