package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	productpb "api-gateway/client-service/proto/productpb"
	cartpb "api-gateway/shopping-cart-service/proto/cartpb"
)

func main() {

	productConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to product-service: %v", err)
	}
	defer productConn.Close()

	productClient := productpb.NewProductServiceClient(productConn)

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to cart service: %v", err)
	}
	defer conn.Close()

	client := cartpb.NewCartServiceClient(conn)

	r := gin.Default()

	r.POST("/cart", func(c *gin.Context) {
		var req struct {
			UserID string `json:"user_id"`
			Items  []struct {
				ProductID string `json:"product_id"`
				Quantity  int32  `json:"quantity"`
			} `json:"items"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var items []*cartpb.CartItem
		for _, item := range req.Items {
			items = append(items, &cartpb.CartItem{
				ProductId: item.ProductID,
				Quantity:  item.Quantity,
			})
		}

		resp, err := client.AddToCart(context.Background(), &cartpb.AddToCartRequest{
			UserId: req.UserID,
			Items:  items,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	r.GET("/cart/:user_id", func(c *gin.Context) {
		userID := c.Param("user_id")

		res, err := client.GetCart(context.Background(), &cartpb.GetCartRequest{
			UserId: userID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.DELETE("/cart/:user_id/:product_id", func(c *gin.Context) {
		userID := c.Param("user_id")
		productID := c.Param("product_id")

		res, err := client.RemoveFromCart(context.Background(), &cartpb.RemoveFromCartRequest{
			UserId:    userID,
			ProductId: productID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/product/:id", func(c *gin.Context) {
		id := c.Param("id")
		resp, err := productClient.GetProduct(context.Background(), &productpb.GetProductRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/product", func(c *gin.Context) {
		var req productpb.Product
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		resp, err := productClient.CreateProduct(context.Background(), &productpb.CreateProductRequest{Product: &req})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.DELETE("/product/:id", func(c *gin.Context) {
		id := c.Param("id")
		resp, err := productClient.DeleteProduct(context.Background(), &productpb.DeleteProductRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/products", func(c *gin.Context) {
		resp, err := productClient.ListProducts(context.Background(), &productpb.ListProductsRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/products/category/:category", func(c *gin.Context) {
		category := c.Param("category")
		resp, err := productClient.GetProductsByCategory(context.Background(), &productpb.GetProductsByCategoryRequest{
			Category: category,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	log.Println("API Gateway started at http://localhost:8080")
	r.Run(":8080")
}
