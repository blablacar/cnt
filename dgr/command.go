package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/blablacar/dgr/dgr/common"
	"github.com/n0rad/go-erlog/logs"
	"github.com/spf13/cobra"
)

type DgrCommand interface {
	Build() error
	Update() error
	CleanAndBuild() error
	CleanAndTry() error
	Clean()
	Push() error
	Install() ([]string, error)
	Test() error
	Graph() error
	Sign() error
	Init() error
	Cp([]string) error
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean build",
	Long:  `clean build, including rootfs`,
	Run: func(cmd *cobra.Command, args []string) {
		checkNoArgs(args)

		NewAciOrPod(workPath, Args).Clean()
	},
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "generate dependency graph",
	Long:  `generate dependency graph`,
	Run: func(cmd *cobra.Command, args []string) {
		checkNoArgs(args)

		if err := NewAciOrPod(workPath, Args).Graph(); err != nil {
			logs.WithE(err).Fatal("Install command failed")
		}
	},
}

var aciVersion = &cobra.Command{
	Use:   "aci-version file",
	Short: "display version of aci",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}
		im, err := common.ExtractManifestFromAci(args[0])
		if err != nil {
			logs.WithE(err).Fatal("Failed to get manifest from file")
		}
		val, _ := im.Labels.Get("version")
		fmt.Println(val)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "display config elements",
	Run: func(cmd *cobra.Command, args []string) {
		templateStr := "{{.}}"
		if len(args) > 1 {
			for i := 0; i < len(args); i++ {
				if len(args[i]) > 0 {
					args[i] = strings.ToUpper(args[i][0:1]) + args[i][1:]
				}
			}
			res := strings.Join(args, ".")
			templateStr = "{{." + res + "}}"
		} else if len(args) == 1 {
			args := strings.Split(args[0], ".")
			for i := 0; i < len(args); i++ {
				if len(args[i]) > 0 {
					args[i] = strings.ToUpper(args[i][0:1]) + args[i][1:]
				}
			}
			res := strings.Join(args, ".")
			templateStr = "{{." + res + "}}"
		}
		tmpl, err := template.New("").Parse(templateStr)
		if err != nil {
			logs.WithE(err).Fatal("Failed to parse config template")
		}
		b := bytes.Buffer{}
		w := bufio.NewWriter(&b)
		if err := tmpl.Execute(w, Home.Config); err != nil {
			logs.WithE(err).Fatal("Failed to process config templating")
		}
		w.Flush()
		fmt.Println(b.String())
	},
}

var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "cp file",
	Long:  `cp file or directory from aci to host`,
	Run: func(cmd *cobra.Command, args []string) {
		command := NewAciOrPod(workPath, Args)
		if err := command.Cp(args); err != nil {
			logs.WithE(err).Fatal("Cp command failed")
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of dgr",
	Long:  `Print the version number of dgr`,
	Run: func(cmd *cobra.Command, args []string) {
		displayVersionAndExit()
	},
}

var buildCmd = newBuildCommand(false)
var installCmd = newInstallCommand(false)
var pushCmd = newPushCommand(false)
var testCmd = newTestCommand(false)
var tryCmd = newTryCommand(false)
var signCmd = newSignCommand(false)
var updateCmd = newUpdateCommand(false)

///////////////////////////////////////////////////////////////

func checkNoArgs(args []string) {
	if len(args) > 0 {
		logs.WithField("args", args).Fatal("Unknown arguments")
	}
}

func newTryCommand(userClean bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "try",
		Short: "try templater (experimental)",
		Long:  `try templater (experimental)`,
		Run: func(cmd *cobra.Command, args []string) {
			checkNoArgs(args)

			if err := NewAciOrPod(workPath, Args).CleanAndTry(); err != nil {
				logs.WithE(err).Fatal("Try command failed")
			}
		},
	}
	return cmd
}

