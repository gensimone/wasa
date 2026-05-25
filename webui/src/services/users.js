import api from "@/services/axios"
import { user } from "@/state/user"

export async function getUserById(userId) {
    const response = await api.get(`/users/${userId}`,
        { headers: { Authorization: user.userId } }
    )

    return response.data
}

export async function deleteUser() {
    return await api.delete(`/users/${user.userId}`,
        { headers: { Authorization: user.userId } }
    )
}

export async function setMyUserName(name) {
    const response = await api.put(`/users/${user.userId}/name`,
        { name: name },
        { headers: { Authorization: user.userId } }
    )

    return response.data.name
}

export async function setMyPhoto(photo) {
    const formData = new FormData()

    formData.append("photo", photo)

    const response = await api.put(`/users/${user.userId}/photo`,
        formData,
        {
            headers: {
                Authorization: user.userId,
                "Content-Type": "multipart/form-data"
            }
        }
    )

    return response.data.photoUrl
}

export async function deleteMyPhoto() {
    return await api.delete(`/users/${user.userId}/photo`,
        { headers: { Authorization: user.userId } }
    )
}

export async function getUsers() {
    const response = await api.get(`/users`,
        { headers: { Authorization: user.userId } }
    )

    return response.data.users
}
