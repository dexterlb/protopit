package builder

import "github.com/DexterLB/protopit/site/translator"

func Build(variants []string, contentDir string, translatorFile string) {
	tran, err := translator.Load(translatorFile)
	noerr("cannot load translator", err)

	sites := make(map[string]*Site)
	for _, variant := range variants {
		s := Init(variant, contentDir, tran)
		s.Clean()
		s.RenderCss()
		s.LoadPages()
		s.AllVariants = sites
		sites[variant] = s
	}
	for _, s := range sites {
		s.Render()
	}

	noerr("cannot store translator", tran.Store(translatorFile))
}
