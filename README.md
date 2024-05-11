# AI Music Fun
This simple prototype act as a wrapper to [Replicate](https://replicate.com/)'s feature [MusicGen](https://replicate.com/meta/musicgen).
MusicGen allows you to create music from a description or from an existing song.

This API is a wrapper for MusicGen.

## Endpoints

Request: 
```
POST /song
{
  "input": {
    "prompt": "some cool description"
  }
}
```
Response
```
{
"message": "Your song is being processed, get the song with the link below",
    "song": "http://localhost:[port]/song/songID",
    "input": {
        "model_version": "stereo-large",
        "normalization_strategy": "peak",
        "output_format": "mp3",
        "prompt": "some cool description"
    },
    "status": "starting",
    "created_at": "2024-05-11T00:03:08.162Z"
}
```

In the response you will find a link that will be received by the API, and that link will redirect to the song (if the process was finished)

## Run locally
You can run this code locally. Node.js is mandatory and you must create an account and get a token that you must replace it in the code.


## Song example
I created a song that you can listen [here](https://soundcloud.com/bruno-giulianetti/replicate-prediction-test?si=c94bdcb3243a43d48559d6ef3167825a)
