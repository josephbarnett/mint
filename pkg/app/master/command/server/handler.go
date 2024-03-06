package server

import (
	"github.com/mintoolkit/mint/pkg/app"
	"github.com/mintoolkit/mint/pkg/app/master/command"
	"github.com/mintoolkit/mint/pkg/app/master/version"
	cmd "github.com/mintoolkit/mint/pkg/command"
	"github.com/mintoolkit/mint/pkg/docker/dockerclient"
	"github.com/mintoolkit/mint/pkg/report"
	"github.com/mintoolkit/mint/pkg/util/errutil"
	"github.com/mintoolkit/mint/pkg/util/fsutil"
	v "github.com/mintoolkit/mint/pkg/version"

	log "github.com/sirupsen/logrus"
)

const appName = command.AppName

type ovars = app.OutVars

// OnCommand implements the 'server' command
func OnCommand(
	xc *app.ExecutionContext,
	gparams *command.GenericParams) {
	logger := log.WithFields(log.Fields{"app": appName, "cmd": Name})

	viChan := version.CheckAsync(gparams.CheckVersion, gparams.InContainer, gparams.IsDSImage)

	cmdReport := report.NewServerCommand(gparams.ReportLocation, gparams.InContainer)
	cmdReport.State = cmd.StateStarted

	xc.Out.State(cmd.StateStarted)

	client, err := dockerclient.New(gparams.ClientConfig)
	if err == dockerclient.ErrNoDockerInfo {
		exitMsg := "missing Docker connection info"
		if gparams.InContainer && gparams.IsDSImage {
			exitMsg = "make sure to pass the Docker connect parameters to the slim app container"
		}

		xc.Out.Info("docker.connect.error",
			ovars{
				"message": exitMsg,
			})

		exitCode := command.ECTCommon | command.ECCNoDockerConnectInfo
		xc.Out.State("exited",
			ovars{
				"exit.code": exitCode,
				"version":   v.Current(),
				"location":  fsutil.ExeDir(),
			})
		xc.Exit(exitCode)
	}
	errutil.FailOn(err)

	if gparams.Debug {
		version.Print(xc, Name, logger, client, false, gparams.InContainer, gparams.IsDSImage)
	}

	xc.Out.State(cmd.StateCompleted)
	cmdReport.State = cmd.StateCompleted
	xc.Out.State(cmd.StateDone)

	vinfo := <-viChan
	version.PrintCheckVersion(xc, "", vinfo)

	cmdReport.State = cmd.StateDone
	if cmdReport.Save() {
		xc.Out.Info("report",
			ovars{
				"file": cmdReport.ReportLocation(),
			})
	}
}
