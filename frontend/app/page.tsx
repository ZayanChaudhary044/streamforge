"use client";

import { useEffect, useState } from "react";
import { fetchVideos, createVideo, Video } from "./lib/api";

export default function HomePage() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [title, setTitle] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function loadVideos() {
    try {
      setVideos(await fetchVideos());
    } catch {
      setError("Failed to load videos");
    }
  }

  useEffect(() => {
    loadVideos();
  }, []);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!title.trim()) return;

    setLoading(true);
    try {
      await createVideo(title);
      setTitle("");
      await loadVideos();
    } catch {
      setError("Failed to create video");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="space-y-8">
      <section className="rounded-xl border border-slate-800 p-6">
        <h2 className="mb-4 text-lg font-semibold">Create Video</h2>
        <form onSubmit={onSubmit} className="flex gap-3">
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Video title"
            className="flex-1 rounded-lg bg-slate-900 px-3 py-2 text-sm"
          />
          <button
            disabled={loading}
            className="rounded-lg bg-indigo-500 px-4 py-2 text-sm font-medium text-white"
          >
            {loading ? "Adding..." : "Add"}
          </button>
        </form>
        {error && <p className="mt-2 text-sm text-red-400">{error}</p>}
      </section>

      <section className="rounded-xl border border-slate-800 p-6">
        <h2 className="mb-4 text-lg font-semibold">Videos</h2>
        {videos.length === 0 ? (
          <p className="text-sm text-slate-400">No videos yet</p>
        ) : (
          <ul className="space-y-3">
            {videos.map((v) => (
              <li
                key={v.id}
                className="rounded-lg border border-slate-800 px-4 py-3"
              >
                <p className="font-medium">{v.title}</p>
                <p className="text-xs text-slate-400">
                  {new Date(v.created_at).toLocaleString()}
                </p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </main>
  );
}
