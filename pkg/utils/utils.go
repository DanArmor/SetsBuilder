package utils

type Rule struct {
	Number    int    `json:"-"`
	Left      string `json:"left"`
	Right     string `json:"right"`
	IsEpsylon bool
}

type RulesJSON struct {
	Rules []Rule `json:"rules"`
}

type RulesTable struct {
	Mapped      map[int]Rule
	Ordered     []Rule
	LeftMapped  map[string][]Rule
	RightMapped map[string][]Rule
	Size        int
}

// NewRulesTable представляет таблицу правил
func NewRulesTable() RulesTable {
	rt := RulesTable{
		Mapped:      make(map[int]Rule),      // Номера правил -> Правила
		Ordered:     make([]Rule, 0, 15),     // Массив в порядке добавления
		LeftMapped:  make(map[string][]Rule), // Левая часть -> правый чести
		RightMapped: make(map[string][]Rule), // Правая часть -> левый части
		Size:        0,                       // Размер
	}
	return rt
}

func (rt *RulesTable) AddRule(left string, right string) {
	rt.Size++
	rule := Rule{Number: rt.Size, Left: left, Right: right, IsEpsylon: false}
	if right == EpsilonSymbol {
		rule.IsEpsylon = true
	}
	rt.Mapped[rt.Size] = rule
	rt.Ordered = append(rt.Ordered, rule)
	rt.LeftMapped[left] = append(rt.LeftMapped[left], rule)
	rt.RightMapped[right] = append(rt.RightMapped[right], rule)
}

func (rt *RulesTable) CanBeEpsylon(str string) bool{
	for _, v := range rt.LeftMapped[str]{
		if v.IsEpsylon{
			return true
		}
	}
	return false
}

const (
	EndLineSymbol  = "⊣"
	EndStackSymbol = "∆"
	EpsilonSymbol  = "ε"
)
