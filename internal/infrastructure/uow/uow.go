package uow

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultMaxRetries    = 3
	defaultBaseDelay     = 100 * time.Millisecond
	mysqlDeadlockErrCode = 1213
)

type ctxKey struct{}

// UnitOfWork manages database transactions via context propagation.
type UnitOfWork interface {
	// Begin starts a new transaction and stores it in the returned context.
	Begin(ctx context.Context) (context.Context, error)
	// Commit commits the transaction stored in the context.
	Commit(ctx context.Context) error
	// Rollback rolls back the transaction stored in the context.
	// Safe to call multiple times; no-op after commit.
	Rollback(ctx context.Context) error
	// Do executes fn inside a transaction with automatic retry on deadlock.
	// The entire fn is re-executed on each retry with a fresh transaction.
	Do(ctx context.Context, fn func(txCtx context.Context) error) error
}

type unitOfWork struct {
	db         *gorm.DB
	maxRetries int
	baseDelay  time.Duration
}

// Option configures the UnitOfWork.
type Option func(*unitOfWork)

// WithMaxRetries sets the maximum number of retry attempts on deadlock.
func WithMaxRetries(n int) Option {
	return func(u *unitOfWork) {
		if n > 0 {
			u.maxRetries = n
		}
	}
}

// WithBaseDelay sets the base delay for exponential backoff between retries.
func WithBaseDelay(d time.Duration) Option {
	return func(u *unitOfWork) {
		if d > 0 {
			u.baseDelay = d
		}
	}
}

func NewUnitOfWork(db *gorm.DB, opts ...Option) UnitOfWork {
	u := &unitOfWork{
		db:         db,
		maxRetries: defaultMaxRetries,
		baseDelay:  defaultBaseDelay,
	}
	for _, opt := range opts {
		opt(u)
	}
	return u
}

func (u *unitOfWork) Begin(ctx context.Context) (context.Context, error) {
	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return ctx, tx.Error
	}
	return context.WithValue(ctx, ctxKey{}, tx), nil
}

func (u *unitOfWork) Commit(ctx context.Context) error {
	tx, ok := ctx.Value(ctxKey{}).(*gorm.DB)
	if !ok {
		return errors.New("no transaction found in context")
	}
	return tx.Commit().Error
}

func (u *unitOfWork) Rollback(ctx context.Context) error {
	tx, ok := ctx.Value(ctxKey{}).(*gorm.DB)
	if !ok {
		return nil // no-op if no transaction
	}
	return tx.Rollback().Error
}

// Do executes fn inside a transaction. If a MySQL deadlock (error 1213) is
// detected, the entire operation is retried up to maxRetries times with
// exponential backoff. The context is checked between retries so callers
// can cancel via context if needed.
func (u *unitOfWork) Do(ctx context.Context, fn func(txCtx context.Context) error) error {
	var lastErr error

	for attempt := 0; attempt < u.maxRetries; attempt++ {
		lastErr = u.doOnce(ctx, fn)
		if lastErr == nil {
			return nil
		}

		if !isDeadlock(lastErr) {
			return lastErr
		}

		// Exponential backoff: baseDelay * 2^attempt
		delay := u.baseDelay * (1 << attempt)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// retry
		}
	}

	return lastErr
}

// doOnce runs fn in a single transaction attempt.
func (u *unitOfWork) doOnce(ctx context.Context, fn func(txCtx context.Context) error) error {
	txCtx, err := u.Begin(ctx)
	if err != nil {
		return err
	}
	defer u.Rollback(txCtx) //nolint:errcheck

	if err := fn(txCtx); err != nil {
		return err
	}

	return u.Commit(txCtx)
}

// isDeadlock checks whether the error is a MySQL deadlock (error code 1213).
func isDeadlock(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == mysqlDeadlockErrCode
	}
	return false
}

// GetTx extracts the transaction from context. Falls back to the provided db
// if no transaction is present, allowing repositories to work both inside and
// outside a Unit of Work.
func GetTx(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(ctxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return db.WithContext(ctx)
}

/*
// Example usage:
// Di service atau handler
err := u.uow.Do(ctx, func(txCtx context.Context) error {
    // semua operasi repository di sini (otomatis pakai tx)
    return u.storeRepo.Create(txCtx, store)
})
*/
