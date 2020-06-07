FB CreateDate Approximator
==========================

This was created in relation to the issue where Filipinos found dummy accounts in their names.

It even caught the attention of the national news. [Users report duplicate, dummy Facebook accounts in PH](https://www.rappler.com/nation/263121-users-report-duplicate-facebook-dummy-accounts-philippines)

The purpose of this script it to approximate when these dummy accounts were created to add context to them.

There are a few limitations to this though:

# Authentication

In order to parse information from facebook properly one must be authenticated on facebook.
I didn't want people to supply their username and passwords so I decided to resort to using cookies. I understand that it isn't better or more secure but at-least only those who actually knows what their doing can use the script.

# Usage Limits

Facebook spam blocks you for a few minutes when doing a lot of requests even via the frontend so use this scarcely.

# Approximation Only

Since we cannot actually find out when an account is created we can only approximate. One of the ways to do this is to find the earliest activity on that account. First post? first profile picture?


How it works?
=============
Facebook uses mysql and with any sql data store we know that most id columns are auto incremented. Basically, facebook user ids are sequencial. With that in mind we deduce that the user with user ids before your user id was created before your account, and users with users ids after your user id was created after yours. 

It also helps to confirm that the first account on facebook was of Mark Zuckerberg with the user id of 4. 1 to 3 might have been test accounts. But going to www.facebook.com/4 will redirect your to Mark's account.

How to use it?
==============
- Just download this repo
- Create a text file named `cookie.txt` and copy paste your facebook cookies there.
- Run the `webservice.exe` executable. This will run a web server at port 8080. 
- On your browser go to `http://localhost:8080/fb/fb_username_or_ID`, example: [http://localhost:8080/fb/madziikoy](http://localhost:8080/fb/madziikoy) and wait for awhile and it should return the year the account was created. 

There is a lot of room for improvement but I did this in a hurry so this is what it is. Feel free to fork and improve.

Legal
=====
This is for educational purposes only.