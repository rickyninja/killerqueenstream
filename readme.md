# Streaming setup

## cab setup
* Cab PC runs ffmpeg commands for each cab-cam, and streams over ethernet to PC running crtmpserver (or similar service, possibly even OBS).
* crtmpserver runs on laptop or cab PC, yet undecided.
* This is all assuming it won't use too many resources on cab PC.

## laptop setup
* The capture card will pass along cab screen capture.
* Can OBS take stream input directly from ffmpeg over ethernet?
* Can OBS handle overlaying the cab cam feeds onto the screen capture?

## Diagram
<img src="images/stream-diagram.png" />
