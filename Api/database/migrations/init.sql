CREATE TABLE IF NOT EXISTS guests
    (
        guest_name VARCHAR(255) NOT NULL,
        table_number VARCHAR(255) NOT NULL,
        accompanying_guests INTEGER NOT NULL,
        time_arrived varchar(255),
        time_left varchar(255),
        PRIMARY KEY (guest_name)
    );