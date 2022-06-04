// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package controller

import (
	"context"
	"github.com/aquaproj/aqua/pkg/config"
	"github.com/aquaproj/aqua/pkg/config-finder"
	"github.com/aquaproj/aqua/pkg/config-reader"
	exec2 "github.com/aquaproj/aqua/pkg/controller/exec"
	"github.com/aquaproj/aqua/pkg/controller/generate"
	"github.com/aquaproj/aqua/pkg/controller/initcmd"
	"github.com/aquaproj/aqua/pkg/controller/install"
	"github.com/aquaproj/aqua/pkg/controller/list"
	"github.com/aquaproj/aqua/pkg/controller/which"
	"github.com/aquaproj/aqua/pkg/download"
	"github.com/aquaproj/aqua/pkg/exec"
	"github.com/aquaproj/aqua/pkg/github"
	"github.com/aquaproj/aqua/pkg/github/archive"
	"github.com/aquaproj/aqua/pkg/github/content"
	"github.com/aquaproj/aqua/pkg/github/release"
	"github.com/aquaproj/aqua/pkg/install-registry"
	"github.com/aquaproj/aqua/pkg/installpackage"
	"github.com/aquaproj/aqua/pkg/link"
	"github.com/aquaproj/aqua/pkg/pkgtype"
	"github.com/aquaproj/aqua/pkg/pkgtype/githubarchive"
	"github.com/aquaproj/aqua/pkg/pkgtype/githubcontent"
	"github.com/aquaproj/aqua/pkg/pkgtype/githubrelease"
	"github.com/aquaproj/aqua/pkg/pkgtype/gobuild"
	"github.com/aquaproj/aqua/pkg/pkgtype/goinstall"
	http2 "github.com/aquaproj/aqua/pkg/pkgtype/http"
	"github.com/aquaproj/aqua/pkg/runtime"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/go-osenv/osenv"
	"net/http"
)

// Injectors from wire.go:

func InitializeListCommandController(ctx context.Context, param *config.Param, httpClient *http.Client) *list.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs)
	repositoryService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	registryDownloader := download.NewRegistryDownloader(repositoryService, httpDownloader)
	installer := registry.New(param, registryDownloader, fs)
	controller := list.NewController(configFinder, configReader, installer)
	return controller
}

func InitializeInitCommandController(ctx context.Context, param *config.Param) *initcmd.Controller {
	repositoryService := github.New(ctx)
	fs := afero.NewOsFs()
	controller := initcmd.New(repositoryService, fs)
	return controller
}

func InitializeGenerateCommandController(ctx context.Context, param *config.Param, httpClient *http.Client) *generate.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs)
	repositoryService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	registryDownloader := download.NewRegistryDownloader(repositoryService, httpDownloader)
	installer := registry.New(param, registryDownloader, fs)
	fuzzyFinder := generate.NewFuzzyFinder()
	controller := generate.New(configFinder, configReader, installer, repositoryService, fs, fuzzyFinder)
	return controller
}

func InitializeInstallCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *install.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs)
	repositoryService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	registryDownloader := download.NewRegistryDownloader(repositoryService, httpDownloader)
	installer := registry.New(param, registryDownloader, fs)
	linker := link.New()
	executor := exec.New()
	archiveClient := github.NewArchiveClient(repositoryService)
	downloader := archive.New(archiveClient, httpDownloader)
	githubarchiveInstaller := githubarchive.New(param, fs, downloader)
	client := content.New(repositoryService, httpDownloader)
	githubcontentInstaller := githubcontent.New(param, fs, client)
	releaseClient := release.New(repositoryService, httpDownloader)
	githubreleaseInstaller := githubrelease.New(param, fs, rt, releaseClient)
	goBuilder := gobuild.NewGoBuilder(executor)
	gobuildInstaller := gobuild.New(param, fs, goBuilder, downloader, rt)
	goInstaller := goinstall.NewGoInstaller(executor)
	goinstallInstaller := goinstall.New(param, goInstaller, fs)
	httpInstaller := http2.New(param, fs, rt)
	v := pkgtype.New(githubarchiveInstaller, githubcontentInstaller, githubreleaseInstaller, gobuildInstaller, goinstallInstaller, httpInstaller)
	installpackageInstaller := installpackage.New(param, rt, fs, linker, executor, v)
	controller := install.New(param, configFinder, configReader, installer, installpackageInstaller, fs)
	return controller
}

func InitializeWhichCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) which.Controller {
	fs := afero.NewOsFs()
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs)
	repositoryService := github.New(ctx)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	registryDownloader := download.NewRegistryDownloader(repositoryService, httpDownloader)
	installer := registry.New(param, registryDownloader, fs)
	osEnv := osenv.New()
	linker := link.New()
	archiveClient := github.NewArchiveClient(repositoryService)
	downloader := archive.New(archiveClient, httpDownloader)
	githubarchiveInstaller := githubarchive.New(param, fs, downloader)
	client := content.New(repositoryService, httpDownloader)
	githubcontentInstaller := githubcontent.New(param, fs, client)
	releaseClient := release.New(repositoryService, httpDownloader)
	githubreleaseInstaller := githubrelease.New(param, fs, rt, releaseClient)
	executor := exec.New()
	goBuilder := gobuild.NewGoBuilder(executor)
	gobuildInstaller := gobuild.New(param, fs, goBuilder, downloader, rt)
	goInstaller := goinstall.NewGoInstaller(executor)
	goinstallInstaller := goinstall.New(param, goInstaller, fs)
	httpInstaller := http2.New(param, fs, rt)
	v := pkgtype.New(githubarchiveInstaller, githubcontentInstaller, githubreleaseInstaller, gobuildInstaller, goinstallInstaller, httpInstaller)
	controller := which.New(param, configFinder, configReader, installer, rt, osEnv, fs, linker, v)
	return controller
}

func InitializeExecCommandController(ctx context.Context, param *config.Param, httpClient *http.Client, rt *runtime.Runtime) *exec2.Controller {
	fs := afero.NewOsFs()
	linker := link.New()
	executor := exec.New()
	repositoryService := github.New(ctx)
	archiveClient := github.NewArchiveClient(repositoryService)
	httpDownloader := download.NewHTTPDownloader(httpClient)
	downloader := archive.New(archiveClient, httpDownloader)
	installer := githubarchive.New(param, fs, downloader)
	client := content.New(repositoryService, httpDownloader)
	githubcontentInstaller := githubcontent.New(param, fs, client)
	releaseClient := release.New(repositoryService, httpDownloader)
	githubreleaseInstaller := githubrelease.New(param, fs, rt, releaseClient)
	goBuilder := gobuild.NewGoBuilder(executor)
	gobuildInstaller := gobuild.New(param, fs, goBuilder, downloader, rt)
	goInstaller := goinstall.NewGoInstaller(executor)
	goinstallInstaller := goinstall.New(param, goInstaller, fs)
	httpInstaller := http2.New(param, fs, rt)
	v := pkgtype.New(installer, githubcontentInstaller, githubreleaseInstaller, gobuildInstaller, goinstallInstaller, httpInstaller)
	installpackageInstaller := installpackage.New(param, rt, fs, linker, executor, v)
	configFinder := finder.NewConfigFinder(fs)
	configReader := reader.New(fs)
	registryDownloader := download.NewRegistryDownloader(repositoryService, httpDownloader)
	registryInstaller := registry.New(param, registryDownloader, fs)
	osEnv := osenv.New()
	controller := which.New(param, configFinder, configReader, registryInstaller, rt, osEnv, fs, linker, v)
	execController := exec2.New(installpackageInstaller, controller, executor, osEnv, fs)
	return execController
}
