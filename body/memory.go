package body

import (
	"fmt"
	"main/global"
	"os"
	"strconv"
)

func Memory(when_to_remember_channel chan global.Command_data, neural_net global.Body_channels,
	snapshots map[int]*global.Snapshot,
	last_sent []int,
	last_listened []int,
	my_id int) {
	defer global.Wg.Done()

	commands_left := true
	for commands_left {
		command, ok := <-when_to_remember_channel
		if !ok {
			commands_left = false
		} else {
			memory_command_handler(command, neural_net, snapshots, last_sent, last_listened, my_id)
		}
	}
}

func memory_command_handler(message global.Command_data, neural_net global.Body_channels, snapshots map[int]*global.Snapshot, last_sent []int, last_listened []int, my_id int) {
	message_sender := message.Participant_id
	candy_or_garlic_or_marker := message.Action // 1: garlic, 0: marker, -1: candy

	if candy_or_garlic_or_marker == -1 { // this message came because of a candy
		last_listened[message_sender]++
	} else if candy_or_garlic_or_marker == 1 { // this message came because of a garlic
		last_listened[message_sender]++
		for _, snapshot := range snapshots {
			if !snapshot.Ajo_allready_received {
				snapshot.Ajo_at_marker = [2]int{1, message_sender}
				snapshot.Ajo_allready_received = true
			}
		}
		neural_net.Second_finished_task_channel <- true
	} else { // this message came because of a marker
		snapshot_id := message.Value
		if _, ok := snapshots[snapshot_id]; !ok { // new snapshot_id
			add_new_snapshot(snapshots, snapshot_id, last_sent, last_listened, my_id)
			// mark this participant as received
			snapshots[snapshot_id].Markers_received[message_sender] = true
			for i := range global.Food_channels {
				if i != my_id {
					global.Food_channels[i] <- global.Food{Name: "marker", Id: snapshot_id, Sender: my_id}
				}
			}
			neural_net.Second_finished_task_channel <- true
		} else {
			if snapshots[snapshot_id].Number_of_markers_received == len(global.Command_channels) {
				neural_net.Second_finished_task_channel <- true
				return
			}
			snapshot_is_full := update_channel_record(snapshots, snapshot_id, message_sender, last_listened)
			if snapshot_is_full {
				global.Wg.Add(1)
				go write_and_delete_snapshot(snapshots, snapshot_id, my_id)
			}
			neural_net.Second_finished_task_channel <- true
		}
	}
}

func add_new_snapshot(snapshots map[int]*global.Snapshot, snapshot_id int, last_sent []int, last_listened []int, my_id int) {
	stateCopy := make([]int, len(last_sent))
	copy(stateCopy, last_sent)

	messageCopy := make([]int, len(last_listened))
	copy(messageCopy, last_listened)

	newSnapshot := &global.Snapshot{
		State_at_record:            stateCopy,
		Message_at_record:          messageCopy,
		Ajo_at_record:              my_id == global.Garlic_owner,
		Channel_at_record:          make([]int, len(global.Command_channels)),
		Ajo_at_marker:              [2]int{-1, -1},
		Markers_received:           make([]bool, len(global.Command_channels)),
		Number_of_markers_received: 1,
	}

	snapshots[snapshot_id] = newSnapshot
}

func update_channel_record(snapshots map[int]*global.Snapshot, snapshot_id int, marker_sender int, last_listened []int) (full bool) {
	snapshot := snapshots[snapshot_id]

	first_marker_from_this_participant := !snapshot.Markers_received[marker_sender]
	if first_marker_from_this_participant {
		channelRecord := snapshot.Channel_at_record
		messageRecord := snapshot.Message_at_record
		sentDifference := last_listened[marker_sender] - messageRecord[marker_sender]
		channelRecord[marker_sender] = sentDifference

		snapshot.Markers_received[marker_sender] = true
		snapshot.Number_of_markers_received++
	}

	maximum_markers := len(global.Command_channels) - 1
	snapshot_full_of_markers := (snapshot.Number_of_markers_received == maximum_markers)

	return snapshot_full_of_markers
}

func write_and_delete_snapshot(snapshots map[int]*global.Snapshot, snapshot_id int, my_id int) {
	defer global.Wg.Done()

	copy_snapshot_to(snapshots, snapshot_id, my_id)

	if len(global.Final_snapshots[snapshot_id]) == len(global.Command_channels) {
		write_snapshot_to_file(snapshot_id)
	}
}

func copy_snapshot_to(snapshots map[int]*global.Snapshot, snapshot_id int, my_id int) {
	snap := snapshots[snapshot_id]

	copySnapshot := global.Snapshot_to_write{
		State_at_record:   make([]int, len(snap.State_at_record)),
		Message_at_record: make([]int, len(snap.Message_at_record)),
		Ajo_at_record:     snap.Ajo_at_record,
		Channel_at_record: make([]int, len(snap.Channel_at_record)),
		Ajo_at_marker:     snap.Ajo_at_marker,
	}

	copy(copySnapshot.State_at_record, snap.State_at_record)
	copy(copySnapshot.Message_at_record, snap.Message_at_record)
	copy(copySnapshot.Channel_at_record, snap.Channel_at_record)

	global.Mu.Lock()
	if _, exists := global.Final_snapshots[snapshot_id]; !exists {
		global.Final_snapshots[snapshot_id] = make(map[int]global.Snapshot_to_write)
	}
	global.Final_snapshots[snapshot_id][my_id] = copySnapshot
	global.Mu.Unlock()
}

func write_snapshot_to_file(snapshot_id int) error {
	file, err := os.Create(global.FileName + strconv.Itoa(snapshot_id) + ".txt")
	if err != nil {
		return err
	}
	defer file.Close()

	for i := range len(global.Final_snapshots[snapshot_id]) {
		snap := global.Final_snapshots[snapshot_id][i]
		fmt.Fprintf(file, "%d:\n", i)
		fmt.Fprintf(file, "stateatRecord: %v\n", snap.State_at_record)
		fmt.Fprintf(file, "messageatRecord: %v\n", snap.Message_at_record)
		fmt.Fprintf(file, "ajoatRecord: %v\n", snap.Ajo_at_record)
		fmt.Fprintf(file, "channelatRecord: %v\n", snap.Channel_at_record)
		true_or_false := "false"
		if snap.Ajo_at_marker[0] == 1 {
			true_or_false = "true"
		}
		if i < len(global.Command_channels)-1 {
			fmt.Fprintf(file, "ajoatMarker: [%v %v]\n", true_or_false, snap.Ajo_at_marker[1])
			fmt.Fprintln(file) // Extra line
		} else {
			fmt.Fprintf(file, "ajoatMarker: [%v %v]", true_or_false, snap.Ajo_at_marker[1])
		}
	}
	return nil
}
