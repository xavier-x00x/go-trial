package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"go-trial/internal/delivery/http/dto"
	"go-trial/internal/domain/entity"
	"go-trial/internal/domain/repository"
	"go-trial/internal/infrastructure/uow"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

var (
	ErrCOANotFound     = errors.New("chart of account not found")
	ErrCOAHasChildren  = errors.New("chart of account has child accounts")
)

type ChartOfAccountUseCase interface {
	Create(ctx context.Context, req dto.CreateCOARequest) (*dto.ChartOfAccountResponse, error)
	GetByID(ctx context.Context, id string) (*dto.ChartOfAccountResponse, error)
	GetAll(ctx context.Context) ([]dto.ChartOfAccountResponse, error)
	GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.ChartOfAccountResponse, *entity.Meta, error)
	GetByType(ctx context.Context, accountType string) ([]dto.ChartOfAccountResponse, error)
	GetTree(ctx context.Context) ([]dto.ChartOfAccountTreeResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateCOARequest) (*dto.ChartOfAccountResponse, error)
	Delete(ctx context.Context, id string) error
	Import(ctx context.Context, file io.Reader, filename string) (*dto.AccountImportResult, error)
	GenerateTemplate(ctx context.Context) ([]byte, error)
}

type chartOfAccountUseCase struct {
	coaRepo repository.ChartOfAccountRepository
	uow     uow.UnitOfWork
}

func NewChartOfAccountUseCase(coaRepo repository.ChartOfAccountRepository, uow uow.UnitOfWork) ChartOfAccountUseCase {
	return &chartOfAccountUseCase{
		coaRepo: coaRepo,
		uow:     uow,
	}
}

func (u *chartOfAccountUseCase) Create(ctx context.Context, req dto.CreateCOARequest) (*dto.ChartOfAccountResponse, error) {
	var fe FieldErrors

	existing, err := u.coaRepo.FindByCode(ctx, req.AccountCode)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		fe.Add("account_code", "kode akun sudah digunakan")
	}

	if req.ParentID != nil && *req.ParentID != "" {
		parent, err := u.coaRepo.FindByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			fe.Add("parent_id", "parent account tidak ditemukan")
		}
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
	}

	coa := &entity.ChartOfAccount{}
	if err := coa.GenerateID(); err != nil {
		return nil, err
	}
	coa.AccountCode = req.AccountCode
	coa.Name = req.Name
	coa.AccountType = req.AccountType
	coa.NormalBalance = req.NormalBalance
	coa.IsActive = true

	if req.ParentID != nil && *req.ParentID != "" {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, err
		}
		coa.ParentID = &pid
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.coaRepo.Create(txCtx, coa); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toCOAResponse(coa)
	return &resp, nil
}

func (u *chartOfAccountUseCase) GetByID(ctx context.Context, id string) (*dto.ChartOfAccountResponse, error) {
	coa, err := u.coaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if coa == nil {
		return nil, ErrCOANotFound
	}

	resp := toCOAResponse(coa)
	return &resp, nil
}

func (u *chartOfAccountUseCase) GetAll(ctx context.Context) ([]dto.ChartOfAccountResponse, error) {
	coas, err := u.coaRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var resp []dto.ChartOfAccountResponse
	for _, c := range coas {
		resp = append(resp, toCOAResponse(&c))
	}
	return resp, nil
}

