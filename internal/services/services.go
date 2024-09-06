package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"

	"wb-tech/internal/models"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type repo interface {
	GetOrders(ctx context.Context) (map[string]models.Order, error)
	AddOrder(ctx context.Context, order models.Order) error
	AddDelivery(ctx context.Context, uuid string, delivery models.Delivery) error
	AddPayment(ctx context.Context, uuid string, payment models.Payment) error
	AddOrderItems(ctx context.Context, uuid string, items []models.Items) error
}

type Service struct {
	stan  stan.Conn
	repo  repo
	cache sync.Map
	db    *sqlx.DB
}

func New(stan stan.Conn, repo repo, db *sqlx.DB) *Service {
	return &Service{
		stan:  stan,
		repo:  repo,
		cache: sync.Map{},
		db:    db,
	}
}

func (s *Service) GetCache(uuid string) (models.Order, error) {
	orders, ok := s.cache.Load(uuid)
	if !ok {
		return models.Order{}, fmt.Errorf("%w", ErrOrderNotFound)
	}

	return orders.(models.Order), nil
}

func (s *Service) LoadCache(ctx context.Context) error {
	cache, err := s.repo.GetOrders(ctx)
	if err != nil {
		return err
	}

	for uuid, order := range cache {
		s.cache.Store(uuid, order)
	}

	return nil
}

func (s *Service) AddFromChannel(ctx context.Context, nameChannel string) error {
	var err error
	trManager := manager.Must(trmsqlx.NewDefaultFactory(s.db))

	sub, err := s.stan.Subscribe(nameChannel, func(msg *stan.Msg) {
		order := models.Order{}
		if err = json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("failed to unmarshal order: %v", err)
			return
		}

		for i, orderItem := range order.Items {
			order.Items[i].OrderUid = orderItem.OrderUid
		}

		err = trManager.Do(ctx, func(ctxTx context.Context) error {
			err = s.repo.AddOrder(ctxTx, order)
			if err != nil {
				log.Printf("failed to add order: %v", err)
				return err
			}

			if err = s.repo.AddDelivery(ctxTx, order.OrderUid, order.Delivery); err != nil {
				log.Printf("failed to add delivery: %v", err)
				return err
			}

			if err = s.repo.AddPayment(ctxTx, order.OrderUid, order.Payment); err != nil {
				log.Printf("failed to add payment: %v", err)
				return err
			}

			if err = s.repo.AddOrderItems(ctxTx, order.OrderUid, order.Items); err != nil {
				log.Printf("failed to add order items: %v", err)
				return err
			}

			s.cache.Store(order.OrderUid, order)
			return nil
		})
		if err != nil {
			log.Printf("transaction failed: %v", err)
		}
	})
	if err != nil {
		return err
	}

	defer func() {
		if err := sub.Unsubscribe(); err != nil {
			log.Printf("failed to unsubscribe: %v", err)
		}

		err = s.stan.Close()
		if err != nil {
			log.Printf("failed to close stan: %v", err)
		}
	}()

	<-ctx.Done()

	return nil
}
