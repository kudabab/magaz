package service

import (
	"context"
	"fmt"

	"github.com/kudabab/market-s/db"
	"github.com/kudabab/market-s/entity"
)

func AddToCart(ctx context.Context, userID, ProductID, quantity int, category string) (*entity.Cart, error) {

	product, err := GetProductByID(category, ProductID)

	var existItem entity.CartItem
	err = db.DB.QueryRow(context.Background(),
		"SELECT product_id, quantity, price, category FROM cart_items WHERE user_id=$1 AND product_id=$2",
		userID, ProductID).Scan(
		&existItem.ProductID,
		&existItem.Quantity,
		&existItem.Price,
		&existItem.Category,
	)

	if err == nil {
		_, err = db.DB.Exec(context.Background(),
			"UPDATE cart_items SET quantity = quantity + $1, price = (quantity + $1) * $2 WHERE user_id = $3 AND product_id = $4",
			quantity, product.Price, userID, ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed update cart item: %v", err)

		}

	} else {
		_, err = db.DB.Exec(context.Background(),
			"INSERT INTO cart_items (user_id, product_id, quantity, price, category) VALUES ($1, $2, $3, $4, $5)",
			userID, ProductID, quantity, float64(quantity)*product.Price, category)
		if err != nil {
			return nil, fmt.Errorf("failed add cart item: %v", err)

		}
	}
	cart, err := GetCart(userID)

	return cart, nil
}

func GetCart(userID int) (*entity.Cart, error) {
	rows, err := db.DB.Query(context.Background(),
		"SELECT product_id, category, quantity, price FROM cart_items WHERE user_id = $1",
		userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %v", err)
	}
	defer rows.Close()

	var cart entity.Cart
	cart.UserID = userID
	var total float64

	for rows.Next() {
		var item entity.CartItem
		if err := rows.Scan(&item.ProductID, &item.Category, &item.Quantity, &item.Price); err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %v", err)
		}
		cart.Items = append(cart.Items, item)
		total += item.Price
	}

	cart.Total = total
	return &cart, nil
}

func RemoveFromCart(userID, productID int) (*entity.Cart, error) {
	// Удаляем товар из корзины
	_, err := db.DB.Exec(context.Background(),
		"DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2",
		userID, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove cart item: %v", err)
	}

	// Получаем обновленную корзину
	cart, err := GetCart(userID)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func CreateOrder(userID int) (*entity.Order, error) {
	// Получаем корзину пользователя
	cart, err := GetCart(userID)
	if err != nil {
		return nil, err
	}

	// Создаем заказ
	var orderID int
	err = db.DB.QueryRow(context.Background(),
		"INSERT INTO orders (user_id, total, status) VALUES ($1, $2, 'pending') RETURNING order_id",
		userID, cart.Total).Scan(&orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	// Переносим элементы корзины в заказ
	for _, item := range cart.Items {
		_, err = db.DB.Exec(context.Background(),
			"INSERT INTO order_items (order_id, product_id, category, quantity, price) VALUES ($1, $2, $3, $4, $5)",
			orderID, item.ProductID, item.Category, item.Quantity, item.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to add order item: %v", err)
		}
	}

	// Очистка корзины после заказа
	_, err = db.DB.Exec(context.Background(),
		"DELETE FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear cart: %v", err)
	}

	// Создаем и возвращаем заказ
	order := &entity.Order{
		OrderID: orderID,
		UserID:  userID,
		Items:   cart.Items,
		Total:   cart.Total,
		Status:  "pending",
	}

	return order, nil
}
