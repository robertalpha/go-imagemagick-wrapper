# Go-imagemagick-wrapper
A simple helper wrapper to generate progressive web thumbnails using imagemagick within go. 

## Getting Started
Usage:
```
err, imagemagick := imwrapper.New()

if err != nil {
    t.Fatalf("could not instantiate imagemagick: %v", err)
}
quality := 95
maxDimension := 2048
autoRotate := true
imagemagick.Convert("input.jpg", "output.jpg", quality, maxDimension, autoRotate)
```

### Prerequisites
You need to have `imagemagick` installed.

## Running the tests
The tests are dockerized, to ensure they run with imagemagick `convert` and `identify` available.
```
docker run --rm -it $(docker build -q .)
```
coverage: `85.4% of statements`

### Test description
These tests make sure the commands to imagemagick `convert` and `identify` are working.

## Versioning
`v0` for now as I'll probably only be using it. 

## License
Distributed under the MIT License. See `LICENSE` for more information.

## Authors
* **Robert van Alphen** - [robertalpha](https://github.com/robertalpha)

## Acknowledgments
* Thanks to the creators of the excellent open source project `imagemagick`
