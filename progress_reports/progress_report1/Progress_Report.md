#   Progress Report 1 sdu.store



We created standard CRUD system for the first week, get User, Get All Users, Create User
![img.png](img.png)

![img_1.png](img_1.png)

![img_2.png](img_2.png)

What we did
---

The sdu.store platform uses a relational database management system to store and manage customer information. The database consists of three main tables: User, Session, and Userdata.

The User table stores information about the platform's registered users, including their login credentials (username and password), and unique ID. The Login and Username fields provide a way to identify users, while the Password field is used to securely store users' passwords.

The Session table stores information about the current session of each user. This includes the user's ID, a unique UUID that identifies the session, the time when the session was created, and the time of the user's last login. The DeletedAt field is used to store the time when the session has ended, allowing for easy management of expired sessions.

The Userdata table stores additional information about each user, including their first and last name, phone number, country code, ZIP code, and birthday. The UserId field is used to link the user's data to their account in the User table.

By using these tables, the sdu.store platform can efficiently manage and store customer information while ensuring data privacy and security. The use of a relational database management system allows for efficient querying and manipulation of the data, providing a robust and scalable solution for the platform's needs.

We met as a team and discussed the main idea, how and what we will do and divided the tasks

Daryn is responsible for commits, merging, and whether the application works correctly
Zangar for database and server side

Nurali for the server part

Nurdaulet frontend

Aruzhan frontend

We've split up for now. And made simple requests on postman


upd.15.02 21:47 sorry for changing the first progress, we didn't add anything new just created a readme file and sent the project idea there. And decided to write who did what. Since it is not clear from the text / screenshot