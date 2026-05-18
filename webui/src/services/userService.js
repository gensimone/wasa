import api from "@/services/axios"
import { setName, setPhotoUrl } from "@/state/user"

export async function updateName(userId, name) {
    const cleanName = name?.trim()

    if (!cleanName) {
        throw new Error("Invalid name")
    }

    await api.put(
        `/users/${userId}/name`,
        { name: cleanName },
        { headers: { Authorization: userId } }
    )

    setName(cleanName)
    return cleanName
}

export async function updatePhoto(userId, file) {
    if (!file) {
        throw new Error("No file provided")
    }

    const formData = new FormData()
    formData.append("photo", file)

    const { data } = await api.put(
        `/users/${userId}/photo`,
        formData,
        {
            headers: {
                Authorization: userId,
                "Content-Type": "multipart/form-data"
            }
        }
    )

    setPhotoUrl(`${__API_URL__}${data.photoUrl}`)
    return data.photoUrl
}

export async function deletePhoto(userId) {
    await api.delete(`/users/${userId}/photo`, {
        headers: { Authorization: userId }
    })

    const { data } = await api.get(`/users/${userId}`, {
        headers: { Authorization: userId }
    })

    setPhotoUrl(`${__API_URL__}${data.photoUrl}`)
    return data.photoUrl
}
