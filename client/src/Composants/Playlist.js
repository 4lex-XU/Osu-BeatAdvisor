import React, { useState, useEffect } from 'react';
import axios from 'axios';
import {
  faTrash,
  faHeart as faHeartSolid,
  faFlag,
} from '@fortawesome/free-solid-svg-icons';
import {
  faHeart as faHeartRegularIcon,
  faComment,
} from '@fortawesome/free-regular-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import SaisieMessage from './SaisieMessage';
import ListeMessages from './ListeMessages';

export default function Playlist(props) {
  const [playlist, setPlaylist] = useState(props.playlist);
  const [mettreAJour, setMettreAJour] = useState(false);
  const [afficher, setAfficher] = useState(false);
  const [erreur, setErreur] = useState(null);
  const [isLiked, setIsLiked] = useState(false);
  const [likes, setLikes] = useState([]);
  const [showLikes, setShowLikes] = useState(false);
  const [hover, setHover] = useState(false);
  const [afficherCommentaires, setAfficherCommentaires] = useState(false);
  const [commentaires, setCommentaires] = useState([]);
  const [saisir, setSaisir] = useState(false);
  const [editPlaylist, setEditPlaylist] = useState(false);
  const [title, setTitle] = useState(props.title);
  const [genres, setGenres] = useState([]);
  const [languages, setLanguages] = useState([]);
  const [description, setDescription] = useState('');

  const toggleDetails = () => {
    setAfficher(!afficher);

    // if (playlist) {
    //   setAfficher(!afficher);
    // } else {
    //   axios
    //   .get(`/playlist/${props.playlistId}`, {
    //     headers: {
    //       'Content-Type': 'application/json',
    //     },
    //     withCredentials: true,
    //     credentials: 'include',
    //   })
    //   .then((res) => {
    //     console.log(res.data);
    //     setPlaylist(res.data);
    //     if (res.data.Likes) {
    //       setLikes(res.data.Likes);
    //       if (res.data.Likes.includes(props.myLogin)) {
    //         setIsLiked(true);
    //       }
    //     }
    //     if (res.data.Comments) {
    //       setCommentaires(res.data.Comments);
    //     }
    //     setAfficher(true);
    //   })
    //   .catch((err) => {
    //     console.log(err.response.data);
    //     setErreur(err.response.data);
    //   });
    // }
  };

  const fetchPlaylist = () => {
    axios
        .get(`/playlist/${props.playlistId}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          withCredentials: true,
          credentials: 'include',
        })
        .then((res) => {
          // console.log("playlist :", res.data); // 🔍 Vérifier si l'API renvoie bien les données

          setPlaylist(res.data); // ✅ Mettre à jour la playlist avant d'utiliser ses valeurs

          setLikes(res.data.Likes || []);
          setIsLiked(res.data.Likes?.includes(props.myLogin) || false);
          setCommentaires(res.data.Comments || []);

          setLanguages(res.data.Languages || []);
          setGenres(res.data.Genres || []);
          setTitle(res.data.Title || "Sans Titre");
          setDescription(res.data.Description || "Pas de description disponible.");
        })
        .catch((err) => {
          console.log(err.response?.data || "Erreur lors du chargement de la playlist");
          setErreur(err.response?.data);
        });
  };

  const handlerSaisie = (evt) => {
    if (saisir) setSaisir(false);
    else setSaisir(true);
  };

  const handleRemove = (beatmapId) => {
    console.log(props.playlistId);
    console.log(beatmapId);
    axios
      .delete(`/playlist/delete/beatmap/${props.playlistId}/${beatmapId}`, {
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
      })
      .then((response) => {
        console.log(response.data);
        if (playlist && playlist.Beatmaps) {
          const updatedBeatmaps = playlist.Beatmaps.filter(
            (map) => map.Id !== beatmapId
          );
          setPlaylist({ ...playlist, Beatmaps: updatedBeatmaps });
        }
        setMettreAJour(!mettreAJour);
      })
      .catch((error) => {
        console.log(error.response.data);
        setErreur(error.response.data);
      });
  };

  // PERMET DE LIKE UN MESSAGE
  const like = (evt) => {
    evt.preventDefault();
    axios
      .post(
        `/playlist/like/${props.playlistId}`,
        {},
        {
          headers: {
            'Content-Type': 'application/json',
          },
          withCredentials: true,
          credentials: 'include',
        }
      )
      .then((res) => {
        console.log(res.data);
        setIsLiked(true);
        setLikes((likes) => [...likes, props.myLogin]);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  };

  // PERMET DE DISLIKE UN MESSAGE
  const dislike = (evt) => {
    evt.preventDefault();
    axios
      .post(
        `/playlist/like/delete/${props.playlistId}`,
        {},
        {
          headers: {
            'Content-Type': 'application/json',
          },
          withCredentials: true,
          credentials: 'include',
        }
      )
      .then((res) => {
        console.log(res.data);
        setIsLiked(false);
        setLikes((likes) => likes.filter((liker) => liker !== props.myLogin));
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  };

  const toggleLikersList = () => {
    setShowLikes(!showLikes);
  };

  // PERMET D'AFFICHER LES COMMENTAIRES
  const handlerCommentaire = (evt) => {
    evt.preventDefault();
    setAfficherCommentaires(!afficherCommentaires);
  };

  // PERMET DE CHANGER DE PAGE
  const pageProfilHandler = (evt, name) => {
    evt.preventDefault();
    props.setCurrentPage(name);
  };

  const handleSetEdit = () => {
    setEditPlaylist(editPlaylist? false : true);
    fetchPlaylist();
  };

  const handleEdit = (playlistId) => {
    const data = {
      title: title,
      genres: genres,
      languages: languages,
    };
    axios
      .put(`/playlist/edit/${playlistId}`, data, {
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((res) => {
        console.log(res.data);
        setEditPlaylist(false);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  }

  const handleAnnuler = () => {
    setEditPlaylist(false);
    fetchPlaylist();
  }

  return (
    <div style={{width: '100%'}}>
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center'}}>
        <h3
          onClick={toggleDetails}
          style={{ cursor: 'pointer', textDecoration: 'underline' }}
        >
          {title || 'Sans Titre'}
        </h3>
        {props.myPage && (
          <>
            <button
            onClick={() => props.handleDelete(props.playlistId)}
            style={{ marginLeft: '10px' }}
            >
              Supprimer
            </button>
            <button
            onClick={() => handleSetEdit()}
            style={{ marginLeft: '10px' }}
            >
              Editer
            </button>
          </>
        )}
        {editPlaylist && (
          <div>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
            <input
              type="text"
              value={genres.join(', ')}
              onChange={(e) => setGenres(e.target.value.split(', '))}
            />
            <input
              type="text"
              value={languages.join(', ')}
              onChange={(e) => setLanguages(e.target.value.split(', '))}
            />
            <button onClick={() => handleEdit(props.playlistId)}>Valider</button>
            <button onClick={() => handleAnnuler()}>Annuler</button>
          </div>
          )
          }
      </div>
      {afficher && (
        <>
          <br />
          <div>
            {!isLiked ? (
              <button onClick={like}>
                <FontAwesomeIcon icon={faHeartRegularIcon} />
              </button>
            ) : (
              <button onClick={dislike}>
                <FontAwesomeIcon icon={faHeartSolid} />
              </button>
            )}
            <span
              onClick={toggleLikersList}
              onMouseEnter={() => setHover(true)}
              onMouseLeave={() => setHover(false)}
              style={{
                cursor: 'pointer',
                textDecoration: hover ? 'underline' : 'none',
              }}
            >
              {' '}
              {likes.length} J'aime {' '}
            </span>
            <button onClick={handlerSaisie}>
              <FontAwesomeIcon icon={faComment} />
            </button>{' '}
            {commentaires.length}{' '}
            {!afficherCommentaires ? (
              <button onClick={handlerCommentaire}>
                Afficher <FontAwesomeIcon icon={faComment} />
              </button>
            ) : (
              <button onClick={handlerCommentaire}>
                Masquer <FontAwesomeIcon icon={faComment} />
              </button>
            )}
            <br />
            {showLikes && (
              <ul>
                {likes.map((liker, index) => (
                  <li key={index}>
                    <a href="a" onClick={(e) => pageProfilHandler(e, liker)}>
                      {liker}
                    </a>
                  </li>
                ))}
              </ul>
            )}
            {saisir && (
              <SaisieMessage
                myLogin={props.myLogin}
                playlistId={props.playlistId}
                setCommentaires={setCommentaires}
                commentaires={commentaires}
              />
            )}
            {afficherCommentaires && (
              <ListeMessages
                commentaires={commentaires}
                myLogin={props.myLogin}
                setCurrentPage={props.setCurrentPage}
                setCommentaires={setCommentaires}
              />
            )}
          </div>
          <br />

          <p>
            Créé par{' '}
            <a href="a" onClick={(e) => pageProfilHandler(e, playlist.Author)}>
              {playlist.Author}
            </a>
          </p>
          <p>Description : {playlist.Description}</p>
          <p>Modes : {playlist.Modes.join(', ')}</p>
          <p>Difficulté : {playlist.Difficulty}</p>
          <p>Langues : {playlist.Languages.join(', ')}</p>
          <p>Genres : {playlist.Genres.join(', ')}</p>
          <p>Statuts : {Array.from(new Set(playlist.Beatmaps.map(map => map.Status))).join(', ')}</p>
          <p>Nombres de Beatmaps : {playlist.Beatmaps.length}</p>
          <ul>
            {playlist.Beatmaps != null &&
              playlist.Beatmaps.map((map) => (
                <li key={map.Id}>
                  <div
                    style={{
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                    }}
                  >
                    <a href={map.Url} target="_blank" rel="noopener noreferrer">
                      <h3>{map.Beatmapset.Title}</h3>
                    </a>
                  </div>
                  {erreur && <p style={{ color: 'red' }}>{erreur}</p>}
                  <p>Mode : {map.Mode}</p>
                  <p>Difficulté : {map.Difficulty_rating}</p>
                  <p>Langue : {map.Beatmapset.Language}</p>
                  <p>Genre : {map.Beatmapset.Genre}</p>
                  <p>Status : {map.Status}</p>
                  <button onClick={(e) => handleRemove(map.Id)}>
                      Retirer
                    </button>
                </li>
              ))}
          </ul>
        </>
      )}
    </div>
  );
}
