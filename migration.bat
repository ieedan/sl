@echo off
SET DATABASE=database.db
SET SQLFILE=migration.sql

echo Running migration from %SQLFILE% on database %DATABASE%
sqlite3 %DATABASE% < %SQLFILE%

echo All done.