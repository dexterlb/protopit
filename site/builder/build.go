package builder

func BuildVariant(variant string, contentDir string) {
	s := Init(variant, contentDir)
	s.RenderCss()
	s.LoadPages()
	s.Render()
}
