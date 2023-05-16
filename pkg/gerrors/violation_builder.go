package gerrors

type violationBuilder struct {
	violations []Violation
}

func NewViolationBuilder() violationBuilder {
	return violationBuilder{}
}

func (b *violationBuilder) Add(field, description string) {
	b.violations = append(b.violations, Violation{
		Field:       field,
		Description: description,
	})
}

func (b *violationBuilder) Make() []Violation {
	return b.violations
}

func (b *violationBuilder) IsEmpty() bool {
	return len(b.violations) == 0
}
