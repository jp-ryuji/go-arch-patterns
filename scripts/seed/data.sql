-- Sample data for the car rental platform

-- Create a sample tenant if it doesn't already exist
INSERT INTO tenants (id, code)
VALUES ('cmft9obzd0000cvndc8kb2asz', 'sample-tenant')
ON CONFLICT (code) DO NOTHING;

-- Create sample cars if they don't already exist
INSERT INTO cars (id, tenant_id, model)
VALUES
  ('cmft9obzn0001cvndwgpqx2td', 'cmft9obzd0000cvndc8kb2asz', 'Toyota Camry'),
  ('cmft9obzn0002cvndabe63bxi', 'cmft9obzd0000cvndc8kb2asz', 'Honda Civic'),
  ('cmft9obzn0003cvndsdcz4157', 'cmft9obzd0000cvndc8kb2asz', 'Ford Mustang')
ON CONFLICT DO NOTHING;
