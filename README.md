# Gyro
A simple log-rotation library for persistent log files in Go.

The reason for creating Gyro was to provide file-logging where you could easily have a list of persistent log files divided up in equally big time frames, e.g

```
myproject_2015-10-09T00.log
myproject_2015-10-09T01.log
myproject_2015-10-09T02.log
myproject_2015-10-09T03.log
myproject_2015-10-09T04.log
myproject_2015-10-09T05.log
myproject_2015-10-09T06.log
...
myproject_2015-10-09T23.log
myproject_2015-10-10T00.log
```

## Installation
`go get github.com/slimmy/gyro`

## Quick start

```go
logger, _ := gyro.New("/tmp")
logger.Write([]byte("Writing some data\n")
logger.WriteString("Some more data\n")
```

## Defaults

By default Gyro writes logs in the format of `YYYY-MM-DDTh.log` with the current time in UTC.

## Standard library integration

Gyro implementents the `io.Writer` interface so it can be used with the stdlib logger by calling `log.SetOutput(*gyro.Logger)`.
Gyro was initially implemented to be used as a stand alone file-logger though, for logging data that needed to persist indefinately (hence there are no deleting of files when it rotates).

## Configuration

Gyro supports a range of configuration options. All examples given below assume the time and date is `1970-01-01 00:00:00`

### Field separator
```go
logger.SetSeparator("_")
```
Sets the separator between the prefix/suffix and the timestamp string

### Prefix/suffix

Adding a prefix and/or suffix to the generated filename is done by calling `SetPrefix`/`SetSuffix`.
```go
logger, _ := gyro.New("/tmp")
logger.FileName() // -> 1970-01-01T00.log (default)

logger.SetSeparator("_")

logger.SetPrefix("pre")
logger.FileName() // -> pre_1970-01-01T00.log

logger.SetSuffix("suf")
logger.FileName() // -> pre_1970-01-01T00_suf.log
```

### Extension
```go
logger.SetExtension("txt")
logger.FileName() // -> 1970-01-01T00.txt
```

### Layout
```go
logger.SetLayout("2006.01.02_15")
logger.FileName() // -> 1970.01.01_00.log
```
Sets the layout format that will be fed in to `time.Format`

### Time function
```go
logger.SetTimeFunction(func() time.Time {
    return time.Now().UTC() //default
})
```
Set the function that Gyro uses to get the time for the filename. If you want to rotate the logs in your local time just use `time.Now()`.

If you want to log in another location it can be done like so
```go
logger.SetTimeFunction(func() time.Time {
    loc, _ := time.LoadLocation("US/Alaska")
    return time.Now().In(loc)
})
```
or shifting the time
```go
logger.SetTimeFunction(func() time.Time {
    return time.Now().Add(-2 * time.Hour)
})
```

### Full example
Using the above to set up a logger
```go
    logger, _ := gyro.New("/tmp")
    logger.SetSeparator("_")
    logger.SetPrefix("prefix")
    logger.SetSuffix("suffix")
    logger.SetExtension("txt")
    logger.SetTimeFunction(func() time.Time {
        loc, _ := time.LoadLocation("Europe/Stockholm")
        return time.Now().In(loc)
    })
    fmt.Println(logger)
```
prints
```
gyro.Logger
  path: /tmp
  prefix: "prefix"
  suffix: "suffix"
  separator: "_"
  extension: "txt"
  layout: 2006-01-02T15
  format: prefix_%s_suffix.txt
  current filename: prefix_2015-11-12T18_suffix.txt
```

## Rotation

How often the logs are rotated are defined by the element with the highest resolution in the layout string. By default Gyro rotates logs hourly since the default layout is `2006-01-02T15`. If you want to rotate logs daily the format would be `2006-01-02`, or for each minute (who would want to do that?) `2006-01-02T15:04`.

Gyro creates files lazily so if there are no calls to `Write`/`WriteString` no file for that time slot will be created.

## Concurrent writes

Writes to the log files are protected with a mutex to guard against concurrent writes to the same file.

## Tests

Tests and benchmarks performs writes in the `test` folder and cleans up after the tests/benchmarks are done.
