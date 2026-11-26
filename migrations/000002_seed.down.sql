BEGIN;

DELETE FROM schedules
WHERE movie_id IN (
    SELECT id FROM movies
    WHERE title IN ('Interstellar', 'Parasite', 'Spider-Man: No Way Home')
);

DELETE FROM studios
WHERE name IN ('Studio 1', 'Studio 2')
  AND branch_id IN (
      SELECT id FROM branches WHERE city IN ('Jakarta', 'Bandung')
  );

DELETE FROM branches WHERE city IN ('Jakarta', 'Bandung');

DELETE FROM movies
WHERE title IN ('Interstellar', 'Parasite', 'Spider-Man: No Way Home');

DELETE FROM users WHERE email = 'admin@example.com';

COMMIT;
