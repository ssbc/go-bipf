# go-bipf

This a Go implementation of the [BIPF format][spec] based on
[`github.com/json-iterator/go`][jsoniter].

## Examples

### Marshal

    import (
        "encoding/hex"
        "fmt"
        "github.com/boreq/go-bipf"
    )

    func ExampleMarshal() {
        type ColorGroup struct {
            ID     int
            Name   string
            Colors []string
        }

        group := ColorGroup{
            ID:     1,
            Name:   "Reds",
            Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
        }

        b, err := bipf.Marshal(group)
        if err != nil {
            fmt.Println("error:", err)
        }

        fmt.Println(hex.EncodeToString(b))
        // Output:
        // 9d031049442201000000204e616d65205265647330436f6c6f7273c401384372696d736f6e185265642052756279304d61726f6f6e
    }

### Unmarshal

    import (
        "encoding/hex"
        "fmt"
        "github.com/boreq/go-bipf"
    )

    func ExampleUnmarshal() {
        bipfBlob, err := hex.DecodeString("9d031049442201000000204e616d65205265647330436f6c6f7273c401384372696d736f6e185265642052756279304d61726f6f6e")
        if err != nil {
            fmt.Println("error:", err)
        }

        type ColorGroup struct {
            ID     int
            Name   string
            Colors []string
        }

        var group ColorGroup

        err = bipf.Unmarshal(bipfBlob, &group)
        if err != nil {
            fmt.Println("error:", err)
        }
        fmt.Printf("%+v", group)
        // Output:
        // {ID:1 Name:Reds Colors:[Crimson Red Ruby Maroon]}
    }

[spec]: https://github.com/ssbc/bipf-spec
[jsoniter]: github.com/json-iterator/go
