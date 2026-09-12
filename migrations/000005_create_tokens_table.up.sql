create table if not exists tokens(
    hash bytea PRIMARY KEY,
    user_id bigint not null REFERENCES users ON DELETE CASCADE,
    expires timestamp(0) with time zone NOT NULL,
    scope text not null
);