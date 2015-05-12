package infrared

import (
	"github.com/SierraSoftworks/Infrared/lib/store"
	"github.com/ant0ine/go-json-rest/rest"
	"net/http"
)

func NodeHeartbeat(res rest.ResponseWriter, req *rest.Request) {
	node_type := req.PathParam("node_type")

	if node_type == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node type in your request.").ToResponse(res)
		return
	}

	id := req.PathParam("id")
	if id == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node identifier in your request.").ToResponse(res)
		return
	}

	err := store.UpdateNodeEntryLastSeen(node_type, id)

	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		apiError := APIError{}
		apiError.FromQueryError(err).ToResponse(res)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}

func NodeGet(res rest.ResponseWriter, req *rest.Request) {
	node_type := req.PathParam("node_type")
	if node_type == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node type in your request.").ToResponse(res)
		return
	}

	id := req.PathParam("id")
	if id == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node identifier in your request.").ToResponse(res)
		return
	}

	node, err := store.GetNodeEntry(node_type, id)

	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		apiError := APIError{}
		apiError.FromQueryError(err).ToResponse(res)
	} else {
		res.WriteHeader(http.StatusOK)
		res.WriteJson(node)
	}
}

func NodeList(res rest.ResponseWriter, req *rest.Request) {
	node_type := req.PathParam("node_type")
	if node_type == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node type in your request.").ToResponse(res)
		return
	}

	nodes, err := store.GetNodeEntries(node_type)

	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		apiError := APIError{}
		apiError.FromQueryError(err).ToResponse(res)
	} else {
		res.WriteHeader(http.StatusOK)
		res.WriteJson(nodes)
	}
}

func NodeCreate(res rest.ResponseWriter, req *rest.Request) {
	node_type := req.PathParam("node_type")
	if node_type == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node type in your request.").ToResponse(res)
		return
	}

	createRequest := store.NodeCreateRequest{}

	req.DecodeJsonPayload(&createRequest)

	if createRequest.Hostname == "" {
		createRequest.Hostname = req.RemoteAddr
	}

	node, err := store.CreateNodeEntry(node_type, createRequest)

	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		apiError := APIError{}
		apiError.FromQueryError(err).ToResponse(res)
	} else {
		res.WriteHeader(http.StatusOK)
		res.WriteJson(node)
	}
}

func NodeUpdate(res rest.ResponseWriter, req *rest.Request) {
	node_type := req.PathParam("node_type")
	if node_type == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node type in your request.").ToResponse(res)
		return
	}

	id := req.PathParam("id")
	if id == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node identifier in your request.").ToResponse(res)
		return
	}

	createRequest := store.NodeCreateRequest{}

	req.DecodeJsonPayload(&createRequest)

	if createRequest.Hostname == "" {
		createRequest.Hostname = req.RemoteAddr
	}

	err := store.UpdateNodeEntry(node_type, id, createRequest)

	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		apiError := APIError{}
		apiError.FromQueryError(err).ToResponse(res)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}

func NodeRemove(res rest.ResponseWriter, req *rest.Request) {
	node_type := req.PathParam("node_type")
	if node_type == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node type in your request.").ToResponse(res)
		return
	}

	id := req.PathParam("id")
	if id == "" {
		NewAPIError().code(400).title("Bad Request").message("You failed to provide a valid node identifier in your request.").ToResponse(res)
		return
	}

	err := store.RemoveNodeEntry(node_type, id)

	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		apiError := APIError{}
		apiError.FromQueryError(err).ToResponse(res)
	} else {
		res.WriteHeader(http.StatusOK)
	}
}
