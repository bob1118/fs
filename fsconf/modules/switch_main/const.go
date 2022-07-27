package switch_main

const MOD_NAME = `mod_switch`
const MOD_CONF_NAME = `switch.conf.xml`
const CORE_DB_DSN = `<param name="core-db-dsn" value="$${pg_handle}"/>`
const SWITCH_CONF = `
<configuration name="switch.conf" description="Core Configuration">

  <cli-keybindings>
  </cli-keybindings> 
  
  <default-ptimes>
  </default-ptimes>
  
  <settings>
  </settings>

</configuration>
`
const PRE_LOAD_SWITCH_CONF = `
<configuration name="pre_load_switch.conf" description="Core Configuration">

  <cli-keybindings>
  </cli-keybindings> 
  
  <default-ptimes>
  </default-ptimes>
  
  <settings>
  </settings>

</configuration>
`
const POST_LOAD_SWITCH_CONF = `
<configuration name="post_load_switch.conf" description="Core Configuration">

  <cli-keybindings>
  </cli-keybindings> 
  
  <default-ptimes>
  </default-ptimes>
  
  <settings>
  </settings>

</configuration>
`
