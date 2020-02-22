package elasticsearch_test

import (
	"context"
	"testing"

	testing2 "elastic-search/pkg/testing"
	"elastic-search/pkg/todo"
	"elastic-search/platform/elasticsearch"

	"github.com/rs/xid"
	"github.com/tj/assert"
)

const localAddr = "http://localhost:9200"

func TestService_Index(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	sess := testing2.LocalSession()

	t.Run("success", func(t *testing.T) {

		service, err := elasticsearch.NewService(sess, localAddr)
		if err != nil {
			t.Fatalf(err.Error())
		}

		td := todo.Todo{
			ID:      xid.New().String(),
			Content: "content",
		}

		err = service.Index(ctx, td)
		assert.NoError(t, err)
	})
}

func TestService_Search(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	sess := testing2.LocalSession()

	t.Run("many matches", func(t *testing.T) {
		service, err := elasticsearch.NewService(sess, localAddr)
		if err != nil {
			t.Fatalf(err.Error())
		}

		tdin1 := todo.Todo{
			ID:      xid.New().String(),
			Content: "ct1",
		}

		tdin2 := todo.Todo{
			ID:      xid.New().String(),
			Content: "ct2",
		}

		err = service.Index(ctx, tdin1)
		if err != nil {
			t.Fatalf(err.Error())
		}

		err = service.Index(ctx, tdin2)
		if err != nil {
			t.Fatalf(err.Error())
		}

		found, err := service.Search(ctx, "ct")

		assert.NoError(t, err)
		assert.Len(t, found, 2)
	})
}

func TestService_Remove(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	sess := testing2.LocalSession()

	t.Run("success", func(t *testing.T) {
		service, err := elasticsearch.NewService(sess, localAddr)
		if err != nil {
			t.Fatalf(err.Error(), nil)
		}

		tdID := xid.New().String()
		tdIn := todo.Todo{
			Content: "some content",
			ID:      tdID,
		}

		err = service.Index(ctx, tdIn)
		if err != nil {
			t.Fatalf(err.Error(), nil)
		}

		err = service.Remove(ctx, tdID)
		assert.NoError(t, err)
	})

	t.Run("with non-existing item", func(t *testing.T) {
		service, err := elasticsearch.NewService(sess, localAddr)
		if err != nil {
			t.Fatalf(err.Error(), nil)
		}

		err = service.Remove(ctx, xid.New().String())
		assert.Error(t, err, nil)
	})
}
