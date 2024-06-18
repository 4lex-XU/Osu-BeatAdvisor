import { Navbar, Nav } from 'react-bootstrap';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import ListePlaylistIcon from '../Images/music-note-list.svg';

import {
  faHome,
  faUser,
  faEdit,
  faPowerOff,
} from '@fortawesome/free-solid-svg-icons';

export default function NavigationPanel(props) {
  return (
    <Navbar expand="lg">
      <Navbar.Toggle />
      <Navbar.Collapse id="navbar">
        <Nav className="navigation">
          {props.isConnected ? (
            <>
              <Nav.Link
                href="#"
                onClick={() => props.setCurrentPage('home_page')}
              >
                <FontAwesomeIcon icon={faHome} /> Accueil
              </Nav.Link>
              <Nav.Link
                href="#"
                onClick={() => props.setCurrentPage(props.myLogin)}
              >
                <FontAwesomeIcon icon={faUser} /> Profil
              </Nav.Link>
              <Nav.Link
                href="#"
                onClick={() => props.setCurrentPage('all_playlists_page')}
              >
                <img
                  src={ListePlaylistIcon}
                  alt="Icône Liste de Lecture"
                  style={{ width: 24, height: 24 }}
                />
                <span> En ce moment</span>
              </Nav.Link>
              <Nav.Link href="#" onClick={props.setLogout}>
                <FontAwesomeIcon icon={faPowerOff} /> Déconnexion
              </Nav.Link>
            </>
          ) : (
            <>
              <Nav.Link
                href="#"
                onClick={() => props.setCurrentPage('register_page')}
              >
                Inscription
              </Nav.Link>
              <Nav.Link
                href="#"
                onClick={() => props.setCurrentPage('login_page')}
              >
                Connexion
              </Nav.Link>
            </>
          )}
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
}
