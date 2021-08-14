    CREATE DATABASE mistgen;
CREATE USER mistgen with encrypted password 'mistgen';
ALTER DATABASE mistgen OWNER TO mistgen;
GRANT ALL PRIVILEGES ON DATABASE mistgen to mistgen;
