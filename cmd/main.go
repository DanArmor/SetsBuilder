package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"unicode"

	"github.com/DanArmor/taEkzUtils/pkg/sets"
	"github.com/DanArmor/taEkzUtils/pkg/utils"
	godsUtils "github.com/emirpasic/gods/utils"
)

func remove(slice []interface{}, s int) []interface{} {
	if s == 0{
		return slice[1:]
	}
	return append(slice[:s], slice[s+1:]...)
}

func filterOutEpsylon(slice []interface{}) []interface{} {
	for i := 0; i < len(slice); i++ {
		if string(slice[i].(rune)) == utils.EpsilonSymbol {
			slice = remove(slice, i)
			i = 0
		}
	}
	return slice
}

func filterOutNonTerms(slice []interface{}) []interface{} {
	for i := 0; i < len(slice); i++ {
		if unicode.IsUpper(slice[i].(rune)) {
			slice = remove(slice, i)
			i = 0
		}
	}
	return slice
}

func filterOutFull(slice []interface{}) []interface{} {
	slice = filterOutEpsylon(slice)
	slice = filterOutNonTerms(slice)
	return slice
}


func main() {
	rt := utils.NewRulesTable()
	jsonFile := "rules4.json"
	byteValue, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Print("Error during reading of file: ", err.Error())
		os.Exit(1)
	}
	var rulesJSON utils.RulesJSON
	json.Unmarshal(byteValue, &rulesJSON)
	for _, v := range rulesJSON.Rules {
		rt.AddRule(v.Left, v.Right)
	}
	fmt.Printf("\033[34m DEBUG \033[0m\n")
	for _, v := range rt.Ordered {
		fmt.Printf("%d. %s -> %s\n", v.Number, v.Left, v.Right)
	}
	fmt.Printf("\033[34m ========== \033[0m\n")
	result := sets.BuildSets(rt)
	for k, v := range result {
		fmt.Printf("\033[31mKey: %s\033[0m\n", k)
		vals := v.First.Values()
		vals = filterOutEpsylon(vals)
		fmt.Printf("\033[34mFirst(%d): \033[0m", len(vals))
		godsUtils.Sort(vals, godsUtils.RuneComparator)

		for i := 0; i < len(vals); i++ {
			fmt.Printf("%s ", string(vals[i].(rune)))
		}

		forVals := v.Forward.Values()
		forVals = filterOutEpsylon(forVals)
		fmt.Printf("\n\033[36mForward(%d): \033[0m", len(forVals))
		godsUtils.Sort(forVals, godsUtils.RuneComparator)
		for i := 0; i < len(forVals); i++ {
			fmt.Printf("%s ", string(forVals[i].(rune)))
		}
		fmt.Printf("\n\033[32m=============\033[0m\n")
	}
	fmt.Printf("\n\033[34m SAME WITHOUT NON TERMINALS \033[0m\n")
	fmt.Printf("\033[34m ========== \033[0m\n")
	for k, v := range result {
		fmt.Printf("\033[31mKey: %s\033[0m\n", k)

		vals := v.First.Values()
		vals = filterOutFull(vals)

		fmt.Printf("\033[34mFirst(%d): \033[0m", len(vals))
		godsUtils.Sort(vals, godsUtils.RuneComparator)

		for i := 0; i < len(vals); i++ {
			fmt.Printf("%s ", string(vals[i].(rune)))
		}

		forVals := v.Forward.Values()
		forVals = filterOutFull(forVals)

		fmt.Printf("\n\033[36mForward(%d): \033[0m", len(forVals))
		godsUtils.Sort(forVals, godsUtils.RuneComparator)
		for i := 0; i < len(forVals); i++ {
			fmt.Printf("%s ", string(forVals[i].(rune)))
		}
		fmt.Printf("\n\033[32m=============\033[0m\n")
	}

}
