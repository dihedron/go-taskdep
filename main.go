package main

import (
	"fmt"
	"os"

	"github.com/dihedron/go-taskdep/tasks"
	"github.com/fako1024/topo"
)

// List of all simple strings (to be sorted)
var stringsToSort = []string{
	"A", "B", "C", "D", "E", "F", "G", "H",
}

// List of dependencies
var stringDependencies = []topo.Dependency{
	topo.Dependency{Child: "B", Parent: "A"},
	topo.Dependency{Child: "B", Parent: "C"},
	topo.Dependency{Child: "B", Parent: "D"},
	topo.Dependency{Child: "A", Parent: "E"},
	topo.Dependency{Child: "D", Parent: "C"},
}

var tks = []tasks.Task{
	tasks.Task{
		ID: "postgresql-jdbc-driver",
		Instructions: []string{
			`echo "pre-installing postgresql-jdbc-driver"`,
			`echo "post-installing postgresql-jdbc-driver"`,
		},
	},
	tasks.Task{
		ID: "terze-valute-settings",
		Instructions: []string{
			`echo "install system properties"`,
			`echo "install datasource"`,
		},
		Dependencies: []string{
			"postgresql-jdbc-driver",
		},
	},
	tasks.Task{
		ID: "terze-valute-rest",
		Instructions: []string{
			`echo "install terze-valute-rest.war"`,
		},
		Dependencies: []string{
			"terze-valute-settings",
		},
	},
	tasks.Task{
		ID: "terze-valute-spa",
		Instructions: []string{
			`echo "install terze-valute-spa.war"`,
		},
		Dependencies: []string{
			"terze-valute-rest",
			"terze-valute-settings",
		},
	},
}

func main() {
	// TODO: start working on this

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
}
