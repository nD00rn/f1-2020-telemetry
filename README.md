# f1-2020-telemetry

This is a personal project in which I capture the F1 2020 telemetry game data and use that when my friend and I race against each other.

# Current state

All current features shown in the terminal itself and there is no user interaction possible from the terminal itself.

* Shows lap data of all participants
* Highlight users who have DRS active
* Show delta between users (leader and to next car)
* Show delta between the two players
* Show players who are currently in the pit lane
* Show flag colours in different marshal zones

# Documentation

* [F1 2020 UDP specification](https://forums.codemasters.com/topic/50942-f1-2020-udp-specification/?tab=comments#comment-515239)

# Ideas for the future

- [ ] Better delta calculations
  - Current implementation is working okay but has some flaws. If people know where to have documentation of how delta calculation could be implemented please let me know.

- [x] REST API support
  - For simple data retrieval of the game while is active
  - Returns the current state of the state machine in JSON format via the `/rest` endpoint
  - Has ability to fetch rendering of screenshot via `/terminal` end point

- [x] WebSocket support
  - This with the meaning that this program parses all the data, and a web browser could make it look pretty.
  - Returns the current state of the state machine in JSON format via the `/socket` endpoint

- [ ] Add support for F1 2021

# Screenshots

![Screenshot](images/screenshot-01.png "terminal screenshot")
