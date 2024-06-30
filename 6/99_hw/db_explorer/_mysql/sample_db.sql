-- Adminer 4.2.1 MySQL dump
set
  NAMES utf8;

set
  time_zone
= '+00:00';

set
  foreign_key_checks
= 0;

set
  sql_mode
= 'NO_AUTO_VALUE_ON_ZERO';

drop table if exists `items`;

create table
  `items`
(
    `id` int
(11) not null AUTO_INCREMENT,
    `title` varchar
(255) not null,
    `description` text not null,
    `updated` varchar
(255) default null,
    primary key
(`id`)
  ) ENGINE = InnoDB default CHARSET = utf8;

insert into
  `items` (`
id`,
`title
`, `description`, `updated`)
values
(
    1,
    'database/sql',
    'Рассказать про базы данных',
    'rvasily'
  ),
(
    2,
    'memcache',
    'Рассказать про мемкеш с примером использования',
    null
  );

drop table if exists `users`;

create table
  `users`
(
    `user_id` int
(11) not null AUTO_INCREMENT,
    `login` varchar
(255) not null,
    `password` varchar
(255) not null,
    `email` varchar
(255) not null,
    `info` text not null,
    `updated` varchar
(255) default null,
    primary key
(`user_id`)
  ) ENGINE = InnoDB default CHARSET = utf8;

insert into
  `users` (
    `
user_id`,
`login
`,
    `password`,
    `email`,
    `info`,
    `updated`
  )
values
(
    1,
    'rvasily',
    'love',
    'rvasily@example.com',
    'none',
    null
  );

-- 2017-11-22 23:33:12