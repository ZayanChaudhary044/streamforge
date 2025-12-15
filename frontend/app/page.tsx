"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { fetchVideos, Video } from "./lib/api";

export default function HomePage() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [title, setTitle] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);

  async function loadVideos() {
    try {
      setLoading(true);
      const data = await fetchVideos();

      // ✅ GUARANTEE array
      setVideos(Array.isArray(data) ? data : []);
      setError(null);
    } catch {
      setVideos([]);
      setError("Failed to load videos");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadVideos();
  }, []);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!title || !file) return;

    setUploading(true);
    setError(null);

    try {
      const data = new FormData();
      data.append("title", title);
      data.append("file", file);

      const res = await fetch("http://localhost:8080/videos/upload", {
        method: "POST",
        body: data,
      });

      if (!res.ok) throw new Error();

      setTitle("");
      setFile(null);
      await loadVideos();
    } catch {
      setError("Failed to upload video");
    } finally {
      setUploading(false);
    }
  }

  return (
    <main className="space-y-8">
      {/* Upload */}
      <section className="rounded-xl border border-slate-800 p-6">
        <h2 className="mb-4 text-lg font-semibold">Upload Video</h2>

        <form onSubmit={onSubmit} className="space-y-3">
          <input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Video title"
            className="w-full rounded-lg bg-slate-900 px-3 py-2 text-sm"
          />

          <input
            type="file"
            accept="video/*"
            onChange={(e) => setFile(e.target.files?.[0] || null)}
            className="text-sm"
          />

          <button
            disabled={uploading}
            className="rounded-lg bg-indigo-500 px-4 py-2 text-sm font-medium text-white disabled:opacity-50"
          >
            {uploading ? "Uploading..." : "Upload"}
          </button>
        </form>

        {error && <p className="mt-2 text-sm text-red-400">{error}</p>}
      </section>

      {/* Videos */}
      <section className="rounded-xl border border-slate-800 p-6">
        <h2 className="mb-4 text-lg font-semibold">Videos</h2>

        {loading ? (
          <p className="text-sm text-slate-400">Loading videos…</p>
        ) : videos.length === 0 ? (
          <p className="text-sm text-slate-400">No videos uploaded yet</p>
        ) : (
          <ul className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {videos.map((v) => (
              <li
                key={v.id}
                className="rounded-xl border border-slate-800 overflow-hidden"
              >
                {v.thumbnail ? (
                  <Link href={`/watch/${v.id}`}>
                    <img
                      src={`http://localhost:8080/thumbnails/${v.thumbnail}`}
                      alt={v.title}
                      className="h-40 w-full object-cover cursor-pointer hover:opacity-90 transition"
                    />
                  </Link>
                ) : (
                  <div className="h-40 w-full bg-slate-900 animate-pulse flex items-center justify-center text-xs text-slate-500">
                    Processing…
                  </div>
                )}

                <div className="p-3 space-y-1">
                  <p className="font-medium truncate">{v.title}</p>
                  <p className="text-xs text-slate-400">
                    {new Date(v.created_at).toLocaleString()}
                  </p>

                  {v.status !== "ready" && (
                    <span className="inline-block rounded bg-yellow-500/10 px-2 py-0.5 text-xs text-yellow-400">
                      {v.status}
                    </span>
                  )}
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>
    </main>
  );
}
