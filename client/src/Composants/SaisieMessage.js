import { useState } from 'react';
import axios from 'axios';

export default function SaisieMessage(props) {
  const [text, setText] = useState('');
  const [error, setError] = useState('');
  const currentDate = new Date();
  const date = currentDate.toLocaleDateString('fr');
  const clock = currentDate.getHours() + ':' + currentDate.getMinutes();

  const getContent = (evt) => setText(evt.target.value);

  const submissionHandler = (evt) => {
    evt.preventDefault();
    setError(null);
    const data = {
      text: text,
      date: date + ' ' + clock,
    };
    axios
      .post(`playlist/comment/${props.playlistId}`, data, {
        headers: {
          'Content-Type': 'application/json',
        },
        withCredentials: true,
        credentials: 'include',
      })
      .then((res) => {
        console.log(res.data);
        setText('');
        props.setCommentaires([...props.commentaires, data]);
      })
      .catch((err) => {
        console.log(err.response.data);
        setError(err.response.data);
      });
  };

  return (
    <div className="message">
      <form>
        <textarea
          className="texte-msg"
          rows="5"
          cols="33"
          placeholder="Entrez votre message ici"
          onChange={getContent}
          value={text}
        />
        <p className="date">
          {clock}
          {' · '}
          {date}
        </p>
        {error && (
          <p style={{ color: 'red', fontSize: '12px' }}>
            {error.message} {error.detail}
          </p>
        )}
        <button onClick={submissionHandler}>Post</button>
      </form>
    </div>
  );
}
