CT=$(docker run -d -v $(pwd)../:/pefi postgres:alpine)

sleep 2

docker exec --user postgres -it $CT psql -f pefi/test/database_setup.sql

docker exec --user postgres -it $CT psql -f pefi/test/database_teardown.sql

docker rm -f $CT
