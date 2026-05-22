import api from "@/services/axios"

export async function doLogin(name) {
    const response = await api.post(`/session`,
        { name: name }
    )

    return response.data
}
