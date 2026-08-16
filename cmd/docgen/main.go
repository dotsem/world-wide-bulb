// Package main generates and synchronizes Just command documentation in Markdown files.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type recipe struct {
	name string
	doc  string
}

type module struct {
	name    string
	title   string
	path    string
	recipes []recipe
}

var (
	modRegex    = regexp.MustCompile(`^mod\s+([a-zA-Z0-9_-]+)\s+['"]([^'"]+)['"]`)
	recipeRegex = regexp.MustCompile(`^([a-zA-Z0-9_-]+)(?:\s+[^:]*?)?:(?:\s+.*)?$`)
)

func main() {
	targetsFlag := flag.String("targets", "README.md,GUIDELINES.md", "comma-separated markdown files to update")
	justfileFlag := flag.String("justfile", "justfile", "path to root justfile")
	markerFlag := flag.String("marker", "JUST_COMMANDS", "marker name for injection comments")
	checkFlag := flag.Bool("check", false, "verify docs without writing (fails if out of sync)")
	flag.Parse()

	modules, err := parseRootAndModules(*justfileFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing justfiles: %v\n", err)
		os.Exit(1)
	}

	generatedMarkdown := formatModulesMarkdown(modules)
	targetFiles := strings.Split(*targetsFlag, ",")
	startMarker := fmt.Sprintf("<!-- %s_START -->", *markerFlag)
	endMarker := fmt.Sprintf("<!-- %s_END -->", *markerFlag)

	driftDetected := false

	for _, target := range targetFiles {
		target = strings.TrimSpace(target)
		if target == "" {
			continue
		}

		cleanTarget := filepath.Clean(target)
		content, err := os.ReadFile(cleanTarget) //nolint:gosec // Target files are CLI arguments for local code generation.
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", target, err)
			os.Exit(1)
		}

		updated, changed, err := injectMarkdown(string(content), generatedMarkdown, startMarker, endMarker)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error updating %s: %v\n", target, err)
			os.Exit(1)
		}

		if !changed {
			continue
		}

		if *checkFlag {
			fmt.Fprintf(os.Stderr, "drift detected in %s (documentation is out of sync)\n", target)
			driftDetected = true
		} else {
			if err := os.WriteFile(cleanTarget, []byte(updated), 0o600); err != nil { //nolint:gosec // Target files are CLI arguments for local code generation.
				fmt.Fprintf(os.Stderr, "error writing %s: %v\n", target, err)
				os.Exit(1)
			}
			fmt.Printf("updated %s\n", target)
		}
	}

	if driftDetected {
		fmt.Fprintln(os.Stderr, "\nrun 'just docgen' locally to regenerate command documentation.")
		os.Exit(1)
	}

	if *checkFlag {
		fmt.Println("documentation is in sync with just recipes.")
	}
}

func parseRootAndModules(rootPath string) ([]module, error) {
	rootDir := filepath.Dir(rootPath)
	rootModule, submods, err := parseJustfile(rootPath, "", "Root Commands")
	if err != nil {
		return nil, err
	}

	modules := []module{rootModule}

	for _, sub := range submods {
		subPath := filepath.Join(rootDir, sub.path)
		var title string
		switch sub.name {
		case "go":
			title = "Backend (`just go <cmd>`)"
		case "web":
			title = "Frontend (`just web <cmd>`)"
		default:
			title = fmt.Sprintf("%s (`just %s <cmd>`)", strings.ToUpper(sub.name[:1])+sub.name[1:], sub.name)
		}

		subMod, _, err := parseJustfile(subPath, sub.name, title)
		if err != nil {
			return nil, err
		}
		modules = append(modules, subMod)
	}

	return modules, nil
}

type submodRef struct {
	name string
	path string
}

func parseJustfile(filePath, modName, title string) (module, []submodRef, error) {
	cleanPath := filepath.Clean(filePath)
	file, err := os.Open(cleanPath)
	if err != nil {
		return module{}, nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var (
		recipes    []recipe
		submods    []submodRef
		docComment []string
		scanner    = bufio.NewScanner(file)
	)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			docComment = nil
			continue
		}

		if match := modRegex.FindStringSubmatch(trimmed); len(match) > 2 {
			submods = append(submods, submodRef{name: match[1], path: match[2]})
			docComment = nil
			continue
		}

		if after, ok := strings.CutPrefix(trimmed, "#"); ok {
			comment := strings.TrimSpace(after)
			docComment = append(docComment, comment)
			continue
		}

		if strings.Contains(trimmed, ":=") {
			docComment = nil
			continue
		}

		if match := recipeRegex.FindStringSubmatch(trimmed); len(match) > 1 {
			name := match[1]
			if !strings.HasPrefix(name, "_") {
				doc := strings.Join(docComment, " ")
				if doc == "" {
					doc = "—"
				}
				recipes = append(recipes, recipe{name: name, doc: doc})
			}
			docComment = nil
			continue
		}

		if trimmed != "" {
			docComment = nil
		}
	}

	if err := scanner.Err(); err != nil {
		return module{}, nil, err
	}

	return module{
		name:    modName,
		title:   title,
		path:    cleanPath,
		recipes: recipes,
	}, submods, nil
}

func formatModulesMarkdown(modules []module) string {
	var builder strings.Builder

	for i, mod := range modules {
		if i > 0 {
			builder.WriteString("\n")
		}

		fmt.Fprintf(&builder, "### %s\n\n", mod.title)
		builder.WriteString("| Command | Description |\n")
		builder.WriteString("| :--- | :--- |\n")

		for _, r := range mod.recipes {
			cmdStr := fmt.Sprintf("`just %s`", r.name)
			if mod.name != "" {
				cmdStr = fmt.Sprintf("`just %s %s`", mod.name, r.name)
			}
			fmt.Fprintf(&builder, "| %s | %s |\n", cmdStr, r.doc)
		}
	}

	return builder.String()
}

func injectMarkdown(source, generated, startMarker, endMarker string) (string, bool, error) {
	startIdx := strings.Index(source, startMarker)
	if startIdx == -1 {
		return "", false, fmt.Errorf("marker %q not found in document", startMarker)
	}

	endIdx := strings.Index(source, endMarker)
	if endIdx == -1 {
		return "", false, fmt.Errorf("marker %q not found in document", endMarker)
	}

	if endIdx < startIdx {
		return "", false, fmt.Errorf("invalid marker order (%q found after %q)", startMarker, endMarker)
	}

	replacement := fmt.Sprintf("%s\n\n%s\n%s", startMarker, generated, endMarker)
	existingBlock := source[startIdx : endIdx+len(endMarker)]

	if existingBlock == replacement {
		return source, false, nil
	}

	updated := source[:startIdx] + replacement + source[endIdx+len(endMarker):]
	return updated, true, nil
}
