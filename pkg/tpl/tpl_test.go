package tpl_test

import (
	"testing"

	"github.com/jrnd-io/jrv2/pkg/tpl"
)

const (
	templateString = "Hello, {{.Name}}!"
)

func TestNew(t *testing.T) {
	// Test case for successful template creation
	t.Run("successful creation", func(t *testing.T) {
		name := "test_success"
		tmpl := templateString
		fmap := map[string]interface{}{}
		ctx := struct{ Name string }{"World"}

		tpl, err := tpl.New(name, tmpl, fmap, ctx)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if tpl.Template == nil {
			t.Fatal("Expected Template to be non-nil")
		}
		if tpl.Context != ctx {
			t.Fatalf("Expected Context to be %v, got %v", ctx, tpl.Context)
		}
	})

	// Test case for invalid template
	t.Run("invalid template", func(t *testing.T) {
		name := "test_invalid_template"
		tmpl := "Hello, {{.Name}!" // Missing closing brace
		fmap := map[string]interface{}{}
		ctx := struct{ Name string }{"World"}

		_, err := tpl.New(name, tmpl, fmap, ctx)

		if err == nil {
			t.Fatal("Expected an error, got nil")
		}
	})
}

func TestExecute(t *testing.T) {
	name := "test_execute"
	tmpl := templateString
	fmap := map[string]interface{}{}
	ctx := struct{ Name string }{"World"}

	tpl, err := tpl.New(name, tmpl, fmap, ctx)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	result := tpl.Execute()
	expected := "Hello, World!"

	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}

func TestExecuteWith(t *testing.T) {
	name := "test_executewith"
	tmpl := templateString
	fmap := map[string]interface{}{}
	ctx := struct{ Name string }{"World"}

	tpl, err := tpl.New(name, tmpl, fmap, ctx)
	if err != nil {
		t.Fatalf("Failed to create template: %v", err)
	}

	newCtx := struct{ Name string }{"Go"}
	result := tpl.ExecuteWith(newCtx)
	expected := "Hello, Go!"

	if result != expected {
		t.Fatalf("Expected %q, got %q", expected, result)
	}
}
