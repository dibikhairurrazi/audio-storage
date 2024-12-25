INSERT INTO users (id, name, email)
VALUES
    (1, 'Foo Bar', 'foobar@something.com'),
    (2, 'John Doe', 'johndoe@some.com'),
    (3, 'Mary Jane', 'mjane@something.com')
ON CONFLICT do nothing;