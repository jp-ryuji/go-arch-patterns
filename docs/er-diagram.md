# ER Diagram

This system represents a SaaS platform for car rental companies.

Key characteristics of this system:

- **SaaS Platform**: Multi-tenant architecture where each tenant is a separate car rental company
- **Class Table Inheritance**: Renter is implemented using Class Table Inheritance pattern where Company and Individual are specialized types of Renter
- **Many-to-Many Association**: Rental and Option entities are connected through the RentalOption entity, with a composite unique index applied to RentalID and OptionID to ensure that the same option cannot be attached to a rental more than once

> **Note**: For simplicity, common columns such as `ID`, `CreatedAt`, and `UpdatedAt` have been omitted from the diagram below. Additionally, the explicit associations with the Tenant entity have been removed, though in the actual implementation all entities are associated with a Tenant in a multi-tenant architecture.

```mermaid
erDiagram
    companies ||--o{ renters : "can be"
    individuals ||--o{ renters : "can be"
    cars ||--o{ rentals : has
    renters ||--o{ rentals : has
    options ||--o{ rental_options : has
    rentals ||--o{ rental_options : has

    companies {
        string id
        string name
        string company_size
    }

    individuals {
        string id
        string email
        string first_name
        string last_name
    }

    renters {
        string id
        string type
    }

    cars {
        string id
        string model
    }

    rentals {
        string id
        string car_id "FK"
        string renter_id "FK"
        time starts_at
        time ends_at
    }

    options {
        string id
        string name
    }

    rental_options {
        string id
        string rental_id
        string option_id
        int count
    }
```

# ER Diagram (Full Version)

```mermaid
erDiagram
    tenants ||--o{ renters : owns
    tenants ||--o{ cars : owns
    tenants ||--o{ rentals : owns
    tenants ||--o{ options : owns
    tenants ||--o{ rental_options : owns

    renters ||--o{ companies : "class table inheritance"
    renters ||--o{ individuals : "class table inheritance"

    cars ||--o{ rentals : has
    renters ||--o{ rentals : places

    rentals ||--o{ rental_options : includes
    options ||--o{ rental_options : included_in

    tenants {
        string id PK
        string code UK
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    renters {
        string id PK
        string tenant_id FK
        string type
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    companies {
        string id PK
        string renter_id FK
        string tenant_id FK
        string name
        string company_size
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    individuals {
        string id PK
        string renter_id FK
        string tenant_id FK
        string email
        string first_name
        string last_name
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    cars {
        string id PK
        string tenant_id FK
        string model
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    rentals {
        string id PK
        string tenant_id FK
        string car_id FK
        string renter_id FK
        timestamp starts_at
        timestamp ends_at
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    options {
        string id PK
        string tenant_id FK
        string name
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    rental_options {
        string id PK
        string tenant_id FK
        string rental_id FK
        string option_id FK
        integer count
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }
```
