package infrared

import (
	"github.com/SierraSoftworks/Infrared/lib/config"
	"github.com/SierraSoftworks/Infrared/lib/store"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"net/http"
)

type InfraredServer struct {
	config        *config.Server
	http          *rest.Api
	udpConnection *net.UDPConn
}

func Setup(config *config.Server) *InfraredServer {
	server := InfraredServer{}

	server.config = config
	store.SetConfig(config)

	if server.config == nil {
		log.Fatal("Server configuration not populated")
	}

	err := SetupHttpEndpoint(&server)

	if err != nil {
		log.Fatal(err)
	}

	err = SetupUdpEndpoint(&server)

	if err != nil {
		log.Fatal(err)
	}

	if server.http == nil {
		log.Fatal("Server HTTP endpoint not populated")
	}

	if server.udpConnection == nil {
		log.Fatal("Server UDP endpoint not populated")
	}

	return &server
}

func SetupHttpEndpoint(server *InfraredServer) error {
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
		return err
	}

	api.SetApp(router)
	server.http = api

	return nil
}

func SetupUdpEndpoint(server *InfraredServer) error {
	addr, err := net.ResolveUDPAddr("udp", server.config.ListenOn)
	if err != nil {
		return err
	}

	connection, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	server.udpConnection = connection

	return nil
}

func (server InfraredServer) Start() {
	httpComplete := make(chan bool)
	udpComplete := make(chan bool)

	go StartHttpEndpoint(&server, httpComplete)
	go StartUdpEndpoint(&server, udpComplete)

	<-httpComplete
	<-udpComplete
}

func StartHttpEndpoint(server *InfraredServer, complete chan bool) {
	err := http.ListenAndServe(server.config.ListenOn, server.http.MakeHandler())
	if err != nil {
		log.Fatalf("Failed to start HTTP endpoint: %s", err)
	}

	complete <- true
}

func StartUdpEndpoint(server *InfraredServer, complete chan bool) {
	buf := make([]byte, 1024)
	heartbeat := &Heartbeat{}
	for {
		n, _, err := server.udpConnection.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Failed to read UDP packet: %s", err)
		}

		err = proto.Unmarshal(buf[0:n], heartbeat)
		if err != nil {
			log.Printf("Failed to deserialize protbuf packet '%s': %s", string(buf[0:n]), err)
		}

		err = store.UpdateNodeEntryLastSeen(heartbeat.GetNodeType(), heartbeat.GetId())
		if err != nil {
			log.Printf("Failed to update node last seen entry for %s: %s", heartbeat.GetId(), err)
		}
	}
	complete <- true
}
