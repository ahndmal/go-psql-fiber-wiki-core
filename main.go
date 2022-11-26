package main

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v11"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/websocket"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
	"log"
	"net/http"
	"os"
	"time"
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
		User:     os.Getenv("DB_USER"),
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
		hundrPages := make([]Page, 100)
		for i := 0; i < 100; i++ {
			hundrPages[i] = pages[i]
		}
		return ctx.JSON(hundrPages)
	})

	app.Get("/all-pages", func(ctx *fiber.Ctx) error {
		return ctx.JSON(pages)
	})

	listErr := app.Listen(":4000")
	if err != nil {
		log.Fatalln(listErr)
	}
	rSocketInit()
	initSockets()
}

func rSocketInit() {
	log.Println(">> initiating RSocket")
	PORT := ":7878"
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

func initSockets() {
	log.Println(">> initiating Sockets")
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	http.HandleFunc("/echo", func(writer http.ResponseWriter, req *http.Request) {
		conn, _ := upgrader.Upgrade(writer, req, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":8090", nil)
}
