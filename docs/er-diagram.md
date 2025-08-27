# Entity Relationship Diagram

This system represents a SaaS platform for car rental companies.

Key characteristics of this system:

- **SaaS Platform**: Multi-tenant architecture where each tenant is a separate car rental company
- **Polymorphic Association**: Renter can be either a Company or an Individual (demonstrated through the Renter entity)
- **Many-to-Many Association**: Rental and Option entities are connected through the RentalOption entity, with a composite unique index applied to RentalID and OptionID to ensure that the same option cannot be attached to a rental more than once

> **Note**: For simplicity, common columns such as `ID`, `CreatedAt`, and `UpdatedAt` have been omitted from the diagram below. Additionally, the explicit associations with the Tenant entity have been removed, though in the actual implementation all entities are associated with a Tenant in a multi-tenant architecture.

```mermaid
erDiagram
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
        string Model
    }

    Company {
        string Name
        string CompanySize
    }

    Individual {
        string Email
        string FirstName
        string LastName
    }

    Renter {
        string RenterEntityID
        string RenterEntityType
    }

    Rental {
        string CarID
        string RenterID
        time StartsAt
        time EndsAt
    }

    Option {
        string Name
    }

    RentalOption {
        string RentalID
        string OptionID
        int Count
    }
```
