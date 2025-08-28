# Database Schema Updates

This project uses GORM's AutoMigrate feature to manage database schema changes. When you modify domain models or add new entities, you need to update the corresponding database models and run migrations.

This project also uses [GORM Gen](https://gorm.io/gen/index.html) to generate type-safe database query code. After making changes to the database schema, you should regenerate the type-safe query code.

## Current Database Schema

For a complete view of the current database schema, see:

- [Database Schema Documentation](./database_schema.md)
- [ER Diagram](./er-diagram.md)

## Process for Updating Database Schema

1. **Update Domain Models**: Modify the structs in `internal/domain/model/` to reflect your changes.

2. **Update Database Models**: Update the corresponding structs in `internal/infrastructure/postgres/dbmodel/` to match your domain model changes. Ensure you:
   - Maintain the GORM annotations for table structure
   - Update the `ToDomain()` and `FromDomain()` conversion methods
   - Keep the `TableName()` method current
   - Update relationship definitions and foreign key constraints

3. **Register Models in Migration**: Add your updated or new model to the `AutoMigrate` call in `internal/infrastructure/postgres/migrate/main.go`:

   ```go
   if err := db.AutoMigrate(
       &dbmodel.Car{},
       &dbmodel.Company{},
       // ... other models
       &dbmodel.YourNewModel{}, // Add your model here
   ); err != nil {
       log.Fatal("failed to migrate database:", err)
   }
   ```

4. **Handle Foreign Key Constraints**: If you're working with polymorphic relationships or need to modify foreign key constraints, you may need to manually drop or add constraints in the migration script. The migration script includes examples of dropping problematic constraints for polymorphic relationships.

5. **Run Migration**: Execute the migration using the Makefile target:

   ```bash
   make migrate
   ```

   GORM's AutoMigrate will automatically handle:

   - Creating tables for new models
   - Adding columns for new fields
   - Modifying column types when possible
   - Creating foreign key constraints for defined relationships

   Note: AutoMigrate won't remove unused columns or tables to prevent data loss. For destructive changes, manual database operations may be required.

6. **Generate Type-Safe Query Code**: After the database schema is updated, generate the type-safe query code:

   ```bash
   make gen.gorm
   ```

   This command generates Go code that provides type-safe methods for querying the database based on your model definitions. The generated code is placed in `internal/infrastructure/postgres/query/` and should be committed to version control.
