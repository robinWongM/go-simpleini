# simpleini
![GitHub Actions status](https://img.shields.io/github/workflow/status/robinWongM/go-simpleini/Test)

Package simpleini provides a simple way to read and reload `.ini` configuration files.

## Usage

```bash
go get github.com/robinWongM/go-simpleini/simpleini
```

```golang
func reloadConf(conf simpleini.Configuration) {
  myNewValue := conf.Get("firstSection", "firstKey")
}

conf, err := simpleini.Watch("test.ini", reloadConf)
if err != nil {
  // your error handling
}

myValue := conf.Get("firstSection", "firstKey")
```

For full example, please refer to [examples/main.go](https://github.com/robinWongM/go-simpleini/blob/main/examples/main.go)

## License
This project is under MIT License. See the LICENSE file for the full license text.

This project includes code from [fsnotify](https://github.com/fsnotify/fsnotify). See the fsnotify.LICENSE file for corresponding license text.