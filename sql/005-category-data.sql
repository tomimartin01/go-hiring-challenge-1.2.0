-- Insert 3 categories
INSERT INTO categories (code, name) VALUES
('CLOTHING', 'Clothing'),
('SHOES', 'Shoes'),
('ACCESSORIES', 'Accessories');

-- Insert product-category relationships
-- PROD001, PROD004, PROD007 belong to "Clothing"
INSERT INTO product_categories (product_id, category_id) VALUES
((SELECT id FROM products WHERE code = 'PROD001'), (SELECT id FROM categories WHERE code = 'CLOTHING')),
((SELECT id FROM products WHERE code = 'PROD004'), (SELECT id FROM categories WHERE code = 'CLOTHING')),
((SELECT id FROM products WHERE code = 'PROD007'), (SELECT id FROM categories WHERE code = 'CLOTHING'));

-- PROD002, PROD006 belong to "Shoes"
INSERT INTO product_categories (product_id, category_id) VALUES
((SELECT id FROM products WHERE code = 'PROD002'), (SELECT id FROM categories WHERE code = 'SHOES')),
((SELECT id FROM products WHERE code = 'PROD006'), (SELECT id FROM categories WHERE code = 'SHOES'));

-- PROD003, PROD005, PROD008 belong to "Accessories"
INSERT INTO product_categories (product_id, category_id) VALUES
((SELECT id FROM products WHERE code = 'PROD003'), (SELECT id FROM categories WHERE code = 'ACCESSORIES')),
((SELECT id FROM products WHERE code = 'PROD005'), (SELECT id FROM categories WHERE code = 'ACCESSORIES')),
((SELECT id FROM products WHERE code = 'PROD008'), (SELECT id FROM categories WHERE code = 'ACCESSORIES'));

