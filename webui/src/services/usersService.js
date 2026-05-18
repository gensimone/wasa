import api from "@/services/axios"
import { user } from "@/state/user"

export async function fetchUser(userId) {
    const response = await api.get(`/users/${userId}`, {
        headers: { Authorization: user.userId }
    })
    return response.data
}

export async function fetchUsers() {
    const response = await api.get(`/users`, {
        headers: { Authorization: user.userId }
    })

    const userPromises = response.data.users
        .filter(id => id !== user.userId)
        .map(id => fetchUser(id))

    return await Promise.all(userPromises)
}
