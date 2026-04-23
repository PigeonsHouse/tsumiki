package helper

import (
	"tsumiki/media"
	"tsumiki/schema"
)

func ResolveTsumikiURLs(t *schema.Tsumiki, svc media.MediaService) {
	t.User.AvatarUrl = svc.ResolveURL(t.User.AvatarUrl)
	if t.Thumbnail != nil {
		t.Thumbnail.Url = svc.ResolveURL(t.Thumbnail.Url)
	}
	if t.Work != nil {
		ResolveWorkURLs(t.Work, svc)
	}
}

func ResolveWorkURLs(w *schema.Work, svc media.MediaService) {
	w.Owner.AvatarUrl = svc.ResolveURL(w.Owner.AvatarUrl)
	if w.Thumbnail != nil {
		w.Thumbnail.Url = svc.ResolveURL(w.Thumbnail.Url)
	}
}

func ResolveBlockViewsURLs(blocks []schema.TsumikiBlockView, svc media.MediaService) {
	for i := range blocks {
		for j := range blocks[i].Medias {
			blocks[i].Medias[j].Url = svc.ResolveURL(blocks[i].Medias[j].Url)
		}
	}
}
