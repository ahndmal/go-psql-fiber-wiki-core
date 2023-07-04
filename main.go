package main

import (
	"context"
	"github.com/go-pg/pg/v11"
	//"github.com/go-pg/pg/v11/orm"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
	"log"
	"fmt"
	"time"
	"strconv"
	"os"
)

const ALL_PAGES = "SELECT * FROM pages"

type Page struct {
	Id          int64 //`pg:",discard_unknown_columns"`
	Title       string
	Body        string
	SpaceKey    string    `pg:"space_key"`
	ParentId    int64     `pg:"parent_id"`
	AuthorId    int64     `pg:"author_id"`
	CreatedAt   time.Time `pg:"created_at"`
	LastUpdated time.Time `pg:"-"`
}

func main() {
	app := fiber.New()

	pages := make([]Page, 0)

	db := pg.Connect(&pg.Options{
		User:     "dev",
		Password: os.Getenv("DB_PASS"),
		Database: "wiki1",
		Addr:     os.Getenv("DB_HOST"),
	})
	
	defer db.Close(context.Background())

	_, err := db.Query(context.Background(), &pages, ALL_PAGES)
	if err != nil {
		log.Fatalln(err)
	}

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

		//reqMethod := ctx.Method()

		hundrPages := make([]Page, limit)
		
		for i := 0; i < limit; i++ {
			hundrPages[i] = pages[i]
		}
		return ctx.JSON(hundrPages)
	})

	app.Get("/all-pages", func(ctx *fiber.Ctx) error {
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
	log.Println( fmt.Sprintf(">> Fiber started on port %d", 4000))
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

