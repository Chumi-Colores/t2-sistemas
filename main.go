package main

import (
	"fmt"
	"main/body"
	"main/global"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Por favor, proporciona una palabra como argumento.")
		return
	}

	filename := os.Args[1]
	actions, err := ConvertFileToIntLists(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	number_of_participants := actions[0][0]
	global.Buffer_size = actions[0][1]
	global.Garlic_owner = actions[0][2]

	make_participants_and_channels(number_of_participants)

	send_actions(actions)

	global.Wg.Wait()
}

func make_participants_and_channels(number_of_participants int) {
	global.Food_channels = make([]chan global.Food, number_of_participants)
	global.Command_channels = make([]chan global.Command_data, number_of_participants)
	global.Wg.Add(number_of_participants)
	for i := 0; i < number_of_participants; i++ {
		global.Food_channels[i] = make(chan global.Food, global.Buffer_size)
		global.Command_channels[i] = make(chan global.Command_data, global.Buffer_size)
		go participant(i)
	}
}

func create_snapshots_map() map[int]*global.Snapshot {
	var snapshots = make(map[int]*global.Snapshot)
	return snapshots
}

func participant(id int) {
	global.Wg.Add(4)
	defer global.Wg.Done()

	snapshots := create_snapshots_map()

	what_to_say_channel := make(chan global.Command_data)
	who_to_listen_channel := make(chan global.Command_data)
	when_to_remember_channel := make(chan global.Command_data)
	finished_task_channel := make(chan bool)
	second_finished_task_channel := make(chan bool)
	neural_net := global.Body_channels{
		What_to_say_channel:          what_to_say_channel,
		Who_to_listen_channel:        who_to_listen_channel,
		When_to_remember_channel:     when_to_remember_channel,
		Finished_task_channel:        finished_task_channel,
		Second_finished_task_channel: second_finished_task_channel,
	}

	last_said := make([]int, len(global.Command_channels))
	last_listened := make([]int, len(global.Command_channels))
	food_listened_too_early := make([][]global.Food, len(global.Command_channels))
	for i := range food_listened_too_early {
		food_listened_too_early[i] = make([]global.Food, 0)
	}

	go body.Mouth(what_to_say_channel, last_said, finished_task_channel)
	go body.Ear(who_to_listen_channel, neural_net, last_listened, food_listened_too_early)
	go body.Brain(global.Command_channels[id], neural_net)
	go body.Memory(when_to_remember_channel, neural_net, snapshots, last_said, last_listened, id)
}

func send_actions(actions [][]int) {
	for i := 1; i < len(actions); i++ {
		protagonist_id := actions[i][0]
		action := actions[i][1]
		secondary_character := actions[i][2]

		global.Command_channels[protagonist_id] <- global.Command_data{
			Participant_id: protagonist_id,
			Action:         action,
			Value:          secondary_character}
	}
	for i := range global.Command_channels {
		close(global.Command_channels[i])
	}
}
