package base

// Location : 标示IP所在地
type Location string

// China :
const China Location = "CN"

// Japan :
const Japan Location = "JP"

// Others :
const Others Location = "Others"

func (l Location) String() string {
	return string(l)
}
