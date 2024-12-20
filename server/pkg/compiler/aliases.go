package compiler

import (
	"errors"
	"fmt"
	"github.com/Argentusz/MTP_coursework/pkg/consts"
	"strings"
)

func (cmp *Compiler) setDefaultAliases() {
	cmp.Aliases = map[string]string{
		"$define":   "",
		"$main":     "0",
		"$signone":  fmt.Sprint(consts.SIGNONE),
		"$sigfpe":   fmt.Sprint(consts.SIGFPE),
		"$sigtrace": fmt.Sprint(consts.SIGTRACE),
		"$sigsegv":  fmt.Sprint(consts.SIGSEGV),
		"$sigterm":  fmt.Sprint(consts.SIGTERM),
		"$sigint":   fmt.Sprint(consts.SIGINT),
		"$sigiie":   fmt.Sprint(consts.SIGIIE),
		"$sigill":   fmt.Sprint(consts.SIGILL),
		"$m8":       fmt.Sprint(consts.MAX_WORD8),
		"$m16":      fmt.Sprint(consts.MAX_WORD16),
		"$m32":      fmt.Sprint(consts.MAX_WORD32),
	}
}

func (cmp *Compiler) deAlias(line string) (string, error) {
	line = prepLine(line)
	if strings.HasPrefix(line, "$define") {
		return "", cmp.define(line)
	}

	cmd := strings.Split(line, " ")

	dealiased := ""
	for _, v := range cmd {
		switch strings.HasPrefix(v, "$") {
		case false:
			dealiased = fmt.Sprintf("%s %s", dealiased, v)
		case true:
			app, found := cmp.Aliases[v]
			if !found {
				return "", errors.New(fmt.Sprintf("alias \"%s\" not found", v))
			}

			dealiased = fmt.Sprintf("%s %s", dealiased, app)
		}
	}

	return strings.Trim(dealiased, " "), nil
}

func (cmp *Compiler) define(line string) error {
	cmd := strings.SplitN(line, " ", 3)

	if len(cmd) != 3 {
		return errors.New("$define expected 2 arguments")
	}

	if !strings.HasPrefix(cmd[1], "$") {
		return errors.New(fmt.Sprintf("alias \"%s\" should have $ prefix", cmd[1]))
	}

	if cmd[1] == "$" {
		return errors.New("can not have \"$\" alias")
	}

	var err error
	cmd[2], err = cmp.deAlias(cmd[2])
	if err != nil {
		return err
	}

	if cmd[2] == "" {
		return errors.New("can not alias empty string nor define a define")
	}

	_, found := cmp.Aliases[cmd[1]]
	if found {
		return errors.New(fmt.Sprintf("can not overwrite alias \"%s\"", cmd[1]))
	}

	cmp.Aliases[cmd[1]] = cmd[2]
	return nil
}
