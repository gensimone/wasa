import { ref } from "vue";

export function useSettingsForm(initialText = "", defaultText = "...") {
  const text = ref(initialText);
  const placeholder = ref(defaultText);
  const loading = ref(false);

  function setText(value) {
    text.value = value === "" ? placeholder.value : value;
  }

  function validate() {
    const value = text.value?.trim();
    return value || null;
  }

  async function submit(handler) {
    const value = validate();

    if (!value) {
      throw new Error("EMPTY_NAME");
    }

    loading.value = true;

    try {
      await handler(value);
    } finally {
      loading.value = false;
    }
  }

  return {
    text,
    loading,
    placeholder,

    setText,
    submit,
  };
}
