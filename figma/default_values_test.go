package figma

import (
	"slices"
	"testing"
)

type Config struct {
	Name  string `default:"guest"`
	Age   int    `default:"30"`
	Admin bool   `default:"true"`
}

type ConfigAdvance struct {
	Name     string   `default:"guest"`
	Age      int      `default:"30"`
	Admin    *bool    `default:"true"`
	Settings Settings // nested struct
	Roles    []Role   // slice of struct (will default each element if not nil)
}

type Settings struct {
	Theme     string `default:"dark"`
	PageLimit int    `default:"20"`
}

type Role struct {
	Name  string `default:"user"`
	Level int    `default:"1"`
}

func TestStructDefaults(t *testing.T) {
	cfg := Config{}

	SetDefaults(&cfg)
	if cfg.Name != "guest" || cfg.Age != 30 || !cfg.Admin {
		t.Errorf("%+v does not have required defaults", cfg)
	}

	cfg = Config{
		Name:  "John Doe",
		Age:   40,
		Admin: false,
	}

	SetDefaults(&cfg)
	// cfd.Admin will be always true for this struct, the bool needs to be a pointer
	if cfg.Name != "John Doe" || cfg.Age != 40 || !cfg.Admin {
		t.Errorf("%+v does not have expected values", cfg)
	}

	cfga := ConfigAdvance{}
	settings := Settings{
		Theme:     "dark",
		PageLimit: 20,
	}
	roles := []Role{}

	SetDefaults(&cfga)
	if cfga.Name != "guest" || cfga.Age != 30 || !*cfga.Admin || cfga.Settings != settings || !slices.Equal(cfga.Roles, roles) {
		t.Errorf("%+v does not have required defaults", cfga)
	}

	falseVal := false
	settings = Settings{
		Theme:     "light",
		PageLimit: 50,
	}
	roles = []Role{
		{},              // empty role should get defaults
		{Name: "admin"}, // only partial data
	}

	cfga = ConfigAdvance{
		Name:     "John Doe",
		Age:      40,
		Admin:    &falseVal,
		Settings: settings,
		Roles:    roles,
	}

	SetDefaults(&cfga)
	if cfga.Name != "John Doe" || cfga.Age != 40 || *cfga.Admin || cfga.Settings != settings || !slices.Equal(cfga.Roles, roles) {
		t.Errorf("%+v does not have required defaults", cfga)
	}
}
