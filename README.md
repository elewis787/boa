# Boa

Boa implements [Cobra's](https://github.com/spf13/cobra) help and usage functions to provided an interactive user experience. User's no longer need to spend time running multiple help commands to see how nested sub commands work!

## Install 

Use go get to install the latest version of the library.

`go get -u github.com/elewis787/boa@latest`

Next, include Boa in your application:

`import "github.com/elewis787/boa"`

## Usage

Using Boa is very simple. Below is an example on how to set the help/usage functions on a root command defined using Cobra. 

```go
	rootCmd := &cobra.Command{
		Version: "v0.0.1",
		Use:     "Example",
		Long:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		Short:   "example command",
		Example: "example [sub command]",
		RunE: func(cmd *cobra.Command, args []string) error {
		    return nil
		},
	}

	rootCmd.SetUsageFunc(boa.UsageFunc)
	rootCmd.SetHelpFunc(boa.HelpFunc)

```

The key lines are: 

```go
	rootCmd.SetUsageFunc(boa.UsageFunc)
	rootCmd.SetHelpFunc(boa.HelpFunc)
```

## Demo 

![demo](demo.gif)

## Used by 
- [rkl](https://github.com/elewis787/rkl)

## Future work 
- Eval how styles are exported. Goal is to make it easy to customize the layout without needing to build a cmd parser form cobra 
- Add back button 
- Option to execute sub command 
- Other ideas ? - Open a feature request or submit a PR ! 