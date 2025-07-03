package dns

import "fmt"

const (
	OperationCreateRecord = "create record"
	OperationUpdateRecord = "update record"
	OperationListRecords  = "list records"
)

type DNSProviderError struct {
	Provider  string
	Operation string
	Err       error
}

func (e *DNSProviderError) Error() string {
	return fmt.Sprintf("%s provider error during %s: %v", e.Provider, e.Operation, e.Err)
}

func (e *DNSProviderError) Unwrap() error {
	return e.Err
}
