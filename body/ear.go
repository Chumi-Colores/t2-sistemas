package body

import (
	"main/global"
)

func Ear(who_to_listen_channel chan global.Command_data, neural_net global.Body_channels, last_listened []int, food_listened_too_early [][]global.Food) {
	defer global.Wg.Done()
	commands_left := true
	for commands_left {
		command, ok := <-who_to_listen_channel
		if !ok {
			commands_left = false
		} else {
			ear_command_handler(command, neural_net, last_listened, food_listened_too_early)
			neural_net.Finished_task_channel <- true
		}
	}
	close(neural_net.When_to_remember_channel)
}

func ear_command_handler(command global.Command_data, neural_net global.Body_channels, last_listened []int, food_listened_too_early [][]global.Food) {
	my_id := command.Participant_id
	expected_sender := command.Value
	when_to_remember_channel := neural_net.When_to_remember_channel

	real_food := global.Food{}
	expected_food_not_found := true
	for expected_food_not_found {
		if len(food_listened_too_early[expected_sender]) > 0 {
			real_food = food_listened_too_early[expected_sender][0]
			expected_food_not_found = false
			food_listened_too_early[expected_sender] = food_listened_too_early[expected_sender][1:]
		} else {
			current_food := <-global.Food_channels[my_id]
			if current_food.Sender == expected_sender {
				real_food = current_food
				expected_food_not_found = false
			} else {
				food_listened_too_early[current_food.Sender] = append(food_listened_too_early[current_food.Sender], current_food)
			}
		}
	}

	if real_food.Name == "garlic" {
		when_to_remember_channel <- global.Command_data{Participant_id: expected_sender, Action: 1, Value: 0}
		<-neural_net.Second_finished_task_channel
		global.Garlic_owner = my_id
		last_listened[expected_sender]++
	} else if real_food.Name == "candy" {
		last_listened[expected_sender]++
	} else if real_food.Name == "marker" {
		when_to_remember_channel <- global.Command_data{Participant_id: expected_sender, Action: 0, Value: real_food.Id}
		<-neural_net.Second_finished_task_channel
	}
}
