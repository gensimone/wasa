import { getUserById } from "@/services/users"
import { handleError } from "@/utils/errors"
import { userConversations, groupConversations } from "@/state/conversations"
import { getMyConversations } from "@/services/conversations"

export async function isValidUserId(id) {
    let isValid = userConversations.value.has(id)

    if (!isValid) { // We must make sure the user exists.
        try {
            const user = await getUserById(id)
            isValid = user?.userId === id

        } catch (e) {
            if (e.response?.status !== 404) {
                handleError(e)
            }

            isValid = false
        }
    }

    return isValid
}

export async function isValidConversationId(id) {
    let isValid = groupConversations.value.has(id)

    if (groupConversations.value.size !== 0) return isValid

    // groupConversations may not have been loaded yet.
    try {
        const fetched = (await getMyConversations()) || []
        isValid = fetched
            .filter(c => c.isGroup)
            .map(c => c.id)
            .includes(id)

    } catch (e) {
        handleError(e)
        isValid = false

    } finally {
        return isValid
    }
}

export async function isValidGroupRoute(route) {
    const id = Number(route.params.id)
    const validId = Number.isInteger(id) && id > 0

    return validId ? isValidConversationId(id) : false
}

export async function isValidConversationRoute(route) {
    const id = Number(route.params.id)
    const validId = Number.isInteger(id) && id > 0

    const direct = String(route.query.direct).toLowerCase()
    const validDirect = direct === "true" || direct === "false"

    if (!validId || !validDirect) return false

    return direct === "true" ? isValidUserId(id) : isValidConversationId(id)
}
