## Description

`golo` prevents an application from running twice. This is useful when
you launch your task thanks to `cron`.

The basic ideas come from the original `Perl` application [solo](http://github.com/timkay/solo)
and its `Ruby` version [rolo](http://github.com/icy/rolo).

## Example usage

Try to create `ssh` port forwarding from local host to remote server.
As we don't want to run this command multiple twice, `golo` helps

```
$ go run golo.go -timeout 10  -port 4040 --no-bind -- /usr/bin/ssh MyServer -o "LocalForward localhost:4040 localhost:8888" -fN
:: Port is available. App is not running
:: Now staring application '/usr/bin/ssh' from .

$ go run golo.go -timeout 10  -port 4040 --no-bind -- /usr/bin/ssh MyServer -o "LocalForward localhost:4040 localhost:8888" -fN
:: Port is not available. App is running?
```

Below the crontab settings to create some `ssh` tunnels. If any tunnel is
broken due to network issue, `cron` will try to restart them in within 1 minute.
If the tunnel is still work, `cron` will simply exit.

```
$ crontab -l

*/1 * * * * golo -port 6432 --address 127.0.0.1 --no-bind /usr/bin/ssh zproxydev -fN
*/1 * * * * golo -port 6442 --address 127.0.0.1 --no-bind /usr/bin/ssh zproxystaging -fN
*/1 * * * * golo -port 6452 --address 127.0.0.1 --no-bind /usr/bin/ssh zproxyproduction -fN
```
