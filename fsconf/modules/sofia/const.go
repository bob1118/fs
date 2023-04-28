package sofia

const MOD_NAME = `mod_sofia`

const MOD_CONF_NAME = `sofia.conf.xml`

const SOFIA_CONF_XML = `
<configuration name="sofia.conf" description="sofia Endpoint">

  <global_settings>
    <param name="log-level" value="0"/>
    <!-- <param name="abort-on-empty-external-ip" value="true"/> -->
    <!-- <param name="auto-restart" value="false"/> -->
    <param name="debug-presence" value="0"/>
    <!-- <param name="capture-server" value="udp:homer.domain.com:5060"/> -->
    
    <!-- 
    	the new format for HEPv2/v3 and capture ID    
	
	protocol:host:port;hep=2;capture_id=200;

    -->
    
    <!-- <param name="capture-server" value="udp:homer.domain.com:5060;hep=3;capture_id=100"/> -->
  </global_settings>

  <!--
      The rabbit hole goes deep.  This includes all the
      profiles in the sip_profiles directory that is up
      one level from this directory.
  -->
  <profiles>
    <X-PRE-PROCESS cmd="include" data="../sip_profiles/*.xml"/>
  </profiles>

</configuration>
`
const SOFIA_CONF_XML_WITH_PROFILE = `
<configuration name="sofia.conf" description="sofia Endpoint">
  <global_settings>
  </global_settings>
  <profiles>
%s
  </profiles>
</configuration>
`
const SOFIA_PROFILE_GATEWAY = `
      <gateway name="%s">
        <param name="username" value="%s"/>
        <param name="realm" value="%s"/>
        <param name="from-user" value="%s"/>
        <param name="from-domain" value="%s"/>
        <param name="password" value="%s"/>
        <param name="extension" value="%s"/>
        <param name="proxy" value="%s"/>
        <param name="register-proxy" value="%s"/>
        <param name="expire-seconds" value="%s"/>
        <param name="register" value="%s"/>
        <!--<param name="register-transport" value="udp"/>-->
        <!--<param name="retry-seconds" value="30"/>-->
        <param name="caller-id-in-from" value="%s"/>
        <!--<param name="contact-params" value=""/>-->
        <param name="extension-in-contact" value="%s"/>
        <param name="ping" value="%s"/>
        <!--<param name="cid-type" value="rpid"/>-->
        <!--<param name="rfc-5626" value="true"/>-->
        <!--<param name="reg-id" value="1"/>-->
      </gateway>
`
