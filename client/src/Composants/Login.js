import { useState } from "react";
import axios from "axios";

export default function Login(props) {
  const [login, setLogin] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState(null);

  const getLogin = (evt) => {
    setLogin(evt.target.value);
  };
  const getPassword = (evt) => {
    setPassword(evt.target.value);
  };

  const submissionHandler = (evt) => {
    evt.preventDefault();
    setError(null);
    const data = new URLSearchParams();
    data.append("login", login);
    data.append("password", password);
    axios
      .post("/user/login", data, {
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        withCredentials: true,
        credentials: "include",
      })
      .then((response) => {
        console.log(response.data);
        props.setMyLogin(login);
        props.getConnected();
        props.setCurrentPage("home_page");
      })
      .catch((error) => {
        console.log(error.response.data);
        setError(error.response.data);
      });
  };
  const registerHandler = (evt) => {
    evt.preventDefault();
    props.setCurrentPage("register_page");
  };
  return (
    <div style={{ marginTop: "30px" }}>
      <form>
        <label htmlFor="login">Nom d'utilisateur</label>
        <input id="login" onChange={getLogin} />
        <label htmlFor="mdp">Mot de passe</label>
        <input type="password" id="mdp" onChange={getPassword} />
        <div>
          <button onClick={submissionHandler}>Me connecter</button>
        </div>
        {error && (
          <p style={{ color: "red" }}>
            {error.message} {error.detail}
          </p>
        )}
        <span>
          Pas encore inscrit ?
          <a className="register" href="a" onClick={registerHandler}>
            S'inscrire
          </a>
        </span>
      </form>
    </div>
  );
}
