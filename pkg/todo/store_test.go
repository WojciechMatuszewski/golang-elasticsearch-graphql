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

func TestStore_Remove(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		prefix := xid.New().String()
		todoID := xid.New().String()
		db := testing2.SetupDynamoTest(t, prefix)
		store := todo.NewStore(prefix+env.TODO_TABLE, db)

		inTodo := todo.Todo{ID: todoID, Content: "Content"}
		err := store.Save(inTodo)
		if err != nil {
			t.Fatalf(err.Error(), nil)
		}

		td, err := store.GetByID(inTodo.ID)
		if err != nil {
			t.Fatalf(err.Error(), nil)
		}
		assert.Equal(t, inTodo, td)

		err = store.Remove(inTodo.ID)
		assert.NoError(t, err)

		td, err = store.GetByID(inTodo.ID)
		if err != nil {
			t.Fatalf(err.Error(), nil)
		}

		assert.Empty(t, td, nil)
	})

	t.Run("removing non existing todo", func(t *testing.T) {
		prefix := xid.New().String()
		db := testing2.SetupDynamoTest(t, prefix)
		store := todo.NewStore(prefix+env.TODO_TABLE, db)

		err := store.Remove(prefix)
		assert.Error(t, err, nil)
	})
}
