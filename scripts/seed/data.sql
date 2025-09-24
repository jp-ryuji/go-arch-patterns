-- Sample data for the car rental platform

-- Create a sample tenant if it doesn't already exist
INSERT INTO tenants (id, code, created_at, updated_at)
VALUES ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'sample-tenant', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Create sample car options if they don't already exist
INSERT INTO car_options (id, tenant_id, name, created_at, updated_at)
VALUES
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z1', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'GPS Navigation', NOW(), NOW()),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z2', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'Child Safety Seat', NOW(), NOW()),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z3', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'Winter Tires', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Create sample cars if they don't already exist
INSERT INTO cars (id, tenant_id, model, created_at, updated_at)
VALUES
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z4', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'Toyota Camry', NOW(), NOW()),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z5', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'Honda Civic', NOW(), NOW()),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z6', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'Ford Mustang', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Create sample renters if they don't already exist
INSERT INTO renters (id, tenant_id, type, created_at, updated_at)
VALUES
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z7', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'individual', NOW(), NOW()),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z8', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'company', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Create sample individuals if they don't already exist
INSERT INTO individuals (id, renter_id, tenant_id, email, first_name, last_name, created_at, updated_at)
VALUES
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z9', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z7', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'john.doe@example.com', 'John', 'Doe', NOW(), NOW()),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZA', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z7', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'jane.smith@example.com', 'Jane', 'Smith', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Create sample companies if they don't already exist
INSERT INTO companies (id, renter_id, tenant_id, name, company_size, created_at, updated_at)
VALUES
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZB', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z8', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'Tech Innovations Inc.', 'small', NOW(), NOW()),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZC', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z8', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', 'Global Solutions Ltd.', 'medium', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Create sample rentals if they don't already exist
INSERT INTO rentals (id, tenant_id, car_id, renter_id, starts_at, ends_at, created_at, updated_at)
VALUES
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZD', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z4', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z7', NOW(), NOW() + INTERVAL '7 days', NOW(), NOW()),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZE', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z5', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z8', NOW(), NOW() + INTERVAL '5 days', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Create sample rental options if they don't already exist
INSERT INTO rental_options (id, tenant_id, rental_id, option_id, count, created_at, updated_at)
VALUES
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZF', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZD', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z1', 1, NOW(), NOW()),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZG', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZE', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z2', 2, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Create sample outbox entries if they don't already exist
INSERT INTO outboxes (id, aggregate_type, aggregate_id, event_type, payload, created_at, status, version)
VALUES
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZH', 'car', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z4', 'car_created', '{"model": "Toyota Camry", "tenant_id": "01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z0"}', NOW(), 'pending', 1),
  ('01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZI', 'rental', '01GQMF65J0Z0Z0Z0Z0Z0Z0Z0ZD', 'rental_created', '{"car_id": "01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z4", "renter_id": "01GQMF65J0Z0Z0Z0Z0Z0Z0Z0Z7", "starts_at": "2025-09-23T09:00:00Z"}', NOW(), 'processed', 1)
ON CONFLICT DO NOTHING;
