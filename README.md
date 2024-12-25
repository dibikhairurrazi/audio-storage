# Audio Storage

## Description
This is an API that can be used to store and retrieve Audio files. API accept mp3 file and then the audio will be stored as .wav file.
This API use
- ffmpeg for converting audio
- mockgen to generate mock

## How to Use
There are two APIs that you can use.

### Store
`POST localhost:8080/audio/user/:user_id/phrase/:phrase_id`

This API is used to store file by provisioning user_id and phrase_id, if phrase_id already exist in the database, it will replace the existing entry.

### Retrieve
`GET localhost:8080/audio/user/:user_id/phrase/:phrase_id/:extension`

This API is used to retrieve stored audio file. Right now it only support mp3

## Running the Server

### Local
If you already have golang (1.22) and postgresql installed 
1. copy env `cp env.sample .env`
2. create `audio_storage` database on postgre-sql
3. run `go run app/api/main.go`

### Docker Compose
1. copy env `cp env.sample .env`
1. run `docker-compose build`
2. run `docker-compose up`

## Development
1. Create a new branch `git checkout -b "branch-name"`
2. Update the code
3. Run `make generate-mock` to update mock

## Improvement
- Implement User Login, Authentication, and Authorization
- Extend the interface to be able to store on other storage
- Cron to delete outdated file
- hash the ID and phraseID when saving the file
- more unit test
