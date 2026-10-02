# SQLtester

## Background

So I had a long script that tested connections to SQL Server with every thinkable driver in Windows. It was great, but what if I needed a simple tester that was not relying on the drivers of Windows, and what if I needed to test the connectivity from different OS 🤔.

That is why I came up with **sqltester**, a simple utility built in **go** that easily tests connectivity and receives back a *real* answer from the database engine. Everything needed is a connection string 😄.

## Usage

sqltester --conn "sqlserver://user:password@host[:port][/instance][?option1=value&option2=value]"

sqltester --conn "sqlserver://sa:password@host"

sqltester --conn "sqlserver://host?trusted_connection=true"

For driver related information check out: https://github.com/microsoft/go-mssqldb