use billing;
alter table app_event change createtime starttime timestamp not null default CURRENT_TIMESTAMP;
alter table app_event modify starttime timestamp not null default CURRENT_TIMESTAMP;
alter table app_event modify endtime timestamp not null default CURRENT_TIMESTAMP;
