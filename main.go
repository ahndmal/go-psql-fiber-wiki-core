package main

import (
	"context"
	"fmt"
	"go-wiki-core/db"

	//"github.com/go-pg/pg/v11/orm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
	"log"
	"strconv"
)

const allPages = "SELECT * FROM pages"

func main() {
	pageRepo := db.PagesRepo{}

	app := fiber.New()

	// REST api
	app.Get("/pages", func(ctx *fiber.Ctx) error {
		pLimit := ctx.Query("limit")

		var limit int

		if pLimit == "" {
			limit = 100
		} else {
			i, err := strconv.Atoi(pLimit)
			if err != nil {
				log.Println(err)
			}
			limit = i
		}

		pages := pageRepo.GetNPages(context.Background(), limit)

		return ctx.JSON(pages)
	})

	app.Get("/all-pages", func(ctx *fiber.Ctx) error {
		pages := pageRepo.GetPages(context.Background())
		return ctx.JSON(pages)
	})

	// websockets
	app.Use("/ws", func(ctx *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(ctx) {
			ctx.Locals("allowed", true)
			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	//rSocketInit()

	log.Fatal(app.Listen(":4000"))
	log.Println(fmt.Sprintf(">> Fiber started on port %d", 4000))
}

func rSocketInit() {
	PORT := ":7878"
	log.Printf(">> initiating RSocket on port %s", PORT)
	err := rsocket.Receive().
		Acceptor(func(ctx context.Context, setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(func(msg payload.Payload) mono.Mono {
					return mono.Just(msg)
				}),
			), nil
		}).
		Transport(rsocket.TCPServer().SetAddr(PORT).Build()).
		Serve(context.Background())
	log.Fatalln(err)
}
