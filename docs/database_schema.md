# Database Schema

This document provides an overview of the database schema for the Go Sample application.

## Tables

### tenants

- **id** (varchar(36), primary key): Unique identifier for the tenant
- **code** (varchar(50), unique): Tenant code
- **created_at** (timestamp with time zone): Creation timestamp
- **updated_at** (timestamp with time zone): Last update timestamp
- **deleted_at** (timestamp with time zone): Soft delete timestamp

### cars

- **id** (varchar(36), primary key): Unique identifier for the car
- **tenant_id** (varchar(36), foreign key to tenants.id): Tenant that owns the car
- **model** (varchar(255)): Car model name
- **created_at** (timestamp with time zone): Creation timestamp
- **updated_at** (timestamp with time zone): Last update timestamp
- **deleted_at** (timestamp with time zone): Soft delete timestamp

### companies

- **id** (varchar(36), primary key): Unique identifier for the company
- **renter_id** (varchar(36), foreign key to renters.id): Reference to the parent renter entity
- **tenant_id** (varchar(36), foreign key to tenants.id): Tenant that owns the company
- **name** (varchar(255)): Company name
- **company_size** (varchar(50)): Size of the company
- **created_at** (timestamp with time zone): Creation timestamp
- **updated_at** (timestamp with time zone): Last update timestamp
- **deleted_at** (timestamp with time zone): Soft delete timestamp

### individuals

- **id** (varchar(36), primary key): Unique identifier for the individual
- **renter_id** (varchar(36), foreign key to renters.id): Reference to the parent renter entity
- **tenant_id** (varchar(36), foreign key to tenants.id): Tenant that owns the individual
- **email** (varchar(255), unique): Individual's email address
- **first_name** (varchar(100)): Individual's first name (nullable)
- **last_name** (varchar(100)): Individual's last name (nullable)
- **created_at** (timestamp with time zone): Creation timestamp
- **updated_at** (timestamp with time zone): Last update timestamp
- **deleted_at** (timestamp with time zone): Soft delete timestamp

### renters

- **id** (varchar(36), primary key): Unique identifier for the renter
- **tenant_id** (varchar(36), foreign key to tenants.id): Tenant that owns the renter
- **type** (varchar(20)): Type of renter ("company" or "individual")
- **created_at** (timestamp with time zone): Creation timestamp
- **updated_at** (timestamp with time zone): Last update timestamp
- **deleted_at** (timestamp with time zone): Soft delete timestamp

### rentals

- **id** (varchar(36), primary key): Unique identifier for the rental
- **tenant_id** (varchar(36), foreign key to tenants.id): Tenant that owns the rental
- **car_id** (varchar(36), foreign key to cars.id): Car being rented
- **renter_id** (varchar(36), foreign key to renters.id): Renter who is renting the car
- **starts_at** (timestamp with time zone): Start time of the rental
- **ends_at** (timestamp with time zone): End time of the rental
- **created_at** (timestamp with time zone): Creation timestamp
- **updated_at** (timestamp with time zone): Last update timestamp
- **deleted_at** (timestamp with time zone): Soft delete timestamp

### options

- **id** (varchar(36), primary key): Unique identifier for the option
- **tenant_id** (varchar(36), foreign key to tenants.id): Tenant that owns the option
- **name** (varchar(255)): Option name
- **created_at** (timestamp with time zone): Creation timestamp
- **updated_at** (timestamp with time zone): Last update timestamp
- **deleted_at** (timestamp with time zone): Soft delete timestamp

### rental_options

- **id** (varchar(36), primary key): Unique identifier for the rental option
- **tenant_id** (varchar(36), foreign key to tenants.id): Tenant that owns the rental option
- **rental_id** (varchar(36), foreign key to rentals.id): Rental that has this option
- **option_id** (varchar(36), foreign key to options.id): Option being added to the rental
- **count** (integer): Number of this option in the rental
- **created_at** (timestamp with time zone): Creation timestamp
- **updated_at** (timestamp with time zone): Last update timestamp
- **deleted_at** (timestamp with time zone): Soft delete timestamp

## Relationships

1. **Tenant ↔ Cars** (One-to-Many)
   - One tenant can own many cars
   - Foreign key: cars.tenant_id → tenants.id

2. **Tenant ↔ Companies** (One-to-Many)
   - One tenant can own many companies
   - Foreign key: companies.tenant_id → tenants.id

3. **Tenant ↔ Individuals** (One-to-Many)
   - One tenant can own many individuals
   - Foreign key: individuals.tenant_id → tenants.id

4. **Tenant ↔ Renters** (One-to-Many)
   - One tenant can own many renters
   - Foreign key: renters.tenant_id → tenants.id

5. **Tenant ↔ Rentals** (One-to-Many)
   - One tenant can own many rentals
   - Foreign key: rentals.tenant_id → tenants.id

6. **Tenant ↔ Options** (One-to-Many)
   - One tenant can own many options
   - Foreign key: options.tenant_id → tenants.id

7. **Tenant ↔ RentalOptions** (One-to-Many)
   - One tenant can own many rental options
   - Foreign key: rental_options.tenant_id → tenants.id

8. **Renter ↔ Company** (One-to-One)
   - One renter has one company (when the renter is a company)
   - Foreign key: companies.renter_id → renters.id

9. **Renter ↔ Individual** (One-to-One)
   - One renter has one individual (when the renter is an individual)
   - Foreign key: individuals.renter_id → renters.id

10. **Car ↔ Rentals** (One-to-Many)
    - One car can be rented many times
    - Foreign key: rentals.car_id → cars.id

11. **Renter ↔ Rentals** (One-to-Many)
    - One renter can have many rentals
    - Foreign key: rentals.renter_id → renters.id

12. **Option ↔ RentalOptions** (One-to-Many)
    - One option can be used in many rental options
    - Foreign key: rental_options.option_id → options.id

13. **Rental ↔ RentalOptions** (One-to-Many)
    - One rental can have many rental options
    - Foreign key: rental_options.rental_id → rentals.id

## Class Table Inheritance

The renters table uses Class Table Inheritance pattern with companies and individuals tables:

- companies.renter_id is a foreign key referencing renters.id
- individuals.renter_id is a foreign key referencing renters.id

This relationship enforces that each company or individual record is associated with exactly one renter record, implementing a true inheritance hierarchy at the database level.
