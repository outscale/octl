/*
SPDX-FileCopyrightText: 2026 Outscale SAS <opensource@outscale.com>

SPDX-License-Identifier: BSD-3-Clause
*/
package alias

import (
	"bytes"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/outscale/octl/pkg/config"
	"github.com/outscale/octl/pkg/debug"
	"github.com/outscale/octl/pkg/flags"
	"github.com/outscale/octl/pkg/messages"
	"github.com/outscale/octl/pkg/output"
	"github.com/outscale/octl/pkg/runner"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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
		withoutUserArgs := userArgs(cmd, nil, true)
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
func userArgs(cmd *cobra.Command, fs config.FlagSet, skipUserFlags bool) []string {
	// compute # of iterated fields to build in target
	// expectation: the target flags iterate on the same prefix
	n := 1
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		nf, found := fs.Get(f.Name)
		if !found || !f.Changed {
			return
		}
		// flag should forced as a container and target must be iterated
		if nf.ContainerKind != reflect.Slice || !strings.Contains(nf.AliasTo, ".#.") {
			return
		}
		if svalue, ok := f.Value.(pflag.SliceValue); ok {
			n = max(n, len(svalue.GetSlice()))
		}
	})
	debug.Println("Generating", n, "iterated fields")
	// set flags
	var userArgs []string
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && !flags.HasDefault(f) {
			return
		}
		targetFlag := f.Name
		nf, found := fs.Get(targetFlag)
		switch {
		case targetFlag == "verbose" || targetFlag == "config" || targetFlag == "profile":
		case flags.IsNoForward(f):
			return
		case !found && skipUserFlags:
			return
		case found:
			targetFlag = nf.AliasTo
		}
		if targetFlag == "" {
			return
		}
		newFlagName := func(oldFlag string, i int) string {
			newFlag := strings.Replace(oldFlag, ".#.", "."+strconv.Itoa(i)+".", 1)
			return strings.ReplaceAll(newFlag, ".#.", ".0.")
		}
		// This is a natural flag, types will match, repeat the value to all sets
		if nf.ContainerKind != reflect.Slice {
			done := map[string]bool{}
			for i := range n {
				newFlag := newFlagName(targetFlag, i)
				if done[newFlag] {
					continue
				}
				debug.Println(i, newFlag, f.Value.String())
				if svalue, ok := f.Value.(pflag.SliceValue); ok {
					userArgs = append(userArgs, "--"+newFlag+"="+strings.Join(svalue.GetSlice(), ","))
				} else {
					userArgs = append(userArgs, "--"+newFlag+"="+f.Value.String())
				}
				done[newFlag] = true
			}
			return
		}
		// The field was remapped from single value to list
		// Expectation: target field is iterated - should be enforced by config, noop here to avoid a potential crash.
		if !strings.Contains(targetFlag, ".#.") {
			debug.Println("WARNING - Slice field is not mapped to iterated flag -", f.Name, "ignored")
			return
		}
		svalue, ok := f.Value.(pflag.SliceValue)
		if !ok {
			debug.Println("WARNING - Slice field is a slice -", f.Name, "ignored")
			return
		}
		vals := svalue.GetSlice()
		if len(vals) == 0 {
			return
		}
		debug.Println(targetFlag, vals)
		for i := range n {
			// if one value is missing, the last value will be repeated
			val := vals[min(i, len(vals)-1)]
			newFlag := newFlagName(targetFlag, i)
			debug.Println(i, newFlag, val)
			userArgs = append(userArgs, "--"+newFlag+"="+val)
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
				},
			) {
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
