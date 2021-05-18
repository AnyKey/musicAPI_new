package graphquery

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"musicAPI/graphql"
	"musicAPI/graphql/delivery/graphquery/graph"
	"musicAPI/graphql/delivery/graphquery/graph/generated"
	"musicAPI/graphql/delivery/graphquery/graph/model"
	"net/http"
)

const defaultPort = "8080"

type graphHandler struct {
	usecase graphql.UseCase
}

func GraphHandlers(graphUC graphql.UseCase) {

	graphH := graphHandler{graphUC}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Tracks: graphH.getTracks()}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}

func (gh *graphHandler) getTracks() []*model.Tracks {
	tracks, err := gh.usecase.GetTracks()
	if err != nil {
		log.Println(err)
	}
	return forGraph(tracks)
}

func forGraph(trackSelect []graphql.TrackSelect) []*model.Tracks {
	var newTracks []*model.Tracks
	for i, _ := range trackSelect {
		newTracks = append(newTracks, &model.Tracks{
			Name:   trackSelect[i].Name,
			Artist: trackSelect[i].Artist,
			Album:  trackSelect[i].Album,
		})
	}
	return newTracks
}
