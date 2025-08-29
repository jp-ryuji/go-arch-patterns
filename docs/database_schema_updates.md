# Database Schema Updates

This project uses Ent's migration feature to manage database schema changes. When you modify domain models or add new entities, you need to update the corresponding Ent schema definitions and run migrations.

This project also uses [Ent](https://entgo.io/) to generate type-safe database query code. After making changes to the database schema, you should regenerate the type-safe query code.

## Current Database Schema

For a complete view of the current database schema, see:

- [Database Schema Documentation](./database_schema.md)
- [ER Diagram](./er-diagram.md)

## Process for Updating Database Schema

1. **Update Domain Models**: Modify the structs in `internal/domain/model/` to reflect your changes.

2. **Update Ent Schema Definitions**: Update the corresponding schema definitions in `internal/infrastructure/postgres/ent/schema/` to match your domain model changes. Ensure you:
   - Maintain the Ent schema definitions for table structure
   - Update relationship definitions and foreign key constraints

3. **Run Migration**: Execute the migration using the Makefile target:

   ```bash
   make migrate
   ```

   Ent's migration will automatically handle:

   - Creating tables for new models
   - Adding columns for new fields
   - Modifying column types when possible
   - Creating foreign key constraints for defined relationships

   Note: Ent's migration won't remove unused columns or tables to prevent data loss. For destructive changes, manual database operations may be required.

4. **Generate Type-Safe Query Code**: After the database schema is updated, generate the type-safe query code:

   ```bash
   make gen.ent
   ```

   This command generates Go code that provides type-safe methods for querying the database based on your schema definitions. The generated code is placed in `internal/infrastructure/postgres/entgen/` and should be committed to version control.
