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
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/go-osenv/osenv"
	"io"
	"net/http"
)

// Injectors from wire.go:

func InitializeListCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *list.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReaderImpl := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifierImpl := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifierImpl := slsa.New(downloader, fs, executorImpl)
	installerImpl := registry.New(param, gitHubContentFileDownloader, fs, rt, verifierImpl, slsaVerifierImpl)
	controller := list.NewController(configFinder, configReaderImpl, installerImpl, fs)
	return controller
}

func InitializeGenerateRegistryCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, stdout io.Writer) *genrgst.Controller {
	fs := afero.NewOsFs()
	repositoriesService := github.New(ctx)
	outputter := output.New(stdout, fs)
	clientImpl := cargo.NewClientImpl(httpClient)
	controller := genrgst.NewController(fs, repositoriesService, outputter, clientImpl)
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
	configReaderImpl := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifierImpl := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifierImpl := slsa.New(downloader, fs, executorImpl)
	installerImpl := registry.New(param, gitHubContentFileDownloader, fs, rt, verifierImpl, slsaVerifierImpl)
	fuzzyfinderFinder := fuzzyfinder.New()
	versionSelector := fuzzyfinder.NewVersionSelector()
	clientImpl := cargo.NewClientImpl(httpClient)
	controller := generate.New(configFinder, configReaderImpl, installerImpl, repositoriesService, fs, fuzzyfinderFinder, versionSelector, clientImpl)
	return controller
}

func InitializeInstallCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *install.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReaderImpl := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifierImpl := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifierImpl := slsa.New(downloader, fs, executorImpl)
	installerImpl := registry.New(param, gitHubContentFileDownloader, fs, rt, verifierImpl, slsaVerifierImpl)
	linker := link.New()
	checksumDownloaderImpl := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	calculator := checksum.NewCalculator()
	unarchiverImpl := unarchive.New(executor, fs)
	checker := policy.NewChecker(param)
	goInstallInstallerImpl := installpackage.NewGoInstallInstallerImpl(executor)
	goBuildInstallerImpl := installpackage.NewGoBuildInstallerImpl(executor)
	cargoPackageInstallerImpl := installpackage.NewCargoPackageInstallerImpl(executor, fs)
	installpackageInstallerImpl := installpackage.New(param, downloader, rt, fs, linker, checksumDownloaderImpl, calculator, unarchiverImpl, checker, verifierImpl, slsaVerifierImpl, goInstallInstallerImpl, goBuildInstallerImpl, cargoPackageInstallerImpl)
	validatorImpl := policy.NewValidator(param, fs)
	configFinderImpl := policy.NewConfigFinder(fs)
	policyConfigReaderImpl := policy.NewConfigReader(fs)
	readerImpl := policy.NewReader(fs, validatorImpl, configFinderImpl, policyConfigReaderImpl)
	controller := install.New(param, configFinder, configReaderImpl, installerImpl, installpackageInstallerImpl, fs, rt, readerImpl, configFinderImpl)
	return controller
}

func InitializeWhichCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *which.ControllerImpl {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReaderImpl := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifierImpl := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifierImpl := slsa.New(downloader, fs, executorImpl)
	installerImpl := registry.New(param, gitHubContentFileDownloader, fs, rt, verifierImpl, slsaVerifierImpl)
	osEnv := osenv.New()
	linker := link.New()
	controllerImpl := which.New(param, configFinder, configReaderImpl, installerImpl, rt, osEnv, fs, linker)
	return controllerImpl
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
	unarchiverImpl := unarchive.New(executor, fs)
	checker := policy.NewChecker(param)
	verifierImpl := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifierImpl := slsa.New(downloader, fs, executorImpl)
	goInstallInstallerImpl := installpackage.NewGoInstallInstallerImpl(executor)
	goBuildInstallerImpl := installpackage.NewGoBuildInstallerImpl(executor)
	cargoPackageInstallerImpl := installpackage.NewCargoPackageInstallerImpl(executor, fs)
	installerImpl := installpackage.New(param, downloader, rt, fs, linker, checksumDownloaderImpl, calculator, unarchiverImpl, checker, verifierImpl, slsaVerifierImpl, goInstallInstallerImpl, goBuildInstallerImpl, cargoPackageInstallerImpl)
	configFinder := finder.NewConfigFinder(fs)
	configReaderImpl := reader.New(fs, param)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	registryInstallerImpl := registry.New(param, gitHubContentFileDownloader, fs, rt, verifierImpl, slsaVerifierImpl)
	osEnv := osenv.New()
	controllerImpl := which.New(param, configFinder, configReaderImpl, registryInstallerImpl, rt, osEnv, fs, linker)
	validatorImpl := policy.NewValidator(param, fs)
	configFinderImpl := policy.NewConfigFinder(fs)
	policyConfigReaderImpl := policy.NewConfigReader(fs)
	readerImpl := policy.NewReader(fs, validatorImpl, configFinderImpl, policyConfigReaderImpl)
	controller := exec2.New(param, installerImpl, controllerImpl, executor, osEnv, fs, readerImpl, configFinderImpl)
	return controller
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
	unarchiverImpl := unarchive.New(executor, fs)
	checker := policy.NewChecker(param)
	verifierImpl := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifierImpl := slsa.New(downloader, fs, executorImpl)
	goInstallInstallerImpl := installpackage.NewGoInstallInstallerImpl(executor)
	goBuildInstallerImpl := installpackage.NewGoBuildInstallerImpl(executor)
	cargoPackageInstallerImpl := installpackage.NewCargoPackageInstallerImpl(executor, fs)
	installerImpl := installpackage.New(param, downloader, rt, fs, linker, checksumDownloaderImpl, calculator, unarchiverImpl, checker, verifierImpl, slsaVerifierImpl, goInstallInstallerImpl, goBuildInstallerImpl, cargoPackageInstallerImpl)
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
	unarchiverImpl := unarchive.New(executor, fs)
	checker := policy.NewChecker(param)
	verifierImpl := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifierImpl := slsa.New(downloader, fs, executorImpl)
	goInstallInstallerImpl := installpackage.NewGoInstallInstallerImpl(executor)
	goBuildInstallerImpl := installpackage.NewGoBuildInstallerImpl(executor)
	cargoPackageInstallerImpl := installpackage.NewCargoPackageInstallerImpl(executor, fs)
	installerImpl := installpackage.New(param, downloader, rt, fs, linker, checksumDownloaderImpl, calculator, unarchiverImpl, checker, verifierImpl, slsaVerifierImpl, goInstallInstallerImpl, goBuildInstallerImpl, cargoPackageInstallerImpl)
	configFinder := finder.NewConfigFinder(fs)
	configReaderImpl := reader.New(fs, param)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	registryInstallerImpl := registry.New(param, gitHubContentFileDownloader, fs, rt, verifierImpl, slsaVerifierImpl)
	osEnv := osenv.New()
	controllerImpl := which.New(param, configFinder, configReaderImpl, registryInstallerImpl, rt, osEnv, fs, linker)
	validatorImpl := policy.NewValidator(param, fs)
	configFinderImpl := policy.NewConfigFinder(fs)
	policyConfigReaderImpl := policy.NewConfigReader(fs)
	readerImpl := policy.NewReader(fs, validatorImpl, configFinderImpl, policyConfigReaderImpl)
	controller := install.New(param, configFinder, configReaderImpl, registryInstallerImpl, installerImpl, fs, rt, readerImpl, configFinderImpl)
	cpController := cp.New(param, installerImpl, fs, rt, controllerImpl, controller, readerImpl, configFinderImpl)
	return cpController
}

func InitializeUpdateChecksumCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *updatechecksum.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReaderImpl := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifierImpl := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifierImpl := slsa.New(downloader, fs, executorImpl)
	installerImpl := registry.New(param, gitHubContentFileDownloader, fs, rt, verifierImpl, slsaVerifierImpl)
	checksumDownloaderImpl := download.NewChecksumDownloader(repositoriesService, rt, httpDownloader)
	controller := updatechecksum.New(param, configFinder, configReaderImpl, installerImpl, fs, rt, checksumDownloaderImpl, downloader, gitHubContentFileDownloader)
	return controller
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
	configReaderImpl := reader.New(fs, param)
	repositoriesService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	gitHubContentFileDownloader := download.NewGitHubContentFileDownloader(repositoriesService, httpDownloader)
	executor := exec.New()
	downloader := download.NewDownloader(repositoriesService, httpDownloader)
	verifierImpl := cosign.NewVerifier(executor, fs, downloader, param)
	executorImpl := slsa.NewExecutor(executor, param)
	slsaVerifierImpl := slsa.New(downloader, fs, executorImpl)
	installerImpl := registry.New(param, gitHubContentFileDownloader, fs, rt, verifierImpl, slsaVerifierImpl)
	fuzzyfinderFinder := fuzzyfinder.New()
	controller := remove.New(param, fs, rt, configFinder, configReaderImpl, installerImpl, fuzzyfinderFinder)
	return controller
}
