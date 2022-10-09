package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"log"
	"os"
	"strings"
	"text/template"
)

type TemplatedChart struct {
	ChartName      string
	ChartReference string
	ChartUrl       string
}

type TemplatedFlow struct {
	FlowName       string
	ChartReference string
	ChartVersion   string
	Namespace      string
}

func main() {
	answers := struct {
		FlowName string // survey will match the question and field names
		FlowType string
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	helmAnswers := struct {
		ChartName      string // survey will match the question and field names
		ChartUrl       string
		ChartReference string
		ChartVersion   string
		Namespace      string
	}{}

	if answers.FlowType == "helm" {
		err := survey.Ask(chartQs, &helmAnswers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	flowName := strings.TrimSpace(answers.FlowName)
	if err := os.Mkdir(flowName, os.ModePerm); err != nil {
		fmt.Println("Unable to create Flow Folder", flowName)
		log.Fatal(err)
	}

	templates, err := template.ParseFiles("templates/helm.properties.tmpl")
	if err != nil {
		fmt.Println("template processing error ")
	}

	tc := getChart(helmAnswers.ChartName, helmAnswers.ChartReference, helmAnswers.ChartUrl)
	f, err := os.Create(flowName + "/" + "helm.properties")
	if err != nil {
		fmt.Println("Unable to create helm properties", flowName)
		log.Fatal(err)
	}
	err = templates.Execute(f, tc)
	if err != nil {
		log.Println("executing template:", err)
	}

	flowTemplate, err := template.ParseFiles("templates/flow.tmpl")
	if err != nil {
		fmt.Println("template processing error ")
	}
	ft := getFlow(flowName, helmAnswers.Namespace, helmAnswers.ChartReference, helmAnswers.ChartVersion)
	f, err = os.Create(flowName + "/" + "ck8s-" + flowName + ".concord.yaml")
	if err != nil {
		fmt.Println("Unable to create helm properties", flowName)
		log.Fatal(err)
	}
	err = flowTemplate.Execute(f, ft)
	if err != nil {
		log.Println("executing template:", err)
	}

}

func getFlow(flowName string, namespace string, chartReference string, chartVersion string) TemplatedFlow {
	return TemplatedFlow{
		FlowName:       flowName,
		Namespace:      namespace,
		ChartReference: chartReference,
		ChartVersion:   chartVersion,
	}
}

func getChart(chartName string, chartReference string, chartUrl string) TemplatedChart {
	if chartName == "" {
		fmt.Println("chartName passed is empty")
	}

	return TemplatedChart{
		ChartName:      strings.TrimSpace(chartName),
		ChartUrl:       strings.TrimSpace(chartUrl),
		ChartReference: strings.TrimSpace(chartReference),
	}
}

var qs = []*survey.Question{
	{
		Name:     "flowName",
		Prompt:   &survey.Input{Message: "What is name of the flow?"},
		Validate: survey.Required,
	},
	{
		Name: "flowType",
		Prompt: &survey.Select{
			Message: "Choose the type of flow:",
			Options: []string{"helm", "kubectl"},
			Default: "helm",
		},
	},
}

var chartQs = []*survey.Question{
	{
		Name:     "chartName",
		Prompt:   &survey.Input{Message: "What is name of the chart ?"},
		Validate: survey.Required,
	},
	{
		Name:     "chartUrl",
		Prompt:   &survey.Input{Message: "Enter Chart Url?"},
		Validate: survey.Required,
	},
	{
		Name:     "chartReference",
		Prompt:   &survey.Input{Message: "Enter Chart Reference?"},
		Validate: survey.Required,
	},

	{
		Name:     "chartVersion",
		Prompt:   &survey.Input{Message: "Enter Chart Version?"},
		Validate: survey.Required,
	},
	{
		Name:     "namespace",
		Prompt:   &survey.Input{Message: "Enter the namespace ?"},
		Validate: survey.Required,
	},
}
