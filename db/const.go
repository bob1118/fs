package db

// const USER_CREATE = "create user fsdba with password 'fsdba'"
// const DB_CREATE = "create database freeswitch with owner fsdba"
// const DBUSER_AUTH = "grant all privileges on database freeswitch to fsdba"
// const DATABASE_USER_DROP = `drop database if exists freeswitch; drop user if exists fsdba;`
const USER_CREATE = "create user %s with password '%s'"
const DB_CREATE = "create database %s with owner %s"
const DBUSER_AUTH = "grant all privileges on database %s to %s"
const DATABASE_USER_DROP = `drop database if exists %s; drop user if exists %s;`

//CREATE TABLE if not exists table_name(...)

//CDR_LEG mod_odbc_cdr table define
const CDR_LEG = `
CREATE TABLE IF NOT EXISTS %s (
	uuid varchar NOT NULL,
	otheruuid varchar,
	othertype varchar,
	ch_name varchar,
	ch_profile varchar,
	ch_direction varchar,
	ch_domain varchar,
	ch_gateway varchar,
	ch_calleridname varchar,
	ch_calleridnumber varchar,
	ch_calleeidname varchar,
	ch_calleeidnumber varchar,
	ch_destination varchar,
	ch_app varchar,
	ch_appdata varchar,
	ch_dialstatus varchar,
	hangup_cause varchar,
	hangup_q850 varchar,
	hangup_disposition varchar,
	hangup_protocause varchar,
	hangup_phrase varchar,
	start_epoch varchar,
	answer_epoch varchar,
	end_epoch varchar,
	waitsec varchar,
	billsec varchar,
	duration varchar
);
`

//fs gateway -run table define
const ACCOUNTS = `
CREATE TABLE IF NOT EXISTS %s (
	account_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	account_id varchar NOT NULL,
	account_name varchar NULL,
	account_auth varchar NULL,
	account_password varchar NULL,
	account_a1hash varchar NULL,
	account_group varchar NULL,
	account_domain varchar NULL,
	account_proxy varchar NULL,
	account_cacheable varchar NULL,
	CONSTRAINT accounts_pkey PRIMARY KEY (account_uuid),
	CONSTRAINT accounts_un UNIQUE (account_id, account_domain)
);
COMMENT ON TABLE %s IS 'sofia internal useragent account';
`
const GATEWAYS = `
CREATE TABLE IF NOT EXISTS %s (
	gateway_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	gateway_name varchar NOT NULL,
	gateway_username varchar NULL,
	gateway_realm varchar NULL,
	gateway_fromuser varchar NULL,
	gateway_fromdomain varchar NULL,
	gateway_password varchar NULL,
	gateway_extension varchar NULL,
	gateway_proxy varchar NULL,
	gateway_registerproxy varchar NULL,
	gateway_expire varchar NULL,
	gateway_register varchar NULL,
	gateway_calleridinfrom varchar NULL,
	gateway_extensionincontact varchar NULL,
	gateway_optionping varchar NULL,
	CONSTRAINT gateways_pkey PRIMARY KEY (gateway_uuid),
	CONSTRAINT gateways_un UNIQUE (gateway_name)
);
COMMENT ON TABLE %s IS 'sofia external gateways gateway';
`
const E164S = `
CREATE TABLE IF NOT EXISTS %s (
	e164_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	gateway_name varchar NOT NULL DEFAULT '',
	e164_number varchar NOT NULL,
	e164_enable bool NULL DEFAULT true,
	e164_lockin bool NULL DEFAULT false,
	e164_lockout bool NULL DEFAULT false,
	CONSTRAINT e164s_pkey PRIMARY KEY (e164_uuid),
	CONSTRAINT e164s_un UNIQUE (e164_number)
);
COMMENT ON TABLE %s IS 'phone numbers of external gateway';
`
const ACCE164 = `
CREATE TABLE IF NOT EXISTS %s (
	acce164_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	account_id varchar NOT NULL,
	account_domain varchar NOT NULL,
	gateway_name varchar NOT NULL,
	e164_number varchar NOT NULL,
	acce164_isdefault bool NOT NULL DEFAULT false
);
COMMENT ON TABLE %s IS 'account e164 number for outgoing call';
ALTER TABLE %s ADD CONSTRAINT acce164_fk FOREIGN KEY (account_id,account_domain) REFERENCES %s(account_id,account_domain) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE %s ADD CONSTRAINT acce164_fk_1 FOREIGN KEY (gateway_name) REFERENCES %s(gateway_name) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE %s ADD CONSTRAINT acce164_fk_2 FOREIGN KEY (e164_number) REFERENCES %s(e164_number) ON DELETE CASCADE ON UPDATE CASCADE;
`
const FIFOS = `
CREATE TABLE IF NOT EXISTS %s (
	fifo_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	fifo_name varchar NOT NULL,
	fifo_importance varchar NULL DEFAULT 0,
	fifo_announce varchar NULL DEFAULT '',
	fifo_holdmusic varchar NULL DEFAULT '',
	CONSTRAINT fifos_un UNIQUE (fifo_name)
);
COMMENT ON TABLE %s IS 'mod_fifo fifos';
`
const FIFOMEMBERS = `
CREATE TABLE IF NOT EXISTS %s (
	fifomember_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	fifo_name varchar NOT NULL,
	member_string varchar NOT NULL,
	member_simo varchar NULL DEFAULT 1,
	member_timeout varchar NULL DEFAULT 10,
	member_lag varchar NULL DEFAULT 10
);
COMMENT ON TABLE %s IS 'mod_fifo fifo members';
ALTER TABLE %s ADD CONSTRAINT fifomember_fk FOREIGN KEY (fifo_name) REFERENCES %s(fifo_name) ON DELETE CASCADE ON UPDATE CASCADE;
`

