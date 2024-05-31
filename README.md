# Gallery organiser / splitter

This small utility can copy your files and split them according to their type and date:

For example:

```
-pictures
    -jpg
        -2020
            -01
                -img001.jpg
                -img005.jpg
                ...
            -02
                ...
            -03
                ...
        -2021
            ...
        -2022
            ...
    -mp4
        -2020
            -01
                -video001.mp4
                -video002.mp4
                ...
            -02
            ...
        -2021
        -2022

```

Usage:
```
filesplitter --src=./sourcedir --dst=./destinationdir

alternatively from golang:
go run .--src=./sourcedir --dst=./destinationdir

The utility will not overwrite, but skip existing file, so ther is a fast way to continue where you left off, unless you switch overwrite on
filesplitter --src=./sourcedir --dst=./destinationdir --overwrite

```

