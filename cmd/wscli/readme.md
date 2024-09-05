The usage of the `wscli` command is as follows:

```
wscli [ws://][wss://][username:password@]remote:port
```

The `wscli` command supports both `ws://` and `wss://` protocols. The `username` and `password` parameters are optional and can be included if required for authentication.

Example usage:

```
wscli ws://remote:port
wscli wss://remote:port
wscli username:password@remote:port
```