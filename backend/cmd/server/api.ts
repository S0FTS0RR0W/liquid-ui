export interface Device {
  index: number;
  name: string;
  type: string;
}

export interface Status {
  temperature: number;
  fanRpm: number;
  pumpRpm: number;
}

const API_BASE = "http://localhost:8765";

export async function getDevices(): Promise<Device[]> {
  const res = await fetch(`${API_BASE}/devices`);
  if (!res.ok) throw new Error("Failed to fetch devices");
  return res.json();
}

export async function getStatus(index: number): Promise<Status> {
  const res = await fetch(`${API_BASE}/status?device=${index}`);
  if (!res.ok) throw new Error("Failed to fetch status");
  return res.json();
}