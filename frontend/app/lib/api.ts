const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:8080";

export interface Video {
  id: string;
  title: string;
  filename: string;
  thumbnail: string;
  created_at: string;
}

export async function fetchVideos(): Promise<Video[]> {
  const res = await fetch(`${API_BASE}/videos`);
  if (!res.ok) throw new Error("Failed to fetch videos");
  return res.json();
}

export async function createVideo(title: string): Promise<Video> {
  const res = await fetch(`${API_BASE}/videos`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title }),
  });

  if (!res.ok) throw new Error("Failed to create video");
  return res.json();
}
