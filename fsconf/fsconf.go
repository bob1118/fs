package fsconf

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Fsconf struct {
	confDir      string   //conf direcotry.
	defaultExt   string   //old conf files named file.defaultExt,eg: vars.xml.default
	changedFiles []string //new conf files who changed by function buildBootableConf()
	v            *viper.Viper
}

func New(conf string, ext string) *Fsconf {
	var dir = conf
	var changedfiles = make([]string, 0, 10)

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
		return &Fsconf{confDir: dir, defaultExt: ext, changedFiles: changedfiles, v: viper.GetViper()}
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

func (p *Fsconf) IsInited() bool {
	return len(p.changedFiles) > 0
}

func (p *Fsconf) Init() error {
	var err error
	if p != nil {
		if p.IsInited() {
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
		if p.IsInited() {
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

// buildBootableConf function.
func (p *Fsconf) buildBootableConf() error {
	var allerrors string

	/////////////////////////////p.dir/*.xml//////////////////////////////////////
	//p.Dir()/vars.xml
	vars := fmt.Sprintf(`%s/vars.xml`, p.Dir())
	if err := p.buildVars(vars); err != nil {
		allerrors += fmt.Sprintf("%s: %s\n", vars, err.Error())
	}

	// ////////////////////////////p.dir/autoload_configs/*.conf.xml/////////////////
	// //p.dir()/autoload_configs/switch.conf.xml
	// autoloadSwitch := fmt.Sprintf(`%s/autoload_configs/switch.conf.xml`, p.Dir())
	// if err := p.buildAutoloadSwitch(autoloadSwitch); err != nil {
	// 	allerrors += fmt.Sprintf("%s: %s\n", autoloadSwitch, err.Error())
	// }
	//p.dir()/autoload_configs/modules.conf.xml
	autoloadModules := fmt.Sprintf(`%s/autoload_configs/modules.conf.xml`, p.Dir())
	if err := p.buildAutoloadModules(autoloadModules); err != nil {
		allerrors += fmt.Sprintf("%s: %s\n", autoloadModules, err.Error())
	}
	//p.dir()/autoload_configs/xml_curl.conf.xml
	autoloadXmlcurl := fmt.Sprintf(`%s/autoload_configs/xml_curl.conf.xml`, p.Dir())
	if err := p.buildAutoloadXmlcurl(autoloadXmlcurl); err != nil {
		allerrors += fmt.Sprintf("%s: %s\n", autoloadXmlcurl, err.Error())
	}

	// /////////////////////////////p.dir/sip_profiles/*.xml//////////////////////////
	// //p.dir()/sip_profiles/internal.xml
	// sipInternal := fmt.Sprintf(`%s/sip_profiles/internal.xml`, p.Dir())
	// if err := p.buildInternal(sipInternal); err != nil {
	// 	allerrors += fmt.Sprintf("%s: %s\n", sipInternal, err.Error())
	// }
	// //p.dir()/sip_profiles/internal-ipv6.xml
	// sipInternalv6 := fmt.Sprintf(`%s/sip_profiles/internal-ipv6.xml`, p.Dir())
	// if err := p.buildInternalv6(sipInternalv6); err != nil {
	// 	allerrors += fmt.Sprintf("%s: %s\n", sipInternalv6, err.Error())
	// }
	// //p.dir()/sip_profiles/external.xml
	// sipExternal := fmt.Sprintf(`%s/sip_profiles/external.xml`, p.Dir())
	// if err := p.buildExternal(sipExternal); err != nil {
	// 	allerrors += fmt.Sprintf("%s: %s\n", sipExternal, err.Error())
	// }
	// //p.dir()/sip_profiles/external-ipv6.xml
	// sipExternalv6 := fmt.Sprintf(`%s/sip_profiles/external-ipv6.xml`, p.Dir())
	// if err := p.buildExternalv6(sipExternalv6); err != nil {
	// 	allerrors += fmt.Sprintf("%s: %s\n", sipExternalv6, err.Error())
	// }

	if len(allerrors) > 0 {
		return errors.New(allerrors)
	}
	return nil
}

func (p *Fsconf) buildVars(in string) error {
	var err error
	var old, new string
	var file = in
	defaultfile := fmt.Sprintf("%s%s", file, p.defaultExt)
	if _, e := os.Stat(defaultfile); os.IsNotExist(e) {
		if data, e := os.ReadFile(file); e != nil {
			err = e
		} else {
			os.WriteFile(defaultfile, data, 0644)
			//`  <X-PRE-PROCESS cmd="set" data="default_password=1234"/>`
			//`  <X-PRE-PROCESS cmd="set" data="pg_handle=pgsql://hostaddr=127.0.0.1 dbname=freeswitch user=fsdba password=fsdba"/>`
			//`  <X-PRE-PROCESS cmd="set" data="json_db_handle=$${pg_handle}"/>`
			old = `  <X-PRE-PROCESS cmd="set" data="default_password=1234"/>`
			new_content := `  <X-PRE-PROCESS cmd="set" data="default_password=D_e_f_a_u_l_t_P_a_s_s_w_o_r_d"/>
  <X-PRE-PROCESS cmd="set" data="pg_handle=%s"/>
  <X-PRE-PROCESS cmd="set" data="json_db_handle=$${pg_handle}"/>`
			pghandle := fmt.Sprintf("pgsql://hostaddr=%s dbname=%s user=%s password=%s", p.v.GetString(`switch.db.host`), p.v.GetString(`switch.db.name`), p.v.GetString(`switch.db.user`), p.v.GetString(`switch.db.password`))
			new = fmt.Sprintf(new_content, pghandle)
			p.Update(file, []byte(old), []byte(new))
			//`  <X-PRE-PROCESS cmd="stun-set" data="external_sip_ip=stun:stun.freeswitch.org"/>`
			old = `  <X-PRE-PROCESS cmd="stun-set" data="external_sip_ip=stun:stun.freeswitch.org"/>`
			external_sip_ip := `  <X-PRE-PROCESS cmd="stun-set" data="external_sip_ip=%s"/>`
			new = fmt.Sprintf(external_sip_ip, p.v.GetString(`switch.vars.external_sip_ip`))
			p.Update(file, []byte(old), []byte(new))
			//`  <X-PRE-PROCESS cmd="stun-set" data="external_rtp_ip=stun:stun.freeswitch.org"/>`
			old = `  <X-PRE-PROCESS cmd="stun-set" data="external_rtp_ip=stun:stun.freeswitch.org"/>`
			external_rtp_ip := `  <X-PRE-PROCESS cmd="stun-set" data="external_rtp_ip=%s"/>`
			new = fmt.Sprintf(external_rtp_ip, p.v.GetString(`switch.vars.external_rtp_ip`))
			p.Update(file, []byte(old), []byte(new))
		}
	}
	return err
}

func (p *Fsconf) buildAutoloadSwitch(in string) error {
	var err error
	var file = in
	defaultfile := fmt.Sprintf("%s%s", file, p.defaultExt)
	if _, e := os.Stat(defaultfile); os.IsNotExist(e) {
		if data, e := os.ReadFile(file); e != nil {
			err = e
		} else {
			os.WriteFile(defaultfile, data, 0644)
			//<!-- <param name="core-db-dsn" value="dsn:username:password" /> -->
			old := `<!-- <param name="core-db-dsn" value="dsn:username:password" /> -->`
			new := `<param name="core-db-dsn" value="$${pg_handle}"/>`
			p.Update(file, []byte(old), []byte(new))
		}
	}
	return err
}

func (p *Fsconf) buildAutoloadModules(in string) error {
	var err error
	var file = in
	defaultfile := fmt.Sprintf("%s%s", file, p.defaultExt)
	if _, e := os.Stat(defaultfile); os.IsNotExist(e) {
		if data, e := os.ReadFile(file); e != nil {
			err = e
		} else {
			os.WriteFile(defaultfile, data, 0644)
			//<load module="mod_enum"/>
			//p.Comment(file, []byte(`<load module="mod_enum"/>`))
			//<!-- <load module="mod_xml_curl"/> -->
			p.Uncomment(file, []byte(`<!-- <load module="mod_xml_curl"/> -->`))
			//<load module="mod_cdr_csv"/>
			modname := p.v.GetString(`switch.cdr.modname`)
			if strings.EqualFold(modname, `mod_odbc_cdr`) {
				p.Update(file, []byte(`<load module="mod_cdr_csv"/>`), []byte(`<load module="mod_odbc_cdr"/>`))
			} else {
				p.Comment(file, []byte(`<load module="mod_cdr_csv"/>`))
			}
			//<load module="mod_loopback"/>
			//p.Comment(file, []byte(`<load module="mod_loopback"/>`))
			//<load module="mod_rtc"/>
			//p.Comment(file, []byte(`<load module="mod_rtc"/>`))
			//<load module="mod_verto"/>
			//p.Comment(file, []byte(`<load module="mod_verto"/>`))
			//<load module="mod_signalwire"/>
			p.Comment(file, []byte(`<load module="mod_signalwire"/>`))
			//<load module="mod_httapi"/>
			//p.Comment(file, []byte(`<load module="mod_httapi"/>`))
			//<load module="mod_dialplan_asterisk"/>
			//p.Comment(file, []byte(`<load module="mod_dialplan_asterisk"/>`))
			//<load module="mod_spandsp"/>
			//p.Comment(file, []byte(`<load module="mod_spandsp"/>`))
			//<load module="mod_b64"/>
			p.Comment(file, []byte(`<load module="mod_b64"/>`))
			//<load module="mod_lua"/>
			//p.Comment(file, []byte(`<load module="mod_say_en"/>`))
			//<load module="mod_say_en"/>
			//p.Comment(file, []byte(`<load module="mod_say_en"/>`))
		}
	}
	return err
}

func (p *Fsconf) buildAutoloadXmlcurl(in string) error {
	var err error
	var file = in
	defaultfile := fmt.Sprintf("%s%s", file, p.defaultExt)
	if _, e := os.Stat(defaultfile); os.IsNotExist(e) {
		if data, e := os.ReadFile(file); e != nil {
			err = e
		} else {
			os.WriteFile(defaultfile, data, 0644)
			//<!-- <param name="gateway-url" value="http://www.freeswitch.org/gateway.xml" bindings="dialplan"/> -->
			//<param name="gateway-url" value="http://localhost/fsapi" bindings="dialplan|configuration|directory|phrases"/>
			old := `<!-- <param name="gateway-url" value="http://www.freeswitch.org/gateway.xml" bindings="dialplan"/> -->`
			new := fmt.Sprintf(`<param name="gateway-url" value="%s" bindings="%s"/>`,
				p.v.GetString(`switch.xml_curl.url`), p.v.GetString(`switch.xml_curl.bindings`))
			p.Update(file, []byte(old), []byte(new))
		}
	}
	return err
}

func (p *Fsconf) buildInternal(in string) error {
	var err error
	var file = in
	defaultfile := fmt.Sprintf("%s%s", file, p.defaultExt)
	if _, e := os.Stat(defaultfile); os.IsNotExist(e) {
		if data, e := os.ReadFile(file); e != nil {
			err = e
		} else {
			os.WriteFile(defaultfile, data, 0644)
			//<!--<param name="odbc-dsn" value="dsn:user:pass"/>-->
			old := `<!--<param name="odbc-dsn" value="dsn:user:pass"/>-->`
			new := `<param name="odbc-dsn" value="$${pg_handle}"/>`
			p.Update(file, []byte(old), []byte(new))
			//<param name="force-register-domain" value="$${domain}"/>
			p.Comment(file, []byte(`<param name="force-register-domain" value="$${domain}"/>`))
			//<param name="force-subscription-domain" value="$${domain}"/>
			p.Comment(file, []byte(`<param name="force-subscription-domain" value="$${domain}"/>`))
			//<param name="force-register-db-domain" value="$${domain}"/>
			p.Comment(file, []byte(`<param name="force-register-db-domain" value="$${domain}"/>`))
		}
	}
	return err
}

func (p *Fsconf) buildInternalv6(in string) error { return p.buildInternal(in) }

func (p *Fsconf) buildExternal(in string) error {
	var err error
	var file = in
	defaultfile := fmt.Sprintf("%s%s", file, p.defaultExt)
	if _, e := os.Stat(defaultfile); os.IsNotExist(e) {
		if data, e := os.ReadFile(file); e != nil {
			err = e
		} else {
			os.WriteFile(defaultfile, data, 0644)
			//<!-- ************************************************* -->
			old := `<!-- ************************************************* -->`
			new := `<param name="odbc-dsn" value="$${pg_handle}"/>`
			p.Update(file, []byte(old), []byte(new))
		}
	}
	return err
}

func (p *Fsconf) buildExternalv6(in string) error { return p.buildExternal(in) }
