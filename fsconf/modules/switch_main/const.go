package switch_main

const MOD_NAME = `mod_switch`
const MOD_CONF_NAME = `switch.conf.xml`
const CORE_DB_DSN = `<param name="core-db-dsn" value="$${pg_handle}"/>`
const MOD_CONF_XML = `
<configuration name="switch.conf" description="Core Configuration">

  <cli-keybindings>
  </cli-keybindings> 
  
  <default-ptimes>
  </default-ptimes>
  
  <settings>
  </settings>

</configuration>
`
const MOD_CONF_SETTINGS = `<settings>`
