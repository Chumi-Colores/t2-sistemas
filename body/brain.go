package body

import (
	"main/global"
	"time"
)

func Brain(command_channel chan global.Command_data, neural_net global.Body_channels) {
	defer global.Wg.Done()
	commands_left := true
	for commands_left {
		command, ok := <-command_channel
		if !ok {
			commands_left = false
		} else {
			brain_command_handler(command, neural_net)
		}
	}
	close(neural_net.What_to_say_channel)
	close(neural_net.Who_to_listen_channel)
}

func brain_command_handler(command global.Command_data, neural_net global.Body_channels) {
	what_to_say_channel := neural_net.What_to_say_channel
	who_to_listen_channel := neural_net.Who_to_listen_channel
	action := command.Action
	switch action {
	case 0: // 0: SEND
		what_to_say_channel <- command
		<-neural_net.Finished_task_channel
	case 1: // 1: RECEIVE
		who_to_listen_channel <- command
		<-neural_net.Finished_task_channel
	case 2: // 2: WAIT
		seconds := command.Value
		wait_seconds(seconds)
	case 3: // 3: SNAPSHOT
		snapshot_id := command.Value
		for i := range global.Command_channels {
			if i != command.Participant_id {
				global.Food_channels[i] <- global.Food{Name: "marker", Id: snapshot_id, Sender: command.Participant_id}
			}
		}
	}
}

func wait_seconds(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
