package builder

func BuildVariant(variant string, contentDir string) {
	s := Init(variant, contentDir)
	s.Clean()
	s.RenderCss()
	s.LoadPages()
	s.Render()
}
