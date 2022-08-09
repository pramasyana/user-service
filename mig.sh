#!/bin/bash
# This bash for aliasing goose command 

# Documentation:
# https://github.com/pressly/goose
# 
# Note:
# Make sure you have installed goose and there is `.env` file in your directory
# 
# Example:
# ./migrate {goose command}
# ./migrate create create_your_tables
# ./migrate status
# ./migrate up
# ./migrate down
# 
source .env
echo "RUN:"
if [ $1 = 'create' ]
then
	echo "goose -dir ./schema create $2"
	echo ""
	goose -dir ./schema create $2 sql 
else
	echo "goose -dir ./schema postgres \"user=$WRITE_DB_USER password=$WRITE_DB_PASSWORD host=$WRITE_DB_HOST dbname=$DBW_WRITE_DB_NAMENAME sslmode=disable\" $1"
	echo ""
	goose -dir ./schema postgres "user=$WRITE_DB_USER password=$WRITE_DB_PASSWORD host=$WRITE_DB_HOST dbname=$WRITE_DB_NAME sslmode=disable" $1
fi
