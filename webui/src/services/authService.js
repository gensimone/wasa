import api from "@/services/axios"
import { setUserId, setName, setPhotoUrl } from "@/state/user"

export async function login(rawName) {
    const name = rawName?.trim()
    if (!name) {
        throw new Error("Invalid name")
    }

    const response = await api.post("/session", {
        name: name
    })

    const data = response.data

    setUserId(data.userId)
    setName(data.name)
    setPhotoUrl(data.photoUrl)
}
