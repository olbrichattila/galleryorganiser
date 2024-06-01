# Gallery Organizer / Splitter

This small utility can copy your files and split them according to their type and date:
If you have sub folders they all goint to be processed and the below breakdown will be per sub folder
The application copies 10 fiels at the same time to be able to cope with large amount of files effectively.
If your computer can handle more, or cannot handle as much, change constant: ```paralellFileCount = 10``` in files go

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

## Parameters

### --src
The source file name

### --dst
The destination file name

### --overwrite
Overwrite if file exists, otherwise it can continue wher it was left off

### --flat
Flattern destination folder structure


## Usage:
```
filesplitter --src=./sourcedir --dst=./destinationdir
```

alternatively from golang:
```

go run .--src=./sourcedir --dst=./destinationdir
```

The utility will not overwrite, but skip existing file, so there is a fast way to continue where you left off, unless you switch overwrite on.
```
filesplitter --src=./sourcedir --dst=./destinationdir --overwrite
```
Flattern the result (ignore subdirectory names in destinatio folder)
```
filesplitter --src=./sourcedir --dst=./destinationdir --flat
```
- Note, using the overwrite and the flat together will keep the last file copied, without overwrite it keeps the firs one if the file names are matching
