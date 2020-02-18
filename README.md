# streamscreen

Save screen to image, video, network, DLNA.

# environment

- D=n - display number, default 0 (means first).
- X=x - start from x when saving screen (counts from top left corner to the right), default 0.
- Y=y - start from y when saving screen (counts from top left corner down), default 0.
- W=w - width, default - selected screen width, max = screen width - x.
- H=h - height, default - selected screen height, max = screen height - y.
- SS=1 - save single screenshot and exit.
- SV=1 - save video until CTRL+C until signal (CTRL^C), saves as much as possible screenshots and creates video from them on CTRL^C handling.
- F=n - in SV mode, set number of frames to generate, default 0 which means infinite.
- FPS=fps - set FPS, it will generate no more than `fps` frames per second, default 0 which means unlimited.
