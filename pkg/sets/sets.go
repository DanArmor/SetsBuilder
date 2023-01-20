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
		for _, rule := range v{
			runes := []rune(rule.Right)
			results[k].First.Add(runes[0])
		}
	}

	for{
		wasChanged := false
		for k, r := range results {
			for _, v := range r.First.Values(){
				if unicode.IsUpper(v.(rune)){
					oldSize := results[k].First.Size()
					r.First = results[k].First.Union(results[string(v.(rune))].First)
					results[k] = r
					if oldSize != results[k].First.Size(){
						wasChanged = true
					}
				}
			}
		}
		if !wasChanged{
			break
		}
	}
	return results
}

func BuildSets(Rules utils.RulesTable) map[string]Result{
	results := BuildSetFirst(Rules)
	results["S"].Forward.Add([]rune(utils.EndLineSymbol)[0])
	for{
		wasChanged := false
		for k, v := range Rules.LeftMapped{
			for _, rule := range v{
				length := utf8.RuneCountInString(rule.Right)
				runes := []rune(rule.Right)
				for pos, sym := range runes{
					stringSym := string(sym)
					if unicode.IsUpper(sym){
						oldSize := results[stringSym].Forward.Size()
						if pos != length-1{
							if unicode.IsUpper(runes[pos+1]){
								r := results[stringSym]
								forwardSym := runes[pos+1]
								r.Forward = r.Forward.Union(results[string(forwardSym)].First)
								results[stringSym] = r
							} else{
								results[stringSym].Forward.Add(runes[pos+1])
							}
						} else{
							r := results[stringSym]
							r.Forward = r.Forward.Union(results[k].Forward)
							results[stringSym] = r
						}
						if oldSize != results[stringSym].Forward.Size(){
							wasChanged = true
						}
					}
				}
			}
		}
		if !wasChanged{
			break
		}
	}
	return results
}