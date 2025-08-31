-- backend/database/seeds/foods.sql
INSERT INTO foods (name, brand, calories, protein, carbohydrate, fat) VALUES
('鶏むね肉', '国産', 116, 25, 0, 1.9),
('ごはん', '白米', 168, 2.5, 37, 0.3),
('納豆', 'おかめ納豆', 100, 8.5, 6, 5),
('ブロッコリー', '国産', 33, 4.3, 5.2, 0.5),
('卵', '全卵', 151, 12.3, 1.2, 10.3);
-- 他にもデータを追加すると、後の開発が楽になります