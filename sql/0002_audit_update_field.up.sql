use billing;
alter table app_event modify createtime timestamp not null default CURRENT_TIMESTAMP;
alter table app_event modify endtime timestamp not null default CURRENT_TIMESTAMP;
