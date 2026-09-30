# Play List Maker

Local music libary scanner to create and update playlist files

## Motivation

I created this to make it easier for me to manage my local music library playlists, update them when I add music to my digital library.

## Functionalities

Tool is built to run on linux system.
Application configurations are stored in plm_conf.json file in users profile .confg/playlistmaker/ directory. Parameters expect full paths to the directories.

Extracts IDv3 metadata from MP3 files. When information is missing then file name na directory path information is used as an alternative although it can cause faulty information due to data quality issues.

Playlyst output is written into .xspf xml formatted playlist file.
