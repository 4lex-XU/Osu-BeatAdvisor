import React, { useEffect, useRef, useState } from 'react';

function BackgroundAudio({ sources }) {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [playStarted, setPlayStarted] = useState(false); // État pour gérer si la lecture a commencé
  const audioRef = useRef(null);

  useEffect(() => {
    if (playStarted && audioRef.current) {
      audioRef.current.play().catch(error => console.error("Error playing audio:", error));
      const handleEnded = () => {
        setCurrentIndex((prevIndex) => (prevIndex + 1) % sources.length);
      };

      audioRef.current.addEventListener('ended', handleEnded);

      return () => audioRef.current.removeEventListener('ended', handleEnded);
    }
  }, [currentIndex, playStarted, sources]);

  const startPlayback = () => {
    setPlayStarted(true);
  };

  return (
    <div style={{ position: 'absolute', left: '-9999px' }}>
      <audio ref={audioRef} src={sources[currentIndex]} loop={sources.length === 1} />
      {!playStarted && (
        <div className="modal-music">
          <div className="modal-content-music">
            <h4>Bienvenue sur BeatAdvisor!</h4>
            <p>Ici, vous pourrez générer de nombreuses playlists pour vos sessions sur Osu!</p>
            <button onClick={startPlayback}>Autoriser la musique de fond </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default BackgroundAudio;
