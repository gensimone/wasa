import api from "@/services/axios";
import { user } from "@/state/user";

export async function createGroup(name, photo) {
  const formData = new FormData();

  if (name) formData.append("name", name);

  if (photo) {
    formData.append("photo", photo);
  }

  const response = await api.post(`/groups`, formData, {
    headers: {
      Authorization: user.userId,
      "Content-Type": "multipart/form-data",
    },
  });

  return response.data;
}

export async function deleteGroup(groupId) {
  return await api.delete(`/groups/${groupId}`, {
    headers: { Authorization: user.userId },
  });
}

export async function getGroup(groupId) {
  const response = await api.get(`/groups/${groupId}`, {
    headers: { Authorization: user.userId },
  });

  return response.data;
}

export async function addToGroup(groupId, userId) {
  const response = await api.post(
    `/groups/${groupId}`,
    { userId: userId },
    { headers: { Authorization: user.userId } },
  );

  return response.data;
}

export async function setGroupName(groupId, rawName) {
  const name = rawName?.trim();
  if (!name) {
    throw new Error("Invalid name");
  }

  const response = await api.put(
    `/groups/${groupId}/name`,
    { name: name },
    { headers: { Authorization: user.userId } },
  );

  return response.data.name;
}

export async function setGroupPhoto(groupId, photo) {
  const formData = new FormData();

  formData.append("photo", photo);

  const response = await api.put(`/groups/${groupId}/photo`, formData, {
    headers: {
      Authorization: user.userId,
      "Content-Type": "multipart/form-data",
    },
  });

  return response.data.photoUrl;
}

export async function deleteGroupPhoto(groupId) {
  return await api.delete(`/groups/${groupId}/photo`, {
    headers: { Authorization: user.userId },
  });
}

export async function leaveGroup(groupId) {
  return await api.delete(`groups/${groupId}/user`, {
    headers: { Authorization: user.userId },
  });
}

export async function removeUser(groupId, userId) {
  return await api.delete(`groups/${groupId}/user/${userId}`, {
    headers: { Authorization: user.userId },
  });
}

export async function getMemberIds(groupId) {
  const response = await api.get(`/groups/${groupId}/members`, {
    headers: { Authorization: user.userId },
  });

  return response.data.userIds;
}
