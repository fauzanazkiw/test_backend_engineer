BEGIN;

-- User default untuk login: admin@example.com / password123
INSERT INTO users (name, email, password_hash)
VALUES (
    'Administrator MKP',
    'admin@example.com',
    '$2a$10$XZyrHcPVNdfe8tLq01g87uA4.jmYOyblkhBn3WgCN2y6lqOW5Deb6'
)
ON CONFLICT (email) DO NOTHING;

INSERT INTO movies (title, description, genre, duration_minutes, release_date) VALUES
    ('Interstellar', 'Sekelompok astronaut melakukan perjalanan melewati wormhole untuk mencari rumah baru bagi umat manusia.', 'Sci-Fi', 169, DATE '2014-11-07'),
    ('Parasite', 'Sebuah keluarga miskin menyusup ke rumah keluarga kaya dengan konsekuensi yang tak terduga.', 'Thriller', 132, DATE '2019-05-30'),
    ('Spider-Man: No Way Home', 'Identitas Spider-Man terungkap dan membuka jalan bagi para villain dari semesta lain.', 'Action', 148, DATE '2021-12-15')
ON CONFLICT (title) DO NOTHING;

INSERT INTO branches (name, city, address, phone) VALUES
    ('MKP Cinema Central', 'Jakarta', 'Jl. Jenderal Sudirman No. 1', '021-5550101'),
    ('MKP Cinema Bandung', 'Bandung', 'Jl. Asia Afrika No. 8', '022-5550202')
ON CONFLICT (name, city) DO NOTHING;

INSERT INTO studios (branch_id, name, total_seats)
SELECT b.id, v.studio_name, v.total_seats
FROM (VALUES
        ('MKP Cinema Central', 'Jakarta', 'Studio 1', 120),
        ('MKP Cinema Central', 'Jakarta', 'Studio 2', 96),
        ('MKP Cinema Bandung', 'Bandung', 'Studio 1', 80)
     ) AS v (branch_name, city, studio_name, total_seats)
JOIN branches b ON b.name = v.branch_name AND b.city = v.city
ON CONFLICT (branch_id, name) DO NOTHING;

INSERT INTO schedules (movie_id, branch_id, studio_id, show_date, show_time, ticket_price)
SELECT m.id, b.id, s.id, v.show_date, v.show_time, v.ticket_price
FROM (VALUES
        ('Interstellar', 'MKP Cinema Central', 'Jakarta', 'Studio 1', DATE '2026-11-30', TIME '13:00:00', 45000::BIGINT),
        ('Interstellar', 'MKP Cinema Central', 'Jakarta', 'Studio 1', DATE '2026-11-30', TIME '19:00:00', 50000::BIGINT),
        ('Parasite', 'MKP Cinema Central', 'Jakarta', 'Studio 2', DATE '2026-12-01', TIME '16:30:00', 40000::BIGINT),
        ('Spider-Man: No Way Home', 'MKP Cinema Bandung', 'Bandung', 'Studio 1', DATE '2026-12-01', TIME '20:00:00', 55000::BIGINT)
     ) AS v (movie_title, branch_name, city, studio_name, show_date, show_time, ticket_price)
JOIN movies m ON m.title = v.movie_title
JOIN branches b ON b.name = v.branch_name AND b.city = v.city
JOIN studios s ON s.branch_id = b.id AND s.name = v.studio_name
ON CONFLICT (studio_id, show_date, show_time) DO NOTHING;

COMMIT;
