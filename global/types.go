package global

type Command_data struct {
	Participant_id int
	Action         int
	Value          int
}

type Body_channels struct {
	What_to_say_channel          chan Command_data
	Who_to_listen_channel        chan Command_data
	When_to_remember_channel     chan Command_data
	Finished_task_channel        chan bool
	Second_finished_task_channel chan bool
}

type Snapshot struct {
	State_at_record            []int
	Message_at_record          []int
	Ajo_at_record              bool
	Channel_at_record          []int
	Ajo_at_marker              [2]int
	Ajo_allready_received      bool
	Markers_received           []bool
	Number_of_markers_received int
}

type Snapshot_to_write struct {
	State_at_record   []int
	Message_at_record []int
	Ajo_at_record     bool
	Channel_at_record []int
	Ajo_at_marker     [2]int
}

type Food struct {
	Name   string
	Id     int
	Sender int
}
