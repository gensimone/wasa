import api from "@/services/axios";
import { user } from "@/state/user";

export async function commentMessage(
  messageId,
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

  const response = await api.post(`/comments/${messageId}`, formData, {
    headers: {
      Authorization: user.userId,
      "Content-Type": "multipart/form-data",
    },
  });

  return response.data;
}

export async function uncommentMessage(messageId) {
  return await api.delete(`/comments/${messageId}`, {
    headers: { Authorization: user.userId },
  });
}
