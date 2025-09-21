-- Sample data for the car rental platform

-- Create a sample tenant if it doesn't already exist
INSERT INTO tenants (id, code, created_at, updated_at)
VALUES ('cmft9obzd0000cvndc8kb2asz', 'sample-tenant', NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- Create sample cars if they don't already exist
INSERT INTO cars (id, tenant_id, model, created_at, updated_at)
VALUES
  ('cmft9obzn0001cvndwgpqx2td', 'cmft9obzd0000cvndc8kb2asz', 'Toyota Camry', NOW(), NOW()),
  ('cmft9obzn0002cvndabe63bxi', 'cmft9obzd0000cvndc8kb2asz', 'Honda Civic', NOW(), NOW()),
  ('cmft9obzn0003cvndsdcz4157', 'cmft9obzd0000cvndc8kb2asz', 'Ford Mustang', NOW(), NOW())
ON CONFLICT DO NOTHING;