func (u *chartOfAccountUseCase) GetAllWithPagination(ctx context.Context, meta *dto.MetaRequest) ([]dto.ChartOfAccountResponse, *entity.Meta, error) {
	allowedOrder := []string{"id", "account_code", "name", "account_type", "updated_at"}
	searchColumns := []string{"id", "account_code", "name"}

	filter := BuildQueryFilter(meta, allowedOrder, searchColumns)
	filter.Conditions["deleted_at"] = nil

	data, resMeta, err := u.coaRepo.FindAllWithPagination(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	var resp []dto.ChartOfAccountResponse
	for _, c := range data {
		resp = append(resp, toCOAResponse(&c))
	}
	return resp, resMeta, nil
}

func (u *chartOfAccountUseCase) GetByType(ctx context.Context, accountType string) ([]dto.ChartOfAccountResponse, error) {
	coas, err := u.coaRepo.FindByType(ctx, accountType)
	if err != nil {
		return nil, err
	}

	var resp []dto.ChartOfAccountResponse
	for _, c := range coas {
		resp = append(resp, toCOAResponse(&c))
	}
	return resp, nil
}

func (u *chartOfAccountUseCase) GetTree(ctx context.Context) ([]dto.ChartOfAccountTreeResponse, error) {
	all, err := u.coaRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	byParent := make(map[string][]entity.ChartOfAccount)
	for _, c := range all {
		pid := ""
		if c.ParentID != nil {
			pid = c.ParentID.String()
		}
		byParent[pid] = append(byParent[pid], c)
	}

	var build func(parentID string) []dto.ChartOfAccountTreeResponse
	build = func(parentID string) []dto.ChartOfAccountTreeResponse {
		var res []dto.ChartOfAccountTreeResponse
		for _, c := range byParent[parentID] {
			node := dto.ChartOfAccountTreeResponse{
				ID:            c.ID.String(),
				AccountCode:   c.AccountCode,
				Name:          c.Name,
				AccountType:   c.AccountType,
				NormalBalance: c.NormalBalance,
				IsActive:      c.IsActive,
			}
			node.Children = build(c.ID.String())
			res = append(res, node)
		}
		return res
	}

	return build(""), nil
}

func (u *chartOfAccountUseCase) Update(ctx context.Context, id string, req dto.UpdateCOARequest) (*dto.ChartOfAccountResponse, error) {
	coa, err := u.coaRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if coa == nil {
		return nil, ErrCOANotFound
	}

	var fe FieldErrors

	if req.AccountCode != nil && *req.AccountCode != coa.AccountCode {
		existing, err := u.coaRepo.FindByCode(ctx, *req.AccountCode)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			fe.Add("account_code", "kode akun sudah digunakan")
		} else {
			coa.AccountCode = *req.AccountCode
		}
	}

	if req.Name != nil {
		coa.Name = *req.Name
	}
	if req.AccountType != nil {
		coa.AccountType = *req.AccountType
	}
	if req.NormalBalance != nil {
		coa.NormalBalance = *req.NormalBalance
	}
	if req.IsActive != nil {
		coa.IsActive = *req.IsActive
	}

	if req.ParentID != nil {
		if *req.ParentID == "" {
			coa.ParentID = nil
		} else {
			if *req.ParentID == id {
				fe.Add("parent_id", "parent_id tidak boleh sama dengan ID akun itu sendiri")
			}
			parent, err := u.coaRepo.FindByID(ctx, *req.ParentID)
			if err != nil {
				return nil, err
			}
			if parent == nil {
				fe.Add("parent_id", "parent account tidak ditemukan")
			}
		}
	}

	if len(fe.Errors) > 0 {
		return nil, &fe
	}

	if req.ParentID != nil && *req.ParentID != "" {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, err
		}
		coa.ParentID = &pid
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.coaRepo.Update(txCtx, coa); err != nil {
		return nil, err
	}

	if err := u.uow.Commit(txCtx); err != nil {
		return nil, err
	}

	resp := toCOAResponse(coa)
	return &resp, nil
}

func (u *chartOfAccountUseCase) Import(ctx context.Context, file io.Reader, filename string) (*dto.AccountImportResult, error) {
	var fe FieldErrors

	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".xlsx" && ext != ".xls" {
		fe.Add("file", "file harus berformat .xlsx atau .xls")
		return nil, &fe
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	xlsx, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer xlsx.Close()

	rows, err := xlsx.GetRows("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet: %w", err)
	}

	if len(rows) < 2 {
		return &dto.AccountImportResult{
			TotalRows:   0,
			SuccessRows: 0,
			ErrorRows:   0,
		}, nil
	}

	validAccountTypes := map[string]bool{"ASSET": true, "LIABILITY": true, "EQUITY": true, "REVENUE": true, "EXPENSE": true}
	validNormalBalances := map[string]bool{"DEBIT": true, "CREDIT": true}

	var created []entity.ChartOfAccount
	var rowErrors []dto.ImportRowError
	successCount := 0

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNum := i + 1

		if len(row) < 4 {
			rowErrors = append(rowErrors, dto.ImportRowError{
				RowNumber: rowNum,
				Message:   "row must have at least 4 columns (account_code, name, account_type, normal_balance)",
			})
			continue
		}

		accountCode := strings.TrimSpace(row[0])
		name := strings.TrimSpace(row[1])
		accountType := strings.ToUpper(strings.TrimSpace(row[2]))
		normalBalance := strings.ToUpper(strings.TrimSpace(row[3]))

		if accountCode == "" {
			rowErrors = append(rowErrors, dto.ImportRowError{
				RowNumber: rowNum,
				Message:   "account_code is required",
			})
			continue
		}
		if name == "" {
			rowErrors = append(rowErrors, dto.ImportRowError{
				RowNumber: rowNum,
				Message:   "name is required",
			})
			continue
		}
		if !validAccountTypes[accountType] {
			rowErrors = append(rowErrors, dto.ImportRowError{
				RowNumber: rowNum,
				Message:   fmt.Sprintf("invalid account_type '%s', must be one of: ASSET, LIABILITY, EQUITY, REVENUE, EXPENSE", accountType),
			})
			continue
		}
		if !validNormalBalances[normalBalance] {
			rowErrors = append(rowErrors, dto.ImportRowError{
				RowNumber: rowNum,
				Message:   fmt.Sprintf("invalid normal_balance '%s', must be one of: DEBIT, CREDIT", normalBalance),
			})
			continue
		}

		existing, err := u.coaRepo.FindByCode(ctx, accountCode)
		if err != nil {
			rowErrors = append(rowErrors, dto.ImportRowError{
				RowNumber: rowNum,
				Message:   fmt.Sprintf("error checking duplicate: %s", err.Error()),
			})
			continue
		}
		if existing != nil {
			rowErrors = append(rowErrors, dto.ImportRowError{
				RowNumber: rowNum,
				Message:   fmt.Sprintf("account code '%s' already exists", accountCode),
			})
			continue
		}

		coa := &entity.ChartOfAccount{}
		if err := coa.GenerateID(); err != nil {
			rowErrors = append(rowErrors, dto.ImportRowError{
				RowNumber: rowNum,
				Message:   fmt.Sprintf("failed to generate id: %s", err.Error()),
			})
			continue
		}
		coa.AccountCode = accountCode
		coa.Name = name
		coa.AccountType = accountType
		coa.NormalBalance = normalBalance
		coa.IsActive = true

		if len(row) >= 5 {
			parentCode := strings.TrimSpace(row[4])
			if parentCode != "" {
				parent, err := u.coaRepo.FindByCode(ctx, parentCode)
				if err != nil {
					rowErrors = append(rowErrors, dto.ImportRowError{
						RowNumber: rowNum,
						Message:   fmt.Sprintf("error finding parent '%s': %s", parentCode, err.Error()),
					})
					continue
				}
				if parent == nil {
					rowErrors = append(rowErrors, dto.ImportRowError{
						RowNumber: rowNum,
						Message:   fmt.Sprintf("parent account code '%s' not found", parentCode),
					})
					continue
				}
				coa.ParentID = &parent.ID
			}
		}

		created = append(created, *coa)
		successCount++
	}

	if len(created) > 0 {
		txCtx, err := u.uow.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer u.uow.Rollback(txCtx)

		if err := u.coaRepo.BulkCreate(txCtx, created); err != nil {
			return nil, err
		}

		if err := u.uow.Commit(txCtx); err != nil {
			return nil, err
		}
	}

	return &dto.AccountImportResult{
		TotalRows:   len(rows) - 1,
		SuccessRows: successCount,
		ErrorRows:   len(rowErrors),
		Errors:      rowErrors,
	}, nil
}

