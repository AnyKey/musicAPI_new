package graphquery

import (
	"log"
	"musicAPI/graphql"
	"musicAPI/graphql/delivery/graphquery/graph"
	"musicAPI/graphql/delivery/graphquery/graph/generated"
	"musicAPI/graphql/delivery/graphquery/graph/model"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

type graphHandler struct {
	usecase graphql.UseCase
}

func GraphHandlers(graphUC graphql.UseCase) {

	graphH := graphHandler{graphUC}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	tracks := graphH.getTracks()
	newTracks := forGraph(tracks)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{newTracks}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (uh *graphHandler) getTracks() []graphql.TrackSelect {
	tracks, err := uh.usecase.GetTracks()
	if err != nil {
		log.Println(err)
	}
	return tracks
}

func forGraph(trackSelect []graphql.TrackSelect) []*model.Tracks {
	var newTracks []*model.Tracks
	for i, _ := range trackSelect {
		newTracks = append(newTracks, &model.Tracks{trackSelect[i].Name, trackSelect[i].Artist, trackSelect[i].Album})
	}
	return newTracks
}
