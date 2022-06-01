package which

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/aquaproj/aqua/pkg/config"
	finder "github.com/aquaproj/aqua/pkg/config-finder"
	reader "github.com/aquaproj/aqua/pkg/config-reader"
	registry "github.com/aquaproj/aqua/pkg/install-registry"
	"github.com/aquaproj/aqua/pkg/link"
	"github.com/aquaproj/aqua/pkg/runtime"
	"github.com/aquaproj/aqua/pkg/validate"
	constraint "github.com/aquaproj/aqua/pkg/expr"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/go-osenv/osenv"
	"github.com/suzuki-shunsuke/logrus-error/logerr"
)

type controller struct {
	stdout            io.Writer
	rootDir           string
	configFinder      finder.ConfigFinder
	configReader      reader.ConfigReader
	registryInstaller registry.Installer
	runtime           *runtime.Runtime
	osenv             osenv.OSEnv
	fs                afero.Fs
	linker            link.Linker
}

func (ctrl *controller) Which(ctx context.Context, param *config.Param, exeName string, logE *logrus.Entry) (*Which, error) {
	for _, cfgFilePath := range ctrl.configFinder.Finds(param.PWD, param.ConfigFilePath) {
		pkg, pkgInfo, file, err := ctrl.findExecFile(ctx, cfgFilePath, exeName, logE)
		if err != nil {
			return nil, err
		}
		if pkg != nil {
			return ctrl.whichFile(pkg, pkgInfo, file)
		}
	}

	for _, cfgFilePath := range param.GlobalConfigFilePaths {
		if _, err := ctrl.fs.Stat(cfgFilePath); err != nil {
			continue
		}
		pkg, pkgInfo, file, err := ctrl.findExecFile(ctx, cfgFilePath, exeName, logE)
		if err != nil {
			return nil, err
		}
		if pkg != nil {
			return ctrl.whichFile(pkg, pkgInfo, file)
		}
	}

	if exePath := ctrl.lookPath(ctrl.osenv.Getenv("PATH"), exeName); exePath != "" {
		return &Which{
			ExePath: exePath,
		}, nil
	}
	return nil, logerr.WithFields(errCommandIsNotFound, logrus.Fields{ //nolint:wrapcheck
		"exe_name": exeName,
	})
}

func (ctrl *controller) whichFileGo(pkg *config.Package, pkgInfo *config.PackageInfo, file *config.File) (*Which, error) {
	return &Which{
		Package: pkg,
		PkgInfo: pkgInfo,
		File:    file,
		ExePath: filepath.Join(ctrl.rootDir, "pkgs", pkgInfo.GetType(), "github.com", pkgInfo.RepoOwner, pkgInfo.RepoName, pkg.Version, "bin", file.Name),
	}, nil
}

func (ctrl *controller) whichFile(pkg *config.Package, pkgInfo *config.PackageInfo, file *config.File) (*Which, error) {
	if pkgInfo.Type == "go" {
		return ctrl.whichFileGo(pkg, pkgInfo, file)
	}
	fileSrc, err := pkgInfo.GetFileSrc(pkg, file, ctrl.runtime)
	if err != nil {
		return nil, fmt.Errorf("get file_src: %w", err)
	}
	pkgPath, err := pkgInfo.GetPkgPath(ctrl.rootDir, pkg, ctrl.runtime)
	if err != nil {
		return nil, fmt.Errorf("get pkg install path: %w", err)
	}
	return &Which{
		Package: pkg,
		PkgInfo: pkgInfo,
		File:    file,
		ExePath: filepath.Join(pkgPath, fileSrc),
	}, nil
}

func (ctrl *controller) findExecFile(ctx context.Context, cfgFilePath, exeName string, logE *logrus.Entry) (*config.Package, *config.PackageInfo, *config.File, error) {
	cfg := &config.Config{}
	if err := ctrl.configReader.Read(cfgFilePath, cfg); err != nil {
		return nil, nil, nil, err //nolint:wrapcheck
	}
	if err := validate.Config(cfg); err != nil {
		return nil, nil, nil, fmt.Errorf("configuration is invalid: %w", err)
	}

	registryContents, err := ctrl.registryInstaller.InstallRegistries(ctx, cfg, cfgFilePath, logE)
	if err != nil {
		return nil, nil, nil, err //nolint:wrapcheck
	}
	for _, pkg := range cfg.Packages {
		if pkgInfo, file := ctrl.findExecFileFromPkg(registryContents, exeName, pkg, logE); pkgInfo != nil {
			return pkg, pkgInfo, file, nil
		}
	}
	return nil, nil, nil, nil
}

func (ctrl *controller) findExecFileFromPkg(registries map[string]*config.RegistryContent, exeName string, pkg *config.Package, logE *logrus.Entry) (*config.PackageInfo, *config.File) {
	logE = logE.WithFields(logrus.Fields{
		"registry_name": pkg.Registry,
		"package_name":  pkg.Name,
	})
	registry, ok := registries[pkg.Registry]
	if !ok {
		logE.Warn("registry isn't found")
		return nil, nil
	}

	m, err := registry.PackageInfos.ToMap()
	if err != nil {
		logerr.WithError(logE, err).Warn("registry is invalid")
		return nil, nil
	}

	pkgInfo, ok := m[pkg.Name]
	if !ok {
		logE.Warn("package isn't found")
		return nil, nil
	}

	pkgInfo, err = pkgInfo.Override(pkg.Version, ctrl.runtime)
	if err != nil {
		logerr.WithError(logE, err).Warn("version constraint is invalid")
		return nil, nil
	}

	if pkgInfo.SupportedIf != nil {
		supported, err := constraint.EvaluateSupportedIf(pkgInfo.SupportedIf, ctrl.runtime)
		if err != nil {
			logerr.WithError(logE, err).WithField("supported_if", *pkgInfo.SupportedIf).Error("check if the package is supported")
			return nil, nil
		}
		if !supported {
			logE.WithField("supported_if", *pkgInfo.SupportedIf).Debug("the package isn't supported on this environment")
			return nil, nil
		}
	}

	for _, file := range pkgInfo.GetFiles() {
		if file.Name == exeName {
			return pkgInfo, file
		}
	}
	return nil, nil
}
