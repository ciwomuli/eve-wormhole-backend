package main

import (
	"eve-wormhole-backend/go/routes"
	"eve-wormhole-backend/go/service/ESI"

	"github.com/gomodule/redigo/redis"
)

//	"service/ESI"
// Update the import path below to the correct location of your ESI package, for example:

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		panic(err) // Handle error appropriately in production code
	}
	defer conn.Close() // Ensure the connection is closed when done
	ESI.InitESI(
		conn,
		"your-client-id",
		"your-client-secret",
		"http://your-callback-url.com/callback",
		[]string{"esi-skills.read_skills.v1", "esi-characters.read_contacts.v1"},
	)

	r := routes.SetRouter()

	r.Run(":8081")
}
