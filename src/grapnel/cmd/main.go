package main
/*
Copyright (c) 2014 Eric Anderton <eric.t.anderton@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

import (
  . "grapnel"
  . "grapnel/flag"
  "os"
  "fmt"
  log "grapnel/log"
)

// application configurables w/default settings
var (
  defaultConfigFileName string = "/etc/grapnel-config.toml"
  configFileName string

  defaultPackageFileName string = "./grapnel.toml"
  packageFileName string

  defaultLockFileName string = "./grapnel-lock.toml"
  lockFileName string

  defaultTargetPath string = "./src"
  targetPath string

  flagQuiet bool
  flagVerbose bool
  flagDebug bool
)

func getResolver() *Resolver{
  resolver := NewResolver()
  resolver.LibSources["git"] = &GitSCM{}

  //TODO: get rules from config file
  resolver.AddRewriteRules(BasicRewriteRules)
  resolver.AddRewriteRules(GitRewriteRules)
  return resolver
}

func configureLogging() {
  if flagDebug {
    log.SetGlobalLogLevel(log.DEBUG)
  } else if flagQuiet {
    log.SetGlobalLogLevel(log.ERROR)
  } else if flagVerbose {
    log.SetGlobalLogLevel(log.INFO)
  } else {
    log.SetGlobalLogLevel(log.WARN)
  }
}

func ShowVersion() error {
  fmt.Printf("%s v%s\n", PROGRAM_NAME, VERSION)
  return nil
}

var rootCmd = &Command{
  Alias: PROGRAM_NAME,
  Desc: "Manages dependencies for Go projects",
  Help: "Use 'grapnel help [command]' for more information about that command.",
  Flags: FlagMap {
    "quiet": &Flag {
      Alias: "q",
      Desc: "Quiet output",
      Fn: BoolFlagFn(&flagQuiet),
    },
    "verbose": &Flag {
      Alias: "v",
      Desc: "Verbose output",
      Fn: BoolFlagFn(&flagVerbose),
    },
    "debug": &Flag {
      Desc: "Debug output",
      Fn: BoolFlagFn(&flagDebug),
    },
    "config": &Flag {
      Alias: "c",
      Desc: "Configuration file",
      ArgDesc: "[filename]",
      Fn: StringFlagFn(&configFileName),
    },
    "version": &Flag {
      Desc: "Displays version information",
      Fn: SimpleFlagFn(ShowVersion),
    },
  },
  Commands: CommandMap {
    "install": &installCmd,
    "update": &updateCmd,
    "version": &Command{
      Desc: "Version information",
      Fn: SimpleCommandFn(ShowVersion),
    },
  },
}

func main() {
  log.SetFlags(0)
  rootCmd.Help =
    fmt.Sprintf("Defaults:\n") +
    fmt.Sprintf("  Config file = %s\n", defaultConfigFileName) +
    fmt.Sprintf("  Lock file = %s\n", defaultLockFileName) +
    fmt.Sprintf("  Package file = %s\n", defaultPackageFileName) +
    "\n" + rootCmd.Help
  if err := rootCmd.Execute(os.Args...); err != nil {
    log.Error(err)
    rootCmd.ShowHelp(os.Args[0])
  }
}
