package main

import "fmt"

type Color struct {
	Disable bool
}

func (c *Color) surround(text, color string) string {
	if c.Disable {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, "\033[0m")
}

func (c *Color) Red(text string) string {
	return c.surround(text, "\033[0;31m")
}

func (c *Color) Yellow(text string) string {
	return c.surround(text, "\033[0;33m")
}

func (c *Color) Blue(text string) string {
	return c.surround(text, "\033[0;34m")
}

func (c *Color) BBlue(text string) string {
	return c.surround(text, "\033[1;34m")
}

func (c *Color) Cyan(text string) string {
	return c.surround(text, "\033[0;36m")
}

func (c *Color) Gray(text string) string {
	return c.surround(text, "\033[0;37m")
}
