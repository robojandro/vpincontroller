package main

import (
	"fmt"
	"time"

	"machine"
	"machine/usb/hid/keyboard"
	"tinygo.org/x/drivers/lsm6ds3"
)

const (
	// buttons - pins resolve to uint8's
	bLShift = machine.D2 // left flipper
	bLCtrl  = machine.D3 // left magna save, secondary flipper
	bRShift = machine.D4 // right flipper
	bRCtrl  = machine.D5 // right magna, secondary
	bLaunch = machine.D7 // plunger, launch

	// meta buttons
	bEsc  = machine.D8  // terminate execution of a table
	bOne  = machine.D9  // player 1 start
	bFive = machine.D10 // insert coin

	// nudge
	NUDGE_THRESHOLD   = int32(400_000)
	NUDGE_SKIP_BUDGET = 5
)

func main() {
	machine.I2C0.Configure(machine.I2CConfig{})
	accel := lsm6ds3.New(machine.I2C0)
	accelConf := lsm6ds3.Configuration{
		AccelRange:      lsm6ds3.ACCEL_2G,
		AccelSampleRate: lsm6ds3.ACCEL_SR_13,
	}
	if err := accel.Configure(accelConf); err != nil {
		for { // interesting sleeps a second while trying to connect
			fmt.Printf("failed to configure: %s\n", err.Error())
			time.Sleep(time.Second)
		}
	}

	var prevPos int32
	prevPos, err := readAccelerometer(accel)
	if err != nil {
		fmt.Print(err.Error())
	}
	kb := keyboard.Port()
	fmt.Printf("kb: %+v\n", kb)

	bCnf := machine.PinConfig{Mode: machine.PinInputPullup}
	bLShift.Configure(bCnf)
	bLCtrl.Configure(bCnf)
	bRShift.Configure(bCnf)
	bRCtrl.Configure(bCnf)
	bLaunch.Configure(bCnf)
	bEsc.Configure(bCnf)
	bOne.Configure(bCnf)
	bFive.Configure(bCnf)

	lShiftPress := false
	lCtrlPress := false
	rShiftPress := false
	rCtrlPress := false
	launchPress := false
	escPress := false
	onePress := false
	fivePress := false

	nudgeLeftPress := false
	nudgeRightPress := false

	skips := 0
	initialNudge := false
	for {
		// NUDGE
		y, err := readAccelerometer(accel)
		if err != nil {
			fmt.Print(err.Error())
		}
		nudge := detectNudge(prevPos, y, skips, initialNudge)
		skips = nudge.Skips - 1
		prevPos = y

		if nudge.Move == 1 { // tap on right side - positive
			nudgeRightPress = true
			fmt.Println("KeySlash")
			if err := kb.Down(keyboard.KeySlash); nil != err {
				fmt.Printf("error pressing forward slash: %s\n", err.Error())
			}
		} else if nudgeRightPress {
			kb.Up(keyboard.KeySlash)
			nudgeRightPress = false

		} else if nudge.Move == -1 { // tap on left side - negative
			nudgeLeftPress = true
			fmt.Println("KeyZ")
			if err := kb.Down(keyboard.KeyZ); nil != err {
				fmt.Printf("error pressing Z: %s\n", err.Error())
			}
		} else if nudgeLeftPress {
			kb.Up(keyboard.KeyZ)
			nudgeLeftPress = false
		}

		// BUTTONS
		// enter / plunger / launch
		if bLaunch.Get() {
			if launchPress {
				kb.Up(keyboard.KeyEnter)
			}
			launchPress = false
		} else if !launchPress {
			launchPress = true
			fmt.Println("KeyEnter")
			if err := kb.Down(keyboard.KeyEnter); nil != err {
				fmt.Printf("error pressing enter: %s\n", err.Error())
			}
		}

		// left flipper
		if bLShift.Get() {
			if lShiftPress {
				kb.Up(keyboard.KeyLeftShift)
			}
			lShiftPress = false
		} else if !lShiftPress {
			lShiftPress = true
			fmt.Println("KeyLeftShift")
			if err := kb.Down(keyboard.KeyLeftShift); nil != err {
				fmt.Printf("error pressing left shift: %s\n", err.Error())
			}
		}

		// left magna/secondary
		if bLCtrl.Get() {
			if lCtrlPress {
				kb.Up(keyboard.KeyLeftCtrl)
			}
			lCtrlPress = false
		} else if !lCtrlPress {
			lCtrlPress = true
			fmt.Println("KeyLeftCtrl")
			if err := kb.Down(keyboard.KeyLeftCtrl); nil != err {
				fmt.Printf("error pressing left ctrl: %s\n", err.Error())
			}
		}

		// right flipper
		if bRShift.Get() {
			if rShiftPress {
				kb.Up(keyboard.KeyRightShift)
			}
			rShiftPress = false
		} else if !rShiftPress {
			rShiftPress = true
			fmt.Println("KeyRightShift")
			if err := kb.Down(keyboard.KeyRightShift); nil != err {
				fmt.Printf("error pressing right shift: %s\n", err.Error())
			}
		}

		// right magna/secondary
		if bRCtrl.Get() {
			if rCtrlPress {
				kb.Up(keyboard.KeyRightCtrl)
			}
			rCtrlPress = false
		} else if !rCtrlPress {
			rCtrlPress = true
			fmt.Println("KeyRightCtrl")
			if err := kb.Down(keyboard.KeyRightCtrl); nil != err {
				fmt.Printf("error pressing right ctrl: %s\n", err.Error())
			}
		}

		// escape
		if bEsc.Get() {
			if escPress {
				kb.Up(keyboard.KeyEsc)
			}
			escPress = false
		} else if !escPress {
			escPress = true
			fmt.Println("KeyEsc")
			if err := kb.Down(keyboard.KeyEsc); nil != err {
				fmt.Printf("error pressing escape: %s\n", err.Error())
			}
		}

		// one
		if bOne.Get() {
			if onePress {
				kb.Up(keyboard.Key1)
			}
			onePress = false
		} else if !onePress {
			onePress = true
			fmt.Println("Key1")
			if err := kb.Down(keyboard.Key1); nil != err {
				fmt.Printf("error pressing 1: %s\n", err.Error())
			}
		}

		// five
		if bFive.Get() {
			if fivePress {
				kb.Up(keyboard.Key5)
			}
			fivePress = false
		} else if !fivePress {
			fivePress = true
			fmt.Println("Key5")
			if err := kb.Down(keyboard.Key5); nil != err {
				fmt.Printf("error pressing 5: %s\n", err.Error())
			}
		}

		// time.Sleep(time.Millisecond * 50)
		// try lowering the sleep cycle
		time.Sleep(time.Millisecond * 16)
	}
}

