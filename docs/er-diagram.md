# Entity Relationship Diagram

This system represents a SaaS platform for car rental companies.

Key characteristics of this system:

- **SaaS Platform**: Multi-tenant architecture where each tenant is a separate car rental company
- **Polymorphic Association**: Renter can be either a Company or an Individual (demonstrated through the Renter entity)
- **Many-to-Many Association**: Rental and Option entities are connected through the RentalOption entity, with composite indexing applied for data consistency

> **Note**: For simplicity, common columns such as `ID`, `CreatedAt`, and `UpdatedAt` have been omitted from the diagram below.

```mermaid
erDiagram
    Tenant ||--o{ Car : has
    Tenant ||--o{ Company : has
    Tenant ||--o{ Individual : has
    Tenant ||--o{ Renter : has
    Tenant ||--o{ Rental : has
    Tenant ||--o{ Option : has
    Tenant ||--o{ RentalOption : has

    Company ||--o{ Renter : "can be"
    Individual ||--o{ Renter : "can be"
    Car ||--o{ Rental : has

    Renter ||--o{ Rental : has
    Option ||--o{ RentalOption : has
    Rental ||--o{ RentalOption : has

    Tenant {
        string Code
    }

    Car {
        string TenantID
        string Model
    }

    Company {
        string TenantID
        string Name
        string CompanySize
    }

    Individual {
        string TenantID
        string Email
        string FirstName
        string LastName
    }

    Renter {
        string TenantID
        string ConcreteRenterID
        string ConcreteRenterModel
    }

    Rental {
        string TenantID
        string CarID
        string RenterID
        time StartsAt
        time EndsAt
    }

    Option {
        string TenantID
        string Name
    }

    RentalOption {
        string TenantID
        string RentalID
        string OptionID
        int Count
    }
```
