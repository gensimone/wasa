import api from "@/services/axios";
import { user } from "@/state/user";

export async function commentMessage(messageId, text, photo) {
  const formData = new FormData();

  formData.append("text", text);

  if (photo) {
    formData.append("file", photo);
    formData.append("mediaType", "photo");
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
