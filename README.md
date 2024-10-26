CSE 312 Group Project 

## the_cookie_jar

### Description

A Learning Management System (LMS) designed to have ease-of-use and basic professor/student interactions.

### Running the Server

1. Clone this repo, go to branch of choice (`main` is production, `dev` is close to production).
2. Create an environment file [`.env`] in root directory with contents:
```
MONGODB_URI=mongodb://database:27017/
DB_NAME=the_cookie_jar
```
3. Run server with Docker from base directory with command `docker compose up --build --force-recreate` and expect:
```
cookie_lms  | Checking availability of database:27017
cookie_lms  | ...
cookie_lms  | Host database:27017 is now available
cookie_lms  | ...
cookie_lms  | -----------------------------------
cookie_lms  | the_cookie_jar server is running...
cookie_lms  | -----------------------------------
cookie_lms  | ...
```
4. You may now go to `http://localhost:8080/` on your webbrowser and interact with the server.

### Folder Layout

| Folder | Description |
|--------|-------------|
| cmd | Applications that make use of `pkg` libraries.
| docs | Any documentation related to the project.
| pkg | Our libraries used in the project.
| public | Files to send to clients based on requests (i.e. the webpage, CSS, etc.)

*Note: `main.go` is the only application not located in `cmd` folder, this is because it is the **main** server application that is expected to interact with clients for production purposes.*