// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package controller

import (
	"context"
	"github.com/aquaproj/aqua/v2/pkg/cargo"
	"github.com/aquaproj/aqua/v2/pkg/checksum"
	"github.com/aquaproj/aqua/v2/pkg/config"
	"github.com/aquaproj/aqua/v2/pkg/config-finder"
	"github.com/aquaproj/aqua/v2/pkg/config-reader"
	"github.com/aquaproj/aqua/v2/pkg/controller/allowpolicy"
	"github.com/aquaproj/aqua/v2/pkg/controller/cp"
	"github.com/aquaproj/aqua/v2/pkg/controller/denypolicy"
	exec2 "github.com/aquaproj/aqua/v2/pkg/controller/exec"
	"github.com/aquaproj/aqua/v2/pkg/controller/generate"
	"github.com/aquaproj/aqua/v2/pkg/controller/generate-registry"
	"github.com/aquaproj/aqua/v2/pkg/controller/generate/output"
	"github.com/aquaproj/aqua/v2/pkg/controller/info"
	"github.com/aquaproj/aqua/v2/pkg/controller/initcmd"
	"github.com/aquaproj/aqua/v2/pkg/controller/initpolicy"
	"github.com/aquaproj/aqua/v2/pkg/controller/install"
	"github.com/aquaproj/aqua/v2/pkg/controller/list"
	"github.com/aquaproj/aqua/v2/pkg/controller/remove"
	"github.com/aquaproj/aqua/v2/pkg/controller/update"
	"github.com/aquaproj/aqua/v2/pkg/controller/updateaqua"
	"github.com/aquaproj/aqua/v2/pkg/controller/updatechecksum"
	"github.com/aquaproj/aqua/v2/pkg/controller/which"
	"github.com/aquaproj/aqua/v2/pkg/cosign"
	"github.com/aquaproj/aqua/v2/pkg/download"
	"github.com/aquaproj/aqua/v2/pkg/exec"
	"github.com/aquaproj/aqua/v2/pkg/fuzzyfinder"
	"github.com/aquaproj/aqua/v2/pkg/github"
	"github.com/aquaproj/aqua/v2/pkg/install-registry"
	"github.com/aquaproj/aqua/v2/pkg/installpackage"
	"github.com/aquaproj/aqua/v2/pkg/link"
	"github.com/aquaproj/aqua/v2/pkg/policy"
	"github.com/aquaproj/aqua/v2/pkg/runtime"
	"github.com/aquaproj/aqua/v2/pkg/slsa"
	"github.com/aquaproj/aqua/v2/pkg/unarchive"
	"github.com/aquaproj/aqua/v2/pkg/versiongetter"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/go-osenv/osenv"
	"io"
	"net/http"
)

// Injectors from wire.go:

func InitializeListCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *list.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	controller := list.NewController(configFinder, configReader, installer, fs)
	return controller
}

func InitializeGenerateRegistryCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, stdout io.Writer) *genrgst.Controller {
	fs := afero.NewOsFs()
	repositoriesService := github.New(ctx)
	outputter := output.New(stdout, fs)
	client := cargo.NewClient(httpClient)
	controller := genrgst.NewController(fs, repositoriesService, outputter, client)
	return controller
}

func InitializeInitCommandController(ctx context.Context, param *config.Param) *initcmd.Controller {
	repositoriesService := github.New(ctx)
	fs := afero.NewOsFs()
	controller := initcmd.New(repositoriesService, fs)
	return controller
}

func InitializeInitPolicyCommandController(ctx context.Context) *initpolicy.Controller {
	fs := afero.NewOsFs()
	controller := initpolicy.New(fs)
	return controller
}

func InitializeGenerateCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *generate.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	fuzzyfinderFinder := fuzzyfinder.New()
	client := cargo.NewClient(httpClient)
	cargoVersionGetter := versiongetter.NewCargo(client)
	gitHubTagVersionGetter := versiongetter.NewGitHubTag(repositoriesService)
	gitHubReleaseVersionGetter := versiongetter.NewGitHubRelease(repositoriesService)
	generalVersionGetter := versiongetter.NewGeneralVersionGetter(cargoVersionGetter, gitHubTagVersionGetter, gitHubReleaseVersionGetter)
	fuzzyGetter := versiongetter.NewFuzzy(fuzzyfinderFinder, generalVersionGetter)
	controller := generate.New(configFinder, configReader, installer, repositoriesService, fs, fuzzyfinderFinder, fuzzyGetter)
	return controller
}

func InitializeInstallCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *install.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	linker := link.New()
	checksumDownloaderImpl := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	unarchiver := unarchive.New(executor, fs)
	goInstallInstallerImpl := installpackage.NewGoInstallInstallerImpl(executor)
	goBuildInstallerImpl := installpackage.NewGoBuildInstallerImpl(executor)
	cargoPackageInstallerImpl := installpackage.NewCargoPackageInstallerImpl(executor, fs)
	installerImpl := installpackage.New(param, downloader, rt, fs, linker, checksumDownloaderImpl, calculator, unarchiver, verifier, slsaVerifier, goInstallInstallerImpl, goBuildInstallerImpl, cargoPackageInstallerImpl)
	validatorImpl := policy.NewValidator(param, fs)
	configFinderImpl := policy.NewConfigFinder(fs)
	configReaderImpl := policy.NewConfigReader(fs)
	policyReader := policy.NewReader(fs, validatorImpl, configFinderImpl, configReaderImpl)
	controller := install.New(param, configFinder, configReader, installer, installerImpl, fs, rt, policyReader, configFinderImpl)
	return controller
}

func InitializeWhichCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *which.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	osEnv := osenv.New()
	linker := link.New()
	controller := which.New(param, configFinder, configReader, installer, rt, osEnv, fs, linker)
	return controller
}

func InitializeExecCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *exec2.Controller {
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	fs := afero.NewOsFs()
	linker := link.New()
	checksumDownloaderImpl := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	executor := exec.New()
	unarchiver := unarchive.New(executor, fs)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	goInstallInstallerImpl := installpackage.NewGoInstallInstallerImpl(executor)
	goBuildInstallerImpl := installpackage.NewGoBuildInstallerImpl(executor)
	cargoPackageInstallerImpl := installpackage.NewCargoPackageInstallerImpl(executor, fs)
	installerImpl := installpackage.New(param, downloader, rt, fs, linker, checksumDownloaderImpl, calculator, unarchiver, verifier, slsaVerifier, goInstallInstallerImpl, goBuildInstallerImpl, cargoPackageInstallerImpl)
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	osEnv := osenv.New()
	controller := which.New(param, configFinder, configReader, installer, rt, osEnv, fs, linker)
	validatorImpl := policy.NewValidator(param, fs)
	configFinderImpl := policy.NewConfigFinder(fs)
	configReaderImpl := policy.NewConfigReader(fs)
	policyReader := policy.NewReader(fs, validatorImpl, configFinderImpl, configReaderImpl)
	execController := exec2.New(param, installerImpl, controller, executor, osEnv, fs, policyReader, configFinderImpl)
	return execController
}

func InitializeUpdateAquaCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *updateaqua.Controller {
	fs := afero.NewOsFs()
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	linker := link.New()
	checksumDownloaderImpl := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	executor := exec.New()
	unarchiver := unarchive.New(executor, fs)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	goInstallInstallerImpl := installpackage.NewGoInstallInstallerImpl(executor)
	goBuildInstallerImpl := installpackage.NewGoBuildInstallerImpl(executor)
	cargoPackageInstallerImpl := installpackage.NewCargoPackageInstallerImpl(executor, fs)
	installerImpl := installpackage.New(param, downloader, rt, fs, linker, checksumDownloaderImpl, calculator, unarchiver, verifier, slsaVerifier, goInstallInstallerImpl, goBuildInstallerImpl, cargoPackageInstallerImpl)
	controller := updateaqua.New(param, fs, rt, repositoriesService, installerImpl)
	return controller
}

func InitializeCopyCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *cp.Controller {
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	fs := afero.NewOsFs()
	linker := link.New()
	checksumDownloaderImpl := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	executor := exec.New()
	unarchiver := unarchive.New(executor, fs)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	goInstallInstallerImpl := installpackage.NewGoInstallInstallerImpl(executor)
	goBuildInstallerImpl := installpackage.NewGoBuildInstallerImpl(executor)
	cargoPackageInstallerImpl := installpackage.NewCargoPackageInstallerImpl(executor, fs)
	installerImpl := installpackage.New(param, downloader, rt, fs, linker, checksumDownloaderImpl, calculator, unarchiver, verifier, slsaVerifier, goInstallInstallerImpl, goBuildInstallerImpl, cargoPackageInstallerImpl)
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	osEnv := osenv.New()
	controller := which.New(param, configFinder, configReader, installer, rt, osEnv, fs, linker)
	validatorImpl := policy.NewValidator(param, fs)
	configFinderImpl := policy.NewConfigFinder(fs)
	configReaderImpl := policy.NewConfigReader(fs)
	policyReader := policy.NewReader(fs, validatorImpl, configFinderImpl, configReaderImpl)
	installController := install.New(param, configFinder, configReader, installer, installerImpl, fs, rt, policyReader, configFinderImpl)
	cpController := cp.New(param, installerImpl, fs, rt, controller, installController, policyReader)
	return cpController
}

func InitializeUpdateChecksumCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *updatechecksum.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	checksumDownloaderImpl := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	controller := updatechecksum.New(param, configFinder, configReader, installer, fs, rt, checksumDownloaderImpl, downloader, gitHubContentFileDownloader)
	return controller
}

func InitializeUpdateCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *update.Controller {
	repositoriesService := github.New(ctx)
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	fuzzyfinderFinder := fuzzyfinder.New()
	client := cargo.NewClient(httpClient)
	cargoVersionGetter := versiongetter.NewCargo(client)
	gitHubTagVersionGetter := versiongetter.NewGitHubTag(repositoriesService)
	gitHubReleaseVersionGetter := versiongetter.NewGitHubRelease(repositoriesService)
	generalVersionGetter := versiongetter.NewGeneralVersionGetter(cargoVersionGetter, gitHubTagVersionGetter, gitHubReleaseVersionGetter)
	fuzzyGetter := versiongetter.NewFuzzy(fuzzyfinderFinder, generalVersionGetter)
	osEnv := osenv.New()
	linker := link.New()
	controller := which.New(param, configFinder, configReader, installer, rt, osEnv, fs, linker)
	updateController := update.New(param, repositoriesService, configFinder, configReader, installer, fs, rt, fuzzyGetter, fuzzyfinderFinder, controller)
	return updateController
}

func InitializeAllowPolicyCommandController(ctx context.Context, param *config.Param) *allowpolicy.Controller {
	fs := afero.NewOsFs()
	configFinderImpl := policy.NewConfigFinder(fs)
	validatorImpl := policy.NewValidator(param, fs)
	controller := allowpolicy.New(fs, configFinderImpl, validatorImpl)
	return controller
}

func InitializeDenyPolicyCommandController(ctx context.Context, param *config.Param) *denypolicy.Controller {
	fs := afero.NewOsFs()
	configFinderImpl := policy.NewConfigFinder(fs)
	validatorImpl := policy.NewValidator(param, fs)
	controller := denypolicy.New(fs, configFinderImpl, validatorImpl)
	return controller
}

func InitializeInfoCommandController(ctx context.Context, param *config.Param, rt *runtime.Runtime) *info.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	controller := info.New(fs, configFinder, rt)
	return controller
}

func InitializeRemoveCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *remove.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifier := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifier := slsa.New(downloader, fs, executorImpl)
	installer := registry.New(param, gitHubContentFileDownloader, fs, rt, verifier, slsaVerifier)
	fuzzyfinderFinder := fuzzyfinder.New()
	osEnv := osenv.New()
	linker := link.New()
	controller := which.New(param, configFinder, configReader, installer, rt, osEnv, fs, linker)
	removeController := remove.New(param, fs, rt, configFinder, configReader, installer, fuzzyfinderFinder, controller)
	return removeController
}
