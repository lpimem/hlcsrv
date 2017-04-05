
/* HLC function tables */

create table if not exists `hlc_user`(
  `id` integer primary key,
  `name` varchar(256) not null,
  `email` varchar(1024) not null,
  `password` varchar(256) not null,
  `_slt` varchar(256) not null,
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
  `page` integer not null,
  `author` integer not null, 
  `ctime` timestamp default current_timestamp,
  `mtime` timestamp default current_timestamp
);

create table if not exists `hlc_comments`(
  `range` integer primary key,
  `comment`  TEXT
);