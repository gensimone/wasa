import api from "@/services/axios"
import { user } from "@/state/user"

export async function updateName(rawName) {
    const name = rawName?.trim()

    if (!name) {
        throw new Error("Invalid name")
    }

    await api.put(
        `/users/${user.userId}/name`,
        { name: name },
        { headers: { Authorization: user.userId } }
    )

    return name
}

export async function updatePhoto(file) {
    if (!file) {
        throw new Error("No file provided")
    }

    const formData = new FormData()
    formData.append("photo", file)

    const { data } = await api.put(
        `/users/${user.userId}/photo`,
        formData,
        {
            headers: {
                Authorization: user.userId,
                "Content-Type": "multipart/form-data"
            }
        }
    )

    return data.photoUrl
}

export async function deletePhoto() {
    // FIXME: Do this in one single request.
    await api.delete(`/users/${user.userId}/photo`, {
        headers: { Authorization: user.userId }
    })

    const { data } = await api.get(`/users/${user.userId}`, {
        headers: { Authorization: user.userId }
    })

    return data.photoUrl
}
