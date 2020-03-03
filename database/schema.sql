CREATE TABLE products (
	id INT8 NOT NULL DEFAULT unique_rowid(),
	title VARCHAR(100) NULL,
	FAMILY "primary" (id, title, rowid)
);

CREATE TABLE schema_lock (
	lock_id INT8 NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (lock_id ASC),
	FAMILY "primary" (lock_id)
);

CREATE TABLE schema_migrations (
	version INT8 NOT NULL,
	dirty BOOL NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (version ASC),
	FAMILY "primary" (version, dirty)
);

INSERT INTO schema_migrations (version, dirty) VALUES
	(1, false);
