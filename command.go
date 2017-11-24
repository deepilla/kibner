package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	flag "github.com/ogier/pflag"
)

var (
	ErrBadArgs = errors.New("bad args")
)

type Env struct {
	Stdin  io.Writer
	Stdout io.Writer
	Stderr io.Writer
	LogTo  io.Writer
}

type CommandSet struct {
	Name     string
	commands []*Command
}

func NewCommandSet(name string, commands ...*Command) *CommandSet {
	return &CommandSet{
		Name:     name,
		commands: commands,
	}
}

func (cs *CommandSet) Run(name string, args []string) error {
	return cs.RunWithEnv(name, args, &Env{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

func (cs *CommandSet) RunWithEnv(name string, args []string, env *Env) error {

	c := cs.findByName(name)
	if c == nil {
		cs.Usage()
		return errors.New("no such command " + name)
	}

	var help bool

	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = c.usage
	fs.BoolVarP(&help, "help", "h", false, "Show this help page")
	c.opts.addToFlagSet(fs)

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if help {
		c.usage()
		return nil
	}

	err = c.run(c.opts, fs.Args(), env)
	if err != nil {
		if err == ErrBadArgs {
			c.usage()
		}
		return err
	}

	return nil
}

func (cs *CommandSet) Usage() {

	width := 0
	marginLeft := 2
	marginRight := 7

	for _, c := range cs.commands {
		if w := len(c.name); w > width {
			width = w
		}
	}

	fmt.Println("Usage:", cs.Name, "[command] [...]")
	fmt.Println()
	fmt.Println("Available Commands:")
	for _, c := range cs.commands {
		fmt.Printf("%*s%-*s %s\n", marginLeft, "", width+marginRight, c.name, c.desc)
	}
	fmt.Println()
	fmt.Printf("Use \"%s <command> --help\" for help with individual commands.\n", cs.Name)
}

func (cs *CommandSet) findByName(name string) *Command {
	for _, c := range cs.commands {
		if c.name == name || c.alias == name {
			return c
		}
	}
	return nil
}

type Command struct {
	name   string
	run    func(Options, []string, *Env) error
	desc   string
	alias  string
	syntax string
	opts   Options
}

func NewCommand(name string, fn func(Options, []string, *Env) error, modifiers ...func(*Command)) *Command {

	if fn == nil {
		fn = func(Options, []string, *Env) error {
			return errors.New(name + " is not implemented")
		}
	}

	c := &Command{
		name:   name,
		run:    fn,
		desc:   "No description provided for " + name,
		syntax: "unknown",
	}

	for _, fn := range modifiers {
		fn(c)
	}

	return c
}

func WithDescription(val string) func(*Command) {
	return func(c *Command) {
		c.desc = val
	}
}

func WithAlias(val string) func(*Command) {
	return func(c *Command) {
		c.alias = val
	}
}

func WithSyntax(val string) func(*Command) {
	return func(c *Command) {
		c.syntax = val
	}
}

func WithOption(name string, usage string, value interface{}) func(*Command) {
	return WithOptionAlias(name, "", usage, value)
}

func WithOptionAlias(name string, alias string, usage string, value interface{}) func(*Command) {
	return func(c *Command) {
		c.opts = append(c.opts, newOption(name, alias, value, usage))
	}
}

func (c *Command) usage() {

	fmt.Println("Usage:", c.syntax)
	fmt.Println(c.desc + ".")

	if len(c.opts) > 0 {
		fmt.Println("\nOptions:")
		c.opts.print()
	}
}

type Options []*Option

func (opts Options) Get(name string) *Option {

	for _, o := range opts {
		if o.Name == name {
			return o
		}
	}

	return nil
}

func (opts Options) addToFlagSet(fs *flag.FlagSet) {

	for _, o := range opts {
		o.addToFlagSet(fs)
	}
}

func (opts Options) print() {

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	for _, o := range opts {
		o.addToFlagSet(fs)
	}
	fs.PrintDefaults()
}

type Option struct {
	Name  string
	alias string
	usage string
	value interface{}
}

func newOption(name string, alias string, value interface{}, usage string) *Option {

	if name == "" {
		panic("Invalid name for flag: name cannot be blank")
	}

	if value == nil {
		panic(fmt.Sprintf("Invalid default for flag %q: value cannot be nil", name))
	}

	// Create a pointer to a new value of the desired type.
	p := reflect.New(reflect.TypeOf(value))
	// Set the new value to the provided value.
	p.Elem().Set(reflect.ValueOf(value))

	return &Option{
		Name:  name,
		alias: alias,
		usage: usage,
		value: p.Interface(),
	}
}

func (o *Option) Bool() bool {
	return *o.value.(*bool)
}

func (o *Option) Float() float64 {
	return *o.value.(*float64)
}

func (o *Option) Int() int {
	return *o.value.(*int)
}

func (o *Option) String() string {
	return *o.value.(*string)
}

func (o *Option) Uint() uint {
	return *o.value.(*uint)
}

// The flag.Getter interface is not defined in the
// pflag package so we have to use our own.
type getter interface {
	Get() interface{}
}

func (o *Option) Value() interface{} {

	if getter, ok := o.value.(getter); ok {
		return getter.Get()
	}

	// Remember to dereference the ptr.
	return reflect.ValueOf(o.value).Elem().Interface()
}

func (o *Option) addToFlagSet(fs *flag.FlagSet) {

	switch v := o.value.(type) {
	case *bool:
		fs.BoolVarP(v, o.Name, o.alias, *v, o.usage)
	case *float64:
		fs.Float64VarP(v, o.Name, o.alias, *v, o.usage)
	case *int:
		fs.IntVarP(v, o.Name, o.alias, *v, o.usage)
	case *string:
		fs.StringVarP(v, o.Name, o.alias, *v, o.usage)
	case *uint:
		fs.UintVarP(v, o.Name, o.alias, *v, o.usage)
	case flag.Value:
		fs.VarP(v, o.Name, o.alias, o.usage)
	default:
		panic(fmt.Sprintf("Invalid default for flag %q: unsupported type %T", o.Name, o.value))
	}
}