//fs server --run table define
const BLACKLISTS = `
CREATE TABLE IF NOT EXISTS %s (
	blacklist_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	blacklist_caller varchar NOT NULL,
	blacklist_callee varchar NOT NULL
);
COMMENT ON TABLE %s IS 'call filter blacklist include caller and callee';
`
const BGJOBS = `
CREATE TABLE IF NOT EXISTS %s (
	job_uuid uuid NOT NULL,
	job_cmd varchar,
	job_cmdarg varchar,
	job_content varchar
);
COMMENT ON TABLE %s IS 'eslclient execute bgapi command then receive EVENT BACKGROUND_JOB ';
`

//////////////////////////////////////////////////////////////
///////////////////tables default values./////////////////////
//////////////////////////////////////////////////////////////

const DEFAULT_ACCOUNTS = `insert into %s(account_id,account_name,account_auth,account_password,account_a1hash,account_group,account_domain,account_proxy,account_cacheable) values('%s','%s','%s','%s','','default','%s','%s','');`

const DEFAULT_GATEWAYS = `
insert into %s(gateway_name,gateway_username,gateway_realm,gateway_fromuser,gateway_fromdomain,gateway_password,gateway_extension,gateway_proxy,gateway_registerproxy,gateway_expire,gateway_register,gateway_calleridinfrom,gateway_extensionincontact,gateway_optionping) values
('p2p','','p2p.ip','','','','','','','','false','true','',''),
('myfsgateway','1000','10.10.10.200','1000','10.10.10.200','1234','1000','10.10.10.200','10.10.10.200','3600','false','true','true',''),
('vos_in','username','vos.ip','','','password','','','','','false','true','true',''),
('vos_out','username','vos.ip','','','password','','','','','false','true','true','')
`
const DEFAULT_E164S = `
insert into %s(e164_number)values
('1000'),
('10010'),
('10086')
`
const DEFAULT_ACCE164 = `
insert into %s(account_id,account_domain,gateway_name,e164_number, acce164_isdefault) values
('8000','mydomain','myfsgateway','1000',true),
('8001','mydomain','myfsgateway','1000',true),
('8002','mydomain','myfsgateway','1000',true),
('8003','mydomain','myfsgateway','1000',true),
('8004','mydomain','myfsgateway','1000',true)
`
const DEFAULT_FIFOS = `
insert into %s(fifo_name)values
('fifomember@fifos'),
('fifoconsumer@fifos')
`
const DEFAULT_FIFOMEMBERS = `
insert into %s(fifo_name,member_string)values
('fifomember@fifos','sofia/mydomain/8000'),
('fifomember@fifos','sofia/mydomain/8001'),
('fifomember@fifos','sofia/mydomain/8002'),
('fifomember@fifos','sofia/mydomain/8003'),
('fifomember@fifos','sofia/mydomain/8004')
`
const DEFAULT_BLACKLISTS = `
insert into %s(blacklist_caller,blacklist_callee)values
('13012345678','1000'),
('1000','13012345678')
`
