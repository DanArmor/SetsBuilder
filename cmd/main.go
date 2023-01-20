package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/DanArmor/taEkzUtils/pkg/sets"
	"github.com/DanArmor/taEkzUtils/pkg/utils"
	godsUtils "github.com/emirpasic/gods/utils"
)

func main(){
	rt := utils.NewRulesTable()
	jsonFile := "rules.json"
	byteValue, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Print("Error during reading of file: ", err.Error())
		os.Exit(1)
	}
	var rulesJSON utils.RulesJSON
	json.Unmarshal(byteValue, &rulesJSON)
	for _, v := range rulesJSON.Rules{
		rt.AddRule(v.Left, v.Right)
	}
	fmt.Printf("\033[34m DEBUG \033[0m\n")
	for _, v := range rt.Ordered{
		fmt.Printf("%d. %s -> %s\n", v.Number, v.Left, v.Right)
	}
	fmt.Printf("\033[34m ========== \033[0m\n")
	result := sets.BuildSets(rt)
	for k, v := range result {
		fmt.Printf("\033[32mFirst\033[0m\n")
		fmt.Printf("\033[31mKey: %s\033[0m\n", k)
		vals := v.First.Values()
		godsUtils.Sort(vals, godsUtils.RuneComparator)
		
		fmt.Printf("Values: ")
		for i := 0; i < len(vals); i++{
			fmt.Printf("%s ", string(vals[i].(rune)))
		}
		fmt.Printf("\n\033[32m=============\033[0m\n")

		fmt.Printf("\033[32mForward\033[0m\n")
		fmt.Printf("\033[31mKey: %s\033[0m\n", k)
		forVals := v.Forward.Values()
		godsUtils.Sort(forVals, godsUtils.RuneComparator)
		fmt.Printf("Values: ")
		for i := 0; i < len(forVals); i++{
			fmt.Printf("%s ", string(forVals[i].(rune)))
		}
		fmt.Printf("\n\033[32m=============\033[0m\n")
	}
}