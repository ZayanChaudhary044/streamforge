"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";

type Video = {
  id: string;
  title: string;
  filename: string;
};

export default function WatchPage() {
  const { id } = useParams();
  const router = useRouter();
  const [video, setVideo] = useState<Video | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/videos")
      .then((res) => {
        if (!res.ok) throw new Error();
        return res.json();
      })
      .then((videos: Video[]) => {
        const found = videos.find((v) => v.id === id);
        if (!found) throw new Error();
        setVideo(found);
      })
      .catch(() => setError("Failed to load video"));
  }, [id]);

  if (error) {
    return (
      <div className="p-10 text-center text-red-400">
        {error}
      </div>
    );
  }

  if (!video) {
    return (
      <div className="p-10 text-center text-slate-400">
        Loading…
      </div>
    );
  }

  return (
    <main className="max-w-6xl mx-auto px-6 py-10 space-y-6">
      {/* Back button */}
      <button
        onClick={() => router.push("/")}
        className="inline-flex items-center gap-2 text-sm text-slate-400 hover:text-white transition"
      >
        <span className="text-lg">←</span>
        Back to uploads
      </button>

      {/* Video player */}
      <div className="aspect-video w-full overflow-hidden rounded-2xl bg-black shadow-lg">
        <video
          controls
          className="h-full w-full"
          src={`http://localhost:8080/media/${video.filename}`}
        />
      </div>

      {/* Metadata */}
      <div className="space-y-2">
        <h1 className="text-2xl font-semibold leading-snug">
          {video.title}
        </h1>
        <p className="text-sm text-slate-400">
          Streamed via Go backend · Adaptive-ready
        </p>
      </div>
    </main>
  );
}
