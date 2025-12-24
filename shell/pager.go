package shell

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// Pager умеет показать текст через внешний пейджер.
type Pager struct {
	cmd  string
	args []string
}

// NewPager конструирует Pager, опираясь на $PAGER или на команду "less".
func NewPager() *Pager {
	pager := os.Getenv("PAGER")
	if pager == "" {
		pager = "less"
	}
	// -R: сохранить ANSI-цвета, -F: сразу выйти, если текст умещается
	return &Pager{
		cmd:  pager,
		args: []string{"-R", "-F"},
	}
}

// Page запускает пейджер и впихивает в его stdin текст.
func (p *Pager) Page(text string) error {
	c := exec.Command(p.cmd, p.args...)
	stdin, err := c.StdinPipe()
	if err != nil {
		return err
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Start(); err != nil {
		stdin.Close()
		return err
	}
	// копируем весь текст
	io.Copy(stdin, bytes.NewBufferString(text))
	stdin.Close()
	return c.Wait()
}
