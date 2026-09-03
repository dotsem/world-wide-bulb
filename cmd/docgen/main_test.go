package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseJustfile(t *testing.T) {
	t.Run("parses recipes with single and multi-line comments", func(t *testing.T) {
		content := `
# Install dependencies
# and git hooks
install:
    pnpm install

# Run linters
lint:
    golangci-lint run
`
		tmpFile := filepath.Join(t.TempDir(), "justfile")
		require.NoError(t, os.WriteFile(tmpFile, []byte(content), 0o600))

		mod, submods, err := parseJustfile(tmpFile, "", "Root Commands")
		require.NoError(t, err)
		assert.Empty(t, submods)
		assert.Equal(t, "Root Commands", mod.title)
		require.Len(t, mod.recipes, 2)
		assert.Equal(t, "install", mod.recipes[0].name)
		assert.Equal(t, "Install dependencies and git hooks", mod.recipes[0].doc)
		assert.Equal(t, "lint", mod.recipes[1].name)
		assert.Equal(t, "Run linters", mod.recipes[1].doc)
	})

	t.Run("handles parameters, assignments, and ignores private recipes", func(t *testing.T) {
		content := `
set dotenv-load := true
NAME := "test"

# Run with optional args
build *args='':
    go build {{args}}

# Private helper
_internal:
    echo "secret"

undoc:
    echo "undoc"
`
		tmpFile := filepath.Join(t.TempDir(), "justfile")
		require.NoError(t, os.WriteFile(tmpFile, []byte(content), 0o600))

		mod, _, err := parseJustfile(tmpFile, "", "Root")
		require.NoError(t, err)
		require.Len(t, mod.recipes, 2)
		assert.Equal(t, "build", mod.recipes[0].name)
		assert.Equal(t, "Run with optional args", mod.recipes[0].doc)
		assert.Equal(t, "undoc", mod.recipes[1].name)
		assert.Equal(t, "—", mod.recipes[1].doc)
	})

	t.Run("extracts submodule definitions", func(t *testing.T) {
		content := `
mod go 'just/go.just'
mod web "just/web.just"

# Top-level recipe
test:
    just go test
`
		tmpFile := filepath.Join(t.TempDir(), "justfile")
		require.NoError(t, os.WriteFile(tmpFile, []byte(content), 0o600))

		_, submods, err := parseJustfile(tmpFile, "", "Root")
		require.NoError(t, err)
		require.Len(t, submods, 2)
		assert.Equal(t, submodRef{name: "go", path: "just/go.just"}, submods[0])
		assert.Equal(t, submodRef{name: "web", path: "just/web.just"}, submods[1])
	})

	t.Run("returns error on non-existent file", func(t *testing.T) {
		_, _, err := parseJustfile("/non/existent/path/justfile", "", "Root")
		assert.Error(t, err)
	})
}

func TestParseRootAndModules(t *testing.T) {
	t.Run("resolves root and submodules with correct titles", func(t *testing.T) {
		dir := t.TempDir()
		justDir := filepath.Join(dir, "just")
		require.NoError(t, os.MkdirAll(justDir, 0o750))

		rootFile := filepath.Join(dir, "justfile")
		goFile := filepath.Join(justDir, "go.just")
		customFile := filepath.Join(justDir, "tools.just")

		require.NoError(t, os.WriteFile(rootFile, []byte("mod go 'just/go.just'\nmod tools 'just/tools.just'\n# Root task\nroot:\n  echo 1"), 0o600))
		require.NoError(t, os.WriteFile(goFile, []byte("# Go test\ntest:\n  go test"), 0o600))
		require.NoError(t, os.WriteFile(customFile, []byte("# Custom tool\nrun:\n  echo 2"), 0o600))

		modules, err := parseRootAndModules(rootFile)
		require.NoError(t, err)
		require.Len(t, modules, 3)

		assert.Equal(t, "Root Commands", modules[0].title)
		assert.Equal(t, "Backend (`just go <cmd>`)", modules[1].title)
		assert.Equal(t, "Tools (`just tools <cmd>`)", modules[2].title)
	})

	t.Run("returns error when submodule file is missing", func(t *testing.T) {
		dir := t.TempDir()
		rootFile := filepath.Join(dir, "justfile")
		require.NoError(t, os.WriteFile(rootFile, []byte("mod missing 'just/missing.just'\n"), 0o600))

		_, err := parseRootAndModules(rootFile)
		assert.Error(t, err)
	})
}

func TestFormatModulesMarkdown(t *testing.T) {
	modules := []module{
		{
			name:  "",
			title: "Root Commands",
			recipes: []recipe{
				{name: "install", doc: "Install dependencies"},
			},
		},
		{
			name:  "go",
			title: "Backend (`just go <cmd>`)",
			recipes: []recipe{
				{name: "test", doc: "Run tests"},
			},
		},
	}

	markdown := formatModulesMarkdown(modules)
	expected := `### Root Commands

| Command | Description |
| :--- | :--- |
| ` + "`just install`" + ` | Install dependencies |

### Backend (` + "`just go <cmd>`" + `)

| Command | Description |
| :--- | :--- |
| ` + "`just go test`" + ` | Run tests |
`
	assert.Equal(t, expected, markdown)
}

