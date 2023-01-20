package sets

import (
	"github.com/DanArmor/taEkzUtils/pkg/utils"
	"github.com/emirpasic/gods/sets/linkedhashset"
	"unicode"
	"unicode/utf8"
)

type Result struct {
	NonTerminal string
	First       *linkedhashset.Set
	Forward     *linkedhashset.Set
}

type Choose struct {
	Rule utils.Rule
	Cs   *linkedhashset.Set
}

func NewChoose(rule utils.Rule) Choose {
	c := Choose{
		Rule: rule,
		Cs:   linkedhashset.New(),
	}
	return c
}

func CanBeEpsylonStr(str string, Rules utils.RulesTable) bool{
	if str == utils.EpsilonSymbol{
		return true
	}
	runes := []rune(str)
	canBeEpsylon := false
	pos := 0
	for ; pos < len(runes); pos++{
		canBeEpsylon = false
		if !unicode.IsUpper(runes[pos]) && string(runes[pos]) != utils.EpsilonSymbol{
			break
		}
		canBeEpsylon = Rules.CanBeEpsylon(string(runes[pos]))
		if !canBeEpsylon{
			break
		}
	}
	return canBeEpsylon && pos == len(runes)
}

func BuildChooseSet(Rules utils.RulesTable, Sets map[string]Result) []Choose{
	results := make([]Choose, 0)
	for _, v := range Rules.Ordered{
		results = append(results, NewChoose(v))
	}
	for i, v := range results{
		if CanBeEpsylonStr(v.Rule.Right, Rules){
			results[i].Cs = results[i].Cs.Union(Sets[v.Rule.Left].Forward)
		}
		for _, sym := range v.Rule.Right{
			if unicode.IsUpper(sym){
				results[i].Cs = results[i].Cs.Union(Sets[string(sym)].First)
				if !Rules.CanBeEpsylon(string(sym)){
					break
				}
			} else{
				if string(sym) != utils.EpsilonSymbol{
					results[i].Cs.Add(sym)
				}
				break
			}
		}
	}
	return results
}

func NewResult(nonTerm string) Result {
	r := Result{
		NonTerminal: nonTerm,
		First:       linkedhashset.New(),
		Forward:     linkedhashset.New(),
	}
	r.First.Add([]rune(nonTerm)[0])
	return r
}

func BuildSetFirst(Rules utils.RulesTable) map[string]Result {
	// Инициализируем
	results := make(map[string]Result)
	for k, v := range Rules.LeftMapped {
		results[k] = NewResult(k)
		for _, rule := range v {
			runes := []rune(rule.Right)
			results[k].First.Add(runes[0])
			if Rules.CanBeEpsylon(string(runes[0])) {
				for i := 1; i < len(runes); i++ {
					if !unicode.IsUpper(runes[i]) {
						break
					}
					results[k].First.Add(runes[i])
					if !Rules.CanBeEpsylon(string(runes[i])) {
						break
					}
				}
			}
		}
	}

	for {
		wasChanged := false
		for k, r := range results {
			for _, v := range r.First.Values() {
				if unicode.IsUpper(v.(rune)) {
					oldSize := results[k].First.Size()
					r.First = results[k].First.Union(results[string(v.(rune))].First)
					results[k] = r
					if oldSize != results[k].First.Size() {
						wasChanged = true
					}
				}
			}
		}
		if !wasChanged {
			break
		}
	}
	return results
}

func BuildSets(Rules utils.RulesTable) map[string]Result {
	results := BuildSetFirst(Rules)
	results["S"].Forward.Add([]rune(utils.EndLineSymbol)[0])
	for {
		wasChanged := false
		for k, v := range Rules.LeftMapped {
			for _, rule := range v {
				length := utf8.RuneCountInString(rule.Right)
				runes := []rune(rule.Right)
				for pos, sym := range runes {
					stringSym := string(sym)
					if unicode.IsUpper(sym) {
						oldSize := results[stringSym].Forward.Size()
						if pos != length-1 {
							if unicode.IsUpper(runes[pos+1]) {
								r := results[stringSym]
								forwardSym := runes[pos+1]
								r.Forward = r.Forward.Union(results[string(forwardSym)].First)
								results[stringSym] = r
								if Rules.CanBeEpsylon(string(runes[pos+1])) {
									posNext := pos + 2
									canBeEpsylon := false
									for ; posNext < len(runes); posNext++ {
										if !unicode.IsUpper(runes[posNext]) {
											break
										}
										r := results[stringSym]
										forwardSym := runes[pos+1]
										r.Forward = r.Forward.Union(results[string(forwardSym)].First)
										results[stringSym] = r
										canBeEpsylon = Rules.CanBeEpsylon(string(runes[posNext]))
										if !canBeEpsylon {
											break
										}
									}
									if posNext == len(runes) && canBeEpsylon {
										r := results[stringSym]
										forwardSym := k
										r.Forward = r.Forward.Union(results[string(forwardSym)].Forward)
										results[stringSym] = r
									}
								}
							} else {
								results[stringSym].Forward.Add(runes[pos+1])
							}
						} else {
							r := results[stringSym]
							r.Forward = r.Forward.Union(results[k].Forward)
							results[stringSym] = r
						}
						if oldSize != results[stringSym].Forward.Size() {
							wasChanged = true
						}
					}
				}
			}
		}
		if !wasChanged {
			break
		}
	}
	return results
}
