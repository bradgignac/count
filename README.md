# count

`count` is a distributed counter built with AWS IoT, Lambda, and DynamoDB. It accepts click events from an AWS IoT button, executes a Lambda that increments a counter stored in DynamoDB, and sends push notifications with the current counter value to all connected clients.

## System Requirements

- [Apex](http://apex.run)
- [Go](https://golang.org)