func newBuildCommand(userClean bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build aci or pod",
		Long:  `build an aci or a pod`,
		Run: func(cmd *cobra.Command, args []string) {
			checkNoArgs(args)

			if err := NewAciOrPod(workPath, Args).CleanAndBuild(); err != nil {
				logs.WithE(err).Fatal("Build command failed")
			}
		},
	}
	cmd.Flags().BoolVarP(&Args.KeepBuilder, "keep-builder", "k", false, "Keep builder container after exit")
	cmd.Flags().BoolVarP(&Args.CatchOnError, "catch-on-error", "c", false, "Catch a shell on build* runlevel fail") // TODO This is builder dependent and should be pushed by builder ? or find a way to become generic
	cmd.Flags().BoolVarP(&Args.CatchOnStep, "catch-on-step", "C", false, "Catch a shell after each build* runlevel")
	return cmd
}

func newUpdateCommand(userClean bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update attributes/templates/prestart in aci",
		Long:  `update attributes/templates/prestart in aci`,
		Run: func(cmd *cobra.Command, args []string) {
			checkNoArgs(args)

			if err := NewAciOrPod(workPath, Args).Update(); err != nil {
				logs.WithE(err).Fatal("Update command failed")
			}
		},
	}
	cmd.Flags().BoolVarP(&Args.KeepBuilder, "keep-builder", "k", false, "Keep builder container after exit")
	cmd.Flags().BoolVarP(&Args.CatchOnError, "catch-on-error", "c", false, "Catch a shell on build* runlevel fail") // TODO This is builder dependent and should be pushed by builder ? or find a way to become generic
	cmd.Flags().BoolVarP(&Args.CatchOnStep, "catch-on-step", "C", false, "Catch a shell after each build* runlevel")
	return cmd
}

func newSignCommand(underClean bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "sign image",
		Long:  `sign image`,
		Run: func(cmd *cobra.Command, args []string) {
			checkNoArgs(args)

			command := NewAciOrPod(workPath, Args)
			if underClean {
				command.Clean()
			}
			if err := command.Sign(); err != nil {
				logs.WithE(err).Fatal("Sign command failed")
			}
		},
	}

	return cmd
}

func newInstallCommand(underClean bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install image(s)",
		Long:  `install image(s) to rkt local storage`,
		Run: func(cmd *cobra.Command, args []string) {
			checkNoArgs(args)

			command := NewAciOrPod(workPath, Args)
			if underClean {
				command.Clean()
			}
			if _, err := command.Install(); err != nil {
				logs.WithE(err).Fatal("Install command failed")
			}
		},
	}

	cmd.Flags().BoolVarP(&Args.NoTestFail, "no-test-fail", "T", false, "Fail if no tests found")
	cmd.Flags().BoolVarP(&Args.Test, "test", "t", false, "Run tests before install")
	return cmd
}

func newPushCommand(underClean bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "push image(s)",
		Long:  `push images to repository`,
		Run: func(cmd *cobra.Command, args []string) {
			checkNoArgs(args)

			command := NewAciOrPod(workPath, Args)
			if underClean {
				command.Clean()
			}
			if err := command.Push(); err != nil {
				logs.WithE(err).Fatal("Push command failed")
			}
		},
	}
	cmd.Flags().BoolVarP(&Args.NoTestFail, "no-test-fail", "T", false, "Fail if no tests found")
	cmd.Flags().BoolVarP(&Args.Test, "test", "t", false, "Run tests before push")
	return cmd
}

func newTestCommand(underClean bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "test image(s)",
		Long:  `test image(s)`,
		Run: func(cmd *cobra.Command, args []string) {
			checkNoArgs(args)

			command := NewAciOrPod(workPath, Args)
			if underClean {
				command.Clean()
			}
			if err := command.Test(); err != nil {
				logs.WithE(err).Fatal("Test command failed")
			}
		},
	}
	cmd.Flags().BoolVarP(&Args.NoTestFail, "no-test-fail", "T", false, "Fail if no tests found")
	cmd.Flags().BoolVarP(&Args.KeepBuilder, "keep-builder", "k", false, "Keep aci & test builder container after exit")
	return cmd
}

func init() {
	cleanCmd.AddCommand(newInstallCommand(true))
	cleanCmd.AddCommand(newPushCommand(true))
	cleanCmd.AddCommand(newTestCommand(true))
	cleanCmd.AddCommand(newBuildCommand(true))
	cleanCmd.AddCommand(newTryCommand(true))
	cleanCmd.AddCommand(newSignCommand(true))
	cleanCmd.AddCommand(newUpdateCommand(true))
}
