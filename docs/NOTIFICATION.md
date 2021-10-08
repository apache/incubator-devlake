# Notification


## Request
example request
```
POST /lake/notify?nouce=3-FDXxIootApWxEVtz&sign=424c2f6159bd9e9828924a53f9911059433dc14328a031e91f9802f062b495d5

{"TaskID":39,"PluginName":"jenkins","CreatedAt":"2021-09-30T15:28:00.389+08:00","UpdatedAt":"2021-09-30T15:28:00.785+08:00"}
```

## Configuration
If you want to use the notification feature, you should add two configuration key to `.env` file. 
```shell
# .env
# endpoint is the notification request url, eg: http://example.com/lake/notify
NOTIFICATION_ENDPOINT=
# screte is used to calculate signature
NOTIFICATION_SECRET=
```

## Signature
You should check the signature before accepting the notification request. We use sha256 algorithm to calculate the checksum.  
```go
// calculate checksum
sum := sha256.Sum256([]byte(requestBody + NOTIFICATION_SECRET + nouce))
return hex.EncodeToString(sum[:])
```
