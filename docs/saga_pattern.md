# Saga Pattern Implementation

This document explains the experimental implementation of the Saga pattern in this project, which ensures data consistency across multiple services (PostgreSQL and OpenSearch).

## Overview

The Saga pattern is a sequence of local transactions where each transaction updates data within a single service. If one transaction fails, the Saga executes compensating transactions to undo the changes made by previous transactions.

In our implementation, we coordinate between:

1. PostgreSQL (primary database)
2. OpenSearch (search engine)

## Implementation Details

### Key Files

1. **Domain Layer**:
   - `internal/domain/repository/opensearch_car.go` - OpenSearch repository interface

2. **Infrastructure Layer**:
   - `internal/infrastructure/opensearch/repository/car_repository.go` - OpenSearch repository implementation
   - `internal/infrastructure/opensearch/client/client.go` - OpenSearch client configuration

3. **Application Layer**:
   - `internal/infrastructure/usecase/car.go` - Car usecase interface (extended with Update method)
   - `internal/infrastructure/usecase/car_impl.go` - Car usecase implementation with saga pattern
   - `internal/infrastructure/usecase/input/update_car.go` - Input struct for Update operation

4. **Tests**:
   - `internal/infrastructure/usecase/car_impl_test.go` - Comprehensive tests for saga pattern implementation

### Saga Flow

The saga pattern is implemented in two key methods in the car usecase:

#### Register Method (Create)

1. Save car data to PostgreSQL
2. Save car data to OpenSearch
3. If OpenSearch operation fails, rollback PostgreSQL insert

#### Update Method

1. Update car data in PostgreSQL
2. Update car data in OpenSearch
3. If OpenSearch operation fails, handle the error appropriately

### Code Examples

Here's how the saga pattern is implemented in the Register method:

```go
// Register registers a new car with saga pattern
func (uc *carUsecase) Register(ctx context.Context, input input.RegisterCarInput) (*model.Car, error) {
    // ... validation code ...

    now := time.Now()
    car := model.NewCar(input.TenantID, input.Model, now)

    // Step 1: Save to PostgreSQL
    if err := uc.carRepo.Create(ctx, car); err != nil {
        return nil, fmt.Errorf("failed to create car in database: %w", err)
    }

    // Step 2: Save to OpenSearch (saga pattern)
    if err := uc.opensearchRepo.Create(ctx, car); err != nil {
        // Compensate: Rollback PostgreSQL insert
        _ = uc.carRepo.Delete(ctx, car.ID) // Best effort rollback
        return nil, fmt.Errorf("failed to create car in opensearch: %w", err)
    }

    return car, nil
}
```

## Testing

The implementation includes comprehensive tests that verify:

- Successful saga execution
- Proper error handling when PostgreSQL operations fail
- Proper error handling when OpenSearch operations fail
- Input validation

Tests can be found in `internal/infrastructure/usecase/car_impl_test.go`, specifically:

- `TestCarUsecase_Register_Success`
- `TestCarUsecase_Register_RepositoryError`
- `TestCarUsecase_Update_Success`
- `TestCarUsecase_Update_GetError`
- `TestCarUsecase_Update_DatabaseError`

## Future Improvements

1. **Enhanced Compensation**: Implement more sophisticated rollback mechanisms
2. **Retry Logic**: Add retry mechanisms for transient failures
3. **Dead Letter Queue**: Implement a mechanism to handle failed compensating transactions
4. **Monitoring**: Add metrics and monitoring for saga execution

## References

- [Saga Pattern - Microsoft Docs](https://docs.microsoft.com/en-us/azure/architecture/reference-architectures/saga/saga)
- [Saga Pattern - Microservices.io](https://microservices.io/patterns/data/saga.html)
