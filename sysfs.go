package gpio

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

type direction uint

const (
	inDirection direction = iota
	outDirection
)

type edge uint

const (
	edgeNone edge = iota
	edgeRising
	edgeFalling
	edgeBoth
)

// Helper function to execute a shell command
func execCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}
	return nil
}

// exportGPIO uses 'gpioset' to configure a pin as an output initially
func exportGPIO(p Pin) {
	err := execCommand("gpiofind", []string{strconv.Itoa(int(p.Number))})
	if err != nil {
		fmt.Printf("failed to find gpio %d\n", p.Number)
		fmt.Println(err)
		return
	}
	// Small delay to ensure settings are applied
	time.Sleep(100 * time.Millisecond)
}

// setDirection sets the direction and initial value of the GPIO
func setDirection(p Pin, d direction, initialValue uint) {
	initial := "0"
	if initialValue != 0 {
		initial = "1"
	}
	args := []string{"gpiochip0", fmt.Sprintf("%d=%s", p.Number, initial)}
	if d == inDirection {
		args = []string{"gpiochip0", fmt.Sprintf("%d", p.Number)}
	}
	err := execCommand("gpioset", args)
	if err != nil {
		fmt.Printf("failed to set gpio %d direction\n", p.Number)
		fmt.Println(err)
	}
}

// setEdgeTrigger configures edge detection settings for GPIO interrupts
func setEdgeTrigger(p Pin, e edge) {
	edge := "none"
	switch e {
	case edgeRising:
		edge = "rising"
	case edgeFalling:
		edge = "falling"
	case edgeBoth:
		edge = "both"
	}
	err := execCommand("gpioget", []string{"gpiochip0", fmt.Sprintf("%d=%s", p.Number, edge)})
	if err != nil {
		fmt.Printf("failed to set gpio %d edge detection\n", p.Number)
		fmt.Println(err)
	}
}

func openPin(p Pin, write bool) Pin {
	// No file needs to be opened with the gpiod approach
	return p
}

func readPin(p Pin) (val uint, err error) {
	cmd := exec.Command("gpioget", "gpiochip0", strconv.Itoa(int(p.Number)))
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to read gpio: %w", err)
	}
	if len(output) > 0 && output[0] == '1' {
		return 1, nil
	}
	return 0, nil
}

func writePin(p Pin, v uint) error {
	value := "0"
	if v == 1 {
		value = "1"
	}
	return execCommand("gpioset", []string{"gpiochip0", fmt.Sprintf("%d=%s", p.Number, value)})
}
