import api from "@/services/axios";
import { user } from "@/state/user";

export async function getMyConversations() {
  const response = await api.get(`/conversations`, {
    headers: { Authorization: user.userId },
  });

  return response.data.conversations;
}

export async function getConversation(id, direct) {
  const response = await api.get(`/conversations/${id}`, {
    headers: { Authorization: user.userId },
    params: { direct },
  });

  return response.data.messageIds;
}

export async function sendMessage(
  id,
  direct,
  text,
  attachment,
  mediaType = "image",
) {
  const formData = new FormData();

  if (text) formData.append("text", text);

  if (attachment) {
    formData.append("file", attachment);
    formData.append("mediaType", mediaType);
  }

  const response = await api.post(`/conversations/${id}/message`, formData, {
    headers: {
      Authorization: user.userId,
      "Content-Type": "multipart/form-data",
    },
    params: { direct },
  });

  return response.data;
}

export async function forwardMessage(id, direct, messageId) {
  const response = await api.post(
    `/conversations/${id}/fmessage`,
    { messageId },
    {
      headers: { Authorization: user.userId },
      params: { direct },
    },
  );

  return response.data;
}

// Not needed?
// export async function getLastMessage(conversationId) {
//     const response = await api.get(`/conversations/${conversationId}/last`,
//         { headers: { Authorization: user.userId } }
//     )
//
//     return response.data
// }
