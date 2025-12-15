"use client";

import { useEffect, useState } from "react";
import { fetchVideos, Video } from "./lib/api";
import Link from "next/link";

export default function HomePage() {
  const [videos, setVideos] = useState<Video[]>([]);
  const [title, setTitle] = useState("");
  const [file, setFile] = useState<File | null>(null);
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
    if (!title || !file) return;

    setLoading(true);
    setError(null);

    try {
      const data = new FormData();
      data.append("title", title);
      data.append("file", file);

      const res = await fetch("http://localhost:8080/videos/upload", {
        method: "POST",
        body: data,
      });

      if (!res.ok) throw new Error("Upload failed");

      setTitle("");
      setFile(null);
      await loadVideos();
    } catch {
      setError("Failed to upload video");
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className="space-y-8">
      {/* Upload section */}
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
            disabled={loading}
            className="rounded-lg bg-indigo-500 px-4 py-2 text-sm font-medium text-white"
          >
            {loading ? "Uploading..." : "Upload"}
          </button>
        </form>

        {error && <p className="mt-2 text-sm text-red-400">{error}</p>}
      </section>

      {/* Videos section */}
      <section className="rounded-xl border border-slate-800 p-6">
        <h2 className="mb-4 text-lg font-semibold">Videos</h2>

        {videos.length === 0 ? (
          <p className="text-sm text-slate-400">No videos yet</p>
        ) : (
          <ul className="space-y-3">
            {videos.map((v) => (
              <li key={v.id}>
                {/* 🔴 THIS IS THE IMPORTANT PART */}
                <Link
                  href={`/watch/${v.id}`}
                  className="block rounded-lg border border-slate-800 p-4 space-y-2 hover:border-slate-600 transition cursor-pointer"
                >
                  <img
                    src={`http://localhost:8080/thumbnails/${v.thumbnail}`}
                    alt={v.title}
                    className="w-full h-40 object-cover rounded-lg"
                  />

                  <div>
                    <p className="font-medium">{v.title}</p>
                    <p className="text-xs text-slate-400">
                      {new Date(v.created_at).toLocaleString()}
                    </p>
                  </div>
                </Link>
              </li>
            ))}
          </ul>
        )}
      </section>
    </main>
  );
}
