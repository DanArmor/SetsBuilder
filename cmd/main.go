package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
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
	jsonFile := "rules.json"
	file, err := os.Open(jsonFile)
	if err != nil {
		fmt.Print("Error during reading of file: ", err.Error())
		os.Exit(1)
	}
	defer file.Close()
	byteValue, err := io.ReadAll(file)
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

	fmt.Printf("\n\033[34m CHOOSE FOR RULES \033[0m\n")
	fmt.Printf("\033[34m ========== \033[0m\n")
	for k, v := range result {
		values := v.First.Values()
		for _, r := range values{
			if unicode.IsUpper(r.(rune)){
				v.First.Remove(r)
			}
		}
		result[k] = v
	}
	for k, v := range result {
		values := v.Forward.Values()
		for _, r := range values{
			if unicode.IsUpper(r.(rune)){
				v.Forward.Remove(r)
			}
		}
		result[k] = v
	}
	choose := sets.BuildChooseSet(rt, result)
	for _, v := range choose{
		strs := make([]string, 0)
		values := v.Cs.Values()
		for _, iv := range values {
			strs = append(strs, string(iv.(rune)))
		}
		fmt.Printf("\033[32m%d.\033[0m %s -> %s   {\033[34m %s \033[0m} \n", v.Rule.Number, v.Rule.Left, v.Rule.Right, strings.Join(strs, " "))
	}
	hasIntersect := false
	for i := 0; i < len(choose); i++{
		for j := i+1; j < len(choose); j++{
			if choose[i].Rule.Left == choose[j].Rule.Left{
				if inter := choose[i].Cs.Intersection(choose[j].Cs); inter.Size() != 0{
					values := inter.Values()
					strs := make([]string, 0)
					for _, iv := range values {
						strs = append(strs, string(iv.(rune)))
					}
					hasIntersect = true
					fmt.Printf("\033[31mПравила %d и %d имеют пересечение множеств ВЫБОР: { \033[33m %s \033[31m }\033[0m\n", choose[i].Rule.Number, choose[j].Rule.Number, strings.Join(strs, " "))
				}
			}
		}
	}
	if !hasIntersect{
		fmt.Printf("\033[32mНет пересечений множеств ВЫБОР\033[0m\n")
	}
}
