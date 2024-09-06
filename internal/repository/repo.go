package repository

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"wb-tech/internal/models"
)

var (
	SelectItemsQuery = `
		json_agg (
			json_build_object(
				'chrt_id', oit.chrt_id,
				'track_number', oit.track_number,
				'price', oit.price,
				'rid', oit.rid,
				'name', oit.name,
				'sale', oit.sale,
				'size', oit.size,
				'total_price', oit.total_price,
				'nm_id', oit.nm_id,
				'brand', oit.brand,
				'status', oit.status
			)
		) AS items
	`
	SelectPaymentQuery = `
		json_build_object(
			'transaction', p.transaction,
			'request_id', p.request_id,
			'currency', p.currency,
			'provider', p.provider,
			'amount', p.amount,
			'payment_dt', p.payment_dt,
			'bank', p.bank,
			'delivery_cost', p.delivery_cost,
			'goods_total', p.goods_total,
			'custom_fee', p.custom_fee
		) AS payment
	`
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetOrders(ctx context.Context) (map[string]models.Order, error) {
	query := `
		SELECT o.*,
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		(
			SELECT ` + SelectItemsQuery + `
			FROM order_items oit
			WHERE oit.order_uid = o.order_uid
		) AS items,
		(
			SELECT ` + SelectPaymentQuery + `
			FROM payment p
			WHERE p.order_uid = o.order_uid
		) AS payment
		FROM orders o
		JOIN delivery d ON o.order_uid = d.order_uid
	`

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	orders := make(map[string]models.Order)

	for rows.Next() {
		var order models.Order
		var d models.Delivery
		var itemsJSON, paymentJSON string

		if err := rows.Scan(
			&order.OrderUid,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerId,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmId,
			&order.DateCreated,
			&order.OofShard,
			&d.Name,
			&d.Phone,
			&d.Zip,
			&d.City,
			&d.Address,
			&d.Region,
			&d.Email,
			&itemsJSON,
			&paymentJSON,
		); err != nil {
			return nil, errors.WithStack(err)
		}

		var items []models.Items
		if err := json.Unmarshal([]byte(itemsJSON), &items); err != nil {
			return nil, errors.WithStack(err)
		}
		order.Items = items

		var payment models.Payment
		if err := json.Unmarshal([]byte(paymentJSON), &payment); err != nil {
			return nil, errors.WithStack(err)
		}
		order.Payment = payment
		order.Delivery = d
		orders[order.OrderUid] = order
	}

	if err := rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	return orders, nil
}

func (r *Repo) AddOrder(ctx context.Context, order models.Order) error {
	query := `
		INSERT INTO orders (order_uid, track_number,
		entry, locale, internal_signature, customer_id, delivery_service,
		shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query, order.OrderUid, order.TrackNumber, order.Entry,
		order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService,
		order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) AddDelivery(ctx context.Context, uuid string, delivery models.Delivery) error {
	query := `
		INSERT INTO delivery (order_uid, name, phone, zip,
		city, address, region, email)
		Values ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query, uuid, delivery.Name, delivery.Phone, delivery.Zip,
		delivery.City, delivery.Address, delivery.Region, delivery.Email)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) AddPayment(ctx context.Context, uuid string, payment models.Payment) error {
	query := `
		INSERT INTO payment (order_uid, transaction, request_id, currency, provider,
		amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.ExecContext(ctx, query, uuid, payment.Transaction, payment.RequestId,
		payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt,
		payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Repo) AddOrderItems(ctx context.Context, uuid string, items []models.Items) error {
	query := `
		INSERT INTO order_items (order_uid, chrt_id, track_number,
		price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	for _, item := range items {
		_, err := r.db.ExecContext(ctx, query, uuid, item.ChrtId, item.TrackNumber,
			item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmId, item.Brand, item.Status)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
