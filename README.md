# egoist
individual twitter notification

```console
$ egoist -c config.json
```

config.json

```json
{
  "AccessToken": "00000000-0xxxxxxxxxx0xxxxx00xxxxxx0xxxxxxxxxxxx0x0",
  "AccessTokenSecret": "xx0xx00xxxxxxxxxxx0xxxxxxx00xx0xxxxxxxx0xxxxx",
  "ConsumerKey": "xxxxxx0xxxxxxxxxxxxxx0xxx",
  "ConsumerSecret": "xxxxxx0xxxxxxxxxxxxxxxxxxxx0xxxxxxxxxxxxxxxxx0x0xx"
}
```

## dependencies

```
github.com/podhmo/egoist #=0
  github.com/ChimeraCoder/anaconda #=1
    github.com/ChimeraCoder/anaconda/vendor/github.com/azr/backoff #=2
    github.com/ChimeraCoder/anaconda/vendor/github.com/dustin/go-jsonpointer #=96
      github.com/ChimeraCoder/anaconda/vendor/github.com/dustin/gojson #=97
    github.com/ChimeraCoder/anaconda/vendor/github.com/garyburd/go-oauth/oauth #=98
    github.com/ChimeraCoder/anaconda/vendor/github.com/ChimeraCoder/tokenbucket #=99
  github.com/gen2brain/beeep #=102
    github.com/godbus/dbus #=103
  gopkg.in/alecthomas/kingpin.v2 #=106
    github.com/alecthomas/template #=115
      github.com/alecthomas/template/parse #=116
    github.com/alecthomas/units #=117
```
