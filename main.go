package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/dihedron/go-taskdep/tasks"
	"github.com/fako1024/topo"
)

var items = map[string]tasks.Task{
	"postgresql-jdbc-driver": tasks.Task{
		ID: "postgresql-jdbc-driver",
		Instructions: []string{
			`echo "pre-installing postgresql-jdbc-driver"`,
			`echo "post-installing postgresql-jdbc-driver"`,
		},
	},
	"my-app-settings": tasks.Task{
		ID: "my-app-settings",
		Instructions: []string{
			`echo "install system properties"`,
			`echo "install datasource"`,
		},
		Dependencies: []string{
			"postgresql-jdbc-driver",
		},
	},
	"my-app-rest": tasks.Task{
		ID: "my-app-rest",
		Instructions: []string{
			`echo "install my-app-rest.war"`,
		},
		Dependencies: []string{
			"my-app-settings",
		},
	},
	"my-app-spa": tasks.Task{
		ID: "my-app-spa",
		Instructions: []string{
			`echo "install my-app-spa.war"`,
		},
		Dependencies: []string{
			"my-app-rest",
			"my-app-settings",
		},
	},
}

func main() {
	// TODO: start working on this

	js, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Original nodes:\n%s\n", string(js))

	list := []string{}
	for key := range items {
		list = append(list, key)
	}
	js, err = json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("list ot items:\n%s\n", string(js))

	// List of dependencies
	dependencies := []topo.Dependency{}
	for _, item := range items {
		for _, dependency := range item.Dependencies {
			if _, ok := items[dependency]; !ok {
				fmt.Printf("error: dependency on non-existing item: %s\n", dependency)
				os.Exit(1)
			}
			dependencies = append(dependencies, topo.Dependency{
				Parent: dependency,
				Child:  item.ID,
			})
		}
	}

	js, err = json.MarshalIndent(dependencies, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Dependencies:\n%s\n", string(js))

	// Getter function to convert original elements to a generic type
	getter := func(i int) topo.Type {
		return list[i]
	}

	// Setter function to restore the original type of the data
	setter := func(i int, val topo.Type) {
		list[i] = val.(string)
	}

	// Perform topological sort
	if err := topo.Sort(list, dependencies, getter, setter); err != nil {
		fmt.Printf("Error performing topological sort on tasks: %v\n", err)
		os.Exit(1)
	}

	// Print resulting Slice in order
	fmt.Println("Sorted list of items:", list)
	fmt.Println("The following dependencies were taken into account:")
	for _, dependency := range dependencies {
		fmt.Println(dependency)
	}

	os.Exit(0)

	/*
		// Getter function to convert original elements to a generic type
		getter := func(i int) topo.Type {
			return stringsToSort[i]
		}

		// Setter function to restore the original type of the data
		setter := func(i int, val topo.Type) {
			stringsToSort[i] = val.(string)
		}

		// Perform topological sort
		if err := topo.Sort(stringsToSort, stringDependencies, getter, setter); err != nil {
			fmt.Printf("Error performing topological sort on slice of strings: %s\n", err)
			os.Exit(1)
		}

		// Print resulting Slice in order
		fmt.Println("Sorted list of strings:", stringsToSort)
		fmt.Println("The following dependencies were taken into account:")
		for _, dep := range stringDependencies {
			fmt.Println(dep)
		}
	*/
}
