package main

/**
 * Commnad reviewdatamodel is added to help reviewing generated data model with sysl
 * file produced by command import. Generate data model diagrams using the following command:
 * sysl reviewdata --root=/Users/guest/data -t Test -o Test.png Test
 * sysl reviewdatamodel --root=/Users/guest/data -t Test -o Test.png Test.sysl
 */

import (
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/anz-bank/sysl/pkg/sysl"
	"github.com/sirupsen/logrus"
)

func GenerateDataModelsView(datagenParams *CmdContextParamDatagen,
	model *sysl.Module, logger *logrus.Logger) (map[string]string, error) {
	outmap := make(map[string]string)

	logger.Debugf("title: %s\n", datagenParams.title)
	logger.Debugf("output: %s\n", datagenParams.output)

	spclass := constructFormatParser("", datagenParams.classFormat)

	apps := model.GetApps()
	for appName := range apps {
		app := apps[appName]
		outputDir := datagenParams.output
		if strings.Contains(outputDir, "%(appname)") {
			of := MakeFormatParser(datagenParams.output)
			outputDir = of.FmtOutput(appName, "", app.GetLongName(), app.GetAttrs())
		}
		var stringBuilder strings.Builder
		if app != nil {
			dataParam := &DataModelParam{
				mod:   model,
				app:   app,
				title: datagenParams.title,
			}
			v := MakeDataModelView(spclass, dataParam.mod, &stringBuilder, dataParam.title, "")
			outmap[outputDir] = v.GenerateDataView(dataParam)
		}
	}

	return outmap, nil
}

// Process pure Sysl datamodel file produced by import cmd
type reviewDatamodelCmd struct {
	plantumlmixin
	CmdContextParamDatagen
}

func (p *reviewDatamodelCmd) Name() string       { return "reviewdatamodel" }
func (p *reviewDatamodelCmd) MaxSyslModule() int { return 1 }

func (p *reviewDatamodelCmd) Configure(app *kingpin.Application) *kingpin.CmdClause {
	cmd := app.Command(p.Name(), "Generate data models for review from pure sysl file produced by command import")
	cmd.Alias("reviewdata")
	cmd.Flag("class_format",
		"Specify the format string for data diagram participants. "+
			"May include %%(appname) and %%(@foo) for attribute foo (default: %(classname))",
	).Default("%(classname)").StringVar(&p.classFormat)

	cmd.Flag("title", "diagram title").Short('t').StringVar(&p.title)

	cmd.Flag("output",
		"output file (default: %(appname).png)",
	).Default("%(appname).png").Short('o').StringVar(&p.output)

	p.AddFlag(cmd)

	EnsureFlagsNonEmpty(cmd)
	return cmd
}

func (p *reviewDatamodelCmd) Execute(args ExecuteArgs) error {
	outmap, err := GenerateDataModelsView(&p.CmdContextParamDatagen, args.Modules[0], args.Logger)
	if err != nil {
		return err
	}
	return p.GenerateFromMap(outmap, args.Filesystem)
}
