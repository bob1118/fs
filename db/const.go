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
const CDR_LEG = `
CREATE TABLE %s (
	id serial NOT NULL,
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
	duration varchar,
	CONSTRAINT cdr_table_a_leg_pkey PRIMARY KEY (id)
);
`
const CC_ACCOUNTS = `
CREATE TABLE cc_accounts (
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
	CONSTRAINT cc_accounts_pkey PRIMARY KEY (account_uuid),
	CONSTRAINT cc_accounts_un UNIQUE (account_id, account_domain)
);
COMMENT ON TABLE public.cc_accounts IS 'sofia internal useragent account';
`
const CC_GATEWAYS = `
CREATE TABLE cc_gateways (
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
	CONSTRAINT cc_gateways_pkey PRIMARY KEY (gateway_uuid),
	CONSTRAINT cc_gateways_un UNIQUE (gateway_name)
);
COMMENT ON TABLE public.cc_gateways IS 'sofia external gateways gateway';
`
const CC_E164S = `
CREATE TABLE cc_e164s (
	e164_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	gateway_name varchar NOT NULL DEFAULT '',
	e164_number varchar NOT NULL,
	e164_enable bool NULL DEFAULT true,
	e164_lockin bool NULL DEFAULT false,
	e164_lockout bool NULL DEFAULT false,
	CONSTRAINT cc_e164s_pkey PRIMARY KEY (e164_uuid),
	CONSTRAINT cc_e164s_un UNIQUE (e164_number)
);
COMMENT ON TABLE public.cc_e164s IS 'phone numbers of external gateway';
`
const CC_ACCE164 = `
CREATE TABLE cc_acce164 (
	acce164_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	account_id varchar NOT NULL,
	account_domain varchar NOT NULL,
	gateway_name varchar NOT NULL,
	e164_number varchar NOT NULL,
	acce164_isdefault bool NOT NULL DEFAULT false
);
COMMENT ON TABLE public.cc_acce164 IS 'useragent dial out through the gateway';
ALTER TABLE cc_acce164 ADD CONSTRAINT cc_acce164_fk FOREIGN KEY (account_id,account_domain) REFERENCES cc_accounts(account_id,account_domain) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE cc_acce164 ADD CONSTRAINT cc_acce164_fk_1 FOREIGN KEY (gateway_name) REFERENCES cc_gateways(gateway_name) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE cc_acce164 ADD CONSTRAINT cc_acce164_fk_2 FOREIGN KEY (e164_number) REFERENCES cc_e164s(e164_number) ON DELETE CASCADE ON UPDATE CASCADE;
`
const CC_FIFOS = `
CREATE TABLE cc_fifos (
	fifo_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	fifo_name varchar NOT NULL,
	fifo_importance varchar NULL DEFAULT 0,
	fifo_announce varchar NULL DEFAULT '',
	fifo_holdmusic varchar NULL DEFAULT '',
	CONSTRAINT cc_fifos_un UNIQUE (fifo_name)
);
COMMENT ON TABLE public.cc_fifos IS 'mod_fifo fifo';
`
const CC_FIFOMEMBER = `
CREATE TABLE cc_fifomember (
	fifomember_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	fifo_name varchar NOT NULL,
	member_string varchar NOT NULL,
	member_simo varchar NULL DEFAULT 1,
	member_timeout varchar NULL DEFAULT 10,
	member_lag varchar NULL DEFAULT 10
);
COMMENT ON TABLE public.cc_fifomember IS 'fifo''s member';
ALTER TABLE public.cc_fifomember ADD CONSTRAINT cc_fifomember_fk FOREIGN KEY (fifo_name) REFERENCES cc_fifos(fifo_name) ON DELETE CASCADE ON UPDATE CASCADE;
`
const CC_BLACKLIST = `
CREATE TABLE cc_blacklist (
	blacklist_uuid uuid NOT NULL DEFAULT gen_random_uuid(),
	blacklist_caller varchar NOT NULL,
	blacklist_callee varchar NOT NULL
);
COMMENT ON TABLE public.cc_blacklist IS 'call filter blacklist include caller and callee';
`
const CC_BGJOBS = `
CREATE TABLE cc_bgjobs (
	job_uuid uuid NOT NULL,
	job_cmd varchar,
	job_cmdarg varchar,
	job_content varchar
);
COMMENT ON TABLE public.cc_bgjobs IS 'eslclient execute bgapi command then receive BACKGROUND_JOB ';
`

//freeswitdh db tables default values.

const DEFAULT_ACCOUNTS = "insert into cc_accounts(account_id,account_name,account_auth,account_password,account_a1hash,account_group,account_domain,account_proxy,account_cacheable) values('%s','%s','%s','%s','','default','%s','%s','');"

const DEFAULT_GATEWAYS = `
insert into cc_gateways(gateway_name,gateway_username,gateway_realm,gateway_fromuser,gateway_fromdomain,gateway_password,gateway_extension,gateway_proxy,gateway_registerproxy,gateway_expire,gateway_register,gateway_calleridinfrom,gateway_extensionincontact,gateway_optionping) values
('p2pgatewayname','','p2p.ip','','','','','','','','false','true','',''),
('myfsgateway','1000','10.10.10.200','1000','10.10.10.200','1234','1000','10.10.10.200','10.10.10.200','3600','false','true','true',''),
('vos_in','username','vos.ip','','','password','','','','','false','true','true',''),
('vos_out','username','vos.ip','','','password','','','','','false','true','true','')
`
const DEFAULT_E164S = `
insert into cc_e164s (e164_number)values
('10000'),
('10010'),
('10086')
`
const DEFAULT_ACCE164 = `
insert into cc_acce164(account_id,account_domain,gateway_name,e164_number, acce164_isdefault) values
('8000','1.domain','myfsgateway','10000',true),
('8001','1.domain','myfsgateway','10000',true),
('8002','1.domain','myfsgateway','10000',true),
('8003','1.domain','myfsgateway','10000',true),
('8004','1.domain','myfsgateway','10000',true)
`
const DEFAULT_FIFOS = `
insert into cc_fifos(fifo_name)values
('fifomember@fifos'),
('fifoconsumer@fifos')
`
const DEFAULT_FIFOMEMBER = `
insert into cc_fifomember(fifo_name,member_string)values
('fifomember@fifos','sofia/1.domain/8000'),
('fifomember@fifos','sofia/1.domain/8001'),
('fifomember@fifos','sofia/1.domain/8002'),
('fifomember@fifos','sofia/1.domain/8003'),
('fifomember@fifos','sofia/1.domain/8004')
`
const DEFAULT_BLACKLIST = `
insert into cc_blacklist(blacklist_caller,blacklist_callee)values
('13012345678','10000'),
('10000','13012345678')
`
