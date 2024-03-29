package db

// const USER_CREATE = "create user fsdba with password 'fsdba'"
// const DB_CREATE = "create database freeswitch with owner fsdba"
// const DBUSER_AUTH = "grant all privileges on database freeswitch to fsdba"
// const DATABASE_USER_DROP = `drop database if exists freeswitch; drop user if exists fsdba;`
const USER_CREATE = "create user %s with password '%s'"
const DB_CREATE = "create database %s with owner %s"
const DBUSER_AUTH = "grant all privileges on database %s to %s"
const DATABASE_DROP = `drop database if exists %s;`
const USER_DROP = `drop user if exists %s;`

// CDR_LEG mod_odbc_cdr table define
const CDR_LEG = `
CREATE TABLE IF NOT EXISTS %s (
	uuid varchar NOT NULL,
	otheruuid varchar NOT NULL DEFAULT '',
	bonduuid varchar NOT NULL DEFAULT '',
	name varchar NOT NULL DEFAULT '',
	direction varchar NOT NULL DEFAULT '',
	sofiaprofile varchar NOT NULL DEFAULT '',
	indomain varchar NOT NULL DEFAULT '',
	ingateway varchar NOT NULL DEFAULT '',
	outdomain varchar NOT NULL DEFAULT '',
	outgateway varchar NOT NULL DEFAULT '',
	ani varchar NOT NULL DEFAULT '',
	destination varchar NOT NULL DEFAULT '',
	calleridname varchar NOT NULL DEFAULT '',
	calleridnumber varchar NOT NULL DEFAULT '',
	calleeidname varchar NOT NULL DEFAULT '',
	calleeidnumber varchar NOT NULL DEFAULT '',
	app varchar NOT NULL DEFAULT '',
	appdata varchar NOT NULL DEFAULT '',
	dialstatus varchar NOT NULL DEFAULT '',
	cause varchar NOT NULL DEFAULT '',
	q850 varchar NOT NULL DEFAULT '',
	disposition varchar NOT NULL DEFAULT '',
	protocause varchar NOT NULL DEFAULT '',
	phrase varchar NOT NULL DEFAULT '',
	startepoch varchar NOT NULL DEFAULT '',
	answerepoch varchar NOT NULL DEFAULT '',
	endepoch varchar NOT NULL DEFAULT '',
	waitsec varchar NOT NULL DEFAULT '',
	billsec varchar NOT NULL DEFAULT '',
	duration varchar NOT NULL DEFAULT ''
);
`

// fs gateway -run table define
const CONFS = `
CREATE TABLE IF NOT EXISTS %s (
	conf_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	conf_filename varchar NOT NULL,
	conf_profile varchar NULL,
	conf_content varchar NOT NULL,
	conf_newcontent varchar NULL,
	CONSTRAINT confs_pkey PRIMARY KEY (conf_uuid)
);
COMMENT ON TABLE %s IS 'switch config files that requested by xml_curl';
`

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
	profile_name varchar NOT NULL DEFAULT 'external',
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
const ACCE164S = `
CREATE TABLE IF NOT EXISTS %s (
	acce164_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	account_id varchar NOT NULL,
	account_domain varchar NOT NULL,
	gateway_name varchar NOT NULL,
	e164_number varchar NOT NULL,
	acce164_isdefault bool NOT NULL DEFAULT false,
	CONSTRAINT acce164_pkey PRIMARY KEY (acce164_uuid)
);
COMMENT ON TABLE %s IS 'account receive incoming call, dial out witch gateway e164 number for outgoing call';
ALTER TABLE %s ADD CONSTRAINT acce164_fk FOREIGN KEY (account_id,account_domain) REFERENCES %s(account_id,account_domain) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE %s ADD CONSTRAINT acce164_fk_1 FOREIGN KEY (gateway_name) REFERENCES %s(gateway_name) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE %s ADD CONSTRAINT acce164_fk_2 FOREIGN KEY (e164_number) REFERENCES %s(e164_number) ON DELETE CASCADE ON UPDATE CASCADE;
`
const E164ACCS = `
CREATE TABLE IF NOT EXISTS %s (
	e164acc_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	gateway_name varchar NOT NULL,
	e164_number varchar NOT NULL,
	account_id varchar NULL DEFAULT '',
	account_domain varchar NULL DEFAULT '',
	fifo_name varchar NULL DEFAULT '',
	e164acc_isfifo bool NOT NULL DEFAULT false,
	CONSTRAINT e164acc_pkey PRIMARY KEY (e164acc_uuid)
);
COMMENT ON TABLE %s IS 'gateway e164 receive incoming call, bridge account@domain or fifo fifo@fifoname in';
ALTER TABLE %s ADD CONSTRAINT e164acc_fk FOREIGN KEY (account_id,account_domain) REFERENCES %s(account_id,account_domain) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE %s ADD CONSTRAINT e164acc_fk_1 FOREIGN KEY (gateway_name) REFERENCES %s(gateway_name) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE %s ADD CONSTRAINT e164acc_fk_2 FOREIGN KEY (e164_number) REFERENCES %s(e164_number) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE %s ADD CONSTRAINT e164acc_fk_3 FOREIGN KEY (fifo_name) REFERENCES %s(fifo_name) ON DELETE CASCADE ON UPDATE CASCADE;
`
const FIFOS = `
CREATE TABLE IF NOT EXISTS %s (
	fifo_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	fifo_name varchar NOT NULL,
	fifo_importance varchar NULL DEFAULT 0,
	fifo_announce varchar NULL DEFAULT '',
	fifo_holdmusic varchar NULL DEFAULT '',
	CONSTRAINT fifos_pkey PRIMARY KEY (fifo_uuid),
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
	member_lag varchar NULL DEFAULT 10,
	CONSTRAINT fifomembers_pkey PRIMARY KEY (fifomember_uuid)
);
COMMENT ON TABLE %s IS 'mod_fifo fifo members';
ALTER TABLE %s ADD CONSTRAINT fifomember_fk FOREIGN KEY (fifo_name) REFERENCES %s(fifo_name) ON DELETE CASCADE ON UPDATE CASCADE;
`

