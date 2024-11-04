package body

import (
	"main/global"
)

func Mouth(what_to_say_channel chan global.Command_data, last_said []int, finished_task_channel chan bool) {
	defer global.Wg.Done()
	commands_left := true
	for commands_left {
		command, ok := <-what_to_say_channel
		if !ok {
			commands_left = false
		} else {
			mouth_command_handler(command, last_said)
			finished_task_channel <- true
		}
	}
}

func mouth_command_handler(command global.Command_data, last_said []int) {
	my_id := command.Participant_id
	food_receiver := command.Value

	if global.Garlic_owner == my_id {
		global.Garlic_owner = -1
		global.Food_channels[food_receiver] <- global.Food{Name: "garlic", Id: my_id}
	} else {
		global.Food_channels[food_receiver] <- global.Food{Name: "candy", Id: my_id}
	}
	last_said[food_receiver]++
}
