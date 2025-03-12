import React, { useEffect, useState } from "react";
import { faTrash } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import axios from "axios";

export default function Message(props) {
  const [isMenuOpen, setMenuOpen] = useState(false);
  const [isAbonne, setIsAbonne] = useState(null);

  // PERMET DE CHANGER DE PAGE
  const pageProfilHandler = (evt) => {
    evt.preventDefault();
    props.setCurrentPage(props.author);
  };

  // PERMET D'OUVRIR UN MENU D'OPTION SUR UN MESSAGE
  const toggleMenu = () => {
    setMenuOpen(!isMenuOpen);
  };

  // PERMET DE SUPPRIMER UN MESSAGE
  const deleteMessage = (evt) => {
    evt.preventDefault();
    axios
      .delete(`/user/${props.id}/messages`, {
        headers: {
          "Content-Type": "application/json",
        },
        withCredentials: true,
        credentials: "include",
      })
      .then((res) => {
        console.log(res.data);
        props.setCommentaires(
          props.commentaires.filter(
            (commentaire) => commentaire.Id !== props.Id,
          ),
        );
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  };

  // Au chargement de la page, on détermine si l'utilisateur est abonné au profil
  useEffect(() => {
    axios
      .get(`/user/${props.myLogin}/friends/${props.login}`, {
        headers: {
          "Content-Type": "application/json",
        },
        withCredentials: true,
        credentials: "include",
      })
      .then((res) => {
        setIsAbonne(res.data);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  }, [isAbonne, props.author]);

  // Permet de suivre un profil
  const Follow = (evt) => {
    evt.preventDefault();
    const data = {
      friend_login: props.author,
    };
    axios
      .put(`/user/${props.myLogin}/newfriend`, data, {
        headers: {
          "Content-Type": "application/json",
        },
        withCredentials: true,
        credentials: "include",
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
      .delete(`/user/${props.myLogin}/friends/${props.author}`, {
        headers: {
          "Content-Type": "application/json",
        },
        withCredentials: true,
        credentials: "include",
      })
      .then((res) => {
        console.log(res.data);
        setIsAbonne(false);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  };

  return (
    <article className="message">
      <div>
        <a className="user" href="a" onClick={pageProfilHandler}>
          {props.author}
        </a>
        <button className="menuBurger" onClick={toggleMenu}>
          ...
        </button>
        {isMenuOpen && (
          <div className="menuList">
            {props.myLogin === props.author ? (
              <div>
                <button onClick={deleteMessage}>
                  <FontAwesomeIcon icon={faTrash} />
                </button>
              </div>
            ) : (
              <div>
                {isAbonne === false ? (
                  <button onClick={Follow}>Suivre</button>
                ) : (
                  <button onClick={unFollow}>Ne plus suivre</button>
                )}
              </div>
            )}
          </div>
        )}
      </div>
      <div className="content">
        <textarea
          className="texte-msg"
          rows="5"
          cols="33"
          readOnly="readonly"
          value={props.text}
        />
      </div>
      <p className="date">{props.date}</p>
    </article>
  );
}
