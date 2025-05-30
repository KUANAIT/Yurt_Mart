package handler

import (
	"net/http"
	"time"

	"api-gateway/internal/proxy"
	"github.com/gin-gonic/gin"
	orderPb "github.com/hsibAD/order-service/proto"
	paymentPb "github.com/hsibAD/payment-service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderHandler struct {
	orderClient   *proxy.OrderServiceClient
	paymentClient *proxy.PaymentServiceClient
}

func NewOrderHandler(orderClient *proxy.OrderServiceClient, paymentClient *proxy.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		orderClient:   orderClient,
		paymentClient: paymentClient,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request struct {
		CartID          string                  `json:"cart_id"`
		DeliveryAddress orderPb.DeliveryAddress `json:"delivery_address"`
		DeliveryTime    int64                   `json:"delivery_time"`
		PaymentMethod   string                  `json:"payment_method"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	request.DeliveryAddress.UserId = userID.(string)

	// Create order
	orderReq := &orderPb.CreateOrderRequest{
		CartId:          request.CartID,
		DeliveryAddress: &request.DeliveryAddress,
		DeliveryTime:    timestamppb.New(time.Unix(request.DeliveryTime, 0)),
	}

	order, err := h.orderClient.CreateOrder(c.Request.Context(), orderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If order is created successfully, initiate payment
	paymentReq := &paymentPb.InitiatePaymentRequest{
		OrderId:       order.Id,
		UserId:        userID.(string),
		Amount:        order.TotalPrice,
		Currency:      order.Currency,
		PaymentMethod: paymentPb.PaymentMethod(paymentPb.PaymentMethod_value[request.PaymentMethod]),
	}

	payment, err := h.paymentClient.InitiatePayment(c.Request.Context(), paymentReq)
	if err != nil {
		// If payment initiation fails, we should still return the order
		c.JSON(http.StatusCreated, gin.H{
			"order": order,
			"error": "Payment initiation failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"order":   order,
		"payment": payment,
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	req := &orderPb.GetOrderRequest{
		OrderId: orderID,
	}

	order, err := h.orderClient.GetOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := &orderPb.UpdateOrderStatusRequest{
		OrderId: orderID,
		Status:  request.Status,
	}

	order, err := h.orderClient.UpdateOrderStatus(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) AddDeliveryAddress(c *gin.Context) {
	var address orderPb.DeliveryAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	address.UserId = userID.(string)

	result, err := h.orderClient.AddDeliveryAddress(c.Request.Context(), &address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *OrderHandler) ListDeliveryAddresses(c *gin.Context) {
	userID, _ := c.Get("user_id")

	req := &orderPb.ListAddressesRequest{
		UserId: userID.(string),
	}

	result, err := h.orderClient.ListDeliveryAddresses(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *OrderHandler) GetAvailableDeliverySlots(c *gin.Context) {
	postalCode := c.Query("postal_code")
	date := c.Query("date")

	// Parse date string to time.Time
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	req := &orderPb.DeliverySlotsRequest{
		PostalCode: postalCode,
		Date:       timestamppb.New(t),
	}

	result, err := h.orderClient.GetAvailableDeliverySlots(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
