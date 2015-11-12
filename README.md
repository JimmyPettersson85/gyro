# Gyro
A simple log-rotation library for Go.

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

## Rotation

How often the logs are rotated are defined by the element with the highest resolution in the layout string. By default Gyro rotates logs hourly since the default layout is `2006-01-02T15`. If you want to rotate logs daily the format would be `2006-01-02`, or for each minute (who would want to do that?) `2006-01-02T15:04`.

Gyro creates files lazily so if there are no calls to `Write`/`WriteString` no file for that time slot will be created.

## Concurrent writes

Writes to the log files are protected with a mutex to guard against concurrent writes to the same file.

## Tests

Tests and benchmarks performs writes in the `test` folder and cleans up after the tests/benchmarks are done.
