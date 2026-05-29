<script>
import Item from "@/components/Users/Item.vue";

export default {
  components: { Item },

  props: {
    items: { type: Array, required: true },
    canSelectMultiple: { type: Boolean, required: true },
  },

  emits: ["select"],

  data() {
    return {
      selectedItems: [],
      query: "",
    };
  },

  computed: {
    itemsToShow() {
      const unselectedItems = this.items.filter((i) => !this.isSelected(i));

      if (!this.query.trim()) return unselectedItems;

      return unselectedItems.filter((i) =>
        i.name.toLowerCase().includes(this.query.toLowerCase()),
      );
    },
  },

  methods: {
    onSelect(item) {
      if (!this.canSelectMultiple) {
        this.$emit("select", item);
        return;
      }

      if (this.isSelected(item)) {
        this.selectedItems = this.selectedItems.filter(
          (i) => i.id !== item.id || i.isGroup !== item.isGroup,
        );
      } else {
        this.selectedItems.push(item);
      }
    },

    confirmSelection() {
      this.$emit("select", this.selectedItems);
    },

    isSelected(item) {
      return this.selectedItems.some(
        (i) => i.id === item.id && i.isGroup === item.isGroup,
      );
    },
  },
};
</script>

<template>
  <div class="item-picker">
    <input
      name="search-item"
      class="input-bar"
      placeholder="Search.."
      @input="query = $event.target.value"
    />
    <div class="item-picker-list">
      <Item
        v-for="item in itemsToShow"
        :key="item.isGroup ? `g-${item.id}` : `u-${item.id}`"
        :item="{
          name: item.name,
          photoUrl: item.photoUrl,
          id: item.id,
          isGroup: item.isGroup,
        }"
        :selected="isSelected(item)"
        @select="onSelect"
      />
    </div>

    <div
      v-if="canSelectMultiple && !query.trim().length"
      class="item-picker-sidebar"
    >
      <div v-if="this.selectedItems.length" class="selected-title">
        Selected
      </div>

      <Item
        v-for="item in selectedItems"
        :key="item.isGroup ? `g-${item.id}` : `u-${item.id}`"
        :item="{
          name: item.name,
          photoUrl: item.photoUrl,
          id: item.id,
          isGroup: item.isGroup,
        }"
        :selected="isSelected(item)"
        @select="onSelect"
      />

      <button
        v-if="selectedItems.length"
        class="submit-button"
        @click="confirmSelection"
      >
        Confirm
      </button>
    </div>
  </div>
</template>

<style scoped>
.item-picker {
  width: min(600px, 100%);
}

.selected-title {
  font-size: 1.5rem;
  color: var(--text);
  text-transform: uppercase;
}
</style>