func (u *chartOfAccountUseCase) GenerateTemplate(_ context.Context) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	f.SetCellValue("Sheet1", "A1", "account_code")
	f.SetCellValue("Sheet1", "B1", "name")
	f.SetCellValue("Sheet1", "C1", "account_type")
	f.SetCellValue("Sheet1", "D1", "normal_balance")
	f.SetCellValue("Sheet1", "E1", "parent_code")

	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"D9E2F3"}},
	})
	_ = f.SetCellStyle("Sheet1", "A1", "E1", style)

	_ = f.SetColWidth("Sheet1", "A", "A", 20)
	_ = f.SetColWidth("Sheet1", "B", "B", 40)
	_ = f.SetColWidth("Sheet1", "C", "C", 20)
	_ = f.SetColWidth("Sheet1", "D", "D", 20)
	_ = f.SetColWidth("Sheet1", "E", "E", 20)

	dvRange := excelize.NewDataValidation(false)
	dvRange.SetSqref("C2:C1048576")
	dvRange.SetDropList([]string{"ASSET", "LIABILITY", "EQUITY", "REVENUE", "EXPENSE"})
	_ = f.AddDataValidation("Sheet1", dvRange)

	dvRange2 := excelize.NewDataValidation(false)
	dvRange2.SetSqref("D2:D1048576")
	dvRange2.SetDropList([]string{"DEBIT", "CREDIT"})
	_ = f.AddDataValidation("Sheet1", dvRange2)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write excel: %w", err)
	}

	return buf.Bytes(), nil
}

func (u *chartOfAccountUseCase) Delete(ctx context.Context, id string) error {
	coa, err := u.coaRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if coa == nil {
		return ErrCOANotFound
	}

	children, err := u.coaRepo.FindByParentID(ctx, &id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return ErrCOAHasChildren
	}

	txCtx, err := u.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer u.uow.Rollback(txCtx)

	if err := u.coaRepo.Delete(txCtx, id); err != nil {
		return err
	}

	return u.uow.Commit(txCtx)
}

func toCOAResponse(coa *entity.ChartOfAccount) dto.ChartOfAccountResponse {
	var parentID *string
	if coa.ParentID != nil {
		s := coa.ParentID.String()
		parentID = &s
	}
	return dto.ChartOfAccountResponse{
		ID:            coa.ID.String(),
		AccountCode:   coa.AccountCode,
		Name:          coa.Name,
		AccountType:   coa.AccountType,
		NormalBalance: coa.NormalBalance,
		IsActive:      coa.IsActive,
		ParentID:      parentID,
		CreatedAt:     coa.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     coa.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
