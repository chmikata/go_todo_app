create database todo;
\c todo
create schema todoapp;
create role todo with login password 'todo';
grant all privileges on schema todoapp to todo;