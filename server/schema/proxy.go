package schema

type ProxyDetails struct {
	Id int64 `json:"id"`
	UserId string `json:"userId"`
	Proto string `json:"proto"`
	Host string `json:"host"`
	Port string `json:"port"`
	Name string `json:"name"`
	Password string `json:"password"`
	IsInUse bool `json:"isInUse"`
	IsEnabled bool `json:"isEnabled"`
}
