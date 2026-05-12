package usecase

import (
	"context"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure"
	"time"
)

type NotificationServiceConfig struct {
	WAEndpoint string
	SMTP       SMTPConfig
}

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromName  string
	FromEmail string
}

type NotificationService interface {
	QueueNotification(ctx context.Context, poID string) error
	ProcessQueue(ctx context.Context) error
	MarkAsSent(ctx context.Context, poID string) error
	MarkAsFailed(ctx context.Context, poID string, reason string) error
}

type notificationServiceImpl struct {
	poRepo            repository.PurchaseOrderRepository
	notificationQueue repository.NotificationQueueRepository
	waClient          WhatsAppClient
	emailClient      EmailClient
	pdfGenerator     infrastructure.PDFGenerator
}

type WhatsAppClient interface {
	SendWhatsApp(ctx context.Context, phone, message string) error
}

type EmailClient interface {
	SendEmailWithAttachment(to, subject, body string, attachment []byte, filename string) error
}

func NewNotificationService(
	cfg NotificationServiceConfig,
	poRepo repository.PurchaseOrderRepository,
	queueRepo repository.NotificationQueueRepository,
) NotificationService {
	return &notificationServiceImpl{
		poRepo:             poRepo,
		notificationQueue: queueRepo,
		waClient:          infrastructure.NewWhatsAppClient(cfg.WAEndpoint),
		emailClient: infrastructure.NewEmailClient(
			cfg.SMTP.Host,
			cfg.SMTP.Port,
			cfg.SMTP.Username,
			cfg.SMTP.Password,
			cfg.SMTP.FromEmail,
			cfg.SMTP.FromName,
		),
		pdfGenerator: infrastructure.NewPurchaseOrderPDFGenerator(),
	}
}

func (s *notificationServiceImpl) QueueNotification(ctx context.Context, poID string) error {
	return s.notificationQueue.Push(ctx, poID)
}

func (s *notificationServiceImpl) ProcessQueue(ctx context.Context) error {
	for {
		poID, err := s.notificationQueue.Pop(ctx)
		if err != nil || poID == "" {
			break
		}

		if err := s.processOne(ctx, poID); err != nil {
			s.MarkAsFailed(ctx, poID, err.Error())
		}
	}
	return nil
}

func (s *notificationServiceImpl) processOne(ctx context.Context, poID string) error {
	po, err := s.poRepo.FindByIDWithSupplier(ctx, poID)
	if err != nil || po == nil {
		return err
	}

	if po.NotificationStatus != entity.NotificationStatusPending {
		return nil
	}

	var pdfBytes []byte
	if po.NotificationMethod == "EMAIL" {
		pdfBytes, err = s.pdfGenerator.GeneratePO(po)
		if err != nil {
			return err
		}
	}

	message := s.buildMessage(po.PONumber)

	if po.NotificationMethod == "WHATSAPP" {
		if po.Supplier.PhoneNumber != nil {
			err = s.waClient.SendWhatsApp(ctx, *po.Supplier.PhoneNumber, message)
		}
	} else {
		err = s.emailClient.SendEmailWithAttachment(
			*po.Supplier.Email,
			"Purchase Order "+po.PONumber,
			message,
			pdfBytes,
			po.PONumber+".pdf",
		)
	}

	if err != nil {
		return err
	}

	return s.MarkAsSent(ctx, poID)
}

func (s *notificationServiceImpl) buildMessage(poNumber string) string {
	return "Yth. Supplier,\n\nPurchase Order No. " + poNumber + " telah disetujui.\nMohon persiapan pengiriman barang.\n\nTerima kasih."
}

func (s *notificationServiceImpl) MarkAsSent(ctx context.Context, poID string) error {
	now := time.Now()
	po, err := s.poRepo.FindByID(ctx, poID)
	if err != nil || po == nil {
		return err
	}

	po.NotificationStatus = entity.NotificationStatusSent
	po.SentAt = &now

	return s.poRepo.Update(ctx, po)
}

func (s *notificationServiceImpl) MarkAsFailed(ctx context.Context, poID string, reason string) error {
	po, err := s.poRepo.FindByID(ctx, poID)
	if err != nil || po == nil {
		return err
	}

	po.NotificationStatus = entity.NotificationStatusFailed

	return s.poRepo.Update(ctx, po)
}