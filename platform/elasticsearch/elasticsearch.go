// Package elasticsearch implements methods for interacting with elastic search service
package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"

	"elastic-search/pkg/todo"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olivere/elastic/v7"
	elasticAws "github.com/olivere/elastic/v7/aws/v4"
)

const indexKey = "todo"

// Service represents elasticsearch service
type Service struct {
	client *elastic.Client
}

// NewService returns elasticsearch service
func NewService(sess *session.Session, url string) (Service, error) {
	signingClient := elasticAws.NewV4SigningClient(sess.Config.Credentials, "eu-central-1")

	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(signingClient),
	)
	if err != nil {
		log.WithError(err).Error("could not connect to es")
		return Service{}, err
	}

	log.Info("connected to es service!")

	return Service{client: client}, nil
}

// Index indexes given todo in elasticsearch service
func (s Service) Index(ctx context.Context, td todo.Todo) error {
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

// Search does something
func (s Service) Search(ctx context.Context, query string) ([]todo.Todo, error) {
	var found []todo.Todo
	fmt.Println(query)

	res, err := s.client.Search().Index(indexKey).Query(elastic.NewMatchQuery("content", query).Fuzziness("2")).Do(ctx)
	if err != nil {
		return found, err
	}

	fmt.Println(res.Hits)

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

// Remove removes given item from elastic search
func (s Service) Remove(ctx context.Context, ID string) error {
	_, err := s.client.Delete().Index(indexKey).Id(ID).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
