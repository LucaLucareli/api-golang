package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Users struct {
	ent.Schema
}

func (Users) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("name"),
		field.String("email").Unique(),
	}
}

