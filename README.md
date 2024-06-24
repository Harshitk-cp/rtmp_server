# RTMP Server

## Overview

This project is an RTMP server inspired by LiveKit Ingress Service. It receives an RTMP stream from the user in a room, transcodes the audio from AAC to Opus (making it WebRTC compatible) and the video to H264, and then pushes it to WebRTC tracks connected to clients. The server acts as a peer, maintaining a P2P connection with each client.

## Features

- **RTMP to WebRTC:** Receives RTMP streams and delivers them to WebRTC clients.
- **Audio Transcoding:** Transcodes AAC audio to Opus for WebRTC compatibility.
- **Video Transcoding:** Ensures video is encoded in H264 for WebRTC delivery.
- **Webhook Notifications:** Uses webhooks to notify the publishing state of the stream to different rooms.
- **WebSocket Signaling:** Establishes WebRTC connections using WebSockets for offer/answer exchange.
- **Concurrency for Performance:** Utilizes Go's concurrency patterns and channels to enhance streaming performance and reduce latency.

## Core Libraries and Packages

- **Pion WebRTC:** Used for handling WebRTC connections.
- **Yuptopp RTMP:** Used for handling RTMP streams.
- **fdkaac:** For AAC decoding.
- **gopkg.in/hraban/opus.v2:** For Opus encoding.
- **go-chi/chi:** Lightweight, idiomatic, and composable router for building Go HTTP services.
- **logrus:** For logging.

## Installation

Clone the repository:

    ```sh
    git clone https://github.com/Harshitk-cp/rtmp_server.git
    cd rtmp_server




  ## How It Works
  
  ### RTMP Server
  
  The RTMP server listens for incoming RTMP streams. When a stream is published:
  
  - **Audio Processing:** Decodes AAC audio and encodes it into Opus format using `fdkaac` and `opus`.
  - **Video Processing:** Ensures the video stream is in H264 format.
  - **WebRTC Integration:** Sends processed audio and video to WebRTC tracks connected to clients.
  
  ### WebRTC Connection
  
  The WebRTC connection is established via WebSockets:
  
  - **WebSocket Handler:** Manages WebRTC signaling (offer/answer exchange) using WebSockets.
  - **Peer Connection:** Each client establishes a peer connection with the server.
  - **Track Delivery:** Delivers audio and video tracks to clients via WebRTC.
  
  ### Webhooks
  
  Webhooks lietens to the audio and video channel and notify the state of streams to it's subscribers:
  
  - **Notifications:** Sent when streams start or stop using a webhook manager.
  
  ## Performance Enhancements
  
  - **Concurrency:** Utilizes Go's concurrency patterns to efficiently handle multiple streams.
  - **Channels:** Uses channels for buffering video and audio data, ensuring smooth delivery to WebRTC tracks.
  - **Optimized Transcoding:** Efficiently transcodes audio and video to minimize latency.

