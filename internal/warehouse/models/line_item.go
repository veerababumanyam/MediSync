package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// LineItem represents a line item from an invoice or a transaction from a bank statement.
type LineItem struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	DocumentID       uuid.UUID  `json:"document_id" db:"document_id"`
	ExtractedFieldID uuid.UUID  `json:"extracted_field_id" db:"extracted_field_id"`
	LineNumber       int        `json:"line_number" db:"line_number"`
	Description      string     `json:"description" db:"description"`
	Quantity         float64    `json:"quantity" db:"quantity"`
	UnitPrice        float64    `json:"unit_price" db:"unit_price"`
	Amount           float64    `json:"amount" db:"amount"`
	TaxRate          float64    `json:"tax_rate" db:"tax_rate"`
	TransactionDate  *time.Time `json:"transaction_date" db:"transaction_date"`
	Reference        string     `json:"reference" db:"reference"`
	DebitAmount      float64    `json:"debit_amount" db:"debit_amount"`
	CreditAmount     float64    `json:"credit_amount" db:"credit_amount"`
	Balance          float64    `json:"balance" db:"balance"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

// Validate checks if the LineItem has valid field values.
func (l *LineItem) Validate() error {
	var errs []error

	if l.ID == uuid.Nil {
		errs = append(errs, errors.New("id is required"))
	}

	if l.DocumentID == uuid.Nil {
		errs = append(errs, errors.New("document_id is required"))
	}

	if l.LineNumber < 1 {
		errs = append(errs, errors.New("line_number must be positive"))
	}

	if l.Quantity < 0 {
		errs = append(errs, errors.New("quantity cannot be negative"))
	}

	if l.UnitPrice < 0 {
		errs = append(errs, errors.New("unit_price cannot be negative"))
	}

	if l.Amount < 0 {
		errs = append(errs, errors.New("amount cannot be negative"))
	}

	if l.TaxRate < 0 || l.TaxRate > 100 {
		errs = append(errs, errors.New("tax_rate must be between 0 and 100"))
	}

	if l.DebitAmount < 0 {
		errs = append(errs, errors.New("debit_amount cannot be negative"))
	}

	if l.CreditAmount < 0 {
		errs = append(errs, errors.New("credit_amount cannot be negative"))
	}

	// Validate that debit and credit are mutually exclusive for transactions
	if l.DebitAmount > 0 && l.CreditAmount > 0 {
		errs = append(errs, errors.New("debit_amount and credit_amount cannot both be positive"))
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %w", errors.Join(errs...))
	}

	return nil
}

// IsInvoiceLineItem returns true if this is an invoice line item.
func (l *LineItem) IsInvoiceLineItem() bool {
	return l.Quantity > 0 || l.UnitPrice > 0 || l.Amount > 0
}

// IsBankTransaction returns true if this is a bank statement transaction.
func (l *LineItem) IsBankTransaction() bool {
	return l.TransactionDate != nil || l.DebitAmount > 0 || l.CreditAmount > 0
}

// IsDebit returns true if this is a debit transaction.
func (l *LineItem) IsDebit() bool {
	return l.DebitAmount > 0
}

// IsCredit returns true if this is a credit transaction.
func (l *LineItem) IsCredit() bool {
	return l.CreditAmount > 0
}

// CalculateTotal calculates the line item total (amount + tax).
func (l *LineItem) CalculateTotal() float64 {
	if l.Amount == 0 {
		return l.Quantity * l.UnitPrice
	}
	taxAmount := l.Amount * (l.TaxRate / 100)
	return l.Amount + taxAmount
}

// NewInvoiceLineItem creates a new invoice line item.
func NewInvoiceLineItem(documentID uuid.UUID, lineNumber int, description string, quantity, unitPrice, amount, taxRate float64) *LineItem {
	return &LineItem{
		ID:         uuid.New(),
		DocumentID: documentID,
		LineNumber: lineNumber,
		Description: description,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		Amount:     amount,
		TaxRate:    taxRate,
		CreatedAt:  time.Now(),
	}
}

// NewBankTransaction creates a new bank statement transaction.
func NewBankTransaction(documentID uuid.UUID, lineNumber int, transactionDate time.Time, description, reference string, debit, credit, balance float64) *LineItem {
	return &LineItem{
		ID:              uuid.New(),
		DocumentID:      documentID,
		LineNumber:      lineNumber,
		TransactionDate: &transactionDate,
		Description:     description,
		Reference:       reference,
		DebitAmount:     debit,
		CreditAmount:    credit,
		Balance:         balance,
		CreatedAt:       time.Now(),
	}
}

// LineItemList represents a list of line items with total calculations.
type LineItemList struct {
	Items      []LineItem `json:"items"`
	Subtotal   float64    `json:"subtotal"`
	TaxAmount  float64    `json:"tax_amount"`
	Total      float64    `json:"total"`
	ItemCount  int        `json:"item_count"`
}

// CalculateTotals calculates the subtotal, tax, and total from the line items.
func (l *LineItemList) CalculateTotals() {
	l.Subtotal = 0
	l.TaxAmount = 0

	for _, item := range l.Items {
		amount := item.Amount
		if amount == 0 {
			amount = item.Quantity * item.UnitPrice
		}
		l.Subtotal += amount
		l.TaxAmount += amount * (item.TaxRate / 100)
	}

	l.Total = l.Subtotal + l.TaxAmount
	l.ItemCount = len(l.Items)
}
