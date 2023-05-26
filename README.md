# Forklift

A simple utility that can take stdin and redirect it to templatized paths on S3.

## Installing

```shell
go install github.com/dacort/forklift/cmd/forklift@latest
```

Then pipe the sample file to a bucket!

```shell
curl -o - \
    "https://raw.githubusercontent.com/dacort/forklift/main/sample_data.json" \
    | forklift -w 's3://forklift-demo/{{json "event_type"}}/{{today}}.json'
```

## Overview

Usage is pretty simple - pipe some content to `forklift` and it will upload it to the desired S3 bucket and path.

```shell
echo "Hello Damon" | forklift -w s3://bucket/some/file.txt
```

While that in itself isn't too exciting (_you could just use `aws s3 cp -` !_), where it gets interesting is when you want to pipe JSON data and have it uploaded to a dynamic location based on the content of the data itself. For example, imagine a JSON file with the following content:

- `sample_data.json`
```json
{"event_type": "click", "data": {"uid": 1234, "path": "/signup"}}
{"event_type": "login", "data": {"uid": 1234, "referer": "yak.shave"}}
```

And imagine we want to pipe this to S3, but split it by `event_type`. Well, `forklift` can do that for us!

```shell
cat sample_data.json | forklift -w 's3://bucket/{{json "event_type"}}/{{today}}.json'
```

That will upload two different files:

1. `s3://bucket/click/2021-02-18.json`
2. `s3://bucket/login/2021-02-18.json`

### Default behavior

Note that the default behavior of `forklift` is to simply echo whatever is passed to it to stdout. This is partially because I build `forklift` into another project, as noted in the section below.

## Advanced Usage

Again, while not terribly interesting as a standalone CLI, where this becomes particularly useful is with `cargo-crates`. This is a sample project that makes it easy to captial-e Extract data from third-party services without having to be a data engineering wizard. 

For example, I've got an Oura ring and want to extract my sleep data. With the Oura Crate, I can simply do:

```shell
docker run -e OURA_PAT ghcr.io/dacort/crates-oura sleep
```

And that'll return a JSON blob with my sleep data for the past 7 days. But let's say I want to drop that sleep data into a location on S3 based on when I went to bed:

```shell
docker run -e OURA_PAT ghcr.io/dacort/crates-oura sleep | forklift  -w 's3://bucket/{{json "bedtime_start" | ymdFromTimestamp }}/sleep_data.json'
```

Cool. Now imagine I want to drop a single Docker container into an ETL workflow that does both of these for me. Well, `forklift` is integrated into Cargo Crates.

```shell
docker run \
    -e OURA_PAT \
    -e FORKLIFT_URI='s3://bucket/{{json "bedtime_start" | ymdFromTimestamp }}/sleep_data.json' \
    ghcr.io/dacort/crates-oura sleep
```

That will automatically take any stdout of the Docker container and pipe it to that location!

## Why?

This seems like a lot of work to just ... upload a file. Well, a few reasons.

1. I started [playing around](https://twitter.com/dacort/status/1359638593812140032) with the idea of Docker containers that could _very simply_ extract data from an API giving the consumer nothing else to worry about except having Docker and the proper authentication tokens. 
2. Then I wanted to upload the data to S3. But I wanted the Docker containers to remain as lightweight as possible. 
3. It's just a fun experiment. 🤷

## Resources

These resources came in handy while building this:
- [Linux pipes in Golang](https://dev.to/napicella/linux-pipes-in-golang-2e8j)
- [Using Goroutines and channels with AWS](https://maptiks.com/blog/using-go-routines-and-channels-with-aws-2/)
