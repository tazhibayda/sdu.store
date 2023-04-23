# Report

[Link to video explanation](https://youtu.be/zDDMd2p2Df4)

[Link to midterm commit](https://github.com/tazhibayda/sdu.store/tree/105c5cf6915d60514f77e7c9d392b54ace8009c8)


## Introduction:
SDU store is an online store that offers various products for sale. The store has three types of users, including clients, staff, and admin. Each user has different levels of access to the store's functionalities.

Client User Access:
The client user is the basic user of the SDU store. Clients can browse through the store's product catalog, view product details, add items to the cart, and make purchases. They can also track their order status, view their order history, and manage their account settings.

Staff User Access:
Staff users have additional privileges compared to client users. They can access the store's backend system, allowing them to change product data, add or remove categories, and manage delivery settings. They can also view orders and customer details, and they have the ability to cancel orders and issue refunds.

Admin User Access:
Admin users have the highest level of access to the SDU store's functionalities. They can perform all the tasks that staff users can do and more. Admin users have access to user data, and they can change user roles, such as changing a client user to a staff user or an admin user. They can also access the store's financial reports and view detailed sales data.

## What We Did?

## Team members


- 200103210 Daryn Tazhibay
- 200103371 Nurali Umirzak
- 200103287 Zanngar Zhumagiyev
- 200103251 Diyarova Aruzhan
- 200103429 Nurdaulet Kalidolla

### What each member did

## How To Run The Code?
Instructions for Running the Project:

1. Create the PostgreSQL database.
   Before you can run the project, you need to create a PostgreSQL database. You can use a tool like pgAdmin or the psql command-line tool to create the database.

2. Create dbInfo.go file
   Create a new file called dbInfo.go in the root directory of the project. This file should contain the information required to connect to the database. Use the following format for the file:

![img.png](img.png)

3. Make sure to replace password, username and database name with your actual PostgreSQL database password and database name, respectively.Make sure to replace your_password_here and your_database_name_here with your actual PostgreSQL database password and database name, respectively.

4. Run the code
   To run the code, open a terminal and navigate to the project's root directory. Then, run the following command:
```
go run . -dbRestart
```

5. This command will start the server and automatically create the necessary database tables. If you make any changes to the database schema, you can use the -dbRestart flag to drop and recreate the tables.

Once the server is running, you can access it by opening a web browser and navigating to http://localhost:9090. From there, you can browse the products, add items to your cart, and complete orders. The admin site can be accessed by logging in with an admin or staff account.

## Explanation Of Each Feature
