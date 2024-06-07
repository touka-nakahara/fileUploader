CREATE DATABASE file;

CREATE TABLE file.Metadata(
	id int unsigned not null auto_increment,
	name varchar(255) not null,
	size int,
	extension varchar(255),
	description text,
  uuid varchar(255),
	password varchar(255),
	thumbnail BLOB,
	is_available datetime not null default '9999-12-31 23:59:59',
	update_date datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	upload_date datetime DEFAULT CURRENT_TIMESTAMP,
	primary key (id)
);

-- CREATE TABLE file.Data (
--   id int unsigned not null auto_increment,
--   file_id int unsigned unique,
--   data LONGBLOB,

--   foreign key file_id_foreign_key (file_id) references file.File (id),

--   primary key(id)
-- );


