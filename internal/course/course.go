package course

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/manifoldco/promptui/list"
	"github.com/pathbird/pbauthor/internal/graphql"
	"github.com/pkg/errors"
	"strings"
)

var (
	promptCourseTemplate        = newSelectTemplate("Course", "🎓", ".Course.Name")
	promptCodexCategoryTemplate = newSelectTemplate("Codex Category", "😻", ".Name")
)

func PromptCourse(courses []graphql.CourseEdge) (*graphql.Course, error) {
	if len(courses) == 0 {
		return nil, errors.New("no courses to select from")
	}
	searcher := newSubstrSearcher(func(i int) string {
		return courses[i].Course.Name
	})
	prompt := promptui.Select{
		Label:     "Course",
		Items:     courses,
		Templates: promptCourseTemplate,
		Searcher:  searcher,
	}
	n, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &courses[n].Course, err
}

func PromptCodexCategory(cats []graphql.CodexCategory) (*graphql.CodexCategory, error) {
	// Meow!
	if len(cats) == 0 {
		return nil, errors.New("course has categories")
	}
	if len(cats) == 1 {
		return &cats[0], nil
	}

	searcher := newSubstrSearcher(func(i int) string {
		return cats[i].Name
	})
	prompt := promptui.Select{
		Label:     "Codex Category",
		Items:     cats,
		Templates: promptCodexCategoryTemplate,
		Searcher:  searcher,
	}
	n, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &cats[n], nil
}

// Construct a promptui.SelectTemplates.
// This helper function exists to have uniformity across different select prompts.
func newSelectTemplate(label string, emoji string, attr string) *promptui.SelectTemplates {
	// Note: emoji's (in my terminal) render as two monospace characters wide
	return &promptui.SelectTemplates{
		Label:    fmt.Sprintf(`%s {{ "%s:" | bold }}`, emoji, label),
		Active:   fmt.Sprintf(` > {{ %s | cyan | underline }}`, attr),
		Inactive: fmt.Sprintf("   {{ %s }}", attr),
		Selected: fmt.Sprintf("%s {{ %s | bold }}", emoji, attr),
	}
}

// Create a new searcher function that simply checks if the query is a substring.
// The resolver argument should resolve the index of the item to a string.
func newSubstrSearcher(resolver func(int) string) list.Searcher {
	return func(input string, index int) bool {
		val := resolver(index)
		return strings.Contains(strings.ToLower(val), strings.ToLower(input))
	}
}
