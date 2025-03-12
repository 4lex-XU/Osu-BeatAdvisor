import Playlist from "./Playlist";
import axios from "axios";

export default function ListePlaylists(props) {
  const playlists = props.playlists;
  const setPlaylists = props.setPlaylists;

  const handleDelete = (playlistId) => {
    axios
      .delete(`/playlist/delete/${playlistId}`, {
        headers: {
          "Content-Type": "application/json",
        },
        withCredentials: true,
        credentials: "include",
      })
      .then((response) => {
        const updatedPlaylists = playlists.filter(
          (playlist) => playlist.Playlist_id !== playlistId,
        );
        setPlaylists(updatedPlaylists);
      })
      .catch((error) => {
        console.error("Erreur lors de la suppression de la playlist", error);
      });
  };

  return (
    <div className="list-container">
      <ul>
        {playlists.map((playlist) => (
          <li key={playlist.Playlist_id}>
            <Playlist
              playlistId={playlist.Playlist_id}
              playlist={playlist}
              title={playlist.Title}
              handleDelete={handleDelete}
              myLogin={props.myLogin}
              userProfil={props.userProfil}
              setCurrentPage={props.setCurrentPage}
              myPage={props.myPage}
            />
          </li>
        ))}
      </ul>
    </div>
  );
}
