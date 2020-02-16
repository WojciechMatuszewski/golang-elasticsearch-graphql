package elasticsearch

import (
	"context"
	"encoding/json"

	"elastic-search/pkg/todo"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olivere/elastic/v7"
)

const indexKey = "todo"

// Service represents elasticsearch service
type Service struct {
	client *elastic.Client
}

// NewService returns elasticsearch service
func NewService(sess *session.Session, url string) (Service, error) {
	//signingClient := elasticAws.NewV4SigningClient(sess.Config.Credentials, "local")

	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return Service{}, err
	}

	return Service{client: client}, nil
}

// Index indexes given todo in elasticsearch service
func (s *Service) Index(ctx context.Context, td todo.Todo) error {
	tdB, err := json.Marshal(td)
	if err != nil {
		return err
	}

	_, err = s.client.Index().Index(indexKey).Id(td.ID).BodyJson(string(tdB)).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Search(ctx context.Context, query string) ([]todo.Todo, error) {
	var found []todo.Todo

	// implement proper query here
	res, err := s.client.Search().Index(indexKey).Query(elastic.NewTermQuery("content", query)).Do(ctx)
	if err != nil {
		return found, err
	}

	for _, hit := range res.Hits.Hits {
		var td todo.Todo

		err = json.Unmarshal(hit.Source, &td)
		if err != nil {
			return found, err
		}

		found = append(found, td)
	}

	return found, nil
}
