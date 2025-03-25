import { useEffect, useRef } from "react";

function BackgroundAudio({ sources, isPlaying }) {
  const audioRef = useRef(null);
  const currentIndexRef = useRef(0);

  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const handleEnded = () => {
      currentIndexRef.current =
          (currentIndexRef.current + 1) % sources.length;
      audio.src = sources[currentIndexRef.current];
      audio.play().catch(console.error);
    };

    audio.addEventListener("ended", handleEnded);
    return () => {
      audio.removeEventListener("ended", handleEnded);
    };
  }, [sources]);

  useEffect(() => {
    const audio = audioRef.current;
    if (isPlaying && audio) {
      audio.src = sources[currentIndexRef.current];
      audio.play().catch(console.error);
    }
  }, [isPlaying, sources]);

  return (
      <audio ref={audioRef} loop={sources.length === 1} />
  );
}

export default BackgroundAudio;
