package command

import (
	"Snapshot/app"
	"fmt"
	"strconv"
	"strings"
)

type Info_command struct{}

func (cmd Info_command) CommandName() string {
	return "info"
}

func (cmd Info_command) Execute(args string) {
	// Format the servent info.
	formatedPrint := fmt.Sprintf("My info:\n\tId: %d\n\tIp: %s\n\tPort: %d\n\tNeighbors: %s", app.MyServentInfo.Id, app.MyServentInfo.Ip, app.MyServentInfo.Port, sliceToString(app.MyServentInfo.Neighbors))

	// Print the formatted info message.
	fmt.Println(formatedPrint)
}

func sliceToString(slice []int) string {
	// Convert each integer to string and join them with ", ".
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = strconv.Itoa(v)
	}
	return strings.Join(strSlice, ", ")
}
