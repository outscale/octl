/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package alias

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/outscale/octl/pkg/config"
	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/output"
	"github.com/outscale/octl/pkg/runner"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	defaultValueAnnotation = "default.octl.outscale.com"
)

func SetDefault(f *pflag.Flag, value string) {
	f.DefValue = value
	if f.Annotations == nil {
		f.Annotations = map[string][]string{}
	}
	f.Annotations[defaultValueAnnotation] = []string{value}
}

func Reset(f *pflag.Flag) {
	if f == nil {
		return
	}
	f.Changed = false
	if !HasDefault(f) {
		return
	}
	f.DefValue = ""
	delete(f.Annotations, defaultValueAnnotation)
}

func HasDefault(f *pflag.Flag) bool {
	return f != nil && f.Annotations != nil && len(f.Annotations[defaultValueAnnotation]) > 0
}

func Default(f *pflag.Flag) string {
	if f.Annotations == nil {
		return ""
	}
	values, found := f.Annotations[defaultValueAnnotation]
	if !found || len(values) == 0 {
		return ""
	}
	return values[0]
}

var prompts = map[config.Action]string{
	config.ActionDelete: "Are you sure you want to delete these resource(s) ?",
}

// RunFunc returns the func to run an alias.
func RunFunc(rootPath string, a config.Alias) func(cmd *cobra.Command, args []string) {
	switch {
	case slices.Contains(a.Command, "|"):
		return runFuncWithPipe(rootPath, a)
	case a.Prompt != nil:
		return runFuncWithPrompt(rootPath, a)
	default:
		return runSingle(rootPath, a.Command, a.Flags, false)
	}
}

func runFuncWithPipe(rootPath string, a config.Alias) func(cmd *cobra.Command, args []string) {
	first, last, _ := lo.Cut(a.Command, []string{"|"})
	return func(cmd *cobra.Command, args []string) {
		// build args based on flags
		withUserArgs := userArgs(cmd, a.Flags, false)
		withoutUserArgs := userArgs(cmd, nil, true)
		// only one of the commands will inherit flags.
		firstHasUserArgs := slices.Contains(first, a.AliasTo)
		buf := &bytes.Buffer{}
		output.InjectOutput(buf)
		if firstHasUserArgs {
			_ = execAlias(rootPath, first, withUserArgs, cmd, args)
		} else {
			_ = execAlias(rootPath, first, withoutUserArgs, cmd, args)
		}
		messages.Info("Piping output")
		runner.InjectStdin(buf.Bytes())
		output.InjectOutput(os.Stdout)
		if firstHasUserArgs {
			_ = execAlias(rootPath, last, withoutUserArgs, cmd, args)
		} else {
			_ = execAlias(rootPath, last, withUserArgs, cmd, args)
		}
	}
}

// runFuncWithPrompt returns the fun to run a non piped alias, potentially having a prompt.
func runFuncWithPrompt(rootPath string, a config.Alias) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		withUserArgs := userArgs(cmd, a.Flags, false)
		withoutUserArgs := userArgs(cmd, a.Prompt.Flags, true)
		if yes, _ := cmd.Flags().GetBool("yes"); !yes {
			if len(a.Prompt.DisplayCommand) > 0 {
				_ = execAlias(rootPath, a.Prompt.DisplayCommand, withoutUserArgs, cmd, args)
			}
			if !messages.Prompt(prompts[a.Prompt.Action]) {
				return
			}
		}
		exec := func(cmd *cobra.Command, args []string) int {
			return execAlias(rootPath, a.Command, withUserArgs, cmd, args)
		}
		if lo.CountBy(a.Command, func(arg string) bool { return strings.HasPrefix(arg, "%") }) == 1 {
			iterate(exec, cmd, args)
		} else {
			once(exec, cmd, args)
		}
	}
}

