package image

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/types"
	perrors "github.com/pkg/errors"
)

func manifestSchema2FromManifestList(ctx context.Context, sys *types.SystemContext, src types.ImageSource, manblob []byte) (genericManifest, error) {
	list, err := manifest.Schema2ListFromManifest(manblob)
	if err != nil {
		return nil, perrors.Wrapf(err, "parsing schema2 manifest list")
	}
	targetManifestDigest, err := list.ChooseInstance(sys)
	if err != nil {
		return nil, perrors.Wrapf(err, "choosing image instance")
	}
	manblob, mt, err := src.GetManifest(ctx, &targetManifestDigest)
	if err != nil {
		return nil, perrors.Wrapf(err, "fetching target platform image selected from manifest list")
	}

	matches, err := manifest.MatchesDigest(manblob, targetManifestDigest)
	if err != nil {
		return nil, perrors.Wrap(err, "computing manifest digest")
	}
	if !matches {
		return nil, fmt.Errorf("Image manifest does not match selected manifest digest %s", targetManifestDigest)
	}

	return manifestInstanceFromBlob(ctx, sys, src, manblob, mt)
}
