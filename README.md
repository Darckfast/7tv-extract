![7tv-extract-icon](icon.ico)

# 7tv extract

Given a user id, downloads all emotes locally, converting them to GIFs and PNGs

```
.\7tv-extract <user-id>
```

All emotes will be compress to 64x64 (the ratio will be maintained) and GIFs will be reduce by half of frames

## Requirements

The auto-conversion requires ImageMagick to be pre-installed https://imagemagick.org/script/download.php
For higher GIF compression first install https://www.lcdf.org/gifsicle/
