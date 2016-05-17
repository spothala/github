# github
Github API in GO

This Library supports most features in release,webhook functionalities. I'll be adding more _features_ soon.

## Installing

### _go get_
```
go get github.com/spothala/slack
```

## Example
```
import (
    "fmt"

    "github.com/spothala/github"
)

func main() {
    api := github.New("YOUR_TOKEN_HERE", true)
    // true turns-on debbuging, it will log all requests to the console
    // Useful when encountering issues
    release := api.CreateRelease(ORGREPO, RELNAME, TITLE, BODY)
    fmt.Println("Release Created with Name: "+relUrl["name"].(string)+" & Tag: "+relUrl["tag_name"].(string))
}
```

## Contribution
You guys are more than welcome to send pull requests if you see any errors/issues in my code.

## License
BSD 2 Clause license
