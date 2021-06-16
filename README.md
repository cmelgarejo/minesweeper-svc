# minesweeper-svc

A Minesweeper API

## Development Log

- Start design of the game engine
- Implemented game engine and service structure
- Implemented game creation and start
- Persistance in DB (pgsql)
- Migrations
- Fiber as underlying http server
- Simple auth, using API keys
- Swagger doc with examples and testing enabled

## Roadmap

- Implement the reveal of adjacent squares that are not in a mine proximity
- REST API
- Redis persistance (to replace in-memory persistance with a proper database interface, able to mock any persistance engine)
- User accounts
- Bonus: Multiplayer!

## What to Implement, in priority order

- ✅ Design and implement a documented RESTful API for the game (think of a mobile app for your API)
- ✔️ Implement an API client library for the API designed above. Ideally, in a different language, of your preference, to the one used for the API (swagger docs act as a client of sorts)
- ✅ When a cell with no adjacent mines is revealed, all adjacent squares will be revealed (and repeat)
- ✅ Ability to 'flag' a cell with a question mark or red flag
- ✅ Detect when game is over
- ✅ Persistence
- ✅ Time tracking
- ✅ Ability to start a new game and preserve/resume the old ones
- ✅ Ability to select the game parameters: number of rows, columns, and mines
- ✅ Ability to support multiple users/accounts

---

## Development

### Swagger docs

- Just add documentation to the handlers in `web/handlers/**` then run:

        swag init -g fiber.go -d web/

    to generate the new swagger information to test your changes

---

## Tech Stack

- [Go](https://golang.org)
- [Ginkgo](https://onsi.github.io/ginkgo/)
- [Swaggo](https://github.com/swaggo/swag)
- [Fiber](https://gofiber.io/)
- [Redis](https://redis.io)
- [Docker](https://docker.com)
