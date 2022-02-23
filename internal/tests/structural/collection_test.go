package structural

import (
	"context"

	conflictresolution "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/conflictResolution"
	"github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection"
	colletionSubCollection "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection/subcollection"
	colletionSubSimple "github.com/PapaCharlie/go-restli/internal/tests/testdata/generated/testsuite/collection/subsimple"
	"github.com/PapaCharlie/go-restli/protocol"
)

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.collection
var _ = collection.Client(new(collectionClient))

type collectionClient int

func (c *collectionClient) Create(*conflictresolution.Message) (*protocol.CreatedEntity[int64], error) {
	panic(nil)
}

func (c *collectionClient) CreateWithContext(context.Context, *conflictresolution.Message) (*protocol.CreatedEntity[int64], error) {
	panic(nil)
}

func (c *collectionClient) Get(int64) (*conflictresolution.Message, error) {
	panic(nil)
}

func (c *collectionClient) GetWithContext(context.Context, int64) (*conflictresolution.Message, error) {
	panic(nil)
}

func (c *collectionClient) Update(int64, *conflictresolution.Message) error {
	panic(nil)
}

func (c *collectionClient) UpdateWithContext(context.Context, int64, *conflictresolution.Message) error {
	panic(nil)
}

func (c *collectionClient) PartialUpdate(int64, *conflictresolution.Message_PartialUpdate) error {
	panic(nil)
}

func (c *collectionClient) PartialUpdateWithContext(context.Context, int64, *conflictresolution.Message_PartialUpdate) error {
	panic(nil)
}

func (c *collectionClient) Delete(int64) error {
	panic(nil)
}

func (c *collectionClient) DeleteWithContext(context.Context, int64) error {
	panic(nil)
}

func (c *collectionClient) BatchCreate([]*conflictresolution.Message) ([]*protocol.CreatedEntity[int64], error) {
	panic(nil)
}

func (c *collectionClient) BatchCreateWithContext(context.Context, []*conflictresolution.Message) ([]*protocol.CreatedEntity[int64], error) {
	panic(nil)
}

func (c *collectionClient) BatchDelete([]int64) (map[int64]*protocol.BatchEntityUpdateResponse, error) {
	panic(nil)
}

func (c *collectionClient) BatchDeleteWithContext(context.Context, []int64) (map[int64]*protocol.BatchEntityUpdateResponse, error) {
	panic(nil)
}

func (c *collectionClient) BatchGet([]int64) (map[int64]*conflictresolution.Message, error) {
	panic(nil)
}

func (c *collectionClient) BatchGetWithContext(context.Context, []int64) (map[int64]*conflictresolution.Message, error) {
	panic(nil)
}

func (c *collectionClient) BatchUpdate(map[int64]*conflictresolution.Message) (map[int64]*protocol.BatchEntityUpdateResponse, error) {
	panic(nil)
}

func (c *collectionClient) BatchUpdateWithContext(context.Context, map[int64]*conflictresolution.Message) (map[int64]*protocol.BatchEntityUpdateResponse, error) {
	panic(nil)
}

func (c *collectionClient) BatchPartialUpdate(map[int64]*conflictresolution.Message_PartialUpdate) (map[int64]*protocol.BatchEntityUpdateResponse, error) {
	panic(nil)
}

func (c *collectionClient) BatchPartialUpdateWithContext(context.Context, map[int64]*conflictresolution.Message_PartialUpdate) (map[int64]*protocol.BatchEntityUpdateResponse, error) {
	panic(nil)
}

func (c *collectionClient) FindBySearch(*collection.FindBySearchParams) (*protocol.FinderResults[*conflictresolution.Message], error) {
	panic(nil)
}

func (c *collectionClient) FindBySearchWithContext(context.Context, *collection.FindBySearchParams) (*protocol.FinderResults[*conflictresolution.Message], error) {
	panic(nil)
}

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.collection.subCollection
var _ = colletionSubCollection.Client(new(colletionSubCollectionClient))

type colletionSubCollectionClient int

func (s *colletionSubCollectionClient) Get(int64, int64) (*conflictresolution.Message, error) {
	panic(nil)
}

func (s *colletionSubCollectionClient) GetWithContext(context.Context, int64, int64) (*conflictresolution.Message, error) {
	panic(nil)
}

// STRUCTURAL TEST ENSURE INTERFACE DOES NOT DRIFT FOR testsuite.collection.subSimple
var _ = colletionSubSimple.Client(new(colletionSubSimpleClient))

type colletionSubSimpleClient int

func (c *colletionSubSimpleClient) Get(int64) (*conflictresolution.Message, error) {
	panic(nil)
}

func (c *colletionSubSimpleClient) GetWithContext(context.Context, int64) (*conflictresolution.Message, error) {
	panic(nil)
}
