package symlink

import (
	"acmp/models"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func createSymlink(target, linkPath string) error {
	exists, isSymlink, err := inspectPath(linkPath)
	if err != nil {
		return err
	}

	if isSymlink {
		return nil
	}

	if exists {
		return fmt.Errorf("path exists but is not a symlink: %s", linkPath)
	}

	err = os.Symlink(target, linkPath)
	if err != nil {
		return fmt.Errorf("creating symlink %s -> %s: %w", linkPath, target, err)
	}
	return nil
}

func DeleteSymlink(linkpath string) error {
	exists, isSymlink, err := inspectPath(linkpath)
	if err != nil {
		return err
	}

	if !exists || !isSymlink {
		return nil
	}

	err = os.Remove(linkpath)
	if err != nil {
		return fmt.Errorf("removing symlink %s: %w", linkpath, err)
	}
	return nil
}

func inspectPath(path string) (bool, bool, error) {
	info, err := os.Lstat(path)

	if os.IsNotExist(err) {
		return false, false, nil
	}

	if err != nil {
		return false, false, fmt.Errorf("lstat: %w", err)
	}

	isSymlink := info.Mode()&os.ModeSymlink != 0
	return true, isSymlink, nil
}

func BuildSymlinkPath(category, modName, assettoPath string) string {
	return filepath.Join(assettoPath, "content", strings.ToLower(category), modName)
}

func ReconcileSymlinks(mods []models.Mod, assettoPath string) error {
	g, ctx := errgroup.WithContext(context.Background())

	for _, mod := range mods {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			linkPath := BuildSymlinkPath(mod.Category, mod.Name, assettoPath)

			var err error
			if mod.Active {
				err = ensureSymlink(mod.Dir, linkPath)
			} else {
				err = DeleteSymlink(linkPath)
			}

			if err != nil {
				fmt.Printf("warning: symlink failed for %s: %v\n", linkPath, err)
				return nil
			}
			return nil
		})
	}
	return g.Wait()
}

func ensureSymlink(modDir, linkPath string) error {
	_, err := os.Stat(modDir)
	if os.IsNotExist(err) {
		return DeleteSymlink(linkPath)
	}
	if err != nil {
		return fmt.Errorf("checking mod folder %s: %w", modDir, err)
	}

	return createSymlink(modDir, linkPath)
}
