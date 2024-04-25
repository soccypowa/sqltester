# SQLtester

## Description

Simple test tool to tes connectivity to MSSQL instances (default + named). The background for the tool is to have a "windows independent" test tool for connections issues that does not rely on windows drivers. Second goal being a tool that can test if the sql connection is brooken (the *pings* feature).

Uses Windows Authentication and can be forced to use encryption with optional flag. 

The flag *pings* can be used to send a number of pings with a 30 second delay in-between. Ex: -pings 2 = ping - 30s - ping

## Usage

sqltester.exe [OPTIONS] instance

sqltester.exe -encrypt myserver

sqltester.exe -pings 2 myserver


sqltester.exe myserver\myinstance