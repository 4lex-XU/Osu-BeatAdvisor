import React, { useState } from 'react';
import 'react-bootstrap-range-slider/dist/react-bootstrap-range-slider.css';
import osuButton from '../Images/osuButton.png';
import '../CSS/osu_button.css';
import CheckboxModes from './CheckboxModes';
import RangeSliderDiffilculty from './RangeSliderDifficulty';
import CheckboxGenres from './CheckboxGenres';
import CheckboxLanguages from './CheckboxLanguages';
import CheckboxCategories from './CheckboxCategories';
import Modal from './Modal';
import Playlist from './Playlist';
import axios from 'axios';

export default function HomePage(props) {
  const [showOptions, setShowOptions] = useState(false);
  const [genres, setGenres] = useState([]);
  const [languages, setLanguages] = useState([]);
  const [size, setSize] = useState('0');
  const [title, setTitle] = useState('');
  const [minDifficulty, setMinDifficulty] = useState(0);
  const [maxDifficulty, setMaxDifficulty] = useState(10);
  const [status, setStatus] = useState([]);
  const [modes, setModes] = useState([]);
  const [error, setError] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const toggleOptions = () => {
    setShowOptions(!showOptions);
  };

  const handleGenerate = (evt) => {
    evt.preventDefault();
    if (size === '0' || size === '') {
      setError({
        message: 'Erreur',
        detail: 'Le nombre de maps doit être supérieur à 0',
      });
      return;
    }
    if (minDifficulty > maxDifficulty) {
      setError({
        message: 'Erreur',
        detail:
          'La difficulté minimale doit être inférieure à la difficulté maximale',
      });
      return;
    }
    setError(null);
    const data = {
      genres: genres,
      languages: languages,
      size: size,
      title: title,
      difficulty: minDifficulty + '-' + maxDifficulty,
      status: status,
      modes: modes,
    };
    console.log(data);
    axios
      .post('/playlist/create', data, {
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((response) => {
        console.log(response.data);
        setIsModalOpen(true);
      })
      .catch((error) => {
        console.log(error.response.data);
        setError(error.response.data);
      });
  };

  const handleProfile = (evt) => {
    evt.preventDefault();
    props.setCurrentPage(props.myLogin);
  };

  return (
    <div >
      <div
        className={`button-container-osu ${
          showOptions ? 'shift-left' : 'shift-right'
        }`}
      >
        <button
          className="osu-button"
          onClick={toggleOptions}
          style={{ border: 'none', background: 'none', padding: 0 }}
          onMouseOver={(e) => (e.currentTarget.style.transform = 'scale(1.2)')}
          onMouseOut={(e) => (e.currentTarget.style.transform = 'scale(1)')}
        >
          <img src={osuButton} alt="osu!" />
        </button>
        {showOptions && (
          <form className="options-container">
            <div className="row-osu-form">
              <input
                type="text"
                placeholder="Titre de la playlist"
                onChange={(e) => setTitle(e.target.value)}
              />
              <input
                type="number"
                placeholder="Nombre de maps"
                onChange={(e) => setSize(e.target.value)}
              />
            </div>
            <CheckboxModes setIsCheck={setModes} isCheck={modes} />
            <RangeSliderDiffilculty
              min={minDifficulty}
              max={maxDifficulty}
              setMin={setMinDifficulty}
              setMax={setMaxDifficulty}
            />
            <CheckboxCategories setIsCheck={setStatus} isCheck={status} />
            <CheckboxGenres setIsCheck={setGenres} isCheck={genres} />
            <CheckboxLanguages setIsCheck={setLanguages} isCheck={languages} />
            <button type="submit" onClick={handleGenerate}>
              Générer
            </button>
            {error && (
              <p style={{ color: 'red' }}>
                {error.message} {error.detail}
              </p>
            )}
          </form>
        )}
      </div>
      {isModalOpen && (
        <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)}>
          <div>
            <p>
              Votre playlist a été générée avec succès. Voir votre{' '}
              <a href="a" onClick={handleProfile}>
                bibliothèque
              </a>
            </p>
          </div>
        </Modal>
      )}
    </div>
  );
}
