package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type TxManager struct {
	client *mongo.Client
}

func NewTxManager(client *mongo.Client) *TxManager {
	return &TxManager{
		client: client,
	}
}

func (tm *TxManager) RunTx(ctx context.Context, fn func(sessionContext mongo.SessionContext) error) error {
	session, err := tm.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		err := fn(sessionContext)
		return nil, err
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return err
	}

	return nil
}
