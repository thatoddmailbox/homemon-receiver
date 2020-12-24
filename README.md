# homemon-receiver
This program receives packets sent over UDP by [homemon-daemon](https://github.com/thatoddmailbox/homemon-daemon).

Note that this is designed to work together with [homemon-server](https://github.com/thatoddmailbox/homemon-server). To that end, it inserts data into homemon-server's database, and so you will need to run the migrations in that repository, even if you're only using this code.