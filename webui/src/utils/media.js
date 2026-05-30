import { defaultUserPhotoUrl, defaultGroupPhotoUrl } from "../assets/default";

export function expandUrl(url) {
  const isLocal =
    url === defaultUserPhotoUrl ||
    url === defaultGroupPhotoUrl ||
    url?.startsWith("blob:");

  return isLocal ? url : `${__API_URL__}${url}`;
}