func TestInjectMarkdown(t *testing.T) {
	startMarker := "<!-- TEST_START -->"
	endMarker := "<!-- TEST_END -->"
	generated := "### Generated Section\n\nContent"

	t.Run("successfully replaces content between markers", func(t *testing.T) {
		source := "# Document\n\n<!-- TEST_START -->\nold content\n<!-- TEST_END -->\n\n## Footer"
		updated, changed, err := injectMarkdown(source, generated, startMarker, endMarker)

		require.NoError(t, err)
		assert.True(t, changed)
		expected := "# Document\n\n<!-- TEST_START -->\n\n### Generated Section\n\nContent\n<!-- TEST_END -->\n\n## Footer"
		assert.Equal(t, expected, updated)
	})

	t.Run("returns unchanged when content already matches", func(t *testing.T) {
		source := "# Document\n\n<!-- TEST_START -->\n\n### Generated Section\n\nContent\n<!-- TEST_END -->\n\n## Footer"
		updated, changed, err := injectMarkdown(source, generated, startMarker, endMarker)

		require.NoError(t, err)
		assert.False(t, changed)
		assert.Equal(t, source, updated)
	})

	t.Run("fails when start marker is missing", func(t *testing.T) {
		source := "# Document\n\n<!-- TEST_END -->"
		_, _, err := injectMarkdown(source, generated, startMarker, endMarker)
		assert.ErrorContains(t, err, "marker \"<!-- TEST_START -->\" not found")
	})

	t.Run("fails when end marker is missing", func(t *testing.T) {
		source := "# Document\n\n<!-- TEST_START -->"
		_, _, err := injectMarkdown(source, generated, startMarker, endMarker)
		assert.ErrorContains(t, err, "marker \"<!-- TEST_END -->\" not found")
	})

	t.Run("fails when markers are in wrong order", func(t *testing.T) {
		source := "<!-- TEST_END -->\n<!-- TEST_START -->"
		_, _, err := injectMarkdown(source, generated, startMarker, endMarker)
		assert.ErrorContains(t, err, "invalid marker order")
	})
}

func TestEndToEndDocgenSync(t *testing.T) {
	dir := t.TempDir()
	justDir := filepath.Join(dir, "just")
	require.NoError(t, os.MkdirAll(justDir, 0o750))

	rootJustfile := filepath.Join(dir, "justfile")
	webJustfile := filepath.Join(justDir, "web.just")
	docFile := filepath.Join(dir, "README.md")

	require.NoError(t, os.WriteFile(rootJustfile, []byte("mod web 'just/web.just'\n# Build all\nbuild:\n  echo build"), 0o600))
	require.NoError(t, os.WriteFile(webJustfile, []byte("# Lint frontend\nlint:\n  pnpm lint"), 0o600))
	require.NoError(t, os.WriteFile(docFile, []byte("# Title\n\n<!-- JUST_COMMANDS_START -->\n<!-- JUST_COMMANDS_END -->\n"), 0o600))

	modules, err := parseRootAndModules(rootJustfile)
	require.NoError(t, err)

	markdown := formatModulesMarkdown(modules)
	cleanDocFile := filepath.Clean(docFile)
	content, err := os.ReadFile(cleanDocFile) //nolint:gosec // Test file in temporary directory.
	require.NoError(t, err)

	updated, changed, err := injectMarkdown(string(content), markdown, "<!-- JUST_COMMANDS_START -->", "<!-- JUST_COMMANDS_END -->")
	require.NoError(t, err)
	assert.True(t, changed)

	require.NoError(t, os.WriteFile(cleanDocFile, []byte(updated), 0o600)) //nolint:gosec // Test file in temporary directory.

	// Re-run to verify idempotency
	contentAfter, err := os.ReadFile(cleanDocFile) //nolint:gosec // Test file in temporary directory.
	require.NoError(t, err)
	_, changedAfter, err := injectMarkdown(string(contentAfter), markdown, "<!-- JUST_COMMANDS_START -->", "<!-- JUST_COMMANDS_END -->")
	require.NoError(t, err)
	assert.False(t, changedAfter)
}

func TestRunDocgen(t *testing.T) {
	t.Run("detects drift in check mode", func(t *testing.T) {
		dir := t.TempDir()
		justFile := filepath.Join(dir, "justfile")
		docFile := filepath.Join(dir, "README.md")

		require.NoError(t, os.WriteFile(justFile, []byte("# Root task\nroot:\n  echo 1"), 0o600))
		require.NoError(t, os.WriteFile(docFile, []byte("<!-- TEST_START -->\n<!-- TEST_END -->"), 0o600))

		drift, err := runDocgen(justFile, docFile, "TEST", true)
		assert.True(t, drift)
		assert.ErrorIs(t, err, errDriftDetected)
	})

	t.Run("returns error when target file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		justFile := filepath.Join(dir, "justfile")

		require.NoError(t, os.WriteFile(justFile, []byte("# Root task\nroot:\n  echo 1"), 0o600))

		drift, err := runDocgen(justFile, "/non/existent/doc.md", "TEST", false)
		assert.False(t, drift)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error reading")
	})

	t.Run("returns error when marker is missing in target file", func(t *testing.T) {
		dir := t.TempDir()
		justFile := filepath.Join(dir, "justfile")
		docFile := filepath.Join(dir, "README.md")

		require.NoError(t, os.WriteFile(justFile, []byte("# Root task\nroot:\n  echo 1"), 0o600))
		require.NoError(t, os.WriteFile(docFile, []byte("# No markers here"), 0o600))

		drift, err := runDocgen(justFile, docFile+",", "TEST", false)
		assert.False(t, drift)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error updating")
	})
}
