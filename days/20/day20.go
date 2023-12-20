package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const DAY = "20"

type ModuleType int

const (
	Broadcaster ModuleType = iota + 1
	FlipFlop
	Conjunction
)

type Pulse struct {
	source string
	target string
	isHigh bool
}

type Simulation struct {
	modules              map[string]Module
	pulses               []Pulse
	highPulsesSent       int
	lowPulsesSent        int
	conjunctionsFiredLow []string
}

func (simulation Simulation) String() string {

	return fmt.Sprintf("(modules=%#v)", simulation.modules)
}

func (module Module) String() string {

	return fmt.Sprintf("(input=%s, lastPulseREceivedByInput=%#v, output=%s, type=%d, isOn=%t, label=%s)", module.input, module.lastPulseReceivedByInput, module.output, module.moduleType, module.isOn, module.label)
}

func (pulse Pulse) String() string {

	return fmt.Sprintf("(source=%s, target=%s, isHigh=%t)", pulse.source, pulse.target, pulse.isHigh)
}

type Module struct {
	input                    []string
	lastPulseReceivedByInput map[string]Pulse
	output                   []string
	moduleType               ModuleType
	isOn                     bool
	label                    string
}

func parseSimulation(input string) Simulation {

	var simulation Simulation

	simulation.modules = make(map[string]Module)

	for _, line := range strings.Split(strings.ReplaceAll(input, "\r", ""), "\n") {

		module := parseModule(line)

		_, moduleExisted := simulation.modules[module.label]

		if moduleExisted {

			simulation.modules[module.label] = Module{
				input:                    simulation.modules[module.label].input,
				lastPulseReceivedByInput: simulation.modules[module.label].lastPulseReceivedByInput,
				output:                   module.output,
				moduleType:               module.moduleType,
				isOn:                     false,
				label:                    module.label,
			}

		} else {

			simulation.modules[module.label] = module
		}

		for _, outputModuleLabel := range module.output {

			_, outputModuleExisted := simulation.modules[outputModuleLabel]

			if outputModuleExisted {

				simulation.modules[outputModuleLabel] = Module{
					input: append(simulation.modules[outputModuleLabel].input, module.label),
					lastPulseReceivedByInput: cloneMapAndInsert(simulation.modules[outputModuleLabel].lastPulseReceivedByInput, module.label, Pulse{
						source: module.label,
						target: outputModuleLabel,
						isHigh: false,
					}),
					output:     simulation.modules[outputModuleLabel].output,
					moduleType: simulation.modules[outputModuleLabel].moduleType,
					isOn:       simulation.modules[outputModuleLabel].isOn,
					label:      simulation.modules[outputModuleLabel].label,
				}

			} else {

				simulation.modules[outputModuleLabel] = Module{
					input: []string{module.label},
					lastPulseReceivedByInput: cloneMapAndInsert(make(map[string]Pulse), module.label, Pulse{
						source: module.label,
						target: outputModuleLabel,
						isHigh: false,
					}),
					output:     []string{},
					moduleType: 0,
					isOn:       false,
					label:      outputModuleLabel,
				}
			}
		}
	}

	return simulation
}

func cloneModule(module Module) Module {

	return Module{
		input:                    module.input,
		lastPulseReceivedByInput: cloneMap(module.lastPulseReceivedByInput),
		output:                   module.output,
		moduleType:               module.moduleType,
		isOn:                     module.isOn,
		label:                    module.label,
	}
}

func cloneMapAndInsert(mapping map[string]Pulse, key string, value Pulse) map[string]Pulse {

	result := cloneMap(mapping)
	result[key] = value
	return result
}

func cloneMap(mapping map[string]Pulse) map[string]Pulse {

	var result = make(map[string]Pulse)

	for key, value := range mapping {

		result[key] = value
	}

	return result
}

func parseModule(line string) Module {

	chunks := strings.Split(line, " -> ")

	var module Module

	if strings.Contains(chunks[0], "broadcaster") {

		module.moduleType = Broadcaster

	} else if strings.Contains(chunks[0], "%") {

		module.moduleType = FlipFlop

	} else if strings.Contains(chunks[0], "&") {

		module.moduleType = Conjunction

	}

	module.output = strings.Split(strings.ReplaceAll(chunks[1], " ", ""), ",")

	module.isOn = false
	module.label = strings.ReplaceAll(strings.ReplaceAll(chunks[0], "%", ""), "&", "")

	return module
}

