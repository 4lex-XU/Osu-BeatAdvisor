Alex XU & Antoine LUONG

# Osu! - BeatAdvisor

Le site web *Osu! - BeatAdvidor* vous permet de générer des playlists basées sur vos préferences ! 
Il contient une composante sociale, vous permettant de partager vos playslists favorites, de commenter celles des autres et de les liker.

# API Web : osu!web

Accès à l'API ici : https://osu.ppy.sh/docs/index.html#

*Contextualisation : osu! est un jeu de rythme, dans lequel les joueurs lancent des stages qu'on appelle "beatmap" dont le but est d'obtenir le meilleur score possible*

Les beatmaps sont composées par les joueurs, ainsi les données de l'API sont régulièrement mis à jour.

L'APi propose des requêtes get sur :
- les beatmaps (Get Beatmaps, Get a User Beatmap score, Get Beatmap scores, etc...)
- les joueurs (Get User, Get User Scores, Get User Beatmaps, etc...)

Les informations des beatmaps que nous utiliserons pour créer des playlists concerneront principalement :
- le mode
- la difficulté 
- les tags
- le genre
- la langue 

# Fonctionnalités

- Les utilisateurs pourront filtrer selon le genre musical, le niveau de difficulté des chansons, les artistes préférés pour la génération de la playlist.
- La génération tirera aléatoirement les beatmaps dont les préférences correspondent. 
- Chaque jour le de nouvelles beatmaps seront ajoutées au site web.
- Les utilisateurs pourront partager, commenter, liker les playlists.

# Cas d'utilisation 

- L'utilisateur se connecte, il attérit sur la page d'accueil dans laquelle il peut générer une playlist en spécifiant ses préférences, en cochant les filtres proposés ou bien en ajoutant des tags.
- Le site pioche ensuite x beatmaps correspondant aux préférences de l'utilisateur, et propose la playlist à l'utilisateur. 
- L'utilisateur peut ensuite l'enregistrer (pour les télécharger en jeu), la partager.
- Après la génération d'une playlist, l'utilisateur peut refuser la playlist et demander d'en générer une autre et peut sélectionner les maps à garder lors de la nouvelle génération.
- Après la génération d'une playlist, l'utilisateur peut manuellement retirer des beatmaps et/ou en ajouter d'autres.
- L'utilisateur après s'être connecté peut observer toutes ses playlists.
- (L'utilisateur peut réorganiser l'ordre des musiques dans sa playlist selon ses préférences.)

# Données 

## User

| id_user | email                              | pseudo | password | playlists | friends |
| ------- | ---------------------------------- | ------ | -------- | --------- | ------- |
| 0       | alex.xu@etu.sorobnne-universite.fr | Alex   | "4Fd8Gh" | [0, 1 ]   | [0,1,2] |

## Playlist

| id_playlist |   title   |  url  | length| mode     | genres      | languages                | tags              | difficulty | 
| ----------- | --------- | ----- |------ | -------- | ----------- | ------------------------ | ----------------- | ------------ |
| 0           |  salle    | url1  |50     | all      | [Jeu vidéo] | [Instumental]            | [pokemon, ..]     | 1.02 - 3.13  |
| 1           |   eau     | url2  |20     | fruits   | [Anime]     | [Japanese, Instrumental] | [dragon ball, ..] | 2.11 - 4.51  |
| 2           |  2prohf   | url3  |10     | mania    | [Rock]      | [English]                | [pink floyd, ..]  | 1.51 - 5.34  |

On peut cliquer sur le titre qui amène vers la playlist téléchargement.

# Mise à jour des données et appel à l'API externe

Les données sont mises à jour en temps réel lorsqu'un utilisateur crée une nouvelle playlist ou modifie ses préférences. L'API externe est appelée pour alimenter la base de données mettant à jour toutes les beatmaps du site. On pourrait appeler l'API tous les jours à 5:00 UTC+1.

# Description du Serveur

