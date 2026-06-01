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
          (i) => i.id !== item.id || i.isDirect !== item.isDirect,
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
        (i) => i.id === item.id && i.isDirect === item.isDirect,
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
        :key="item.isDirect ? `d-${item.id}` : `g-${item.id}`"
        :item="{
          name: item.name,
          photoUrl: item.photoUrl,
          id: item.id,
          isDirect: item.isDirect,
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
        :key="item.isDirect ? `d-${item.id}` : `g-${item.id}`"
        :item="{
          name: item.name,
          photoUrl: item.photoUrl,
          id: item.id,
          isDirect: item.isDirect,
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

.item-picker {
  width: min(600px, 100%);

  display: flex;
  flex-direction: column;

  max-height: 70vh;
  min-height: 0;
}

.input-bar {
  flex-shrink: 0;
}

.item-picker-list {
  flex: 1;
  min-height: 0;

  overflow-y: auto;
  overflow-x: hidden;

  padding-right: 6px;
}

.item-picker-sidebar {
  flex-shrink: 0;

  max-height: 40%;
  overflow-y: auto;
  overflow-x: hidden;

  border-top: 1px solid var(--border);
  margin-top: 10px;
  padding-top: 10px;
}

.submit-button {
  flex-shrink: 0;
  margin-top: 10px;
}
</style>
