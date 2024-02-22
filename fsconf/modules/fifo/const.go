package fifo

const MOD_NAME = `mod_fifo`
const MOD_CONF_NAME = `fifo.conf.xml`
const FIFO_CONF_XML = `
<configuration name="fifo.conf" description="FIFO Configuration">
  <settings>
    <param name="delete-all-outbound-member-on-startup" value="false"/>
  </settings>
  <fifos>
    <fifo name="cool_fifo@$${domain}" importance="0">
      <!--<member timeout="60" simo="1" lag="20">{member_wait=nowait}user/1005@$${domain}</member>-->
    </fifo>
  </fifos>
</configuration>
`
const OUTBOUND_STRATEGY_RINGALL = `    <param name="outbound-strategy" value="ringall"/>`
const OUTBOUND_STRATEGY_ENTERPRISE = `    <param name="outbound-strategy" value="enterprise"/>`
const ODBC_DSN = `    <param name="odbc-dsn" value="$${pg_handle}"/>`
const DEFAULT_FIFO = `
    <fifo name="cool_fifo@$${domain}" importance="0">
      <!--<member timeout="60" simo="1" lag="20">{member_wait=nowait}user/1005@$${domain}</member>-->
    </fifo>`
const FIFO = `    <fifo name="%s" importance="0"></fifo>`
