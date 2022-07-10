// mod_odbc_cdr
// apt-get install freeswitch-mod-odbc-cdr
// default install without configuration file "odbc_cdr.conf.xml"

// request:
// hostname=D1130&section=configuration&tag_name=configuration&key_name=name&key_value=odbc_cdr.conf
// response:
// <document type="freeswitch/xml">
//   <section name="configuration">
//     <configuration name="odbc_cdr.conf" description="ODBC CDR Configuration">
//       <settings>
//          <!--ADD your parameters here-->
//       </settings>
//     </configuration>
//   </section>
// </document>

package odbc_cdr

import "errors"

func init() {}

func Default() (string, error) { return CONFXML_MOD_ODBC_CDR, nil }

func Read(file string) (string, error) {
	//content, err := os.ReadFile(file)
	//return string(content), err
	return "", errors.New("odbc_cdr.Read nothing")
}

func Build(old string) (string, error) {
	//nothing todo.
	return "", errors.New("content build nothing")
}
