create table todoapp.users (
    id          serial primary key,
    name        text not null unique,
    password    text not null,
    role        text not null,
    created     timestamp not null,
    modified    timestamp not null
);
comment on table todoapp.users is            'ユーザ';
comment on column todoapp.users.id is        'ユーザの識別子';
comment on column todoapp.users.name is      'ユーザの名前';
comment on column todoapp.users.password is  'パスワードハッシュ';
comment on column todoapp.users.role is      'ロール';
comment on column todoapp.users.created is   'レコード作成日時';
comment on column todoapp.users.modified is  'レコード更新日時';

create table todoapp.tasks (
    id          serial primary key,
    title       text not null,
    stat        text not null,
    created     timestamp not null,
    modified    timestamp not null
);
comment on table todoapp.tasks is            'タスク';
comment on column todoapp.tasks.id is        'タスクの識別子';
comment on column todoapp.tasks.title is     'タスクのタイトル';
comment on column todoapp.tasks.stat is      'タスクの状態';
comment on column todoapp.tasks.created is   'タスク識別子';
comment on column todoapp.tasks.modified is  'タスク識別子';
