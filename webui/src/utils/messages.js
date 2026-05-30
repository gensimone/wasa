// export function getTime(sentAt) {
//   const date = new Date(sentAt);
//
//   const hh = String(date.getHours()).padStart(2, "0");
//   const mm = String(date.getMinutes()).padStart(2, "0");
//
//   return `${hh}:${mm}`;
// }

export function getTime(sentAt) {
  const date = new Date(sentAt);
  const now = new Date();

  const isToday =
    date.getDate() === now.getDate() &&
    date.getMonth() === now.getMonth() &&
    date.getFullYear() === now.getFullYear();

  if (isToday) {
    const hh = String(date.getHours()).padStart(2, "0");
    const mm = String(date.getMinutes()).padStart(2, "0");
    return `${hh}:${mm}`;
  }

  const dd = String(date.getDate()).padStart(2, "0");
  const mm = String(date.getMonth() + 1).padStart(2, "0");

  return `${dd}/${mm}`;
}

export function getCheckIcon(receipts) {
  if (!receipts) return null;

  const statuses = receipts.map((r) => r.status);

  const hasReceived = statuses.includes("received");
  const allRead = statuses.length > 0 && statuses.every((s) => s === "read");

  let icon = "check-sent";
  if (hasReceived) icon = "check-received";
  else if (allRead) icon = "check-read";

  return icon;
}
