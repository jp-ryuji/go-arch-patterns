# Entity Relationship Diagram

This system represents a SaaS platform for car rental companies.

Key characteristics of this system:

- **SaaS Platform**: Multi-tenant architecture where each tenant is a separate car rental company
- **Polymorphic Association**: Renter can be either a Company or an Individual (demonstrated through the Renter entity)
- **Many-to-Many Association**: Rental and Option entities are connected through the RentalOption entity, with composite indexing applied for data consistency

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
        string ID
        string Code
        time CreatedAt
        time UpdatedAt
    }

    Car {
        string ID
        string TenantID
        string Model
        time CreatedAt
        time UpdatedAt
    }

    Company {
        string ID
        string TenantID
        string Name
        string CompanySize
        time CreatedAt
        time UpdatedAt
    }

    Individual {
        string ID
        string TenantID
        string Email
        string FirstName
        string LastName
        time CreatedAt
        time UpdatedAt
    }

    Renter {
        string ID
        string TenantID
        string ConcreteRenterID
        string ConcreteRenterModel
        time CreatedAt
        time UpdatedAt
    }

    Rental {
        string ID
        string TenantID
        string CarID
        string RenterID
        time StartsAt
        time EndsAt
        time CreatedAt
        time UpdatedAt
    }

    Option {
        string ID
        string TenantID
        string Name
        time CreatedAt
        time UpdatedAt
    }

    RentalOption {
        string ID
        string TenantID
        string RentalID
        string OptionID
        int Count
        time CreatedAt
        time UpdatedAt
    }
```
