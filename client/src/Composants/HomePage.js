import React, {useState} from "react";
import "react-bootstrap-range-slider/dist/react-bootstrap-range-slider.css";
import osuButton from "../Images/osuButton.png";
import "../CSS/osu_button.css";
import CheckboxModes from "./CheckboxModes";
import RangeSliderDiffilculty from "./RangeSliderDifficulty";
import CheckboxGenres from "./CheckboxGenres";
import CheckboxLanguages from "./CheckboxLanguages";
import CheckboxCategories from "./CheckboxCategories";
import Modal from "./Modal";
import axios from "axios";

export default function HomePage(props) {
  const [showOptions, setShowOptions] = useState(false);
  const [genres, setGenres] = useState([]);
  const [languages, setLanguages] = useState([]);
  const [size, setSize] = useState("0");
  const [title, setTitle] = useState("");
  const [minDifficulty, setMinDifficulty] = useState(0);
  const [maxDifficulty, setMaxDifficulty] = useState(10);
  const [status, setStatus] = useState([]);
  const [modes, setModes] = useState([]);
  const [error, setError] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isToggleActive, setIsToggleActive] = useState(false);
  const [description, setDescription] = useState("");
  const handleToggle = (event) => {
    event.preventDefault(); // Empêche le comportement par défaut du bouton
    setIsToggleActive(!isToggleActive);
  };

  const toggleOptions = () => {
    setShowOptions(!showOptions);
  };

  const isGenerateDisabled = () => {
    return (
      title.trim() === "" || // Vérifie si le titre est vide
      size === "0" ||
      size === "" || // Vérifie si la taille est vide ou 0
      modes.length === 0 || // Vérifie s'il y a au moins un mode
      status.length === 0 || // Vérifie s'il y a au moins un status (ex: ranked)
      genres.length === 0 || // Vérifie s'il y a au moins un genre
      languages.length === 0 // Vérifie s'il y a au moins une langue
    );
  };

  const handleGenerate = (evt) => {
    evt.preventDefault();
    if (size === "0" || size === "") {
      setError({
        message: "Erreur",
        detail: "Le nombre de maps doit être supérieur à 0",
      });
      return;
    }
    if (minDifficulty > maxDifficulty) {
      setError({
        message: "Erreur",
        detail:
          "La difficulté minimale doit être inférieure à la difficulté maximale",
      });
      return;
    }
    setError(null);
    const data = {
      genres: genres,
      languages: languages,
      size: size,
      title: title,
      difficulty: minDifficulty + "-" + maxDifficulty,
      status: status,
      modes: modes,
      ...(isToggleActive ? {} : { description }), // Ajoute 'description' uniquement si le toggle est désactivé
    };
    console.log(data);
    axios
      .post("/playlist/create", data, {
        headers: {
          "Content-Type": "application/json",
        },
        withCredentials: true,
        credentials: "include",
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

  const handleGenerateFromMyPlaylists = (evt) => {
    evt.preventDefault();
    axios
      .post(
        "/playlist/generate-from-my-playlists",
        {},
        {
          headers: {
            "Content-Type": "application/json",
          },
          withCredentials: true,
          credentials: "include",
        },
      )
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
    <div>
      <div
        className={`button-container-osu ${
          showOptions ? "shift-left" : "shift-right"
        }`}
      >
        <button
          className="osu-button"
          onClick={toggleOptions}
          style={{ border: "none", background: "none", padding: 0 }}
          onMouseOver={(e) => (e.currentTarget.style.transform = "scale(1.2)")}
          onMouseOut={(e) => (e.currentTarget.style.transform = "scale(1)")}
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
            <div style={{ display: 'flex',gap: '10px' }}>
              <input
                type="text"
                placeholder="Description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                style={{
                  flex: 1,
                  padding: '8px',
                  borderRadius: '4px',
                  border: '1px solid #ccc',
                  backgroundColor: isToggleActive ? '#f0f0f0' : '#fff', // Grisé si toggle activé
                  color: isToggleActive ? '#a0a0a0' : '#000', // Texte grisé si toggle activé
                  cursor: isToggleActive ? 'not-allowed' : 'text', // Curseur interdit si toggle activé
                }}
                disabled={isToggleActive}
              />
              <button
                onClick={handleToggle}
                style={{
                  padding: '8px 16px',
                  borderRadius: '20px',
                  border: 'none',
                  backgroundColor: isToggleActive ? '#007bff' : '#007bff',
                  color: '#fff',
                  cursor: 'pointer',
                  transition: 'background-color 0.3s ease',
                }}
              >
                {isToggleActive ? 'Désactiver' : 'Activer'} la génération automatique
              </button>
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
            <button type="submit" onClick={handleGenerate}
            >
              Générer
            </button>
            {error && (
              <p style={{ color: "red" }}>
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
              Votre playlist a été générée avec succès. Voir votre{" "}
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
