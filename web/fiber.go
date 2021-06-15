package server

import (
	"net/http"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/cmelgarejo/minesweeper-svc/database"
	"github.com/cmelgarejo/minesweeper-svc/database/repo"
	_ "github.com/cmelgarejo/minesweeper-svc/docs"
	"github.com/cmelgarejo/minesweeper-svc/utils"
	"github.com/cmelgarejo/minesweeper-svc/utils/config"
	"github.com/cmelgarejo/minesweeper-svc/utils/logger"
	"github.com/cmelgarejo/minesweeper-svc/web/handlers/game"
	"github.com/cmelgarejo/minesweeper-svc/web/handlers/ping"
	"github.com/cmelgarejo/minesweeper-svc/web/middleware"
	"github.com/cmelgarejo/minesweeper-svc/web/services"
	"github.com/cmelgarejo/minesweeper-svc/web/services/common"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/loopcontext/msgcat"
)

// @title Minesweeper API
// @version 1.0.8
// @description A minesweeper game API
// @contact.name Christian Melgarejo
// @contact.email cmelgarejo.dev@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
func InitFiberServer(cfg *config.Config, log *logger.Logger,
	catalog *msgcat.MessageCatalog, db *database.DB) (app *fiber.App, err error) {
	app = fiber.New()
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "${time} - ${ip} - ${pid} - ${locals:requestid} - ${ua} - ${status} ${method} ${path}​ [url: ${url}]​​ \n",
	}))

	app.Use("/swagger", swagger.Handler) // swagger

	err = setupRoutes(app, cfg, log, *catalog, db)

	return app, err
}

func setupRoutes(app *fiber.App, cfg *config.Config, log *logger.Logger, catalog msgcat.MessageCatalog, db *database.DB) (err error) {
	// Default handler
	pingHandler := adaptor.HTTPHandlerFunc(ping.Ping)
	// Repos
	authRepo := repo.NewAuthRepoSvc(db, *log, catalog)
	gameRepo := repo.NewGameRepoSvc(db)
	// HTTP server's request and response service
	requestHelperSvc := common.NewRequestHelperSvc(*log, catalog)
	responseHelperSvc := common.NewResponseHelperSvc(*log, catalog, cfg.Server.ResponseContentType)
	// HTTP server's helpers
	authSvc := services.NewAuthSvcImpl(catalog)
	// Auth middleware
	authMiddleware := middleware.NewAuthMiddlewareSvc(*log, catalog, authRepo, responseHelperSvc)
	apiKeyMiddleware := HTTPMiddleware(authMiddleware.CheckUserAPIKey)
	// Game
	gameHandler := game.NewGameHandlerSvc(*log, catalog, gameRepo, authSvc, requestHelperSvc, responseHelperSvc)
	gameCreate := adaptor.HTTPHandlerFunc(gameHandler.Create)
	gameRead := adaptor.HTTPHandlerFunc(gameHandler.Read)
	gameUpdate := adaptor.HTTPHandlerFunc(gameHandler.Update)

	app.Get("/", pingHandler)
	app.Get("/ping", pingHandler)
	app.Get("/health", pingHandler)

	// Middleware makes sure only authorized users are allowed to use these resources
	api := app.Group("/api", apiKeyMiddleware)

	// Game
	gameRoute := api.Group("/game")
	gameRoute.Put("/", gameCreate)
	gameRoute.Get("/:id", gameRead)
	gameRoute.Put("/:id/click/:clickType/:row/:col", gameUpdate)

	return err
}

// HTTPMiddleware wraps net/http middleware to fiber middleware
func HTTPMiddleware(mw func(http.Handler) http.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var next bool
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next = true
			// Convert again in case request may modify by middleware
			c.Locals(string(utils.CurrrentUserCtxKey), r.Context().Value(utils.CurrrentUserCtxKey))
			c.Request().Header.SetMethod(r.Method)
			c.Request().SetRequestURI(r.RequestURI)
			c.Request().SetHost(r.Host)
			for key, val := range r.Header {
				for _, v := range val {
					c.Request().Header.Set(key, v)
				}
			}
		})
		_ = adaptor.HTTPHandler(mw(nextHandler))(c)
		if next {
			return c.Next()
		}
		return nil
	}
}
