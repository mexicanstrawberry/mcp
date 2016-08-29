package events

const (
	P_MAX Priority = iota
	P_NORM
	P_MIN
)

var priorityMap = map[Priority]string{
	P_MAX:  "Highest priority",
	P_NORM: "Normal priority",
	P_MIN:  "Lowest priority",
}

var Channel *Event

type Event chan interface{}

type MqttRecive struct {
	Map map[string]interface{}
}

type Ctl map[string]int

type Priority int

func (p *Priority) String() string {

	if s, ok := priorityMap[*p]; ok {
		return s
	}
	return "unknown priority"

}

type MqttCommand struct {
	CommandID uint8
	Actuator  string
	Value     float64
	Priority  Priority
}

type MqttCommandAck struct {
	CommandID uint8
	Accepted  bool
}
