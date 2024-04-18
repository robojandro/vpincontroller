# Overview

This is a TinyGo application written for the Arduino Nano 33 IoT microcontroller using the USB HID keyboard libraries.

It implements flippers, magna save buttons, and a plunger.

It also does a lazy attempt at supporting nudge detection.

## Nudging

Nudging is implemented as keypressed of the following keys:

```
    z - nudge table to the right
        like shoving the table on it's left side

    / - nudge table to the right
        like shoving the table on it's left side
```

The algorithm for the detection comes down to comparing accelerometer readings over time filtering out small movements.

I took a stab at trying to determine a base position by average the last 10 or so readings but this proved tricky and supprisingly the more simple approach yielded better results.

It is probably about 70% accurate now, which is almost good enough as I was simply interested in altering the trajectory of a ball headed into the drain and not trying to pull off tricks.

## TODO

- Implement a couple of more buttons to support 'insert coins' and the 'player 1 start'.
- Improve the nudge detection
