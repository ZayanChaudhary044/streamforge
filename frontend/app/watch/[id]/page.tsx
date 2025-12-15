"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

type Video = {
  id: string;
  title: string;
  filename: string;
};

export default function WatchPage() {
  const { id } = useParams();
  const [video, setVideo] = useState<Video | null>(null);
  const [error, setError] = useState("");

  useEffect(() => {
    fetch("http://localhost:8080/videos")
      .then((res) => {
        if (!res.ok) throw new Error("Failed");
        return res.json();
      })
      .then((videos: Video[]) => {
        const found = videos.find((v) => v.id === id);
        if (!found) throw new Error("Not found");
        setVideo(found);
      })
      .catch(() => setError("Failed to load video"));
  }, [id]);

  if (error) return <p className="text-red-500 p-4">{error}</p>;
  if (!video) return <p className="p-4">Loading...</p>;

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-2xl font-semibold mb-4">
        {video.title}
      </h1>

      <video
        controls
        className="w-full rounded bg-black"
        src={`http://localhost:8080/media/${video.filename}`}
      />
    </div>
  );
}