Le site utilise GO avec le package net/http, avec une approche basée sur des ressources, avec un ensemble de services RESTful pour gérer les utilisateurs, les playlists, et les interactions avec l'API d'Osu!. Chaque service est conçu pour gérer un aspect spécifique de l'application, comme la création de playlists, la gestion des utilisateurs, et le recueil des feedbacks.

## Ressources et Fonctionnalités Associées

- **Authentification (/auth)**:
  
    - POST /auth/login: Permet à un utilisateur de se connecter en envoyant son identifiant et mot de passe.
    - POST /auth/logout: Permet à un utilisateur de se déconnecter.
    - POST /auth/register: Permet à un nouvel utilisateur de s'inscrire en fournissant les informations nécessaires.

- **Utilisateurs (/users)**:

    - GET /users/{id}: Récupère les informations d'un utilisateur spécifique.
    - POST /users: Crée un nouvel utilisateur.
    - PUT /users/{id}: Met à jour les informations d'un utilisateur.
    - DELETE /users/{id}: Supprime un utilisateur.

- **Playlists (/playlists)**:

    - GET /playlists/{id}: Récupère une playlist spécifique.
    - GET /playlists: Récupère toutes les playlists selon les filtres appliqués (genre, langue, difficulté, etc.).
    - POST /playlists: Crée une nouvelle playlist basée sur les préférences de l'utilisateur.
    - PUT /playlists/{id}: Met à jour une playlist (ajouter/enlever des beatmaps, changer le titre, etc.).
    - DELETE /playlists/{id}: Supprime une playlist.

- **Beatmaps (/beatmaps)**:

    - GET /beatmaps: Récupère les beatmaps depuis l'API externe osu!web selon les critères spécifiés (mode, difficulté, tags, genre, langue).

- **Commentaires (/comments)**:

    - GET /comments: Récupère les commentaires sur une playlist spécifique.
    - PUT /comments: Ajoute un commentaire à une playlist.
    - DELETE /comments/{id}: Supprime un commentaire.

- **Likes (/likes)**:

    - POST /likes: Ajoute un like à une playlist.
    - DELETE /likes/{id}: Supprime un like d'une playlist.

# Description du Client

Le site est conçu avec React.js pour offrir une expérience utilisateur fluide. Les écrans incluent une page d'accueil, une interface de création de playlists, une page de découverte, et une section communauté. Les appels au serveur sont effectués à partir de ces différentes sections pour récupérer ou envoyer des données.

## Plan du Site et Contenu des Écrans

- **Page de connexion**:
    - Formulaire de connexion avec champs pour l'identifiant et le mot de passe avec bouton pour se connecter.
    - Lien vers la page d'inscription pour les nouveaux utilisateurs qui n'ont pas encore de compte.
    - Lien vers la récupération de mot de passe pour les utilisateurs qui ont oublié leur mot de passe.
    
- **Page d'Accueil**:
    - Génération de playlists basée sur des filtres spécifiés par l'utilisateur.
    - Affichage des playlists générées avec options pour sauvegarder, partager, et commenter.

