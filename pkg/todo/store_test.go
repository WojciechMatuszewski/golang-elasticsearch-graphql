package todo_test

import (
	"testing"

	"elastic-search/pkg/env"
	testing2 "elastic-search/pkg/testing"
	"elastic-search/pkg/todo"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		prefix := xid.New().String()
		todoID := xid.New().String()
		db := testing2.SetupDynamoTest(t, prefix)
		store := todo.NewStore(prefix+env.TODO_TABLE, db)

		inTodo := todo.Todo{ID: todoID, Content: "Content"}
		err := store.Save(inTodo)
		if err != nil {
			t.Fatalf(err.Error())
		}

		outTodo, err := store.GetByID(inTodo.ID)

		assert.NoError(t, err)
		assert.Equal(t, inTodo, outTodo)
	})

}