func runSingle(rootPath string, command []string, flags config.FlagSet, skipUserFlags bool) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		userArgs := userArgs(cmd, flags, skipUserFlags)
		exec := func(cmd *cobra.Command, args []string) int {
			return execAlias(rootPath, command, userArgs, cmd, args)
		}
		if lo.CountBy(command, func(arg string) bool { return strings.HasPrefix(arg, "%") }) == 1 {
			iterate(exec, cmd, args)
		} else {
			once(exec, cmd, args)
		}
	}
}

func resetFlags(fs *pflag.FlagSet) {
	fs.VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			if svalue, ok := f.Value.(pflag.SliceValue); ok {
				_ = svalue.Replace(nil)
			} else {
				_ = f.Value.Set(f.DefValue)
			}
			f.Changed = false
		}
	})
}

// once executes a run func once
func once(fn func(cmd *cobra.Command, args []string) int, cmd *cobra.Command, args []string) {
	_ = fn(cmd, args)
}

// iterate executes a run func multiple times, until all args have been consumed.
func iterate(fn func(cmd *cobra.Command, args []string) int, cmd *cobra.Command, args []string) {
	for {
		consumed := fn(cmd, args)
		debug.Println("consumed", consumed, "len", len(args))
		if consumed <= 0 || len(args) == consumed {
			break
		}
		args = args[consumed:]
	}
}

// userArgs returns the list of args for the underlying command, including flags mapped from the alias flags.
func userArgs(cmd *cobra.Command, flags config.FlagSet, skipUserFlags bool) []string {
	var userArgs []string
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		newFlag := f.Name
		nf, found := flags.Get(newFlag)
		switch {
		case newFlag == "verbose" || newFlag == "config" || newFlag == "profile":
		case !found && skipUserFlags:
			return
		case found:
			newFlag = nf.AliasTo
		}
		switch {
		case f.Changed:
			if svalue, ok := f.Value.(pflag.SliceValue); ok {
				userArgs = append(userArgs, "--"+newFlag+"="+strings.Join(svalue.GetSlice(), ","))
				return
			}
			userArgs = append(userArgs, "--"+newFlag+"="+f.Value.String())
		case HasDefault(f):
			userArgs = append(userArgs, "--"+newFlag+"="+Default(f))
		}
	})
	return userArgs
}

// execAlias executes the run func, based on inherited flags.
func execAlias(rootPath string, command []string, userArgs []string, cmd *cobra.Command, args []string) int {
	nargs := make([]string, 2, len(command)+2)
	nargs[0] = "octl"
	nargs[1] = rootPath
	skipnextvalue := false
	consumed := -1
	for _, arg := range command {
		if !strings.HasPrefix(arg, "%") {
			// skip flags already present in user flags
			isFlag := strings.HasPrefix(arg, "--")
			if isFlag && slices.ContainsFunc(
				userArgs, func(uf string) bool {
					return strings.HasPrefix(uf, arg+"=")
				}) {
				skipnextvalue = true
				continue
			}
			// skip value present after skipped flag
			if skipnextvalue && !isFlag {
				skipnextvalue = false
				continue
			}
			skipnextvalue = false
			nargs = append(nargs, arg)
			continue
		}
		if arg == "%*" {
			if len(args) == 0 {
				_ = cmd.Usage()
				os.Exit(1)
			}
			nargs = append(nargs, strings.Join(args, ","))
			consumed = len(args) - 1
			continue
		}
		idx, err := strconv.Atoi(arg[1:])
		if err != nil {
			messages.Warn(err.Error())
			continue
		}
		if idx >= len(args) {
			_ = cmd.Usage()
			os.Exit(1)
		}
		nargs = append(nargs, args[idx])
		consumed = max(consumed, idx)
	}
	nargs = append(nargs, userArgs...)
	nargs, err := runner.TemplateArgs(nargs)
	if err != nil {
		messages.ExitErr(err)
	}
	messages.Info("Resolving alias to %v", nargs)
	// no need to check for an update a second time
	nargs = append(nargs, "--no-upgrade")

	resetFlags(cmd.Flags())
	os.Args = nargs
	err = cmd.Execute()
	if err != nil {
		messages.ExitErr(err)
	}
	return consumed + 1
}