- **Page de Playlist**:
    - Détails de la playlist sélectionnée avec liste des beatmaps, possibilité de modifier la playlist (ajouter/enlever des beatmaps, changer l'ordre, etc).

- **Profil Utilisateur**:
    - Affichage des informations de l'utilisateur, liste de ses playlists, et historique des commentaires.

## Appels Serveur

- Requête POST pour vérifier les informations de connexion de l'utilisateur. Cette requête vérifiera les informations saisies dans le formulaire de connexion et authentifiera l'utilisateur:
    - Si les informations de connexion sont valides, une réponse réussie (code 200) sera renvoyée, et l'utilisateur sera redirigé vers la page d'accueil.
    - En cas d'échec de la connexion, une réponse d'erreur (code 401 ou 403) sera renvoyée, et un message d'erreur approprié sera affiché à l'utilisateur sur la page de connexion.
- Les playlists sont générées en envoyant une requête POST à /playlists avec les préférences de l'utilisateur.
- Les détails d'une playlist sont récupérés via GET /playlists/{id}.
- Les commentaires sont ajoutés en envoyant une requête POST à /comments.
- Les likes sont gérés par POST et DELETE sur /likes et /likes/{id} respectivement.

# Requêtes et Réponses
  
| Nom du web service | URL du web service | Description du service | Paramètres en entrée | Format de sortie | Exemple de sortie |   Erreurs possibles   | Avancement du Service | Classes/Fichiers GO | Informations additionnelles |
| ---------------- | ------------------- | ---------------- | ----------------- | ----------------- | ---------------- |  -----------------------------------------------------  | -------------------- | -------------------- | -------------------- |
| Login | auth/login/ avec POST | Permet de récupérer une clef de connexion valide pendant un certain temps | login; password | JSON | - {"status": 200, "message": "Connexion réussie", "userid": "6442546cd354647be214"};<br/> - {"status": 400, "message": "Requête invalide: login et password nécessaires"};<br/>  - {"status": 401, "message": "Utilisateur inconnu"};<br/>  - {"status": 403, "message": "Login et/ou mot de passe invalide(s)"};<br/>  - {"status": 500, "message": "Erreur interne"};<br/> | - Champs manquants (400)<br/> - Utilisateur inconnu (401)<br/> - Login et/ou mot de passe invalide(s) (403)<br/> - Erreur interne (500) | En cours | api.go | 1. Si tout se passe bien renvoyer le code 200 et Connexion réussie, et création de la session avec req.session.userid<br/> 2. Tous les champs doivent être complétés, si non -> 400<br/> 3. Vérification de l’existence du login, si non -> 401<br/> 4. Vérification du login et mot de passe, si non -> 403 |
| Générer Playlist | /playlists avec POST | Permet de générer une playlist | les préférences (genre, difficulté, tags, etc) | JSON | - {"status": 200, "message": "Playlist générée avec succès"};<br/> - {"status": 500, "message": "Erreur interne"}; | - Utilisateur inconnu (401) | En cours | api.go | 1. Si tout se passe bien renvoyer le code 200 et Playlist générée, et création de la playlist |
| Récupérer une Playlist | /playlists/{id} avec GET | Permet de récupérer une Playlist | l'id de la Playlist | JSON, Tableau | - [{"date" : "17/04/2002", "auteur": "1d5s1s88d5f4s64f", "likes": 5, "commentaires": {"comm1: "}, "beatmaps" : {"beatmap1": "Symphony of the night"}, {beatmap2": ...}, {...}},  , ;<br/> - {"status": 404, "message": Playlist pas trouvée<br/> - {"status": 500, "message": "Erreur interne"}; | - Utilisateur inconnu (401)<br/> - Playlist pas trouvée (404) | En cours | api.go |  1. Si tout se passe bien renvoyer le code 200 et Playlist récupérée<br/> |
| Ajouter un commentaire | comments/ avec PUT | Permet de poster un commentaire pour une Playlist | parentId, login, date, clock, content | JSON | - {"status": 200, "message": "Commentaire posté avec succès"};<br/> - {"status": 400, "message": Commentaire vide"}; | - Non connecté (401)<br/> - Champs manquants (400)<br/> - Commentaire vide (400)<br/> - Erreur interne (500) | En cours | Messages/messages.go | - Vérification de la connexion -> 401<br/> - Vérification des champs manquants -> 400<br/> - Création du commentaire avec la fonction createComment<br/> - Si tout se passe bien -> 200 “Commentaire créé avec succès” |

# Schéma global du système

/osu-beatadvisor/
│
├── /api/                     # Centralized API interaction and user management
│   ├── api.go                # Main handler for all API requests, routing to specific functionalities
│
├── /models/                  # Data models for the application
│   ├── user.go               # User model
│   ├── playlist.go           # Playlist model
│   ├── beatmap.go            # Beatmap model
│   ├── comment.go            # Comment model
│   ├── like.go               # Like model
|   ├── share.go              # Share model
│
├── /util/                    # Utility functions and common definitions
│   ├── session.go            # Session management utilities
│   ├── error.go              # Error handling utilities
│   ├── response.go           # Common response formatting utilities
│
└── /osu!api/                 # Database interaction extracted from the osu!api website
    ├── osu.go                # Database connection and common operations