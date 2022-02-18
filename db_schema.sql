DROP SCHEMA public CASCADE;
CREATE SCHEMA public;

CREATE TABLE Roles (
	role_id		INT PRIMARY KEY,
	role_name	TEXT,
	description	TEXT
);

CREATE TYPE permission_type as ENUM ();
CREATE TABLE Permissions (
	permission_id		INT PRIMARY KEY,
	permission_title	TEXT,
	permission_type		permission_type,
	description			TEXT
);

CREATE TABLE Role_Permission (
	role_id			INT REFERENCES Roles(role_id) ON DELETE CASCADE,
	permission_id	INT	REFERENCES Permissions(permission_id) ON DELETE CASCADE,
	PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE Status (
	status_id		INT PRIMARY KEY,
	status_name		TEXT,
	description		TEXT
);

CREATE TABLE Process (
	process_id		INT PRIMARY KEY,
	process_name	TEXT,
	min_amount		FLOAT,
	max_amount		FLOAT,
	status			INT REFERENCES Status(status_id),
	approvers		TEXT --NOT SURE (I think we dont need this?)
);

CREATE TABLE Approvers (
	role_id		INT REFERENCES Roles(role_id) ON DELETE CASCADE,
	process_id	INT REFERENCES Process(process_id) ON DELETE CASCADE,
	precedence	TEXT, --NOT SURE
	PRIMARY KEY (role_id, process_id)
);

CREATE TABLE Users (
	user_id			INT PRIMARY KEY,
	name			TEXT,
	password		TEXT, -- can store password here?
	salt			TEXT,
	email			TEXT,
	contact_no		TEXT,	
	roles			INT,
	subcomms		TEXT, -- not sure
	account_status	TEXT,
	FOREIGN KEY (roles) REFERENCES Roles (role_id)
);

CREATE TABLE Suppliers (
	supplier_id			INT,
	quotation_id		INT,
	supplier_name		TEXT,
	contact_person		TEXT,
	contact_number		TEXT,
	unit_price			FLOAT,
	total_price			FLOAT,
	remarks				TEXT,
 	PRIMARY KEY (supplier_id, quotation_id)
);


CREATE TABLE Quotations (
	quotation_id		INT PRIMARY KEY,
	event_name			TEXT,
	item_description	TEXT,
	item_quantity		INT,
	student_name		TEXT,
	obtained_date		DATE,
	student_contact_no	TEXT,
	suppliers			TEXT,
	reason				TEXT,
	selected_supplier	INT,
	sum					INT,
	process				INT REFERENCES Process(process_id) ON DELETE SET NULL,
	status				INT,
	created_by			TEXT, -- NOT SURE
	creation_time		TIMESTAMP,
	modified_by			TEXT, -- NOT SURE
	modification_time	TIMESTAMP
);

CREATE TYPE approval_stage AS ENUM ();
CREATE TABLE Quotation_Status (
	status_id		INT REFERENCES Status(status_id) ON DELETE CASCADE,
	quotation_id	INT REFERENCES Quotations(quotation_id) ON DELETE CASCADE,
	approval_stage	approval_stage,
	PRIMARY KEY(status_id, quotation_id)
);

CREATE TABLE Email (
	email_id		INT PRIMARY KEY,
	email_subject	TEXT,
	email_body		TEXT,
	recipient		TEXT
);

CREATE TYPE trigger_status AS ENUM ();
CREATE TABLE Process_Email (
	email_id		INT,
	process_id		INT,
	trigger_status	trigger_status,
	PRIMARY KEY(email_id, process_id, trigger_status)
);
