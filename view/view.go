package view

import (
	"bytes"

	"github.com/CloudyKit/jet"

	"github.com/s0kil/ESF-Products-Manager/fault"
	"github.com/s0kil/ESF-Products-Manager/model"
)

func ProductIndex(v *jet.Set, p []model.Product) (writer bytes.Buffer) {
	vars := make(jet.VarMap).Set("products", []model.Product{})

	view, err := v.GetTemplate("product/index.jet")
	if err != nil {
		fault.Report(err, "Could Not Get Template")
	}

	err = view.Execute(&writer, vars, &p)
	if err != nil {
		fault.Report(err, "Could Not Execute Template")
	}

	return
}
