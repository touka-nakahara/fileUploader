CREATE DATABASE file;

CREATE TABLE file.File(
	id int unsigned not null auto_increment,
	name varchar(255) not null,
	size int,
	extension varchar(255),
	description text,
	password varchar(255),
	UUID varchar(255) not null,
	thumbnail BLOB,
	is_available datetime not null default '9999-12-31 23:59:59',
	update_date datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	upload_date datetime DEFAULT CURRENT_TIMESTAMP,
	primary key (id)
);

CREATE TABLE file.Data (
  id int unsigned not null auto_increment,
  file_id int unsigned unique,
  data LONGBLOB,

  foreign key file_id_foreign_key (file_id) references file.File (id),

  primary key(id)
);

INSERT INTO file.File (name, size, extension, description, password, UUID, thumbnail) VALUES
('example1.txt', 1234, 'txt', 'This is an example text file.', 'password1', '123e4567-e89b-12d3-a456-426614174000', NULL),
('example2.jpg', 5678, 'jpg', 'This is an example image file.', 'password2', '123e4567-e89b-12d3-a456-426614174001', NULL),
('example3.pdf', 91011, 'pdf', 'This is an example PDF file.', 'password3', '123e4567-e89b-12d3-a456-426614174002', NULL);

INSERT INTO file.Data (file_id, data) VALUES
(1, 0x5468697320697320736f6d652064756d6d79206461746120666f722066696c653120696e2068657820666f726d6174), -- "This is some dummy data for file1" in hex
(2, 0x5468697320697320736f6d652064756d6d79206461746120666f722066696c653220696e2068657820666f726d6174), -- "This is some dummy data for file2" in hex
(3, 0x5468697320697320736f6d652064756d6d79206461746120666f722066696c653320696e2068657820666f726d6174); -- "This is some dummy data for file3" in hex
