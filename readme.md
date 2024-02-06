# dog (dump out graphics) 🐶

A friendly, `cat`-like command to view images in your terminal.

The main goal with this tool is not to render images with prefect accuracy, but to be able to
somewhat get a quick and dirty preview of what an image looks like when you're on the go.

## Installing
```shell
go install github.com/tomasruud/dog@latest
```

## Usage
```shell
# You can use it with local files
dog <image file>

# Or pipe something from a remote url
curl -s <image url> | dog

# Cats and dogs play together just fine
cat <image file> | dog
```

## Similar tools
Inspired by [imgcat by danielgatis](https://github.com/danielgatis/imgcat).
