import {useEffect, useState} from "react";
import axios from "axios";

import '../CSS/editer-profil.css'

export default function EditerProfil(props) {
  const [newLogin, setNewLogin] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [newPasswordConfirm, setNewPasswordConfirm] = useState("");
  const [ville, setVille] = useState("");
  const [naissance, setNaissance] = useState("");
  const [description, setDescription] = useState("");
  const [PassOk, setPassOk] = useState(true);
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const getNewLogin = (evt) => {
    setNewLogin(evt.target.value);
  };
  const getNewPassword = (evt) => {
    setNewPassword(evt.target.value);
  };
  const getNewPasswordConfirm = (evt) => {
    setNewPasswordConfirm(evt.target.value);
  };
  const getNaissance = (evt) => {
    setNaissance(evt.target.value);
  };
  const getDescription = (evt) => {
    setDescription(evt.target.value);
  };
  const getVille = (evt) => {
    setVille(evt.target.value);
  };

  useEffect(() => {
    axios
      .get(`/user/${props.myLogin}`, {
        headers: {
          "Content-Type": "application/json",
        },
        withCredentials: true,
        credentials: "include",
      })
      .then((res) => {
        setLogin(res.data.Pseudo);
        setPassword(res.data.Password);
        setNaissance(res.data.Naissance);
        setDescription(res.data.Description);
        setVille(res.data.Ville);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  }, []);

  const Edit = (evt) => {
    evt.preventDefault();
    if (newPassword === newPasswordConfirm) {
      const data = new URLSearchParams();
      data.append("login", newLogin);
      data.append("password", newPassword);
      data.append("ville", ville);
      data.append("naissance", naissance);
      data.append("description", description);
      console.log(data);
      axios
        .put(`/user/edit`, data, {
          headers: {
            "Content-Type": "application/x-www-form-urlencoded",
          },
          withCredentials: true,
          credentials: "include",
        })
        .then((res) => {
          console.log(res.data);
          if (newLogin !== "") {
            props.setMyLogin(newLogin);
            props.setCurrentPage(newLogin);
            return;
          }
          props.setCurrentPage(props.myLogin);
        })
        .catch((err) => {
          console.log(err.response.data);
          setError(err.response.data);
        });
    } else {
      setPassOk(false);
    }
  };

  const pageProfilHandler = (evt) => {
    evt.preventDefault();
    props.setCurrentPage(props.myLogin);
  };

  return (
      <div className="EditerProfil">
        <form onSubmit={Edit}>
          <div className="form-columns">
            {/* Colonne Gauche */}
            <div className="form-left">
              <label htmlFor="newLogin">Pseudo</label>
              <input type="text" id="newLogin" onChange={getNewLogin} placeholder={login} />

              <label htmlFor="newDesciption">Bio</label>
              <textarea id="newDesciption" onChange={getDescription} placeholder={description} />

              <label htmlFor="newCountry">Ville</label>
              <input type="text" id="newCountry" onChange={getVille} placeholder={ville} />

              <label htmlFor="newDateBirth">Date de naissance</label>
              <input type="date" id="newDateBirth" onChange={getNaissance} placeholder={naissance} />
            </div>

            {/* Colonne Droite */}
            <div className="form-right">
              <label htmlFor="newPassword">Mot de passe</label>
              <input type="password" id="newPassword" onChange={getNewPassword} placeholder="******" />

              <label htmlFor="newPasswordConfirm">Confirmation</label>
              <input type="password" id="newPasswordConfirm" onChange={getNewPasswordConfirm} placeholder="******" />

              <button type="submit">Valider</button>

              {!PassOk && <p style={{ color: "red" }}>Les mots de passe ne correspondent pas</p>}
              {error && <p style={{ color: "red" }}>{error.message} {error.detail}</p>}
              <a className="pageProfil" href="a" onClick={pageProfilHandler}>Retour</a>
            </div>
          </div>
        </form>
      </div>
  );

}
