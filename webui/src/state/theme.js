import { ref, watch } from "vue";

export const isDark = ref(localStorage.getItem("theme") === "dark");

watch(
  isDark,
  (val) => {
    document.documentElement.classList.toggle("dark", val);
    localStorage.setItem("theme", val ? "dark" : "light");
    applyTheme(isDark.value);
  },
  { immediate: true },
);

export function applyTheme(value) {
  document.documentElement.setAttribute("data-theme", value ? "dark" : "light");
}

export function getIcon(name) {
  const dir = isDark.value ? "light" : "dark"; // Inverted
  const ext = name?.endsWith(".svg") ? "" : ".svg";

  return `/icons/${dir}/${name}${ext}`;
}
