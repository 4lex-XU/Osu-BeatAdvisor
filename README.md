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

# Données 

## User

| id_user | email                              | pseudo | password | playlists |
| ------- | ---------------------------------- | ------ | -------- | --------- |
| 0       | alex.xu@etu.sorobnne-universite.fr | Alex   | "4Fd8Gh" | [0, 1 ]   |

## Playlist

| id_playlist | length | mode     | genres      | languages                | tags              | difficulty | 
| ----------- | ------ | -------- | ----------- | ------------------------ | ----------------- | ------------ |
| 0           | 50     | all      | [Jeu vidéo] | [Instumental]            | [pokemon, ..]     | 1.02 - 3.13  |
| 1           | 20     | fruits   | [Anime]     | [Japanese, Instrumental] | [dragon ball, ..] | 2.11 - 4.51  |
| 2           | 10     | mania    | [Rock]      | [English]                | [pink floyd, ..]  | 1.51 - 5.34  |

