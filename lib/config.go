package infrared

import (
	"github.com/SierraSoftworks/Infrared/lib/store"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

func ConfigGet(res rest.ResponseWriter, req *rest.Request) {
	node_type := req.PathParam("node_type")
	config, err := store.GetConfigEntry(node_type)

	if err != nil {
		apiError := APIError{}
		apiError.FromQueryError(err).ToResponse(res)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.WriteJson(&config)
	}
}

func ConfigSet(res rest.ResponseWriter, req *rest.Request) {
	node_type := req.PathParam("node_type")

	var config interface{}
	req.DecodeJsonPayload(&config)

	err := store.SetConfigEntry(node_type, config)

	if err != nil {
		apiError := APIError{}
		apiError.FromQueryError(err).ToResponse(res)
	} else {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.WriteJson(map[string]interface{}{"code": 200, "message": "Configuration entry persisted to database."})
	}
}