// fs server --run table define
const BLACKLISTS = `
CREATE TABLE IF NOT EXISTS %s (
	blacklist_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	blacklist_caller varchar NOT NULL,
	blacklist_callee varchar NOT NULL,
	CONSTRAINT blacklists_pkey PRIMARY KEY (blacklist_uuid)
);
COMMENT ON TABLE %s IS 'call filter blacklist include caller and callee';
`
const BGJOBS = `
CREATE TABLE IF NOT EXISTS %s (
	job_uuid uuid NOT NULL,
	job_cmd varchar NOT NULL,
	job_cmdarg varchar NULL DEFAULT '',
	job_content varchar NULL DEFAULT '',
	CONSTRAINT bgjobs_pkey PRIMARY KEY (job_uuid)
);
COMMENT ON TABLE %s IS 'eslclient execute bgapi command then receive EVENT BACKGROUND_JOB ';
`
const OUTGOINGCALLS = `
CREATE TABLE IF NOT EXISTS %s (
	uuidjob uuid NOT NULL,
	uuida uuid NOT NULL,
	uuidb uuid NOT NULL,
	id varchar NOT NULL DEFAULT '',
	domain varchar NOT NULL DEFAULT '',
	e164 varchar NOT NULL DEFAULT '',
	gateway varchar NOT NULL DEFAULT '',
	ani varchar NOT NULL DEFAULT '',
	destination varchar NOT NULL DEFAULT '',
	CONSTRAINT outgoingcalls_pkey PRIMARY KEY (uuidjob)
);
COMMENT ON TABLE %s IS 'outgoingcalls bgapi command originate a bridge b';
`

//////////////////////////////////////////////////////////////
///////////////////tables default values./////////////////////
//////////////////////////////////////////////////////////////

const DEFAULT_ACCOUNTS = `insert into %s(account_id,account_name,account_auth,account_password,account_a1hash,account_group,account_domain,account_proxy,account_cacheable) values('%s','%s','%s','%s','','default','%s','%s','');`

const DEFAULT_GATEWAYS = `
insert into %s(gateway_name,gateway_username,gateway_realm,gateway_fromuser,gateway_fromdomain,gateway_password,gateway_extension,gateway_proxy,gateway_registerproxy,gateway_expire,gateway_register,gateway_calleridinfrom,gateway_extensionincontact,gateway_optionping) values
('p2p','','p2p.ip','','','','','','','','false','true','',''),
('myfsgateway','1000','10.10.10.200','1000','10.10.10.200','1234','1000','10.10.10.200','10.10.10.200','3600','true','true','true',''),
('vos_in','username','vos.ip','','','password','','','','','false','true','true',''),
('vos_out','username','vos.ip','','','password','','','','','false','true','true','')
`
const DEFAULT_E164S = `
insert into %s(gateway_name,e164_number)values
('myfsgateway','1000'),
('','10010'),
('','10086')
`
const DEFAULT_ACCE164S = `
insert into %s(account_id,account_domain,gateway_name,e164_number, acce164_isdefault) values
('8000','mydomain','myfsgateway','1000',true),
('8001','mydomain','myfsgateway','1000',true),
('8002','mydomain','myfsgateway','1000',true),
('8003','mydomain','myfsgateway','1000',true),
('8004','mydomain','myfsgateway','1000',true)
`
const DEFAULT_E164ACCES = `
insert into %s(gateway_name,e164_number,account_id,account_domain,fifo_name,e164acc_isfifo) values
('myfsgateway','1000','8001','mydomain','fifomember@fifos',false)
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
