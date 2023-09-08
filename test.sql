
CREATE TABLE IF NOT EXISTS data(
  id serial,
  data_id int,
  link text null,
  describe text null,
  pdf_link text null,
    primary key (id),
  foreign key (data_id)
      references under_cells (id) on delete cascade
);


CREATE TABLE IF NOT EXISTS under_cells(
   id serial,
   cell_id int,
   name varchar(200) not null,
   primary key (id),
   foreign key (cell_id)
    references cell (id) on delete cascade
);



CREATE TABLE IF NOT EXISTS cell(
    id serial,
    name varchar(200) not null,
    user_id int,
    primary key (id),
    foreign key (user_id)
    references "user" (id) on delete cascade
);



CREATE TABLE IF NOT EXISTS  "user"(
    id serial,
    nickname varchar(100) unique not null,
    first_name varchar(100) not null,
    last_name varchar(100) not null,
    role varchar(5) default 'user',
    primary key (id)
);