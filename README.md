# Paring Test

ICC (<https://play.chessclub.com/>) is a player versus player (PvP) chess platform that matches players based on their chess game ratings. Each player has ratings for different game types: bullet, blitz, rapid, and classic (games lasting more than 30 minutes). Each game type has its time controls. For example, `blitz 3+2` means a 3-minute game with a 2-second increment (each player gains an extra 2 seconds on their clock after every move).

Player information is represented as follows:

```json
{
  "id": "", // unique identifier
  "name": "",
  "bullet_rating": 1500,
  "blitz_rating": 1500,
  "rapid_rating": 1500,
  "classic_rating": 1500
}
```

## Pairing requirements

### User flow

The user clicks one of the game buttons (refer to the screenshot below). When the user clicks on a button, a pulsing animation appears on the corresponding button, indicating that the system is searching for an opponent for the user within that specific time control. A second click on the same button cancels the request. When an opponent is found, the user starts a game with the chosen opponent.

### Requirements

1. Users can have multiple pairing requests.
2. Users can create or cancel requests at any given moment.
3. Users can only be paired/play with one user at a time.
4. Users are free to be paired again only when their previous game has ended.
5. Users will be paired with other users within a ±150 rating for the corresponding game type rating.
6. All users start with 1500 ratings in all game types.

### Assumptions

1. The chess logic already exists, and once we find a pair, the users are directed to another service called "game-service" where they can play the game.
2. When the user wins or loses a game, their rating goes up or down by 10 points.

## Task

Design a system to pair users based on the given requirements. Your design should include data structure design, APIs/Endpoints, and diagrams to explain the design