type Nudge struct {
	Move         int
	Skips        int
	InitialNudge bool
}

func readAccelerometer(accel *lsm6ds3.Device) (int32, error) {
	if !accel.Connected() {
		fmt.Println("LSM6DS3 not connected")
		time.Sleep(time.Second)
	}

	_, y, _, err := accel.ReadAcceleration()
	if err != nil {
		return 0, fmt.Errorf("failed to read acceleration: %s\n", err.Error())
	}

	return y, nil
}

func detectNudge(previous, y int32, skips int, initialNudge bool) Nudge {
	move := 0
	abs := y
	if abs < 0 {
		abs = abs * -1
	}

	absPrev := previous
	if absPrev < 0 {
		absPrev = absPrev * -1
	}
	diff := abs - absPrev

	absDiff := diff
	if diff < 0 {
		absDiff = -1 * diff
	}
	if skips <= 0 && !initialNudge {
		if absDiff > NUDGE_THRESHOLD {
			fmt.Printf("y %d   previous %d   diff %d   ", y, previous, absDiff)
			if y > previous {
				fmt.Printf("POSITIVE\n")
				move = 1
				skips = NUDGE_SKIP_BUDGET
				initialNudge = true
			} else if y < previous {
				fmt.Printf("NEGATIVE\n")
				move = -1
				skips = NUDGE_SKIP_BUDGET
				initialNudge = true
			}
		}
	} else {
		// fmt.Printf(" skips %d    initial: %t\n", skips, initialNudge)
		if skips == 0 {
			initialNudge = false
		}
	}

	return Nudge{
		Move:         move,
		Skips:        skips,
		InitialNudge: initialNudge,
	}
}
