import { useEffect, useState } from "react";
import "../CSS/styles.css";
import "../CSS/formulaire.css";
import "../CSS/message.css";
import logo from "../Images/BeatAdvisorLogo.png";
import fond from "../Images/fond.jpg";
import axios from "axios";

import NavigationPanel from "./NavigationPanel.js";
import Login from "./Login";
import Register from "./Register";
import HomePage from "./HomePage";
import PageProfil from "./PageProfil";
import EditerProfil from "./EditerProfil";
import ListePlaylists from "./ListePlaylists";
import BackgroundAudio from "./BackgroundAudio.js";
import ost0 from "../Musiques/ost0.mp3";
import ost1 from "../Musiques/ost1.mp3";
import ost2 from "../Musiques/ost2.mp3";
import ost3 from "../Musiques/ost3.mp3";
import ost4 from "../Musiques/ost4.mp3";
import ost5 from "../Musiques/ost5.mp3";
import ost6 from "../Musiques/ost6.mp3";
import ost7 from "../Musiques/ost7.mp3";
import ost8 from "../Musiques/ost8.mp3";
import ost9 from "../Musiques/ost9.mp3";
import ost10 from "../Musiques/ost10.mp3";

export default function MainPage(props) {
  // states
  const [isConnected, setConnect] = useState(false);
  const [currentPage, setCurrentPage] = useState("login_page");
  const [myLogin, setMyLogin] = useState("");
  const [playlists, setPlaylists] = useState([]);
  const audios = [
    ost0,
    ost1,
    ost2,
    ost3,
    ost4,
    ost5,
    ost6,
    ost7,
    ost8,
    ost9,
    ost10,
  ];

  // comportements
  const getConnected = () => {
    setConnect(true);
    setCurrentPage("home_page");
  };
  const setLogout = () => {
    setConnect(false);
    setCurrentPage("login_page");
  };

  const handler = (evt) => {
    evt.preventDefault();
    if (isConnected) {
      setCurrentPage("home_page");
    } else if (currentPage === "register_page") {
      setCurrentPage("login_page");
    } else {
      setCurrentPage("register_page");
    }
  };

  useEffect(() => {
    fetchPlaylists();
  }, [currentPage]);

  const fetchPlaylists = () => {
    axios
      .get(`playlist/get/all`, {
        headers: {
          "Content-Type": "application/json",
        },
        withCredentials: true,
        credentials: "include",
      })
      .then((res) => {
        console.log(res.data);
        setPlaylists(res.data.playlists);
      })
      .catch((err) => {
        console.log(err.response.data);
      });
  };

  const backgroundStyle = {
    backgroundImage: `url(${fond})`,
    backgroundSize: "cover",
    backgroundPosition: "center",
    backgroundRepeat: "no-repeat",
    backgroundAttachment: "fixed",
    minHeight: "100vh",
  };

  // render
  return (
    <div style={backgroundStyle}>
      <header className="beatadvisor">
        <BackgroundAudio sources={audios} />
        <a className="logo" href="a" onClick={handler}>
          <img src={logo} />
        </a>
        <aside>
          <NavigationPanel
            myLogin={myLogin}
            setCurrentPage={setCurrentPage}
            currentPage={currentPage}
            setLogout={setLogout}
            isConnected={isConnected}
          />
        </aside>
      </header>
      <main>
        {currentPage === "register_page" ? (
          <Register setCurrentPage={setCurrentPage} />
        ) : isConnected ? (
          currentPage === "home_page" ? (
            <HomePage myLogin={myLogin} setCurrentPage={setCurrentPage} />
          ) : currentPage === "edit_page" ? (
            <EditerProfil
              myLogin={myLogin}
              setMyLogin={setMyLogin}
              setCurrentPage={setCurrentPage}
            />
          ) : currentPage === "all_playlists_page" ? (
            <div className="liste-playlists">
              <h2>En ce moment :</h2>
              <ListePlaylists
                userProfil={props.userProfil}
                myLogin={props.myLogin}
                setCurrentPage={setCurrentPage}
                playlists={playlists}
                setPlaylists={setPlaylists}
              />
            </div>
          ) : (
            <PageProfil
              myLogin={myLogin}
              userProfil={currentPage}
              setCurrentPage={setCurrentPage}
              setLogout={setLogout}
              myPage={false}
            />
          )
        ) : (
          <Login
            setMyLogin={setMyLogin}
            getConnected={getConnected}
            setCurrentPage={setCurrentPage}
          />
        )}
      </main>
    </div>
  );
}
