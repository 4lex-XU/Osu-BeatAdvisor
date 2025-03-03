import { useState, useEffect } from 'react';
import ListeProfils from './ListeProfils';
import avatar from '../Images/avatar.png';
import entete from '../Images/entete.png';
import '../CSS/profil.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {
  faCake,
  faMapMarkerAlt,
} from '@fortawesome/free-solid-svg-icons';

import axios from 'axios';
import ListePlaylists from './ListePlaylists';

export default function PageProfil(props) {
  const [isAbonne, setIsAbonne] = useState(false);
  const [afficherAmis, setAfficherAmis] = useState(false);
  const [amis, setAmis] = useState([]);
  const [email, setEmail] = useState('');
  const [naissance, setNaissance] = useState('');
  const [description, setDescription] = useState('');
  const [ville, setVille] = useState('');
  const [playlists, setPlaylists] = useState([]);

  useEffect(() => {
    fetchPlaylists();
  }, [props.userProfil]);

  const fetchPlaylists = () => {
    axios
      .get(`playlist/get/${props.userProfil}`, {
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((res) => {
        console.log("playlists: ",res.data);
        setPlaylists(res.data.playlists);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  };

  // Lors de l'ouverture de la page, on récupère les informations du profil
  useEffect(() => {
    axios
      .get(`/user/${props.userProfil}`, {
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((res) => {
        console.log(res.data);
        setNaissance(res.data.Naissance);
        setDescription(res.data.Description);
        setVille(res.data.Ville);
        setEmail(res.data.Email);
        setAmis(res.data.Friends);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  }, [props.userProfil]);

  // Au chargement de la page, on détermine si l'utilisateur est abonné au profil
  useEffect(() => {
    axios
      .get(`/user/friends/${props.userProfil}`, {
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((res) => {
        setIsAbonne(res.data.isFriend);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  }, [isAbonne, props.userProfil]);

  // Permet de suivre un profil
  const Follow = (evt) => {
    evt.preventDefault();
    const data = new URLSearchParams();
      data.append('friend', props.userProfil);
    
    axios
      .put(`/user/friends`, data, {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((res) => {
        console.log(res.data);
        setIsAbonne(true);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  };

  // Permet de ne plus suivre un profil
  const unFollow = (evt) => {
    evt.preventDefault();
    axios
      .delete(`/user/friends/delete/${props.userProfil}`, {
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((res) => {
        console.log(res.data);
        setIsAbonne(false);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  };

  // Permet de passer à la page d'édition du profil
  const handleEdit = (evt) => {
    evt.preventDefault();
    props.setCurrentPage('edit_page');
  };

  // Permet de supprimer le compte
  const handleDelete = () => {
    axios
      .delete(`/user/delete/${props.myLogin}`, {
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((res) => {
        props.setLogout();
        console.log(res.data);
      })
      .catch((err) => {
        props.setLogout();
        console.log(err.response.data);
      });
  };

  // Permet de récupérer la liste des amis
  const getListAmis = () => {
    setAfficherAmis(!afficherAmis);
  };

  return (
    <div className="profil">
      <div className="headerProfil" style={{ minHeight: '81.2vh' }}>
        <div className="entete">
          <img src={entete} />
        </div>
        <div className="avatar">
          <img src={avatar} />
        </div>
        <div className="header-info">
          <div className="prenom-nom">{props.userProfil}</div>
          <div className="tag">{email}</div>
        </div>
        <div className="list-info">
          <p>{description}</p>
          <FontAwesomeIcon icon={faCake} /> {naissance}
          {'  '}
          <FontAwesomeIcon icon={faMapMarkerAlt} /> {ville}
        </div>
        <div className="btn-group">
          {props.myLogin === props.userProfil ? (
            <div className="button-container">
              <button className="btn btn-secondary" onClick={getListAmis}>
                Amis
              </button>
              <button className="btn btn-secondary" onClick={handleEdit}>
                Editer le profil
              </button>
              <button className="btn btn-secondary" onClick={handleDelete}>
                Supprimer mon compte
              </button>
            </div>
          ) : isAbonne === false ? (
            <div className="button-container">
              <button className="btn btn-secondary" onClick={getListAmis}>
                Amis
              </button>
              <button className="btn btn-secondary" onClick={Follow}>
                Suivre
              </button>
            </div>
          ) : isAbonne === true ? (
            <div className="button-container">
              <button className="btn btn-secondary" onClick={getListAmis}>
                Amis
              </button>
              <button className="btn btn-secondary" onClick={unFollow}>
                Ne plus suivre
              </button>
            </div>
          ) : (
            <div className="button-container">
              <button className="btn btn-secondary" onClick={getListAmis}>
                Amis
              </button>
              <button className="btn btn-secondary" onClick={unFollow}>
                Ne plus suivre
              </button>
            </div>
          )}
        </div>
      </div>

      {afficherAmis && (
        <div className="liste-amis">
          <h2>Amis</h2>
          <ListeProfils profils={amis} setCurrentPage={props.setCurrentPage} />
        </div>
      )}

      <div className="liste-playlists">
        <h2>Playlists :</h2>
        <ListePlaylists
          userProfil={props.userProfil}
          myLogin={props.myLogin}
          setCurrentPage={props.setCurrentPage}
          playlists={playlists}
          setPlaylists={setPlaylists}
          myPage={true}
        />
      </div>
    </div>
  );
}
