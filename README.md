## Usage
Clone repository and build the project:
```console
git clone https://github.com/l1raf/tg-messages.git
cd tg-messages/cmd/tgm
go build
```
Run the app like this:
```console
APP_ID=1 APP_HASH=4b3589b8f0ab99129f89768354ac7135 PHONE=89262296263 CHATS="1358715993,1005684212" DB_URI="host=localhost user=root dbname=db password=1234 port=5432" PORT=8080 N=100 ./tgm
```
Arguments:
* APP_ID (required)
* APP_HASH (required)
* PHONE (required)
* PASSWORD
* DB_URI - PostgreSQL connection string (required)
* N - number of messages to save
* PORT (8080 by default)
* CHATS - channel, group chat or user ids