func processPulses(simulation *Simulation) {

	simulation.highPulsesSent = 0
	simulation.lowPulsesSent = 0

	simulation.conjunctionsFiredLow = []string{}

	simulation.pulses = append(simulation.pulses, Pulse{
		target: "broadcaster",
		isHigh: false,
	})

	for len(simulation.pulses) > 0 {

		currentPulse := simulation.pulses[0]
		simulation.pulses = simulation.pulses[1:]

		processPulse(simulation, currentPulse)
	}
}

func processPulse(simulation *Simulation, pulse Pulse) {

	if pulse.isHigh {
		simulation.highPulsesSent += 1
	} else {
		simulation.lowPulsesSent += 1
	}

	module := cloneModule(simulation.modules[pulse.target])

	switch module.moduleType {
	case Broadcaster:
		for _, outputModuleLabel := range module.output {

			simulation.pulses = append(simulation.pulses, Pulse{
				source: module.label,
				target: outputModuleLabel,
				isHigh: pulse.isHigh,
			})
		}
	case FlipFlop:
		if !pulse.isHigh {

			if module.isOn {

				module.isOn = false
				for _, outputModuleLabel := range module.output {

					simulation.pulses = append(simulation.pulses, Pulse{
						source: module.label,
						target: outputModuleLabel,
						isHigh: false,
					})
				}
			} else {

				module.isOn = true
				for _, outputModuleLabel := range module.output {

					simulation.pulses = append(simulation.pulses, Pulse{
						source: module.label,
						target: outputModuleLabel,
						isHigh: true,
					})
				}
			}

			simulation.modules[module.label] = module
		}
	case Conjunction:

		module.lastPulseReceivedByInput[pulse.source] = Pulse{
			source: pulse.source,
			target: pulse.target,
			isHigh: pulse.isHigh,
		}

		sendHighPulse := false

		for _, pulse := range module.lastPulseReceivedByInput {

			if !pulse.isHigh {

				sendHighPulse = true
				break
			}
		}

		if !sendHighPulse {

			simulation.conjunctionsFiredLow = append(simulation.conjunctionsFiredLow, module.label)
		}

		for _, outputModuleLabel := range module.output {

			simulation.pulses = append(simulation.pulses, Pulse{
				source: module.label,
				target: outputModuleLabel,
				isHigh: sendHighPulse,
			})
		}

		simulation.modules[module.label] = module
	}
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

func LCM(numbers []int) int {
	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		result = lcm(result, numbers[i])
	}
	return result
}

func Part1(input string) string {

	content := GetContent(input)

	simulation := parseSimulation(content)

	limit := 1000

	highPulsesSent := 0
	lowPulsesSent := 0
	runsCompleted := 0

	for runsCompleted < limit {

		processPulses(&simulation)

		highPulsesSent += simulation.highPulsesSent
		lowPulsesSent += simulation.lowPulsesSent
		runsCompleted += 1
	}

	return strconv.Itoa(lowPulsesSent * highPulsesSent)
}

func Part2(input string) string {

	content := GetContent(input)

	// hardcoded required conjunctions that must fire a low pulse simultaneously based on input
	conjunctionsFiringLow := []string{"dc", "qm", "jh", "zq"}

	var runsUntilOn = make(map[string]int)

	simulation := parseSimulation(content)

	runsCompleted := 0

	for len(runsUntilOn) < len(conjunctionsFiringLow) {

		processPulses(&simulation)

		runsCompleted += 1

		for _, conjunctionRequired := range conjunctionsFiringLow {

			if _, alreadyFound := runsUntilOn[conjunctionRequired]; !alreadyFound {

				if slices.Contains(simulation.conjunctionsFiredLow, conjunctionRequired) {

					runsUntilOn[conjunctionRequired] = runsCompleted
				}
			}
		}
	}

	var runs []int

	for _, runCount := range runsUntilOn {

		runs = append(runs, runCount)
	}

	// calculate the run count until all conjunctions fire in one run
	return strconv.Itoa(LCM(runs))
}

func GetContent(filepath string) string {

	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY))))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
