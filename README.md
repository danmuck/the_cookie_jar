CSE 312 Group Project

## the_cookie_jar

### Description

A Learning Management System (LMS) designed to have ease-of-use and basic professor/student interactions.

## Project Part 2: Production Server

#### Objective #1 Multimedia Uploads

After signing up and logging in you can update your profile picture by right clicking your name in the top right corner and navigating to account settings

#### Objective #2 WebSocket Interactions

Click the cookie in the top left corner to be brought back to the create classroom screen.

Create a new classroom and then click on the name of your newly created classroom and pick either of two options as both are WebSocket driven.

#### Objective #3 Deployment and Encryption

# [Our app is deployed here](https://thecookiejar.gensosekai.com)

## Project Part 3: Finishing Up

Our app is deployed at the link above ^

#### Objective #1 A Sense of Time

After signing up and logging in you can create a classroom followed by clicking on the new classroom name. Start a game by clicking 'Class Game' and then 'Start'. The rounds are timed and this is displayed to the users in real time using websockets.

#### Objective #2 DoS Protection (IP rate limiting)

Implementated as specified in the handout.

#### Objective #3 Creativity and Documentation

Implemented ReCaptcha to verify that users are created by humans and not bots.

##### Testing Procedure

1. Navigate to this [link](https://thecookiejar.gensosekai.com) (public deployment).
2. Register an account with a valid username/password, without clicking ReCaptcha.
3. Click register and verify that you are given an error about not doing ReCaptcha.
4. Enter the same username and password for registering again, click ReCaptcha.
5. Click register and verify that you are registered.
6. Log in with account you just registered, without clicking ReCaptcha.
7. Click login and verify that you are given an error about not doing ReCaptcha.
8. Log in with account again, click ReCaptcha.
9. Verify that you were redirected to the classrooms page indicating a successful login.

---

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

| Folder | Description                                                              |
| ------ | ------------------------------------------------------------------------ |
| cmd    | Applications that make use of `pkg` libraries.                           |
| docs   | Any documentation related to the project.                                |
| pkg    | Our libraries used in the project.                                       |
| public | Files to send to clients based on requests (i.e. the webpage, CSS, etc.) |

_Note: `main.go` is the only application not located in `cmd` folder, this is because it is the **main** server application that is expected to interact with clients for production purposes._
