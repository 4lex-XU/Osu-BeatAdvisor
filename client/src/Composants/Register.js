import { useState } from 'react';
import axios from 'axios';

export default function Register(props) {
  const [email, setEmail] = useState('');
  const [pseudo, setPseudo] = useState('');
  const [passOK, setPassOK] = useState(false);
  const [pass, setPass] = useState('');
  const [passVerif, setPassVerif] = useState('');
  const [error, setError] = useState(null);

  const getEmail = (evt) => {
    setEmail(evt.target.value);
  };
  const getPseudo = (evt) => {
    setPseudo(evt.target.value);
  };
  const getPass = (evt) => {
    setPass(evt.target.value);
  };
  const getPassVerif = (evt) => {
    setPassVerif(evt.target.value);
  };
  const submissionHandler = (evt) => {
    if (pass === passVerif) setPassOK(false);
    else setPassOK(true);
    evt.preventDefault();
    setError(null);
    if (!passOK) {
      const data = new URLSearchParams();
      data.append('email', email);
      data.append('pseudo', pseudo);
      data.append('password', pass);
      console.log(data);
      axios
        .put('/user/register', data, {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
          },
        })
        .then((response) => {
          console.log(response.data);
          props.setCurrentPage('login_page');
        })
        .catch((error) => {
          console.log(error.response.data);
          setError(error.response.data);
        });
    }
  };

  const loginHandler = (evt) => {
    evt.preventDefault();
    props.setCurrentPage('login_page');
  };

  return (
    <div style={{ marginTop: '30px' }}>
      <form name="register">
        <label htmlFor="pseudo">Pseudo</label>
        <input id="pseudo" onChange={getPseudo} />
        <label htmlFor="register_login">Email</label>
        <input id="register_login" onChange={getEmail} />
        <label htmlFor="register_mdp1">Password</label>
        <input type="password" id="register_mdp1" onChange={getPass} />
        <label htmlFor="register_mdp2">Confirm Password</label>
        <input type="password" id="register_mdp2" onChange={getPassVerif} />
        {passOK && (
          <p style={{ color: 'red' }}>Veuillez reconfirmer le mot de passe</p>
        )}
        {error && (
          <p style={{ color: 'red' }}>
            {error.message} {error.detail}
          </p>
        )}
        <div className="register-input">
          <button onClick={submissionHandler}>Sign In</button>
          <button type="reset">Reset</button>
        </div>
        <a className="connexion" href="a" onClick={loginHandler}>
          déja inscrit ?
        </a>
      </form>
    </div>
  );
}
