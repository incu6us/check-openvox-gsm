# check-openvox-gsm

Tool to check GSM signal on OpenVox. Compatible with [Sensu](https://sensuapp.org) 

### Requirements
 - [Glide](https://github.com/Masterminds/glide)
 
To install all dependent packages, run on root of the project:
```
glide i
```

### Build
```
go build .
```
 
##### Parameters:
```
  -crit int
        critical signal (default 7)
  -h    -h for help
  -host string
        192.168.1.254 (required)
  -modem string
        min: 1, max: 4 (default "1")
  -pass string
        password for digest auth (default "admin")
  -schema string
        http or https (default "http")
  -slot string
        like: /2/service?action=get_gsminfo (2 - is a slot)
  -user string
        username for digest auth (default "admin")
```