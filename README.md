### CSE 312 Group Project 

# the_cookie_jar API

### the_cookie_jar is a Learning Management System designed to be both modular and extensible, using a plugin based architecture


clone this repo

create file [.env] with the following variables set for local developement:
```
  MONGODB_URI=mongodb://database:27017/
  DB_NAME=the_cookie_jar
```

from project directory: 

  ```
    clear; docker compose up --build --force-recreate
  ```

EXPECT:
```
  app-1       | >> [db] Pinged your deployment. You successfully connected to MongoDB!
```