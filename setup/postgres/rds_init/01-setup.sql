CREATE DATABASE waypoint
  ENCODING = 'UTF8'
  LOCALE_PROVIDER = 'icu'
  ICU_LOCALE = 'et-EE'
  TEMPLATE = template0;

CREATE USER dbuser WITH password 'CHANGE_ME' createrole valid until 'infinity';
GRANT ALL ON database waypoint TO dbuser;
ALTER database waypoint owner TO dbuser;