CREATE TABLE pairstats (
  first_id bigint NOT NULL,
  second_id bigint NOT NULL,
  correlation FLOAT (8) NOT NULL,
  correlation_variance FLOAT (8) NOT NULL,
  date DATE NOT NULL
);