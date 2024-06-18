import Message from './Message';

export default function ListeMessages(props) {
  const commentaires = props.commentaires;
  return (
    <div className="listeMsg">
      {commentaires &&
        commentaires.map((commentaire, index) => (
          <Message
            key={index}
            id={commentaire.Id}
            author={commentaire.author}
            text={commentaire.text}
            date={commentaire.date}
            setCurrentPage={props.setCurrentPage}
            commentaires={commentaires}
            setCommentaires={props.setCommentaires}
            myLogin={props.myLogin}
          />
        ))}
    </div>
  );
}
