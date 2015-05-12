package infrared

import (
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

func Setup() http.Handler {
	api := rest.NewApi()
	api.Use(&rest.ContentTypeCheckerMiddleware{})
	api.Use(&rest.JsonIndentMiddleware{})
	api.Use(&rest.PoweredByMiddleware{"Sierra Softworks Technology"})

	router, err := rest.MakeRouter(
		rest.Get("/api/v1/#node_type/config", ConfigGet),
		rest.Put("/api/v1/#node_type/config", ConfigSet),

		rest.Get("/api/v1/#node_type", NodeList),
		rest.Post("/api/v1/#node_type", NodeCreate),

		rest.Get("/api/v1/#node_type/#id/heartbeat", NodeHeartbeat),
		rest.Get("/api/v1/#node_type/#id", NodeGet),
		rest.Put("/api/v1/#node_type/#id", NodeUpdate),
		rest.Delete("/api/v1/#node_type/#id", NodeRemove),
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)
	return api.MakeHandler()
}
