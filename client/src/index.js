import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import axios from 'axios';
import 'bootstrap/dist/css/bootstrap.min.css';
import MainPage from './Composants/MainPage';

const rootElement = document.getElementById('root');
const root = createRoot(rootElement);
axios.defaults.baseURL = 'http://localhost:8080/api';

root.render(
  <StrictMode>
    <div className="mainpage">
      <MainPage />
    </div>
  </StrictMode>
);
