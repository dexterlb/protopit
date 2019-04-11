package builder

func Build(variants []string, contentDir string) {
	sites := make(map[string]*Site)
	for _, variant := range variants {
		s := Init(variant, contentDir)
		s.Clean()
		s.RenderCss()
		s.LoadPages()
		s.AllVariants = sites
		sites[variant] = s
	}
	for _, s := range sites {
		s.Render()
	}
}
