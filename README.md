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

### Pin usage layout


```
    |  LEFT | PCB | RIGHT |
    -----------------------
    |  D13  |     |  D12  |
    |  3V3  |     |  D11  |
    |  AREF |     |  D10  | - Insert Coin (Number 5)
    |  D14  |     |  D9   | - Player 1 Start (Number 1)
    |  D15  |     |  D8   | - Terminate (ESC)
    |  D16  |     |  D7   | - Launch (Enter)
    |  D17  |     |  D6   | 
    |  D18  |     |  D5   | - Right Magna Save (Right CTRL)
    |  D19  |     |  D4   | - Right Flipper (Right Shift)
    |  D20  |     |  D3   | - Left Magna Save (Left CTRL)
    |  D21  |     |  D2   | - Left Flipper (Left Shift)
    |  5V   |     | GND1  |
    | RESET |     | RESET |
    |  GND2 |     |  RX   |
    |  VIN  |     |  TX   |
```

In my own wiring I used right column ground, 'GND1', to connect the flipper and plunger buttons and left column 'GND2' to connect the remaining inputs - ESC, 1, 5.


### Editing the code

It is highly recommened to use tinygo-edit as it informs `gopls` about how to reference the custom libraries that tinygo overlays on top of go itself:

https://github.com/sago35/tinygo-edit

  go install github.com/sago35/tinygo-edit@latest

  tinygo-edit --editor vi --wait --target arduino-nano33

  tinygo-edit --editor nvim --wait --target arduino-nano33

then :e main.go'

#### Manual testing

While the device is plugged in you can execute `tinygo monitor` to see what might be logged to STDOUT.

#### Flashing the arduino-nano33

Flash your device with:

  tinygo flash -target=arduino-nano33

