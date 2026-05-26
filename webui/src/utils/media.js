import { defaultUserPhotoUrl, defaultGroupPhotoUrl } from "../assets/default"

export function expandUrl(url) {
    const isDefault =
        url === defaultUserPhotoUrl ||
        url === defaultGroupPhotoUrl

    const isLocal = url?.startsWith("blob:")

    return isLocal || isDefault ? url : `${__API_URL__}${url}`
}
