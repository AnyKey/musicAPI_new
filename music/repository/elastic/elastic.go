package elastic

import (
	"github.com/elastic/go-elasticsearch/v7"
	"musicAPI/model"
)

type Repository struct {
	Es *elasticsearch.Client
}

func New(es *elasticsearch.Client) *Repository {
	return &Repository{
		Es: es,
	}
}

func (repo Repository) ElasticAdd(tracks []model.TrackSelect) error {
	return nil
}
func (repo Repository) ElasticGet(tracks []model.TrackSelect) bool {
	return false
}
