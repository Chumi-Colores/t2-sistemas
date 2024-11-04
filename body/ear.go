package body

import (
	"fmt"
	"main/global"
)

func Ear(who_to_listen_channel chan global.Command_data, neural_net global.Body_channels, last_listened []int) {
	defer global.Wg.Done()
	commands_left := true
	for commands_left {
		command, ok := <-who_to_listen_channel
		if !ok {
			commands_left = false
		} else {
			ear_command_handler(command, neural_net, last_listened)
			neural_net.Finished_task_channel <- true
		}
	}
	close(neural_net.When_to_remember_channel)
}

func ear_command_handler(command global.Command_data, neural_net global.Body_channels, last_listened []int) {
	my_id := command.Participant_id
	food_sender := command.Value
	when_to_remember_channel := neural_net.When_to_remember_channel

	food := <-global.Food_channels[my_id]
	if food.Name == "garlic" {
		when_to_remember_channel <- global.Command_data{Participant_id: food_sender, Action: 1, Value: 0}
		<-neural_net.Second_finished_task_channel
		global.Garlic_owner = my_id
		last_listened[food_sender]++
	} else if food.Name == "candy" {
		last_listened[food_sender]++
	} else if food.Name == "marker" {
		fmt.Println(my_id, "received marker:", global.Command_data{Participant_id: food_sender, Action: 0, Value: food.Id})
		when_to_remember_channel <- global.Command_data{Participant_id: food_sender, Action: 0, Value: food.Id}
		<-neural_net.Second_finished_task_channel
	}
}
