# Mural
*Batch pixel-sorting images in Go*

## Converting Videos to Individual Frames
Use `ffmpeg` to convert a video to individual frames, saved as pngs:

    sudo apt-get install ffmpeg
    ffmpeg -i input-video.mov frame-%04d.png -hide_banner
