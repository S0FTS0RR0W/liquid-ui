const API_BASE = "http://127.0.0.1:8765";

export async function apiGet(path: string) {
  const res = await fetch(`${API_BASE}${path}`);
  if (!res.ok) throw new Error(`GET ${path} failed`);
  return res.json();
}

export async function apiPost(path: string, body: any) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });
  if (!res.ok) throw new Error(`POST ${path} failed`);
  return res.json();
}

// Convenience wrappers
export const getDevices = () => apiGet("/devices");
export const getStatus = (deviceIndex: number) =>
  apiGet(`/status?device=${deviceIndex}`);
export const getProfiles = () => apiGet("/profiles");
export const saveProfile = (profile: any) =>
  apiPost("/profiles/save", profile);
export const applyProfile = (deviceIndex: number, profileName: string) =>
  apiPost("/profiles/apply", { deviceIndex, profileName });
