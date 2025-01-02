package pb

type LoginBody struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type PBEntry struct {
	Id      string `json:"id"`
	Updated string `json:"updated"`
}

type ListSearch struct {
	Items []PBEntry `json:"items"`
}