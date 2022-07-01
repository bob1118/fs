package fsconf

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Fsconf struct {
	confDir      string   //conf direcotry.
	defaultExt   string   //old conf files named file.ext `./vars.xml.default`
	changedFiles []string //new conf files who changed by function buildBootableConf()
}

func New(conf string, ext string) *Fsconf {
	var dir string
	var changedfiles = make([]string, 0, 10)

	if _, err := os.Stat(conf); err != nil { //defaultConfigDirectory
		os := runtime.GOOS
		switch os {
		case `linux`:
			dir = `/etc/freeswitch`
		case `windows`:
			dir = `C:/Program Files/FreeSWITCH/conf`
		default:
		}
	} else {
		dir = conf
	}

	//changedConfigFiles
	errWalk := filepath.WalkDir(dir,
		func(path string, d fs.DirEntry, e error) error {
			if e == nil {
				if filepath.Ext(path) == ext {
					changedfile := strings.TrimSuffix(path, ext)
					if _, e := os.Stat(changedfile); e == nil {
						changedfiles = append(changedfiles, changedfile)
					}
				}
			} else {
				log.Println(e)
			}
			return e
		})
	if errWalk != nil {
		return nil
	} else {
		return &Fsconf{confDir: dir, defaultExt: ext, changedFiles: changedfiles}
	}
}

func Newconf(conf string) *Fsconf {
	return New(conf, `.default`)
}

func (p *Fsconf) Dir() string { return p.confDir }

func (p *Fsconf) Update(file string, old []byte, new []byte) error {
	var err error
	if data, e := os.ReadFile(file); e != nil {
		err = e
	} else {
		if len(old) > 0 && len(new) > 0 {
			newdata := bytes.ReplaceAll(data, old, new)
			os.WriteFile(file, newdata, 0644)
		} else {
			err = errors.New(`param empty`)
		}
	}
	return err
}

func (p *Fsconf) Comment(file string, content []byte) error {
	var buffer bytes.Buffer
	buffer.Write([]byte(`<!-- `))
	buffer.Write(content)
	buffer.Write([]byte(` -->`))
	new := buffer.Bytes()
	return p.Update(file, []byte(content), []byte(new))
}

func (p *Fsconf) Uncomment(filename string, content []byte) error {
	var old, new, tmp []byte
	old = content
	if bytes.HasPrefix(old, []byte(`<!-- `)) {
		tmp = bytes.TrimPrefix(old, []byte(`<!-- `))
	} else {
		if bytes.HasPrefix(old, []byte(`<!--`)) {
			tmp = bytes.TrimPrefix(old, []byte(`<!--`))
		} else {
			return errors.New(`Uncomment content not contain "<!-- " or "<!--"`)
		}
	}
	if bytes.HasSuffix(tmp, []byte(` -->`)) {
		new = bytes.TrimSuffix(tmp, []byte(` -->`))
	} else {
		if bytes.HasSuffix(tmp, []byte(`-->`)) {
			new = bytes.TrimSuffix(tmp, []byte(`-->`))
		} else {
			return errors.New(`Uncomment content not contain " -->" or "-->"`)
		}
	}

	return p.Update(filename, []byte(old), []byte(new))
}

func (p *Fsconf) Init() error {
	var err error
	if p != nil {
		if len(p.changedFiles) > 0 {
			err = errors.New(`conf already inited`)
		} else {
			//do init
			err = p.buildBootableConf()
		}
	}
	return err
}

func (p *Fsconf) Reset() error {
	var err error
	if p != nil {
		if len(p.changedFiles) > 0 {
			//do reset
			for _, filechanged := range p.changedFiles {
				defaultFile := fmt.Sprintf("%s%s", filechanged, p.defaultExt)
				if data, errRead := os.ReadFile(defaultFile); errRead == nil {
					os.WriteFile(filechanged, data, 0664)
					os.Remove(defaultFile)
				}
			}
		} else {
			err = errors.New("conf not inited")
		}
	}
	return err
}

func (p *Fsconf) buildBootableConf() error {
	var allerrors string

	//p.Dir()/*.xml
	vars := fmt.Sprintf(`%s/vars.xml`, p.Dir())
	if err := p.buildVars(vars); err != nil {
		allerrors += fmt.Sprintf("vars.xml: %s\n", err.Error())
	}
	return nil
}

func (p *Fsconf) buildVars(vars string) error {
	var err error
	var file = vars
	defaultfile := fmt.Sprintf("%s%s", file, p.defaultExt)
	if _, e := os.Stat(defaultfile); os.IsNotExist(e) {
		if data, e := os.ReadFile(file); e != nil {
			err = e
		} else {
			os.WriteFile(defaultfile, data, 0644)
			//`  <X-PRE-PROCESS cmd="set" data="default_password=1234"/>`
			old := `  <X-PRE-PROCESS cmd="set" data="default_password=1234"/>`
			new := `  <X-PRE-PROCESS cmd="set" data="default_password=D_e_f_a_u_l_t_P_a_s_s_w_o_r_d"/>
  <X-PRE-PROCESS cmd="set" data="pg_handle=pgsql://hostaddr=127.0.0.1 dbname=freeswitch user=fsdba password=fsdba"/>
  <X-PRE-PROCESS cmd="set" data="json_db_handle=$${pg_handle}"/>`
			p.Update(file, []byte(old), []byte(new))
			//`  <X-PRE-PROCESS cmd="stun-set" data="external_sip_ip=stun:stun.freeswitch.org"/>`
			p.Update(file, []byte(`external_sip_ip=stun:stun.freeswitch.org`), []byte(`external_sip_ip=$${local_ip_v4}`))
			//`  <X-PRE-PROCESS cmd="stun-set" data="external_rtp_ip=stun:stun.freeswitch.org"/>`
			p.Update(file, []byte(`external_rtp_ip=stun:stun.freeswitch.org`), []byte(`external_rtp_ip=$${local_ip_v4}`))
		}
	}
	return err
}
