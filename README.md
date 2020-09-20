# Mural
*Batch pixel-sorting images in Go*

## Helpful Commands

Use `convert` to transform JPG to PNG:

```
ls -1 *.jpeg | xargs -n 1 bash -c 'convert "$0" "${0%.jpeg}.png"'
```

Use `ffmpeg` to convert a video to individual frames, saved as pngs:

```
sudo apt-get install ffmpeg
ffmpeg -i input-video.mov frame-%04d.png -hide_banner
```