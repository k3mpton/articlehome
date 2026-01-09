-- +goose Up 
CREATE TABLE Articles ( 
    Id serial primary key, 
    Name varchar(40), 
    Description varchar not null, 
    Rating float default null
    -- created_at Timestamp default now(), 
    -- updated_at Timestamp default now()
);

-- +goose Down
drop table if exists Articles; 