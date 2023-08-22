# set -e

# psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
#     CREATE DATABASE database_name;
#     CREATE USER new_user_name WITH ENCRYPTED PASSWORD 'new_password';
#     GRANT ALL PRIVILEGES ON DATABASE new_user_name TO database_name;
# EOSQL
