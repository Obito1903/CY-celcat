# CYtech Celcat V2.0

## Build

### Dependencies

- golang
  - github.com/arran4/golang-ical
  - github.com/gorilla/mux
- Chromium (For PNG output only)

``` sh
go get github.com/arran4/golang-ical
go get github.com/gorilla/mux
```

### Build project

``` sh
go build ./cmd/cy-celcat/main.go
```

## Usage

### Configuration

rename `example.config.json` into `config.json` then edit it to match your situation.

Options avalaible :

|   JSON field   |  CLI Option  |   Type   | Desc                                                                                                      |
| :------------: | :----------: | :------: | --------------------------------------------------------------------------------------------------------- |
|   `userName`   |    `user`    | `string` | Username for celcat.                                                                                      |
| `userPassword` |    `pass`    | `string` | Password for celcat.                                                                                      |
|  `celcatHost`  |    `host`    | `string` | The host of the celcat instance.                                                                          |
|  `continuous`  |    `loop`    |  `bool`  | Run in continuous mode. Will query the calendar periodicly according to the period defined in the config. |
|  `queryDelay`  |   `delay`    |  `int`   | Time in seconds between each query in daemon mode. Default : `1800`                                       |
|  `chromePath`  | `chromePath` | `string` | Path to the chrome executable. Default : `/usr/bin/chromium`                                              |
|     `png`      |    `png`     |  `bool`  | Enable PNG output (Require Chromium on your computer). Default : `false `                                 |
|   `pngPath`    |   `pngOut`   | `string` | Output directory for the PNG output. Default : `out/calendar/png/`                                        |
|   `pngWidth`   |   `width`    |  `int`   | Width of the PNG output. Default : `1920`                                                                 |
|   `pngHeigh`   |   `height`   |  `int`   | Height of the PNG output. Default : `1080`                                                                |
|     `html`     |    `html`    |  `bool`  | Enable HTML output. Default : `false`                                                                     |
| `htmlTemplate` |  `template`  | `string` | The template used to render the html page. Default : `web/templates/calendar.go.html`                     |
|   `htmlPath`   |  `htmlOut`   | `string` | Output directory for the HTML output. Default : `out/calendar/html/`                                      |
|     `ics`      |    `ics`     |  `bool`  | Enable ICS output. Default : `true`                                                                       |
|   `icsPath`    |   `icsOut`   | `string` | Output directory for the ICS output. Default : `out/calendar/ics/`                                        |
|     `web`      |    `web`     |  `bool`  | Enable the web server. Default : `false`                                                                  |
|   `webPort`    |    `port`    | `string` | Web listen Port. Default : `8080`                                                                         |


#### Add calendars to query from

To track new calendar you need to add them the group they belong to to the config file like so :

```json
{
  //...
  "groupes": [
    {
      "name": "Groupe1", // Name of the first calendar/group
      "id": "22014815" // Id of the calendar/group
    },
    //...
  ]
}
```

### Execution

Just output ICS :

```sh
go run ./cmd/cy-celcat/main.go -user=Someuser -pass=Pass
```

ICS+HTML+PNG in continous mode served with a web server
```sh
go run ./cmd/cy-celcat/main.go -html=1 -png=1 -web=1 -loop=1
```

### Docker

## Web Server

The web server will serve fill following that hierarchy:

- /
  - \*.ics // Ics file
  - \*.png // screenshot of the html page
  - \*     // Html calendar of the current week

Where `*` is the name if each of the calendars specified in the config
