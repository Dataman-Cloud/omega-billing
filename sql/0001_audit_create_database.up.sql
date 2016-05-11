CREATE DATABASE IF NOT EXISTS billing;

CREATE TABLE IF NOT EXISTS app_event (
  id bigint unsigned primary key auto_increment,
  uid bigint(20) not null,
  cid bigint(20) not null,
  appname varchar(64) not null,
  active tinyint(1) not null,
  starttime timestamp not null default CURRENT_TIMESTAMP,
  endtime timestamp not null default CURRENT_TIMESTAMP,
  cpus float(6, 2) not null,
  mem float(6, 2) not null,
  instances int(10) not null
);
