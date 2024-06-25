package app

/*
Class for holding servent info data.
*/
type ServentInfo struct {
	Id        int
	Ip        string
	Port      int
	Neighbors []int
}

/*
Constructor for servent info.
*/
func ConstructServentInfo(id int, ip string, port int, neighbors []int) *ServentInfo {
	newServentInfo := ServentInfo{
		Id:        id,
		Ip:        ip,
		Port:      port,
		Neighbors: neighbors,
	}

	return &newServentInfo
}
