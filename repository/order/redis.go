package order

import ("context"
"errors"
  "txn"
)

type RedisRepo struct {
  Client *redis.Client
}

func orderIDKey(id uint64) string {
  return fmt.Sprintf("order:%d", orderID)
}

func (r *RedisRepo) Insert(ctx context.Context, order model.Order) error {
  data, err := json.Marshal(order)
  if err != nil {
    return fmt.Errorf("marshal error: %w", err)
  }
  key := orderIDKey(order.OrderID)

  txn := r.Client.TxPipeline() 

  res := txn.SetNX(ctx, key, string(data), 0)
  if err := res.Err(); err != nil {
    txn.Discard()
    return fmt.Errorf("redis error failed to set: %w", err)
  }

  if err := txn.SAdd(ctx, "orders", key).Err(); err != nil {
    txn.Discard()
    return fmt.Errorf("redis error failed to add to set: %w", err)
  }
  
  if _, err := txn.Exec(ctx); err != nil {
    return fmt.Errorf("redis error failed to exec: %w", err)
  }

  return nil
}

var ErrOrderNotExist = errors.New("order not exist-")

func (r* RedisRepo) FIndByID(ctx context.Context, id uint64) (model.Order, error) {
  key := orderIDKey(id)
  value, err := r.Client.Get(ctx, key).Result()
  if errors.Is(err, redis.Nil) {
    return model.Order{}, ErrOrderNotExist
  }else if err != nil {
    return model.Order{}, fmt.Errorf("redis error: %w", err)
  }

  var order model.Order
  err = json.Unmarshal([]byte(value), &order)
  if err != nil {
    return model.Order{}, fmt.Errorf("unmarshal error: %w", err)
  }
  return order, nil
}


func (r *RedisRepo) DeleteByID(ctx context.Context, id uint64) error {
 key := orderIDKey(id)
  
  txn := r.Client.TxPipeline()

  err := txn.Del(ctx, key).Err()
  if errors.Is(err, redis.Nil) {
    txn.Discard()
    return ErrOrderNotExist
  }else if err != nil {
    txn.Discard()
    return fmt.Errorf("redis error: %w", err)
  }
  
  if err := txn.SRem(ctx, "orders", key).Err(); err != nil {
    txn.Discard()
    return fmt.Errorf("redis error: %w", err)
  }

  return nil
}

func (r *RedisRepo) Update(ctx context.Context, order model.Order) error {
  data, err := json.Marshal(order)
  if err != nil {
    return fmt.Errorf("marshal error: %w", err)
  }
  key := orderIDKey(order.OrderID)

  err := r.Client.Set(ctx, key, string(data), 0).Err()
  if errors.Is(err, redis.Nil) {
    return ErrOrderNotExist
  }else if err != nil {
    return fmt.Errorf("redis error: %w", err)
  }

  return nil
}

type FindALlPage struct{
  Size unit
  Offset uint
}

type FindResult struct {
  Orders []model.Order
  Cursor uint64
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindALlPage) (FindResult, error) {
  res := r.Client.SScan (ctx, "orders", page.Cursor, "*", int64(page.Size))
  keys, cursor, err := res.Result()
  if err != nil {
    return FindResult{}, fmt.Errorf("redis error failed to get order ids: %w", err)
  }

  if len(keys) == 0 {
    return FindResult{
      Orders: []model.Order{},
    }, nil
  }

  xs, err := r.Client.MGet(ctx, keys...).Result()
  if err != nil {
    return FindResult{}, fmt.Errorf("redis error failed to get orders: %w", err)
  }
  
  orders := make([]model.Order, len(xs))
  for i, x := range xs {
    x := x.(string)
    var order model.Order
    err := json.Unmarshal([]byte(x), &order)
    if err != nil {
      return FindResult{}, fmt.Errorf("unmarshal error: %w", err)
    }
    orders[i] = order
  }
  return FindResult{
    Orders: orders,
    Cursor: cursor,
  }, nil

}
