import { useState } from 'react';
import axios from 'axios';

export default function Login(props) {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
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
    data.append('login', login);
    data.append('password', password);
    axios
      .post('/user/login', data, {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((response) => {
        console.log(response.data);
        props.setMyLogin(login);
        props.getConnected();
        props.setCurrentPage('home_page');
      })
      .catch((error) => {
        console.log(error.response.data);
        setError(error.response.data);
      });
  };
  const registerHandler = (evt) => {
    evt.preventDefault();
    props.setCurrentPage('register_page');
  };
  return (
    <div style={{ marginTop: '30px' }}>
      <form>
        <label htmlFor="login">Pseudo</label>
        <input id="login" onChange={getLogin} />
        <label htmlFor="mdp">Password</label>
        <input type="password" id="mdp" onChange={getPassword} />
        <div>
          <button onClick={submissionHandler}>Log In</button>
          <button type="reset">Reset</button>
        </div>
        {error && (
          <p style={{ color: 'red' }}>
            {error.message} {error.detail}
          </p>
        )}
        <a className="register" href="a" onClick={registerHandler}>
          pas encore inscrit ?
        </a>
      </form>
    </div>
  );
}
