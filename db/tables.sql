/* HLC function tables */
create table if not exists `hlc_user`(
  `id` integer primary key,
  `name` varchar(256) null,
  `email` varchar(1024) not null,
  `password` varchar(512) null,
  `_slt` varchar(512) not null,
  `_status` integer default 1,
  `ctime` timestamp default current_timestamp,
  `mtime` timestamp default current_timestamp
);

create unique index if not exists `idx_user_name` on `hlc_user` (`name`);
create unique index if not exists `idx_user_email` on `hlc_user` (`email`);

create table if not exists `hlc_page`(
  `id` integer primary key,
  `title` varchar(512) not null,
  `url` varchar(512) not null, 
  `version` integer default 1,
  `_type` integer default 1, 
  `_note` varchar(2048) null,
  `ctime` timestamp default current_timestamp,
  `mtime` timestamp default current_timestamp
);

create unique index if not exists `idx_page_url` on `hlc_page` (`url`);

create table if not exists `hlc_range`(
  `id` integer primary key,
  `anchor` varchar(255), 
  `start` varchar(255), 
  `startOffset` integer, 
  `end` varchar(255), 
  `endOffset` integer, 
  `text` varchar(1024),
  `option` varchar(1024) default '',
  `page` integer not null,
  `author` integer not null, 
  `ctime` timestamp default current_timestamp,
  `mtime` timestamp default current_timestamp
);

create table if not exists `hlc_comment`(
  `range` integer primary key,
  `comment`  TEXT
);

/**
 * Apr 11, 2017
 *
 * -- alter table `hlc_user` add column `avatar` varchar(256);
 */
create table if not exists `hlc_google_auth`(
  `google_id` varchar primary key,
  `uid` integer,
  `picture` varchar(32),
  `last_access` timestamp default current_timestamp,
  `ctime` timestamp default current_timestamp
);

create unique index if not exists `idx_google_auth_uid` on `hlc_google_auth` (`uid`);

create table if not exists `hlc_session`(
  `id` varchar PRIMARY KEY,
  `uid` integer,
  `last_access` timestamp default current_timestamp
);

create unique index if not exists `idx_session_uid` on `hlc_google_auth` (`uid`);

create table if not exists `permission` (
  `id` integer primary key,
  `user` integer, 
  `uri` varchar(256),
  `ctime` timestamp default current_timestamp
);

create unique index if not exists `uidx_permission_user_uri` on `permission` (`user`, `uri`);

create table if not exists `restriction` (
  `uri` varchar(256) primary key
);