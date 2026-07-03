# SQLtester

## Bakgrund

Tänk om man hade ett verktyg så man kunde testa anslutningar till SQL server utan att vara beroende av operativsystemets (Windows 😣) nycker vad det gäller drivrutiner, begränsade kommandon via cli som kanske inte ens är installerade eller behöva öppna ODBC-hanteraren.

Därför finns *sqltester* ett litet enkelt program byggt i Go som kan testa anslutningen via ett simpelt litet kommando och du får tillbaka ett riktigt svar från databasmotorn. Allt som behövs är en liten liten connectionsträng.

## Usage

sqltester --conn "sqlserver://user:password@host[:port][/instance][?option1=value&option2=value]"

sqltester --conn "sqlserver://sa:password@host"

sqltester --conn "sqlserver://host?trusted_connection=true"

För ytterligare information ta en titt på: https://github.com/microsoft/go-mssqldb