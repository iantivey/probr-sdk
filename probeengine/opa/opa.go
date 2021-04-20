package opa

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	//"io/ioutil"
	"log"
)

//TODO: the regoFile will need to be packages using Packagr
func Eval(regoFile string, regoPackageName string, regoFuncName string, jsonInput *[]byte) (bool, error) {
	load := make([]string, 1)
	load[0] = regoFile

	r := rego.New(
		rego.Query(fmt.Sprintf("x = data.%s.%s", regoPackageName, regoFuncName)),
		rego.Load(load, nil))

	ctx := context.Background()
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		log.Printf("error line 16")
		return false, err
	}
	var input interface{}

	if err := json.Unmarshal(*jsonInput, &input); err != nil {
		log.Printf("error line 32")
		return false, err
	}

	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		log.Printf("error line 38")
		return false, err
	}

	v, ok := rs[0].Bindings["x"].(bool)
	if !ok {
		log.Printf("Did not get bool")
		log.Printf(fmt.Sprintf("rs[0].Bindings[\"x\"] = %v", rs[0].Bindings["x"]))
		return false, fmt.Errorf("Did not get bool")
	}

	return v, nil
}
