package builder

func BuildVariant(variant string, contentDir string) {
	s := Init(variant, contentDir)
	s.LoadPages()
	s.Render()
}
