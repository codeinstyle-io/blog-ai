package types

type Skill struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type SkillSection struct {
	Section string  `json:"section"`
	Skills  []Skill `json:"skills"`
}
