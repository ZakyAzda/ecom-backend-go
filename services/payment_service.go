package services

import (
	"ecom-backend-go/config"
	"ecom-backend-go/models"
	"ecom-backend-go/repositories"
	"errors"
	"fmt"
	"time"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService struct {
	Repo repositories.PaymentRepository
}

// SnapTokenResult - data yang dikembalikan setelah snap token berhasil dibuat
type SnapTokenResult struct {
	SnapToken   string `json:"snap_token"`
	RedirectURL string `json:"redirect_url"`
	OrderID     uint   `json:"order_id"`
}

// WebhookPayload - data yang diparsing dari notifikasi Midtrans
type WebhookPayload struct {
	OrderID           string `json:"order_id"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
}

// CreateSnapToken - Buat Midtrans Snap Token untuk order tertentu
func (s *PaymentService) CreateSnapToken(userID uint, orderID uint) (SnapTokenResult, error) {
    order, err := s.Repo.GetOrderByIDAndUser(orderID, userID)
    if err != nil {
        return SnapTokenResult{}, err
    }

    if order.Status == "SELESAI" || order.Status == "PENGIRIMAN" {
        return SnapTokenResult{}, errors.New("order sudah dibayar atau sedang dikirim")
    }

    // ✅ Reuse token jika masih valid (belum expired)
    if order.SnapToken != "" && order.SnapTokenExpiredAt != nil {
        if time.Now().Before(*order.SnapTokenExpiredAt) {
            // Token masih valid → langsung return, tidak perlu ke Midtrans
            return SnapTokenResult{
                SnapToken: order.SnapToken,
                OrderID:   order.ID,
            }, nil
        }
        // Token expired → lanjut buat baru di bawah
    }

    user, err := s.Repo.GetUserByID(userID)
    if err != nil {
        return SnapTokenResult{}, errors.New("gagal mengambil data user")
    }

    itemDetails := s.buildItemDetails(order.OrderItems)

    // ✅ Pakai order.ID + expiry timestamp sebagai suffix unik
    expiredAt := time.Now().Add(23 * time.Hour) // sedikit lebih pendek dari 24 jam Midtrans

    snapReq := &snap.Request{
        TransactionDetails: midtrans.TransactionDetails{
            OrderID:  fmt.Sprintf("ORDER-%d-%d", order.ID, expiredAt.Unix()),
            GrossAmt: int64(order.TotalAmount),
        },
        CustomerDetail: &midtrans.CustomerDetails{
            FName: user.Name,
            Email: user.Email,
            Phone: user.WhatsappNumber,
        },
        Items: &itemDetails,
    }

    snapResp, midErr := config.SnapClient.CreateTransaction(snapReq)
    if midErr != nil {
        return SnapTokenResult{}, fmt.Errorf("gagal membuat sesi pembayaran: %s", midErr.GetMessage())
    }

    // ✅ Simpan token + waktu expired
    if err := s.Repo.SaveSnapToken(order.ID, snapResp.Token, expiredAt); err != nil {
        fmt.Printf("Warning: gagal menyimpan snap_token untuk order #%d: %v\n", order.ID, err)
    }

    return SnapTokenResult{
        SnapToken:   snapResp.Token,
        RedirectURL: snapResp.RedirectURL,
        OrderID:     order.ID,
    }, nil
}

// HandleWebhook - Proses notifikasi dari Midtrans dan update status order
func (s *PaymentService) HandleWebhook(payload WebhookPayload) error {
	// 1. Parse internal order ID dari format "ORDER-{id}-{timestamp}"
	internalOrderID := s.parseOrderID(payload.OrderID)
	if internalOrderID == 0 {
		return fmt.Errorf("order_id tidak bisa diparsing: %s", payload.OrderID)
	}

	// 2. Pastikan order ada di DB
	_, err := s.Repo.GetOrderByID(internalOrderID)
	if err != nil {
		return err
	}

	// 3. Tentukan status internal berdasarkan status Midtrans
	newStatus := s.resolveOrderStatus(payload.TransactionStatus, payload.FraudStatus)

	// 4. Update status order
	if err := s.Repo.UpdateOrderStatusAndPayment(internalOrderID, newStatus, payload.PaymentType); err != nil {
		return fmt.Errorf("gagal update status order #%d: %w", internalOrderID, err)
	}

	fmt.Printf("Order #%d diupdate → %s (tx_status: %s, payment: %s)\n",
		internalOrderID, newStatus, payload.TransactionStatus, payload.PaymentType)

	return nil
}

// GetPaymentStatus - Ambil status pembayaran dan snap_token order
func (s *PaymentService) GetPaymentStatus(userID uint, orderID uint) (models.Order, error) {
	return s.Repo.GetOrderByIDAndUser(orderID, userID)
}

// ─── Private helpers ──────────────────────────────────────────────────────────

func (s *PaymentService) buildItemDetails(items []models.OrderItem) []midtrans.ItemDetails {
	var details []midtrans.ItemDetails
	for _, item := range items {
		details = append(details, midtrans.ItemDetails{
			ID:    fmt.Sprintf("PROD-%d", item.ProductID),
			Name:  item.Product.Name,
			Price: int64(item.Price),
			Qty:   int32(item.Quantity),
		})
	}
	return details
}

func (s *PaymentService) parseOrderID(raw string) uint {
	var id uint
	var timestamp int64
	
	// Gunakan format yang sama persis dengan saat CreateSnapToken
	// Format: "ORDER-{id}-{timestamp}"
	_, err := fmt.Sscanf(raw, "ORDER-%d-%d", &id, &timestamp)
	if err == nil && id > 0 {
		return id
	}

	// Fallback: coba parse langsung sebagai angka
	var directID uint
	fmt.Sscanf(raw, "%d", &directID)
	return directID
}

// resolveOrderStatus - Mapping status Midtrans → status internal
// Referensi: https://docs.midtrans.com/reference/transaction-status-callback
func (s *PaymentService) resolveOrderStatus(transactionStatus, fraudStatus string) string {
	switch transactionStatus {
	case "capture":
		// Hanya untuk kartu kredit — perlu cek fraud status
		if fraudStatus == "accept" {
			return "PENGIRIMAN"
		}
		// fraudStatus == "challenge" → tunggu review manual
		return "BELUM_BAYAR"

	case "settlement":
		// Pembayaran dikonfirmasi (transfer bank, QRIS, e-wallet, dll)
		return "PENGIRIMAN"

	case "pending":
		// Menunggu pembayaran dari user
		return "BELUM_BAYAR"

	case "deny", "cancel", "expire", "refund":
		return "DIBATALKAN"

	default:
		return "BELUM_BAYAR"
	}
}

func (s *PaymentService) UpdateStatusAfterPayment(userID uint, orderID uint, status string) error {
	// Validasi order milik user
	order, err := s.Repo.GetOrderByIDAndUser(orderID, userID)
	if err != nil {
		return err
	}

	// Jangan update jika sudah selesai
	if order.Status == "SELESAI" {
		return nil
	}

	return s.Repo.UpdateOrderStatusAndPayment(orderID, status, order.PaymentMethod)
}