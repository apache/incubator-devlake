#!/usr/bin/env bash

if [ -n "$MONGO_INITDB_ROOT_USERNAME" ]; then
    echo "Creating mongo users..."
    mongo admin --host localhost -u $MONGO_INITDB_ROOT_USERNAME \
        -p $MONGO_INITDB_ROOT_PASSWORD --eval "db.createUser({
                user: '$MONGO_INITDB_ADMIN_USERNAME', pwd: '$MONGO_INITDB_ADMIN_PASSWORD',
                roles: [{role: 'userAdminAnyDatabase', db: 'admin'}]
           })"
    echo "Mongo users created."
fi
